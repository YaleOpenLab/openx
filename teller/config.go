package main

import (
	"log"

	utils "github.com/OpenFinancing/openfinancing/utils"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	"github.com/spf13/viper"
)

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
	seedpwd := viper.Get("seedpwd").(string)
	username := viper.Get("username").(string)
	password := utils.SHA3hash(viper.Get("password").(string))
	ApiUrl = viper.Get("apiurl").(string)

	err = LoginToPlatForm(username, password)
	if err != nil {
		return err
	}

	RecpSeed, err = wallet.DecryptSeed(LocalRecipient.U.EncryptedSeed, seedpwd)
	if err != nil {
		return err
	}

	RecpPublicKey, err = wallet.ReturnPubkey(RecpSeed)
	if err != nil {
		return err
	}

	return nil
}
