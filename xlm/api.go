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
