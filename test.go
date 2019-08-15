package main

import (
	"fmt"
	// "github.com/pkg/errors"
	"log"
	"os"

	consts "github.com/YaleOpenLab/openx/consts"
	loader "github.com/YaleOpenLab/openx/loader"
	database "github.com/YaleOpenLab/openx/database"
	// ipfs "github.com/YaleOpenLab/openx/ipfs"
	opensolar "github.com/YaleOpenLab/opensolar/consts"
	rpc "github.com/YaleOpenLab/openx/rpc"
	// scan "github.com/YaleOpenLab/openx/scan"
	// oracle "github.com/YaleOpenLab/openx/oracle"
	// algorand "github.com/YaleOpenLab/openx/chains/algorand"
	// stablecoin "github.com/YaleOpenLab/openx/chains/stablecoin"
	// utils "github.com/Varunram/essentials/utils"
	// scan "github.com/YaleOpenLab/openx/scan"
	// wallet "github.com/YaleOpenLab/openx/wallet"
	// xlm "github.com/YaleOpenLab/openx/xlm"
	// assets "github.com/YaleOpenLab/openx/assets"
	flags "github.com/jessevdk/go-flags"
	// "github.com/spf13/viper"
)

// the backend server powering the openx platform of platforms

var opts struct {
	Insecure  bool `short:"i" description:"Start the API using http. Not recommended"`
	Port      int  `short:"p" description:"The port on which the server runs on. Default: HTTPS/8080"`
	Simulate  bool `short:"t" description:"Simulate the test database with demo values (last updated: April 2019)"`
	Mainnet   bool `short:"m" description:"Switch mainnet mode on"`
	Trustline bool `short:"x" description:"create trustlines from platform seed to anchorUSD"`
	Rescue    bool `short:"r" description:"start rescue mode"`
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

func main() {
	var err error
	insecure, port, err := ParseConfig(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	if consts.Mainnet {
		err = loader.Mainnet()
		if err != nil {
			log.Fatal(err)
		}
		opensolar.SetMnConsts()
		if opts.Trustline {
			loader.StablecoinTrust()
		}
	} else {
		err = loader.Testnet()
		if err != nil {
			log.Fatal(err)
		}
		opensolar.SetTnConsts()
	}

	if opts.Rescue {
		RescueMode()
		os.Exit(1)
	}

	user, err := database.RetrieveUser(1)
	if err != nil {
		log.Fatal(err)
	}
	user.Admin = true
	err = user.Save()
	if err != nil {
		log.Fatal(err)
	}
	// rpc.KillCode = "NUKE" // compile time nuclear code
	// run this only when you need to monitor the tellers. Not required for local testing.
	// go opensolar.MonitorTeller(1)
	fmt.Println(`
		██████╗ ██████╗ ███████╗███╗   ██╗██╗  ██╗
	 ██╔═══██╗██╔══██╗██╔════╝████╗  ██║╚██╗██╔╝
	 ██║   ██║██████╔╝█████╗  ██╔██╗ ██║ ╚███╔╝
	 ██║   ██║██╔═══╝ ██╔══╝  ██║╚██╗██║ ██╔██╗
	 ╚██████╔╝██║     ███████╗██║ ╚████║██╔╝ ██╗
	  ╚═════╝ ╚═╝     ╚══════╝╚═╝  ╚═══╝╚═╝  ╚═╝
		`)

	rpc.StartServer(port, insecure)
}
