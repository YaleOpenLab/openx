package wallet

import (
	"log"

	aes "github.com/OpenFinancing/openfinancing/aes"
	"github.com/stellar/go/keypair"
)

// NewSeed creates a new seed and stores the seed in an encrypted form in the
// specified path
func NewSeed(path string, password string) (string, string, error) {
	// these can store the file ion any path passed to them
	var seed string
	var publicKey string
	var err error

	pair, err := keypair.Random()
	seed = pair.Seed()
	publicKey = pair.Address()
	log.Printf("\nTHE GENERATED SEED IS: %s\nAND YOUR PUBLIC KEY IS: %s\nKEEP IT SUPER SAFE OR YOU MIGHT NOT HAVE ACCESS TO THESE FUNDS AGAIN \n", seed, publicKey)
	StoreSeed(seed, password, path) // store the seed in a secure location
	return publicKey, seed, err
}

// StoreSeed encrypts and stores the seed in a file
func StoreSeed(seed string, password string, path string) error {
	// these can store the file ion any path passed to them
	err := aes.EncryptFile(path, []byte(seed), password)
	if err != nil {
		return err
	}
	_, err = aes.DecryptFile(path, password)
	return err
}

// RetrieveSeed retrieves the seed and the publicket when an encrypted file path
// is passed to it
func RetrieveSeed(path string, password string) (string, string, error) {
	var publicKey string
	var seed string
	data, err := aes.DecryptFile(path, password)
	if err != nil {
		return publicKey, seed, err
	}
	seed = string(data)
	keyp, err := keypair.Parse(seed)
	return keyp.Address(), seed, err
}

// RetrievePubkey restores the publicKey when passed a seed and stores the
// seed in an encrypted format in the specified path
func RetrieveAndStorePubkey(seed string, path string, password string) (string, error) {
	var publicKey string
	keyp, err := keypair.Parse(seed)
	if err != nil {
		return publicKey, err
	} else {
		publicKey = keyp.Address()
	}
	StoreSeed(seed, password, path)
	return publicKey, nil
}

// DecryptSeed decrpyts the encrypted seed and returns the raw unencrypted seed
func DecryptSeed(encryptedSeed []byte, seedpwd string) (string, error) {
	data, err := aes.Decrypt(encryptedSeed, seedpwd)
	return string(data), err
}

func ReturnPubkey(seed string) (string, error) {
	keyp, err := keypair.Parse(seed)
	return keyp.Address(), err
}
