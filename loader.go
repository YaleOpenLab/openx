package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"

	algorand "github.com/YaleOpenLab/openx/chains/algorand"
	stablecoin "github.com/YaleOpenLab/openx/chains/stablecoin"
	// utils "github.com/Varunram/essentials/utils"
	opensolar "github.com/YaleOpenLab/opensolar/core"
	xlm "github.com/YaleOpenLab/openx/chains/xlm"
	assets "github.com/YaleOpenLab/openx/chains/xlm/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	openx "github.com/YaleOpenLab/openx/platforms"
	"github.com/spf13/viper"
)

// imagine the loader like in a retro game, loading mainnet
func MainnetLoader() error {
	log.Println("initializing openx mainnet..")
	consts.SetConsts()

	var err error
	consts.DbDir = consts.HomeDir + "/mainnet/"                           // set mainnet db to open in spearate folder
	consts.PlatformSeedFile = consts.HomeDir + "/mainnetplatformseed.hex" // where the platform's seed is stored
	log.Println("DB DIRL: ", consts.DbDir)
	log.Println("DB DIRL: ", consts.PlatformSeedFile)

	lim, _ := database.RetrieveAllUsersLim()
	if lim == 0 {
		// nothing exists, create dbs and buckets
		log.Println("creating mainnet home dir")
		database.CreateHomeDir()
		err = openx.InitializePlatform()
		if err != nil {
			return err
		}
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading platform email from config file")
		return err
	}

	if !viper.IsSet("platformemail") {
		return errors.New("required param platformemail not found")
	}
	if !viper.IsSet("platformpass") {
		return errors.New("required param platformpass not found")
	}
	if !viper.IsSet("kycapikey") {
		return errors.New("required param kycapikey not found")
	}

	consts.PlatformEmail = viper.GetString("platformemail")
	consts.PlatformEmailPass = viper.GetString("password")
	consts.KYCAPIKey = viper.GetString("kycapikey")

	fmt.Printf("PLATFORM SEED IS: %s\n PLATFORM PUBLIC KEY IS: %s\n", consts.PlatformSeed, consts.PlatformPublicKey)
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

func TestnetLoader() error {
	log.Println("initializing openx testnet..")
	consts.SetConsts()
	database.CreateHomeDir()
	var err error
	// init stablecoin before platform so we don't have to create a stablecoin in case our dbdir is wiped
	consts.StablecoinPublicKey, consts.StablecoinSeed, err = stablecoin.InitStableCoin(consts.Mainnet) // start the stablecoin daemon
	if err != nil {
		log.Println("errored out while starting stablecoin")
		return err
	}

	err = opensolar.InitializePlatform()
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

	allContracts, err := opensolar.RetrieveAllProjects()
	if err != nil {
		log.Println("Error retrieving all projects from the database")
		return err
	}

	if len(allContracts) == 0 {
		log.Println("initialziing openx testnet")
		err = InsertDummyData(opts.Simulate)
		if err != nil {
			return err
		}
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
