package opensolar

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	database "github.com/YaleOpenLab/openx/database"
	notif "github.com/YaleOpenLab/openx/notif"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
	xlm "github.com/YaleOpenLab/openx/xlm"
	"github.com/boltdb/bolt"
)

// the contractor super struct comprises of various entities within it. Its a
// super struct because combining them results in less duplication of code

type Entity struct {
	U database.User
	// inherit the base user class
	Contractor bool
	// the name of the contractor / company that is contracting
	// A contractor is party who proposes a specific some of money towards a
	// particular project. This is the actual amount that the investors invest in.
	// This ideally must include the developer fee within it, so that investors
	// don't have to invest in two things. It would also make sense because the contractors
	// sometimes would hire developers themselves.
	Developer bool
	// A developer is someone who installs the required equipment (Raspberry Pi,
	// network adapters, anti tamper installations and similar) In the initial
	// projects, this will be us, since we'd be installign the pi ourselves, but in
	// the future, we expect third party developers / companies to do this for us
	// and act in a decentralized fashion. This money can either be paid out of chain
	// in fiat or can be a portion of the funds the investors chooses to invest in.
	// a contractor may also employ developers by himself, so this entity is not
	// strictly necessary.
	Originator bool
	// An Originator is an entity that will start a project and get a fixed fee for
	// rendering its service. An Originator's role is not restricted, the originator
	// can also be the developer, contractor or guarantor. The originator should take
	// the responsibility of auditing the requirements of the project - panel size,
	// location, number of panels needed, etc. He then should ideally be able to fill
	// out some kind of form on the website so that the originator's proposal is live
	// and shown to potential investors. The originators get paid only when the project
	// is live, else they can just spam, without any actual investment
	Guarantor bool
	// A guarantor is somebody who can assure investors that the school will get paid
	// on time. This authority should be trusted and either should be vetted by the law
	// or have a multisig paying out to the investors beyond a certain timeline if they
	// don't get paid by the school. This way, the guarantor can be anonymous, like the
	// nice Pineapple Fund guy. This can also be an insurance company, who is willing to
	// guarantee for specific school and the school can pay him out of chain / have
	// that as fee within the contract the originator
	PastContracts []Project
	// list of all the contracts that the contractor has won in the past
	ProposedContracts []Project
	// the Originator proposes a contract which will then be taken up
	// by a contractor, who publishes his own copy of the proposed contract
	// which will be the set of contracts that will be sent to auction
	PresentContracts []Project
	// list of all contracts that the contractor is presently undertaking1
	PastFeedback []Feedback
	// feedback received on the contractor from parties involved in the past
	// What kind of proof do we want from the company? KYC?
	// maybe we could have a photo op like exchanges do these days, with the owner
	// holding up his drivers' license or similar
	Collateral float64
	// the amount of collateral that the entity is willing to hold in case it reneges
	// on a specific contract's details. This is an optional parameter but having collateral
	// would increase investor confidence that a particular entity will keep its word
	// regarding a particular contract.
	CollateralData []string
	// the specific thing(s) which the contractor wants to hold as collateral described
	// as a string (for eg, if a cash bond worht 5000 USD is held as collaterlal,
	// collateral would be set to 5000 USD and CollateralData would be "Cash Bond")
}

// Save stores the entity in the database
func (a *Entity) Save() error {
	db, err := database.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ContractorBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			return err
		}
		return b.Put([]byte(utils.ItoB(a.U.Index)), encoded)
	})
	return err
}

// RetrieveAllEntitiesWithoutRole gets all the proposed contracts for a particular recipient
func RetrieveAllEntitiesWithoutRole() ([]Entity, error) {
	var arr []Entity
	temp, err := database.RetrieveAllUsers()
	if err != nil {
		return arr, err
	}
	limit := len(temp) + 1
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ContractorBucket)
		for i := 1; i < limit; i++ {
			var rContractor Entity
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// might be some other user like an investor or recipient
				continue
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

// RetrieveAllEntities gets all the proposed contracts for a particular recipient
func RetrieveAllEntities(role string) ([]Entity, error) {
	var arr []Entity
	temp, err := database.RetrieveAllUsers()
	if err != nil {
		return arr, err
	}
	limit := len(temp) + 1
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ContractorBucket)
		for i := 1; i < limit; i++ {
			var rContractor Entity
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
			default:
				continue
				// default is to add all contractentities to the array
			}
			arr = append(arr, rContractor)
		}
		return nil
	})
	return arr, err
}

// RetrieveEntity retrieves a specific entity from the database
func RetrieveEntity(key int) (Entity, error) {
	var a Entity
	db, err := database.OpenDB()
	if err != nil {
		return a, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ContractorBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return fmt.Errorf("Retrieving entity returns nil, quitting!")
		}
		return json.Unmarshal(x, &a)
	})
	return a, err
}

// newEntity creates a new entity based on the role passed
func newEntity(uname string, pwd string, seedpwd string, Name string, Address string, Description string, role string) (Entity, error) {
	var a Entity
	var err error
	a.U, err = database.NewUser(uname, pwd, seedpwd, Name)
	if err != nil {
		return a, err
	}

	a.U.Address = Address
	a.U.Description = Description

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
		return a, fmt.Errorf("invalid entity type passed!")
	}

	err = a.Save()
	return a, err
}

// ChangeReputation changes the reputation associated with a particular entity
func ChangeReputation(entityIndex int, reputation float64) error {
	a, err := RetrieveEntity(entityIndex)
	if err != nil {
		return err
	}
	if reputation > 0 {
		err = a.U.IncreaseReputation(reputation)
	} else {
		err = a.U.DecreaseReputation(reputation)
	}
	if err != nil {
		return err
	}
	return a.Save()
}

// TopReputationEntitiesWithoutRole returns the list of all the top reputed entities in descending order
func TopReputationEntitiesWithoutRole() ([]Entity, error) {
	allEntities, err := RetrieveAllEntitiesWithoutRole()
	if err != nil {
		return allEntities, err
	}
	for i, _ := range allEntities {
		for j, _ := range allEntities {
			if allEntities[i].U.Reputation < allEntities[j].U.Reputation {
				tmp := allEntities[i]
				allEntities[i] = allEntities[j]
				allEntities[j] = tmp
			}
		}
	}
	return allEntities, nil
}

// TopReputationEntities returns the list of all the top reputed entities with the specific role in descending order
func TopReputationEntities(role string) ([]Entity, error) {
	// caller knows what role he needs this list for, so directly retrieve and do stuff here
	allEntities, err := RetrieveAllEntities(role)
	if err != nil {
		return allEntities, err
	}
	for i, _ := range allEntities {
		for j, _ := range allEntities {
			if allEntities[i].U.Reputation < allEntities[j].U.Reputation {
				tmp := allEntities[i]
				allEntities[i] = allEntities[j]
				allEntities[j] = tmp
			}
		}
	}
	return allEntities, nil
}

// ValidateEntity validates the entity with the specific name and pwhash and returns true if everything matches the thing on record
func ValidateEntity(name string, pwhash string) (Entity, error) {
	var rec Entity
	user, err := database.ValidateUser(name, pwhash)
	if err != nil {
		return rec, err
	}
	return RetrieveEntity(user.Index)
}

func AgreeToContractConditions(contractHash string, projIndex string,
	debtAssetCode string, entityIndex int, seedpwd string) error {
	// we need to display this on the frontend and once the user presses agree, commit
	// a tx to the blockchain with the outcome
	message := "I agree to the terms and conditions specified in contract " + contractHash +
		"and by signing this message to the blockchain agree that I accept the investment in project " + projIndex +
		"whose debt asset is: " + debtAssetCode
	// hash the message and transmit the message in 5 parts
	// eg.
	// CONTRACTHASH9a768ace36ff3d17
	// 71d5c145a544de3d68343b2e7609
	// 3cb7b2a8ea89ac7f1a20c852e6fc
	// 1d71275b43abffefac381c5b906f
	// 55c3bcff4225353d02f1d3498758

	user, err := database.RetrieveUser(entityIndex)
	if err != nil {
		log.Println(err)
		return err
	}

	seed, err := wallet.DecryptSeed(user.EncryptedSeed, seedpwd)
	if err != nil {
		log.Println(err)
		return err
	}

	messageHash := "CONTRACTHASH" + strings.ToUpper(utils.SHA3hash(message))
	firstPart := messageHash[:28] // higher limit is not included in the slice
	secondPart := messageHash[28:56]
	thirdPart := messageHash[56:84]
	fourthPart := messageHash[84:112]
	fifthPart := messageHash[112:140]

	timeStamp := utils.I64toS(utils.Unix())
	_, firstHash, err := xlm.SendXLM(user.PublicKey, timeStamp, seed, firstPart)
	if err != nil {
		log.Println(err)
		return err
	}

	_, secondHash, err := xlm.SendXLM(user.PublicKey, timeStamp, seed, secondPart)
	if err != nil {
		log.Println(err)
		return err
	}

	_, thirdHash, err := xlm.SendXLM(user.PublicKey, timeStamp, seed, thirdPart)
	if err != nil {
		log.Println(err)
		return err
	}

	_, fourthHash, err := xlm.SendXLM(user.PublicKey, timeStamp, seed, fourthPart)
	if err != nil {
		log.Println(err)
		return err
	}

	_, fifthHash, err := xlm.SendXLM(user.PublicKey, timeStamp, seed, fifthPart)
	if err != nil {
		log.Println(err)
		return err
	}

	//if user.Notification {
	notif.SendContractNotification(firstHash, secondHash, thirdHash, fourthHash, fifthHash, user.Email)
	//}

	return nil
}