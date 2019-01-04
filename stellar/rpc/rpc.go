package rpc

// the rpc package contains functions related to the server which will be interacting
// with the frontend. Not expanding on this too much since this will be changing quite often
// TODO: update RPC package in line with recent changes on master
import (
	"fmt"
	"log"
	// "net"
	// "io/ioutil"
	// "net/rpc"
	"encoding/json"
	"net/http"

	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/stellar/go/keypair"
)

type PingResponse struct {
	// pingresponse returns "alive" when calle,d could be used by services
	// that scan for uptime
	Status string
}

type StatusResponse struct {
	Status int
}

func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Header().Add("Access-Control-Allow-Origin", "localhost")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonString)
}

func checkOrigin(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Origin") != "localhost" {
		// allow only our frontend UI to connect to our RPC instance
		http.Error(w, "404 page not found", http.StatusNotFound)
		return fmt.Errorf("Cross domain request error")
	}
	return nil
}

func checkGet(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		// reject wrong method entries
		http.Error(w, "404 page not found", http.StatusNotFound)
		return fmt.Errorf("Invalid Method error")
	}
	return nil
}

func checkPost(w http.ResponseWriter, r *http.Request) error {
	log.Println("CHEcking POST")
	if r.Method != "POST" {
		// reject wrong method entries
		http.Error(w, "404 page not found", http.StatusNotFound)
		return fmt.Errorf("Invalid Method error")
	}
	return nil
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 Page Not Found")
	}
}

func setupDefaultHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// default to 404 for every application not running on localhost
		errorHandler(w, r, http.StatusNotFound)
		return
	})
}

func setupPingHandler() {

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		err := checkOrigin(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		err = checkGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		var pr PingResponse
		pr.Status = "Alive"
		prJson, err := json.Marshal(pr)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, prJson)
	})
}

func getOpenProjects() {
	http.HandleFunc("/projects/open", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		// make a call to the db to get all projects
		// while making this call, the rpc should not be aware of the db we are using
		// and stuff. So we need to have another route that would open the existing
		// db, without asking for one
		allContracts, err := database.RetrieveAllProjects()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		contractsJson, err := json.Marshal(allContracts)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, contractsJson)
	})
}

func getProject() {
	// we need to read passed the key from the URL that the user calls
	http.HandleFunc("/projects/get", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		URLPath := r.URL.Path
		// so we now have the URL path
		// slice "/project/" off from the URLPath
		keyS := URLPath[7:]
		// we now need to get the project corresponding to keyS
		// the rpc accepts the key as int though, so string -> int

		uKey := utils.StoI(keyS)
		contract, err := database.RetrieveProject(uKey)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		projectJson, err := json.Marshal(contract.Params)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, projectJson)
	})
}

func parseProject(r *http.Request) (database.DBParams, error) {
	// we need to create an instance of the Project
	// and then map values if they do exist
	// note that we just prepare the project here and don't invest in it
	// for that, we need new a new investor struct and a recipient struct
	// Index         int
	// PanelSize     string  required
	// TotalValue    int     required
	// Location      string  required
	// MoneyRaised   int     default to 0
	// Metadata      string  optional
	// Live          bool    will be set to true if this is called
	// INVAssetCode  string  should be set by subsequent calls
	// DEBAssetCode  string  should be set by subsequent calls
	// PBAssetCode   string  should be set by subsequent calls
	// BalLeft       float64 should be equal to  totalValue since this is a new project
	// DateInitiated string  auto
	// DateLastPaid  string  don't set

	var prepProject database.DBParams
	err := r.ParseForm()
	if err != nil {
		return prepProject, err
	}
	// if we're inserting this in, we need to get the next index number
	// so that we can set this without causing some weird bugs
	allContracts, err := database.RetrieveAllProjects()
	if err != nil {
		return prepProject, fmt.Errorf("Error in assigning index")
	}
	prepProject.Index = len(allContracts) + 1
	if r.FormValue("PanelSize") != "" {
		prepProject.PanelSize = r.FormValue("PanelSize")
	} else {
		// if this is not defined, error out
		return prepProject, fmt.Errorf("No PanelSize")
	}

	if r.FormValue("TotalValue") != "" {
		// the totlaValue passed here is a string, we need to convert to an int
		totalValueS := utils.StoI(r.FormValue("TotalValue"))
		prepProject.TotalValue = totalValueS
	} else {
		return prepProject, fmt.Errorf("No TotalValue")
	}

	if r.FormValue("Location") != "" {
		prepProject.Location = r.FormValue("Location")
	} else {
		return prepProject, fmt.Errorf("No Location")
	}

	prepProject.MoneyRaised = 0

	if r.FormValue("Metadata") != "" {
		prepProject.Metadata = r.FormValue("Metadata")
	}

	// set the codes later while setting up stuff, need rpc calls for those as well
	prepProject.BalLeft = float64(0)
	prepProject.DateInitiated = utils.Timestamp()
	return prepProject, nil
}

func insertProject() {
	// this should be a post method since you want to accetp an project and then insert
	// that into the database
	http.HandleFunc("/insert/project", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkPost(w, r)
		var prepContract database.Project
		var prepProject database.DBParams
		prepProject, err := parseProject(r)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		log.Println("Prepared Project:", prepProject)
		prepContract.Params = prepProject
		err = prepContract.Save()
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

func parseInvestor(r *http.Request) (database.Investor, error) {
	// we need to create an instance of the Project
	// Index int auto
	// Name string required
	// PublicKey string optional, need bool "gen"
	// Seed string optional, need bool "gen"
	// AmountInvested float64 0 on init
	// FirstSignedUp string auto
	// InvestedAssets []Project don't set
	// LoginUserName string required
	// LoginPassword string required

	var prepInvestor database.Investor
	err := r.ParseForm()
	if err != nil {
		return prepInvestor, err
	}
	// if we're inserting this in, we need to get the next index number
	// so that we can set this without causing some weird bugs
	allInvestors, err := database.RetrieveAllInvestors()
	if err != nil {
		return prepInvestor, err
	}
	prepInvestor.U.Index = len(allInvestors) + 1
	if r.FormValue("LoginUserName") != "" {
		prepInvestor.U.LoginUserName = r.FormValue("LoginUserName")
	} else {
		// no username, error out
		return prepInvestor, fmt.Errorf("No LoginUserName")
	}

	if r.FormValue("LoginPassword") != "" {
		prepInvestor.U.LoginPassword = r.FormValue("LoginPassword")
	} else {
		// no password, error out
		return prepInvestor, fmt.Errorf("No LoginPassword")
	}

	if r.FormValue("Name") != "" {
		prepInvestor.U.Name = r.FormValue("Name")
	} else {
		return prepInvestor, fmt.Errorf("No Name")
	}

	if r.FormValue("gen") == "true" {
		// we need to generate a seed and pk pair
		pair, err := keypair.Random()
		if err != nil {
			return prepInvestor, fmt.Errorf("Error while generating keypair")
		}
		prepInvestor.U.Seed = pair.Seed()
		prepInvestor.U.PublicKey = pair.Address()
	}

	prepInvestor.AmountInvested = float64(0)
	prepInvestor.U.FirstSignedUp = utils.Timestamp()
	log.Println("Prepared investor: ", prepInvestor)
	return prepInvestor, nil
}

func insertInvestor() {
	// this should be a post method since you want to accetp an project and then insert
	// that into the database
	http.HandleFunc("/investor/insert", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkPost(w, r)
		var prepInvestor database.Investor
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

func investorPassword() {
	http.HandleFunc("/investor/name", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		var prepInvestor database.Investor
		// need to pass the pwhash param here
		if r.URL.Query() == nil || r.URL.Query()["LoginUserName"] == nil || len(r.URL.Query()["LoginUserName"][0]) != 128 { // sha 512 length
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		param := r.URL.Query()["LoginUserName"][0]
		log.Println("The pwhash is: ", param)
		// this is something like /investor/password?hash
		// so we need to remove the /investor/password part
		err := r.ParseForm()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		prepInvestor, err = database.ValidateInvestor(r.URL.Query()["LoginUserName"][0], "password") // TODO: Change this
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

func getAllInvestors() {
	http.HandleFunc("/investor/all", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		investors, err := database.RetrieveAllInvestors()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		log.Println("Retrieved all investors: ", investors)
		investorJson, err := json.Marshal(investors)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, investorJson)
	})
}

func parseRecipient(r *http.Request) (database.Recipient, error) {
	// Index int auto
	// Name string required
	// PublicKey string required
	// Seed string required
	// FirstSignedUp string auto timestamp
	// DebtAssets []string don't set
	// PaybackAssets []string don't set
	// LoginUserName string required
	// LoginPassword string required

	var prepRecipient database.Recipient
	err := r.ParseForm()
	if err != nil {
		return prepRecipient, err
	}

	allInvestors, err := database.RetrieveAllRecipients()
	if err != nil {
		return prepRecipient, err
	}

	prepRecipient.U.Index = len(allInvestors) + 1

	if r.FormValue("Name") != "" {
		prepRecipient.U.Name = r.FormValue("Name")
	} else {
		return prepRecipient, fmt.Errorf("No Name")
	}

	// we need to generate a seed and pk pair
	pair, err := keypair.Random()
	if err != nil {
		return prepRecipient, fmt.Errorf("Error while generating keypair")
	}
	prepRecipient.U.PublicKey = pair.Address()
	prepRecipient.U.Seed = pair.Seed()
	prepRecipient.U.FirstSignedUp = utils.Timestamp()

	if r.FormValue("LoginUserName") != "" {
		prepRecipient.U.LoginUserName = r.FormValue("LoginUserName")
	} else {
		// no username, error out
		return prepRecipient, fmt.Errorf("No LoginUserName")
	}

	if r.FormValue("LoginPassword") != "" {
		prepRecipient.U.LoginPassword = r.FormValue("LoginPassword")
	} else {
		// no password, error out
		return prepRecipient, fmt.Errorf("No LoginPassword")
	}

	log.Println("Prepared recipient: ", prepRecipient)
	return prepRecipient, nil
}

func getAllRecipient() {
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
		var prepRecipient database.Recipient
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

func recipientPassword() {
	http.HandleFunc("/recipient/password", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		var prepRecipient database.Recipient
		// need to pass the pwhash param here
		if r.URL.Query() == nil || r.URL.Query()["LoginUserName"] == nil || len(r.URL.Query()["LoginUserName"][0]) != 128 { // sha 512 length
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		param := r.URL.Query()["LoginUserName"][0]
		log.Println("The pwhash is: ", param)
		// this is something like /investor/password?hash
		// so we need to remove the /investor/password part
		err := r.ParseForm()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		prepRecipient, err = database.ValidateRecipient(r.URL.Query()["LoginUserName"][0], "password") // TODO: change this
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

// collect all hadnlers in one place so that we can aseemble them easily
// there are some repeating RPCs that we would like to avoid and maybe there's some
// nice way to group them together
// setupProjectRPCs sets up all the RPC calls related to projects that might be used
func setupProjectRPCs() {
	getOpenProjects()
	getProject()
	insertProject()
}

// setupInvestorRPCs sets up all RPCs related to the investor
func setupInvestorRPCs() {
	insertInvestor()
	investorPassword()
	// do we want an rpc that returns all investors for use in the backend?
	// right now adding it in but we can reomve this later if this is a feature that is not desired
	// TODO: add RPC to get a single investor from this list based on the index
	// a bigger question is do we index by number after all?
	// see TODO at investors.go for arguments for and against this
	getAllInvestors()
}

// setupBasicHandlerssets up two hadnler functions that can be used to serve a default
// 404 response when we either error out or received input is incorrect.  This is not
// exactly ideal, because we don't expcet the RPC to be exposed and would like some more
// errors when we handle it on the frontend, but this makes for more a bit more
// secure Frontedn implementation which doesn't leak any information to the frontend
func setupBasicHandlers() {
	setupDefaultHandler()
	setupPingHandler()
}

// setupRecipientRPCs sets up all RPCs related to the recipient. Most are similar
// to the investor RPCs, so maybe there's some nice way we can group them together
// to avoid code duplication
func setupRecipientRPCs() {
	getAllRecipient()
	insertRecipient()
	recipientPassword()
}

// StartServer runs on the server side ie the server with the frontend.
// having to define specific endpoints for this because this
// is the system that would be used by the backend, so has to be built secure.
func StartServer(port string) {
	// we have a couple sub handlers for each main handler. these handlers
	// call the relevant internal endpoints and return a status / data.
	// we also have to process data from the pi itself, and that should have its own
	// functions somewhere else that can be accessed by the rpc.

	// also, this is assumed to run on localhost and hence has no authentication mehcanism.
	// in the case we want to expose the API, we must add some stuff that secures this.
	// right now, its just the CORS header, since we want to allwo all localhost processes
	// to access the API
	setupBasicHandlers()
	// setup basic handlers - / and /ping
	setupProjectRPCs()
	// setup project related RPCs
	setupInvestorRPCs()
	// setup investor related RPCs
	setupRecipientRPCs()
	// setup recipient related RPCs
	// TODO: need to add recipient related RPCs
	portString := ":" + port // weird construction, but this should work
	// a potential improvement will be to add an authentication level like macaroons
	// so that we can serve over an authenticated channel.
	log.Fatal(http.ListenAndServe(portString, nil))
}
