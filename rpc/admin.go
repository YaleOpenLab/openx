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

// AdminRPC is the list of all admin RPC endpoints
var AdminRPC = map[int][]string{
	1:  {"/admin/kill", "POST", "nuke"},                                   // POST
	2:  {"/admin/freeze", "GET"},                                          // GET
	3:  {"/admin/gennuke", "POST"},                                        // POST
	5:  {"/admin/platform/all", "GET"},                                    // GET
	6:  {"/admin/list"},                                                   // GET
	7:  {"/admin/platform/new", "POST", "name", "code", "timeout"},        // POST
	8:  {"/admin/sendmessage", "POST", "subject", "message", "recipient"}, // POST
	9:  {"/admin/getallusers", "GET"},                                     // GET
	10: {"/admin/userverify", "POST", "index"},                            // POST
	11: {"/admin/userunverify", "POST", "index"},                          // POST
}

// adminHandlers are a list of all the admin handlers defined by openx
func adminHandlers() {
	killServer()
	freezeServer()
	genNuclearCode()
	retrieveAllPlatforms()
	listAllAdmins()
	addNewPlatform()
	sendNewMessage()
	getallUsersAdmin()
	verifyUser()
	unverifyUser()
}

// KillCode is a code that can immediately shut down the server in case of hacks / crises
var KillCode string

// validateAdmin validates whether a given user is an admin and returns a bool
func validateAdmin(w http.ResponseWriter, r *http.Request, options []string, method string) (database.User, bool) {
	prepUser, err := userValidateHelper(w, r, options, method)
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
	http.HandleFunc(AdminRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		log.Println("kill command received")
		// need to pass the pwhash param here
		_, adminBool := validateAdmin(w, r, AdminRPC[1][2:], AdminRPC[1][1])
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
	http.HandleFunc(AdminRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		// need to pass the pwhash param here
		_, adminBool := validateAdmin(w, r, []string{}, AdminRPC[2][1])
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
	http.HandleFunc(AdminRPC[3][0], func(w http.ResponseWriter, r *http.Request) {
		// need to pass the pwhash param here
		_, adminBool := validateAdmin(w, r, AdminRPC[3][2:], AdminRPC[3][1])
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

// retrieveAllPlatforms retrieves all platforms from the database
func retrieveAllPlatforms() {
	http.HandleFunc(AdminRPC[5][0], func(w http.ResponseWriter, r *http.Request) {
		_, adminBool := validateAdmin(w, r, []string{}, AdminRPC[5][1])
		if !adminBool {
			return
		}

		pfs, err := database.RetrieveAllPlatforms()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, pfs)
	})
}

// listAllAdmins lists all the admin users of openx so users can contact them in case they face any problem
func listAllAdmins() {
	http.HandleFunc(AdminRPC[6][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		admins, err := database.RetrieveAllAdmins()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, admins)
	})
}

func addNewPlatform() {
	http.HandleFunc(AdminRPC[7][0], func(w http.ResponseWriter, r *http.Request) {
		_, adminBool := validateAdmin(w, r, AdminRPC[7][2:], AdminRPC[7][1])
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
			if erpc.Err(w, err, erpc.StatusInternalServerError) {
				return
			}
		} else {
			err = database.NewPlatform(name, code, true)
			if erpc.Err(w, err, erpc.StatusInternalServerError) {
				return
			}
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func sendNewMessage() {
	http.HandleFunc(AdminRPC[8][0], func(w http.ResponseWriter, r *http.Request) {
		_, adminBool := validateAdmin(w, r, AdminRPC[8][2:], AdminRPC[8][1])
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
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

type getUsersAdmin struct {
	Length int
}

func getallUsersAdmin() {
	http.HandleFunc(AdminRPC[9][0], func(w http.ResponseWriter, r *http.Request) {
		_, adminBool := validateAdmin(w, r, AdminRPC[9][2:], AdminRPC[9][1])
		if !adminBool {
			return
		}

		users, err := database.RetrieveAllUsers()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		var x getUsersAdmin
		x.Length = len(users)

		erpc.MarshalSend(w, x)
	})
}

func verifyUser() {
	http.HandleFunc(AdminRPC[10][0], func(w http.ResponseWriter, r *http.Request) {
		admin, adminBool := validateAdmin(w, r, AdminRPC[10][2:], AdminRPC[10][1])
		if !adminBool {
			return
		}

		indexS := r.FormValue("index")

		index, err := utils.ToInt(indexS)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		user, err := database.RetrieveUser(index)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		user.Verified = true
		user.VerifiedBy = admin.Index
		user.VerifiedTime = utils.Timestamp()

		err = user.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func unverifyUser() {
	http.HandleFunc(AdminRPC[11][0], func(w http.ResponseWriter, r *http.Request) {
		admin, adminBool := validateAdmin(w, r, AdminRPC[11][2:], AdminRPC[11][1])
		if !adminBool {
			return
		}

		indexS := r.FormValue("index")

		index, err := utils.ToInt(indexS)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		user, err := database.RetrieveUser(index)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		user.Verified = false
		user.VerifiedBy = admin.Index
		user.VerifiedTime = utils.Timestamp()

		err = user.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
