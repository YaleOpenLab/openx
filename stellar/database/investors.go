package database

// contains the WIP Investor struct which will be st ored in a separate bucket
import (
	"encoding/json"
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
	"github.com/boltdb/bolt"
)

// NewInvestor creates a new investor object when passed the username, password hash,
// name and an option to generate the seed and publicKey. This is done because if
// we decide to allow anonymous investors to invest on our platform, we can easily
// insert their pbulickey into the system and hten have hanlders for them signing
// transactions
// TODO: add anonymous investor signing handlers
func NewInvestor(uname string, pwhash string, Name string, pkgen bool) (Investor, error) {
	// call this after the user has failled in username and password. Store hashed password
	// in the database
	var a Investor

	allInvestors, err := RetrieveAllInvestors()
	if err != nil {
		return a, err
	}

	// the ugly indexing thing again, need to think of something better here
	if len(allInvestors) == 0 {
		a.Index = 1
	} else {
		a.Index = uint32(len(allInvestors) + 1)
	}

	// for investors, we need to index by username, so Index is not that useful
	// except maybe for quick stats
	a.Name = Name
	if pkgen {
		a.Seed, a.PublicKey, err = xlm.GetKeyPair()
		if err != nil {
			return a, err
		}
	}
	a.AmountInvested = float64(0)
	a.FirstSignedUp = utils.Timestamp()
	// don't set InvestedAssets
	a.LoginUserName = uname
	a.LoginPassword = pwhash
	// now we have a new investor, take this and then send this off to be stored in the database
	return a, nil
}

// InsertInvestor inserts a passed Investor object into the database
func InsertInvestor(a Investor) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(InvestorBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.Uint32toB(a.Index)), encoded)
		// but why do we index based on Index?
		// this is because we do want to enumerate through all investors, which can not be done
		// in a name based construction. But this makes search ahrder, since now you
		// all entries to find something as simple as a password. But if this is the
		// only use case that exists, we index by password hash and then get data only
		// when the user requests it. Nice data protection as well
		// TODO: discuss indexing by pwd hash and implications. For small no of entries,
		// we can s till tierate over all the entries.
	})
	return err
}

// RetrieveAllInvestors gets a list of all investor in the database
func RetrieveAllInvestors() ([]Investor, error) {
	var arr []Investor
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(InvestorBucket)
		i := uint32(1)
		for ; ; i++ {
			var rInvestor Investor
			x := b.Get(utils.Uint32toB(i))
			if x == nil {
				// this is where the key does not exist
				return nil
			}
			err := json.Unmarshal(x, &rInvestor)
			//if err != nil && rInvestor.Live == false {
			if err != nil {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			arr = append(arr, rInvestor)
		}
		return nil
	})
	return arr, err
}

// RetrieveInvestor retrieves a particular investor indexed by key from the database
func RetrieveInvestor(key uint32) (Investor, error) {
	var inv Investor
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(InvestorBucket)
		x := b.Get(utils.Uint32toB(key))
		if x == nil {
			return nil
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, nil
}

// SearchForInvestor searches for an investor when passed the investor's name.
// This is useful for checking the user's password while logging in
func SearchForInvestor(name string) (Investor, error) {
	var inv Investor
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(InvestorBucket)
		i := uint32(1)
		for ; ; i++ {
			var rInvestor Investor
			x := b.Get(utils.Uint32toB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rInvestor)
			if err != nil {
				return nil
			}
			// we have the investor class, check names
			if rInvestor.LoginUserName == name {
				inv = rInvestor
			}
		}
		return fmt.Errorf("Not Found")
	})
	return inv, err
}

// TODO: migrate the password checking logic here and we can simply have
// something like ValidateUser()
// Also, have a new user class implemented which can be borrowed by all
// subsequent classes
