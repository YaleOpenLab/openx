package main

import (
	"fmt"

	// "github.com/pkg/errors"
	"log"
	"os"

	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	loader "github.com/YaleOpenLab/openx/loader"

	// ipfs "github.com/YaleOpenLab/openx/ipfs"
	// opensolar "github.com/YaleOpenLab/opensolar/consts"
	rpc "github.com/YaleOpenLab/openx/rpc"
	// scan "github.com/YaleOpenLab/openx/scan"
	// oracle "github.com/YaleOpenLab/openx/oracle"
	// algorand "github.com/Varunram/essentials/algorand"
	// stablecoin "github.com/Varunram/essentials/xlm/stablecoin"
	utils "github.com/Varunram/essentials/utils"
	// scan "github.com/YaleOpenLab/openx/scan"
	// wallet "github.com/YaleOpenLab/openx/wallet"
	xlm "github.com/Varunram/essentials/xlm"
	// assets "github.com/Varunram/essentials/xlm/assets"
	// assets "github.com/YaleOpenLab/openx/assets"
	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
)

// the backend server powering the openx platform of platforms

var opts struct {
	Insecure  bool `short:"i" description:"Start the API using http. Not recommended"`
	Port      int  `short:"p" description:"The port on which the server runs on" default:"0"`
	Simulate  bool `short:"t" description:"Simulate the test database with demo values (last updated: April 2019)"`
	Mainnet   bool `short:"m" description:"Switch mainnet mode on"`
	Trustline bool `short:"x" description:"create trustlines from platform seed to anchorUSD"`
	Rescue    bool `short:"r" description:"start rescue mode"`
}

var InsecureSet bool

// ParseConfFile parses stuff from the config file provided
func ParseConfFile() (bool, int, error) {

	var port int
	var insecure bool
	var err error

	_, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("error while reading platform email from config file")
		return insecure, port, err
	}

	if opts.Port != 0 {
		port = opts.Port
	} else if viper.IsSet("port") {
		port = viper.GetInt("port")
	}

	if opts.Insecure {
		InsecureSet = true
	} else if viper.IsSet("insecure") {
		insecure = viper.GetBool("insecure")
		InsecureSet = insecure
	}

	if viper.IsSet("mainnet") {
		consts.Mainnet = viper.GetBool("mainnet")
	}

	return insecure, port, nil
}

func main() {
	var err error
	insecure, port, err := ParseConfFile()
	if err != nil {
		log.Fatal(err)
	}

	if consts.Mainnet {
		err = loader.Mainnet()
		if err != nil {
			log.Fatal(err)
		}
		// opensolar.SetMnConsts()
		if opts.Trustline {
			loader.StablecoinTrust()
		}
	} else {
		err = loader.Testnet()
		if err != nil {
			log.Fatal(err)
		}
		var admin database.User
		admin.Index = 1
		admin.Username = "admin"
		admin.Pwhash = utils.SHA3hash("password")
		admin.AccessToken = "	"
		admin.AccessTokenTimeout = utils.Unix() + 1000000
		admin.Admin = true
		err = admin.Save()
		if err != nil {
			log.Fatal(err)
		}
		err = admin.GenKeys("x")
		if err != nil {
			log.Fatal(err)
		}
		go xlm.GetXLM(admin.StellarWallet.PublicKey)
	}

	if opts.Rescue {
		RescueMode()
		os.Exit(1)
	}

	// rpc.KillCode = "NUKE" // compile time nuclear code
	// run this only when you need to monitor the tellers. Not required for local testing.
	// go opensolar.MonitorTeller(1)
	fmt.Println(`
	  ██████╗  ██████╗███████╗ ███╗   ██╗██╗  ██╗
	 ██╔═══██╗██╔══██╗██╔════╝████╗  ██║╚██╗██╔╝
	 ██║   ██║██████╔╝█████╗  ██╔██╗ ██║ ╚███╔╝
	 ██║   ██║██╔═══╝ ██╔══╝  ██║╚██╗██║ ██╔██╗
	 ╚██████╔╝██║     ███████╗██║ ╚████║██╔╝ ██╗
	  ╚═════╝ ╚═╝     ╚══════╝╚═╝  ╚═══╝╚═╝  ╚═╝
		`)

	rpc.StartServer(port, insecure)
}
