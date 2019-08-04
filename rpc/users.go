package rpc

import (
	"crypto/tls"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"

	xlm "github.com/Varunram/essentials/crypto/xlm"
	assets "github.com/Varunram/essentials/crypto/xlm/assets"
	wallet "github.com/Varunram/essentials/crypto/xlm/wallet"
	ipfs "github.com/Varunram/essentials/ipfs"
	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	notif "github.com/YaleOpenLab/openx/notif"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	recovery "github.com/bithyve/research/sss"
)

func setupUserRpcs() {
	updateUser()
	validateUser()
	getBalances()
	getXLMBalance()
	getAssetBalance()
	getIpfsHash()
	authKyc()
	sendXLM()
	notKycView()
	kycView()
	askForCoins()
	trustAsset()
	uploadFile()
	platformEmail()
	sendTellerShutdownEmail()
	sendTellerFailedPaybackEmail()
	tellerPing()
	increaseTrustLimit()
	addContractHash()
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
}

const (
	// TellerUrl defines the teller URL to check. In future, would be an array
	TellerUrl = "https://localhost"
)

// we want to pass to the caller whether the user is a recipient or an investor.
// For this, we have an additional param called Role which we can use to classify
// this information and return to the caller

// ValidateParams is a struct used fro validating user params
type ValidateParams struct {
	Role   string
	Entity interface{}
}

// removeSeedRecp removes the encrypted seed from the recipient structure
func removeSeedRecp(recipient database.Recipient) database.Recipient {
	// any field that is private needs to be set to null here. A person using the API
	// knows the username and password anyway, so the route must return all routes
	// that are accessible by a single login (uname + pwhash)
	var dummy []byte
	recipient.U.StellarWallet.EncryptedSeed = dummy
	return recipient
}

// removeSeedInv removes the encrypted seed from the investor structure
func removeSeedInv(investor database.Investor) database.Investor {
	var dummy []byte
	investor.U.StellarWallet.EncryptedSeed = dummy
	return investor
}

// removeSeedEntity removes the encrypted seed from the entity structure
func removeSeedEntity(entity opensolar.Entity) opensolar.Entity {
	var dummy []byte
	entity.U.StellarWallet.EncryptedSeed = dummy
	return entity
}

// CheckReqdParams is a helper that validates a user on the platform
func CheckReqdParams(w http.ResponseWriter, r *http.Request, options ...string) (database.User, error) {
	var prepUser database.User
	var err error
	// need to pass the pwhash param here
	if r.URL.Query() == nil {
		return prepUser, errors.New("url query can't be empty")
	}

	options = append(options, "username", "pwhash")

	for _, option := range options {
		if r.URL.Query()[option] == nil {
			return prepUser, errors.New("required param: " + option + "not specified, quitting")
		}
	}

	if len(r.URL.Query()["pwhash"][0]) != 128 {
		return prepUser, errors.New("pwhash length not 128, quitting")
	}

	if r.URL.Query()["seedpwd"] != nil {
		// check seed pwhash before decryption
		seedpwhash := utils.SHA3hash(r.URL.Query()["seedpwd"][0])
		prepUser, err = database.ValidateSeedpwd(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0], seedpwhash)
	} else {
		// no seedpwhash, normal call
		prepUser, err = database.ValidateUser(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
	}

	// catch the error from the relevant error call
	if err != nil {
		log.Println("did not validate user", err)
		return prepUser, err
	}

	return prepUser, nil
}

func updateUser() {
	/* List of changeable parameters for the user struct
	Name string
	City string
	ZipCode string
	Country string
	RecoveryPhone string
	Address string
	Description string
	Email string
	Notification bool
	*/
	http.HandleFunc("/user/update", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		user, err := CheckReqdParams(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}
		if r.URL.Query()["name"] != nil {
			user.Name = r.URL.Query()["name"][0]
		}
		if r.URL.Query()["city"] != nil {
			user.City = r.URL.Query()["city"][0]
		}
		if r.URL.Query()["zipcode"] != nil {
			user.ZipCode = r.URL.Query()["zipcode"][0]
		}
		if r.URL.Query()["country"] != nil {
			user.Country = r.URL.Query()["country"][0]
		}
		if r.URL.Query()["recoveryphone"] != nil {
			user.RecoveryPhone = r.URL.Query()["recoveryphone"][0]
		}
		if r.URL.Query()["address"] != nil {
			user.Address = r.URL.Query()["address"][0]
		}
		if r.URL.Query()["description"] != nil {
			user.Description = r.URL.Query()["description"][0]
		}
		if r.URL.Query()["email"] != nil {
			user.Email = r.URL.Query()["email"][0]
		}
		if r.URL.Query()["notification"] != nil {
			if r.URL.Query()["notification"][0] != "true" {
				user.Notification = false
			} else {
				user.Notification = true
			}
		}

		err = user.Save()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// check whether given user is an investor or recipient
		investor, err := InvValidateHelper(w, r)
		if err == nil {
			investor.U = &user
			err = investor.Save()
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
		}
		recipient, err := RecpValidateHelper(w, r)
		if err == nil {
			recipient.U = &user
			err = recipient.Save()
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
		// now we have the user, need to check which parts the user has specified
	})
}

// validateUser is a route that helps validate users on the platform
func validateUser() {
	http.HandleFunc("/user/validate", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		// need to pass the pwhash param here
		prepUser, err := CheckReqdParams(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}
		// no we need to see whether this guy is an investor or a recipient.
		var prepInvestor database.Investor
		var prepRecipient database.Recipient
		var prepEntity opensolar.Entity

		var x ValidateParams

		prepInvestor, err = database.RetrieveInvestor(prepUser.Index)
		if err == nil {
			x.Role = "Investor"
			x.Entity = removeSeedInv(prepInvestor)
			erpc.MarshalSend(w, x)
			return
		}

		prepRecipient, err = database.RetrieveRecipient(prepUser.Index)
		if err == nil {
			x.Role = "Recipient"
			x.Entity = removeSeedRecp(prepRecipient)
			erpc.MarshalSend(w, x)
			return
		}

		prepEntity, err = opensolar.RetrieveEntity(prepUser.Index)
		if err == nil {
			x.Role = "Entity"
			x.Entity = removeSeedEntity(prepEntity)
			erpc.MarshalSend(w, x)
			return
		}

	})
}

// getBalances returns a list of all balances (assets and coins) held by the user
func getBalances() {
	http.HandleFunc("/user/balances", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		prepUser, err := CheckReqdParams(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		pubkey := prepUser.StellarWallet.PublicKey
		balances, err := xlm.GetAllBalances(pubkey)
		if err != nil {
			log.Println("did not get all balances", err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
			return
		}
		erpc.MarshalSend(w, balances)
	})
}

// getXLMBalance gets the XLM balance of a user's account
func getXLMBalance() {
	http.HandleFunc("/user/balance/xlm", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		prepUser, err := CheckReqdParams(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		pubkey := prepUser.StellarWallet.PublicKey
		balance, err := xlm.GetNativeBalance(pubkey)
		if err != nil {
			log.Println("did not get native balance", err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
			return
		}
		erpc.MarshalSend(w, balance)
	})
}

// getAssetBalance gets the balance of a specific asset
func getAssetBalance() {
	http.HandleFunc("/user/balance/asset", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		prepUser, err := CheckReqdParams(w, r, "asset")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		pubkey := prepUser.StellarWallet.PublicKey
		asset := r.URL.Query()["asset"][0]
		balance, err := xlm.GetAssetBalance(pubkey, asset)
		if err != nil {
			log.Println("did not get asset balance", err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
			return
		}
		erpc.MarshalSend(w, balance)
	})
}

// getIpfsHash gets the ipfs hash of the passed string
func getIpfsHash() {
	http.HandleFunc("/ipfs/hash", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		_, err := CheckReqdParams(w, r, "string")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		hashString := r.URL.Query()["string"][0]
		hash, err := ipfs.IpfsAddString(hashString)
		if err != nil {
			log.Println("did not add string to ipfs", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		hashCheck, err := ipfs.IpfsGetString(hash)
		if err != nil || hashCheck != hashString {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, hash)
	})
}

// authKyc authenticates a user. Should ideally be part of a callback from the third
// party service that we choose
func authKyc() {
	http.HandleFunc("/user/kyc", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		prepUser, err := CheckReqdParams(w, r, "userIndex")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		uInput, err := utils.ToInt(r.URL.Query()["userIndex"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		err = prepUser.Authorize(uInput)
		if err != nil {
			log.Println("did not authorize user", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// sendXLM sends a given amount of XLM to the destination address specified.
func sendXLM() {
	http.HandleFunc("/user/sendxlm", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		prepUser, err := CheckReqdParams(w, r, "destination", "amount", "seedpwd")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		destination := r.URL.Query()["destination"][0]
		amount, err := utils.ToFloat(r.URL.Query()["amount"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		seedpwd := r.URL.Query()["seedpwd"][0]
		seed, err := wallet.DecryptSeed(prepUser.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not decrypt seed", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		var memo string
		if r.URL.Query()["memo"] != nil {
			memo = r.URL.Query()["memo"][0]
		}

		_, txhash, err := xlm.SendXLM(destination, amount, seed, memo)
		if err != nil {
			log.Println("did not send xlm", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		erpc.MarshalSend(w, txhash)
	})
}

// notKycView returns a list of all the users who have not yet been verified through KYC. Called by KYC Inspectors
func notKycView() {
	http.HandleFunc("/user/notkycview", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		prepUser, err := CheckReqdParams(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}
		if !prepUser.Inspector && !prepUser.Admin {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		users, err := database.RetrieveAllUsersWithoutKyc()
		if err != nil {
			log.Println("did not retrieve all users wihtout kyc", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, users)
	})
}

// kycView returns a list of all the users who have been verified through KYC. Called by KYC Inspectors
func kycView() {
	http.HandleFunc("/user/kycview", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		prepUser, err := CheckReqdParams(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}
		if !prepUser.Inspector && !prepUser.Admin {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		users, err := database.RetrieveAllUsersWithKyc()
		if err != nil {
			log.Println("did not retrieve users", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, users)
	})
}

// askForCoins asks for coins from the testnet faucet. Will be disabled once we move to mainnet
func askForCoins() {
	http.HandleFunc("/user/askxlm", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if consts.Mainnet {
			log.Println("Openx is in mainnet mode, can't ask for coins")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		prepUser, err := CheckReqdParams(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		err = xlm.GetXLM(prepUser.StellarWallet.PublicKey)
		if err != nil {
			log.Println("did not get xlm from friendbot", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// trustAsset creates a trustline for the given limit with a remote peer for receiving that asset.
func trustAsset() {
	http.HandleFunc("/user/trustasset", func(w http.ResponseWriter, r *http.Request) {
		// since this is testnet, give caller coins from the testnet faucet
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		prepUser, err := CheckReqdParams(w, r, "assetCode", "assetIssuer", "limit", "seedpwd")
		if err != nil {
			log.Println("did not validate user", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		assetCode := r.URL.Query()["assetCode"][0]
		assetIssuer := r.URL.Query()["assetIssuer"][0]
		limit, err := utils.ToFloat(r.URL.Query()["limit"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		seedpwd := r.URL.Query()["seedpwd"][0]
		seed, err := wallet.DecryptSeed(prepUser.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not decrypt seed", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		// func TrustAsset(assetCode string, assetIssuer string, limit string, PublicKey string, Seed string) (string, error) {
		txhash, err := assets.TrustAsset(assetCode, assetIssuer, limit, seed)
		if err != nil {
			log.Println("did not trust asset", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, txhash)
	})
}

// uploadFile uploads a file to ipfs and returns the ipfs hash of the uploaded file
// this is a POST request
func uploadFile() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckPost(w, r)
		erpc.CheckOrigin(w, r)
		_, err := CheckReqdParams(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			log.Println("did not parse form", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		defer file.Close()

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
		if err != nil {
			log.Println("did not read returned data", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		hashString, err := ipfs.IpfsAddBytes(data)
		if err != nil {
			log.Println("did not hash data to ipfs", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, hashString)
	})
}

// PlatformEmailResponse is a structure used to contain the platform's email response
type PlatformEmailResponse struct {
	Email string
}

func platformEmail() {
	http.HandleFunc("/platformemail", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		_, err := CheckReqdParams(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		var x PlatformEmailResponse
		x.Email = consts.PlatformEmail
		erpc.MarshalSend(w, x)
	})
}

func sendTellerShutdownEmail() {
	http.HandleFunc("/tellershutdown", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		prepUser, err := CheckReqdParams(w, r, "projIndex", "deviceId", "tx1", "tx2")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		projIndex := r.URL.Query()["projIndex"][0]
		deviceId := r.URL.Query()["deviceId"][0]
		tx1 := r.URL.Query()["tx1"][0]
		tx2 := r.URL.Query()["tx2"][0]
		notif.SendTellerShutdownEmail(prepUser.Email, projIndex, deviceId, tx1, tx2)
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func sendTellerFailedPaybackEmail() {
	http.HandleFunc("/tellerpayback", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		prepUser, err := CheckReqdParams(w, r, "deviceId", "projIndex")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		projIndex := r.URL.Query()["projIndex"][0]
		deviceId := r.URL.Query()["deviceId"][0]
		notif.SendTellerPaymentFailedEmail(prepUser.Email, projIndex, deviceId)
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func tellerPing() {
	http.HandleFunc("/tellerping", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		_, err := CheckReqdParams(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		req, err := http.NewRequest("GET", TellerUrl+"/ping", nil)
		if err != nil {
			log.Println("did not create new GET request", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		req.Header.Set("Origin", "localhost")
		res, err := client.Do(req)
		if err != nil {
			log.Println("did not make request", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var x erpc.StatusResponse

		err = json.Unmarshal(data, &x)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

func increaseTrustLimit() {
	http.HandleFunc("/user/increasetrustlimit", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		prepUser, err := CheckReqdParams(w, r, "trust", "seedpwd")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		// now the user is validated, we need to call the db function to increase the trust limit
		trust, err := utils.ToFloat(r.URL.Query()["trust"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		seedpwd := r.URL.Query()["seedpwd"][0]

		err = prepUser.IncreaseTrustLimit(seedpwd, trust)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// AddContractHash adds a specific contract hash to the database
func addContractHash() {
	http.HandleFunc("/utils/addhash", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		var err error
		_, err = CheckReqdParams(w, r, "projIndex", "choice", "choicestr")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		choice := r.URL.Query()["choice"][0]
		hashString := r.URL.Query()["choicestr"][0]
		projIndex, err := utils.ToInt(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println("passed project index not int, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		project, err := opensolar.RetrieveProject(projIndex)
		if err != nil {
			log.Println("couldn't retrieve prject index from database")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		// there are in total 5 types of hashes: OriginatorMoUHash, ContractorContractHash,
		// InvPlatformContractHash, RecPlatformContractHash, SpecSheetHash
		// lets have a fixed set of strings that we can map on here so we have a single endpoint for storing all these hashes

		// TODO: read from the pending docs map here and store this only if we need to.
		switch choice {
		case "omh":
			if project.Stage == 0 {
				project.StageData = append(project.StageData, hashString)
			}
		case "cch":
			if project.Stage == 2 {
				project.StageData = append(project.StageData, hashString)
			}
		case "ipch":
			if project.Stage == 4 {
				project.StageData = append(project.StageData, hashString)
			}
		case "rpch":
			if project.Stage == 4 {
				project.StageData = append(project.StageData, hashString)
			}
		case "ssh":
			if project.Stage == 5 {
				project.StageData = append(project.StageData, hashString)
			}
		default:
			log.Println("invalid choice passed, quitting!")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		err = project.Save()
		if err != nil {
			log.Println("error while saving project to db, quitting!")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// sendSecrets sends secrets out to the email ids passed. This does not require the seedpwd since one can generate a new seed
// anyway using the username and password, so possessing the secrets does not require seed authentication
func sendSecrets() {
	http.HandleFunc("/user/sendrecovery", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		user, err := CheckReqdParams(w, r, "email1", "email2", "email3")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		// we should distribute the shares and then set them to nil since a person who is in
		// control of the server c ould then reconstruct the seed
		// now send emails out to these three trusted entities with the share
		email1 := r.URL.Query()["email1"][0]
		email2 := r.URL.Query()["email2"][0]
		email3 := r.URL.Query()["email3"][0]

		err = notif.SendSecretsEmail(user.Email, email1, email2, email3, user.RecoveryShares[0], user.RecoveryShares[1], user.RecoveryShares[2])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// set the stored shares to nil since possessing them would enable an attacker to generate the secrets he needs by simply controlling the server
		user.RecoveryShares[0] = ""
		user.RecoveryShares[1] = ""
		user.RecoveryShares[2] = ""

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

type SeedResponse struct {
	Seed string
}

// mergeSecrets takes in two shares in a 2 of 3 Shamir Secret Sharing Scheme and reconstructs the seed
func mergeSecrets() {
	http.HandleFunc("/user/seedrecovery", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		_, err := CheckReqdParams(w, r, "secret1", "secret2")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		var shares []string
		secret1 := r.URL.Query()["secret1"][0]
		secret2 := r.URL.Query()["secret2"][0]
		shares = append(shares, secret1, secret2)
		// now we have 2 out of the 3 secrets needed to reconstruct. Reconstruct the seed.
		secret, err := recovery.Combine(shares)
		if err != nil {
			log.Println("couldn't combine shares: ", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var x SeedResponse
		x.Seed = secret
		erpc.MarshalSend(w, x)
	})
}

// generateNewSecrets generates an ew set of secrets for the given function
func generateNewSecrets() {
	http.HandleFunc("/user/newsecrets", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		user, err := CheckReqdParams(w, r, "seedpwd", "email1", "email2", "email3")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		seedpwd, err := ValidateSeedPwd(w, r, user.StellarWallet.EncryptedSeed, user.StellarWallet.PublicKey)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		seed, err := wallet.DecryptSeed(user.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		// user has validated his seed and identity. Generate new shares and send them out
		shares, err := recovery.Create(2, 3, seed)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		email1 := r.URL.Query()["email1"][0]
		email2 := r.URL.Query()["email2"][0]
		email3 := r.URL.Query()["email3"][0]

		err = notif.SendSecretsEmail(user.Email, email1, email2, email3, shares[0], shares[1], shares[2])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func generateResetPwdCode() {
	http.HandleFunc("/user/resetpwd", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		// the notion here si that the user must have his seedpwd in order to reset the password.
		// we retrieve the user using his email id and lookup his encrypted seed. If the
		// seed can be unlocked using hte seedpwd, we send a pwd reset email. One of two passwords
		// must be remembered
		if r.URL.Query()["email"] == nil || r.URL.Query()["seedpwd"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		email := r.URL.Query()["email"][0]

		rUser, err := database.SearchWithEmailId(email)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		_, err = ValidateSeedPwd(w, r, rUser.StellarWallet.EncryptedSeed, rUser.StellarWallet.PublicKey)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		// now we can verify that this is rellay the user. Now we need to cgenerate a verification code
		// and send it over to the user.
		verificationCode := utils.GetRandomString(16)
		log.Println("VERIFICATION CODE: ", verificationCode)
		rUser.PwdResetCode = verificationCode
		err = rUser.Save()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// now send this verification code to the email we have in the database
		err = notif.SendPasswordResetEmail(rUser.Email, verificationCode)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func resetPassword() {
	http.HandleFunc("/user/pwdreset", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if r.URL.Query()["email"] == nil || r.URL.Query()["seedpwd"] == nil || r.URL.Query()["verificationCode"] == nil ||
			r.URL.Query()["pwhash"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		email := r.URL.Query()["email"][0]
		vCode := r.URL.Query()["verificationCode"][0]
		pwhash := r.URL.Query()["pwhash"][0]

		rUser, err := database.SearchWithEmailId(email)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		_, err = ValidateSeedPwd(w, r, rUser.StellarWallet.EncryptedSeed, rUser.StellarWallet.PublicKey)
		if err != nil {
			log.Println("bad req1")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if vCode != rUser.PwdResetCode || vCode == "INVALID" {
			log.Println("bad req2")
			log.Println(rUser.PwdResetCode == vCode, vCode == "INVALID")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		// reset the user's password
		rUser.Pwhash = pwhash
		rUser.PwdResetCode = "INVALID" // invalidate the pwd reset code to avoid replay attacks
		err = rUser.Save()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// sweepFunds tries to sweep all funds that we have from one account to another. Requires
// the seedpwd. Can't transfre assets automatically since platform does not know the list
// of issuer publickeys
func sweepFunds() {
	http.HandleFunc("/user/sweep", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepUser, err := CheckReqdParams(w, r, "seedpwd", "destination")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		transferAddress := r.URL.Query()["destination"][0]
		if !xlm.AccountExists(transferAddress) {
			log.Println("Can only transfer to existing accounts, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		seedpwd, err := ValidateSeedPwd(w, r, prepUser.StellarWallet.EncryptedSeed, prepUser.StellarWallet.PublicKey)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		seed, err := wallet.DecryptSeed(prepUser.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		// validated the user, so now proceed to sweep funds
		xlmBalance, err := xlm.GetNativeBalance(prepUser.StellarWallet.PublicKey)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		xlmBalanceF, err := utils.ToFloat(xlmBalance)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// reduce 0.05 xlm and then sweep funds
		if xlmBalanceF < 5 {
			log.Println("xlm balance for user too small to sweep funds, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		xlmBalanceF -= 5
		// now we have the xlm balance, shift funds to the other account as requested by the user.
		sweepAmt := math.Round(xlmBalanceF)
		_, txhash, err := xlm.SendXLM(transferAddress, sweepAmt, seed, "sweep funds")
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		log.Println("sweep funds txhash: ", txhash)
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// sweepAsset sweeps a given asset from one account to another. Can't transfer multiple
// assets since we require the issuer pubkey
func sweepAsset() {
	http.HandleFunc("/user/sweepasset", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepUser, err := CheckReqdParams(w, r, "seedpwd", "destination", "assetName", "issuerPubkey")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		assetName := r.URL.Query()["assetName"][0]
		destination := r.URL.Query()["destination"][0]
		issuerPubkey := r.URL.Query()["issuerPubkey"][0]

		seedpwd, err := ValidateSeedPwd(w, r, prepUser.StellarWallet.EncryptedSeed, prepUser.StellarWallet.PublicKey)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		seed, err := wallet.DecryptSeed(prepUser.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		// validated the user, so now proceed to sweep funds
		assetBalance, err := xlm.GetAssetBalance(prepUser.StellarWallet.PublicKey, assetName)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		assetBalanceF, err := utils.ToFloat(assetBalance)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// reduce 0.05 xlm and then sweep funds
		if assetBalanceF < 5 {
			log.Println("asset balance for user too smal lto sweep funds, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		} else {
			assetBalanceF -= 5
		}

		sweepAmt := math.Round(assetBalanceF)
		_, txhash, err := assets.SendAsset(assetName, issuerPubkey, destination, sweepAmt, seed, "sweeping funds")
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		log.Println("txhash: ", txhash)
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

type KycResponse struct {
	Status string // the status whether the kyc verification request was succcessful or not
	Reason string // the reason why the person was rejected (OFAC blacklist, sanctioned individual, etc)
}

func validateKYC() {
	http.HandleFunc("/user/verifykyc", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		// we first need to check the user params here
		prepUser, err := CheckReqdParams(w, r, "selfie")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		var isId bool
		var idType string
		var id string
		var verif bool

		prepUser.KYC.PersonalPhoto = r.URL.Query()["selfie"][0]

		if r.URL.Query()["passport"] != nil {
			isId = true
			idType = "passport"
			id = r.URL.Query()["passport"][0]
			prepUser.KYC.PassportPhoto = id
		}

		if r.URL.Query()["dlicense"] != nil {
			isId = true
			idType = "dlicense"
			id = r.URL.Query()["dlicense"][0]
			prepUser.KYC.DriversLicense = id
		}

		if r.URL.Query()["idcard"] != nil {
			isId = true
			idType = "idcard"
			id = r.URL.Query()["idcard"][0]
			prepUser.KYC.IDCardPhoto = id
		}

		if !isId {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		var response KycResponse
		var apikey = consts.KYCAPIKey
		apiUrl := "https://api.complyadvantage.com"
		body := apiUrl + "/" + apikey

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
		if err != nil {
			log.Println("error while saving user credentials to database, quitting")
			erpc.MarshalSend(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, response)
	})
}

func giveStarRating() {
	http.HandleFunc("/user/giverating", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepUser, err := CheckReqdParams(w, r, "feedback", "userIndex")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		feedbackStr := r.URL.Query()["feedback"][0]
		uIndex := r.URL.Query()["userIndex"][0]

		feedback, err := utils.ToInt(feedbackStr)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if feedback > 5 || feedback < 0 {
			log.Println("given feedback doesn't fall witin prescribed limits, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		userIndex, err := utils.ToInt(uIndex)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = prepUser.GiveFeedback(userIndex, feedback)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

type TwoFAResponse struct {
	ImageData string
}

func new2fa() {
	http.HandleFunc("/user/2fa/generate", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepUser, err := CheckReqdParams(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
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
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}

			if !result {
				erpc.ResponseHandler(w, erpc.StatusUnauthorized)
				return
			}
			// now the old 2fa account is verified, we can proceed with creating a new 2fa secret
		}

		otpString, err := prepUser.Generate2FA()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var x TwoFAResponse
		x.ImageData = otpString

		erpc.MarshalSend(w, x)
	})
}

func auth2fa() {
	http.HandleFunc("/user/2fa/authenticate", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepUser, err := CheckReqdParams(w, r, "password")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		password := r.URL.Query()["password"][0]
		result, err := prepUser.Authenticate2FA(password)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		if !result {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		} else {
			erpc.ResponseHandler(w, erpc.StatusOK)
		}
	})
}

func addAnchorKYCInfo() {
	http.HandleFunc("/user/anchorusd/kyc", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepUser, err := CheckReqdParams(w, r, "name", "bdaymonth", "bdayday", "bdayyear", "taxcountry",
			"taxid", "addrstreet", "addrcity", "addrpostal", "addrregion", "addrcountry", "addrphone", "primaryphone", "gender")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		prepUser.AnchorKYC.Name = r.URL.Query()["name"][0]
		prepUser.AnchorKYC.Birthday.Month = r.URL.Query()["bdaymonth"][0]
		prepUser.AnchorKYC.Birthday.Day = r.URL.Query()["bdayday"][0]
		prepUser.AnchorKYC.Birthday.Year = r.URL.Query()["bdayyear"][0]
		prepUser.AnchorKYC.Tax.Country = r.URL.Query()["taxcountry"][0]
		prepUser.AnchorKYC.Tax.Id = r.URL.Query()["taxid"][0]
		prepUser.AnchorKYC.Address.Street = r.URL.Query()["addrstreet"][0]
		prepUser.AnchorKYC.Address.City = r.URL.Query()["addrcity"][0]
		prepUser.AnchorKYC.Address.Postal = r.URL.Query()["addrpostal"][0]
		prepUser.AnchorKYC.Address.Region = r.URL.Query()["addrregion"][0]
		prepUser.AnchorKYC.Address.Country = r.URL.Query()["addrcountry"][0]
		prepUser.AnchorKYC.Address.Phone = r.URL.Query()["addrphone"][0]
		prepUser.AnchorKYC.PrimaryPhone = r.URL.Query()["primaryphone"][0]
		prepUser.AnchorKYC.Gender = r.URL.Query()["gender"][0]

		err = prepUser.Save()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// changeReputationInv can be used to change the reputation of a sepcific investor on the platform
// on completion of a contract or on evaluation of feedback proposed by other entities on the system
func changeReputation() {
	http.HandleFunc("/user/reputation", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckReqdParams(w, r, "reputation")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		reputation, err := strconv.ParseFloat(r.URL.Query()["reputation"][0], 32) // same as StoI but we need to catch the error here
		if err != nil {
			log.Println("could not parse float", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		err = user.ChangeReputation(reputation)
		if err != nil {
			log.Println("did not change investor reputation", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

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
