package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
)

func setupPlatformRoutes() {
	pfGetConsts()
	pfGetUser()
	pfValidateUser()
	pfNewUser()
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

func pfGetConsts() {
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

func pfGetUser() {
	http.HandleFunc("/platform/user/retrieve", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if r.URL.Query()["code"] == nil || r.URL.Query()["key"] == nil {
			log.Println("code missing")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		code := r.URL.Query()["code"][0]

		if code == "OPENSOLARTEST" {
			keyInt, err := utils.ToInt(r.URL.Query()["key"][0])
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}
			user, err := database.RetrieveUser(keyInt)
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}
			erpc.MarshalSend(w, user)
		}
	})
}

func pfValidateUser() {
	http.HandleFunc("/platform/user/validate", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if r.URL.Query()["name"] == nil || r.URL.Query()["pwhash"] == nil {
			log.Println("code missing")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		name := r.URL.Query()["name"][0]
		pwhash := r.URL.Query()["pwhash"][0]

		user, err := database.ValidateUser(name, pwhash)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		erpc.MarshalSend(w, user)
	})
}

func pfNewUser() {
	http.HandleFunc("/platform/user/new", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if r.URL.Query()["name"] == nil || r.URL.Query()["pwhash"] == nil || r.URL.Query()["seedpwd"] == nil ||
			r.URL.Query()["realname"] == nil {
			log.Println("code missing")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		name := r.URL.Query()["name"][0]
		pwhash := r.URL.Query()["pwhash"][0]
		seedpwd := r.URL.Query()["seedpwd"][0]
		realname := r.URL.Query()["realname"][0]

		user, err := database.NewUser(name, pwhash, seedpwd, realname)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		log.Println("USER+", user)
		erpc.MarshalSend(w, user)
	})
}
