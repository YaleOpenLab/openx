package main

import (
	"log"
	"time"

	consts "github.com/YaleOpenLab/openx/consts"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
	"github.com/spf13/viper"
)

func RefreshLogin(username string, pwhash string) error {
	// refresh login runs once every 5 minutes in order to fetch the latest recipient details
	// for eg, if the recipient loads his balance on the platform, we need it to be reflected on
	// the teller
	var err error
	for {
		err = LoginToPlatForm(username, pwhash)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("UPDATED RECIPIENT")
		}
		time.Sleep(consts.LoginRefreshInterval * time.Minute)
	}
}

// SetupConfig reads required values from the config file
func SetupConfig() error {
	var err error
	viper.SetConfigType("yaml")
	viper.SetConfigName("tellerconfig")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading email values from config file")
		return err
	}

	PlatformPublicKey = viper.Get("platformPublicKey").(string)
	LocalSeedPwd = viper.Get("seedpwd").(string)                   // seed password used to unlock the seed of the recipient on the platform
	username := viper.Get("username").(string)                 // username of the recipient on the platform
	password := utils.SHA3hash(viper.Get("password").(string)) // password of the recipient on the platform
	ApiUrl = viper.Get("apiurl").(string)                      // ApiUrl of the remote / local openx node
	mapskey := viper.Get("mapskey").(string)                   // google maps API key. Need to activate it

	err = LoginToPlatForm(username, password)
	if err != nil {
		log.Println("Error while logging on to the platform", err)
		return err
	}

	RecpSeed, err = wallet.DecryptSeed(LocalRecipient.U.EncryptedSeed, LocalSeedPwd)
	if err != nil {
		log.Println("Error while decrypting seed", err)
		return err
	}

	RecpPublicKey, err = wallet.ReturnPubkey(RecpSeed)
	if err != nil {
		log.Println("Error while returning publickey", err)
		return err
	}

	// check for device id and set it if none is set
	err = CheckDeviceID()
	if err != nil {
		return err
	}

	DeviceId, err = GetDeviceID() // Stores DeviceId
	if err != nil {
		return err
	}

	err = StoreStartTime()
	if err != nil {
		return err
	}

	err = StoreLocation(mapskey) // stores DeviceLocation
	if err != nil {
		return err
	}

	return nil
}
