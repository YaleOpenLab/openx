package accounts

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
)

type Account struct {
	Seed      string
	PublicKey string
}

// Setup Account is a handler to setup a new account (issuer / investor / school)
func SetupAccount() Account {
	a, err := New()
	if err != nil {
		log.Fatal(err)
	}
	return a
}

func New() (Account, error) {
	var a Account
	pair, err := keypair.Random()
	// so key value pairs over here are ed25519 key pairs instead of bitcoin style key pairs
	// they also seem to sue al lcaps, which I don't know why
	if err != nil {
		return a, err
	}
	log.Println("MY SEED IS: ", pair.Seed())
	a.Seed = pair.Seed()
	a.PublicKey = pair.Address()
	return a, nil
}

func (a *Account) GetCoins() error {
	// get some coins from the stellar robot for testing
	// gives only a constant amount of stellar, so no need to pass it a coin param
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + a.PublicKey)
	if err != nil || resp == nil {
		log.Println("ERRORED OUT while calling friendbot, no coins for us")
		return err
	}
	return nil
}

func (a *Account) GetAssetBalance(assetCode string) (string, error) {

	account, err := horizon.DefaultTestNetClient.LoadAccount(a.PublicKey)
	if err != nil {
		return "", nil
	}

	for _, balance := range account.Balances {
		if balance.Asset.Code == assetCode {
			return balance.Balance, nil
		}
	}

	return "", nil
}

func (a *Account) GetAllBalances() ([]horizon.Balance, error) {

	account, err := horizon.DefaultTestNetClient.LoadAccount(a.PublicKey)
	if err != nil {
		return nil, nil
	}

	return account.Balances, nil
}

func (a *Account) SendCoins(destination string, amount string) (int32, string, error) {

	if _, err := horizon.DefaultTestNetClient.LoadAccount(destination); err != nil {
		// if destination doesn't exist, do nothing
		// returning -11 since -1 maybe returned for unconfirmed tx or something like that
		return -11, "", err
	}

	passphrase := network.TestNetworkPassphrase

	tx, err := build.Transaction(
		build.Network{passphrase},
		build.SourceAccount{a.Seed},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.Payment(
			build.Destination{destination},
			build.NativeAmount{amount},
		),
	)

	if err != nil {
		return -11, "", err
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(a.Seed)
	if err != nil {
		return -11, "", err
	}

	txeB64, err := txe.Base64()
	if err != nil {
		return -11, "", err
	}
	// And finally, send it off to Stellar!
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		return -11, "", err
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)
	return resp.Ledger, resp.Hash, nil
}

func (a *Account) CreateAsset(assetName string) build.Asset {
	// need to set a couple flags here
	return build.CreditAsset(assetName, a.PublicKey)
}

func (a *Account) TrustAsset(asset build.Asset, limit string) error {
	// TRUST is FROM recipient TO issuer
	trustTx, err := build.Transaction(
		build.SourceAccount{a.PublicKey},
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.TestNetwork,
		build.Trust(asset.Code, asset.Issuer, build.Limit(limit)),
	)

	if err != nil {
		return err
	}

	trustTxe, err := trustTx.Sign(a.Seed)
	if err != nil {
		return err
	}

	trustTxeB64, err := trustTxe.Base64()
	if err != nil {
		return err
	}

	tx, err := horizon.DefaultTestNetClient.SubmitTransaction(trustTxeB64)
	if err != nil {
		return err
	}

	log.Println("Trusted asset tx: ", tx.Hash)
	return nil
}

func (a *Account) SendAsset(assetName string, destination string, amount string) (int32, string, error) {
	// this transaction is FROM issuer TO recipient
	paymentTx, err := build.Transaction(
		build.SourceAccount{a.PublicKey},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.Payment(
			build.Destination{AddressOrSeed: destination},
			build.CreditAmount{assetName, a.PublicKey, amount},
		),
	)

	if err != nil {
		return -11, "", err
	}

	paymentTxe, err := paymentTx.Sign(a.Seed)
	if err != nil {
		return -11, "", err
	}

	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {
		return -11, "", err
	}

	tx, err := horizon.DefaultTestNetClient.SubmitTransaction(paymentTxeB64)
	if err != nil {
		return -11, "", err
	}

	log.Println("Sent asset tx is: ", tx.Hash)
	return tx.Ledger, tx.Hash, nil
}
