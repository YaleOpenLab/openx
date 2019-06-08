package xlm

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"

	clients "github.com/stellar/go/clients/horizon"
	protocols "github.com/stellar/go/protocols/horizon"
)

// xlm is a set of functions that interface with stellar-core without needing the
// horizon API that stellar provides, which is incomplete

// TestNetClient defines the horizon client to connect to
var TestNetClient = &clients.Client{
	// URL: "http://35.192.122.229:8080",
	URL:  "https://horizon-testnet.stellar.org",
	HTTP: http.DefaultClient,
}

// GetLedgerData gets the latest data from the ledger
func GetLedgerData(blockNumber string) ([]byte, error) {
	var err error
	var data []byte
	resp, err := http.Get(TestNetClient.URL + "/ledgers/" + blockNumber)
	if err != nil || resp.Status != "200 OK" {
		return data, errors.New("API Request did not succeed")
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	return data, err
}

// GetBlockHash gets the block hash corresponding to the passed block number
func GetBlockHash(blockNumber string) (string, error) {
	var err error
	var hash string
	b, err := GetLedgerData(blockNumber)
	if err != nil {
		return hash, errors.Wrap(err, "could not get updated ledger data")
	}
	var x protocols.Ledger
	err = json.Unmarshal(b, &x)
	hash = x.Hash
	log.Printf("The block hash for block %d is: %s and the prev hash is %s", x.Sequence, hash, x.PrevHash)
	return hash, err
}

// GetLatestBlockHash gets the lastest block hash
func GetLatestBlockHash() (string, error) {
	url := TestNetClient.URL + "/ledgers?cursor=now&order=desc&limit=1"
	resp, err := http.Get(url)
	if err != nil || resp.Status != "200 OK" {
		return "", errors.New("API Request did not succeed")
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "could not read response body")
	}
	// hacks below follow because of stellar's incomplete go sdk support
	var x map[string]*json.RawMessage
	err = json.Unmarshal(data, &x)
	if err != nil {
		return "", errors.Wrap(err, "could not unmarshal json")
	}

	var y map[string]*json.RawMessage
	err = json.Unmarshal(*x["_embedded"], &y)
	if err != nil {
		return "", errors.Wrap(err, "could not unmarshal json")
	}

	var z []protocols.Ledger
	err = json.Unmarshal(*y["records"], &z)
	if err != nil {
		return "", errors.Wrap(err, "could not unmarshal json")
	}

	return z[0].Hash, nil
}

// GetAccountData gets the account data
func GetAccountData(a string) ([]byte, error) {
	var err error
	var data []byte
	resp, err := http.Get(TestNetClient.URL + "/accounts/" + a)
	if err != nil {
		return data, errors.Wrap(err, "could not get /accounts/ endpoint from API")
	}
	if resp.Status != "200 OK" {
		return data, errors.New("Request did not succeed")
	}

	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	return data, err
}

// GetNativeBalance gets the xlm balance of a specific account
func GetNativeBalance(publicKey string) (string, error) {
	var balance string
	var err error
	b, err := GetAccountData(publicKey)
	if err != nil {
		// error where account does not exist at all
		// so don't return error and hope its caught later on
		return balance, errors.New("Account does not exist yet, get funds!")
	}
	var x protocols.Account
	err = json.Unmarshal(b, &x)
	if err != nil {
		return balance, errors.Wrap(err, "could not unmarshal data")
	}
	for _, balance := range x.Balances {
		if balance.Asset.Type == "native" {
			return balance.Balance, nil
		}
	}
	// technically accounts on stellar can't exist without a balance, so it should
	// never come here
	return balance, errors.New("Native balance not found")
}

// GetAssetBalance gets the balance of the user in the specific asset
func GetAssetBalance(publicKey string, assetName string) (string, error) {
	var balance string
	var err error
	b, err := GetAccountData(publicKey)
	if err != nil {
		return balance, errors.Wrap(err, "could not get account data")
	}
	var x protocols.Account
	err = json.Unmarshal(b, &x)
	if err != nil {
		return balance, errors.Wrap(err, "could not unmarshal data")
	}
	for _, balance := range x.Balances {
		if balance.Asset.Code == assetName {
			return balance.Balance, nil
		}
	}
	return balance, errors.New("Asset balance not found")
}

// GetAssetTrustLimit gets the trust limit that the user has with the issue for a specific asset
func GetAssetTrustLimit(publicKey string, assetName string) (string, error) {
	var balance string
	var err error
	b, err := GetAccountData(publicKey)
	if err != nil {
		return balance, errors.Wrap(err, "could not get account data")
	}
	var x protocols.Account
	err = json.Unmarshal(b, &x)
	if err != nil {
		return balance, errors.Wrap(err, "could not unmarshal data")
	}
	for _, balance := range x.Balances {
		if balance.Asset.Code == assetName {
			return balance.Limit, nil
		}
	}
	return balance, errors.New("Asset limit not found")
}

// GetAllBalances calls the stellar testnet API to get all the balances associated
// with a certain account.
func GetAllBalances(publicKey string) ([]protocols.Balance, error) {
	account, err := TestNetClient.LoadAccount(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not load account")
	}
	return account.Balances, nil
}

// HasStableCoin checks whether the PublicKey has a stablecoin balance associated
// with it in the first place, if not, returns false
func HasStableCoin(PublicKey string) bool {
	account, err := TestNetClient.LoadAccount(PublicKey)
	if err != nil {
		// account does not exist
		return false
	}

	for _, balance := range account.Balances {
		if balance.Asset.Code == "STABLEUSD" {
			return true
		}
	}
	return false
}

// GetTransactionData gets tx data
func GetTransactionData(txhash string) ([]byte, error) {
	var err error
	var data []byte
	resp, err := http.Get(TestNetClient.URL + "/transactions/" + txhash)
	if err != nil || resp.Status != "200 OK" {
		// check here since if we don't, we need to check the body of the unmarshalled
		// response to see if we have 0
		return data, errors.New("API Request did not succeed")
	}

	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	return data, err
}

// GetTransactionHeight gets tx height
func GetTransactionHeight(txhash string) (int, error) {
	var err error
	var txheight int

	b, err := GetTransactionData(txhash)
	if err != nil {
		return txheight, errors.Wrap(err, "could not get transaction data")
	}
	var x protocols.Transaction
	err = json.Unmarshal(b, &x)
	return int(x.Ledger), err
}