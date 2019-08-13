package rpc

import (
	"log"
	"net/http"
	"os"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
)

// admin contains a list of all the functions that will hopefully never be used in practice
// but if needed are incredibly powerful

// adminHandlers are a list of all the admin handlers defined by openx
func adminHandlers() {
	killServer()
	freezeServer()
	genNuclearCode()
	newPlatform()
	retrieveAllPlatforms()
}

// KillCode is a code that can immediately shut down the server in case of hacks / crises
var KillCode string

// validateAdmin validates whether a given user is an admin and returns a bool
func validateAdmin(w http.ResponseWriter, r *http.Request, options ...string) bool {
	erpc.CheckGet(w, r)
	erpc.CheckOrigin(w, r)

	prepUser, err := CheckReqdParams(w, r, []string{})
	if err != nil {
		return false
	}
	if !prepUser.Admin {
		return false
	}

	for _, option := range options {
		if r.URL.Query()[option] == nil {
			return false
		}
	}

	return true
}

// killServer instantly kills the server. Recovery possible only with server access
func killServer() {
	http.HandleFunc("/admin/kill", func(w http.ResponseWriter, r *http.Request) {
		log.Println("kill command received")
		// need to pass the pwhash param here
		if !validateAdmin(w, r, "nuke", "username") {
			// admin account not accessible
			if r.URL.Query()["nuke"][0] == KillCode {
				log.Println("nuclear code activated, killing server")
				os.Exit(1)
			}
		}

		if r.URL.Query()["username"][0] == "martin" {
			// only certain admins can access this endpoint, can be compiled at runtime
			log.Println("Activating kill switch")
			os.Exit(1)
		}
	})
}

// freezeServer freezes the server to make all transactions void. The easiest way to do that
// is to set the Mainnet const to false.
func freezeServer() {
	http.HandleFunc("/admin/freeze", func(w http.ResponseWriter, r *http.Request) {
		// need to pass the pwhash param here
		if !validateAdmin(w, r) {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		consts.SetConsts(false) // runtime const migration
		log.Println("Server frozen, state reverted to mainnet. Restart server to unfreeze")
	})
}

// genNuclearCode generates a nuclear code capable of instantly killing the platform. Can only
// be called by certain admins
func genNuclearCode() {
	http.HandleFunc("/admin/gennuke", func(w http.ResponseWriter, r *http.Request) {
		// need to pass the pwhash param here
		if !validateAdmin(w, r, "username") {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		if r.URL.Query()["username"][0] == "martin" {
			// only authorized users, can change at compile time
			log.Println("generating new nuclear code")
			KillCode = utils.GetRandomString(64)
			w.Write([]byte(KillCode))
		}
	})
}

// newPlatform creates a new platform code
func newPlatform() {
	http.HandleFunc("/admin/platform/new", func(w http.ResponseWriter, r *http.Request) {
		// need to pass the pwhash param here
		if !validateAdmin(w, r, "name", "code") {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		name := r.URL.Query()["name"][0]
		code := r.URL.Query()["code"][0]

		log.Println("Creating new platform code: ", code, " for: ", name)

		err := database.NewPlatform(name, code)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}
		return
	})
}

// retrieveAllPlatforms retrieves all platforms from the database
func retrieveAllPlatforms() {
	http.HandleFunc("/admin/platform/all", func(w http.ResponseWriter, r *http.Request) {
		if !validateAdmin(w, r) {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		pfs, err := database.RetrieveAllPlatforms()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, pfs)
	})
}
