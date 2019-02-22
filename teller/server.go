package main

import (
	"encoding/json"
	//"fmt"
	"log"
	"net/http"

	rpc "github.com/YaleOpenLab/openx/rpc"
	utils "github.com/YaleOpenLab/openx/utils"
)

// server starts a local server which would inform us about the uptime of the teller and provide a data endpoint

// Data defines a structure to be used while capturing data
type Data struct {
	// the data that is going to be streamed
	// TODO: define what goes in here
	Timestamp string
	Info      string
}

func checkGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}

func responseHandler(w http.ResponseWriter, r *http.Request, status int) {
	var response rpc.StatusResponse
	response.Code = status
	switch status {
	case rpc.StatusOK:
		response.Status = "OK"
	case rpc.StatusNotFound:
		response.Status = "404 Error Not Found!"
	case rpc.StatusInternalServerError:
		response.Status = "Internal Server Error"
	}
	rpc.MarshalSend(w, r, response)
}

// setupDefaultHandler sets up the default handler (ie returns 404 for invalid routes)
func setupDefaultHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		responseHandler(w, r, rpc.StatusNotFound)
		return
	})
}

// pingHandler can be used on the frontend to try checking whether the teller is still up.
// maybe have a button or something and pressing that would call this endpoint
func pingHandler() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		var pr rpc.StatusResponse
		pr.Code = 200
		pr.Status = "HEALTH OK"
		prJson, err := json.Marshal(pr)
		if err != nil {
			responseHandler(w, r, rpc.StatusInternalServerError)
			return
		}
		WriteToHandler(w, prJson)
	})
}

// HCHeaderResponse defines the hash chain header's response
type HCHeaderResponse struct {
	Hash string
}

// hashChainHeaderHandler returns the header of the ipfs hash chain
// clients who want historicasl record of all activities can record the latest hash
// and then derive all the other files from it. This avoids a need for a direct endpoint
// that will serve data directly while leveraging ipfs.
func hashChainHeaderHandler() {
	http.HandleFunc("/hash", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		var x HCHeaderResponse
		x.Hash = HashChainHeader
		xJson, err := json.Marshal(x)
		if err != nil {
			responseHandler(w, r, rpc.StatusInternalServerError)
			return
		}
		WriteToHandler(w, xJson)
	})
}

func setupRoutes() {
	setupDefaultHandler()
	pingHandler()
	hashChainHeaderHandler()
}

// curl https://localhost/ping --insecure {"Code":200,"Status":""}
// TODO: have an alternative for using insecure certificates here.
func startServer(port int) {
	setupRoutes()
	err := http.ListenAndServeTLS(":"+utils.ItoS(port), "ssl/server.crt", "ssl/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
