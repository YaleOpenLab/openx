package loader

import (
	"fmt"
	"github.com/pkg/errors"
	"log"

	"github.com/spf13/viper"

	// utils "github.com/Varunram/essentials/utils"
	// xlm "github.com/Varunram/essentials/xlm"
	xlm "github.com/Varunram/essentials/xlm"
	assets "github.com/Varunram/essentials/xlm/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	openx "github.com/YaleOpenLab/openx/platforms"
)

// Mainnet loads the stuff needed for mainnet. Ordering is very important since some consts need the others
// to function correctly
func Mainnet() error {
	log.Println("initializing openx mainnet..")
	var err error
	consts.SetConsts(true) // set in house  consts

	lim, _ := database.RetrieveAllUsersLim()
	if lim == 0 {
		// nothing exists, create dbs and buckets
		log.Println("creating mainnet home dir")
		database.CreateHomeDir()
	}

	// Initialize platform stuff like the platform seed
	err = openx.InitializePlatform()
	if err != nil {
		return err
	}

	// read from the config file
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

// StablecoinTrust creates a trustline with AnchorUSD on mainnet. We can't do this automatically since
// we need to wait for the platform to be funded before doing stuff on mainnet
func StablecoinTrust() error {
	_, txhash, err := xlm.SetAuthImmutable(consts.PlatformSeed)
	log.Println("TX HASH FOR SETOPTIONS: ", txhash)
	if err != nil {
		return errors.Wrap(err, "ERROR WHILE SETTING OPTIONS")
	}
	log.Println("TX HASH FOR SETTING AUTH IMMUTABLE: ", txhash)

	txhash, err = assets.TrustAsset(consts.AnchorUSDCode, consts.AnchorUSDAddress, 10000000000, consts.PlatformSeed)
	if err != nil {
		return errors.Wrap(err, "error while trusting stablecoin")
	}
	log.Println("TX HASH FOR TRUSTING ANCHORUSD: ", txhash)
	return nil
}
