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
	blockNumber := "1000000"
	hash, err := GetBlockHash(blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	if hash != "8267f8beb7c461bbc00af5e153d0d28fdf332b2b21c421d8423a59b995c958b1" {
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
	if x.ID != "8267f8beb7c461bbc00af5e153d0d28fdf332b2b21c421d8423a59b995c958b1" {
		t.Fatal(err)
	}
	//t.Fatal(hash)
}
