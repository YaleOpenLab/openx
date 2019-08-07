// +build all travis

package multisig

import (
	"log"
	"testing"

	xlm "github.com/Varunram/essentials/crypto/xlm"
)

var stablecoinPublicKey = "GDPCLB35E4JBVCL2OI6GCM7XK6PLTSKD5EDLRRKFHEI5L4FDKGL4CLIS"

func TestMultisig2of2(t *testing.T) {
	xlm.SetConsts(10, false)
	seed1, pubkey1, err := xlm.GetKeyPair()
	if err != nil {
		log.Println(err)
	}

	seed2, pubkey2, err := xlm.GetKeyPair()
	if err != nil {
		log.Println(err)
	}

	// setup both accounts
	err = xlm.GetXLM(pubkey1)
	if err != nil {
		t.Fatal(err)
	}

	err = xlm.GetXLM(pubkey2)
	if err != nil {
		t.Fatal(err)
	}

	err = Convert2of2(pubkey1, seed1, pubkey2)
	if err != nil {
		t.Fatal(err)
	}

	// now account 1 requires the signature of both seed1 and seed2 in order to be able to send a tx.
	// we need to check that.
	destination := pubkey1
	err = Tx2of2(pubkey1, destination, seed1, seed2, 1, "test")
	if err != nil {
		t.Fatal(err)
	}

	err = AuthImmutable2of2(pubkey1, seed1, seed2)
	if err != nil {
		t.Fatal(err)
	}

	err = TrustAssetTx("STABLEUSD", stablecoinPublicKey, "10000", pubkey1, seed1, seed2)
	if err != nil {
		t.Fatal(err)
	}

	err = Tx2of2(pubkey1, pubkey1, seed1, seed2, 1, "cool")
	if err != nil {
		t.Fatal(err)
	}
}

// we're forced to hav separate tests because we can't use the same tests (they'll eb converted to multisig accounts)
func TestNew2of2MultiSig(t *testing.T) {
	seed1, pubkey1, err := xlm.GetKeyPair()
	if err != nil {
		log.Println(err)
	}

	seed2, pubkey2, err := xlm.GetKeyPair()
	if err != nil {
		log.Println(err)
	}

	// setup both accounts
	err = xlm.GetXLM(pubkey1)
	if err != nil {
		t.Fatal(err)
	}

	err = xlm.GetXLM(pubkey2)
	if err != nil {
		t.Fatal(err)
	}

	pubkey, err := New2of2(pubkey1, pubkey2)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("New 2of2 Multisig Pubkey: ", pubkey)

	destination := pubkey
	err = Tx2of2(pubkey, destination, seed1, seed2, 1, "test")
	if err != nil {
		t.Fatal(err)
	}
}

// now test 1of2 multisig
func TestNew1of2MultiSig(t *testing.T) {
	seed1, pubkey1, err := xlm.GetKeyPair()
	if err != nil {
		log.Println(err)
	}

	seed2, pubkey2, err := xlm.GetKeyPair()
	if err != nil {
		log.Println(err)
	}

	// setup both accounts
	err = xlm.GetXLM(pubkey1)
	if err != nil {
		t.Fatal(err)
	}

	err = xlm.GetXLM(pubkey2)
	if err != nil {
		t.Fatal(err)
	}

	pubkey, err := New1of2(pubkey1, pubkey2)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("New 1of2 Multisig Pubkey: ", pubkey)
	destination := pubkey
	amount := 1.0
	memo := "seed1test"
	// we now have a one of 2 multisig, this means we can send a tx using one of the 2 seeds generated above
	_, _, err = xlm.SendXLM(destination, amount, seed1, memo)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = xlm.SendXLM(destination, amount, seed2, memo)
	if err != nil {
		t.Fatal(err)
	}
}

func Test1ofxMultiSig(t *testing.T) {
	seed1, pubkey1, err := xlm.GetKeyPair()
	if err != nil {
		log.Println(err)
	}

	seed2, pubkey2, err := xlm.GetKeyPair()
	if err != nil {
		log.Println(err)
	}

	// setup both accounts
	err = xlm.GetXLM(pubkey1)
	if err != nil {
		t.Fatal(err)
	}

	err = xlm.GetXLM(pubkey2)
	if err != nil {
		t.Fatal(err)
	}

	pubkey, err := Newxofy(1, 2, pubkey1, pubkey2)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Newxofy(1, 3, pubkey1, pubkey2)
	if err == nil {
		t.Fatalf("not able to catch number of signers error, quitting")
	}

	log.Println("New 1of2 Multisig Pubkey: ", pubkey)
	destination := pubkey
	amount := 1.0
	memo := "seed1test"
	// we now have a one of 2 multisig, this means we can send a tx using one of the 2 seeds generated above
	_, _, err = xlm.SendXLM(destination, amount, seed1, memo)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = xlm.SendXLM(destination, amount, seed2, memo)
	if err != nil {
		t.Fatal(err)
	}
}
