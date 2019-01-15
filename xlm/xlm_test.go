// +build all travis

package xlm

import (
	"testing"
)

func TestXLM(t *testing.T) {
	seed, address, err := GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	err = GetXLM(address)
	if err != nil {
		t.Fatal(err)
	}
	err = GetXLM("blah")
	if err == nil {
		t.Fatal("Invalid Address, shouldn't work!")
	}
	// now we have coins, so we can send to another account
	// need to create an account through
	_, destPubKey, err := GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLMCreateAccount(destPubKey, "2", "wrong seed")
	// create the destiantion account by sendignb some coins to bootstrap
	if err == nil {
		t.Fatalf("Wrong seed, shouldn't work!")
	}
	_, _, err = SendXLMCreateAccount("wrong pubkey", "2", seed)
	// create the destiantion account by sendignb some coins to bootstrap
	if err == nil {
		t.Fatalf("Wrong pubkey, shouldn't work!")
	}
	_, _, err = SendXLMCreateAccount(destPubKey, "2", seed)
	// create the destiantion account by sendignb some coins to bootstrap
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLM(destPubKey, "1", seed, "")
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLM(destPubKey, "1", "wrong seed", "")
	if err == nil {
		t.Fatalf("Wrong seed, shouldn't work!")
	}
	_, _, err = SendXLM("wrong pubkey", "1", seed, "")
	if err == nil {
		t.Fatalf("Wrong pubkey, shouldn't work!")
	}
	if AccountExists("blah") {
		t.Fatalf("Dummy account exists, shouldn't!")
	}
	if !AccountExists(destPubKey) {
		t.Fatalf("Account which should exist doesn't, quitting!")
	}
	_, testPubKey, err := GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	err = RefillAccount(testPubKey, seed)
	if err != nil {
		t.Fatal(err)
	}
	err = RefillAccount("blah", seed)
	if err == nil {
		t.Fatal("Not catching wrong pubkey error, quitting!")
	}
	// don't test the reverse becuase apparently there's some problem in catching the next block
	// or something
	testseed, pk, err :=  GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	err = GetXLM(pk)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLM(destPubKey, "9998", testseed, "")
	if err != nil {
		t.Fatal(err)
	}
	err = RefillAccount(pk, seed)
	if err != nil {
		t.Fatal(err)
	}
}
