package rpc

import (
	"fmt"
	"log"
	"net/http"

	database "github.com/OpenFinancing/openfinancing/database"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	utils "github.com/OpenFinancing/openfinancing/utils"
)

// setupRecipientRPCs sets up all RPCs related to the recipient. Most are similar
// to the investor RPCs, so maybe there's some nice way we can group them together
// to avoid code duplication
// not exporting this function because its being used only within the same package
func setupRecipientRPCs() {
	insertRecipient()
	validateRecipient()
	getAllRecipients()
	payback()
	storeDeviceId()
	storeStartTime()
}

func parseRecipient(r *http.Request) (database.Recipient, error) {
	var prepRecipient database.Recipient
	err := r.ParseForm()
	if err != nil || r.FormValue("LoginUserName") == "" || r.FormValue("LoginPassword") == "" || r.FormValue("Name") == "" || r.FormValue("EPassword") == "" {
		// don't care which type of error because you send 404 anyway
		return prepRecipient, fmt.Errorf("One of required fields missing: LoginUserName, LoginPassword, Name, EPassword")
	}

	prepRecipient.U, err = database.NewUser(r.FormValue("LoginUserName"), r.FormValue("LoginPassword"), r.FormValue("Name"), r.FormValue("EPassword"))
	log.Println("Parsed recipient: ", prepRecipient)
	return prepRecipient, err
}

func getAllRecipients() {
	http.HandleFunc("/recipient/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		recipients, err := database.RetrieveAllRecipients()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		log.Println("Retrieved all recipients: ", recipients)
		MarshalSend(w, r, recipients)
	})
}

func insertRecipient() {
	// this should be a post method since you want to accept an project and then insert
	// that into the database
	http.HandleFunc("/recipient/insert", func(w http.ResponseWriter, r *http.Request) {
		checkPost(w, r)
		prepRecipient, err := parseRecipient(r)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		log.Println("Parsed recipient:", prepRecipient)
		err = prepRecipient.Save()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		var rt StatusResponse
		rt.Status = 200
		MarshalSend(w, r, rt)
	})
}

func validateRecipient() {
	http.HandleFunc("/recipient/validate", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		// need to pass the pwhash param here
		if r.URL.Query() == nil || r.URL.Query()["LoginUserName"] == nil ||
			len(r.URL.Query()["LoginPassword"][0]) != 128 {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		prepRecipient, err := database.ValidateRecipient(r.URL.Query()["LoginUserName"][0], r.URL.Query()["LoginPassword"][0])
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		log.Println("Parsed recipient:", prepRecipient)
		MarshalSend(w, r, prepRecipient)
	})
}

func payback() {
	// func Payback(recpIndex int, projIndex int, assetName string, amount string, recipientSeed string,
	// 	platformPubkey string) error {
	http.HandleFunc("/recipient/payback", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		// this is a get request to make things easier for the teller
		if r.URL.Query() == nil || r.URL.Query()["recpIndex"] == nil ||
			r.URL.Query()["projIndex"] == nil || r.URL.Query()["assetName"] == nil ||
			r.URL.Query()["amount"] == nil || r.URL.Query()["platformPublicKey"] == nil {
			log.Println("PARAM ERROR")
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		recpIndex := utils.StoI(r.URL.Query()["recpIndex"][0])
		projIndex := utils.StoI(r.URL.Query()["projIndex"][0])
		assetName := r.URL.Query()["assetName"][0]
		recipientSeed := r.URL.Query()["recipientSeed"][0]
		amount := r.URL.Query()["amount"][0]
		platformPublicKey := r.URL.Query()["platformPublicKey"][0]

		err := solar.Payback(recpIndex, projIndex, assetName, amount, recipientSeed, platformPublicKey)
		if err != nil {
			log.Println("PAYBACK ERROR: ", err)
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		var rt StatusResponse
		rt.Status = 200
		MarshalSend(w, r, rt)
	})
}

func storeDeviceId() {
	http.HandleFunc("/recipient/deviceId", func(w http.ResponseWriter, r *http.Request) {
		// first validate the recipient or anyone would be able to set device ids
		checkGet(w, r)
		// need to pass the pwhash param here
		if r.URL.Query() == nil || r.URL.Query()["LoginUserName"] == nil ||
			len(r.URL.Query()["LoginPassword"][0]) != 128 || r.URL.Query()["deviceid"] == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		prepRecipient, err := database.ValidateRecipient(r.URL.Query()["LoginUserName"][0], r.URL.Query()["LoginPassword"][0])
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		// we have the recipient ready. Now set the device id
		prepRecipient.DeviceId = r.URL.Query()["deviceid"][0]
		err = prepRecipient.Save()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		var rt StatusResponse
		rt.Status = 200
		MarshalSend(w, r, rt)
	})
}

func storeStartTime() {
	http.HandleFunc("/recipient/startdevice", func(w http.ResponseWriter, r *http.Request) {
		// first validate the recipient or anyone would be able to set device ids
		checkGet(w, r)
		// need to pass the pwhash param here
		if r.URL.Query() == nil || r.URL.Query()["LoginUserName"] == nil ||
			len(r.URL.Query()["LoginPassword"][0]) != 128 || r.URL.Query()["start"] == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		prepRecipient, err := database.ValidateRecipient(r.URL.Query()["LoginUserName"][0], r.URL.Query()["LoginPassword"][0])
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		// we have the recipient ready. Now set the device id
		prepRecipient.DeviceStarts = append(prepRecipient.DeviceStarts, r.URL.Query()["start"][0])
		err = prepRecipient.Save()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		var rt StatusResponse
		rt.Status = 200
		MarshalSend(w, r, rt)
	})
}
