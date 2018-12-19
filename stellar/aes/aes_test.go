package aes

import (
	"log"
	"testing"
)

func TestAes(t *testing.T) {
	password := "Cool"
	ciphertext := Encrypt([]byte("Hello World"), password)
	log.Printf("Encrypted: %x\n", ciphertext)
	plaintext := Decrypt(ciphertext, password)
	if string(plaintext) != "Hello World" {
		t.Errorf("Problem with Decryption")
	}
	log.Printf("Decrypted: %s\n", plaintext)
}
