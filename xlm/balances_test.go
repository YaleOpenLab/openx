// +build all travis

package xlm

import (
	"testing"
)

func TestBalances(t *testing.T) {
	balance, err := GetNativeBalance("GDPCBDVZGJ3WXL2B7YUSTRYPNZ464MCZDAHA3ZTPC36KULDYAMPNX423")
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetNativeBalance("blah")
	if err == nil {
		t.Fatalf("Account doesn't exist, quitting!")
	}
	_, err = GetAccountData("GDPCBDVZGJ3WXL2B7YUSTRYPNZ464MCZDAHA3ZTPC36KULDYAMPNX423")
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
	if balance != "9998.9999600" {
		t.Fatalf("Balance doesn't match with remote API, quitting!")
	}
	balance, err = GetAssetBalance("GDPCBDVZGJ3WXL2B7YUSTRYPNZ464MCZDAHA3ZTPC36KULDYAMPNX423", "YOL9e0b5fa3d")
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetAssetBalance("blah", "YOL9e0b5fa3d")
	if err == nil {
		t.Fatalf("Account doesn't exist, quitting!")
	}
	_, err = GetAssetBalance("GDPCBDVZGJ3WXL2B7YUSTRYPNZ464MCZDAHA3ZTPC36KULDYAMPNX423", "blah")
	if err == nil {
		t.Fatalf("Asset doesn't exist, quitting!")
	}
	if balance != "1000.0000000" {
		t.Fatalf("Balance doesn't match with remote API, quitting!")
	}
	_, err = GetAllBalances("GDPCBDVZGJ3WXL2B7YUSTRYPNZ464MCZDAHA3ZTPC36KULDYAMPNX423")
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetAllBalances("blah")
	if err == nil {
		t.Fatalf("Account doesn't exist, quitting!")
	}
	if !HasStableCoin("GDPCBDVZGJ3WXL2B7YUSTRYPNZ464MCZDAHA3ZTPC36KULDYAMPNX423") {
		// no token balance, should error out
		t.Fatal("Stablecoin not present on an address which should have no stablecoin associated with it")
	}
	if HasStableCoin("blah") {
		t.Fatalf("Balance exists on something that ideally shouldn't have balance!")
	}
}
