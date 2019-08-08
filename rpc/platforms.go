package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	consts "github.com/YaleOpenLab/openx/consts"
)

func setupPlatformRoutes() {
	getConsts()
}

type OpensolarConstReturn struct {
	PlatformPublicKey   string
	PlatformSeed        string
	PlatformEmail       string
	PlatformEmailPass   string
	StablecoinCode      string
	StablecoinPublicKey string
	AnchorUSDCode       string
	AnchorUSDAddress    string
	AnchorUSDTrustLimit float64
	AnchorAPI           string
	Mainnet             bool
	DbDir               string
}

func getConsts() {
	http.HandleFunc("/platform/getconsts", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if r.URL.Query()["code"] == nil {
			log.Println("code missing")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		code := r.URL.Query()["code"][0]

		if code == "OPENSOLARTEST" {
			log.Println("authenticated opensolar platform, sending consts")
			var x OpensolarConstReturn
			x.PlatformPublicKey = consts.PlatformPublicKey
			x.PlatformSeed = consts.PlatformSeed
			x.PlatformEmail = consts.PlatformEmail
			x.PlatformEmailPass = consts.PlatformEmailPass
			x.StablecoinCode = consts.StablecoinCode
			x.StablecoinPublicKey = consts.StablecoinPublicKey
			x.AnchorUSDCode = consts.AnchorUSDCode
			x.AnchorUSDAddress = consts.AnchorUSDAddress
			x.AnchorUSDTrustLimit = consts.AnchorUSDTrustLimit
			x.AnchorAPI = consts.AnchorAPI
			x.Mainnet = consts.Mainnet
			x.DbDir = consts.DbDir
			erpc.MarshalSend(w, x)
			return
		}
	})
}
