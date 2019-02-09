package rpc

import (
	"encoding/json"
	// "log"
	"net/http"

	database "github.com/YaleOpenLab/openx/database"
	bonds "github.com/YaleOpenLab/openx/platforms/ozones"
	utils "github.com/YaleOpenLab/openx/utils"
	xlm "github.com/YaleOpenLab/openx/xlm"
)

func setupCoopRPCs() {
	getCoopDetails()
	InvestInCoop()
	GetAllCoops()
}

// GetAllCoops gets a list of all the coops  that are registered on the platform
func GetAllCoops() {
	http.HandleFunc("/coop/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		allBonds, err := bonds.RetrieveAllBonds()
		if err != nil {
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
		bond, err := bonds.RetrieveCoop(uKey)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		bondJson, err := json.Marshal(bond)
		if err != nil {
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
		var iCoop bonds.Coop
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
			err = xlm.GetXLM(issuerPk)
			if err != nil {
				responseHandler(w, r, StatusInternalServerError)
				return
			}
		}

		invAmount := r.FormValue("MonthlyPayment")
		CoopIndex := utils.StoI(r.FormValue("CoopIndex"))
		invIndex := utils.StoI(r.FormValue("InvIndex"))
		invSeedPwd := r.FormValue("InvSeedPwd")

		iCoop, err = bonds.RetrieveCoop(CoopIndex)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		// pass the investor index, pk and seed
		iInv, err := database.RetrieveInvestor(invIndex)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		_, err = xlm.GetNativeBalance(iInv.U.PublicKey) // get testnet funds if their account is new
		if err != nil {
			err = xlm.GetXLM(iInv.U.PublicKey)
			if err != nil {
				responseHandler(w, r, StatusInternalServerError)
				return
			}
		}
		invSeed, err := iInv.U.GetSeed(invSeedPwd)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		err = iCoop.Invest(issuerPk, issuerSeed, &iInv, invAmount, invSeed)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		bondJson, err := json.Marshal(iCoop)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		WriteToHandler(w, bondJson)
	})
}
