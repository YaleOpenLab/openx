package loader

import (
	"fmt"
	"log"

	"github.com/Varunram/essentials/email"

	"github.com/pkg/errors"

	"github.com/spf13/viper"

	algorand "github.com/Varunram/essentials/algorand"
	stablecoin "github.com/Varunram/essentials/xlm/stablecoin"

	// utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	openx "github.com/YaleOpenLab/openx/platforms"
)

// Testnet loads the stuff needed for testnet. Ordering is very important since some consts need the others
// to function correctly
func Testnet() error {
	log.Println("initializing openx testnet..")
	consts.SetConsts(false)
	database.CreateHomeDir()
	var err error
	// init stablecoin before platform so we don't have to create a stablecoin in case our dbdir is wiped
	consts.StablecoinPublicKey, consts.StablecoinSeed, err = stablecoin.InitStableCoin()
	if err != nil {
		return errors.Wrap(err, "errored out while starting stablecoin")
	}

	// start platform
	err = openx.InitializePlatform()
	if err != nil {
		return err
	}

	// read from consts
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading platform email from config file")
		return err
	}

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
		consts.PlatformEmailPass = viper.GetString("platformpass")
	}
	if !viper.IsSet("kycapikey") {
		log.Println("kyc api key not set, kyc will be disabled")
	} else {
		consts.KYCAPIKey = viper.GetString("kycapikey")
	}

	email.SetConsts(consts.PlatformEmail, consts.PlatformEmailPass)
	fmt.Printf("PLATFORM SEED IS: %s\n PLATFORM PUBLIC KEY IS: %s\n", consts.PlatformSeed, consts.PlatformPublicKey)
	fmt.Printf("STABLECOIN PUBLICKEY IS: %s\nSTABLECOIN SEED is: %s\n\n", consts.StablecoinPublicKey, consts.StablecoinSeed)
	return nil
}
