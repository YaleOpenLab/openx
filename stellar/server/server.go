package server

import (
	"fmt"
	"log"
	// "net"
	"encoding/json"
	"net/http"
	// "net/rpc"
)

type Dummy struct {
	Username string
	Password string
}

type PingResponse struct {
	Status string
}

func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Header().Add("Access-Control-Allow-Origin", "localhost")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonString)
}

func CheckOrigin(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Origin") != "localhost" {
		// allow only our frontend UI to connect to our RPC instance
		http.Error(w, "404 page not found", http.StatusNotFound)
		return fmt.Errorf("Cross domain request error")
	}
	return nil
}

func CheckGet(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
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
		err := CheckOrigin(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		err = CheckGet(w, r)
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

func setupTestHandler() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		CheckOrigin(w, r)
		CheckGet(w, r)
		var dummy Dummy
		log.Println("remote peer calling /cool endpoint")
		dummy.Username = "dummy"
		dummy.Password = "123"
		dummyJson, err := json.Marshal(dummy)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, dummyJson)
	})
}

func SetupServer() {
	// this runs on the server side ie the server with the frontend.
	// having to define specific endpoints for this because this
	// is the system that would be used by the backend, so has to be secure.
	// good thing to lock these calls to localhost as well
	setupDefaultHandler()
	setupPingHandler()
	setupTestHandler()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
