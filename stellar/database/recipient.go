package database

import (
	"encoding/json"
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
	"github.com/stellar/go/keypair"
)

func TestFn() {
	log.Println("Endpoint called! Cool!")
	return
}

func NewRecipient(uname string, pwhash string, Name string) (Recipient, error) {
	var a Recipient

	allRecipients, err := RetrieveAllRecipients()
	if err != nil {
		return a, err
	}

	// the ugly indexing thing again, need to think of something better here
	if len(allRecipients) == 0 {
		a.Index = 1
	} else {
		a.Index = uint32(len(allRecipients) + 1)
	}

	// generate a pk and seed pair and store it
	pair, err := keypair.Random()
	if err != nil {
		return a, err
	}
	a.Name = Name
	a.PublicKey = pair.Address()
	a.Seed = pair.Seed()
	a.FirstSignedUp = utils.Timestamp()
	a.LoginUserName = uname
	a.LoginPassword = pwhash
	// now we have a new Recipient, take this and then send this off to be stored in the database
	log.Println("Created Recipient: ", a)
	return a, nil
}

func NewRecipientWithoutSeed(uname string, pwhash string, Name string) (Recipient, error) {
	// this should be called initially since the recipient's seed is created only
	// if someone decides to ivnest in the order
	var a Recipient

	allRecipients, err := RetrieveAllRecipients()
	if err != nil {
		return a, err
	}

	// the ugly indexing thing again, need to think of something better here
	if len(allRecipients) == 0 {
		a.Index = 1
	} else {
		a.Index = uint32(len(allRecipients) + 1)
	}

	a.Name = Name
	a.FirstSignedUp = utils.Timestamp()
	a.LoginUserName = uname
	a.LoginPassword = pwhash
	// now we have a new Recipient, take this and then send this off to be stored in the database
	return a, nil
}

// all operations are mostly similar to that of the Recipient class
// TODO: merge where possible by adding an extra bucket param
func InsertRecipient(a Recipient) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(RecipientBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.Uint32toB(a.Index)), encoded)
		// but why do we index based on Index?
		// this is because we do want to enumerate through all Recipients, which can not be done
		// in a name based construction. But this makes search ahrder, since now you
		// all entries to find something as simple as a password. But if this is the
		// only use case that exists, we index by password hash and then get data only
		// when the user requests it. Nice data protection as well
		// TODO: discuss indexing by pwd hash and implications. For small no of entries,
		// we can s till tierate over all the entries.
	})
	return err
}

func RetrieveAllRecipients() ([]Recipient, error) {
	var arr []Recipient
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(RecipientBucket)
		i := uint32(1)
		for ; ; i++ {
			var rRecipient Recipient
			x := b.Get(utils.Uint32toB(i))
			if x == nil {
				// this is where the key does not exist
				return nil
			}
			err := json.Unmarshal(x, &rRecipient)
			//if err != nil && rRecipient.Live == false {
			if err != nil {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			arr = append(arr, rRecipient)
		}
		return nil
	})
	return arr, err
}

func RetrieveRecipient(key uint32) (Recipient, error) {
	var inv Recipient
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(RecipientBucket)
		x := b.Get(utils.Uint32toB(key))
		if x == nil {
			return nil
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, nil
}

func SearchForRecipient(name string) (Recipient, error) {
	var inv Recipient
	// this is very ugly, but the only way it works right now (see TODO earlier)
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(RecipientBucket)
		i := uint32(1)
		for ; ; i++ {
			var rRecipient Recipient
			x := b.Get(utils.Uint32toB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rRecipient)
			if err != nil {
				return nil
			}
			// we have the investor class, check password
			if rRecipient.Name == name {
				inv = rRecipient
			}
		}
		return fmt.Errorf("Not Found")
	})
	return inv, err
}
