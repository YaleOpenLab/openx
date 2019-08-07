// +build all travis

package xlm

import (
	"encoding/json"
	"log"
	"testing"

	protocols "github.com/stellar/go/protocols/horizon"
)

var blockNumber = 47730
var blockHash = "e0b2b2cff90312b60e55365914a9e8d550ed05aba50bf926312577d71e08546c"
var txhash = "7532311a4816a4d61eccb6087704880108b57e27628299e0129267f197d3c5f1"
var txHashNumber = 48046

func TestApi(t *testing.T) {
	SetConsts(10, false)
	// test out stuff here
	hash, err := GetBlockHash(blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	if hash != blockHash {
		t.Fatalf("Hashes don't match, quitting!")
	}
	log.Println(hash)
	_, err = GetBlockHash(-1)
	if err == nil {
		t.Fatalf("Can get data for negative block number, quitting!")
	}
	data, err := GetLedgerData(blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetLedgerData(-1)
	if err == nil {
		t.Fatal("Can get data for negative block number, quitting!")
	}
	var x protocols.Ledger
	err = json.Unmarshal(data, &x)
	if err != nil {
		t.Fatal(err)
	}
	if x.ID != blockHash {
		t.Fatal(err)
	}
	_, err = GetLatestBlockHash()
	if err != nil {
		t.Fatal(err)
	}
	oldTc := TestNetClient.HorizonURL
	TestNetClient.HorizonURL = "blah"
	_, err = GetLatestBlockHash()
	if err == nil {
		t.Fatalf("can call with invalid client URL")
	}
	TestNetClient.HorizonURL = oldTc
}

func TestBalances(t *testing.T) {
	account := "GDRGIFNANGDI5CIWAAVA6H2DZJTQQZQMXVRI42NNDQSQCJI575QS4NHE"

	balance, err := GetNativeBalance(account)
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetNativeBalance("blah")
	if err == nil {
		t.Fatalf("Account doesn't exist, quitting!")
	}
	_, err = GetAccountData(account)
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetAccountData("blah")
	if err == nil {
		t.Fatalf("Invalid account exists, quitting!")
	}
	oldTc := TestNetClient.HorizonURL
	TestNetClient.HorizonURL = "blah"
	_, err = GetAccountData("blah")
	if err == nil {
		t.Fatalf("Can return data with invalid url, quitting!")
	}
	TestNetClient.HorizonURL = oldTc
	if balance != 10000 {
		log.Println("CHECKBAL:", balance)
		t.Fatalf("Balance doesn't match with remote API, quitting!")
	}
	balance, err = GetAssetBalance(account, "STABLEUSD")
	if err != nil {
		t.Fatal(err)
	}
	if balance != 10.0000000 {
		log.Println("CHECKBALASSET: ", balance)
		t.Fatalf("Balance doesn't match with remote API, quitting!")
	}
	_, err = GetAssetBalance("blah", "STABLEUSD")
	if err == nil {
		t.Fatalf("Account doesn't exist, quitting!")
	}
	_, err = GetAssetBalance(account, "blah")
	if err == nil {
		t.Fatalf("Asset doesn't exist, quitting!")
	}
	_, err = GetAllBalances(account)
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetAllBalances("blah")
	if err == nil {
		t.Fatalf("Account doesn't exist, quitting!")
	}
	if !HasStableCoin(account) {
		// no token balance, should error out
		t.Fatal("Stablecoin not present on an address which should have stablecoin associated with it")
	}
	if HasStableCoin("blah") {
		t.Fatalf("Balance exists on something that ideally shouldn't have balance!")
	}
}

// This test would call the remote / local API and check whether balances match
// if not, this quits immediately
func TestAPIs(t *testing.T) {
	var err error
	height, err := GetTransactionHeight(txhash)
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetTransactionHeight("blah")
	if err == nil {
		t.Fatalf("Shouldn't work, invalid tx hash")
	}
	if height != txHashNumber {
		t.Fatalf("Heights don't match, quitting!")
	}
	_, err = GetTransactionData(txhash)
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetTransactionData("blah")
	if err == nil {
		t.Fatalf("Shouldn't work, invalid tx hash")
	}
}
