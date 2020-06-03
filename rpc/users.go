package rpc

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	"github.com/pkg/errors"

	aes "github.com/Varunram/essentials/aes"
	ipfs "github.com/Varunram/essentials/ipfs"
	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/Varunram/essentials/xlm"
	assets "github.com/Varunram/essentials/xlm/assets"
	wallet "github.com/Varunram/essentials/xlm/wallet"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	notif "github.com/YaleOpenLab/openx/notif"
	recovery "github.com/bithyve/research/sss"
)

// UserRPC is a collection of all user RPC endpoints and their required params
var UserRPC = map[int][]string{
	0:  {"/token"},                                                                         // POST
	1:  {"/user/validate", "GET"},                                                          // GET
	2:  {"/user/balances", "GET"},                                                          // GET
	3:  {"/user/balance/xlm", "GET"},                                                       // GET
	4:  {"/user/balance/asset", "GET", "asset"},                                            // GET
	5:  {"/ipfs/getdata", "GET", "hash"},                                                   // GET
	6:  {"/user/kyc", "GET", "userIndex"},                                                  // GET
	7:  {"/user/sendxlm", "GET", "destination", "amount", "seedpwd"},                       // GET
	8:  {"/user/notkycview", "GET"},                                                        // GET
	9:  {"/user/kycview", "GET"},                                                           // GET
	10: {"/user/askxlm", "GET"},                                                            // GET
	11: {"/user/trustasset", "GET", "assetCode", "assetIssuer", "limit", "seedpwd"},        // GET
	12: {"/upload", "POST"},                                                                // POST
	13: {"/platformemail", "GET"},                                                          // GET
	17: {"/user/increasetrustlimit", "GET", "trust", "seedpwd"},                            // GET
	19: {"/user/sendrecovery", "GET", "email1", "email2", "email3"},                        // GET
	20: {"/user/seedrecovery", "GET", "secret1", "secret2"},                                // GET
	21: {"/user/newsecrets", "GET", "seedpwd", "email1", "email2", "email3"},               // GET
	22: {"/user/resetpwd", "GET", "seedpwd", "email"},                                      // GET
	23: {"/user/pwdreset", "GET", "pwhash", "email", "verificationCode"},                   // GET
	24: {"/user/sweep", "GET", "seedpwd", "destination"},                                   // GET
	25: {"/user/sweepasset", "GET", "seedpwd", "destination", "assetName", "issuerPubkey"}, // GET
	26: {"/user/verifykyc", "GET", "selfie"},                                               // GET
	27: {"/user/giverating", "GET", "feedback", "userIndex"},                               // GET
	28: {"/user/2fa/generate", "GET"},                                                      // GET
	29: {"/user/2fa/authenticate", "GET", "password"},                                      // GET
	31: {"/user/reputation", "GET", "reputation"},                                          // GET
	32: {"/user/addseed", "GET", "encryptedseed", "seedpwd", "pubkey"},                     // GET
	33: {"/user/latestblockhash", "GET"},                                                   // GET
	34: {"/ipfs/putdata", "POST", "data"},                                                  // POST
	35: {"/user/tc", "POST"},                                                               // POST
	36: {"/user/progress", "POST", "progress"},                                             // POST
	37: {"/user/update", "POST"},                                                           // POST
	38: {"/user/tellerfile", "GET"},                                                        // GET
	39: {"/user/logout", "POST"},                                                           // POST
	40: {"/user/verify", "POST"},                                                           // POST
	41: {"/user/unverify", "POST"},                                                         // POST

	30: {"/user/anchorusd/kyc", "GET", "name", "bdaymonth", "bdayday", "bdayyear", "taxcountry", // GET
		"taxid", "addrstreet", "addrcity", "addrpostal", "addrregion", "addrcountry", "addrphone", "primaryphone", "gender"},
}

// setupUserRpcs sets up user related RPCs
func setupUserRpcs() {
	validateUser()
	getBalances()
	getXLMBalance()
	getAssetBalance()
	getIpfsData()
	putIpfsData()
	authKyc()
	sendXLM()
	notKycView()
	kycView()
	askForCoins()
	trustAsset()
	uploadFile()
	platformEmail()
	increaseTrustLimit()
	sendSecrets()
	mergeSecrets()
	generateNewSecrets()
	generateResetPwdCode()
	resetPassword()
	sweepFunds()
	sweepAsset()
	validateKYC()
	giveStarRating()
	new2fa()
	auth2fa()
	changeReputation()
	addAnchorKYCInfo()
	importSeed()
	genAccessToken()
	getLatestBlockHash()
	acceptTc()
	updateProgress()
	updateUser()
	downloadTeller()
	logout()
	verify()
	unverify()

	// sendTellerShutdownEmail()
	// sendTellerFailedPaybackEmail()
	// tellerPing()
	// addContractHash()
}

func checkReqdParams(w http.ResponseWriter, r *http.Request, options []string, method string) error {
	if method == "GET" {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return err
		}

		if r.URL.Query() == nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return errors.New("url query can't be empty")
		}

		options = append(options, "username", "token") // default for all endpoints

		for _, option := range options {
			if r.URL.Query()[option] == nil {
				erpc.ResponseHandler(w, erpc.StatusUnauthorized)
				return errors.New("required param: " + option + " not specified, quitting")
			}
		}

		if len(r.URL.Query()["token"][0]) != consts.AccessTokenLength {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return errors.New("token length not 32, quitting")
		}

	} else if method == "POST" {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return err
		}

		err = r.ParseForm()
		if erpc.Err(w, err, erpc.StatusUnauthorized) {
			return err
		}

		options = append(options, "username", "token") // default for all endpoints

		for _, option := range options {
			if r.FormValue(option) == "" {
				erpc.ResponseHandler(w, erpc.StatusUnauthorized)
				return errors.New("required param: " + option + " not specified, quitting")
			}
		}

		if len(r.FormValue("token")) != consts.AccessTokenLength {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return errors.New("token length not 32, quitting")
		}

	} else {
		erpc.ResponseHandler(w, erpc.StatusBadRequest)
		return errors.New("invalid method (not GET/POST)")
	}
	return nil
}

// userValidateHelper is a helper that validates the username, pwhash and if passed, seedpwd along with other
// params that are required by each endpoint
func userValidateHelper(w http.ResponseWriter, r *http.Request, options []string, method string) (database.User, error) {
	var prepUser database.User
	var err error
	// need to pass the pwhash param here
	err = checkReqdParams(w, r, options, method)
	if err != nil {
		log.Println("error while checking required params: ", err)
		return prepUser, errors.New("url query can't be empty")
	}

	if method == "GET" {
		if r.URL.Query()["seedpwd"] != nil {
			// check seed pwhash before decryption
			prepUser, err = database.ValidateSeedpwdAuthToken(r.URL.Query()["username"][0], r.URL.Query()["token"][0], r.URL.Query()["seedpwd"][0])
		} else {
			// no seedpwhash, normal call
			prepUser, err = database.ValidateAccessToken(r.URL.Query()["username"][0], r.URL.Query()["token"][0])
		}
	} else if method == "POST" {
		if r.FormValue("seedpwd") != "" && r.FormValue("oldseedpwd") == "" {
			log.Println("validating seedpwd of user")
			prepUser, err = database.ValidateSeedpwdAuthToken(r.FormValue("username"), r.FormValue("token"), r.FormValue("seedpwd"))
		} else {
			prepUser, err = database.ValidateAccessToken(r.FormValue("username"), r.FormValue("token"))
		}
	} else {
		return prepUser, errors.New("invalid method (not GET/POST)")
	}
	// catch the error from the relevant error call
	if erpc.Err(w, err, erpc.StatusUnauthorized, "error while validating user") {
		return prepUser, err
	}

	log.Println("successfully validated: ", prepUser.Name)
	return prepUser, nil
}

// GenAccessTokenReturn is the struct defined for returning access tokens
type GenAccessTokenReturn struct {
	Token string
}

func genAccessToken() {
	http.HandleFunc(UserRPC[0][0], func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		username := r.FormValue("username")
		pwhash := r.FormValue("pwhash")

		if username == "" || pwhash == "" {
			log.Println("required params username or pwhash not found, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		log.Println("username: ", username, " requesting a new access token")
		user, err := database.ValidatePwhash(username, pwhash)
		if erpc.Err(w, err, erpc.StatusUnauthorized) {
			return
		}

		token, err := user.GenAccessToken()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		var x GenAccessTokenReturn
		x.Token = token
		erpc.MarshalSend(w, x)
	})
}

// validateUser validates a user and returns whether the user is an investor or recipient on the opensolar platform
func validateUser() {
	http.HandleFunc(UserRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[1][2:], UserRPC[1][1])
		if err != nil {
			return
		}

		erpc.MarshalSend(w, prepUser)
	})
}

// getBalances returns a list of all balances (assets and XLM) held by the user
func getBalances() {
	http.HandleFunc(UserRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[2][2:], UserRPC[2][1])
		if err != nil {
			return
		}

		pubkey := prepUser.StellarWallet.PublicKey
		balances, err := xlm.GetAllBalances(pubkey)
		if erpc.Err(w, err, erpc.StatusNotFound) {
			return
		}

		erpc.MarshalSend(w, balances)
	})
}

// getXLMBalance gets the XLM balance of the user's primary XLM account
func getXLMBalance() {
	http.HandleFunc(UserRPC[3][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[3][2:], UserRPC[3][1])
		if err != nil {
			return
		}

		pubkey := prepUser.StellarWallet.PublicKey
		balance := xlm.GetNativeBalance(pubkey)
		erpc.MarshalSend(w, balance)
	})
}

// getAssetBalance gets the balance of a specific asset on Stellar
func getAssetBalance() {
	http.HandleFunc(UserRPC[4][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[4][2:], UserRPC[4][1])
		if err != nil {
			return
		}

		asset := r.URL.Query()["asset"][0]
		pubkey := prepUser.StellarWallet.PublicKey

		balance := xlm.GetAssetBalance(pubkey, asset)
		erpc.MarshalSend(w, balance)
	})
}

// getIpfsData gets data from ipfs
func getIpfsData() {
	http.HandleFunc(UserRPC[5][0], func(w http.ResponseWriter, r *http.Request) {
		_, err := userValidateHelper(w, r, UserRPC[5][2:], UserRPC[5][1])
		if err != nil {
			return
		}

		hashString := r.URL.Query()["hash"][0]
		data, err := ipfs.GetString(hashString)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, data)
	})
}

// putIpfsData stores data in ipfs
func putIpfsData() {
	http.HandleFunc(UserRPC[34][0], func(w http.ResponseWriter, r *http.Request) {
		_, err := userValidateHelper(w, r, UserRPC[34][2:], UserRPC[34][1])
		if err != nil {
			return
		}

		data := []byte(r.FormValue("data"))
		hash, err := ipfs.AddBytes([]byte(data))
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		_, err = ipfs.AddBytes(data)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, hash)
	})
}

// authKyc authenticates a user for KYC services
func authKyc() {
	http.HandleFunc(UserRPC[6][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[6][2:], UserRPC[6][1])
		if err != nil {
			return
		}

		uInput, err := utils.ToInt(r.URL.Query()["userIndex"][0])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		err = prepUser.Authorize(uInput)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// sendXLM sends a given amount of XLM to the destination address specified.
func sendXLM() {
	http.HandleFunc(UserRPC[7][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[7][2:], UserRPC[7][1])
		if err != nil {
			return
		}

		destination := r.URL.Query()["destination"][0]
		seedpwd := r.URL.Query()["seedpwd"][0]

		amount, err := utils.ToFloat(r.URL.Query()["amount"][0])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		seed, err := wallet.DecryptSeed(prepUser.StellarWallet.EncryptedSeed, seedpwd)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		var memo string
		if r.URL.Query()["memo"] != nil {
			memo = r.URL.Query()["memo"][0]
		}

		_, txhash, err := xlm.SendXLM(destination, amount, seed, memo)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		erpc.MarshalSend(w, txhash)
	})
}

// notKycView returns a list of all the users who have not yet been verified through KYC. Can be
// called only by KYC Inspectors
func notKycView() {
	http.HandleFunc(UserRPC[8][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[8][2:], UserRPC[8][1])
		if err != nil {
			return
		}

		if !prepUser.Inspector && !prepUser.Admin {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		users, err := database.RetrieveAllUsersWithoutKyc()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, users)
	})
}

// kycView returns a list of all the users who have been KYC verified. Can be called
// only by KYC Inspectors
func kycView() {
	http.HandleFunc(UserRPC[9][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[9][2:], UserRPC[9][1])
		if err != nil {
			return
		}

		if !prepUser.Inspector && !prepUser.Admin {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		users, err := database.RetrieveAllUsersWithKyc()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, users)
	})
}

// askForCoins asks for coins from the Stellar testnet faucet Available only on Stellar testnet
func askForCoins() {
	http.HandleFunc(UserRPC[10][0], func(w http.ResponseWriter, r *http.Request) {
		if consts.Mainnet {
			log.Println("Openx is in mainnet mode, can't ask for coins")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		prepUser, err := userValidateHelper(w, r, UserRPC[10][2:], UserRPC[10][1])
		if err != nil {
			return
		}

		err = xlm.GetXLM(prepUser.StellarWallet.PublicKey)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// trustAsset creates a trustline for the given limit with a remote peer for receiving assets.
func trustAsset() {
	http.HandleFunc(UserRPC[11][0], func(w http.ResponseWriter, r *http.Request) {
		// since this is testnet, give caller coins from the testnet faucet
		prepUser, err := userValidateHelper(w, r, UserRPC[11][2:], UserRPC[11][1])
		if err != nil {
			return
		}

		assetCode := r.URL.Query()["assetCode"][0]
		assetIssuer := r.URL.Query()["assetIssuer"][0]
		limit, err := utils.ToFloat(r.URL.Query()["limit"][0])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		seedpwd := r.URL.Query()["seedpwd"][0]
		seed, err := wallet.DecryptSeed(prepUser.StellarWallet.EncryptedSeed, seedpwd)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		txhash, err := assets.TrustAsset(assetCode, assetIssuer, limit, seed)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, txhash)
	})
}

// uploadFile uploads a file to ipfs and returns the ipfs hash of the uploaded file. This is a POST request
func uploadFile() {
	http.HandleFunc(UserRPC[12][0], func(w http.ResponseWriter, r *http.Request) {
		_, err := userValidateHelper(w, r, UserRPC[12][2:], UserRPC[12][1])
		if err != nil {
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		defer func() {
			if ferr := file.Close(); ferr != nil {
				err = ferr
			}
		}()

		supportedType := false
		header := fileHeader.Header.Get("Content-Type")
		// I guess people could change the content type here and set it to anything they want to, but doesn't
		// matter since we batch this off to ipfs anyway

		switch header {
		case "image/jpeg":
			supportedType = true
		case "image/png":
			supportedType = true
		case "application/pdf":
			supportedType = true
		}

		// can't do anything with extensions, so while decrypting from ipfs, we can attach
		// all three types and return to the user.
		if !supportedType {
			erpc.ResponseHandler(w, erpc.StatusNotAcceptable)
			return
		}

		// file type is supported, store in ipfs
		data, err := ioutil.ReadAll(file)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not read returned data") {
			return
		}

		hashString, err := ipfs.AddBytes(data)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not hash data to ipfs") {
			return
		}
		erpc.MarshalSend(w, hashString)
	})
}

// PlatformEmailResponse is a structure used to contain the platform's email response
type PlatformEmailResponse struct {
	Email string
}

// platformEmail returns the platform's email address
func platformEmail() {
	http.HandleFunc(UserRPC[13][0], func(w http.ResponseWriter, r *http.Request) {
		_, err := userValidateHelper(w, r, UserRPC[13][2:], UserRPC[13][1])
		if err != nil {
			return
		}

		var x PlatformEmailResponse
		x.Email = consts.PlatformEmail
		erpc.MarshalSend(w, x)
	})
}

// increaseTrustLimit increases the trust limit a user has towards a specific asset on stellar
func increaseTrustLimit() {
	http.HandleFunc(UserRPC[17][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[17][2:], UserRPC[17][1])
		if err != nil {
			return
		}

		seedpwd := r.URL.Query()["seedpwd"][0]

		// now the user is validated, we need to call the db function to increase the trust limit
		trust, err := utils.ToFloat(r.URL.Query()["trust"][0])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		err = prepUser.IncreaseTrustLimit(seedpwd, trust)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// sendSecrets sends secrets out to the email ids passed. This does not require the seedpwd since one can generate a new seed
// anyway using the username and password, so possessing the secrets does not require seed authentication
func sendSecrets() {
	http.HandleFunc(UserRPC[19][0], func(w http.ResponseWriter, r *http.Request) {
		user, err := userValidateHelper(w, r, UserRPC[19][2:], UserRPC[19][1])
		if err != nil {
			return
		}

		// we should distribute the shares and then set them to nil since a person who is in
		// control of the server c ould then reconstruct the seed
		// now send emails out to these three trusted entities with the share
		email1 := r.URL.Query()["email1"][0]
		email2 := r.URL.Query()["email2"][0]
		email3 := r.URL.Query()["email3"][0]

		err = notif.SendSecretsEmail(user.Email, email1, email2, email3, user.RecoveryShares[0],
			user.RecoveryShares[1], user.RecoveryShares[2])
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		// set the stored shares to nil since possessing them would enable an attacker to generate the secrets he needs by simply controlling the server
		user.RecoveryShares[0] = ""
		user.RecoveryShares[1] = ""
		user.RecoveryShares[2] = ""

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// SeedResponse is a wrapper around the Seed
type SeedResponse struct {
	Seed string
}

// mergeSecrets takes in two shares in a 2 of 3 Shamir Secret Sharing Scheme and reconstructs the seed
func mergeSecrets() {
	http.HandleFunc(UserRPC[20][0], func(w http.ResponseWriter, r *http.Request) {
		_, err := userValidateHelper(w, r, UserRPC[20][2:], UserRPC[20][1])
		if err != nil {
			return
		}

		var shares []string
		secret1 := r.URL.Query()["secret1"][0]
		secret2 := r.URL.Query()["secret2"][0]
		shares = append(shares, secret1, secret2)
		// now we have 2 out of the 3 secrets needed to reconstruct. Reconstruct the seed.
		secret, err := recovery.Combine(shares)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "couldn't combine shares") {
			return
		}

		var x SeedResponse
		x.Seed = secret
		erpc.MarshalSend(w, x)
	})
}

// generateNewSecrets generates a new set of secrets for the given function
func generateNewSecrets() {
	http.HandleFunc(UserRPC[21][0], func(w http.ResponseWriter, r *http.Request) {
		user, err := userValidateHelper(w, r, UserRPC[21][2:], UserRPC[21][1])
		if err != nil {
			return
		}

		seedpwd := r.URL.Query()["seedpwd"][0] // we've already validated this earlier
		email1 := r.URL.Query()["email1"][0]
		email2 := r.URL.Query()["email2"][0]
		email3 := r.URL.Query()["email3"][0]

		seed, err := wallet.DecryptSeed(user.StellarWallet.EncryptedSeed, seedpwd)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}
		// user has validated his seed and identity. Generate new shares and send them out
		shares, err := recovery.Create(2, 3, seed)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		err = notif.SendSecretsEmail(user.Email, email1, email2, email3, shares[0], shares[1], shares[2])
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// generateResetPwdCode generates a password reset code
func generateResetPwdCode() {
	http.HandleFunc(UserRPC[22][0], func(w http.ResponseWriter, r *http.Request) {
		_, err := userValidateHelper(w, r, UserRPC[22][2:], UserRPC[22][1])
		if err != nil {
			return
		}

		email := r.URL.Query()["email"][0]

		rUser, err := database.SearchWithEmailID(email)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		_, err = ValidateSeedPwd(w, r, rUser.StellarWallet.EncryptedSeed, rUser.StellarWallet.PublicKey)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}
		// now we can verify that this is rellay the user. Now we need to cgenerate a verification code
		// and send it over to the user.
		verificationCode := utils.GetRandomString(16)
		log.Println("VERIFICATION CODE: ", verificationCode)
		rUser.PwdResetCode = verificationCode
		err = rUser.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		// now send this verification code to the email we have in the database
		err = notif.SendPasswordResetEmail(rUser.Email, verificationCode)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// resetPassword is a reset password route that can be called by the user in case they forget their password
func resetPassword() {
	http.HandleFunc(UserRPC[23][0], func(w http.ResponseWriter, r *http.Request) {
		_, err := userValidateHelper(w, r, UserRPC[23][2:], UserRPC[23][1])
		if err != nil {
			return
		}

		email := r.URL.Query()["email"][0]
		vCode := r.URL.Query()["verificationCode"][0]
		pwhash := r.URL.Query()["pwhash"][0]

		rUser, err := database.SearchWithEmailID(email)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		_, err = ValidateSeedPwd(w, r, rUser.StellarWallet.EncryptedSeed, rUser.StellarWallet.PublicKey)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		if vCode != rUser.PwdResetCode || vCode == "INVALID" {
			log.Println(rUser.PwdResetCode == vCode, vCode == "INVALID")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		// reset the user's password
		rUser.Pwhash = pwhash
		rUser.PwdResetCode = "INVALID" // invalidate the pwd reset code to avoid replay attacks
		err = rUser.Save()
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// sweepFunds tries to sweep all XLM that a user has from one account to another. Requires
// the seedpwd. Can't transfer assets automatically since platform does not know the list
// of issuer publickeys
func sweepFunds() {
	http.HandleFunc(UserRPC[24][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[24][2:], UserRPC[24][1])
		if err != nil {
			return
		}

		transferAddress := r.URL.Query()["destination"][0]
		if !xlm.AccountExists(transferAddress) {
			log.Println("Can only transfer to existing accounts, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		seedpwd, err := ValidateSeedPwd(w, r, prepUser.StellarWallet.EncryptedSeed, prepUser.StellarWallet.PublicKey)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		seed, err := wallet.DecryptSeed(prepUser.StellarWallet.EncryptedSeed, seedpwd)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		// validated the user, so now proceed to sweep funds
		xlmBalance := xlm.GetNativeBalance(prepUser.StellarWallet.PublicKey)
		log.Println(xlmBalance)
		// reduce 0.05 xlm and then sweep funds
		if xlmBalance < 5 {
			log.Println("xlm balance for user too small to sweep funds, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		xlmBalance -= 5
		// now we have the xlm balance, shift funds to the other account as requested by the user.
		sweepAmt := math.Round(xlmBalance)
		_, txhash, err := xlm.SendXLM(transferAddress, sweepAmt, seed, "sweep funds")
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		log.Println("sweep funds txhash: ", txhash)
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// sweepAsset sweeps a given asset from one account to another. Can't transfer multiple
// assets since we require the issuer pubkey(s)
func sweepAsset() {
	http.HandleFunc(UserRPC[25][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[25][2:], UserRPC[25][1])
		if err != nil {
			return
		}

		assetName := r.URL.Query()["assetName"][0]
		destination := r.URL.Query()["destination"][0]
		issuerPubkey := r.URL.Query()["issuerPubkey"][0]

		seedpwd, err := ValidateSeedPwd(w, r, prepUser.StellarWallet.EncryptedSeed, prepUser.StellarWallet.PublicKey)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		seed, err := wallet.DecryptSeed(prepUser.StellarWallet.EncryptedSeed, seedpwd)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		// validated the user, so now proceed to sweep funds
		assetBalance := xlm.GetAssetBalance(prepUser.StellarWallet.PublicKey, assetName)
		assetBalanceF, err := utils.ToFloat(assetBalance)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		// reduce 0.05 xlm and then sweep funds
		if assetBalanceF < 5 {
			log.Println("asset balance for user too smal lto sweep funds, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		assetBalanceF -= 5
		sweepAmt := math.Round(assetBalanceF)
		_, txhash, err := assets.SendAsset(assetName, issuerPubkey, destination, sweepAmt, seed, "sweeping funds")
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		log.Println("txhash: ", txhash)
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// KycResponse is a wrapper around status and reason for KYC responses
type KycResponse struct {
	Status string // the status whether the kyc verification request was succcessful or not
	Reason string // the reason why the person was rejected (OFAC blacklist, sanctioned individual, etc)
}

// validateKYC verifies whether a given user has passed kyc
func validateKYC() {
	http.HandleFunc(UserRPC[26][0], func(w http.ResponseWriter, r *http.Request) {
		// we first need to check the user params here
		prepUser, err := userValidateHelper(w, r, UserRPC[26][2:], UserRPC[26][1])
		if err != nil {
			return
		}

		var isID bool
		var idType string
		var id string
		var verif bool

		prepUser.KYC.PersonalPhoto = r.URL.Query()["selfie"][0]

		if r.URL.Query()["passport"] != nil {
			isID = true
			idType = "passport"
			id = r.URL.Query()["passport"][0]
			prepUser.KYC.PassportPhoto = id
		}

		if r.URL.Query()["dlicense"] != nil {
			isID = true
			idType = "dlicense"
			id = r.URL.Query()["dlicense"][0]
			prepUser.KYC.DriversLicense = id
		}

		if r.URL.Query()["idcard"] != nil {
			isID = true
			idType = "idcard"
			id = r.URL.Query()["idcard"][0]
			prepUser.KYC.IDCardPhoto = id
		}

		if !isID {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		var response KycResponse
		var apikey = consts.KYCAPIKey
		apiURL := "https://api.complyadvantage.com"
		body := apiURL + "/" + apikey

		switch idType {
		case "passport":
		case "dlicense":
			verif = true // solely for testing, remove once we add the real kyc provider in
		case "idcard":
			// no default since we check for that earlier
		}

		log.Println("requesting api verification for: " + body)
		// make the api request here, read response

		if verif {
			response.Status = "OK"
			response.Reason = ""
		} else {
			response.Status = "NOTOK"
			response.Reason = "Sanctioned Individual" // read the reason from the API response
		}

		err = prepUser.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}
		erpc.MarshalSend(w, response)
	})
}

// giveStarRating gives a star rating towards another person
func giveStarRating() {
	http.HandleFunc(UserRPC[27][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[27][2:], UserRPC[27][1])
		if err != nil {
			return
		}

		feedbackStr := r.URL.Query()["feedback"][0]
		uIndex := r.URL.Query()["userIndex"][0]

		feedback, err := utils.ToInt(feedbackStr)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		if feedback > 5 || feedback < 0 {
			log.Println("given feedback doesn't fall witin prescribed limits, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		userIndex, err := utils.ToInt(uIndex)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		err = prepUser.GiveFeedback(userIndex, feedback)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// TwoFAResponse is a wrapper around the QRCode data
type TwoFAResponse struct {
	ImageData string
}

// new2fa generates a new 2fa code
func new2fa() {
	http.HandleFunc(UserRPC[28][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[28][2:], UserRPC[28][1])
		if err != nil {
			return
		}

		if len(prepUser.TwoFASecret) != 0 {
			// user already has a 2fa secret, we need that in order to generate a new one
			if r.URL.Query()["password"] == nil {
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}

			password := r.URL.Query()["password"][0]
			result, err := prepUser.Authenticate2FA(password)
			if erpc.Err(w, err, erpc.StatusInternalServerError) {
				return
			}

			if !result {
				erpc.ResponseHandler(w, erpc.StatusUnauthorized)
				return
			}
			// now the old 2fa account is verified, we can proceed with creating a new 2fa secret
		}

		otpString, err := prepUser.Generate2FA()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		var x TwoFAResponse
		x.ImageData = otpString

		erpc.MarshalSend(w, x)
	})
}

// auth2fa authenticates the passed 2fa code
func auth2fa() {
	http.HandleFunc(UserRPC[29][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[29][2:], UserRPC[29][1])
		if err != nil {
			return
		}

		password := r.URL.Query()["password"][0]
		result, err := prepUser.Authenticate2FA(password)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		if !result {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// addAnchorKYCInfo adds anchorKYC info that the user passes to our platform.
func addAnchorKYCInfo() {
	http.HandleFunc(UserRPC[30][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[30][2:], UserRPC[30][1])
		if err != nil {
			return
		}

		prepUser.AnchorKYC.Name = r.URL.Query()["name"][0]
		prepUser.AnchorKYC.Birthday.Month = r.URL.Query()["bdaymonth"][0]
		prepUser.AnchorKYC.Birthday.Day = r.URL.Query()["bdayday"][0]
		prepUser.AnchorKYC.Birthday.Year = r.URL.Query()["bdayyear"][0]
		prepUser.AnchorKYC.Tax.Country = r.URL.Query()["taxcountry"][0]
		prepUser.AnchorKYC.Tax.ID = r.URL.Query()["taxid"][0]
		prepUser.AnchorKYC.Address.Street = r.URL.Query()["addrstreet"][0]
		prepUser.AnchorKYC.Address.City = r.URL.Query()["addrcity"][0]
		prepUser.AnchorKYC.Address.Postal = r.URL.Query()["addrpostal"][0]
		prepUser.AnchorKYC.Address.Region = r.URL.Query()["addrregion"][0]
		prepUser.AnchorKYC.Address.Country = r.URL.Query()["addrcountry"][0]
		prepUser.AnchorKYC.Address.Phone = r.URL.Query()["addrphone"][0]
		prepUser.AnchorKYC.PrimaryPhone = r.URL.Query()["primaryphone"][0]
		prepUser.AnchorKYC.Gender = r.URL.Query()["gender"][0]

		err = prepUser.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// changeReputationInv can be used to change the reputation of a sepcific investor on the platform
// on completion of a contract or on evaluation of feedback proposed by other entities on the system
func changeReputation() {
	http.HandleFunc(UserRPC[31][0], func(w http.ResponseWriter, r *http.Request) {
		user, err := userValidateHelper(w, r, UserRPC[31][2:], UserRPC[31][1])
		if err != nil {
			return
		}

		reputation, err := utils.ToFloat(r.URL.Query()["reputation"][0])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		err = user.ChangeReputation(reputation)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// ValidateSeedPwd validates only the seedpwd and not the username / pwhash
func ValidateSeedPwd(w http.ResponseWriter, r *http.Request, encryptedSeed []byte, userPublickey string) (string, error) {
	seedpwd := r.URL.Query()["seedpwd"][0]
	// we've validated the seedpwd, try decrypting the Encrypted Seed.
	seed, err := wallet.DecryptSeed(encryptedSeed, seedpwd)
	if err != nil {
		return seedpwd, errors.New("could not decrypt seed")
	}

	// now get the pubkey from this seed and match with original pubkey
	pubkey, err := wallet.ReturnPubkey(seed)
	if err != nil {
		return seedpwd, errors.New("could not retrieve pubkey")
	}

	if pubkey != userPublickey {
		return seedpwd, errors.New("pubkeys don't match, quitting")
	}

	return seedpwd, nil
}

// importSeed adds a user provided encrypted hex string to the openx platform. one can create their own keys and then import them onto openx
func importSeed() {
	http.HandleFunc(UserRPC[32][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, UserRPC[32][2:], UserRPC[32][1])
		if err != nil {
			return
		}

		encryptedSeedHex := r.URL.Query()["encryptedSeed"][0] // this will be a hex encoded string of the byte array
		pubkey := r.URL.Query()["pubkey"][0]
		seedpwd := r.URL.Query()["seedpwd"][0]

		encryptedSeed, err := hex.DecodeString(encryptedSeedHex)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		err = prepUser.ImportSeed(encryptedSeed, pubkey, seedpwd)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// getLatestBlockHash gets the latest Stellar blockchain hash from horizon
func getLatestBlockHash() {
	http.HandleFunc(UserRPC[33][0], func(w http.ResponseWriter, r *http.Request) {
		_, err := userValidateHelper(w, r, UserRPC[33][2:], UserRPC[33][1])
		if err != nil {
			return
		}

		hash, err := xlm.GetLatestBlockHash()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, hash)
	})
}

// acceptTc accepts the terms and conditions associated with openx
func acceptTc() {
	http.HandleFunc(UserRPC[35][0], func(w http.ResponseWriter, r *http.Request) {
		user, err := userValidateHelper(w, r, UserRPC[35][2:], UserRPC[35][1])
		if err != nil {
			return
		}

		if user.Legal {
			erpc.ResponseHandler(w, erpc.StatusOK)
			return
		}

		user.Legal = true
		err = user.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// updateProgress updates the profile progress bar on the frontend
func updateProgress() {
	http.HandleFunc(UserRPC[36][0], func(w http.ResponseWriter, r *http.Request) {
		user, err := userValidateHelper(w, r, UserRPC[36][2:], UserRPC[36][1])
		if err != nil {
			return
		}

		progressx := r.FormValue("progress")
		progress, err := utils.ToFloat(progressx)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		if progress > 100 || progress < 0 {
			log.Println("progress can't be greater than 100 or 0, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		user.ProfileProgress = progress
		err = user.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// updateUser updates credentials of the user
func updateUser() {
	http.HandleFunc(UserRPC[37][0], func(w http.ResponseWriter, r *http.Request) {
		user, err := userValidateHelper(w, r, UserRPC[37][2:], UserRPC[37][1])
		if err != nil {
			return
		}

		if r.FormValue("name") != "" {
			user.Name = r.FormValue("name")
		}
		if r.FormValue("city") != "" {
			user.City = r.FormValue("city")
		}
		if r.FormValue("pwhash") != "" {
			if len(r.FormValue("pwhash")) != 128 {
				log.Println("length of pwhash not 128")
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}
			user.Pwhash = r.FormValue("pwhash")
		}
		if r.FormValue("zipcode") != "" {
			user.ZipCode = r.FormValue("zipcode")
		}
		if r.FormValue("country") != "" {
			user.Country = r.FormValue("country")
		}
		if r.FormValue("recoveryphone") != "" {
			user.RecoveryPhone = r.FormValue("recoveryphone")
		}
		if r.FormValue("address") != "" {
			user.Address = r.FormValue("address")
		}
		if r.FormValue("description") != "" {
			user.Description = r.FormValue("description")
		}
		if r.FormValue("email") != "" {
			user.Email = r.FormValue("email")
		}
		if r.FormValue("seedpwd") != "" {
			if r.FormValue("oldseedpwd") == "" {
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}
			oldseedpwd := r.FormValue("oldseedpwd")
			seedpwd := r.FormValue("seedpwd")
			seed, err := wallet.DecryptSeed(user.StellarWallet.EncryptedSeed, oldseedpwd)
			if erpc.Err(w, err, erpc.StatusInternalServerError) {
				return
			}
			user.StellarWallet.EncryptedSeed, err = aes.Encrypt([]byte(seed), seedpwd)
			if erpc.Err(w, err, erpc.StatusInternalServerError) {
				return
			}
		}

		if r.FormValue("notification") != "" {
			if r.FormValue("notification") != "true" {
				user.Notification = false
			} else {
				user.Notification = true
			}
		}

		err = user.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, user)
	})
}

// downloadTeller downloads the teller data file
func downloadTeller() {
	http.HandleFunc(UserRPC[38][0], func(w http.ResponseWriter, r *http.Request) {
		//_, err := userValidateHelper(w, r, UserRPC[38][2:], UserRPC[38][1])
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		http.ServeFile(w, r, "screenlog.0")
	})
}

// logout logs out from all devices
func logout() {
	http.HandleFunc(UserRPC[39][0], func(w http.ResponseWriter, r *http.Request) {
		//_, err := userValidateHelper(w, r, UserRPC[38][2:], UserRPC[38][1])
		user, err := userValidateHelper(w, r, UserRPC[39][2:], UserRPC[39][1])
		if err != nil {
			return
		}

		err = user.AllLogout() // generate a new token to invalidate the old one
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// verify marks a user as verified
func verify() {
	http.HandleFunc(UserRPC[40][0], func(w http.ResponseWriter, r *http.Request) {
		//_, err := userValidateHelper(w, r, UserRPC[38][2:], UserRPC[38][1])
		user, err := userValidateHelper(w, r, UserRPC[40][2:], UserRPC[40][1])
		if err != nil {
			return
		}

		err = user.VerReq()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// unverify marks a user as verified
func unverify() {
	http.HandleFunc(UserRPC[41][0], func(w http.ResponseWriter, r *http.Request) {
		//_, err := userValidateHelper(w, r, UserRPC[38][2:], UserRPC[38][1])
		user, err := userValidateHelper(w, r, UserRPC[41][2:], UserRPC[41][1])
		if err != nil {
			return
		}

		err = user.UnverReq()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
