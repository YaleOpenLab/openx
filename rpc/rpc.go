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

// use these standard error codes to send out to request replies so callers can figure
// out what's going wrong with their requests
const (
	StatusOK                  = http.StatusOK                  //  200 RFC 7231, 6.3.1
	StatusCreated             = http.StatusCreated             //  201 RFC 7231, 6.3.2
	StatusMovedPermanently    = http.StatusMovedPermanently    //  301 RFC 7231, 6.4.2
	StatusBadRequest          = http.StatusBadRequest          //  400 RFC 7231, 6.5.1
	StatusUnauthorized        = http.StatusUnauthorized        //  401 RFC 7235, 3.1
	StatusPaymentRequired     = http.StatusPaymentRequired     //  402 RFC 7231, 6.5.2
	StatusNotFound            = http.StatusNotFound            //  404 RFC 7231, 6.5.4
	StatusInternalServerError = http.StatusInternalServerError //  500 RFC 7231, 6.6.1
	StatusBadGateway          = http.StatusBadGateway          //  502 RFC 7231, 6.6.3
	StatusLocked              = http.StatusLocked              //  423 RFC 4918, 11.3
	StatusTooManyRequests     = http.StatusTooManyRequests     //  429 RFC 6585, 4
	StatusGatewayTimeout      = http.StatusGatewayTimeout      //  504 RFC 7231, 6.6.5
	StatusNotAcceptable       = http.StatusNotAcceptable       //  406 RFC 7231, 6.5.6
	StatusServiceUnavailable  = http.StatusServiceUnavailable  //  503 RFC 7231, 6.6.4
)

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
func GetRequest(url string) ([]byte, error) {
	// make a curl request out to lcoalhost and get the ping response
	var dummy []byte
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dummy, err
	}
	req.Header.Set("Origin", "localhost")
	res, err := client.Do(req)
	if err != nil {
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
		return dummy, err
	}
	// need to add this header or we'll get a negative response
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return dummy, err
	}

	defer res.Body.Close()
	x, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return dummy, err
	}

	return x, nil
}

// responseHandler is teh default response handler that sends out response codes on successful
// completion of certain calls
func responseHandler(w http.ResponseWriter, r *http.Request, status int) {
	var response StatusResponse
	response.Code = status
	switch status {
	case StatusOK:
		response.Status = "OK"
	case StatusCreated:
		response.Status = "Method Created"
	case StatusMovedPermanently:
		response.Status = "Endpoint moved permanently"
	case StatusBadRequest:
		response.Status = "Bad Request error!"
	case StatusUnauthorized:
		response.Status = "You are unauthorized to make this request"
	case StatusPaymentRequired:
		response.Status = "Payment required before you can access this endpoint"
	case StatusNotFound:
		response.Status = "404 Error Not Found!"
	case StatusInternalServerError:
		response.Status = "Internal Server Error"
	case StatusLocked:
		response.Status = "Endpoint locked until further notice"
	case StatusTooManyRequests:
		response.Status = "Too many requests made, try again later"
	case StatusBadGateway:
		response.Status = "Bad Gateway Error"
	case StatusServiceUnavailable:
		response.Status = "Service Unavailable error"
	case StatusGatewayTimeout:
		response.Status = "Gateway Timeout Error"
	case StatusNotAcceptable:
		response.Status = "Not accepted"
	default:
		response.Status = "404 Page Not Found"
	}
	MarshalSend(w, r, response)
}

// setupDefaultHandler sets up the default handler (ie returns 404 for invalid routes)
func setupDefaultHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// default to 404 for every application not running on localhost
		responseHandler(w, r, StatusNotFound)
		return
	})
}

// setupPingHandler is a ping route for remote callers to check if the platform is up
func setupPingHandler() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
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
