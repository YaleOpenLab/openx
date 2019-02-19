package rpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"crypto/tls"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	ipfs "github.com/YaleOpenLab/openx/ipfs"
	notif "github.com/YaleOpenLab/openx/notif"
	platform "github.com/YaleOpenLab/openx/platforms/opensolar"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
	xlm "github.com/YaleOpenLab/openx/xlm"
)

func setupUserRpcs() {
	ValidateUser()
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
	tellerPing()
}

const (
	TellerUrl = "https://localhost"
)

// we want to pass to the caller whether the user is a recipient or an investor.
// For this, we have an additional param called Role which we can use to classify
// this information and return to the caller
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
	recipient.U.EncryptedSeed = dummy
	return recipient
}

// removeSeedInv removes the encrypted seed from the investor structure
func removeSeedInv(investor database.Investor) database.Investor {
	var dummy []byte
	investor.U.EncryptedSeed = dummy
	return investor
}

// removeSeedEntity removes the encrypted seed from the entity structure
func removeSeedEntity(entity platform.Entity) platform.Entity {
	var dummy []byte
	entity.U.EncryptedSeed = dummy
	return entity
}

// UserValidateHelper is a helper that validates a user on the platform
func UserValidateHelper(w http.ResponseWriter, r *http.Request) (database.User, error) {
	checkGet(w, r)
	var prepUser database.User
	var err error
	// need to pass the pwhash param here
	if r.URL.Query() == nil || r.URL.Query()["username"] == nil || r.URL.Query()["pwhash"] == nil || len(r.URL.Query()["pwhash"][0]) != 128 {
		return prepUser, fmt.Errorf("Invalid params passed!")
	}

	prepUser, err = database.ValidateUser(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
	if err != nil {
		log.Println("did not validate user", err)
		return prepUser, err
	}

	return prepUser, nil
}

// ValidateUser is a route that helps validate users on the platform
func ValidateUser() {
	http.HandleFunc("/user/validate", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		// need to pass the pwhash param here
		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}
		// no we need to see whether this guy is an investor or a recipient.
		var prepInvestor database.Investor
		var prepRecipient database.Recipient
		var prepEntity platform.Entity
		rec := false
		entity := false
		prepInvestor, err = database.RetrieveInvestor(prepUser.Index)
		if err != nil {
			log.Println("did not retrieve investor", err)
			// means the user is a recipient, retrieve recipient credentials
			prepRecipient, err = database.ValidateRecipient(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
			if err != nil {
				log.Println("did not validate recipient", err)
				// it is not a recipient either
				prepEntity, err = platform.ValidateEntity(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
				if err != nil {
					log.Println("did not validate entity", err)
					// not an investor, recipient or entity, error
					responseHandler(w, r, StatusBadRequest)
					return
				} else {
					entity = true
				}
			} else {
				rec = true
			}
		}

		// the frontend should read the received response and figure out the role of the person
		var x ValidateParams
		if rec {
			x.Role = "Recipient"
			x.Entity = removeSeedRecp(prepRecipient)
			MarshalSend(w, r, x)
		} else if entity {
			x.Role = "Entity"
			x.Entity = removeSeedEntity(prepEntity)
			MarshalSend(w, r, x)
		} else {
			x.Role = "Investor"
			x.Entity = removeSeedInv(prepInvestor)
			MarshalSend(w, r, x)
		}
	})
}

// getBalances returns a list of all balances (assets and coins) held by the user
func getBalances() {
	http.HandleFunc("/user/balances", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		pubkey := prepUser.PublicKey
		balances, err := xlm.GetAllBalances(pubkey)
		if err != nil {
			log.Println("did not get all balances", err)
			responseHandler(w, r, StatusNotFound)
			return
		}
		MarshalSend(w, r, balances)
	})
}

// getXLMBalance gets the XLM balance of a user's account
func getXLMBalance() {
	http.HandleFunc("/user/balance/xlm", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)

		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		pubkey := prepUser.PublicKey
		balance, err := xlm.GetNativeBalance(pubkey)
		if err != nil {
			log.Println("did not get native balance", err)
			responseHandler(w, r, StatusNotFound)
			return
		}
		MarshalSend(w, r, balance)
	})
}

// getAssetBalance gets the balance of a specific asset
func getAssetBalance() {
	http.HandleFunc("/user/balance/asset", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)

		prepUser, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["asset"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		pubkey := prepUser.PublicKey
		asset := r.URL.Query()["asset"][0]
		balance, err := xlm.GetAssetBalance(pubkey, asset)
		if err != nil {
			log.Println("did not get assset balance", err)
			responseHandler(w, r, StatusNotFound)
			return
		}
		MarshalSend(w, r, balance)
	})
}

// getIpfsHash gets the ipfs hash of the passed string
func getIpfsHash() {
	http.HandleFunc("/ipfs/hash", func(w http.ResponseWriter, r *http.Request) {

		_, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["string"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		hashString := r.URL.Query()["string"][0]
		hash, err := ipfs.AddStringToIpfs(hashString)
		if err != nil {
			log.Println("did not add string to ipfs", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		hashCheck, err := ipfs.GetStringFromIpfs(hash)
		if err != nil || hashCheck != hashString {
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		MarshalSend(w, r, hash)
	})
}

// authKyc authenticates a user. Should ideally be part of a callback from the third
// party service that we choose
func authKyc() {
	http.HandleFunc("/user/kyc", func(w http.ResponseWriter, r *http.Request) {

		prepUser, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["userIndex"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		uInput := utils.StoI(r.URL.Query()["userIndex"][0])
		err = prepUser.Authorize(uInput)
		if err != nil {
			log.Println("did not authorize user", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}

// sendXLM sends a given amount of XLM to the destination address specified.
func sendXLM() {
	http.HandleFunc("/user/sendxlm", func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["destination"] == nil || r.URL.Query()["amount"] == nil ||
			r.URL.Query()["seedpwd"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		destination := r.URL.Query()["destination"][0]
		amount := r.URL.Query()["amount"][0]

		seedpwd := r.URL.Query()["seedpwd"][0]
		seed, err := wallet.DecryptSeed(prepUser.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not decrypt seed", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		var memo string
		if r.URL.Query()["memo"] != nil {
			memo = r.URL.Query()["memo"][0]
		}

		_, txhash, err := xlm.SendXLM(destination, amount, seed, memo)
		if err != nil {
			log.Println("did not send xlm", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}
		MarshalSend(w, r, txhash)
	})
}

// notKycView returns a list of all the users who have not yet been verified through KYC. Called by KYC Inspectors
func notKycView() {
	http.HandleFunc("/user/notkycview", func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		if !prepUser.Inspector {
			responseHandler(w, r, StatusUnauthorized)
			return
		}

		users, err := database.RetrieveAllUsersWithoutKyc()
		if err != nil {
			log.Println("did not retrieve all users wihtout kyc", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		MarshalSend(w, r, users)
	})
}

// kycView returns a list of all the users who have been verified through KYC. Called by KYC Inspectors
func kycView() {
	http.HandleFunc("/user/kycview", func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		if !prepUser.Inspector {
			responseHandler(w, r, StatusUnauthorized)
			return
		}

		users, err := database.RetrieveAllUsersWithKyc()
		if err != nil {
			log.Println("did not retrieve users", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		MarshalSend(w, r, users)
	})
}

// askForCoins asks for coins from the testnet faucet. Will be disabled once we move to testnet
func askForCoins() {
	http.HandleFunc("/user/askxlm", func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		err = xlm.GetXLM(prepUser.PublicKey)
		if err != nil {
			log.Println("did not get xlm from friendbot", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		responseHandler(w, r, StatusOK)
	})
}

// trustAsset creates a trustline for the given limit with a remote peer for receiving that asset.
func trustAsset() {
	http.HandleFunc("/user/trustasset", func(w http.ResponseWriter, r *http.Request) {
		// since this is testnet, give caller coins from the testnet faucet
		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		assetCode := r.URL.Query()["assetCode"][0]
		assetIssuer := r.URL.Query()["assetIssuer"][0]
		limit := r.URL.Query()["limit"][0]

		seedpwd := r.URL.Query()["seedpwd"][0]
		seed, err := wallet.DecryptSeed(prepUser.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not decrypt seed", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		// func TrustAsset(assetCode string, assetIssuer string, limit string, PublicKey string, Seed string) (string, error) {
		txhash, err := assets.TrustAsset(assetCode, assetIssuer, limit, prepUser.PublicKey, seed)
		if err != nil {
			log.Println("did not trust asset", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		MarshalSend(w, r, txhash)
	})
}

// uploadFile uploads a fil to ipfs and returns the ipfs hash of the uploaded file
// this is a POST request
func uploadFile() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

		checkPost(w, r)
		checkOrigin(w, r)

		_, err := UserValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			log.Println("did not parse form", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}
		defer file.Close()

		supportedType := false
		header := fileHeader.Header.Get("Content-Type")

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
			responseHandler(w, r, StatusNotAcceptable)
			return
		}

		// file type is supported, store in ipfs
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("did not  read", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		hashString, err := ipfs.IpfsHashData(data)
		if err != nil {
			log.Println("did not hash data to ipfs", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		MarshalSend(w, r, hashString)
	})
}

type PlatformEmailResponse struct {
	Email string
}

func platformEmail() {
	http.HandleFunc("/platformemail", func(w http.ResponseWriter, r *http.Request) {

		checkGet(w, r)
		checkOrigin(w, r)

		_, err := UserValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		var x PlatformEmailResponse
		x.Email = consts.PlatformEmail
		MarshalSend(w, r, x)
	})
}

func sendTellerShutdownEmail() {
	http.HandleFunc("/tellershutdown", func(w http.ResponseWriter, r *http.Request) {

		checkGet(w, r)
		checkOrigin(w, r)

		prepUser, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["projIndex"] == nil || r.URL.Query()["deviceId"] == nil ||
			r.URL.Query()["tx1"] == nil || r.URL.Query()["tx2"] == nil {
			log.Println("did not validate user", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		projIndex := r.URL.Query()["projIndex"][0]
		deviceId := r.URL.Query()["deviceId"][0]
		tx1 := r.URL.Query()["tx1"][0]
		tx2 := r.URL.Query()["tx2"][0]
		notif.SendTellerShutdownEmail(prepUser.Email, projIndex, deviceId, tx1, tx2)
		responseHandler(w, r, StatusOK)
	})
}

func sendTellerFailedPaybackEmail() {
	http.HandleFunc("/tellerpayback", func(w http.ResponseWriter, r *http.Request) {

		checkGet(w, r)
		checkOrigin(w, r)

		prepUser, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["projIndex"] == nil || r.URL.Query()["deviceId"] == nil {
			log.Println("did not validate user", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		projIndex := r.URL.Query()["projIndex"][0]
		deviceId := r.URL.Query()["deviceId"][0]
		notif.SendTellerPaymentFailedEmail(prepUser.Email, projIndex, deviceId)
		responseHandler(w, r, StatusOK)
	})
}

func tellerPing() {
	http.HandleFunc("/tellerping", func(w http.ResponseWriter, r *http.Request) {

		checkGet(w, r)
		checkOrigin(w, r)

		_, err := UserValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate user", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		req, err := http.NewRequest("GET", TellerUrl+"/ping", nil)
		if err != nil {
			log.Println("did not create new GET request", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		req.Header.Set("Origin", "localhost")
		res, err := client.Do(req)
		if err != nil {
			log.Println("did not make request", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		var x StatusResponse

		err = json.Unmarshal(data, &x)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		MarshalSend(w, r, x)
	})
}
