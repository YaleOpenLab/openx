package platform

import (
	"fmt"
	"log"
	"os"

	consts "github.com/YaleOpenLab/smartPropertyMVP/stellar/consts"
	scan "github.com/YaleOpenLab/smartPropertyMVP/stellar/scan"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	wallet "github.com/YaleOpenLab/smartPropertyMVP/stellar/wallet"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
)

// InitializePlatform returns the platform publickey and seed
func InitializePlatform() (string, string, error) {
	var publicKey string
	var seed string
	var err error

	// now we can be sure we have the directory, check for seed
	if _, err := os.Stat(consts.PlatformSeedFile); !os.IsNotExist(err) {
		// the seed exists
		fmt.Println("ENTER YOUR PASSWORD TO DECRYPT THE SEED FILE")
		password, err := scan.ScanRawPassword()
		if err != nil {
			log.Println(err)
			return publicKey, seed, err
		}
		publicKey, seed, err = wallet.RetrieveSeed(consts.PlatformSeedFile, password)
		return publicKey, seed, err
	}
	// platform doesn't exist or user doesn't have encrypted file. Ask
	fmt.Println("DO YOU HAVE YOUR RAW SEED? IF SO, ENTER SEED. ELSE ENTER N")
	seed, err = scan.ScanForString()
	if err != nil {
		log.Println(err)
		return publicKey, seed, err
	}
	if seed == "N" || seed == "n" {
		// no seed, no file, create new keypair
		// need to pass the password for the  eed file
		fmt.Println("Enter a password to encrypt your master seed. Please store this in a very safe place. This prompt will not ask to confirm your password")
		password, err := scan.ScanRawPassword()
		if err != nil {
			return publicKey, seed, err
		}
		publicKey, seed, err = wallet.NewSeed(consts.PlatformSeedFile, password)
		if err != nil {
			log.Println(err)
			return publicKey, seed, err
		}
		err = xlm.GetXLM(publicKey)
	} else {
		// no file, retrieve pukbey
		// user has given us a seed, validate
		log.Println("ENTER A PASSWORD TO DECRYPT YOUR SEED")
		password, err := scan.ScanRawPassword()
		if err != nil {
			return publicKey, seed, err
		}
		publicKey, err = wallet.RetrievePubkey(seed, consts.PlatformSeedFile, password)
		if err != nil {
			log.Println(err)
			return publicKey, seed, err
		}
	}
	err = xlm.GetXLM(publicKey)
	return publicKey, seed, err
}

// RefillPlatform checks whether the publicKey passed has any xlm and if its balance
// is less than 21 XLM, it proceeds to ask the friendbot for more test xlm
func RefillPlatform(publicKey string) error {
	// check whether the investor has XLM already
	balance, err := xlm.GetNativeBalance(publicKey)
	if err != nil {
		return err
	}
	// balance is in string, convert to int
	balanceI := utils.StoF(balance)
	log.Println("Platform's balance is: ", balanceI)
	if balanceI < 21 { // 1 to account for fees
		// get coins if balance is this low
		log.Println("Refilling platform balance")
		err := xlm.GetXLM(publicKey)
		// TODO: in future, need to refill platform sufficiently well and interact
		// with a cold wallet that we have previously set
		if err != nil {
			return err
		}
	}
	return nil
}
