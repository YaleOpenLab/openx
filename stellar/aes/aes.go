package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/ioutil"
	"log"
	"os"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
)

func Encrypt(data []byte, passphrase string) []byte {
	key := []byte(utils.SHA3hash(passphrase)[96:128]) // last 32 characters in hash
	block, _ := aes.NewCipher(key)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Failed to initialize a new AES GCM while encrypting")
		log.Fatal(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println("Error while reading gcm bytes")
		log.Fatal(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(utils.SHA3hash(passphrase)[96:128]) // last 32 characters in hash
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Error while initializing cipher decryption")
		log.Fatal(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Failed to initialize a new AES GCM while decrypting")
		log.Fatal(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println("Failed to open gcm while decrypting")
		log.Fatal(err.Error())
	}
	return plaintext
}

func EncryptFile(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(Encrypt(data, passphrase))
}

func DecryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return Decrypt(data, passphrase)
}

/*
func Test() {
	password := "Cool"
	ciphertext := Encrypt([]byte("Hello World"), password)
	fmt.Printf("Encrypted: %x\n", ciphertext)
	plaintext := Decrypt(ciphertext, password)
	fmt.Printf("Decrypted: %s\n", plaintext)

	password2 := "cooler"
	EncryptFile("sample.txt", []byte("Hello World"), password2)
	DecryptFile("sample.txt", password2)
}
*/
