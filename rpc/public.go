package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	database "github.com/YaleOpenLab/openx/database"
)

// public contains all the RPC routes that we explicitly intend to make public

// setupPublicRoutes sets up public routes that can be called by anyone without an
// acount on openx
func setupPublicRoutes() {
	getTopReputationPublic()
	getUserInfo()
}

// SnUser defines a sanitized user
type SnUser struct {
	Name       string
	PublicKey  string
	Reputation float64
}

// sanitizeUser sanitizes a particular user
func sanitizeUser(user database.User) SnUser {
	var sanitize SnUser
	sanitize.Name = user.Name
	sanitize.PublicKey = user.StellarWallet.PublicKey
	sanitize.Reputation = user.Reputation
	return sanitize
}

// sanitizeAllUsers sanitizes an arryay of users
func sanitizeAllUsers(users []database.User) []SnUser {
	var arr []SnUser
	for _, elem := range users {
		arr = append(arr, sanitizeUser(elem))
	}
	return arr
}

// getTopReputationPublic returns a list of users sorted by descending order of reputation
func getTopReputationPublic() {
	http.HandleFunc("/public/reputation/top", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		allUsers, err := database.TopReputationUsers()
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not retrive all top reputation users") {
			return
		}
		sUsers := sanitizeAllUsers(allUsers)
		erpc.MarshalSend(w, sUsers)
	})
}

// getUserInfo returns a list of sanitised users of the openx platform
func getUserInfo() {
	http.HandleFunc("/public/user", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		if r.URL.Query()["index"] == nil {
			log.Println("no index retrieved, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		index, err := utils.ToInt(r.URL.Query()["index"][0])
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		user, err := database.RetrieveUser(index)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		sUser := sanitizeUser(user)
		erpc.MarshalSend(w, sUser)
	})
}
