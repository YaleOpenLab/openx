package rpc

import (
	"log"
	"net/http"

	stablecoin "github.com/Varunram/essentials/crypto/stablecoin"
	wallet "github.com/Varunram/essentials/crypto/xlm/wallet"
	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
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
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		user, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["seedpwd"] == nil || r.URL.Query()["amount"] == nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		// we need to validate the user and check if its a part of the platform. If not,
		// we don't allow it to exchange xlm for stablecoin.
		receiverSeed, err := wallet.DecryptSeed(user.StellarWallet.EncryptedSeed, r.URL.Query()["seedpwd"][0])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		amount := r.URL.Query()["amount"][0] // in string
		receiverPubkey, err := wallet.ReturnPubkey(receiverSeed)
		if err != nil {
			log.Println("did not return pubkey", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		err = stablecoin.Exchange(receiverPubkey, receiverSeed, amount)
		if err != nil {
			log.Println("did not exchange for xlm", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

type GetAnchorResponse struct {
	Txhash string // this tx hash is for the sent xlm, not for the received anchorUSD
}

// getAnchorUSD gets anchorUSD from Anchor
func getAnchorUSD() {
	http.HandleFunc("/anchor/get", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		user, err := UserValidateHelper(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		if !user.Kyc || user.Banned {
			// banned  or user without kyc is trying to request stablecoin, don't allow
			log.Println("user who is not verified under kyc / is sanctioned is requesting stablecoin: ", user.Name, user.Pwhash)
			erpc.ResponseHandler(w, erpc.StatusNotAcceptable)
			return
		}

		// there are two ways in which a person can get anchorUSD - wire payments / crypto transfer. This is defined by Anchor
		// and there's nothing we can do to change this.
		if r.URL.Query()["mode"] == nil {
			log.Println("user hasn't specified mode, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		mode := r.URL.Query()["mode"][0]

		switch mode {
		case "wire":
			// wire payments require a set of constraints that have to be stored in the system, put that in here
		case "crypto":
			if r.URL.Query()["seedpwd"] == nil || r.URL.Query()["amount"] == nil {
				log.Println("required params for crypto to anchorUSD transaction not defined, quitting")
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}

			seed, err := wallet.DecryptSeed(user.StellarWallet.EncryptedSeed, r.URL.Query()["seedpwd"][0])
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}

			amount, err := utils.ToFloat(r.URL.Query()["amount"][0]) // amount that the person wants to get. This must be in USD
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}

			txhash, err := stablecoin.GetAnchorUSD(seed, amount)
			if err != nil {
				log.Println("error in fetching stablecoin, quitting")
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}

			var response GetAnchorResponse
			response.Txhash = txhash
			erpc.MarshalSend(w, response)
			// send a tx to anchorUSD's pubkey and hope they send stablecoin back. Since this is not dependent on us, we'd need
			// to wait for the interval that anchorUSD determines in order to be able to proceed further.
		default:
			log.Println("mode not specified for anchorUSD conversion, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
	})
}
