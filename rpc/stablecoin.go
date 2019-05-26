package rpc

import (
	"log"
	"net/http"

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
		checkOrigin(w, r)
		user, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["seedpwd"] == nil || r.URL.Query()["amount"] == nil {
			log.Println(err)
			responseHandler(w, StatusBadRequest)
			return
		}
		// we need to validate the user and check if its a part of the platform. If not,
		// we don't allow it to exchange xlm for stablecoin.
		receiverSeed, err := wallet.DecryptSeed(user.EncryptedSeed, r.URL.Query()["seedpwd"][0])
		if err != nil {
			log.Println(err)
			responseHandler(w, StatusBadRequest)
			return
		}
		amount := r.URL.Query()["amount"][0] // in string
		receiverPubkey, err := wallet.ReturnPubkey(receiverSeed)
		if err != nil {
			log.Println("did not return pubkey", err)
			responseHandler(w, StatusBadRequest)
			return
		}
		err = stablecoin.Exchange(receiverPubkey, receiverSeed, amount)
		if err != nil {
			log.Println("did not exchange for xlm", err)
			responseHandler(w, StatusInternalServerError)
			return
		}
		responseHandler(w, StatusOK)
	})
}
