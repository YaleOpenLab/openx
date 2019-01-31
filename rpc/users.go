package rpc

import (
	"fmt"
	"log"
	"net/http"

	database "github.com/OpenFinancing/openfinancing/database"
	ipfs "github.com/OpenFinancing/openfinancing/ipfs"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
	utils "github.com/OpenFinancing/openfinancing/utils"
)

func setupUserRpcs() {
	ValidateUser()
	getBalances()
	getXLMBalance()
	getAssetBalance()
	getIpfsHash()
	authKyc()
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
			errorHandler(w, r, http.StatusNotFound)
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
					errorHandler(w, r, http.StatusNotFound)
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
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		pubkey := prepUser.PublicKey
		balances, err := xlm.GetAllBalances(pubkey)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
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
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		pubkey := prepUser.PublicKey
		log.Println("PUBKEY: ", pubkey)
		balance, err := xlm.GetNativeBalance(pubkey)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
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
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		pubkey := prepUser.PublicKey
		asset := r.URL.Query()["asset"][0]
		balance, err := xlm.GetAssetBalance(pubkey, asset)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		MarshalSend(w, r, balance)
	})
}

func getIpfsHash() {
	http.HandleFunc("/ipfs/hash", func(w http.ResponseWriter, r *http.Request) {

		_, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["string"] == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		hashString := r.URL.Query()["string"][0]
		hash, err := ipfs.AddStringToIpfs(hashString)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		hashCheck, err := ipfs.GetStringFromIpfs(hash)
		if err != nil || hashCheck != hashString {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		MarshalSend(w, r, hash)
	})
}

func authKyc() {
	http.HandleFunc("/user/kyc", func(w http.ResponseWriter, r *http.Request) {

		prepUser, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["userIndex"] == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		uInput := utils.StoI(r.URL.Query()["userIndex"][0])
		err = prepUser.Authorize(uInput)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		Send200(w, r)
	})
}
