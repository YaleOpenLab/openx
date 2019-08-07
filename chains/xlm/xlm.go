package xlm

import (
	"github.com/pkg/errors"
	"log"

	utils "github.com/Varunram/essentials/utils"
	horizon "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	horizonprotocol "github.com/stellar/go/protocols/horizon"
	build "github.com/stellar/go/txnbuild"
)

// xlm is a package with stellar related handlers which are useful for interacting with horizon

// Generating a keypair on stellar doesn't mean that you can send funds to it
// you need to call the CreateAccount method in project to be able to send funds
// to it

// GetKeyPair gets a keypair that can be used to interact with the stellar blockchain
func GetKeyPair() (string, string, error) {
	pair, err := keypair.Random()
	return pair.Seed(), pair.Address(), err
}

// AccountExists checks whether an account exists
func AccountExists(publicKey string) bool {
	x, err := ReturnSourceAccountPubkey(publicKey)
	if err != nil {
		// error in the horizon api call
		return false
	}
	return x.Sequence != "0" // if the sequence is zero, the account doesn't exist yet. This equals to the ledger number at which the account was created
}

// SendTx signs and broadcasts a given stellar tx
func SendTx(mykp keypair.KP, tx build.Transaction) (int32, string, error) {
	txe, err := tx.BuildSignEncode(mykp.(*keypair.Full))
	if err != nil {
		return -1, "", errors.Wrap(err, "could not build/sign/encode")
	}

	resp, err := TestNetClient.SubmitTransactionXDR(txe)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not submit tx to horizon")
	}

	log.Printf("Propagated Transaction: %s, sequence: %d\n", resp.Hash, resp.Ledger)
	return resp.Ledger, resp.Hash, nil
}

// SendXLMCreateAccount creates and sends XLM to a new account
func SendXLMCreateAccount(destination string, amountx float64, seed string) (int32, string, error) {
	// don't check if the account exists or not, hopefully it does
	sourceAccount, mykp, err := ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not get source account of seed")
	}

	amount, err := utils.ToString(amountx)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not convert amount to string")
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

// ReturnSourceAccount returns the source account of the seed
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

// ReturnSourceAccountPubkey returns the source account of the pubkey
func ReturnSourceAccountPubkey(pubkey string) (horizonprotocol.Account, error) {
	client := horizon.DefaultTestNetClient
	ar := horizon.AccountRequest{AccountID: pubkey}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return sourceAccount, errors.Wrap(err, "could not load client details, quitting")
	}

	return sourceAccount, nil
}

// SendXLM sends xlm to a destination address
func SendXLM(destination string, amountx float64, seed string, memo string) (int32, string, error) {
	// don't check if the account exists or not, hopefully it does
	sourceAccount, mykp, err := ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not return source account")
	}

	amount, err := utils.ToString(amountx)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not convert amount to string")
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
func RefillAccount(publicKey string, refillSeed string) error {
	if Mainnet {
		return errors.New("can't give free xlm on mainnet, quitting")
	}
	var err error
	if !AccountExists(publicKey) {
		// there is no account under the user's name
		// means we need to setup an account first
		log.Println("Account does not exist, creating: ", publicKey)
		_, _, err = SendXLMCreateAccount(publicKey, RefillAmount, refillSeed)
		if err != nil {
			log.Println("Account Could not be created")
			return errors.Wrap(err, "Account Could not be created")
		}
	}
	// balance is in string, convert to float
	balance, err := GetNativeBalance(publicKey)
	if err != nil {
		return errors.Wrap(err, "could not get native balance")
	}
	balanceI, _ := utils.ToFloat(balance)
	if balanceI < 3 { // to setup trustlines
		_, _, err = SendXLM(publicKey, RefillAmount, refillSeed, "Sending XLM to refill")
		if err != nil {
			return errors.Wrap(err, "Account doesn't have funds or invalid seed")
		}
	}
	return nil
}
