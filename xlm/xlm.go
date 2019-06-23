package xlm

import (
	"github.com/pkg/errors"
	"log"

	consts "github.com/YaleOpenLab/openx/consts"
	utils "github.com/YaleOpenLab/openx/utils"
	horizon "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	horizonprotocol "github.com/stellar/go/protocols/horizon"
	build "github.com/stellar/go/txnbuild"
)

// package xlm provides all the necessary handlers in order to interact with the
// Stellar blockchain

// Generating a keypair on stellar doesn't mean that you can send funds to it
// you need to call the CreateAccount method in project to be able to send funds
// to it

// GetKeyPair gets a keypair that can be used to interact with the stellar blockchain
func GetKeyPair() (string, string, error) {
	pair, err := keypair.Random()
	return pair.Seed(), pair.Address(), err
}

// AccountExists checks whether an accoutn exists, not needed now since we do the check ourselves
// in multiple places
func AccountExists(address string) bool {
	_, err := TestNetClient.LoadAccount(address)
	return !(err != nil)
	/*
		if err != nil {
			return false
		}
		return true
	*/
}

func SendTx(mykp keypair.KP, tx build.Transaction) (int32, string, error) {
	txe, err := tx.BuildSignEncode(mykp.(*keypair.Full))
	if err != nil {
		return -1, "", err
	}

	resp, err := TestNetClient.SubmitTransaction(txe)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not submit tx to horizon")
	}

	log.Printf("Propagated Transaction: %s, sequence: %d\n", resp.Hash, resp.Ledger)

	return resp.Ledger, resp.Hash, nil
}

// SendXLMCreateAccount sends XLM to an account and creates the account if it doesn't exist already
func SendXLMCreateAccount(destination string, amount string, seed string) (int32, string, error) {

	// don't check if the account exists or not, hopefully it does
	sourceAccount, mykp, err := ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", err
	}

	op := build.CreateAccount{
		Destination: destination,
		Amount:      amount,
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       Passphrase,
	}

	return SendTx(mykp, tx)
}

func ReturnSourceAccount(seed string) (horizonprotocol.Account, keypair.KP, error) {
	var sourceAccount horizonprotocol.Account
	mykp, err := keypair.Parse(seed)
	if err != nil {
		return sourceAccount, mykp, errors.Wrap(err, "could not parse keypair, quitting")
	}

	client := horizon.DefaultTestNetClient
	ar := horizon.AccountRequest{AccountID: mykp.Address()}
	sourceAccount, err = client.AccountDetail(ar)
	if err != nil {
		return sourceAccount, mykp, errors.Wrap(err, "could not load client details, quitting")
	}

	return sourceAccount, mykp, nil
}

func ReturnSourceAccountPubkey(pubkey string) (horizonprotocol.Account, error) {
	client := horizon.DefaultTestNetClient
	ar := horizon.AccountRequest{AccountID: pubkey}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return sourceAccount, errors.Wrap(err, "could not load client details, quitting")
	}

	return sourceAccount, nil
}

// SendXLM sends _amount_ number of native tokens (XLM) to the specified destination
// address using the stellar testnet API
func SendXLM(destination string, amount string, seed string, memo string) (int32, string, error) {
	// don't check if the account exists or not, hopefully it does
	sourceAccount, mykp, err := ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", err
	}

	op := build.Payment{
		Destination: destination,
		Amount:      amount,
		Asset:       build.NativeAsset{},
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       Passphrase,
		Memo:          build.Memo(build.MemoText(memo)),
	}

	return SendTx(mykp, tx)
}

// RefillAccount refills an account
func RefillAccount(publicKey string, platformSeed string) error {
	var err error
	if !AccountExists(publicKey) {
		// there is no account under the user's name
		// means we need to setup an account first
		log.Println("Account does not exist, creating: ", publicKey)
		_, _, err = SendXLMCreateAccount(publicKey, consts.RefillAmount, platformSeed)
		if err != nil {
			log.Println("Account Could not be created")
			return errors.New("Account Could not be created")
		}
	}
	// balance is in string, convert to float
	balance, err := GetNativeBalance(publicKey)
	if err != nil {
		return errors.Wrap(err, "could not get native balance")
	}
	balanceI := utils.StoF(balance)
	if balanceI < 3 { // to setup trustlines
		_, _, err = SendXLM(publicKey, consts.RefillAmount, platformSeed, "Sending XLM to refill")
		if err != nil {
			return errors.New("Account doesn't have funds or invalid seed")
		}
	}
	return nil
}
