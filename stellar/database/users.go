package database

import (
	"encoding/json"
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
	"github.com/boltdb/bolt"
)

// the user structure houses all entities that are of type "User". This contains
// commonly used functions so that we need not repeat the ssame thing for every instance.
type User struct {
	Index int
	// default index, gets us easy stats on how many people are there and stuff,
	// don't want to omit this
	Seed string
	// Seed is the equivalent of a private key in stellar (stellar doesn't expose private keys)
	// do we make seed optional like that for the Recipient? Couple things to consider
	// here: if the recipient loses the publickey, it can nver send DEBTokens back
	// to the issuer, so it would be as if it reneged on the deal. Do we count on
	// technically less sound people to hold their public keys safely? I suggest
	// this would be  difficult in practice, so maybe enforce that they need to hold|
	// their account on the platform?
	Name string
	// Name of the primary stakeholder involved (principal trustee of school, for eg.)
	PublicKey string
	// PublicKey denotes the public key of the recipient
	LoginUserName string
	// the username you use to login to the platform
	LoginPassword string
	// the password, which you use to authenticate on the platform
	Address string
	// the registered address of the above company
	Description string
	// Does the contractor need to have a seed and a publickey?
	// we assume that it does in this case and proceed.
	// information on company credentials, their experience
	Image string
	// image can be company logo, founder selfie
	// hash of the password in reality
	FirstSignedUp string
	// auto generated timestamp
}

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
		a.Index = len(allUsers) + 1
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
	err = a.Save()
	return a, err // since user is a meta structure, insert it and then return the function
}

// InsertUser inserts a passed User object into the database
func (a *User) Save() error {
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
		return b.Put([]byte(utils.ItoB(a.Index)), encoded)
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
		for i := 1; ; i++ {
			var rUser User
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rUser)
			if err != nil {
				return err
			}
			arr = append(arr, rUser)
		}
		return nil
	})
	return arr, err
}

// RetrieveUser retrieves a particular User indexed by key from the database
func RetrieveUser(key int) (User, error) {
	var inv User
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return nil
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, err
}

func ValidateUser(name string, pwhash string) (User, error) {
	var inv User
	temp, err := RetrieveAllUsers()
	if err != nil {
		return inv, err
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; i < limit; i++ {
			var rUser User
			x := b.Get(utils.ItoB(i))
			if x == nil {
				continue
			}
			err := json.Unmarshal(x, &rUser)
			if err != nil {
				return err
			}
			// we have the User class, check names
			if rUser.LoginUserName == name && rUser.LoginPassword == pwhash {
				inv = rUser
				return nil
			}
		}
		return fmt.Errorf("Not Found")
	})
	return inv, err
}

func (a *User) GenKeys() error {
	var err error
	var dup User
	dup = *a
	dup.Seed, dup.PublicKey, err = xlm.GetKeyPair()
	if err != nil {
		return err
	}
	err = dup.Save()
	// a = &dup
	return err
}
