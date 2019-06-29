package database

import (
	"github.com/pkg/errors"

	utils "github.com/Varunram/essentials/utils"
	"github.com/boltdb/bolt"
)

// Recipient defines the recipient structure
type Recipient struct {
	U *User
	// user related functions are called as an instance directly
	ReceivedSolarProjects       []string
	ReceivedSolarProjectIndices []int
	ReceivedConstructionBonds   []string
	// ReceivedProjects denotes the projects that have been received by the recipient
	// instead of storing the PaybackAssets and the DebtAssets, we store this
	DeviceId string
	// the device ID of the associated solar hub. We don't do much with it here,
	// but we need it on the IoT Hub side to check login stuff
	DeviceStarts []string
	// the start time of the devices recorded for reference. We could monitor unscheduled
	// closes on the platform level as well and send email notifications or similar
	DeviceLocation string
	// the location of the device. Teller gets location using google's geolocation
	// API. Accuracy is of the order ~1km radius. Not great, but enough to detect
	// theft or something
	StateHashes []string
	// StateHashes provides the list of state updates (ipfs hashes) that the teller associated with this
	// particular recipient has communicated.
	TotalEnergyCP float64
	// the total energy produced by the recipient's assets in the current period
	TotalEnergy float64
	// the total energy produced by the recipient's assets over all billed periods
	Autoreload bool
	// a bool to denote whether the recipient wants to reload balance from his secondary account to pay any dues that are remaining
}

// NewRecipient returns a new recipient provided with the function parameters
func NewRecipient(uname string, pwd string, seedpwd string, Name string) (Recipient, error) {
	var a Recipient
	var err error
	user, err := NewUser(uname, pwd, seedpwd, Name)
	if err != nil {
		return a, errors.Wrap(err, "failed to retrieve new user")
	}
	a.U = &user
	err = a.Save()
	return a, err
}

// Save saves a given recipient's details
func (a *Recipient) Save() error {
	db, err := OpenDB()
	if err != nil {
		return errors.Wrap(err, "Error while opening database")
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(RecipientBucket)
		encoded, err := a.MarshalJSON()
		if err != nil {
			return errors.Wrap(err, "Error while marshaling json")
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
		return arr, errors.Wrap(err, "Error while retreiving all users from database")
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "Error while opening database")
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
			err := rRecipient.UnmarshalJSON(x)
			if err != nil {
				return errors.Wrap(err, "Error while unmarshalling json")
			}
			arr = append(arr, rRecipient)
		}
		return nil
	})
	return arr, err
}

// RetrieveRecipientHelper is a helper associated with the RetrieveRecipient function
func RetrieveRecipientHelper(key int) (Recipient, error) {
	var rec Recipient
	db, err := OpenDB()
	if err != nil {
		return rec, errors.Wrap(err, "Error while opening database")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(RecipientBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			// there is no key with the specific details
			return errors.New("Recipient not found!")
		}
		return rec.UnmarshalJSON(x)
	})
	return rec, err
}

// RetrieveRecipient retrieves a specific recipient from the database
func RetrieveRecipient(key int) (Recipient, error) {
	var rec Recipient
	user, err := RetrieveUser(key)
	if err != nil {
		return rec, err
	}
	rec, err = RetrieveRecipientHelper(key)
	if err != nil {
		return rec, err
	}
	rec.U = &user
	return rec, rec.Save()
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

// TopReputationRecipient returns a list of recipients with the best reputation
func TopReputationRecipient() ([]Recipient, error) {
	allRecipients, err := RetrieveAllRecipients()
	if err != nil {
		return allRecipients, errors.Wrap(err, "failed to retrieve all recipients")
	}
	for i := range allRecipients {
		for j := range allRecipients {
			if allRecipients[i].U.Reputation > allRecipients[j].U.Reputation {
				tmp := allRecipients[i]
				allRecipients[i] = allRecipients[j]
				allRecipients[j] = tmp
			}
		}
	}
	return allRecipients, nil
}
