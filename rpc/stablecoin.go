package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	stablecoin "github.com/Varunram/essentials/xlm/stablecoin"
	wallet "github.com/Varunram/essentials/xlm/wallet"
	consts "github.com/YaleOpenLab/openx/consts"
)

// StablecoinRPC is a collection of all stablecoin RPC endpoints and their required params
var StablecoinRPC = map[int][]string{
	1: {"/stablecoin/get", "GET", "seedpwd", "amount"}, // GET
	2: {"/anchor/get", "GET"},                          // GET
}

// setupStableCoinRPCs sets up the endpoints that would be required to interact with
// openx's inhouse stablecoin
func setupStableCoinRPCs() {
	getTestStableCoin()
	getAnchorUSD()
}

// getTestStableCoin gets stablecoin in exchange for xlm
func getTestStableCoin() {
	http.HandleFunc(StablecoinRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		if consts.Mainnet {
			log.Println("test stablecoin not available on mainnet")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		user, err := userValidateHelper(w, r, StablecoinRPC[1][2:], StablecoinRPC[1][1])
		if err != nil {
			return
		}

		receiverSeed, err := wallet.DecryptSeed(user.StellarWallet.EncryptedSeed, r.URL.Query()["seedpwd"][0])
		if err != nil {
			return
		}
		amount, err := utils.ToFloat(r.URL.Query()["amount"][0])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		receiverPubkey, err := wallet.ReturnPubkey(receiverSeed)
		if erpc.Err(w, err, erpc.StatusBadRequest, "did not return pubkey") {
			return
		}

		err = stablecoin.Exchange(receiverPubkey, receiverSeed, amount)
		if err != nil {
			log.Println("did not exchange for xlm", err)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// GetAnchorResponse is a wrapper around the txhash for sent XLM
type GetAnchorResponse struct {
	Txhash string // this tx hash is for the sent xlm, not for the received anchorUSD
}

// getAnchorUSD gets anchorUSD from Anchor
func getAnchorUSD() {
	http.HandleFunc(StablecoinRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		user, err := userValidateHelper(w, r, StablecoinRPC[2][2:], StablecoinRPC[2][1])
		if err != nil {
			return
		}

		if !user.Kyc || user.Banned {
			// banned or user without kyc is trying to request stablecoin, don't allow
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
			if erpc.Err(w, err, erpc.StatusBadRequest) {
				return
			}

			amount, err := utils.ToFloat(r.URL.Query()["amount"][0]) // amount that the person wants to get. This must be in USD
			if erpc.Err(w, err, erpc.StatusBadRequest) {
				return
			}

			txhash, err := stablecoin.GetAnchorUSD(seed, amount)
			if erpc.Err(w, err, erpc.StatusInternalServerError, "error in fetching stablecoin, quitting") {
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
