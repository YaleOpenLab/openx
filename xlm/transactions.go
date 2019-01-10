package xlm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	protocols "github.com/stellar/go/protocols/horizon"
)

func GetTransactionData(txhash string) ([]byte, error) {
	var err error
	var data []byte
	resp, err := http.Get(TestNetClient.URL + "/transactions/" + txhash)
	if err != nil {
		return data, err
	}

	if resp.Status != "200 OK" {
		// check here since if we don't, we need to check the body of the unmarshalled
		// response to see if we have 0
		return data, fmt.Errorf("API Request did not succeed")
	}

	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	return data, err
}

func GetTransactionHeight(txhash string) (int, error) {
	var err error
	var txheight int

	b, err := GetTransactionData(txhash)
	if err != nil {
		return txheight, err
	}
	var x protocols.Transaction
	err = json.Unmarshal(b, &x)
	if err != nil {
		return txheight, err
	}
	log.Printf("Tx height of %s is %d", txhash, x.Ledger)
	return int(x.Ledger), nil
}
