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
	Passphrase = network.TestNetworkPassphrase
	// Mainnet is a bool which decides which chain to connect to
	Mainnet bool
	// TestNetClient defines the horizon client to connect to
	TestNetClient *horizon.Client
)

// RefillAmount defines the default stellar refill amount
var RefillAmount float64

func SetConsts(amount float64, mainnet bool) {
	RefillAmount = amount
	Mainnet = mainnet
	if mainnet {
		log.Println("Pointing horizon to mainnet")
		TestNetClient = &horizon.Client{
			HorizonURL: "https://horizon.stellar.org/", // switch to mainnet horizon
			HTTP:       http.DefaultClient,
		}
	} else {
		TestNetClient = &horizon.Client{
			HorizonURL: "https://horizon-testnet.stellar.org/",
			HTTP:       http.DefaultClient,
		}
	}
}
