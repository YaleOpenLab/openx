package platform

import (
	"log"
	"os"

	scan "github.com/Varunram/essentials/scan"
	xlm "github.com/Varunram/essentials/xlm"
	assets "github.com/Varunram/essentials/xlm/assets"
	wallet "github.com/Varunram/essentials/xlm/wallet"
	consts "github.com/YaleOpenLab/openx/consts"
	"github.com/pkg/errors"
)

// InitializePlatform starts the platform, initializing the platform seed and publickey
func InitializePlatform() error {
	var publicKey string
	var seed string
	var err error

	// check whether the home directory exists
	if _, err := os.Stat(consts.PlatformSeedFile); !os.IsNotExist(err) {
		// home dir exists, ask for password
		log.Println("ENTER YOUR PASSWORD TO DECRYPT THE PLATFORM SEED FILE")
		password, err := scan.RawPassword()
		if err != nil {
			return errors.Wrap(err, "couldn't scan raw password")
		}
		consts.PlatformPublicKey, consts.PlatformSeed, err = wallet.RetrieveSeed(consts.PlatformSeedFile, password)
		if err != nil {
			return err
		}

		log.Printf("PLATFORM SEED IS: %s\n PLATFORM PUBLIC KEY IS: %s\n", consts.PlatformSeed, consts.PlatformPublicKey)

		if consts.Mainnet {
			log.Println("mainnet init, stablecoin disabled")
			if !xlm.AccountExists(consts.PlatformPublicKey) {
				// ie we're on mainnet and the account doesn't have enough funds to start
				return errors.New("please refill the platform with xlm to be able to start openx. Min balance: 0.5XLM")
			}
			balance := xlm.GetNativeBalance(consts.PlatformPublicKey)
			if balance < 1.5 { // 0.5 min + 0.5x2 trustlines
				return errors.New("balance insufficient to run platform")
			}
		}
		return nil
	}
	// the home directory doesn't exist, two cases: seed doesn't exist or user has deleted it
	log.Println("DO YOU HAVE YOUR RAW PLATFORM SEED? IF SO, ENTER SEED. ELSE ENTER N")
	seed, err = scan.String()
	if err != nil {
		return errors.Wrap(err, "couldn't scan raw string")
	}
	if seed == "N" || seed == "n" {
		// no seed, no file, create new keypair
		log.Println("Enter a password to encrypt your master seed. Please store this in a very safe place. This prompt will not ask to confirm your password")
		password, err := scan.RawPassword()
		if err != nil {
			return err
		}
		publicKey, seed, err = wallet.NewSeedStore(consts.PlatformSeedFile, password)
		if err != nil {
			return errors.Wrap(err, "couldn't retrieve seed")
		}
		log.Println("Stored seed in seed file: ", consts.PlatformSeedFile)
		consts.PlatformPublicKey = publicKey
		consts.PlatformSeed = seed

		if consts.Mainnet {
			// in mainnet, don't init stablecoin
			return nil
		}
	} else {
		// no file but user remembers seed, retrieve pukbey
		log.Println("ENTER A PASSWORD TO ENCRYPT YOUR SEED")
		password, err := scan.RawPassword()
		if err != nil {
			return err
		}

		publicKey, err = wallet.ReturnPubkey(seed)
		if err != nil {
			return err
		}

		err = wallet.StoreSeed(seed, password, consts.PlatformSeed)
		if err != nil {
			return err
		}

		return nil
	}

	// comes here only if we're on testnet and we didn't have a seed earlier

	// getXLM and setup the account
	err = xlm.GetXLM(publicKey)
	if err != nil {
		return errors.Wrap(err, "error while getting xlm")
	}

	// set auth immutable on the account
	_, txhash, err := xlm.SetAuthImmutable(seed)
	log.Println("TX HASH FOR SETOPTIONS: ", txhash)
	if err != nil {
		log.Println("ERROR WHILE SETTING OPTIONS")
	}

	// make the platform trust the in house stablecoin for receiving payments
	txhash, err = assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, 10000000000, seed)
	if err != nil {
		log.Println("error while trusting stablecoin", consts.StablecoinCode, consts.StablecoinPublicKey, seed)
		return err
	}

	// send the platform some stablecoin to test if the trustline is setup correctly
	_, _, err = assets.SendAssetFromIssuer(consts.StablecoinCode, publicKey, 10, consts.StablecoinSeed, consts.StablecoinPublicKey)
	if err != nil {
		log.Println("error while sending stablecoin tp platform")
		log.Println("SEED: ", consts.StablecoinSeed)
		return err
	}

	log.Println("Platform trusts stablecoin: ", txhash)
	return err
}

// RefillPlatform asks friendbot for XLM in case the platform's funds are running low (below 20XLM).
// For obvious reasons, available only on testnet
func RefillPlatform(publicKey string) error {
	if consts.Mainnet {
		return errors.New("no provision to refill on mainnet") // refilling platform has to be done manually in the case of mainnet
	}

	balance := xlm.GetNativeBalance(publicKey)

	log.Println("Platform's balance is: ", balance)
	if balance < 21 {
		log.Println("Refilling platform balance")
		err := xlm.GetXLM(publicKey)
		if err != nil {
			return err
		}
	}
	return nil
}
