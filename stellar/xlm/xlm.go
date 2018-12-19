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
		return "", err
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
		return "", err
	}

	for _, balance := range account.Balances {
		if balance.Asset.Code == assetCode {
			return balance.Balance, nil
		}
	}

	return "", nil
}

func GetUSDTokenBalance(PublicKey string, targetBalance string) (error) {
	// the USD token defined here is what is issued by the speciifc bank. Ideally, we
	// could accept a tx hash and check it as well, but since we can query balances,
	// much easier to do it this way
	// probably query balance from a couple of servers to avoid the chance for a mitm attack
	// this also makes it asset independednt, means we can require people to hold a
	// specific amount of "X TOKEN" which can be either be a currency like btc / usd / xlm
	// or can be something like a stablecoin or token
	// we also assume that the assetCode of the USDToken is constant and doesn't change.
	return nil
	account, err := utils.DefaultTestNetClient.LoadAccount(PublicKey)
	if err != nil {
		return err
	}

	for _, balance := range account.Balances {
		if balance.Asset.Code == "USDTokenCode here" && balance.Balance == targetBalance {
			return nil
		}
	}
	return fmt.Errorf("Balance insufficient or token not found on your account")
}

// GetAllBalances calls  the stellar testnet API to get all the balances associated
// with a certain account.
func GetAllBalances(PublicKey string) ([]horizon.Balance, error) {

	account, err := utils.DefaultTestNetClient.LoadAccount(PublicKey)
	if err != nil {
		return nil, err
	}

	return account.Balances, nil
}

func DestinationExists(destination string) error {
	_, err := utils.DefaultTestNetClient.LoadAccount(destination)
	return err
}

// Account.SetupAccount() is a method on the structure Account that
// creates a new account using the stellar build.CreateAccount function and
// sends _amount_ number of stellar lumens to the newly created account.
// Note that the destination must alreayd have a keypair generated for this to work
// or else we'd be burning the coins since we wouldn't have the public key
// associated with it,

// SendXLM sends _amount_ number of native tokens (XLM) to the specified destination
// address using the stellar testnet API
func SendXLMCreateAccount(destination string, amount string, Seed string) (int32, string, error) {

	// destination will not exist yet, so don't check
	passphrase := network.TestNetworkPassphrase

	tx, err := build.Transaction(
		build.SourceAccount{Seed},
		build.AutoSequence{utils.DefaultTestNetClient},
		build.Network{passphrase},
		build.CreateAccount(
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
	// And finally, send it off to Stellar
	resp, err := utils.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		return -1, "", err
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)
	return resp.Ledger, resp.Hash, nil
}

// SendXLM sends _amount_ number of native tokens (XLM) to the specified destination
// address using the stellar testnet API
func SendXLM(destination string, amount string, Seed string) (int32, string, error) {

	// don't check if the account exists or not, hopefully it does

	passphrase := network.TestNetworkPassphrase

	tx, err := build.Transaction(
		build.SourceAccount{Seed},
		build.AutoSequence{utils.DefaultTestNetClient},
		build.Network{passphrase},
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
	// And finally, send it off to Stellar
	resp, err := utils.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		log.Println(resp)
		log.Println("3")
		return -1, "", err
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)
	return resp.Ledger, resp.Hash, nil
}
