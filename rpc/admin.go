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
	listAllAdmins()
	addNewPlatform()
	sendNewMessage()
}

// KillCode is a code that can immediately shut down the server in case of hacks / crises
var KillCode string

// validateAdmin validates whether a given user is an admin and returns a bool
func validateAdmin(w http.ResponseWriter, r *http.Request, options ...string) (database.User, bool) {
	prepUser, err := userValidateHelper(w, r, options)
	if err != nil {
		log.Println(err)
		return prepUser, false
	}

	if !prepUser.Admin {
		erpc.ResponseHandler(w, erpc.StatusUnauthorized)
		return prepUser, false
	}

	return prepUser, true
}

// killServer instantly kills the server. Recovery possible only with server access
func killServer() {
	http.HandleFunc("/admin/kill", func(w http.ResponseWriter, r *http.Request) {
		log.Println("kill command received")
		// need to pass the pwhash param here
		_, adminBool := validateAdmin(w, r, "nuke", "username")
		if !adminBool {
			return
		}

		if r.FormValue("username") == "martin" {
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
		_, adminBool := validateAdmin(w, r)
		if !adminBool {
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
		_, adminBool := validateAdmin(w, r, "username")
		if !adminBool {
			return
		}

		if r.FormValue("username") == "martin" {
			// only authorized users, can change at compile time
			log.Println("generating new nuclear code")
			KillCode = utils.GetRandomString(64)
			w.Write([]byte(KillCode))
		} else {
			erpc.MarshalSend(w, erpc.StatusUnauthorized)
			return
		}
	})
}

// newPlatform creates a new platform code
func newPlatform() {
	http.HandleFunc("/admin/platform/new", func(w http.ResponseWriter, r *http.Request) {
		// need to pass the pwhash param here
		_, adminBool := validateAdmin(w, r, "name", "code", "timeout")
		if !adminBool {
			return
		}

		name := r.FormValue("name")
		code := r.FormValue("code")
		timeout := r.FormValue("timeout") // if specified, timeout is false

		var timeoutBool bool

		if timeout != "false" {
			timeoutBool = true
		}

		log.Println("Creating new platform code: ", code, " for: ", name, " with timeout: ", timeoutBool)

		err := database.NewPlatform(name, code, timeoutBool)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}
		return
	})
}

// retrieveAllPlatforms retrieves all platforms from the database
func retrieveAllPlatforms() {
	http.HandleFunc("/admin/platform/all", func(w http.ResponseWriter, r *http.Request) {
		_, adminBool := validateAdmin(w, r)
		if !adminBool {
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

// listAllAdmins lists all the admin users of openx so users can contact them in case they face any problem
func listAllAdmins() {
	http.HandleFunc("/admin/list", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		admins, err := database.RetrieveAllAdmins()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		erpc.MarshalSend(w, admins)
	})
}

func addNewPlatform() {
	http.HandleFunc("/admin/add/platform", func(w http.ResponseWriter, r *http.Request) {
		_, adminBool := validateAdmin(w, r, "name", "code", "timeout")
		if !adminBool {
			return
		}

		var err error
		name := r.FormValue("name")
		code := r.FormValue("code")
		timeout := r.FormValue("timeout")

		if name == "" || code == "" {
			log.Println("code or desired name empty")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		if timeout == "false" {
			err = database.NewPlatform(name, code, false)
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
		} else {
			err = database.NewPlatform(name, code, true)
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func sendNewMessage() {
	http.HandleFunc("/admin/sendmessage", func(w http.ResponseWriter, r *http.Request) {
		_, adminBool := validateAdmin(w, r, "subject", "message", "recipient")
		if !adminBool {
			return
		}

		subject := r.FormValue("subject")
		message := r.FormValue("message")
		recipient := r.FormValue("recipient")

		if len(subject) == 0 {
			log.Println("length of subject is zero, not sending message")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if len(message) == 0 {
			log.Println("length of message is zero, not sending message")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		user, err := database.CheckUsernameCollision(recipient)
		if err == nil { // ie if there is no user
			log.Println("there is no user with the given username")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = user.AddtoMailbox(subject, message)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
