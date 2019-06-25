package xlm

import (
	"net/http"

	horizon "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
)

// constants that are imported by other packages

// TestNetClient defines the horizon client to connect to
var TestNetClient = &horizon.Client{
	HorizonURL: "https://horizon-testnet.stellar.org/",
	HTTP:       http.DefaultClient,
}

var Passphrase = network.TestNetworkPassphrase
