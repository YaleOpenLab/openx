package database

import (
	"encoding/json"
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
	"github.com/stellar/go/keypair"
)

type Recipient struct {
	Index uint32
	// defauult index, gets us easy stats on how many people are there and stuff,
	// don't want to omit this
	Name string
	// Name of the primary stakeholder involved (principal trustee of school, for eg.)
	PublicKey string
	// PublicKey denotes the public key of the recipient
	Seed string
	// do we make seed optional like that for the Recipient? Couple things to consider
	// here: if the recipient loses the publickey, it can nver send DEBTokens back
	// to the issuer, so it would be as if it reneged on the deal. Do we count on
	// technically less sound people to hold their public keys safely? I suggest
	// this would be  difficult in practice, so maybe enforce that they need to hold|
	// their accounts on the platform?
	FirstSignedUp string
	// auto generated timestamp
	DebtAssets []string
	// DebtAssets denotes the list of all DEBTokens that the recipient possesses
	// this is an array since a single recipient could technically still have multiple
	// projects under its wing which Recipients can invest in.
	PaybackAssets []string
	// Payback Assets denotes the status of all assets that the recipient has received
	// this could be used to easily display payback progress, calculate ratings
	// for a specific school and so on.
	LoginUserName string
	// the thing you use to login to the platform
	LoginPassword string
	// password, which is separate from the generated seed.
}

var RecipientBucket = []byte("Recipients")

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

// all operations are mostly similar to that of the Recipient class
// TODO: merge where possible by adding an extra bucket param
func InsertRecipient(a Recipient) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(RecipientBucket) // the orders bucket contains all our orders
		if err != nil {
			log.Fatal(err)
			return err
		}
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
		b, err := tx.CreateBucketIfNotExists(RecipientBucket)
		if err != nil {
			return err
		}
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
		b, err := tx.CreateBucketIfNotExists(RecipientBucket)
		if err != nil {
			return err
		}
		x := b.Get(utils.Uint32toB(key))
		if x == nil {
			return nil
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, nil
}

func SearchForRecipientPassword(pwhash string) (Recipient, error) {
	var inv Recipient
	// this is very ugly, but the only way it works right now (see TODO earlier)
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(RecipientBucket)
		if err != nil {
			return err
		}
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
			if rRecipient.LoginPassword == pwhash {
				log.Println("FOUDN INVESOTR")
				inv = rRecipient
			}
		}
		return fmt.Errorf("Not Found")
	})
	return inv, err
}
