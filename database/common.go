package database

import (
	// "log"
	"encoding/json"
	"github.com/pkg/errors"

	edb "github.com/Varunram/essentials/database"
	consts "github.com/YaleOpenLab/openx/consts"
)

// this file contains commong methods that are repeated across interfaces
// Save inserts a passed User object into the database
func (a *User) Save() error {
	return edb.Save(consts.DbDir+consts.DbName, UserBucket, a, a.Index)
}

// RetrieveUser retrieves a particular User indexed by key from the database
func RetrieveUser(key int) (User, error) {
	var user User
	x, err := edb.Retrieve(consts.DbDir+consts.DbName, UserBucket, key)
	if err != nil {
		return user, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &user)
	return user, err
}

// RetrieveAllUsers gets a list of all User in the database
func RetrieveAllUsers() ([]User, error) {
	var arr []User
	x, err := edb.RetrieveAllKeys(consts.DbDir+consts.DbName, UserBucket)
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users")
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

func RetrieveAllUsersLim() (int, error) {
	return edb.RetrieveAllKeysLim(consts.DbDir+consts.DbName, UserBucket)
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
