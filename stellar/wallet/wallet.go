package wallet

import (
	"fmt"
	"log"

	aes "github.com/YaleOpenLab/smartPropertyMVP/stellar/aes"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/stellar/go/keypair"
)

// NewSeed creates a new platform and returns the platform struct
func NewSeed(path string) (string, string, error) {
	// this function is used to generate a new keypair which will be assigned to
	// the platform
	var seed string
	var publicKey string
	var err error

	pair, err := keypair.Random()
	seed = pair.Seed()
	publicKey = pair.Address()
	log.Printf("\nTHE STABLECOIN SEED IS: %s\nAND YOUR PUBLIC KEY IS: %s\nKEEP IT SUPER SAFE OR YOU MIGHT NOT HAVE ACCESS TO THESE FUNDS AGAIN \n", seed, publicKey)
	fmt.Println("Enter a password to encrypt your stablecoin's master seed. Please store this in a very safe place. This prompt will not ask to confirm your password")
	password, err := utils.ScanRawPassword()
	if err != nil {
		return publicKey, seed, err
	}
	StoreSeed(seed, password, path) // store the seed in a secure location
	// err = xlm.GetXLM(publicKey) this should be only when we are setting up an account
	return publicKey, seed, err
}

// StoreSeed encrypts and stores the seed of the platform in a file
// named stablecoinseed.hex at the root directory
func StoreSeed(seed string, password string, path string) error {
	// handler to store the seed over at stablecoinseed.hex
	// person either needs to store this file and remember the password or has to
	// remember the seed in order to access the platform again
	aes.EncryptFile(path, []byte(seed), password)
	decrypted, err := aes.DecryptFile(path, password)
	if err != nil {
		return err
	}
	if seed != string(decrypted) {
		// something wrong with encryption, exit
		log.Fatal("Encrypted and decrypted seeds don't match, exiting!")
	}
	fmt.Println("Successfully encrypted your seed at: ", path)
	return nil
}

// RestorePlatformFromFile restores the platform directly from the file
func RetrieveSeed(path string, password string) (string, string, error) {
	var publicKey string
	var seed string
	data, err := aes.DecryptFile(path, password)
	if err != nil {
		return publicKey, seed, err
	}
	seed = string(data)
	keyp, err := keypair.Parse(seed)
	if err != nil {
		return publicKey, seed, err
	} else {
		publicKey = keyp.Address()
	}
	return publicKey, seed, nil
}

// RetrievePubkey restores the platform struct when passed the seed
func RetrievePubkey(seed string, path string) (string, error) {
	// this function should be used when the platform admin remembers the seed but
	// does not possess the encrypted file. The seed is what's needed to access
	// the account, so we don't restrict access
	var publicKey string
	keyp, err := keypair.Parse(seed)
	if err != nil {
		return publicKey, err
	} else {
		publicKey = keyp.Address()
	}
	log.Println("ENTER A PASSWORD TO ENCRYPT YOUR SEED AT: ", path)
	password, err := utils.ScanRawPassword()
	if err != nil {
		return publicKey, err
	}
	StoreSeed(seed, password, path)
	return publicKey, nil
}

func DecryptSeed(encryptedSeed []byte, seedpwd string) (string, error) {
	// need to call the aes decrypt function here
	// func Decrypt(data []byte, passphrase string) ([]byte, error) {
	log.Println("ENCRPYTED SEED: ", encryptedSeed, seedpwd)
	data, err := aes.Decrypt(encryptedSeed, seedpwd)
	log.Println("DECRYPTED DATA: ", data)
	return string(data), err
}
