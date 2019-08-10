package wallet

import (
	"github.com/pkg/errors"
	"log"
	"os"

	aes "github.com/Varunram/essentials/aes"
	"github.com/stellar/go/keypair"
)

// the wallet package contains stellar specific wallet functions

// NewSeedStore creates a new seed and stores the seed in an encrypted form in the passed path
func NewSeedStore(path string, password string) (string, string, error) {
	// these can store the file in any path passed to them
	var seed string
	var publicKey string
	var err error

	pair, err := keypair.Random()
	if err != nil {
		return publicKey, seed, err
	}
	seed = pair.Seed()
	publicKey = pair.Address()
	log.Printf("\nTHE GENERATED SEED IS: %s\nAND YOUR PUBLIC KEY IS: %s\nKEEP IT SUPER SAFE OR YOU MIGHT NOT HAVE ACCESS TO THESE FUNDS AGAIN \n", seed, publicKey)
	err = StoreSeed(seed, password, path) // store the seed in a secure location
	return publicKey, seed, err
}

// StoreSeed encrypts and stores the seed
func StoreSeed(seed string, password string, path string) error {
	// these can store the file ion any path passed to them
	// now we can be sure we have the directory, check for seed
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// file not created before, create
		log.Println("file doesn't exist, creating a new one")
		file, err := os.Create(path)
		if err != nil {
			log.Println("ERROR WHILE CREATING FILE: ", err)
			return err
		}
		file.Close()
	}
	err := aes.EncryptFile(path, []byte(seed), password)
	if err != nil {
		return errors.Wrap(err, "could not encrypt file")
	}
	_, err = aes.DecryptFile(path, password)
	return err
}

// RetrieveSeed retrieves the seed and publickey when an encrypted file path is passed
func RetrieveSeed(path string, password string) (string, string, error) {
	var publicKey string
	var seed string
	data, err := aes.DecryptFile(path, password)
	if err != nil {
		return publicKey, seed, errors.Wrap(err, "could not decrypt file")
	}
	seed = string(data)
	keyp, err := keypair.Parse(seed)
	return keyp.Address(), seed, errors.Wrap(err, "could not parse seed to get keypair")
}

// DecryptSeed decrpyts the encrypted seed and returns the raw unencrypted seed
func DecryptSeed(encryptedSeed []byte, seedpwd string) (string, error) {
	data, err := aes.Decrypt(encryptedSeed, seedpwd)
	return string(data), err
}

// ReturnPubkey returns the pubkey when passed the seed
func ReturnPubkey(seed string) (string, error) {
	if len(seed) == 0 {
		return seed, errors.New("Empty Seed passed!")
	}
	keyp, err := keypair.Parse(seed)
	return keyp.Address(), errors.Wrap(err, "could not parse seed to get keypair")
}
