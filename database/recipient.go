package database

// recipient.go defines all recipient related functions that are not defined on
// the struct itself.
import (
	"encoding/json"
	"fmt"

	utils "github.com/OpenFinancing/openfinancing/utils"
	"github.com/boltdb/bolt"
)

type Recipient struct {
	ReceivedSolarProjects []string
	// ReceivedProjects denotes the projects that have been received by the recipient
	// instead of storing the PaybackAssets and the DebtAssets, we store this
	U User
	// user related functions are called as an instance directly
	// TODO: Consider how effective the name 'recipient' is. Consider more information
	// about recipients to add in the struct, For example, recipients should be associated
	// to sites eligible for projects (eg. a building or land where you can put panels),
	// (and eventually need to show proof of this)
}

// NewRecipient returns a new recipient provided with the function parameters
func NewRecipient(uname string, pwd string, seedpwd string, Name string) (Recipient, error) {
	var a Recipient
	var err error
	a.U, err = NewUser(uname, pwd, seedpwd, Name)
	if err != nil {
		return a, err
	}
	err = a.Save()
	return a, err
}

func (a *Recipient) AddEmail(email string) error {
	// call this function when a user wants to get notifications. Ask on frontend whether
	// it wants to
	a.U.Email = email
	a.U.Notification = true
	return a.Save()
}

// Save() saves a given recipient's details
func (a *Recipient) Save() error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(RecipientBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			return err
		}
		return b.Put([]byte(utils.ItoB(a.U.Index)), encoded)
	})
	return err
}

// RetrieveAllRecipients gets a list of all Recipient in the database
func RetrieveAllRecipients() ([]Recipient, error) {
	var arr []Recipient
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

	err = db.View(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(RecipientBucket)
		i := 1
		for ; i < limit; i++ {
			var rRecipient Recipient
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// this is where the key does not exist
				continue
			}
			err := json.Unmarshal(x, &rRecipient)
			if err != nil {
				return err
			}
			arr = append(arr, rRecipient)
		}
		return nil
	})
	return arr, err
}

// RetrieveRecipient retrieves a specific recipient from the database
func RetrieveRecipient(key int) (Recipient, error) {
	var inv Recipient
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(RecipientBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			// there is no key with the specific details
			return fmt.Errorf("Recipient not found!")
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, err
}

func ValidateRecipient(name string, pwhash string) (Recipient, error) {
	var rec Recipient
	user, err := ValidateUser(name, pwhash)
	if err != nil {
		return rec, err
	}
	return RetrieveRecipient(user.Index)
}

func ChangeRecpReputation(recpIndex int, reputation float64) error {
	a, err := RetrieveRecipient(recpIndex)
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

func TopReputationRecipient() ([]Recipient, error) {
	allRecipients, err := RetrieveAllRecipients()
	if err != nil {
		return allRecipients, err
	}
	for i, _ := range allRecipients {
		for j, _ := range allRecipients {
			if allRecipients[i].U.Reputation < allRecipients[j].U.Reputation {
				tmp := allRecipients[i]
				allRecipients[i] = allRecipients[j]
				allRecipients[j] = tmp
			}
		}
	}
	return allRecipients, nil
}
