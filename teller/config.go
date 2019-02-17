package main

import (
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
	"github.com/spf13/viper"
)

// StartTeller starts the teller
func StartTeller() error {
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
	LocalSeedPwd = viper.Get("seedpwd").(string)               // seed password used to unlock the seed of the recipient on the platform
	username := viper.Get("username").(string)                 // username of the recipient on the platform
	password := utils.SHA3hash(viper.Get("password").(string)) // password of the recipient on the platform
	ApiUrl = viper.Get("apiurl").(string)                      // ApiUrl of the remote / local openx node
	mapskey := viper.Get("mapskey").(string)                   // google maps API key. Need to activate it
	LocalProjIndex = utils.ItoS(viper.Get("projIndex").(int))  // get the project index which should be in the config file
	assetName := viper.Get("assetName").(string)               // used to double check before starting the teller

	projIndex, err := GetProjectIndex(assetName)
	if err != nil {
		return err
	}

	if utils.ItoS(projIndex) != LocalProjIndex {
		log.Println("Project indices don't match, quitting!")
		return fmt.Errorf("Project indices don't match, quitting!")
	}

	// don't allow login before this since that becomes an attack vector where a person can guess
	// multiple passwords
	err = LoginToPlatform(username, password)
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

	if RecpPublicKey != LocalRecipient.U.PublicKey {
		log.Println("PUBLIC KEYS DON'T MATCH, QUITTING!")
		return fmt.Errorf("PUBLIC KEYS DON'T MATCH, QUITTING!")
	}

	LocalProject, err = GetLocalProjectDetails(LocalProjIndex)
	if err != nil {
		return err
	}

	if LocalProject.Stage < 4 {
		log.Println("TRYING TO INSTALL A PROJECT THAT HASN'T BEEN FUNDED YET, QUITTING!")
		return fmt.Errorf("TRYING TO INSTALL A PROJECT THAT HASN'T BEEN FUNDED YET, QUITTING!")
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

	err = GetPlatformEmail()
	if err != nil {
		return err
	}

	DeviceInfo = "Raspberry Pi3 Model B+"
	return nil
}
