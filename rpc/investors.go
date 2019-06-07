package rpc

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	database "github.com/YaleOpenLab/openx/database"
	notif "github.com/YaleOpenLab/openx/notif"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	opzones "github.com/YaleOpenLab/openx/platforms/ozones"
	utils "github.com/YaleOpenLab/openx/utils"
	xlm "github.com/YaleOpenLab/openx/xlm"
	assets "github.com/YaleOpenLab/openx/xlm/assets"
	wallet "github.com/YaleOpenLab/openx/xlm/wallet"
)

// setupInvestorRPCs sets up all RPCs related to the investor
func setupInvestorRPCs() {
	registerInvestor()
	validateInvestor()
	getAllInvestors()
	invest()
	changeReputationInv()
	voteTowardsProject()
	addLocalAssetInv()
	invAssetInv()
	sendEmail()
	investInConstructionBond()
	investInLivingUnitCoop()
}

func registerInvestor() {
	http.HandleFunc("/investor/register", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		// to register, we need the name, username and pwhash
		if r.URL.Query()["name"] == nil || r.URL.Query()["username"] == nil || r.URL.Query()["pwd"] == nil || r.URL.Query()["seedpwd"] == nil {
			log.Println("missing basic set of params that can be used ot validate a user")
			responseHandler(w, StatusBadRequest)
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
			if duplicateUser.Name == name && duplicateUser.Pwhash == utils.SHA3hash(pwd) {
				// this is the same user who wants to register as an investor now, check if encrypted seed decrypts
				seed, err := wallet.DecryptSeed(duplicateUser.StellarWallet.EncryptedSeed, seedpwd)
				if err != nil {
					responseHandler(w, StatusInternalServerError)
					return
				}
				pubkey, err := wallet.ReturnPubkey(seed)
				if err != nil {
					responseHandler(w, StatusInternalServerError)
					return
				}
				if pubkey != duplicateUser.StellarWallet.PublicKey {
					responseHandler(w, StatusUnauthorized)
					return
				}
				var a database.Investor
				a.U = &duplicateUser
				err = a.Save()
				if err != nil {
					responseHandler(w, StatusInternalServerError)
					return
				}
				MarshalSend(w, a)
				return
			}
		}
		user, err := database.NewInvestor(username, pwd, seedpwd, name)
		if err != nil {
			log.Println(err)
			responseHandler(w, StatusInternalServerError)
			return
		}

		MarshalSend(w, user)
	})
}

// validateInvestor retrieves the investor after valdiating if such an ivnestor exists
// by checking the pwhash of the given investor with the stored one
func validateInvestor() {
	http.HandleFunc("/investor/validate", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		if r.URL.Query() == nil || r.URL.Query()["username"] == nil || r.URL.Query()["pwhash"] == nil ||
			len(r.URL.Query()["pwhash"][0]) != 128 { // sha 512 length
			responseHandler(w, StatusBadRequest)
			return
		}
		prepInvestor, err := database.ValidateInvestor(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
		if err != nil {
			log.Println("did not validate investor", err)
			responseHandler(w, StatusBadRequest)
			return
		}
		MarshalSend(w, prepInvestor)
	})
}

// getAllInvestors gets a list of all the investors in the system so that we can
// display it to some entity that is interested to view such stats
func getAllInvestors() {
	http.HandleFunc("/investor/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		investors, err := database.RetrieveAllInvestors()
		if err != nil {
			log.Println("did not retrieve all investors", err)
			responseHandler(w, StatusBadRequest)
			return
		}
		MarshalSend(w, investors)
	})
}

// Invest invests in a specific project of the user's choice
func invest() {
	http.HandleFunc("/investor/invest", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		// need the following params to invest in a project:
		// 1. Seed pwhash (for the investor)
		// 2. project index
		// 3. investment amount
		// 4. Login username (for the investor)
		// 5. Login pwhash (for the investor)

		investor, err := InvValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}
		if r.URL.Query()["seedpwd"] == nil || r.URL.Query()["projIndex"] == nil ||
			r.URL.Query()["amount"] == nil { // sha 512 length
			responseHandler(w, StatusBadRequest)
			return
		}

		seedpwd := r.URL.Query()["seedpwd"][0]
		investorSeed, err := wallet.DecryptSeed(investor.U.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not decrypt seed", err)
			responseHandler(w, StatusBadRequest)
			return
		}

		projIndex := utils.StoI(r.URL.Query()["projIndex"][0])
		amount := r.URL.Query()["amount"][0]
		investorPubkey, err := wallet.ReturnPubkey(investorSeed)
		if err != nil {
			log.Println("did not return pubkey", err)
			responseHandler(w, StatusBadRequest)
			return
		}
		// splitting the conditions into two since in the future we will be returning
		// error codes towards each type
		if !xlm.AccountExists(investorPubkey) {
			responseHandler(w, StatusNotFound)
			return
		}

		// note that while using this route, we can't send the investor assets (maybe)
		// make it so in the UI that only they can accept an investment so we can get their
		// seed and send them assets. By not accepting, they would forfeit their investment,
		// so incentive would be there to unlock the seed.
		err = opensolar.Invest(projIndex, investor.U.Index, amount, investorSeed)
		if err != nil {
			log.Println("did not invest in order", err)
			responseHandler(w, StatusNotFound)
			return
		}
		responseHandler(w, StatusOK)
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
		return prepInvestor, errors.New("invalid params passed")
	}

	prepInvestor, err := database.ValidateInvestor(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
	if err != nil {
		log.Println("did not validate investor", err)
		return prepInvestor, err
	}

	return prepInvestor, nil
}

// changeReputationInv can be used to change the reputation of a sepcific investor on the platform
// on completion of a contract or on evaluation of feedback proposed by other entities on the system
func changeReputationInv() {
	http.HandleFunc("/investor/reputation", func(w http.ResponseWriter, r *http.Request) {
		investor, err := InvValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}
		if r.URL.Query()["reputation"] == nil {
			responseHandler(w, StatusBadRequest)
			return
		}
		reputation, err := strconv.ParseFloat(r.URL.Query()["reputation"][0], 32) // same as StoI but we need to catch the error here
		if err != nil {
			log.Println("could not parse float", err)
			responseHandler(w, StatusBadRequest)
			return
		}
		err = database.ChangeInvReputation(investor.U.Index, reputation)
		if err != nil {
			log.Println("did not change investor reputation", err)
			responseHandler(w, StatusInternalServerError)
			return
		}
		responseHandler(w, StatusOK)
	})
}

// voteTowardsProject votes towards a specific propsoed project of the user's choice.
func voteTowardsProject() {
	http.HandleFunc("/investor/vote", func(w http.ResponseWriter, r *http.Request) {
		investor, err := InvValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}
		if r.URL.Query()["votes"] == nil || r.URL.Query()["projIndex"] == nil {
			responseHandler(w, StatusBadRequest)
			return
		}

		votes := utils.StoI(r.URL.Query()["votes"][0])
		projIndex := utils.StoI(r.URL.Query()["projIndex"][0])
		err = opensolar.VoteTowardsProposedProject(investor.U.Index, votes, projIndex)
		if err != nil {
			log.Println("did not vote towards proposed project", err)
			responseHandler(w, StatusInternalServerError)
			return
		}
		responseHandler(w, StatusOK)
	})
}

// addLocalAssetInv adds a local asset that can be traded in a p2p fashion wihtout direct invlvement
// from the platform. The platform can have a UI that will deal with this or this can be
// made an emualtor only function so that only experienced users use this.
func addLocalAssetInv() {
	http.HandleFunc("/investor/localasset", func(w http.ResponseWriter, r *http.Request) {

		prepInvestor, err := InvValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}
		if r.URL.Query()["assetName"] == nil {
			responseHandler(w, StatusBadRequest)
			return
		}

		assetName := r.URL.Query()["assetName"][0]
		prepInvestor.U.LocalAssets = append(prepInvestor.U.LocalAssets, assetName)
		err = prepInvestor.Save()
		if err != nil {
			log.Println("did not save investor", err)
			responseHandler(w, StatusInternalServerError)
			return
		}

		responseHandler(w, StatusOK)
	})
}

// invAssetInv sends a local asset to a remote peer
func invAssetInv() {
	http.HandleFunc("/investor/sendlocalasset", func(w http.ResponseWriter, r *http.Request) {
		prepInvestor, err := InvValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}
		if r.URL.Query()["assetName"] == nil || r.URL.Query()["seedpwd"] == nil ||
			r.URL.Query()["destination"] == nil || r.URL.Query()["amount"] == nil {
			responseHandler(w, StatusBadRequest)
			return
		}

		assetName := r.URL.Query()["assetName"][0]

		seedpwd := r.URL.Query()["seedpwd"][0]
		seed, err := wallet.DecryptSeed(prepInvestor.U.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not decrypt seed", err)
			responseHandler(w, StatusBadRequest)
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
			responseHandler(w, StatusBadRequest)
			return
		}

		_, txhash, err := assets.SendAssetFromIssuer(assetName, destination, amount, seed, prepInvestor.U.StellarWallet.PublicKey)
		if err != nil {
			log.Println("did not send asset from issuer", err)
			responseHandler(w, StatusInternalServerError)
			return
		}
		MarshalSend(w, txhash)
	})
}

// sendEmail sends an email to a specific entity
func sendEmail() {
	http.HandleFunc("/investor/sendemail", func(w http.ResponseWriter, r *http.Request) {
		prepInvestor, err := InvValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}
		if r.URL.Query()["message"] == nil || r.URL.Query()["to"] == nil {
			responseHandler(w, StatusBadRequest)
			return
		}

		message := r.URL.Query()["message"][0]
		to := r.URL.Query()["to"][0]
		err = notif.SendEmail(message, to, prepInvestor.U.Name)
		if err != nil {
			log.Println("did not send email", err)
			responseHandler(w, StatusBadRequest)
			return
		}
		responseHandler(w, StatusOK)
	})
}

// curl request attached for convenience
// curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -H "Origin: localhost" -H "Cache-Control: no-cache" -d 'InvestmentAmount=1000&BondIndex=1&InvIndex=2&seedpwd=x&recpSeedPwd=x' "http://localhost:8080/bond/invest"
// investInConstructionBond invests a specific amount in a bond of the user's choice
func investInConstructionBond() {
	http.HandleFunc("/constructionbond/invest", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		var err error

		prepInvestor, err := InvValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}
		if r.URL.Query()["amount"] == nil || r.URL.Query()["projIndex"] == nil || r.URL.Query()["seedpwd"] == nil {
			responseHandler(w, StatusBadRequest)
			return
		}

		invAmount := r.URL.Query()["amount"][0]
		projIndex := utils.StoI(r.URL.Query()["projIndex"][0])
		seedpwd := r.URL.Query()["seedpwd"][0]

		invSeed, err := wallet.DecryptSeed(prepInvestor.U.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not get investor seed from password", err)
			responseHandler(w, StatusBadRequest)
			return
		}

		err = opzones.InvestInConstructionBond(projIndex, prepInvestor.U.Index, invAmount, invSeed)
		if err != nil {
			log.Println("did not invest in bond", err)
			responseHandler(w, StatusBadRequest)
			return
		}
		responseHandler(w, StatusOK)
	})
}

// InvestInCoop invests in a coop of the user's choice
func investInLivingUnitCoop() {
	http.HandleFunc("/livingunitcoop/invest", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		var err error

		prepInvestor, err := InvValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}
		if r.URL.Query()["amount"] == nil || r.URL.Query()["projIndex"] == nil || r.URL.Query()["seedpwd"] == nil {
			log.Println("couldn't validate investor", err)
			responseHandler(w, StatusBadRequest)
			return
		}

		invAmount := r.URL.Query()["amount"][0]
		projIndex := utils.StoI(r.URL.Query()["projIndex"][0])
		seedpwd := r.URL.Query()["seedpwd"][0]

		invSeed, err := wallet.DecryptSeed(prepInvestor.U.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not get investor seed from password", err)
			responseHandler(w, StatusBadRequest)
			return
		}

		recpSeed := "SA5LO2G3XR37YY7566K2NHWQCK6PFXMF7UE64WGFBCOAPFHEKNSWT6PE"
		err = opzones.InvestInLivingUnitCoop(projIndex, prepInvestor.U.Index, invAmount, invSeed, recpSeed)
		if err != nil {
			log.Println("did not invest in the coop", err)
			responseHandler(w, StatusInternalServerError)
			return
		}

		responseHandler(w, StatusOK)
	})
}
