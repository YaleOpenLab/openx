package main

import (
	"log"
	"os"

	consts "github.com/OpenFinancing/openfinancing/consts"
	database "github.com/OpenFinancing/openfinancing/database"
	// ipfs "github.com/OpenFinancing/openfinancing/ipfs"
	platform "github.com/OpenFinancing/openfinancing/platforms"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	rpc "github.com/OpenFinancing/openfinancing/rpc"
	// scan "github.com/OpenFinancing/openfinancing/scan"
	stablecoin "github.com/OpenFinancing/openfinancing/stablecoin"
	utils "github.com/OpenFinancing/openfinancing/utils"
	// wallet "github.com/OpenFinancing/openfinancing/wallet"
	// xlm "github.com/OpenFinancing/openfinancing/xlm"
	flags "github.com/jessevdk/go-flags"
)

// the server powering the openfinancing platform. There are two clients that can be used
// with the backned - ofcli and emulator
// refer https://github.com/stellar/go/blob/master/build/main_test.go in case the stellar
// go SDK docs are insufficient.
var opts struct {
	Port int `short:"p" description:"The port on which the server runs on"`
}

func ParseConfig(args []string) (string, error) {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return "", err
	}
	port := utils.ItoS(consts.DefaultRpcPort)
	if opts.Port != 0 {
		port = utils.ItoS(opts.Port)
		log.Println("Starting RPC Server on Port: ", opts.Port)
	}
	return port, nil
}

func StartPlatform() (string, string, error) {
	var publicKey string
	var seed string
	database.CreateHomeDir()

	// InsertDummyData()
	allContracts, err := solar.RetrieveAllProjects()
	if err != nil {
		log.Println("Error retrieving all projects from the database")
		return publicKey, seed, err
	}

	if len(allContracts) == 0 {
		log.Println("Populating database with test values")
		err = InsertDummyData()
		if err != nil {
			return publicKey, seed, err
		}
	}
	publicKey, seed, err = platform.InitializePlatform()
	return publicKey, seed, err
}

func main() {
	var err error
	port, err := ParseConfig(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	consts.PlatformPublicKey, consts.PlatformSeed, err = StartPlatform()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("PLATFORM SEED IS: %s\n PLATFORM PUBLIC KEY IS: %s\n", consts.PlatformSeed, consts.PlatformPublicKey)

	err = stablecoin.InitStableCoin() // start the stablecoin daemon
	if err != nil {
		log.Fatal(err)
	}

	rpc.StartServer(port)
}
