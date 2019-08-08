package loader

import (
	"fmt"
	"github.com/pkg/errors"
	"log"

	"github.com/spf13/viper"

	algorand "github.com/YaleOpenLab/openx/chains/algorand"
	stablecoin "github.com/YaleOpenLab/openx/chains/stablecoin"
	// utils "github.com/Varunram/essentials/utils"
	xlm "github.com/YaleOpenLab/openx/chains/xlm"
	assets "github.com/YaleOpenLab/openx/chains/xlm/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	openx "github.com/YaleOpenLab/openx/platforms"

	opensolarconsts "github.com/YaleOpenLab/opensolar/consts"
)

func Testnet() error {

	opensolarconsts.HomeDir += "/testnet"
	opensolarconsts.DbDir = opensolarconsts.HomeDir + "/database/"                   // the directory where the database is stored (project info, user info, etc)
	opensolarconsts.PlatformSeedFile = opensolarconsts.HomeDir + "/platformseed.hex" // where the platform's seed is stored

	log.Println("initializing openx testnet..")
	consts.SetConsts(false)
	database.CreateHomeDir()
	var err error
	// init stablecoin before platform so we don't have to create a stablecoin in case our dbdir is wiped
	consts.StablecoinPublicKey, consts.StablecoinSeed, err = stablecoin.InitStableCoin(consts.Mainnet) // start the stablecoin daemon
	if err != nil {
		return errors.Wrap(err, "errored out while starting stablecoin")
	}

	err = openx.InitializePlatform()
	if err != nil {
		return err
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading platform email from config file")
		return err
	}

	// alogrand is supported only in testnet mode
	if viper.IsSet("algodAddress") {
		consts.AlgodAddress = viper.GetString("algodAddress")
	}
	if viper.IsSet("algodToken") {
		consts.AlgodToken = viper.GetString("algodToken")
	}
	if viper.IsSet("kmdAddress") {
		consts.KmdAddress = viper.GetString("kmdAddress")
	}
	if viper.IsSet("kmdToken") {
		consts.KmdToken = viper.GetString("kmdToken")
	}

	err = algorand.Init()
	if err != nil {
		return err
	}

	if !viper.IsSet("platformemail") {
		log.Println("platform email not set")
	} else {
		consts.PlatformEmail = viper.GetString("platformemail")
		log.Println("PLATFORM EMAIL: ", consts.PlatformEmail)
	}
	if !viper.IsSet("platformpass") {
		log.Println("platform email password not set")
	} else {
		consts.PlatformEmailPass = viper.Get("password").(string) // interface to string
	}
	if !viper.IsSet("kycapikey") {
		log.Println("kyc api key not set, kyc will be disabled")
	} else {
		consts.KYCAPIKey = viper.GetString("kycapikey")
	}

	fmt.Printf("PLATFORM SEED IS: %s\n PLATFORM PUBLIC KEY IS: %s\n", consts.PlatformSeed, consts.PlatformPublicKey)
	fmt.Printf("STABLECOIN PUBLICKEY IS: %s\nSTABLECOIN SEED is: %s\n\n", consts.StablecoinPublicKey, consts.StablecoinSeed)
	return nil
}

func StablecoinTrust() error {
	_, txhash, err := xlm.SetAuthImmutable(consts.PlatformSeed)
	log.Println("TX HASH FOR SETOPTIONS: ", txhash)
	if err != nil {
		log.Println("ERROR WHILE SETTING OPTIONS")
		return err
	}
	// make the platform trust the stablecoin for receiving payments
	txhash, err = assets.TrustAsset(consts.AnchorUSDCode, consts.AnchorUSDAddress, 10000000000, consts.PlatformSeed)
	if err != nil {
		log.Println("error while trusting stablecoin", consts.AnchorUSDCode, consts.AnchorUSDAddress, consts.PlatformSeed)
		return err
	}
	return nil
}
