package platform

import (
	"fmt"
	"log"
	"os"

	consts "github.com/YaleOpenLab/smartPropertyMVP/stellar/consts"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	wallet "github.com/YaleOpenLab/smartPropertyMVP/stellar/wallet"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
)

// InitializePlatform returns the platform structure and the seed
func InitializePlatform() (string, string, error) {
	var publicKey string
	var seed string
	var err error

	// now we can be sure we have the directory, check for seed
	if _, err := os.Stat(consts.PlatformSeedFile); !os.IsNotExist(err) {
		// the seed exists
		fmt.Println("ENTER YOUR PASSWORD TO DECRYPT THE SEED FILE")
		password, err := utils.ScanRawPassword()
		if err != nil {
			log.Println(err)
			return publicKey, seed, err
		}
		publicKey, seed, err = wallet.RetrieveSeed(consts.PlatformSeedFile, password)
		return publicKey, seed, err
	}
	// platform doesn't exist or user doesn't have encrypted file. Ask
	fmt.Println("DO YOU HAVE YOUR RAW SEED? IF SO, ENTER SEED. ELSE ENTER N")
	seed, err = utils.ScanForString()
	if err != nil {
		log.Println(err)
		return publicKey, seed, err
	}
	if seed == "N" || seed == "n" {
		// no seed, no file, carete new keypair
		publicKey, seed, err = wallet.NewSeed(consts.PlatformSeedFile)
		if err != nil {
			log.Println(err)
			return publicKey, seed, err
		}
		err = xlm.GetXLM(publicKey)
	} else {
		// no file, retrieve pukbey
		// user has given us a seed, validate
		publicKey, err = wallet.RetrievePubkey(seed, consts.PlatformSeedFile)
		if err != nil {
			log.Println(err)
			return publicKey, seed, err
		}
	}
	err = xlm.GetXLM(publicKey)
	return publicKey, seed, err
}

func RefillPlatform(publicKey string) error {
	// when I am creating an account, I will have a PublicKey and Seed, so
	// don't need them here
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
