// +build all travis

package xlm

import (
	"testing"
	"log"
)

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
	oldTc := TestNetClient.URL
	TestNetClient.URL = "blah"
	_, err = GetAccountData("blah")
	if err == nil {
		t.Fatalf("Can return data with invalid url, quitting!")
	}
	TestNetClient.URL = oldTc
	if balance != "5.9999700" {
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
