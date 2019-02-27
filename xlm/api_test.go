// +build all travis

package xlm

import (
	"encoding/json"
	"log"
	"testing"

	protocols "github.com/stellar/go/protocols/horizon"
)

func TestApi(t *testing.T) {
	// test out stuff here
	blockNumber := "2"
	hash, err := GetBlockHash(blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	if hash != "1685117d8d2270aed3cf81b641087b4f8d8e2b0774b4bb82d1de34d9472fb3d5" {
		t.Fatalf("Hashes don't match, quitting!")
	}
	log.Println(hash)
	_, err = GetBlockHash("-1")
	if err == nil {
		t.Fatalf("Can get data for negative block number, quitting!")
	}
	data, err := GetLedgerData(blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetLedgerData("-1")
	if err == nil {
		t.Fatal("Can get data for negative block number, quitting!")
	}
	var x protocols.Ledger
	err = json.Unmarshal(data, &x)
	if err != nil {
		t.Fatal(err)
	}
	if x.ID != "1685117d8d2270aed3cf81b641087b4f8d8e2b0774b4bb82d1de34d9472fb3d5" {
		t.Fatal(err)
	}
	_, err = GetLatestBlockHash()
	if err != nil {
		t.Fatal(err)
	}
	oldTc := TestNetClient.URL
	TestNetClient.URL = "blah"
	_, err = GetLatestBlockHash()
	if err == nil {
		t.Fatalf("can call with invalid client URL")
	}
	TestNetClient.URL = oldTc
}
