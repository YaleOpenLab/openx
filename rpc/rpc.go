package rpc

// the rpc package contains functions related to the server which will be interacting
// with the frontend. Not expanding on this too much since this will be changing quite often
// also evaluate on how easy it would be to rewrite this in nodeJS since the
// frontend is in react. Not many advantages per se and this works fine, so I guess
// we'll stay with this one for a while
import (
	"encoding/json"
	//"fmt"
	"log"
	"net/http"
)

// TODO: have some nice API documentation page for this so that we can easily reference

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
	StatusInternalServerError = http.StatusInternalServerError //  RFC 7231, 6.6.1
	StatusBadGateway          = http.StatusBadGateway          //  RFC 7231, 6.6.3
	StatusLocked              = http.StatusLocked              //  423 RFC 4918, 11.3
	StatusTooManyRequests     = http.StatusTooManyRequests     //  RFC 6585, 4
	StatusGatewayTimeout      = http.StatusGatewayTimeout      //  RFC 7231, 6.6.5
	StatusNotAcceptable       = http.StatusNotAcceptable       // RFC 7231, 6.5.6
	StatusServiceUnavailable  = http.StatusServiceUnavailable  //  RFC 7231, 6.6.4
)

func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Authorization, Cache-Control, Content-Type")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func MarshalSend(w http.ResponseWriter, r *http.Request, x interface{}) {
	xJson, err := json.Marshal(x)
	if err != nil {
		responseHandler(w, r, StatusInternalServerError)
		return
	}
	WriteToHandler(w, xJson)
}

func checkOrigin(w http.ResponseWriter, r *http.Request) {
	// re-enable this function for all private routes
	// if r.Header.Get("Origin") != "localhost" { // allow only our frontend UI to connect to our RPC instance
	// 	http.Error(w, "404 page not found", http.StatusNotFound)
	// }
}

func checkGet(w http.ResponseWriter, r *http.Request) {
	checkOrigin(w, r)
	if r.Method != "GET" {
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}

func checkPost(w http.ResponseWriter, r *http.Request) {
	checkOrigin(w, r)
	log.Println("Checking POST")
	if r.Method != "POST" {
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}

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

func setupDefaultHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// default to 404 for every application not running on localhost
		responseHandler(w, r, StatusNotFound)
		return
	})
}

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

	portString := ":" + port // weird construction, but this should work
	log.Fatal(http.ListenAndServe(portString, nil))
}
