package database

// contains the WIP Investor struct which will be st ored in a separate bucket
import (
	"encoding/json"
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
)

// NewInvestor creates a new investor object when passed the username, password hash,
// name and an option to generate the seed and publicKey. This is done because if
// we decide to allow anonymous investors to invest on our platform, we can easily
// insert their publickey into the system and hten have hanlders for them signing
// transactions
// TODO: add anonymous investor signing handlers
func NewInvestor(uname string, pwd string, Name string) (Investor, error) {
	// call this after the user has failled in username and password. Store hashed password
	// in the database
	var a Investor
	var err error
	a.U, err = NewUser(uname, pwd, Name)
	if err != nil {
		return a, err
	}
	a.AmountInvested = float64(0)
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
		return b.Put([]byte(utils.Uint32toB(a.U.Index)), encoded)
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
	// this route is broken becuase it reads through keys sequentially
	// need to see keys until the lenght of ther users dat abse
	var arr []Investor
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
		b := tx.Bucket(InvestorBucket)
		i := uint32(1)
		for ; i < limit; i++ {
			var rInvestor Investor
			x := b.Get(utils.Uint32toB(i))
			if x == nil {
				// this is where the key does not exist
				continue
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

// ValidateInvestor searches for an investor when passed the investor's name.
// This is useful for checking the user's password while logging in
func ValidateInvestor(uname string, pwd string) (Investor, error) {
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
			if rInvestor.U.LoginUserName == uname && rInvestor.U.LoginPassword == pwd {
				inv = rInvestor
			}
		}
		return fmt.Errorf("Not Found")
	})
	if inv.U.Index == 0 {
		return inv, fmt.Errorf("Investor Not Found")
	}
	return inv, err
}
