package main

import (
	"fmt"
	"log"
	"os"

	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	// ipfs "github.com/YaleOpenLab/openx/ipfs"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	rpc "github.com/YaleOpenLab/openx/rpc"
	// scan "github.com/YaleOpenLab/openx/scan"
	// oracle "github.com/YaleOpenLab/openx/oracle"
	algorand "github.com/Varunram/essentials/crypto/algorand"
	stablecoin "github.com/Varunram/essentials/crypto/stablecoin"
	// utils "github.com/Varunram/essentials/utils"
	// scan "github.com/YaleOpenLab/openx/scan"
	// wallet "github.com/YaleOpenLab/openx/wallet"
	// xlm "github.com/YaleOpenLab/openx/xlm"
	// assets "github.com/YaleOpenLab/openx/assets"
	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
)

// the server powering the openx platform of platforms. There are two clients that can be used
// with the backend - ofcli and emulator
// refer https://github.com/stellar/go/blob/master/build/main_test.go in case the stellar
// go SDK docs are insufficient.
var opts struct {
	Insecure bool `short:"i" description:"Start the API using http. Not recommended"`
	Port     int  `short:"p" description:"The port on which the server runs on. Default: HTTPS/8080"`
	Simulate bool `short:"t" description:"Simulate the test database with demo values (last updated: April 2019)"`
	Mainnet  bool `short:"m" description:"Switch mainnet mode on"`
}

// ParseConfig parses CLI parameters passed
func ParseConfig(args []string) (bool, int, error) {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return false, -1, err
	}
	port := consts.DefaultRpcPort
	if opts.Port != 0 {
		port = opts.Port
	}
	if opts.Mainnet {
		consts.Mainnet = true
	}
	return opts.Insecure, port, nil
}

// StartPlatform starts the platform
func StartPlatform() error {
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

	allContracts, err := opensolar.RetrieveAllProjects()
	if err != nil {
		log.Println("Error retrieving all projects from the database")
		return err
	}

	if len(allContracts) == 0 {
		log.Println("Populating database with test values")
		err = InsertDummyData(opts.Simulate)
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

	if !consts.Mainnet {
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
	}

	return nil
}

func main() {
	var err error
	insecure, port, err := ParseConfig(os.Args) // parseconfig should be before StartPlatform to parse the mainnet bool
	if err != nil {
		log.Fatal(err)
	}

	err = StartPlatform()
	if err != nil {
		log.Fatal(err)
	}

	// run this only when you need to monitor the tellers. Not required for local testing.
	// go opensolar.MonitorTeller(1)
	fmt.Printf("PLATFORM SEED IS: %s\n PLATFORM PUBLIC KEY IS: %s\n", consts.PlatformSeed, consts.PlatformPublicKey)
	fmt.Printf("STABLECOIN PUBLICKEY IS: %s\nSTABLECOIN SEED is: %s\n\n", consts.StablecoinPublicKey, consts.StablecoinSeed)
	rpc.StartServer(port, insecure)
}
