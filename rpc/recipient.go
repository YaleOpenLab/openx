package rpc

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	database "github.com/YaleOpenLab/openx/database"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	opzones "github.com/YaleOpenLab/openx/platforms/ozones"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
	xlm "github.com/YaleOpenLab/openx/xlm"
)

// setupRecipientRPCs sets up all RPCs related to the recipient. Most are similar
// to the investor RPCs, so maybe there's some nice way we can group them together
// to avoid code duplication
// not exporting this function because its being used only within the same package
func setupRecipientRPCs() {
	registerRecipient()
	validateRecipient()
	getAllRecipients()
	payback()
	storeDeviceId()
	storeStartTime()
	storeDeviceLocation()
	changeReputationRecp()
	chooseBlindAuction()
	chooseVickreyAuction()
	chooseTimeAuction()
	unlockOpenSolar()
	addEmail()
	finalizeProject()
	originateProject()
	calculateTrustLimit()
	unlockCBond()
	storeStateHash()
}

// parseRecipient parses a recipient from the passed form data and returns a recipient structure if
// the form data passed was accurate
func parseRecipient(r *http.Request) (database.Recipient, error) {
	var prepRecipient database.Recipient
	err := r.ParseForm()
	if err != nil || r.FormValue("username") == "" || r.FormValue("pwhash") == "" || r.FormValue("Name") == "" || r.FormValue("EPassword") == "" {
		// don't care which type of error because you send 404 anyway
		return prepRecipient, errors.New("one of required fields missing: username, pwhash, Name, EPassword")
	}

	prepRecipient.U, err = database.NewUser(r.FormValue("username"), r.FormValue("pwhash"), r.FormValue("Name"), r.FormValue("EPassword"))
	return prepRecipient, err
}

// getAllRecipients gets a list of all the recipients who have registered on the platform
func getAllRecipients() {
	http.HandleFunc("/recipient/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		recipients, err := database.RetrieveAllRecipients()
		if err != nil {
			log.Println("did not retrieve all recipients", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		MarshalSend(w, r, recipients)
	})
}

func registerRecipient() {
	http.HandleFunc("/recipient/register", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		// to register, we need the name, username and pwhash
		if r.URL.Query()["name"] == nil || r.URL.Query()["username"] == nil || r.URL.Query()["pwd"] == nil || r.URL.Query()["seedpwd"] == nil {
			log.Println("missing basic set of params that can be used ot validate a user")
			responseHandler(w, r, StatusBadRequest)
			return
		}

		name := r.URL.Query()["name"][0]
		username := r.URL.Query()["username"][0]
		pwd := r.URL.Query()["pwd"][0]
		seedpwd := r.URL.Query()["seedpwd"][0]

		// check for username collision here. IF the usernamer already exists, fetch details from that and register as investor
		duplicateUser, err := database.CheckUsernameCollision(username)
		if err != nil {
		  // username collision, check other fields by fetching user details for the collided user
		  if duplicateUser.Name == name && duplicateUser.Pwhash == pwd {
		    // this is the same user who wants to register as an investor now, check if encrypted seed decrypts
		    seed, err := wallet.DecryptSeed(duplicateUser.EncryptedSeed, seedpwd)
		    if err != nil {
		      responseHandler(w, r, StatusInternalServerError)
		      return
		    }
		    pubkey, err := wallet.ReturnPubkey(seed)
		    if err != nil {
		      responseHandler(w, r, StatusInternalServerError)
		      return
		    }
		    if pubkey != duplicateUser.PublicKey {
		      responseHandler(w, r, StatusUnauthorized)
		      return
		    }
		    var a database.Recipient
		    a.U = duplicateUser
		    err = a.Save()
		    if err != nil {
		      responseHandler(w, r, StatusInternalServerError)
		      return
		    }
		    MarshalSend(w, r, a)
		    return
		  }
		}

		user, err := database.NewRecipient(username, pwd, seedpwd, name)
		if err != nil {
			log.Println(err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		MarshalSend(w, r, user)
	})
}
// validateRecipient validates a recipient on the platform
func validateRecipient() {
	http.HandleFunc("/recipient/validate", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		if r.URL.Query() == nil || r.URL.Query()["username"] == nil ||
			len(r.URL.Query()["pwhash"][0]) != 128 {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		prepRecipient, err := database.ValidateRecipient(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
		if err != nil {
			log.Println("did not validate recipient", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}
		MarshalSend(w, r, prepRecipient)
	})
}

// payback pays back towards a specific invested order
func payback() {
	// func Payback(recpIndex int, projIndex int, assetName string, amount string, recipientSeed string,
	// 	platformPubkey string) error {
	http.HandleFunc("/recipient/payback", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		// this is a get request to make things easier for the teller
		prepRecipient, err := RecpValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusUnauthorized)
			return
		}
		if r.URL.Query()["assetName"] == nil || r.URL.Query()["amount"] == nil ||
			r.URL.Query()["seedpwd"] == nil || r.URL.Query()["projIndex"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		recpIndex := prepRecipient.U.Index
		projIndex := utils.StoI(r.URL.Query()["projIndex"][0])
		assetName := r.URL.Query()["assetName"][0]
		seedpwd := r.URL.Query()["seedpwd"][0]
		amount := r.URL.Query()["amount"][0]

		recipientSeed, err := wallet.DecryptSeed(prepRecipient.U.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not decrypt seed", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		err = opensolar.Payback(recpIndex, projIndex, assetName, amount, recipientSeed)
		if err != nil {
			log.Println("did not payback", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}

// RecpValidateHelper is a helper that helps validates recipients in routes
func RecpValidateHelper(w http.ResponseWriter, r *http.Request) (database.Recipient, error) {
	// first validate the recipient or anyone would be able to set device ids
	checkGet(w, r)
	checkOrigin(w, r)
	var prepRecipient database.Recipient
	// need to pass the pwhash param here
	if r.URL.Query() == nil || r.URL.Query()["username"] == nil ||
		len(r.URL.Query()["pwhash"][0]) != 128 {
		return prepRecipient, errors.New("invalid params passed")
	}

	prepRecipient, err := database.ValidateRecipient(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
	if err != nil {
		log.Println("did not validate recipient", err)
		return prepRecipient, err
	}

	return prepRecipient, nil
}

// storeDeviceId st ores the recipient's device id from the teller. Called by the teller
func storeDeviceId() {
	http.HandleFunc("/recipient/deviceId", func(w http.ResponseWriter, r *http.Request) {
		// first validate the recipient or anyone would be able to set device ids
		checkGet(w, r)
		checkOrigin(w, r)
		prepRecipient, err := RecpValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusUnauthorized)
			return
		}
		if r.URL.Query()["deviceid"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		// we have the recipient ready. Now set the device id
		prepRecipient.DeviceId = r.URL.Query()["deviceid"][0]
		err = prepRecipient.Save()
		if err != nil {
			log.Println("did not save recipient", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}

// storeStartTime stores the start time of the remote device installed as part of an invested project.
// Called by the teller
func storeStartTime() {
	http.HandleFunc("/recipient/startdevice", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		// first validate the recipient or anyone would be able to set device ids
		prepRecipient, err := RecpValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusUnauthorized)
			return
		}
		if r.URL.Query()["start"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		prepRecipient.DeviceStarts = append(prepRecipient.DeviceStarts, r.URL.Query()["start"][0])
		err = prepRecipient.Save()
		if err != nil {
			log.Println("did not save recipient", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}

// storeDeviceLocation stores the location of the remote device when it starts up. Called by the teller
func storeDeviceLocation() {
	http.HandleFunc("/recipient/storelocation", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		// first validate the recipient or anyone would be able to set device ids
		prepRecipient, err := RecpValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusUnauthorized)
			return
		}
		if r.URL.Query()["location"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		prepRecipient.DeviceLocation = r.URL.Query()["location"][0]
		err = prepRecipient.Save()
		if err != nil {
			log.Println("did not save recipient", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}

// changeReputation changes the reputation of a specified recipient
func changeReputationRecp() {
	http.HandleFunc("/recipient/reputation", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusUnauthorized)
			return
		}
		if r.URL.Query()["reputation"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		reputation, err := strconv.ParseFloat(r.URL.Query()["reputation"][0], 32) // same as StoI but we need to catch the error here
		if err != nil {
			log.Println("did not parse float", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}
		err = database.ChangeRecpReputation(recipient.U.Index, reputation)
		if err != nil {
			log.Println("did not cahnge reputation", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}

// chooseBlindAuction chooses a blind auction method to choose for the winner. Also commonly
// known as a 1st price auction.
func chooseBlindAuction() {
	http.HandleFunc("/recipient/auction/choose/blind", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate recipient", err)
			responseHandler(w, r, StatusUnauthorized)
			return
		}

		allContracts, err := opensolar.RetrieveRecipientProjects(opensolar.Stage2.Number, recipient.U.Index)
		if err != nil {
			log.Println("did not validate recipient projects", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		bestContract, err := opensolar.SelectContractBlind(allContracts)
		if err != nil {
			log.Println("did not select contract", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		err = bestContract.SetStage(4) // TODO: change to 3
		if err != nil {
			log.Println("did not set final project", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		responseHandler(w, r, StatusOK)
	})
}

// chooseVickreyAuction chooses a vickrey auction method to choose the winning contractor.
// also known as a second price auction
func chooseVickreyAuction() {
	http.HandleFunc("/recipient/auction/choose/vickrey", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate recipient", err)
			responseHandler(w, r, StatusUnauthorized)
			return
		}

		allContracts, err := opensolar.RetrieveRecipientProjects(opensolar.Stage2.Number, recipient.U.Index)
		if err != nil {
			log.Println("did not retrieve recipient projects", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		// the only differing part in the three auction routes. Would be nice if there were
		// some way to avoid repetition like this
		bestContract, err := opensolar.SelectContractVickrey(allContracts)
		if err != nil {
			log.Println("did not select contract", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		err = bestContract.SetStage(4) // change to 3 once done
		if err != nil {
			log.Println("did not set final project", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		responseHandler(w, r, StatusOK)
	})
}

// chooseTimeAuction chooses the winning contractor based on least completion time
func chooseTimeAuction() {
	http.HandleFunc("/recipient/auction/choose/time", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate recipient", err)
			responseHandler(w, r, StatusUnauthorized)
			return
		}

		allContracts, err := opensolar.RetrieveRecipientProjects(opensolar.Stage2.Number, recipient.U.Index)
		if err != nil {
			log.Println("did not retrieve recipient projects", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		// the only differing part in the three auction routes. Would be nice if there were
		// some way to avoid repetition like this
		bestContract, err := opensolar.SelectContractTime(allContracts)
		if err != nil {
			log.Println("did not select contract", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		err = bestContract.SetStage(4) // TODO: change to 3
		if err != nil {
			log.Println("did not set final project", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		responseHandler(w, r, StatusOK)
	})
}

// unlock unlocks a speciifc projectwhich has been invested in, signalling that the recipient
// has accepted the investment.
func unlockOpenSolar() {
	http.HandleFunc("/recipient/unlock/opensolar", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusUnauthorized)
			return
		}
		if r.URL.Query()["seedpwd"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		seedpwd := r.URL.Query()["seedpwd"][0]
		projIndex, err := utils.StoICheck(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println("did not parse to integer", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		err = opensolar.UnlockProject(recipient.U.Username, recipient.U.Pwhash, projIndex, seedpwd)
		if err != nil {
			log.Println("did not unlock project", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		responseHandler(w, r, StatusOK)
	})
}

// addEmail adds an email address to the recipient's profile
func addEmail() {
	http.HandleFunc("/recipient/addemail", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusUnauthorized)
			return
		}
		if r.URL.Query()["email"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		email := r.URL.Query()["email"][0]
		err = recipient.AddEmail(email)
		if err != nil {
			log.Println("did not add email", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}

// finalizeProject finalizes (ie moves from stage 2 to 3) a specific project. usually
// this shouldn't be called directly since tehre would be auctions for choosign the winning
// contractor
func finalizeProject() {
	http.HandleFunc("/recipient/finalize", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		_, err := RecpValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusUnauthorized)
			return
		}
		if r.URL.Query()["projIndex"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		projIndex := utils.StoI(r.URL.Query()["projIndex"][0])
		project, err := opensolar.RetrieveProject(projIndex)
		if err != nil {
			log.Println("did not retrieve project", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		err = project.SetStage(4) // TODO: in the future once this is defined well enough, this must be set to stage 3
		if err != nil {
			log.Println("did not set final project", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		responseHandler(w, r, StatusOK)
	})
}

// originateProject originates (ie moves from stage 0 to 1) a project
func originateProject() {
	http.HandleFunc("/recipient/originate", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusUnauthorized)
			return
		}
		if r.URL.Query()["projIndex"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		projIndex := utils.StoI(r.URL.Query()["projIndex"][0])
		err = opensolar.RecipientAuthorize(projIndex, recipient.U.Index)
		if err != nil {
			log.Println("did not authorize project", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		responseHandler(w, r, StatusOK)
	})
}

// calculateTrustLimit calculates the trust limit associated with a specific asset.
func calculateTrustLimit() {
	http.HandleFunc("/recipient/trustlimit", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusUnauthorized)
			return
		}
		if r.URL.Query()["assetName"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		assetName := r.URL.Query()["assetName"][0]
		trustLimit, err := xlm.GetAssetTrustLimit(recipient.U.PublicKey, assetName)
		if err != nil {
			log.Println("did not get trust limit", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		MarshalSend(w, r, trustLimit)
	})
}

// unlock unlocks a speciifc projectwhich has been invested in, signalling that the recipient
// has accepted the investment.
func unlockCBond() {
	http.HandleFunc("/recipient/unlock/opzones/cbond", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusUnauthorized)
			return
		}
		if r.URL.Query()["seedpwd"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		seedpwd := r.URL.Query()["seedpwd"][0]
		projIndex, err := utils.StoICheck(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println("did not parse to integer", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}

		err = opzones.UnlockProject(recipient.U.Username, recipient.U.Pwhash, projIndex, seedpwd, "constructionbond")
		if err != nil {
			log.Println("did not unlock project", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		responseHandler(w, r, StatusOK)
	})
}

// storeStateHash stores the start time of the remote device installed as part of an invested project.
// Called by the teller
func storeStateHash() {
	http.HandleFunc("/recipient/ssh", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		// first validate the recipient or anyone would be able to set device ids
		prepRecipient, err := RecpValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusUnauthorized)
			return
		}
		if r.URL.Query()["hash"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		prepRecipient.StateHashes = append(prepRecipient.StateHashes, r.URL.Query()["hash"][0])
		err = prepRecipient.Save()
		if err != nil {
			log.Println("did not save recipient", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}
