package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
)

// this file has routes that are to be exclusively used by external platforms in order to call
// data that is exclusive to openx. These platforms can request a code to be generated and they
// can call openx's endpoints

// setupPlatformRoutes sets up routes that are related with external third party platforms
// which need information from openx to operate
func setupPlatformRoutes() {
	mainnetRPC()
	pfGetConsts()
	pfGetUser()
	pfValidateUser()
	pfNewUser()
	pfCollisionCheck()
	retrieveAllPlatformNames()
}

// mainnetRPC is an RPC that reutrns 0 if openx is running on mainnet, 1 if running on testnet
func mainnetRPC() {
	http.HandleFunc("/mainnet", func(w http.ResponseWriter, r *http.Request) {
		// set a single byte response for mainnet / testnet
		// mainnet is 0, testnet is 1
		mainnet := []byte{0}
		testnet := []byte{1}
		if consts.Mainnet {
			w.Write(mainnet)
		} else {
			w.Write(testnet)
		}
		return
	})
}

// OpensolarConstReturn is a struct that can be used to export consts from openx
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

// pfGetConsts is an RPC that returns running constants to platforms which might need this information
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

// pfGetUser retrieves a user from openx's database and returns it to the requesting platform
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

// pfValidateUser validates a given user and returns the user struct
func pfValidateUser() {
	http.HandleFunc("/platform/user/validate", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if r.URL.Query()["code"] == nil {
			log.Println("code missing")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		code := r.URL.Query()["code"][0]

		if code == "OPENSOLARTEST" {

			if r.URL.Query()["username"] == nil || r.URL.Query()["token"] == nil {
				log.Println("token missing")
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
			}

			log.Println("QUERY PARAMS: ", r.URL.Query())
			name := r.URL.Query()["username"][0]
			token := r.URL.Query()["token"][0]

			user, err := database.ValidateAccessToken(name, token)
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}
			erpc.MarshalSend(w, user)
		}
	})
}

// pfNewUser creates a new user
func pfNewUser() {
	http.HandleFunc("/platform/user/new", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if r.URL.Query()["code"] == nil {
			log.Println("code missing")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		code := r.URL.Query()["code"][0]
		if code == "OPENSOLARTEST" {

			if r.URL.Query()["username"] == nil || r.URL.Query()["pwhash"] == nil || r.URL.Query()["seedpwd"] == nil ||
				r.URL.Query()["realname"] == nil {
				log.Println("code missing")
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
			}

			name := r.URL.Query()["username"][0]
			pwhash := r.URL.Query()["pwhash"][0]
			seedpwd := r.URL.Query()["seedpwd"][0]
			realname := r.URL.Query()["realname"][0]

			user, err := database.NewUser(name, pwhash, seedpwd, realname)
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}
			erpc.MarshalSend(w, user)
		}
	})
}

// pfCollisionCheck checks for username collision
func pfCollisionCheck() {
	http.HandleFunc("/platform/user/collision", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if r.URL.Query()["code"] == nil {
			log.Println("code missing")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		code := r.URL.Query()["code"][0]

		if code == "OPENSOLARTEST" {
			if r.URL.Query()["username"] == nil {
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}
			name := r.URL.Query()["username"][0]
			noCollision := []byte{0}
			collision := []byte{1}
			_, err := database.CheckUsernameCollision(name)
			if err != nil {
				w.Write(collision)
			} else {
				w.Write(noCollision)
			}
		}
	})
}

// retrieveAllPlatformNames retrieves all platforms from the database
func retrieveAllPlatformNames() {
	http.HandleFunc("/platforms/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		platforms, err := database.RetrieveAllPlatforms()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var platformNames []string

		for _, platform := range platforms {
			platformNames = append(platformNames, platform.Name)
		}

		erpc.MarshalSend(w, platformNames)
	})
}
