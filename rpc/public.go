package rpc

import (
	"log"
	"net/http"

	database "github.com/YaleOpenLab/openx/database"
)

// SnInvestor defines a sanitized investor
type SnInvestor struct {
	Name                  string
	InvestedSolarProjects []string
	AmountInvested        float64
	InvestedBonds         []string
	InvestedCoops         []string
	PublicKey             string
	Reputation            float64
}

// SnRecipient defines a sanitized recipient
type SnRecipient struct {
	Name                  string
	PublicKey             string
	ReceivedSolarProjects []string
	Reputation            float64
}

// SnUser defines a sanitized user
type SnUser struct {
	Name       string
	PublicKey  string
	Reputation float64
}

func setupPublicRoutes() {
	getAllInvestorsPublic()
	getAllRecipientsPublic()
	getTopReputationPublic()
	getInvTopReputationPublic()
}

// public contains all the RPC routes that we explicitly intend to make public. Other
// routes such as the invest route are things we could make private as well, but that
// doesn't change the security model since we ask for username+pwauth

// sanitizeInvestor removes sensitive fields frm the investor struct in order to be able
// to return the investor field in a public route
func sanitizeInvestor(investor database.Investor) SnInvestor {
	// this is a public route, so we shouldn't ideally return all parameters that are present
	// in the investor struct
	var sanitize SnInvestor
	sanitize.Name = investor.U.Name
	sanitize.InvestedSolarProjects = investor.InvestedSolarProjects
	sanitize.AmountInvested = investor.AmountInvested
	sanitize.InvestedBonds = investor.InvestedBonds
	sanitize.InvestedCoops = investor.InvestedCoops
	sanitize.PublicKey = investor.U.PublicKey
	sanitize.Reputation = investor.U.Reputation
	return sanitize
}

// sanitizeRecipient removes sensitive fields from the recipient struct in order to be
// able to return the recipient fields in a public route
func sanitizeRecipient(recipient database.Recipient) SnRecipient {
	// this is a public route, so we shouldn't ideally return all parameters that are present
	// in the investor struct
	var sanitize SnRecipient
	sanitize.Name = recipient.U.Name
	sanitize.PublicKey = recipient.U.PublicKey
	sanitize.Reputation = recipient.U.Reputation
	sanitize.ReceivedSolarProjects = recipient.ReceivedSolarProjects
	return sanitize
}

// sanitizeAllInvestors sanitizes an array of investors
func sanitizeAllInvestors(investors []database.Investor) []SnInvestor {
	var arr []SnInvestor
	for _, elem := range investors {
		arr = append(arr, sanitizeInvestor(elem))
	}
	return arr
}

// sanitizeUser sanitizes a particular user
func sanitizeUser(user database.User) SnUser {
	var sanitize SnUser
	sanitize.Name = user.Name
	sanitize.PublicKey = user.PublicKey
	sanitize.Reputation = user.Reputation
	return sanitize
}

// sanitizeAllRecipients sanitizes an array of recipients
func sanitizeAllRecipients(recipients []database.Recipient) []SnRecipient {
	var arr []SnRecipient
	for _, elem := range recipients {
		arr = append(arr, sanitizeRecipient(elem))
	}
	return arr
}

// sanitizeAllUsers sanitizes an arryay of users
func sanitizeAllUsers(users []database.User) []SnUser {
	var arr []SnUser
	for _, elem := range users {
		arr = append(arr, sanitizeUser(elem))
	}
	return arr
}

// getAllInvestors gets a list of all the investors in the system so that we can
// display it to some entity that is interested to view such stats
func getAllInvestorsPublic() {
	http.HandleFunc("/public/investor/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		investors, err := database.RetrieveAllInvestors()
		if err != nil {
			log.Println("did not retrieve all investors", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		sInvestors := sanitizeAllInvestors(investors)
		MarshalSend(w, r, sInvestors)
	})
}

// getAllRecipients gets a list of all the investors in the system so that we can
// display it to some entity that is interested to view such stats
func getAllRecipientsPublic() {
	http.HandleFunc("/public/recipient/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		recipients, err := database.RetrieveAllRecipients()
		if err != nil {
			log.Println("did not retrieve all recipients", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		sRecipients := sanitizeAllRecipients(recipients)
		MarshalSend(w, r, sRecipients)
	})
}

// this is to publish a list of the users with the best feedback in the system in order
// to award them badges or something similar
func getTopReputationPublic() {
	http.HandleFunc("/public/reputation/top", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		allUsers, err := database.TopReputationUsers()
		if err != nil {
			log.Println("did not retrive all top reputation users", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		sUsers := sanitizeAllUsers(allUsers)
		MarshalSend(w, r, sUsers)
	})
}

// getRecpTopReputationPublic gets a list of the recipients who have the best reputation on the platform
func getRecpTopReputationPublic() {
	http.HandleFunc("/public/recipient/reputation/top", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		allRecps, err := database.TopReputationRecipient()
		if err != nil {
			log.Println("did not retrieve all top reputaiton recipients", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		sRecipients := sanitizeAllRecipients(allRecps)
		MarshalSend(w, r, sRecipients)
	})
}

// getInvTopReputationPublic gets a lsit of the investors who have the best reputation on the platform
func getInvTopReputationPublic() {
	http.HandleFunc("/public/investor/reputation/top", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		allInvs, err := database.TopReputationInvestors()
		if err != nil {
			log.Println("did not retrieve all top reputation investors", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		sInvestors := sanitizeAllInvestors(allInvs)
		MarshalSend(w, r, sInvestors)
	})
}
