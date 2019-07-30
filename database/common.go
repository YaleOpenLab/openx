package database

import (
	"encoding/json"
	"github.com/pkg/errors"

	edb "github.com/Varunram/essentials/database"
	consts "github.com/YaleOpenLab/openx/consts"
)

// this file contains commong methods that are repeated across interfaces
// Save inserts a passed User object into the database
func (a *User) Save() error {
	return edb.Save(consts.DbDir, UserBucket, a, a.Index)
}

// Save inserts a passed Investor object into the database
func (a *Investor) Save() error {
	return edb.Save(consts.DbDir, InvestorBucket, a, a.U.Index)
}

// Save saves a given recipient's details
func (a *Recipient) Save() error {
	return edb.Save(consts.DbDir, RecipientBucket, a, a.U.Index)
}

// RetrieveUser retrieves a particular User indexed by key from the database
func RetrieveUser(key int) (User, error) {
	var user User
	x, err := edb.Retrieve(consts.DbDir, UserBucket, key)
	if err != nil {
		return user, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &user)
	return user, err
}

// RetrieveInvestor retrieves a particular investor indexed by key from the database
func RetrieveInvestor(key int) (Investor, error) {
	var inv Investor
	user, err := RetrieveUser(key)
	if err != nil {
		return inv, err
	}

	x, err := edb.Retrieve(consts.DbDir, InvestorBucket, key)
	if err != nil {
		return inv, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &inv)
	if err != nil {
		return inv, errors.Wrap(err, "could not unmarshal investor")
	}

	inv.U = &user
	return inv, inv.Save()
}

// RetrieveRecipient retrieves a specific recipient from the database
func RetrieveRecipient(key int) (Recipient, error) {
	var recp Recipient
	user, err := RetrieveUser(key)
	if err != nil {
		return recp, err
	}

	x, err := edb.Retrieve(consts.DbDir, RecipientBucket, key)
	if err != nil {
		return recp, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &recp)
	if err != nil {
		return recp, errors.New("could not unmarshal recipient")
	}

	recp.U = &user
	return recp, recp.Save()
}

// RetrieveAllUsers gets a list of all User in the database
func RetrieveAllUsers() ([]User, error) {
	var arr []User
	x, err := edb.RetrieveAllKeys(consts.DbDir, UserBucket)
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range x {
		var temp User
		err := json.Unmarshal(value, &temp)
		if err != nil {
			return arr, errors.New("error while unmarshalling json, quitting")
		}
		arr = append(arr, temp)
	}

	return arr, nil
}

// RetrieveAllUsers gets a list of all User in the database
func RetrieveAllInvestors() ([]Investor, error) {
	var arr []Investor

	allUsers, err := RetrieveAllUsers()
	if err != nil {
		return arr, errors.Wrap(err, "could not retrieve all users from db")
	}

	lim := len(allUsers)
	x, err := edb.RetrieveAllKeysLim(consts.DbDir, InvestorBucket, lim)
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range x {
		var temp Investor
		err := json.Unmarshal(value, &temp)
		if err != nil {
			return arr, errors.Wrap(err, "error while unmarshalling json, quitting")
		}
		if temp.U.Index != 0 {
			arr = append(arr, temp)
		}
	}

	return arr, nil
}

// RetrieveAllRecipients gets a list of all Recipients in the database
func RetrieveAllRecipients() ([]Recipient, error) {
	var arr []Recipient

	allUsers, err := RetrieveAllUsers()
	if err != nil {
		return arr, errors.Wrap(err, "could not retrieve all users from db")
	}

	lim := len(allUsers)
	x, err := edb.RetrieveAllKeysLim(consts.DbDir, RecipientBucket, lim)
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range x {
		var temp Recipient
		err := json.Unmarshal(value, &temp)
		if err != nil {
			return arr, errors.Wrap(err, "error while unmarshalling json, quitting")
		}
		if temp.U.Index != 0 {
			arr = append(arr, temp)
		}
	}

	return arr, nil
}

// TopReputationUsers gets the users with top reputation
func TopReputationUsers() ([]User, error) {
	// these reputation functions should mostly be used by the frontend through the
	// RPC to display to other users what other users' reputation is.
	arr, err := RetrieveAllUsers()
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users from database")
	}
	for i := range arr {
		for j := range arr {
			if arr[i].Reputation > arr[j].Reputation {
				tmp := arr[i]
				arr[i] = arr[j]
				arr[j] = tmp
			}
		}
	}
	return arr, nil
}

// TopReputationInvestors gets a list of all the investors with top reputation
func TopReputationInvestors() ([]Investor, error) {
	arr, err := RetrieveAllInvestors()
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users from database")
	}
	for i := range arr {
		for j := range arr {
			if arr[i].U.Reputation > arr[j].U.Reputation {
				tmp := arr[i]
				arr[i] = arr[j]
				arr[j] = tmp
			}
		}
	}
	return arr, nil
}

// TopReputationRecipient returns a list of recipients with the best reputation
func TopReputationRecipients() ([]Recipient, error) {
	arr, err := RetrieveAllRecipients()
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users from database")
	}
	for i := range arr {
		for j := range arr {
			if arr[i].U.Reputation > arr[j].U.Reputation {
				tmp := arr[i]
				arr[i] = arr[j]
				arr[j] = tmp
			}
		}
	}
	return arr, nil
}

// ValidateUser validates a particular user
func ValidateUser(name string, pwhash string) (User, error) {
	var dummy User
	users, err := RetrieveAllUsers()
	if err != nil {
		return dummy, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		if user.Username == name && user.Pwhash == pwhash {
			return user, nil
		}
	}
	return dummy, errors.New("could not find user with requested credentials")
}

// ValidateInvestor is a function to validate the investors username and password to log them into the platform, and find the details related to the investor
// This is separate from the publicKey/seed pair (which are stored encrypted in the database); since we can help users change their password, but we can't help them retrieve their seed.
func ValidateInvestor(name string, pwhash string) (Investor, error) {
	var rec Investor
	user, err := ValidateUser(name, pwhash)
	if err != nil {
		return rec, errors.Wrap(err, "failed to validate user")
	}
	return RetrieveInvestor(user.Index)
}

// ValidateRecipient validates a particular recipient
func ValidateRecipient(name string, pwhash string) (Recipient, error) {
	var rec Recipient
	user, err := ValidateUser(name, pwhash)
	if err != nil {
		return rec, errors.Wrap(err, "Error while validating user")
	}
	return RetrieveRecipient(user.Index)
}
