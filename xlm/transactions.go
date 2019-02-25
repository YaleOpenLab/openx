package xlm

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"

	protocols "github.com/stellar/go/protocols/horizon"
)

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
