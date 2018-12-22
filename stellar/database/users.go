package database

import (
	"fmt"
	"log"
	"encoding/json"

	"github.com/boltdb/bolt"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
)

// User is a metastrucutre that contains commonyl used keys within a single umbrella
// so that we can import it wherever needed.
func NewUser(uname string, pwd string, Name string) (User, error) {
	// call this after the user has failled in username and password.
	// Store hashed password in the database
	var a User

	allUsers, err := RetrieveAllUsers()
	if err != nil {
		return a, err
	}

	// the ugly indexing thing again, need to think of something better here
	if len(allUsers) == 0 {
		a.Index = 1
	} else {
		a.Index = uint32(len(allUsers) + 1)
	}

	a.Name = Name
	a.Seed, a.PublicKey, err = xlm.GetKeyPair()
	if err != nil {
		return a, err
	}
	a.LoginUserName = uname
	a.LoginPassword = utils.SHA3hash(pwd) // store tha sha3 hash
	// now we have a new User, take this and then send this struct off to be stored in the database
	a.FirstSignedUp = utils.Timestamp()
	err = InsertUser(a)
	return a, err // since user is a meta structure, insert it and then return the function
}

// InsertUser inserts a passed User object into the database
func InsertUser(a User) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.Uint32toB(a.Index)), encoded)
	})
	return err
}

// RetrieveAllUsers gets a list of all User in the database
func RetrieveAllUsers() ([]User, error) {
	var arr []User
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		i := uint32(1)
		for ; ; i++ {
			var rUser User
			x := b.Get(utils.Uint32toB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rUser)
			if err != nil {
				return nil
			}
			arr = append(arr, rUser)
		}
		return nil
	})
	return arr, err
}

// RetrieveUser retrieves a particular User indexed by key from the database
func RetrieveUser(key uint32) (User, error) {
	var inv User
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		x := b.Get(utils.Uint32toB(key))
		if x == nil {
			return nil
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, nil
}

// SearchForUser searches for an User when passed the User's name.
// This is useful for checking the user's password while logging in
func SearchForUser(name string, pwhash string) (User, error) {
	var inv User
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		i := uint32(1)
		for ; ; i++ {
			var rUser User
			x := b.Get(utils.Uint32toB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rUser)
			if err != nil {
				return nil
			}
			// we have the User class, check names
			if rUser.LoginUserName == name && rUser.LoginPassword == pwhash {
				inv = rUser
			}
		}
		return fmt.Errorf("Not Found")
	})
	return inv, err
}
