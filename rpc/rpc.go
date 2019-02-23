package rpc

// the rpc package contains functions related to the server which will be interacting
// with the frontend. Not expanding on this too much since this will be changing quite often
// also evaluate on how easy it would be to rewrite this in nodeJS since the
// frontend is in react. Not many advantages per se and this works fine, so I guess
// we'll stay with this one for a while
import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// API documentation over at the apidocs repo

// StatusResponse defines a generic status response structure
type StatusResponse struct {
	Code   int
	Status string
}

// setupBasicHandlers sets up two handler functions that can be used to serve a default
// 404 response when we either error out or received input is incorrect.  This is not
// exactly ideal, because we don't expcet the RPC to be exposed and would like some more
// errors when we handle it on the frontend, but this makes for more a bit more
// secure Frontedn implementation which doesn't leak any information to the frontend
func setupBasicHandlers() {
	setupDefaultHandler()
	setupPingHandler()
}

// WriteToHandler constructs a reply to the passed writer
func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Authorization, Cache-Control, Content-Type")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

// MarshalSend marshals and writes a json string into the writer
func MarshalSend(w http.ResponseWriter, r *http.Request, x interface{}) {
	xJson, err := json.Marshal(x)
	if err != nil {
		log.Println("did not marshal json", err)
		responseHandler(w, r, StatusInternalServerError)
		return
	}
	WriteToHandler(w, xJson)
}

// checkOrigin checks the origin of the incoming request
func checkOrigin(w http.ResponseWriter, r *http.Request) {
	// re-enable this function for all private routes
	// if r.Header.Get("Origin") != "localhost" { // allow only our frontend UI to connect to our RPC instance
	// 	http.Error(w, "404 page not found", http.StatusNotFound)
	// }
}

// checkGet checks if the invoming request is a GET request
func checkGet(w http.ResponseWriter, r *http.Request) {
	checkOrigin(w, r)
	if r.Method != "GET" {
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}

// checkPost checks whether the incomign request is a POST request
func checkPost(w http.ResponseWriter, r *http.Request) {
	checkOrigin(w, r)
	if r.Method != "POST" {
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}

// GetRequest is a handler that makes it easy to send out GET requests
// we don't set timeouts here because block times can be variable and a single request
// can sometimes take a long while to complete
func GetRequest(url string) ([]byte, error) {
	var dummy []byte
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("did not create new GET request", err)
		return dummy, err
	}
	req.Header.Set("Origin", "localhost")
	res, err := client.Do(req)
	if err != nil {
		log.Println("did not make request", err)
		return dummy, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// PutRequest is a handler that makes it easy to send out POST requests
func PutRequest(body string, payload io.Reader) ([]byte, error) {

	// the body must be the param that you usually pass to curl's -d option
	var dummy []byte
	req, err := http.NewRequest("PUT", body, payload)
	if err != nil {
		log.Println("did not create new PUT request", err)
		return dummy, err
	}
	// need to add this header or we'll get a negative response
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("did not make request", err)
		return dummy, err
	}

	defer res.Body.Close()
	x, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("did not read from ioutil", err)
		return dummy, err
	}

	return x, nil
}

// setupDefaultHandler sets up the default handler (ie returns 404 for invalid routes)
func setupDefaultHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// default to 404 for every application not running on localhost
		checkGet(w, r)
		checkOrigin(w, r)
		responseHandler(w, r, StatusNotFound)
		return
	})
}

// setupPingHandler is a ping route for remote callers to check if the platform is up
func setupPingHandler() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		responseHandler(w, r, StatusOK)
	})
}

// StartServer runs on the server side ie the server with the frontend.
// having to define specific endpoints for this because this
// is the system that would be used by the backend, so has to be built secure.
func StartServer(port string) {
	// we have a sub handlers for each major entity. These handlers
	// call the relevant internal endpoints and return a StatusResponse message.
	// we also have to process data from the pi itself, and that should have its own
	// functions somewhere else that can be accessed by the rpc.

	// also, this is assumed to run on localhost and hence has no authentication mehcanism.
	// in the case we want to expose the API, we must add some stuff that secures this.
	// right now, its just the CORS header, since we want to allow all localhost processes
	// to access the API
	// a potential improvement will be to add something like macaroons
	// so that we can serve over an authenticated channel
	// setup all related handlers
	setupBasicHandlers()
	setupProjectRPCs()
	setupInvestorRPCs()
	setupRecipientRPCs()
	setupBondRPCs()
	setupCoopRPCs()
	setupUserRpcs()
	setupStableCoinRPCs()
	setupPublicRoutes()
	setupEntityRPCs()
	setupParticleHandlers()

	portString := ":" + port // weird construction, but this should work
	log.Fatal(http.ListenAndServe(portString, nil))
}
