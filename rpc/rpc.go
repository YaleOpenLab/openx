package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
)

// package rpc contains stuff that openx-cli and other third party platforms would use
// to interact with the opensolar platform

// we have a sub handler for each file in the subrepo. These handlers
// call the relevant internal endpoints

// API documentation over at the apidocs repo

// StartServer starts the server on the passed port and mode (http / https)
func StartServer(portx int, insecure bool) {
	erpc.SetupBasicHandlers()
	setupUserRpcs()
	setupStableCoinRPCs()
	setupPublicRoutes()
	setupAnchorHandlers()
	setupCAHandlers()
	adminHandlers()
	setupPlatformRoutes()

	port, err := utils.ToString(portx)
	if err != nil {
		log.Fatal("Port not string")
	}

	log.Println("Starting RPC Server on Port: ", port)
	if insecure {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(":"+port, "certs/server.crt", "certs/server.key", nil))
	}
}
