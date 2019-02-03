package rpc

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	database "github.com/OpenFinancing/openfinancing/database"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	utils "github.com/OpenFinancing/openfinancing/utils"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
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
	storeDeviceLocation()
	changeReputationRecp()
	chooseBlindAuction()
	unlock()
	addEmail()
	finalizeProject()
	originateProject()
	calculateTrustLimit()
}

func parseRecipient(r *http.Request) (database.Recipient, error) {
	var prepRecipient database.Recipient
	err := r.ParseForm()
	if err != nil || r.FormValue("username") == "" || r.FormValue("pwhash") == "" || r.FormValue("Name") == "" || r.FormValue("EPassword") == "" {
		// don't care which type of error because you send 404 anyway
		return prepRecipient, fmt.Errorf("One of required fields missing: username, pwhash, Name, EPassword")
	}

	prepRecipient.U, err = database.NewUser(r.FormValue("username"), r.FormValue("pwhash"), r.FormValue("Name"), r.FormValue("EPassword"))
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

		Send200(w, r)
	})
}

func validateRecipient() {
	http.HandleFunc("/recipient/validate", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		// need to pass the pwhash param here
		if r.URL.Query() == nil || r.URL.Query()["username"] == nil ||
			len(r.URL.Query()["pwhash"][0]) != 128 {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		prepRecipient, err := database.ValidateRecipient(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
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
		prepRecipient, err := RecpValidateHelper(w, r)
		if err != nil || r.URL.Query()["assetName"] == nil || r.URL.Query()["amount"] == nil ||
			r.URL.Query()["platformPublicKey"] == nil || r.URL.Query()["seedpwd"] == nil || r.URL.Query()["projIndex"] == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		recpIndex := prepRecipient.U.Index
		projIndex := utils.StoI(r.URL.Query()["projIndex"][0])
		assetName := r.URL.Query()["assetName"][0]
		seedpwd := r.URL.Query()["seedpwd"][0]
		amount := r.URL.Query()["amount"][0]
		platformPublicKey := r.URL.Query()["platformPublicKey"][0]

		recipientSeed, err := wallet.DecryptSeed(prepRecipient.U.EncryptedSeed, seedpwd)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		err = solar.Payback(recpIndex, projIndex, assetName, amount, recipientSeed, platformPublicKey)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		Send200(w, r)
	})
}

func RecpValidateHelper(w http.ResponseWriter, r *http.Request) (database.Recipient, error) {
	// first validate the recipient or anyone would be able to set device ids
	checkGet(w, r)
	var prepRecipient database.Recipient
	// need to pass the pwhash param here
	if r.URL.Query() == nil || r.URL.Query()["username"] == nil ||
		len(r.URL.Query()["pwhash"][0]) != 128 {
		return prepRecipient, fmt.Errorf("Invalid params passed")
	}

	prepRecipient, err := database.ValidateRecipient(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
	if err != nil {
		return prepRecipient, err
	}

	return prepRecipient, nil
}

func storeDeviceId() {
	http.HandleFunc("/recipient/deviceId", func(w http.ResponseWriter, r *http.Request) {
		// first validate the recipient or anyone would be able to set device ids
		prepRecipient, err := RecpValidateHelper(w, r)
		if err != nil || r.URL.Query()["deviceid"] == nil {
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
		Send200(w, r)
	})
}

func storeStartTime() {
	http.HandleFunc("/recipient/startdevice", func(w http.ResponseWriter, r *http.Request) {
		// first validate the recipient or anyone would be able to set device ids
		prepRecipient, err := RecpValidateHelper(w, r)
		if err != nil || r.URL.Query()["start"] == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		prepRecipient.DeviceStarts = append(prepRecipient.DeviceStarts, r.URL.Query()["start"][0])
		err = prepRecipient.Save()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		Send200(w, r)
	})
}

func storeDeviceLocation() {
	http.HandleFunc("/recipient/storelocation", func(w http.ResponseWriter, r *http.Request) {
		// first validate the recipient or anyone would be able to set device ids
		prepRecipient, err := RecpValidateHelper(w, r)
		if err != nil || r.URL.Query()["location"] == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		prepRecipient.DeviceLocation = r.URL.Query()["location"][0]
		err = prepRecipient.Save()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		Send200(w, r)
	})
}

// changeReputation changes the reputation of a specified recipient
func changeReputationRecp() {
	http.HandleFunc("/recipient/reputation", func(w http.ResponseWriter, r *http.Request) {
		recipient, err := RecpValidateHelper(w, r)
		if err != nil || r.URL.Query()["reputation"] == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		reputation, err := strconv.ParseFloat(r.URL.Query()["reputation"][0], 32) // same as StoI but we need to catch the error here
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		err = database.ChangeRecpReputation(recipient.U.Index, reputation)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		Send200(w, r)
	})
}

func chooseBlindAuction() {
	http.HandleFunc("/recipient/auction/choose/blind", func(w http.ResponseWriter, r *http.Request) {
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		allContracts, err := solar.RetrieveRecipientProjects(solar.ProposedProject, recipient.U.Index)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		bestContract, err := solar.SelectContractBlind(allContracts)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		err = bestContract.SetFinalizedProject()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		Send200(w, r)
	})
}

func chooseVickreyAuction() {
	http.HandleFunc("/recipient/auction/choose/vickrey", func(w http.ResponseWriter, r *http.Request) {
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		allContracts, err := solar.RetrieveRecipientProjects(solar.ProposedProject, recipient.U.Index)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		// the only differing part in the three auction routes. Would be nice if there were
		// some way to avoid repetition like this
		bestContract, err := solar.SelectContractVickrey(allContracts)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		err = bestContract.SetFinalizedProject()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		Send200(w, r)
	})
}

func chooseTimeAuction() {
	http.HandleFunc("/recipient/auction/choose/time", func(w http.ResponseWriter, r *http.Request) {
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		allContracts, err := solar.RetrieveRecipientProjects(solar.ProposedProject, recipient.U.Index)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		// the only differing part in the three auction routes. Would be nice if there were
		// some way to avoid repetition like this
		bestContract, err := solar.SelectContractTime(allContracts)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		err = bestContract.SetFinalizedProject()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		Send200(w, r)
	})
}

func unlock() {
	http.HandleFunc("/recipient/unlock", func(w http.ResponseWriter, r *http.Request) {
		recipient, err := RecpValidateHelper(w, r)
		if err != nil || r.URL.Query()["seedpwd"] == nil {
			log.Println(err)
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		seedpwd := r.URL.Query()["seedpwd"][0]
		log.Println("SEEDPWD: ", seedpwd)
		projIndex, err := utils.StoICheck(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println(err)
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		err = solar.UnlockProject(recipient.U.Username, recipient.U.Pwhash, projIndex, seedpwd)
		if err != nil {
			log.Println(err)
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		Send200(w, r)
	})
}

func addEmail() {
	http.HandleFunc("/recipient/addemail", func(w http.ResponseWriter, r *http.Request) {
		recipient, err := RecpValidateHelper(w, r)
		if err != nil || r.URL.Query()["email"] == nil {
			log.Println(err)
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		email := r.URL.Query()["email"][0]
		err = recipient.AddEmail(email)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		Send200(w, r)
	})
}

func finalizeProject() {
	http.HandleFunc("/recipient/finalize", func(w http.ResponseWriter, r *http.Request) {
		_, err := RecpValidateHelper(w, r)
		if err != nil || r.URL.Query()["projIndex"] == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		projIndex := utils.StoI(r.URL.Query()["projIndex"][0])
		project, err := solar.RetrieveProject(projIndex)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		err = project.SetFinalizedProject()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		Send200(w, r)
	})
}

func originateProject() {
	http.HandleFunc("/recipient/originate", func(w http.ResponseWriter, r *http.Request) {
		recipient, err := RecpValidateHelper(w, r)
		if err != nil || r.URL.Query()["projIndex"] == nil {
			log.Println("ERROR WHILE HANDLIGN RECPS: ", err)
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		projIndex := utils.StoI(r.URL.Query()["projIndex"][0])
		err = solar.RecipientAuthorize(projIndex, recipient.U.Index)
		if err != nil {
			log.Println("ERROR WHILE AUTHORIZING")
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		Send200(w, r)
	})
}

func calculateTrustLimit() {
	http.HandleFunc("/recipient/trustlimit", func(w http.ResponseWriter, r *http.Request) {
		recipient, err := RecpValidateHelper(w, r)
		if err != nil || r.URL.Query()["assetName"] == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		assetName := r.URL.Query()["assetName"][0]
		trustLimit, err := xlm.GetAssetTrustLimit(recipient.U.PublicKey, assetName)
		if err != nil {
			log.Println(err)
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		MarshalSend(w, r, trustLimit)
	})
}
