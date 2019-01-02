package xlm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/stellar/go/protocols/horizon" // using this since hte client/horizon package has some deprecated fields
	protocols "github.com/stellar/go/protocols/horizon"
)

/*
type Account struct {
    Links struct {
        Self         hal.Link `json:"self"`
        Transactions hal.Link `json:"transactions"`
        Operations   hal.Link `json:"operations"`
        Payments     hal.Link `json:"payments"`
        Effects      hal.Link `json:"effects"`
        Offers       hal.Link `json:"offers"`
        Trades       hal.Link `json:"trades"`
        Data         hal.Link `json:"data"`
    }   `json:"_links"`

    HistoryAccount
    Sequence             string            `json:"sequence"`
    SubentryCount        int32             `json:"subentry_count"`
    InflationDestination string            `json:"inflation_destination,omitempty"`
    HomeDomain           string            `json:"home_domain,omitempty"`
    Thresholds           AccountThresholds `json:"thresholds"`
    Flags                AccountFlags      `json:"flags"`
    Balances             []Balance         `json:"balances"`
    Signers              []Signer          `json:"signers"`
    Data                 map[string]string `json:"data"`
}
*/
func GetAccountData(a string) ([]byte, error) {
	var err error
	var data []byte
	resp, err := http.Get(TestNetClient.URL + "/accounts/" + a)
	if err != nil {
		return data, err
	}
	log.Println("StATUS: ", resp.Status)
	if resp.Status != "200 OK" {
		return data, fmt.Errorf("Request did not succeed")
	}

	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	return data, err
}

func GetNativeBalance(publicKey string) (string, error) {
	var balance string
	var err error
	b, err := GetAccountData(publicKey)
	if err != nil {
		return balance, err
	}
	var x protocols.Account
	err = json.Unmarshal(b, &x)
	if err != nil {
		return balance, err
	}
	for _, balance := range x.Balances {
		if balance.Asset.Type == "native" {
			log.Println("Native balance: ", balance.Balance)
			return balance.Balance, nil
		}
	}
	return balance, fmt.Errorf("Native balance not found")
}

func GetAssetBalance(publicKey string, assetName string) (string, error) {
	var balance string
	var err error
	b, err := GetAccountData(publicKey)
	if err != nil {
		return balance, err
	}
	var x protocols.Account
	err = json.Unmarshal(b, &x)
	if err != nil {
		return balance, err
	}
	for _, balance := range x.Balances {
		if balance.Asset.Type == assetName {
			log.Println("Native balance: ", balance.Balance)
			return balance.Balance, nil
		}
	}
	return balance, fmt.Errorf("Native balance not found")
}

// GetAllBalances calls the stellar testnet API to get all the balances associated
// with a certain account.
func GetAllBalances(publicKey string) ([]horizon.Balance, error) {

	account, err := TestNetClient.LoadAccount(publicKey)
	if err != nil {
		return nil, err
	}

	return account.Balances, nil
}

func GetUSDTokenBalance(PublicKey string, targetBalance string) error {
	// the USD token defined here is what is issued by the speciifc bank. Ideally, we
	// could accept a tx hash and check it as well, but since we can query balances,
	// much easier to do it this way
	// probably query balance from a couple of servers to avoid the chance for a mitm attack
	// this also makes it asset independednt, means we can require people to hold a
	// specific amount of "X TOKEN" which can be either be a currency like btc / usd / xlm
	// or can be something like a stablecoin or token
	// we also assume that the assetCode of the USDToken is constant and doesn't change.
	account, err := TestNetClient.LoadAccount(PublicKey)
	if err != nil {
		return err
	}

	for _, balance := range account.Balances {
		if balance.Asset.Code == "STABLEUSD" && balance.Balance >= targetBalance {
			return nil
		}
	}
	return fmt.Errorf("Balance insufficient or STABLEUSD token not found on your account")
}
