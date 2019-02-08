package rpc

import (
	//"log"
	"net/http"
	"strings"

	database "github.com/OpenFinancing/openfinancing/database"
	bonds "github.com/OpenFinancing/openfinancing/platforms/bonds"
	utils "github.com/OpenFinancing/openfinancing/utils"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
)

func setupBondRPCs() {
	InvestInBond()
	getBondDetails()
	Search()
	GetAllBonds()
}

// CreateBond creates a new bond with the parameters passed to it
func CreateBond() {
	// newParams(mdate string, mrights string, stype string, intrate float64, rating string, bIssuer string, uWriter string
	// unitCost float64, itype string, nUnits int, tax string
	var bond1 bonds.ConstructionBond
	var err error
	bond1, err = bonds.NewBond("Dec 21 2049", "Maturation Rights Link", "Security Type", 5.4, "AAA", "Bond Issuer", "underwriter.com",
		100000, "Instrument Type", 100, "No Fed tax for 10 years", 1, "title", "location", "string")
	if err != nil {
		return
	}
	_, err = bonds.RetrieveBond(bond1.Params.Index)
	if err != nil {
		return
	}
}

// curl request attached for convenience
// curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -H "Origin: localhost" -H "Cache-Control: no-cache" -d 'InvestmentAmount=1000&BondIndex=1&InvIndex=2&InvSeedPwd=x&RecSeedPwd=x' "http://localhost:8080/bond/invest"
// InvestInBond invests a specific amount in a bond of the user's choice
func InvestInBond() {
	http.HandleFunc("/bond/invest", func(w http.ResponseWriter, r *http.Request) {
		checkPost(w, r)
		var err error
		var iBond bonds.ConstructionBond
		// need to receive a whole lot of parameters here
		// need the bond index passed so that we can retrieve the bond easily
		if r.FormValue("InvestmentAmount") == "" || r.FormValue("BondIndex") == "" || r.FormValue("InvIndex") == "" || r.FormValue("InvSeedPwd") == "" || r.FormValue("RecSeedPwd") == "" {
			responseHandler(w, r, StatusBadRequest)
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

		invAmount := r.FormValue("InvestmentAmount")
		bondIndex := utils.StoI(r.FormValue("BondIndex"))
		invIndex := utils.StoI(r.FormValue("InvIndex"))
		invSeedPwd := r.FormValue("InvSeedPwd")
		recSeedPwd := r.FormValue("RecSeedPwd")

		iBond, err = bonds.RetrieveBond(bondIndex)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		iRec, err := database.RetrieveRecipient(iBond.RecipientIndex)
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
		_, err = xlm.GetNativeBalance(iRec.U.PublicKey) // get testnet funds if their account is new
		if err != nil {
			err = xlm.GetXLM(iRec.U.PublicKey)
			if err != nil {
				responseHandler(w, r, StatusInternalServerError)
				return
			}
		}

		invSeed, err := iInv.U.GetSeed(invSeedPwd)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		recSeed, err := iRec.U.GetSeed(recSeedPwd)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		err = iBond.Invest(issuerPk, issuerSeed, &iInv, &iRec, invAmount, invSeed, recSeed)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		MarshalSend(w, r, iBond)
	})
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
		bond, err := bonds.RetrieveBond(uKey)
		if err != nil {
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
		allBonds, err := bonds.RetrieveAllBonds()
		if err != nil {
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
			allBonds, err := bonds.RetrieveAllBonds()
			if err != nil {
				responseHandler(w, r, StatusInternalServerError)
				return
			}
			MarshalSend(w, r, allBonds)
			// do bond stuff
		} else if strings.Contains(searchString, "coop") {
			// do coop stuff
			allCoops, err := bonds.RetrieveAllCoops()
			if err != nil {
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
