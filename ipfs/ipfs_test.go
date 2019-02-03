// +build all

package ipfs

import (
	"log"
	"testing"
)

// You need to be running an active ipfs daemon in order for this to succeed
// kind of weird to get ipfs setup for a single tests, but this would be useful
func TestIpfs(t *testing.T) {
	hash, err := AddStringToIpfs("Hello, this is a test from ipfs to see if it works")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("HASH: ", hash)
	string1, err := GetStringFromIpfs(hash)
	if err != nil {
		t.Fatal(err)
	}
	if string1 == "Hellox, this is a test from ipfs to see if it works" {
		t.Fatal("DEcrypted string does not match with original, exiting!")
	}
	_, err = GetStringFromIpfs("blah")
	if err == nil {
		t.Fatal("Can retrieve non existing hash, quitting!")
	}
	err = GetFileFromIpfs("/ipfs/"+hash, "pdf")
	if err != nil {
		t.Fatal(err)
	}
	err = GetFileFromIpfs("blah", "pdf")
	if err == nil {
		t.Fatalf("Can retrieve non existing hash, quitting!")
	}
	_, err = ReadfromFile("files/test.pdf") // get the data from the pdf as a datastream
	if err != nil {
		t.Fatal(err)
	}
	_, err = ReadfromFile("blah") // get the data from the pdf as a datastream
	if err == nil {
		t.Fatal("Can read from non existing pdf.")
	}
	hash, err = IpfsHashFile("files/test.pdf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = IpfsHashFile("blah")
	if err == nil {
		t.Fatal("Can retrieve non existing pdf file")
	}
	err = GetFileFromIpfs(hash, "pdf")
	if err != nil {
		t.Fatal(err)
	}
	err = GetFileFromIpfs("blah", "pdf")
	if err == nil {
		t.Fatal("CAn retrieve non exiting file, quitting")
	}
	dummy := []byte("Hello World")
	_, err = IpfsHashData(dummy)
	if err != nil {
		t.Fatalf("Can't hash ipfs data")
	}
	var dummy2 []byte
	_, err = IpfsHashData(dummy2)
	if err != nil {
		t.Fatalf("Can't hash ipfs data")
	}
}
