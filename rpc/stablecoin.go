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
	getTestStableCoin()
	getAnchorUSD()
}

// getStableCoin gets stablecoin in exchange for xlm
func getTestStableCoin() {
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

type GetAnchorResponse struct {
	Txhash string // this tx hash is for the sent xlm, not for the received anchorUSD
}

// getAnchorUSD gets anchorUSD from Anchor
func getAnchorUSD() {
	http.HandleFunc("/anchor/get", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		user, err := UserValidateHelper(w, r)
		if err != nil {
			log.Println(err)
			responseHandler(w, StatusUnauthorized)
			return
		}

		if !user.Kyc || user.Banned {
			// banned  or user without kyc is trying to request stablecoin, don't allow
			log.Println("user who is not verified under kyc / is sanctioned is requesting stablecoin: ", user.Name, user.Pwhash)
			responseHandler(w, StatusNotAcceptable)
			return
		}

		// there are two ways in which a person can get anchorUSD - wire payments / crypto transfer. This is defined by Anchor
		// and there's nothing we can do to change this.
		if r.URL.Query()["mode"] == nil {
			log.Println("user hasn't specified mode, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		mode := r.URL.Query()["mode"][0]

		switch mode {
		case "wire":
			// wire payments require a set of constraints that have to be stored in the system, put that in here
		case "crypto":
			if r.URL.Query()["seedpwd"] == nil || r.URL.Query()["amount"] == nil {
				log.Println("required params for crypto to anchorUSD transaction not defined, quitting")
				responseHandler(w, StatusBadRequest)
				return
			}

			seed, err := wallet.DecryptSeed(user.EncryptedSeed, r.URL.Query()["seedpwd"][0])
			if err != nil {
				log.Println(err)
				responseHandler(w, StatusBadRequest)
				return
			}

			amount := r.URL.Query()["amount"][0] // amount that the person wants to get. This must be in USD

			txhash, err := stablecoin.GetAnchorUSD(seed, amount)
			if err != nil {
				log.Println("error in fetching stablecoin, quitting")
				responseHandler(w, StatusInternalServerError)
				return
			}

			var response GetAnchorResponse
			response.Txhash = txhash
			MarshalSend(w, response)
			// send a tx to anchorUSD's pubkey and hope they send stablecoin back. Since this is not dependent on us, we'd need
			// to wait for the interval that anchorUSD determines in order to be able to proceed further.
		default:
			log.Println("mode not specified for anchorUSD conversion, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}
	})
}
