package rpc

import (
	"log"
	"net/http"
	"os"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
)

// admin contains a list of all the functions that will hopefully never be used in practice
// but if needed are incredibly powerful

func adminHandlers() {
	killServer()
	freezeServer()
	genNuclearCode()
}

var KillCode string

func validateAdmin(w http.ResponseWriter, r *http.Request) bool {
	erpc.CheckGet(w, r)
	erpc.CheckOrigin(w, r)

	prepUser, err := CheckReqdParams(w, r)
	if err != nil {
		return false
	}
	if !prepUser.Admin {
		return false
	}
	return true
}

// killServer instantly kills the server. Recovery possible only with server access
func killServer() {
	http.HandleFunc("/admin/kill", func(w http.ResponseWriter, r *http.Request) {
		log.Println("kill command received")
		// need to pass the pwhash param here
		if !validateAdmin(w, r) {
			// admin account not accessible
			if r.URL.Query()["nuke"] != nil {
				if r.URL.Query()["nuke"][0] == KillCode {
					log.Println("nuclear code activated, killing server")
					os.Exit(1)
				}
			} else {
				erpc.ResponseHandler(w, erpc.StatusUnauthorized)
				return
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

func genNuclearCode() {
	http.HandleFunc("/admin/gennuke", func(w http.ResponseWriter, r *http.Request) {
		// need to pass the pwhash param here
		if !validateAdmin(w, r) {
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
