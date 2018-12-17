// the xlm package is a package that interacts with the stellar testnet
// API and fetches testnet coins for the user
package xlm

import (
	"fmt"
	"log"
	"net/http"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon" // using this since hte client/horizon package has some deprecated fields
)

func GetKeyPair() (string, string, error) {
	pair, err := keypair.Random()
	return pair.Seed(), pair.Address(), err
}

// GetCoins makes an API call to the friendbot on stellar testnet, which gives
// us 10000 XLM for use. We don't need 10000XLM (we need only ~3 XLM for setting up
// various trustlines), but there's no option to receive less, so we're having to call
// this. On mainnet, we'd be refilling the accoutns manually, so this function
// wouldn't exist.
func GetXLM(PublicKey string) error {
	// get some coins from the stellar robot for testing
	// gives only a constant amount of stellar, so no need to pass it a coin param
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + PublicKey)
	if err != nil || resp == nil {
		log.Println("ERRORED OUT while calling friendbot, no coins for us")
		return err
	}
	return nil
}

func GetXLMBalance(PublicKey string) (string, error) {

	account, err := utils.DefaultTestNetClient.LoadAccount(PublicKey)
	if err != nil {
		return "", nil
	}

	for _, balance := range account.Balances {
		if balance.Asset.Type == "native" {
			return balance.Balance, nil
		}
	}

	return "", fmt.Errorf("Couldn't find native asset balance")
}

// GetAssetBalance calls the stellar testnet API to get all balances
// and then runs through the balances to get the balance of a specific account
func GetAssetBalance(PublicKey string, assetCode string) (string, error) {

	account, err := utils.DefaultTestNetClient.LoadAccount(PublicKey)
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

// GetAllBalances calls  the stellar testnet API to get all the balances associated
// with a certain account.
func GetAllBalances(PublicKey string) ([]horizon.Balance, error) {

	account, err := utils.DefaultTestNetClient.LoadAccount(PublicKey)
	if err != nil {
		return nil, nil
	}

	return account.Balances, nil
}

// Account.SetupAccount() is a method on the structure Account that
// creates a new account using the stellar build.CreateAccount function and
// sends _amount_ number of stellar lumens to the newly created account.
// Note that the destination must alreayd have a keypair generated for this to work
// or else we'd be burning the coins since we wouldn't have the public key
// associated with it,

// SendXLM sends _amount_ number of native tokens (XLM) to the specified destination
// address using the stellar testnet API
func SendXLM(destination string, amount string, Seed string) (int32, string, error) {

	if _, err := utils.DefaultTestNetClient.LoadAccount(destination); err != nil {
		// if destination doesn't exist, do nothing
		// returning -1 since -1 maybe returned for unconfirmed tx or something like that
		return -1, "", err
	}

	passphrase := network.TestNetworkPassphrase

	tx, err := build.Transaction(
		build.Network{passphrase},
		build.SourceAccount{Seed},
		build.AutoSequence{utils.DefaultTestNetClient},
		build.Payment(
			build.Destination{destination},
			build.NativeAmount{amount},
		),
	)

	if err != nil {
		return -1, "", err
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(Seed)
	if err != nil {
		return -1, "", err
	}

	txeB64, err := txe.Base64()
	if err != nil {
		return -1, "", err
	}
	// And finally, send it off to Stellar!
	resp, err := utils.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		return -1, "", err
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)
	return resp.Ledger, resp.Hash, nil
}
