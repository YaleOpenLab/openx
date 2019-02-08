package rpc

import (
	"fmt"
	//"log"
	"net/http"
	"strconv"

	assets "github.com/OpenFinancing/openfinancing/assets"
	database "github.com/OpenFinancing/openfinancing/database"
	notif "github.com/OpenFinancing/openfinancing/notif"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	utils "github.com/OpenFinancing/openfinancing/utils"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
)

// setupInvestorRPCs sets up all RPCs related to the investor
func setupInvestorRPCs() {
	insertInvestor()
	validateInvestor()
	getAllInvestors()
	investInProject()
	changeReputationInv()
	voteTowardsProject()
	addLocalAssetInv()
	invAssetInv()
	sendEmail()
}

// parseInvestor is a helper that can be used to validate POST data and assigns the passed form
// data to an Investor struct
func parseInvestor(r *http.Request) (database.Investor, error) {
	var prepInvestor database.Investor
	err := r.ParseForm()
	if err != nil || r.FormValue("username") == "" || r.FormValue("pwhash") == "" || r.FormValue("Name") == "" || r.FormValue("EPassword") == "" {
		return prepInvestor, fmt.Errorf("One of required fields missing: username, pwhash, Name, EPassword")
	}

	prepInvestor.AmountInvested = float64(0)
	prepInvestor.U, err = database.NewUser(r.FormValue("username"), r.FormValue("pwhash"), r.FormValue("Name"), r.FormValue("EPassword"))
	return prepInvestor, err
}

// insertInvestor inserts an investor in to the main platform database
func insertInvestor() {
	// this should be a post method since you want to accetp an project and then insert
	// that into the database
	http.HandleFunc("/investor/insert", func(w http.ResponseWriter, r *http.Request) {
		checkPost(w, r)
		prepInvestor, err := parseInvestor(r)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		err = prepInvestor.Save()
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		responseHandler(w, r, StatusCreated)
	})
}

// validateInvestor retrieves the investor after valdiating if such an ivnestor exists
// by checking the pwhash of the given investor with the stored one
func validateInvestor() {
	http.HandleFunc("/investor/validate", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		if r.URL.Query() == nil || r.URL.Query()["username"] == nil || r.URL.Query()["pwhash"] == nil ||
			len(r.URL.Query()["pwhash"][0]) != 128 { // sha 512 length
			responseHandler(w, r, StatusBadRequest)
			return
		}
		prepInvestor, err := database.ValidateInvestor(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		MarshalSend(w, r, prepInvestor)
	})
}

// getAllInvestors gets a list of all the investors in the system so that we can
// display it to some entity that is interested to view such stats
func getAllInvestors() {
	http.HandleFunc("/investor/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		investors, err := database.RetrieveAllInvestors()
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		MarshalSend(w, r, investors)
	})
}

// investInProject invests in a specific project of the user's choice
func investInProject() {
	http.HandleFunc("/investor/invest", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		// need the following params to invest in a project:
		// 1. Seed pwhash (for the investor)
		// 2. project index
		// 3. investment amount
		// 4. Login username (for the investor)
		// 5. Login pwhash (for the investor)

		investor, err := InvValidateHelper(w, r)
		if err != nil || r.URL.Query()["seedpwd"] == nil || r.URL.Query()["projIndex"] == nil ||
			r.URL.Query()["amount"] == nil { // sha 512 length
			responseHandler(w, r, StatusBadRequest)
			return
		}

		seedpwd := r.URL.Query()["seedpwd"][0]
		investorSeed, err := wallet.DecryptSeed(investor.U.EncryptedSeed, seedpwd)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		projIndex := utils.StoI(r.URL.Query()["projIndex"][0])
		amount := r.URL.Query()["amount"][0]
		investorPubkey, err := wallet.ReturnPubkey(investorSeed)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		// splitting the conditions into two since in the future we will be returning
		// error codes towards each type
		if !xlm.AccountExists(investorPubkey) {
			responseHandler(w, r, StatusNotFound)
			return
		}

		// note that while using this route, we can't send the investor assets (maybe)
		// make it so in the UI that only they can accept an investment so we can get their
		// seed and send them assets. By not accepting, they would forfeit their investment,
		// so incentive would be there to unlock the seed.
		_, err = solar.InvestInProject(projIndex, investor.U.Index, amount, investorSeed)
		if err != nil {
			responseHandler(w, r, StatusNotFound)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}

// InvValidateHelper is a helper that is used to validate an ivnestor on the platform
func InvValidateHelper(w http.ResponseWriter, r *http.Request) (database.Investor, error) {
	// first validate the investor or anyone would be able to set device ids
	checkGet(w, r)
	var prepInvestor database.Investor
	// need to pass the pwhash param here
	if r.URL.Query() == nil || r.URL.Query()["username"] == nil ||
		len(r.URL.Query()["pwhash"][0]) != 128 {
		return prepInvestor, fmt.Errorf("Invalid params passed")
	}

	prepInvestor, err := database.ValidateInvestor(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
	if err != nil {
		return prepInvestor, err
	}

	return prepInvestor, nil
}

// changeReputationInv can be used to change the reputation of a sepcific investor on the platform
// on completion of a contract or on evaluation of feedback proposed by other entities on the system
func changeReputationInv() {
	http.HandleFunc("/investor/reputation", func(w http.ResponseWriter, r *http.Request) {
		investor, err := InvValidateHelper(w, r)
		if err != nil || r.URL.Query()["reputation"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		reputation, err := strconv.ParseFloat(r.URL.Query()["reputation"][0], 32) // same as StoI but we need to catch the error here
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		err = database.ChangeInvReputation(investor.U.Index, reputation)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}

// voteTowardsProject votes towards a specific propsoed project of the user's choice.
func voteTowardsProject() {
	http.HandleFunc("/investor/vote", func(w http.ResponseWriter, r *http.Request) {
		investor, err := InvValidateHelper(w, r)
		if err != nil || r.URL.Query()["votes"] == nil || r.URL.Query()["projIndex"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		votes := utils.StoI(r.URL.Query()["votes"][0])
		projIndex := utils.StoI(r.URL.Query()["projIndex"][0])
		err = solar.VoteTowardsProposedProject(investor.U.Index, votes, projIndex)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}

// addLocalAssetInv adds a local asset that can be traded in a p2p fashion wihtout direct invlvement
// from the platform. The platform can have a UI that will deal with this or this can be
// made an emualtor only function so that only experienced users use this.
func addLocalAssetInv() {
	http.HandleFunc("/investor/localasset", func(w http.ResponseWriter, r *http.Request) {

		prepInvestor, err := InvValidateHelper(w, r)
		if err != nil || r.URL.Query()["assetName"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		assetName := r.URL.Query()["assetName"][0]
		prepInvestor.U.LocalAssets = append(prepInvestor.U.LocalAssets, assetName)
		err = prepInvestor.Save()
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		responseHandler(w, r, StatusOK)
	})
}

// invAssetInv sends a local asset to a remote peer
func invAssetInv() {
	http.HandleFunc("/investor/sendlocalasset", func(w http.ResponseWriter, r *http.Request) {
		prepInvestor, err := InvValidateHelper(w, r)
		if err != nil || r.URL.Query()["assetName"] == nil || r.URL.Query()["seedpwd"] == nil ||
			r.URL.Query()["destination"] == nil || r.URL.Query()["amount"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		assetName := r.URL.Query()["assetName"][0]

		seedpwd := r.URL.Query()["seedpwd"][0]
		seed, err := wallet.DecryptSeed(prepInvestor.U.EncryptedSeed, seedpwd)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		destination := r.URL.Query()["destination"][0]
		amount := r.URL.Query()["amount"][0]

		found := true
		for _, elem := range prepInvestor.U.LocalAssets {
			if elem == assetName {
				found = true
			}
		}

		if !found {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		_, txhash, err := assets.SendAssetFromIssuer(assetName, destination, amount, seed, prepInvestor.U.PublicKey)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		MarshalSend(w, r, txhash)
	})
}

// sendEmail sends an email to a specific entity
func sendEmail() {
	http.HandleFunc("/investor/sendemail", func(w http.ResponseWriter, r *http.Request) {
		prepInvestor, err := InvValidateHelper(w, r)
		if err != nil || r.URL.Query()["message"] == nil || r.URL.Query()["to"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		message := r.URL.Query()["message"][0]
		to := r.URL.Query()["to"][0]
		err = notif.SendEmail(message, to, prepInvestor.U.Name)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		responseHandler(w, r, StatusOK)
	})
}
