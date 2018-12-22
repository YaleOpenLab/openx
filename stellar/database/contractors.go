package database

import (
	"encoding/json"
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
)

// TODO: most of these entities have some common fields like Index, Name,
// LoginUserName, LoginPassword, FirstSignedUp, Seed, PublicKey
// we should split these into a separate entity called a "User" and have all
// other entities import from this low level entity. That would save lots of
// code duplication on our way forward.
/*
	 Contractor Fields
		 Index uint32 auto
		 Name string required
		 Address string required
		 Description string required
		 Image string optional
		Seed string auto
		PublicKey string auto
		 LoginUserName string required
		 LoginPassword string required
		one of the following four flags is required:
			IsContractor bool
			IsGuarantor bool
			IsDeveloper bool
			IsOriginator bool
		 PastContracts []Contract
		 PresentContracts []Contract
		 PastFeedback []Feedback
		FirstSignedUp string auto
*/
func NewContractor(uname string, pwd string, Name string, Address string, Description string) (Contractor, error) {
	// call this after the user has failled in username and password. Store hashed password
	// in the database
	var a Contractor
	var err error
	a.U, err = NewUser(uname, pwd, Name)
	if err != nil {
		return a, err
	}
	// set all auto fields above
	a.Address = Address
	a.Description = Description
	// insertion into the database will be a separate handler, pass this contractor there
	return a, nil
}

func InsertContractor(a Contractor) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ContractorBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.Uint32toB(a.U.Index)), encoded)
	})
	return err
}

func RetrieveAllContractors() ([]Contractor, error) {
	var arr []Contractor
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ContractorBucket)
		i := uint32(1)
		for ; ; i++ {
			var rContractor Contractor
			x := b.Get(utils.Uint32toB(i))
			if x == nil {
				// no key, return
				return nil
			}
			err := json.Unmarshal(x, &rContractor)
			if err != nil {
				return nil
			}
			arr = append(arr, rContractor)
		}
		return nil
	})
	return arr, err
}

func RetrieveContractor(key uint32) (Contractor, error) {
	var a Contractor
	db, err := OpenDB()
	if err != nil {
		return a, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ContractorBucket)
		x := b.Get(utils.Uint32toB(key))
		if x == nil {
			return nil
		}
		return json.Unmarshal(x, &a)
	})
	return a, nil
}

// search by username for login stuff
// TODO: if two people have the same username, bolt defaults to the alst inserted
// one. So we need to have a function that prevents username collisions
func SearchForContractor(name string, pwhash string) (Contractor, error) {
	var a Contractor
	db, err := OpenDB()
	if err != nil {
		return a, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ContractorBucket)
		i := uint32(1)
		for ; ; i++ {
			var rContractor Contractor
			x := b.Get(utils.Uint32toB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rContractor)
			if err != nil {
				return nil
			}
			// we have the investor class, check names
			if rContractor.U.LoginUserName == name && rContractor.U.LoginPassword == pwhash {
				a = rContractor
			}
		}
		return fmt.Errorf("Not Found")
	})
	return a, err
}
