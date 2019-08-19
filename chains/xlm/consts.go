package xlm

import (
	"log"
	"net/http"

	horizon "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
)

// constants that are imported by other packages

var (
	// Passphrase defines the stellar network passphrase
	Passphrase string
	// Mainnet is a bool which decides which chain to connect to
	Mainnet bool
	// TestNetClient defines the horizon client to connect to
	TestNetClient *horizon.Client
)

// RefillAmount defines the default stellar refill amount
var RefillAmount float64

// SetConsts XLM consts
func SetConsts(amount float64, mainnet bool) {
	RefillAmount = amount
	Mainnet = mainnet
	log.Println("SETTING MAINNET TO: ", mainnet)
	if mainnet {
		Passphrase = network.PublicNetworkPassphrase
		log.Println("Pointing horizon to mainnet")
		TestNetClient = &horizon.Client{
			HorizonURL: "https://horizon.stellar.org/", // switch to mainnet horizon
			HTTP:       http.DefaultClient,
		}
	} else {
		log.Println("Pointing horizon to testnet")
		Passphrase = network.TestNetworkPassphrase
		TestNetClient = &horizon.Client{
			HorizonURL: "https://horizon-testnet.stellar.org/",
			HTTP:       http.DefaultClient,
		}
	}
}
