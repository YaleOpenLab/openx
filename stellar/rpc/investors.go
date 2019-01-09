package rpc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
)

// setupInvestorRPCs sets up all RPCs related to the investor
func setupInvestorRPCs() {
	insertInvestor()
	validateInvestor()
	getAllInvestors()
}

func parseInvestor(r *http.Request) (database.Investor, error) {
	var prepInvestor database.Investor
	err := r.ParseForm()
	if err != nil || r.FormValue("LoginUserName") == "" || r.FormValue("LoginPassword") == "" || r.FormValue("Name") == "" || r.FormValue("EPassword") == "" {
		return prepInvestor, fmt.Errorf("One of required fields missing: LoginUserName, LoginPassword, Name, EPassword")
	}

	prepInvestor.AmountInvested = float64(0)
	prepInvestor.U, err = database.NewUser(r.FormValue("LoginUserName"), r.FormValue("LoginPassword"), r.FormValue("Name"), r.FormValue("EPassword"))
	return prepInvestor, err
}

func insertInvestor() {
	// this should be a post method since you want to accetp an project and then insert
	// that into the database
	http.HandleFunc("/investor/insert", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkPost(w, r)
		prepInvestor, err := parseInvestor(r)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		log.Println("Prepared Investor:", prepInvestor)
		err = prepInvestor.Save()
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

// validateInvestor retreives the investor after valdiating if such an ivnestor exists
// by checking the pwhash of the given investor with the stored one
func validateInvestor() {
	http.HandleFunc("/investor/validate", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		if r.URL.Query() == nil || r.URL.Query()["LoginUserName"] == nil || r.URL.Query()["LoginPassword"] == nil || len(r.URL.Query()["LoginPassword"][0]) != 128 { // sha 512 length
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		prepInvestor, err := database.ValidateInvestor(r.URL.Query()["LoginUserName"][0], r.URL.Query()["LoginPassword"][0]) // TODO: Change this
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		log.Println("Prepared Investor:", prepInvestor)
		investorJson, err := json.Marshal(prepInvestor)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, investorJson)
	})
}

// getAllInvestors gets a list of all the investors in the system so that we can
// display it to some entity that is interested to view such stats
func getAllInvestors() {
	http.HandleFunc("/investor/all", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		investors, err := database.RetrieveAllInvestors()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		investorJson, err := json.Marshal(investors)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, investorJson)
	})
}
