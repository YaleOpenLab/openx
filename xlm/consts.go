package xlm

import (
	"net/http"

	clients "github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/network"
)

// TestNetClient defines the horizon client to connect to
var TestNetClient = &clients.Client{
	// URL: "http://35.192.122.229:8080",
	URL:  "https://horizon-testnet.stellar.org",
	HTTP: http.DefaultClient,
}

var Passphrase = network.TestNetworkPassphrase
