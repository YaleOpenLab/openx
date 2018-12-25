package database
// recipient.go defines all recipient related functions that are not defined on
// the struct itself.
import (
	"encoding/json"
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
)

func TestFn() {
	log.Println("Endpoint called! Cool!")
	return
}

func NewRecipient(uname string, pwd string, Name string) (Recipient, error) {
	var a Recipient
	var err error
	a.U, err = NewUser(uname, pwd, Name)
	if err != nil {
		return a, err
	}
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
		return b.Put([]byte(utils.Uint32toB(a.U.Index)), encoded)
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

// RetrieveAllRecipients gets a list of all Recipient in the database
func RetrieveAllRecipients() ([]Recipient, error) {
	var arr []Recipient
	temp, err := RetrieveAllUsers()
	if err != nil {
		return arr, err
	}
	limit := uint32(len(temp) + 1)
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
		for ; i < limit; i++ {
			var rRecipient Recipient
			x := b.Get(utils.Uint32toB(i))
			if x == nil {
				// this is where the key does not exist
				continue
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

func ValidateRecipient(name string, pwhash string) (Recipient, error) {
	var inv Recipient
	temp, err := RetrieveAllUsers()
	if err != nil {
		return inv, err
	}
	limit := uint32(len(temp) + 1)
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(RecipientBucket)
		i := uint32(1)
		for ; i < limit ; i++ {
			var rRecipient Recipient
			x := b.Get(utils.Uint32toB(i))
			if x == nil {
				// the given key may be an investor
				continue
			}
			err := json.Unmarshal(x, &rRecipient)
			if err != nil {
				return nil
			}
			// we have the Recipient class, check password
			if rRecipient.U.LoginUserName == name && pwhash == rRecipient.U.LoginPassword {
				inv = rRecipient
				return nil
			}
		}
		return fmt.Errorf("Not Found")
	})
	return inv, err
}


// DeleteKeyFromBucket deletes a given key from the bucket _bucketName
func DeleteKeyFromBucket(key uint32, bucketName []byte) error {
	// deleting order might be dangerous since that would mess with the RetrieveAllOrders
	// function, have it in here for now, don't do too much with it / fiox retrieve all
	// to handle this case
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(RecipientBucket)
		err := b.Delete(utils.Uint32toB(key))
		if err != nil {
			return err
		}
		log.Println("Deleted recipient with key: ", key)
		return nil
	})
	return err
}
