package rpc

import (
	"log"
	"net/http"

	database "github.com/OpenFinancing/openfinancing/database"
	stablecoin "github.com/OpenFinancing/openfinancing/stablecoin"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
)

// this file handles the RPCs necessary for converting a fixed amount of XLM into
// stablecoins. Right now, this is hooked to our account which serves the stablecoin
// but in the future, we can have any provider that is willing to provide this
// We have get requests here simply because they're easy to handle.

func setupStableCoinRPCs() {
	getStableCoin()
}

func getStableCoin() {
	http.HandleFunc("/stablecoin/get", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		log.Println("Calling route")
		if r.URL.Query() == nil || r.URL.Query()["seed"] == nil || r.URL.Query()["amount"] == nil ||
			r.URL.Query()["username"] == nil || r.URL.Query()["pwhash"] == nil ||
			len(r.URL.Query()["pwhash"][0]) != 128 {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		_, err := database.ValidateUser(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		// we need to validate the user and check if its a part of the platform. If not,
		// we don't allow it to exchange xlm for stablecoin.
		receiverSeed := r.URL.Query()["seed"][0]
		amount := r.URL.Query()["amount"][0] // in string
		receiverPubkey, err := wallet.ReturnPubkey(receiverSeed)
		if err != nil {
			log.Println("Error while retrieving pubkey")
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		log.Println("Pubkey: ", receiverPubkey)
		err = stablecoin.Exchange(receiverPubkey, receiverSeed, amount)
		if err != nil {
			log.Println("error while exchanging", err)
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		Send200(w, r)
	})
}