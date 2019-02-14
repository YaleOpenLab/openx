package rpc

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	database "github.com/YaleOpenLab/openx/database"
	opzones "github.com/YaleOpenLab/openx/platforms/ozones"
	utils "github.com/YaleOpenLab/openx/utils"
	xlm "github.com/YaleOpenLab/openx/xlm"
)

func setupCoopRPCs() {
	getCoopDetails()
	InvestInCoop()
	GetAllCoops()
}

func setupBondRPCs() {
	getBondDetails()
	Search()
	GetAllBonds()
}

// GetAllCoops gets a list of all the coops  that are registered on the platform
func GetAllCoops() {
	http.HandleFunc("/coop/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		allBonds, err := opzones.RetrieveAllConstructionBonds()
		if err != nil {
			log.Println("did not retireve all bonds", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		MarshalSend(w, r, allBonds)
	})
}

// getCoopDetails gets teh details of a particular coop
func getCoopDetails() {
	http.HandleFunc("/coop/get", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		// get the details of a specific bond by key
		if r.URL.Query()["index"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		uKey := utils.StoI(r.URL.Query()["index"][0])
		bond, err := opzones.RetrieveLivingUnitCoop(uKey)
		if err != nil {
			log.Println("did not retireve coop", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}
		bondJson, err := json.Marshal(bond)
		if err != nil {
			log.Println("did not marhsal json", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		WriteToHandler(w, bondJson)
	})
}

// curl request attached for convenience
// curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -H "Origin: localhost" -H "Cache-Control: no-cache" -d 'MonthlyPayment=1000&CoopIndex=1&InvIndex=2&InvSeedPwd=x' "http://localhost:8080/coop/invest"
// InvestInCoop invests in a coop of the user's choice
func InvestInCoop() {
	http.HandleFunc("/coop/invest", func(w http.ResponseWriter, r *http.Request) {
		checkPost(w, r)
		var err error
		var iCoop opzones.LivingUnitCoop
		// need to receive a whole lot of parameters here
		// need the bond index passed so that we can retrieve the bond easily
		if r.FormValue("MonthlyPayment") == "" || r.FormValue("CoopIndex") == "" || r.FormValue("InvIndex") == "" || r.FormValue("InvSeedPwd") == "" {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		issuerSeed := "SBBYVEI4YNKZANRQEFH35U5GPEJ27MBLL7XHEKX5VC75QLJZWAXGX36Y"
		issuerPk := "GAEY5TVFYWBIIHF7PQCQVNIFTNIF7QSG4IH27HRW3DH476RI4NA2BPV3"
		_, err = xlm.GetNativeBalance(issuerPk)
		if err != nil {
			log.Println("did not get native balance", err)
			err = xlm.GetXLM(issuerPk)
			if err != nil {
				log.Println("did not get xlm from friendbot", err)
				responseHandler(w, r, StatusInternalServerError)
				return
			}
		}

		invAmount := r.FormValue("MonthlyPayment")
		CoopIndex := utils.StoI(r.FormValue("CoopIndex"))
		invIndex := utils.StoI(r.FormValue("InvIndex"))
		invSeedPwd := r.FormValue("InvSeedPwd")

		iCoop, err = opzones.RetrieveLivingUnitCoop(CoopIndex)
		if err != nil {
			log.Println("did not retrieve coop", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		// pass the investor index, pk and seed
		iInv, err := database.RetrieveInvestor(invIndex)
		if err != nil {
			log.Println("did not retrieve investor", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		_, err = xlm.GetNativeBalance(iInv.U.PublicKey) // get testnet funds if their account is new
		if err != nil {
			log.Println("did not get native balance", err)
			err = xlm.GetXLM(iInv.U.PublicKey)
			if err != nil {
				log.Println("did not get xlm from friendbot", err)
				responseHandler(w, r, StatusInternalServerError)
				return
			}
		}
		invSeed, err := iInv.U.GetSeed(invSeedPwd)
		if err != nil {
			log.Println("did not get the investor seed from password", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		err = iCoop.Invest(issuerPk, issuerSeed, &iInv, invAmount, invSeed)
		if err != nil {
			log.Println("did not invest in the coop", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		bondJson, err := json.Marshal(iCoop)
		if err != nil {
			log.Println("did not marshal json", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		WriteToHandler(w, bondJson)
	})
}

// CreateBond creates a new bond with the parameters passed to it
func CreateBond() {
	// newParams(mdate string, mrights string, stype string, intrate float64, rating string, bIssuer string, uWriter string
	// unitCost float64, itype string, nUnits int, tax string
	var bond opzones.ConstructionBond
	var err error
	bond, err = opzones.NewConstructionBond("Dec 21 2049", "Security Type", 5.4, "AAA", "Bond Issuer", "underwriter.com",
		100000, "Instrument Type", 100, "No Fed tax for 10 years", 1, "title", "location", "string")
	if err != nil {
		log.Println("did not create new bond", err)
		return
	}
	_, err = opzones.RetrieveConstructionBond(bond.Index)
	if err != nil {
		log.Println("did not retrieve bond", err)
		return
	}
}

// getBondDetails gets the details of a particular bond
func getBondDetails() {
	http.HandleFunc("/bond/get", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		// get the details of a specific bond by key
		if r.URL.Query()["index"] == nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		uKey := utils.StoI(r.URL.Query()["index"][0])
		bond, err := opzones.RetrieveConstructionBond(uKey)
		if err != nil {
			log.Println("did not retrieve bond", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		MarshalSend(w, r, bond)
	})
}

// GetAllBonds gets the list of all bonds that are registered on the platfomr
func GetAllBonds() {
	http.HandleFunc("/bond/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		allBonds, err := opzones.RetrieveAllConstructionBonds()
		if err != nil {
			log.Println("did not retrieve all bonds", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		MarshalSend(w, r, allBonds)
	})
}

// Search can be used for searching bonds and coops to a limited capacity
func Search() {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		// search for coop / bond  and return accordingly
		if r.URL.Query()["q"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		searchString := r.URL.Query()["q"][0]
		if strings.Contains(searchString, "bond") {
			allBonds, err := opzones.RetrieveAllConstructionBonds()
			if err != nil {
				log.Println("did not retrieve all bonds", err)
				responseHandler(w, r, StatusInternalServerError)
				return
			}
			MarshalSend(w, r, allBonds)
			// do bond stuff
		} else if strings.Contains(searchString, "coop") {
			// do coop stuff
			allCoops, err := opzones.RetrieveAllLivingUnitCoops()
			if err != nil {
				log.Println("did not retrieve bond", err)
				responseHandler(w, r, StatusInternalServerError)
				return
			}
			MarshalSend(w, r, allCoops)
		} else {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
	})
}
