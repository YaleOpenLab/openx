package rpc

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	assets "github.com/OpenFinancing/openfinancing/assets"
	database "github.com/OpenFinancing/openfinancing/database"
	ipfs "github.com/OpenFinancing/openfinancing/ipfs"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	utils "github.com/OpenFinancing/openfinancing/utils"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
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
}

// we want to pass to the caller whether the user is a recipient or an investor.
// For this, we have an additional param called Role which we can use to classify
// this information and return to the caller
type ValidateParams struct {
	Role   string
	Entity interface{}
}

func removeSeedRecp(recipient database.Recipient) database.Recipient {
	// any field that is private needs to be set to null here. A person using the API
	// knows the username and password anyway, so the route must return all routes
	// that are accessible by a single login (uname + pwhash)
	var dummy []byte
	recipient.U.EncryptedSeed = dummy
	return recipient
}

func removeSeedInv(investor database.Investor) database.Investor {
	var dummy []byte
	investor.U.EncryptedSeed = dummy
	return investor
}

func removeSeedEntity(entity solar.Entity) solar.Entity {
	var dummy []byte
	entity.U.EncryptedSeed = dummy
	return entity
}

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
		return prepUser, err
	}

	return prepUser, nil
}

func ValidateUser() {
	http.HandleFunc("/user/validate", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		// need to pass the pwhash param here
		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		// no we need to see whether this guy is an investor or a recipient.
		var prepInvestor database.Investor
		var prepRecipient database.Recipient
		var prepEntity solar.Entity
		rec := false
		entity := false
		prepInvestor, err = database.RetrieveInvestor(prepUser.Index)
		if err != nil {
			// means the user is a recipient, retrieve recipient credentials
			prepRecipient, err = database.ValidateRecipient(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
			if err != nil {
				// it is not a recipient either
				prepEntity, err = solar.ValidateEntity(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
				if err != nil {
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

func getBalances() {
	http.HandleFunc("/user/balances", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		pubkey := prepUser.PublicKey
		balances, err := xlm.GetAllBalances(pubkey)
		if err != nil {
			responseHandler(w, r, StatusNotFound)
			return
		}
		MarshalSend(w, r, balances)
	})
}

func getXLMBalance() {
	http.HandleFunc("/user/balance/xlm", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)

		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		pubkey := prepUser.PublicKey
		log.Println("PUBKEY: ", pubkey)
		balance, err := xlm.GetNativeBalance(pubkey)
		if err != nil {
			responseHandler(w, r, StatusNotFound)
			return
		}
		MarshalSend(w, r, balance)
	})
}

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
			responseHandler(w, r, StatusNotFound)
			return
		}
		MarshalSend(w, r, balance)
	})
}

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
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}

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
			responseHandler(w, r, StatusBadRequest)
			return
		}

		var memo string
		if r.URL.Query()["memo"] != nil {
			memo = r.URL.Query()["memo"][0]
		}

		_, txhash, err := xlm.SendXLM(destination, amount, seed, memo)
		if err != nil {
			log.Println(err)
		}
		MarshalSend(w, r, txhash)
	})
}

func notKycView() {
	http.HandleFunc("/user/notkycview", func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		if !prepUser.Inspector {
			responseHandler(w, r, StatusUnauthorized)
			return
		}

		users, err := database.RetrieveAllUsersWithoutKyc()
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		MarshalSend(w, r, users)
	})
}

func kycView() {
	http.HandleFunc("/user/kycview", func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		if !prepUser.Inspector {
			responseHandler(w, r, StatusUnauthorized)
			return
		}

		users, err := database.RetrieveAllUsersWithKyc()
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		MarshalSend(w, r, users)
	})
}

func askForCoins() {
	http.HandleFunc("/user/askxlm", func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		err = xlm.GetXLM(prepUser.PublicKey)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		responseHandler(w, r, StatusOK)
	})
}

func trustAsset() {
	http.HandleFunc("/user/trustasset", func(w http.ResponseWriter, r *http.Request) {
		// since this is testnet, give caller coins from the testnet faucet
		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		assetCode := r.URL.Query()["assetCode"][0]
		assetIssuer := r.URL.Query()["assetIssuer"][0]
		limit := r.URL.Query()["limit"][0]

		seedpwd := r.URL.Query()["seedpwd"][0]
		seed, err := wallet.DecryptSeed(prepUser.EncryptedSeed, seedpwd)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		// func TrustAsset(assetCode string, assetIssuer string, limit string, PublicKey string, Seed string) (string, error) {
		log.Println("ASSET CODE, ASSET ISSUER, LIMIT, PublicKey, SEED", assetCode, assetIssuer, limit, prepUser.PublicKey, seed)
		txhash, err := assets.TrustAsset(assetCode, assetIssuer, limit, prepUser.PublicKey, seed)
		if err != nil {
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
			responseHandler(w, r, StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
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
			log.Println(err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		hashString, err := ipfs.IpfsHashData(data)
		if err != nil {
			log.Println(err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		MarshalSend(w, r, hashString)
	})
}
