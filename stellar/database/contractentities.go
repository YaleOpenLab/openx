package database

import (
	"encoding/json"
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
)

func newContractEntityHelper(uname string, pwd string, Name string, Address string, Description string, role string) (ContractEntity, error) {
	// call this after the user has failled in username and password. Store hashed password
	// in the database
	var a ContractEntity
	var err error
	a.U, err = NewUser(uname, pwd, Name)
	if err != nil {
		return a, err
	}
	// set all auto fields above
	a.U.Address = Address
	a.U.Description = Description
	// insertion into the database will be a separate handler, pass this ContractEntity there
	switch role {
	case "contractor":
		a.Contractor = true
	case "developer":
		a.Developer = true
	case "originator":
		a.Originator = true
	case "guarantor":
		a.Guarantor = true
	default:
		// nothing, since only we call this function internally, this shouldn't arrive here
	}
	return a, nil
}

func InsertContractEntity(a ContractEntity) error {
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
		return b.Put([]byte(utils.ItoB(a.U.Index)), encoded)
	})
	return err
}

func NewContractEntity(uname string, pwd string, Name string, Address string, Description string, role string) (ContractEntity, error) {
	var dummy ContractEntity
	switch role {
	case "originator":
		return newContractEntityHelper(uname, pwd, Name, Address, Description, "originator")
	case "developer":
		return newContractEntityHelper(uname, pwd, Name, Address, Description, "developer")
	case "contractor":
		return newContractEntityHelper(uname, pwd, Name, Address, Description, "contractor")
	case "guarantor":
		return newContractEntityHelper(uname, pwd, Name, Address, Description, "guarantor")
	}
	return dummy, fmt.Errorf("Invalid entity passed, check again!")
}

// gets all the proposed contracts for a particular recipient
func RetrieveAllContractEntities(role string) ([]ContractEntity, error) {
	var arr []ContractEntity
	temp, err := RetrieveAllUsers()
	if err != nil {
		return arr, err
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ContractorBucket)
		for i := 1 ; i < limit; i++ {
			var rContractor ContractEntity
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// might be some other user like an investor or recipient
				continue
			}
			err := json.Unmarshal(x, &rContractor)
			if err != nil {
				return nil
			}
			switch role {
			case "contractor":
				if !rContractor.Contractor {
					continue
				}
			case "developer":
				if !rContractor.Developer {
					continue
				}
			case "originator":
				if !rContractor.Originator {
					continue
				}
			case "guarantor":
				if !rContractor.Guarantor {
					continue
				}
				// default is to add all contractentities to the array
			}
			arr = append(arr, rContractor)
		}
		return nil
	})
	return arr, err
}

func RetrieveContractEntity(key int) (ContractEntity, error) {
	var a ContractEntity
	db, err := OpenDB()
	if err != nil {
		return a, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ContractorBucket)
		x := b.Get(utils.ItoB(key))
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
func SearchForContractEntity(name string, pwhash string) (ContractEntity, error) {
	var a ContractEntity
	db, err := OpenDB()
	if err != nil {
		return a, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		// TODO: change all similar functions to db.View
		b := tx.Bucket(ContractorBucket)
		for i := 1; ; i++ {
			var rContractor ContractEntity
			x := b.Get(utils.ItoB(i))
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

// you need to have a lock in period beyond which contractors can not post what
// stuff they want. now, how do you choose which contractor wins? Ideally,
// the school would want the most stuff but you need to vet which contracts are good
// and not. In this case, we use prive as the metric, but this can be anything
// or even chosen by the school / demo bidding auction by investors and then
// take the one which has the most demo votes
// Also, in contracts, when contractors are proposing a contract towards something,
// we need to be sure that they are not followign the price and are instead giving
// their best quote possible. In this case, a blind auction method is the best
// and that's what we have right now. If we want this to be an auction as well, we
// need to have a specific date of sorts where all the contractors can propose
// contracts immmediately, without latency.
// Also, have some kind of deposit for Contractors (5% or something) so that they
// don't go back on their investment and slash their ivnestment by 10% if this happens
// and distribute that amount to the recipient directly and reduce everyone's bids
// by that amount to account for the change in underlying Order
// also, a given Contractor right now is allowed only for one final bid for blind
// auction advantages (no price disvocery, etc). If we want to change this, we must
// have an auction handler that will take care of this.
