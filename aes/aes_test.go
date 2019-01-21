// +build all travis

package aes

import (
	"log"
	"os"
	"testing"
)

func TestAes(t *testing.T) {
	password := "Cool"
	ciphertext, err := Encrypt([]byte("Hello World"), password)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Encrypted: %x\n", ciphertext)
	plaintext, err := Decrypt(ciphertext, password)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Decrypted: %s\n", plaintext)
	if string(plaintext) != "Hello World" {
		t.Errorf("Problem with Decryption")
	}
	_, err = Decrypt(ciphertext, "Notcool")
	if err == nil {
		t.Fatalf("Didn't catch error during decrpytion, exiting!")
	}
	_, err = Decrypt([]byte(""), "")
	if err == nil {
		t.Fatalf("Didn't catch error during decrpytion, exiting!")
	}
	data := []byte("This is test data")
	os.MkdirAll("test_files", os.ModePerm)
	err = EncryptFile("test_files/text.txt", data, password)
	if err != nil {
		t.Fatal(err)
	}
	err = EncryptFile("blah/blah.txt", data, password)
	if err == nil {
		t.Fatalf("not erroring out on file creation error, quitting!")
	}
	decryptedSlice, err := DecryptFile("test_files/text.txt", password)
	if err != nil {
		t.Fatal(err)
	}
	if string(decryptedSlice) != "This is test data" {
		t.Fatalf("Can't decrypt file, exiting!")
	}
	_, err = DecryptFile("test_files/text.txt", "Notcool")
	if err == nil {
		t.Fatal(err)
	}
	err = os.RemoveAll("test_files")
	if err != nil {
		t.Fatal(err)
	}
}
