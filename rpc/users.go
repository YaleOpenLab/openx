package rpc

import (
	"net/http"

	database "github.com/OpenFinancing/openfinancing/database"
)

func setupUserRpcs() {
	ValidateUser()
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

func ValidateUser() {
	http.HandleFunc("/user/validate", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		// need to pass the pwhash param here
		if r.URL.Query() == nil || r.URL.Query()["username"] == nil || r.URL.Query()["pwhash"] == nil || len(r.URL.Query()["pwhash"][0]) != 128 {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		prepUser, err := database.ValidateUser(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		// no we need to see whether this guy is an investor or a recipient.
		var prepInvestor database.Investor
		var prepRecipient database.Recipient
		rec := false
		prepInvestor, err = database.RetrieveInvestor(prepUser.Index)
		if err != nil {
			// means the user is a recipient, retrieve recipient credentials
			rec = true
			prepRecipient, err = database.ValidateRecipient(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
			if err != nil {
				errorHandler(w, r, http.StatusNotFound)
				return
			}
		}

		// the frontend should read the received response and figure out the role of the person
		var x ValidateParams
		if rec {
			x.Role = "Recipient"
			x.Entity = removeSeedRecp(prepRecipient)
			MarshalSend(w, r, x)
		} else {
			x.Role = "Investor"
			x.Entity = removeSeedInv(prepInvestor)
			MarshalSend(w, r, x)
		}
	})
}
