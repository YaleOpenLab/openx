package issuer

import (
	"log"
	"os"

	xlm "github.com/Varunram/essentials/crypto/xlm"
	wallet "github.com/Varunram/essentials/crypto/xlm/wallet"
	utils "github.com/Varunram/essentials/utils"
	"github.com/pkg/errors"
)

// issuer contains functions that can be called by asset issuers on Stellar

// GetPath returns the path of a specific project
func GetPath(path string, projIndex int) string {
	piS, _ := utils.ToString(projIndex)
	return path + piS + ".key"
}

// CreateFile creates a new empty keyfile
func CreateFile(issuerPath string, projIndex int) string {
	path := GetPath(issuerPath, projIndex)
	// we need to create this file
	os.Create(path)
	return path
}

// InitIssuer creates a new keypair and stores it in a file
func InitIssuer(issuerPath string, projIndex int, seedpwd string) error {
	seed, _, err := xlm.GetKeyPair()
	if err != nil {
		return errors.Wrap(err, "Error while generating keypair")
	}
	// store this seed in home/projects/projIndex.hex
	// we need a password for encrypting the seed
	path := CreateFile(issuerPath, projIndex)
	err = wallet.StoreSeed(seed, seedpwd, path)
	if err != nil {
		return errors.Wrap(err, "Error while storing seed")
	}
	return nil
}

// DeleteIssuer deletes the keyfile
func DeleteIssuer(issuerPath string, projIndex int) error {
	path := GetPath(issuerPath, projIndex)
	return os.Remove(path)
}

// FundIssuer creates an issuer account and funds it with a second account
func FundIssuer(issuerPath string, projIndex int, seedpwd string, funderSeed string) error {
	// need to read the seed from the file using the seedpwd
	path := GetPath(issuerPath, projIndex)
	pubkey, seed, err := wallet.RetrieveSeed(path, seedpwd)
	if err != nil {
		return errors.Wrap(err, "Error while retrieving seed")
	}
	log.Printf("Project Index: %d, Seed: %s, Address: %s", projIndex, seed, pubkey)
	_, txhash, err := xlm.SendXLMCreateAccount(pubkey, 100, funderSeed)
	if err != nil {
		return errors.Wrap(err, "Error while sending xlm to create account")
	}
	log.Printf("Txhash for setting up Project Issuer for project %d is %s", projIndex, txhash)
	_, txhash, err = xlm.SetAuthImmutable(seed)
	if err != nil {
		return errors.Wrap(err, "Error while setting auth immutable on account")
	}
	log.Printf("Txhash for setting Auth Immutable on project %d is %s", projIndex, txhash)
	return nil
}

// FreezeIssuer freezes the issuer account
func FreezeIssuer(issuerPath string, projIndex int, seedpwd string) (string, error) {
	path := GetPath(issuerPath, projIndex)
	_, seed, err := wallet.RetrieveSeed(path, seedpwd)
	if err != nil {
		return "", errors.Wrap(err, "Error while retrieving seed")
	}
	_, txhash, err := xlm.FreezeAccount(seed)
	if err != nil {
		return "", errors.Wrap(err, "Error while freezing account")
	}
	log.Println("Tx hash for freezing account is: ", txhash)
	return txhash, nil
}
