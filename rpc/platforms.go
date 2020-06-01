package rpc

import (
	"log"
	"net/http"

	"github.com/YaleOpenLab/openx/notif"

	"github.com/Varunram/essentials/email"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	"github.com/pkg/errors"
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
	pfSendEmail()
	pfConfirmUser()
}

// PlatformRPC is a map that stores all handlers related to the platform
var PlatformRPC = map[int][]string{
	0: {"/platform/getconsts"},                                          // GET
	1: {"/platform/user/retrieve", "key"},                               // GET
	2: {"/platform/user/validate", "username", "token"},                 // GET
	3: {"/platform/user/new", "username", "pwhash", "seedpwd", "email"}, // GET
	4: {"/platform/user/collision", "username"},                         // GET
	5: {"/platforms/all"},                                               // GET NOAUTH
	6: {"/platform/email", "body", "to"},                                // POST
	7: {"/platform/user/confirm", "username", "pwhash", "code"},         // GET
}

// mainnetRPC is an RPC that reutrns 0 if openx is running on mainnet, 1 if running on testnet
func mainnetRPC() {
	http.HandleFunc("/mainnet", func(w http.ResponseWriter, r *http.Request) {
		// set a single byte response for mainnet / testnet
		// mainnet is 0, testnet is 1
		mainnet := []byte{1}
		testnet := []byte{0}
		if consts.Mainnet {
			w.Write(mainnet)
		} else {
			w.Write(testnet)
		}
		return
	})
}

func authPlatform(w http.ResponseWriter, r *http.Request) error {
	var code string
	if r.Method == "GET" {
		if r.URL.Query()["code"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return errors.New("required params code and name not found in request")
		}

		code = r.URL.Query()["code"][0]
	} else if r.Method == "POST" {
		log.Println("platform post call")
		err := r.ParseForm()
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return errors.New("required params code and name not found in request")
		}

		code = r.FormValue("code")
		if code == "" {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return errors.New("required param code not found in request")
		}
	}

	log.Println("Platform's code: ", code)

	platforms, err := database.RetrieveAllPlatforms()
	if erpc.Err(w, err, erpc.StatusBadRequest) {
		return err
	}

	for _, platform := range platforms {
		log.Println(platform.Code, platform.Name, platform.Code == code, utils.Unix(), platform.Timeout)
		if platform.Code == code && utils.Unix() < platform.Timeout {
			return nil
		}
	}

	erpc.ResponseHandler(w, erpc.StatusBadRequest)
	return errors.New("could not authenticate platform, quitting")
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
	http.HandleFunc(PlatformRPC[0][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = authPlatform(w, r)
		if err != nil {
			log.Println(err)
			return
		}

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
	})
}

// pfGetUser retrieves a user from openx's database and returns it to the requesting platform
func pfGetUser() {
	http.HandleFunc(PlatformRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = authPlatform(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["key"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		keyInt, err := utils.ToInt(r.URL.Query()["key"][0])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		user, err := database.RetrieveUser(keyInt)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		erpc.MarshalSend(w, user)
	})
}

// pfValidateUser takes in a username and token and returns the user struct
func pfValidateUser() {
	http.HandleFunc(PlatformRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		log.Println("external platform requests validation")
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = authPlatform(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		if r.URL.Query()["username"] == nil || r.URL.Query()["token"] == nil {
			log.Println("token missing")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		name := r.URL.Query()["username"][0]
		token := r.URL.Query()["token"][0]

		user, err := database.ValidateAccessToken(name, token)
		if erpc.Err(w, err, erpc.StatusBadRequest, "error while validating user") {
			return
		}

		erpc.MarshalSend(w, user)
	})
}

// pfNewUser creates a new user
func pfNewUser() {
	http.HandleFunc(PlatformRPC[3][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = authPlatform(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		if r.URL.Query()["username"] == nil || r.URL.Query()["pwhash"] == nil ||
			r.URL.Query()["seedpwd"] == nil || r.URL.Query()["email"] == nil {
			log.Println("required params missing")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		name := r.URL.Query()["username"][0]
		pwhash := r.URL.Query()["pwhash"][0]
		seedpwd := r.URL.Query()["seedpwd"][0]
		email := r.URL.Query()["email"][0]

		user, err := database.NewUser(name, pwhash, seedpwd, email)
		if erpc.Err(w, err, erpc.StatusBadRequest, "error while creating a new user") {
			return
		}

		err = notif.SendUserConfEmail(email, user.ConfToken)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "could not send user conf token") {
			return
		}

		erpc.MarshalSend(w, user)
	})
}

// pfCollisionCheck checks for username collision
func pfCollisionCheck() {
	http.HandleFunc(PlatformRPC[4][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = authPlatform(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["username"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		name := r.URL.Query()["username"][0]
		noCollision := []byte{0}
		collision := []byte{1}
		_, err = database.CheckUsernameCollision(name)
		if err != nil {
			w.Write(collision)
		} else {
			w.Write(noCollision)
		}
	})
}

// retrieveAllPlatformNames retrieves all platforms from the database
func retrieveAllPlatformNames() {
	http.HandleFunc(PlatformRPC[5][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		platforms, err := database.RetrieveAllPlatforms()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		var platformNames []string

		for _, platform := range platforms {
			platformNames = append(platformNames, platform.Name)
		}

		erpc.MarshalSend(w, platformNames)
	})
}

// pfSendEmail sends an email on behalf of another platform from the openx platform
func pfSendEmail() {
	http.HandleFunc(PlatformRPC[6][0], func(w http.ResponseWriter, r *http.Request) {
		log.Println("external platform requests email from openx")
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = authPlatform(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("AUTHENTICATED PLATFORM")
		if r.FormValue("body") == "" || r.FormValue("to") == "" {
			log.Println("reqd param body or code not found")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		body := r.FormValue("body")
		to := r.FormValue("to")

		err = email.SendMail(body, to)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// pfConfirmUser confirms a user's registration on an openx platform
func pfConfirmUser() {
	http.HandleFunc(PlatformRPC[7][0], func(w http.ResponseWriter, r *http.Request) {
		log.Println("external platform requests validation")
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = authPlatform(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		if r.URL.Query()["username"] == nil || r.URL.Query()["pwhash"] == nil || r.URL.Query()["code"] == nil {
			log.Println("username / pwhash missing")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		name := r.URL.Query()["username"][0]
		pwhash := r.URL.Query()["pwhash"][0]
		confcode := r.URL.Query()["confcode"][0]

		user, err := database.ValidatePwhashReg(name, pwhash)
		if erpc.Err(w, err, erpc.StatusBadRequest, "error while validating user") {
			return
		}

		if confcode != user.ConfToken {
			log.Println("provided code does not match with required code")
		}

		user.Conf = true
		user.ConfToken = ""
		err = user.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError, "could not save user") {
			return
		}

		erpc.MarshalSend(w, user)
	})
}
