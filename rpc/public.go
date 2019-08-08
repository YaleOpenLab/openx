package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
)

// SnUser defines a sanitized user
type SnUser struct {
	Name       string
	PublicKey  string
	Reputation float64
}

func setupPublicRoutes() {
	getTopReputationPublic()
	getUserInfo()
}

// public contains all the RPC routes that we explicitly intend to make public. Other
// routes such as the invest route are things we could make private as well, but that
// doesn't change the security model since we ask for username+pwauth

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

// this is to publish a list of the users with the best feedback in the system in order
// to award them badges or something similar
func getTopReputationPublic() {
	http.HandleFunc("/public/reputation/top", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		allUsers, err := database.TopReputationUsers()
		if err != nil {
			log.Println("did not retrive all top reputation users", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		sUsers := sanitizeAllUsers(allUsers)
		erpc.MarshalSend(w, sUsers)
	})
}

func getUserInfo() {
	http.HandleFunc("/public/user", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if r.URL.Query()["index"] == nil {
			log.Println("no index retrieved, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		index, err := utils.ToInt(r.URL.Query()["index"][0])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		user, err := database.RetrieveUser(index)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		sUser := sanitizeUser(user)
		erpc.MarshalSend(w, sUser)
	})
}

func mainnetRPC() {
	http.HandleFunc("/mainnet", func(w http.ResponseWriter, r *http.Request) {
		// set a single byte response for mainnet / testnet
		// mainnet is 0, testnet is 1
		mainnet := []byte{0}
		testnet := []byte{1}
		if consts.Mainnet {
			w.Write(mainnet)
		} else {
			w.Write(testnet)
		}
		return
	})
}
