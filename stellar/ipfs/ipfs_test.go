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
	err = GetFileFromIpfs("/ipfs/"+hash, "pdf")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("DECRYPTED STRING IS: ", string1)
	_, err = ReadfromPdf("files/test.pdf") // get the data from the pdf as a datastream
	if err != nil {
		t.Fatal(err)
	}

	hash, err = IpfsHashPdf("files/test.pdf")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("HASH IS: ", hash)
	err = GetFileFromIpfs(hash, "pdf")
	if err != nil {
		t.Fatal(err)
	}
}
