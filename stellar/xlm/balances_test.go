package xlm

import (
	"testing"
)

func TestBalances(t *testing.T) {
	balance, err := GetNativeBalance("GC6Z2KKU4EDTIHAYTJC3Y3AER4ZS5GDSX7S5IKJRRTHLRMJIMCPKQY34")
	if err != nil {
		t.Fatal(err)
	}
	if balance != "9.9999900" {
		t.Fatalf("Balance doesn't match with remote API, quitting!")
	}
	balance, err = GetAssetBalance("GC6Z2KKU4EDTIHAYTJC3Y3AER4ZS5GDSX7S5IKJRRTHLRMJIMCPKQY34", "YOL77fa301ef")
	if err != nil {
		t.Fatal(err)
	}
	if balance != "40000.0000000" {
		t.Fatalf("Balance doesn't match with remote API, quitting!")
	}
	_, err = GetAllBalances("GC6Z2KKU4EDTIHAYTJC3Y3AER4ZS5GDSX7S5IKJRRTHLRMJIMCPKQY34")
	if err != nil {
		t.Fatal(err)
	}
	balance, err = GetUSDTokenBalance("GC6Z2KKU4EDTIHAYTJC3Y3AER4ZS5GDSX7S5IKJRRTHLRMJIMCPKQY34")
	if err == nil {
		// should error out because there is no stableUSD balance on this account
		t.Fatal("There should be no stablecoin balance on this account")
	}
	if HasStableCoin("GC6Z2KKU4EDTIHAYTJC3Y3AER4ZS5GDSX7S5IKJRRTHLRMJIMCPKQY34") {
		// no token balance, should error out
		t.Fatal("Stablecoin present on an address which should have no stablecoin associated with it")
	}
}
