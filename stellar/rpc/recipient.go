package rpc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
)

// setupRecipientRPCs sets up all RPCs related to the recipient. Most are similar
// to the investor RPCs, so maybe there's some nice way we can group them together
// to avoid code duplication
// not exporting this function because its being used only within the same package
func setupRecipientRPCs() {
	insertRecipient()
	validateRecipient()
	getAllRecipients()
}

func parseRecipient(r *http.Request) (database.Recipient, error) {
	var prepRecipient database.Recipient
	err := r.ParseForm()
	if err != nil || r.FormValue("LoginUserName") == "" || r.FormValue("LoginPassword") == "" || r.FormValue("Name") == "" || r.FormValue("EPassword") == "" {
		// don't care which type of error because you send 404 anyway
		return prepRecipient, fmt.Errorf("One of required fields missing: LoginUserName, LoginPassword, Name, EPassword")
	}

	prepRecipient.U, err = database.NewUser(r.FormValue("LoginUserName"), r.FormValue("LoginPassword"), r.FormValue("Name"), r.FormValue("EPassword"))
	log.Println("Prepared recipient: ", prepRecipient)
	return prepRecipient, err
}

func getAllRecipients() {
	http.HandleFunc("/recipient/all", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		recipients, err := database.RetrieveAllRecipients()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		log.Println("Retrieved all recipients: ", recipients)
		recipientJson, err := json.Marshal(recipients)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, recipientJson)
	})
}

func insertRecipient() {
	// this should be a post method since you want to accept an project and then insert
	// that into the database
	http.HandleFunc("/recipient/insert", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkPost(w, r)
		prepRecipient, err := parseRecipient(r)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		log.Println("Prepared Recipient:", prepRecipient)
		err = prepRecipient.Save()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		var rt StatusResponse
		rt.Status = 200
		rtJson, err := json.Marshal(rt)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, rtJson)
	})
}

func validateRecipient() {
	http.HandleFunc("/recipient/validate", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		// need to pass the pwhash param here
		if r.URL.Query() == nil || r.URL.Query()["LoginUserName"] == nil || len(r.URL.Query()["LoginPassword"][0]) != 128 {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		prepRecipient, err := database.ValidateRecipient(r.URL.Query()["LoginUserName"][0], r.URL.Query()["LoginPassword"][0]) // TODO: change this
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		log.Println("Prepared Recipient:", prepRecipient)
		investorJson, err := json.Marshal(prepRecipient)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, investorJson)
	})
}
