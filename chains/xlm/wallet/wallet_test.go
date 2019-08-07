// +build all travis

package wallet

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestWallet(t *testing.T) {
	// test out stuff here
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filepath := pwd + "/files/testseed.hex"
	filepathdup := pwd + "/files/testseed2.hex"
	// testseed1 already has a seed, we need to test if decrypted seeds and stored seeds match
	// seed: SCLRCFCYJZ7YPQZPNZVR7Q2NXUJBBSHZJAZQZW2UGI57RFP3KA26MGGE
	// pubkey: GDH5HXIPZL435MJZGXZ4DNZFJEC7CFB6EJ44AC4VCGRAZRDS5T7PS3ZN
	pubkey, seed, err := RetrieveSeed(filepath, "password")
	if err != nil {
		t.Fatal(err)
	}
	if seed != "SCLRCFCYJZ7YPQZPNZVR7Q2NXUJBBSHZJAZQZW2UGI57RFP3KA26MGGE" {
		t.Fatalf("Seed doesn't match with test seed")
	}
	if pubkey != "GDH5HXIPZL435MJZGXZ4DNZFJEC7CFB6EJ44AC4VCGRAZRDS5T7PS3ZN" {
		t.Fatalf("Publickey doesn't match with test publickey")
	}
	_, _, err = RetrieveSeed(filepath, "wrongpassword")
	if err == nil {
		t.Fatalf("Decrpytion succeeds with incorrect password")
	}
	_, _, err = RetrieveSeed("blah", "password")
	if err == nil {
		t.Fatalf("Decrpytion succeeds with incorrect path")
	}
	os.Remove(filepathdup)
	// need to read from filepath and decrypt seed to test that function
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}
	seedcheck, err := DecryptSeed(data, "password")
	if err != nil {
		t.Fatal(err)
	}
	if seedcheck != seed {
		t.Fatalf("Seeds doesn't match with stored seed, quitting!")
	}
	_, err = DecryptSeed(data, "wrongpassword")
	if err == nil {
		t.Fatalf("Can decrypt the seed with the wrong password, qutiting!")
	}
	_, seedtest, err := NewSeedStore(filepathdup, "password")
	if err != nil {
		t.Fatal(err)
	}
	err = StoreSeed(seedtest, "password", filepathdup)
	if err != nil {
		t.Fatal(err)
	}
	err = StoreSeed(seedtest, "blah", filepathdup)
	if err != nil {
		t.Fatal(err)
	}
	err = StoreSeed(seedtest, "blah", "./blah/blah.txt")
	if err == nil {
		t.Fatalf("Not able to catch invalid path error")
	}
	rtpk, err := ReturnPubkey(seed)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ReturnPubkey("")
	if err == nil {
		t.Fatalf("Not able to catch empty seed error")
	}
	if rtpk != pubkey {
		t.Fatalf("Pubkeys don't match, quitting!")
	}
	os.Remove(filepathdup)
}
