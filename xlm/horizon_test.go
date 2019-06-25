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
	if hash != "3c46ced6f9bf63bc6c2de5f9a5386445ff04340697c61699d91be92da91b9a45" {
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
	if x.ID != "3c46ced6f9bf63bc6c2de5f9a5386445ff04340697c61699d91be92da91b9a45" {
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
	balance, err := GetNativeBalance("GAMTX6MDG65OFU42WGZH2W73AODEKILBW3IWZYPY5SMIVSDSMYXGTTUH")
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetNativeBalance("blah")
	if err == nil {
		t.Fatalf("Account doesn't exist, quitting!")
	}
	_, err = GetAccountData("GAMTX6MDG65OFU42WGZH2W73AODEKILBW3IWZYPY5SMIVSDSMYXGTTUH")
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
	if balance != "5.9996700" {
		log.Println("CHECKBAL:", balance)
		t.Fatalf("Balance doesn't match with remote API, quitting!")
	}
	balance, err = GetAssetBalance("GAMTX6MDG65OFU42WGZH2W73AODEKILBW3IWZYPY5SMIVSDSMYXGTTUH", "OXAc6989b60f")
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetAssetBalance("blah", "OXAc6989b60f")
	if err == nil {
		t.Fatalf("Account doesn't exist, quitting!")
	}
	_, err = GetAssetBalance("GAMTX6MDG65OFU42WGZH2W73AODEKILBW3IWZYPY5SMIVSDSMYXGTTUH", "blah")
	if err == nil {
		t.Fatalf("Asset doesn't exist, quitting!")
	}
	if balance != "120.0000000" {
		t.Fatalf("Balance doesn't match with remote API, quitting!")
	}
	_, err = GetAllBalances("GAMTX6MDG65OFU42WGZH2W73AODEKILBW3IWZYPY5SMIVSDSMYXGTTUH")
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetAllBalances("blah")
	if err == nil {
		t.Fatalf("Account doesn't exist, quitting!")
	}
	if !HasStableCoin("GAMTX6MDG65OFU42WGZH2W73AODEKILBW3IWZYPY5SMIVSDSMYXGTTUH") {
		// no token balance, should error out
		t.Fatal("Stablecoin not present on an address which should have no stablecoin associated with it")
	}
	if HasStableCoin("blah") {
		t.Fatalf("Balance exists on something that ideally shouldn't have balance!")
	}
}

// This test would call the remote / local API and check whether balances match
// if not, this quits immediately
func TestAPIs(t *testing.T) {
	var err error
	height, err := GetTransactionHeight("bea5f00c6327a2d76dbe427c242c5087230191a9c83778b68f3d1fda5a7534a8")
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetTransactionHeight("blah")
	if err == nil {
		t.Fatalf("Shouldn't work, invalid tx hash")
	}
	if height != 2452 {
		t.Fatalf("Heights don't match, quitting!")
	}
	_, err = GetTransactionData("bea5f00c6327a2d76dbe427c242c5087230191a9c83778b68f3d1fda5a7534a8")
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetTransactionData("blah")
	if err == nil {
		t.Fatalf("Shouldn't work, invalid tx hash")
	}
}
