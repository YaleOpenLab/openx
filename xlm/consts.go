package xlm

import (
	"net/http"

	clients "github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/network"
)

// constants that are imported by other packages

// TestNetClient defines the horizon client to connect to
var TestNetClient = &clients.Client{
	URL:  "https://horizon-testnet.stellar.org",
	HTTP: http.DefaultClient,
}

var Passphrase = network.TestNetworkPassphrase
