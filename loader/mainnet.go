package loader

import (
	"fmt"
	"github.com/pkg/errors"
	"log"

	"github.com/spf13/viper"

	// utils "github.com/Varunram/essentials/utils"
	// xlm "github.com/YaleOpenLab/openx/chains/xlm"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	openx "github.com/YaleOpenLab/openx/platforms"
)

// imagine the loader like in a retro game, loading mainnet
func Mainnet() error {
	log.Println("initializing openx mainnet..")
	consts.SetConsts(true)

	var err error

	lim, _ := database.RetrieveAllUsersLim()
	if lim == 0 {
		// nothing exists, create dbs and buckets
		log.Println("creating mainnet home dir")
		database.CreateHomeDir()
		err = openx.InitializePlatform()
		if err != nil {
			return err
		}
	} else {
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
