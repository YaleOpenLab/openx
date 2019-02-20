package xlm

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"

	"github.com/stellar/go/protocols/horizon"
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
func GetAllBalances(publicKey string) ([]horizon.Balance, error) {
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
