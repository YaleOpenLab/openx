package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"

	algorand "github.com/Varunram/essentials/crypto/algorand"
	stablecoin "github.com/Varunram/essentials/crypto/stablecoin"
	// utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	"github.com/spf13/viper"
)

// imagine the loader like in a retro game, loading mainnet
func MainnetLoader() error {
	log.Println("initializing openx mainnet..")

	var err error
	consts.DbDir = consts.HomeDir + "/database/mainnet/"                  // set mainnet db to open in spearate folder
	consts.PlatformSeedFile = consts.HomeDir + "/mainnetplatformseed.hex" // where the platform's seed is stored

	err = opensolar.InitializePlatform()
	if err != nil {
		return err
	}

	lim, _ := database.RetrieveAllUsersLim()
	if lim == 0 {
		// nothing exists, create dbs and buckets
		log.Println("creating mainnet home dir")
		database.CreateHomeDir()
		log.Println("created mainnet home dir")
		// Create an admin investor

		log.Println("seeding dci as admin investor")
		inv, err := database.NewInvestor("dci@mit.edu", "p", "x", "dci")
		if err != nil {
			return err
		}
		inv.U.Inspector = true
		inv.U.Kyc = true
		inv.U.Admin = true // no handlers for the admin bool, just set it wherever needed.
		inv.U.Reputation = 100000
		inv.U.Notification = true
		err = inv.U.Save()
		if err != nil {
			return err
		}
		err = inv.U.AddEmail("varunramganesh@gmail.com") // change this to something more official later
		if err != nil {
			return err
		}
		err = inv.Save()
		if err != nil {
			return err
		}
		log.Println("Please seed DCI pubkey: ", inv.U.StellarWallet.PublicKey, " with funds")

		// Create an admin recipient
		log.Println("seeding vx as admin investor")
		recp, err := database.NewRecipient("varunramganesh@gmail.com", "p", "x", "vg")
		if err != nil {
			return err
		}
		recp.U.Inspector = true
		recp.U.Kyc = true
		recp.U.Admin = true // no handlers for the admin bool, just set it wherever needed.
		recp.U.Reputation = 100000
		recp.U.Notification = true
		err = recp.U.Save()
		if err != nil {
			return err
		}
		err = recp.U.AddEmail("varunramganesh@gmail.com")
		if err != nil {
			return err
		}
		err = recp.Save()
		if err != nil {
			return err
		}
		log.Println("Please seed Varunram's pubkey: ", recp.U.StellarWallet.PublicKey, " with funds")

		orig, err := opensolar.NewOriginator("martin", "p", "x", "Martin Wainstein", "California", "Project Originator")
		if err != nil {
			return err
		}

		log.Println("Please seed Martin's pubkey: ", orig.U.StellarWallet.PublicKey, " with funds")

		contractor, err := opensolar.NewContractor("samuel", "p", "x", "Samuel Visscher", "Georgia", "Project Contractor")
		if err != nil {
			return err
		}

		log.Println("Please seed Samuel's pubkey: ", contractor.U.StellarWallet.PublicKey, " with funds")

		var project opensolar.Project
		project.Index = 1
		project.TotalValue = 8000
		project.Name = "SU Pasto School, Aibonito"
		project.Metadata = "MIT/Yale Pilot 2"
		project.OriginatorIndex = orig.U.Index
		project.ContractorIndex = contractor.U.Index
		project.EstimatedAcquisition = 5
		project.Stage = 4
		project.MoneyRaised = 0
		// add stuff in here as necessary
		err = project.Save()
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
