package platform

import (
	"fmt"
	"log"
	"os"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	scan "github.com/YaleOpenLab/openx/scan"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
	xlm "github.com/YaleOpenLab/openx/xlm"
)

// the platform structure is the backend representation of the frontend UI.
// on a very low level, this should just be a pubkey + seed pair. Each platform
// needs to be hosted somewhere, so it is necessary that each platform should have
// its own pubkey and seed pair
// InitializePlatform returns the platform publickey and seed
// We have a new model in which we have a new seed for every project that is
// advertised on the platform. The way this would wokr is that it sets up the assets,
// and then we freeze the account to freeze issuance. This would mean we would no longer
// be able to transact with the account although people can still send funds to it
// in this case, they would send us back DebtAssets provided they have sufficient
// stableUSD balance. Else they would not be able to trigger payback.
// TODO: this password could also be agreed upon by the party of investors and the recipient
// so that we act as a trustless entity, which is cool. This has to be done on the frontend preferably
// the main platform still has its pubkey and seed pair and sends funds out to issuers
// but is not directly involved in the setting up of trustlines
func InitializePlatform() error {
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
			return err
		}
		publicKey, seed, err = wallet.RetrieveSeed(consts.PlatformSeedFile, password)
		return err
	}
	// platform doesn't exist or user doesn't have encrypted file. Ask
	fmt.Println("DO YOU HAVE YOUR RAW SEED? IF SO, ENTER SEED. ELSE ENTER N")
	seed, err = scan.ScanForString()
	if err != nil {
		log.Println(err)
		return err
	}
	if seed == "N" || seed == "n" {
		// no seed, no file, create new keypair
		// need to pass the password for the  eed file
		fmt.Println("Enter a password to encrypt your master seed. Please store this in a very safe place. This prompt will not ask to confirm your password")
		password, err := scan.ScanRawPassword()
		if err != nil {
			return err
		}
		publicKey, seed, err = wallet.NewSeed(consts.PlatformSeedFile, password)
		if err != nil {
			log.Println(err)
			return err
		}
		err = xlm.GetXLM(publicKey)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		// no file, retrieve pukbey
		// user has given us a seed, validate
		log.Println("ENTER A PASSWORD TO DECRYPT YOUR SEED")
		password, err := scan.ScanRawPassword()
		if err != nil {
			return err
		}
		publicKey, err = wallet.RetrieveAndStorePubkey(seed, consts.PlatformSeedFile, password)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	_ = xlm.GetXLM(publicKey) // the API request errors out even on success, so
	// don't catch this error
	_, txhash, err := xlm.SetAuthImmutable(seed)
	log.Println("TX HASH FOR SETOPTIONS: ", txhash)
	if err != nil {
		log.Println("ERROR WHILE SETTING OPTIONS")
	}
	// make the platform trust the stablecoin for receiving payments
	txhash, err = assets.TrustAsset(consts.Code, consts.StablecoinPublicKey, "10000000000", publicKey, seed)
	if err != nil {
		return err
	}

	_, _, err = assets.SendAssetFromIssuer(consts.Code, publicKey, "10", consts.StablecoinSeed, consts.StablecoinPublicKey)
	if err != nil {
		log.Println("SEED: ", consts.StablecoinSeed)
		return err
	}

	log.Println("Platform trusts stablecoin: ", txhash)
	consts.PlatformPublicKey = publicKey
	consts.PlatformSeed = seed
	return err
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
		// refill platform sufficiently well and interact with a cold wallet that we
		// have previously set earlier to avoid hacks and similar
		if err != nil {
			return err
		}
	}
	return nil
}
