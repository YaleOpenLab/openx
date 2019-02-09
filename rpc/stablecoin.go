package rpc

import (
	// "log"
	"net/http"

	database "github.com/YaleOpenLab/openx/database"
	stablecoin "github.com/YaleOpenLab/openx/stablecoin"
	wallet "github.com/YaleOpenLab/openx/wallet"
)

// this file handles the RPCs necessary for converting a fixed amount of XLM into
// stablecoins. Right now, this is hooked to our account which serves the stablecoin
// but in the future, we can have any provider that is willing to provide this
// We have get requests here simply because they're easy to handle.

func setupStableCoinRPCs() {
	getStableCoin()
}

// getStableCoin gets stablecoin in exchange for xlm
func getStableCoin() {
	http.HandleFunc("/stablecoin/get", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		_, err := database.ValidateUser(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
		if err != nil || r.URL.Query()["seed"] == nil || r.URL.Query()["amount"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		// we need to validate the user and check if its a part of the platform. If not,
		// we don't allow it to exchange xlm for stablecoin.
		receiverSeed := r.URL.Query()["seed"][0]
		amount := r.URL.Query()["amount"][0] // in string
		receiverPubkey, err := wallet.ReturnPubkey(receiverSeed)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		err = stablecoin.Exchange(receiverPubkey, receiverSeed, amount)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}
