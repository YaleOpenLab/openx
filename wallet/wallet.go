package wallet

import (
	"github.com/pkg/errors"
	"log"

	aes "github.com/YaleOpenLab/openx/aes"
	"github.com/stellar/go/keypair"
)

// NewSeed creates a new seed and stores the seed in an encrypted form in the
// specified path
func NewSeed(path string, password string) (string, string, error) {
	// these can store the file in any path passed to them
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
		return errors.Wrap(err, "could not encrypt file")
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
		return publicKey, seed, errors.Wrap(err, "could not decrypt file")
	}
	seed = string(data)
	keyp, err := keypair.Parse(seed)
	return keyp.Address(), seed, errors.Wrap(err, "could not parse seed to get keypair")
}

func RetrieveEscrowPubkey(path string, password string) (string, error) {
	var pubkey string
	data, err := aes.DecryptFile(path, password)
	if err != nil {
		return pubkey, errors.Wrap(err, "could not decrypt file")
	}
	pubkey = string(data)
	return pubkey, nil
}

// RetrieveAndStorePubkey restores the publicKey when passed a seed and stores the
// seed in an encrypted format in the specified path
func RetrieveAndStorePubkey(seed string, path string, password string) (string, error) {
	var publicKey string
	keyp, err := keypair.Parse(seed)
	if err != nil {
		return publicKey, errors.Wrap(err, "could not parse seed to get keypair")
	}

	publicKey = keyp.Address()
	StoreSeed(seed, password, path)
	return publicKey, nil
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
