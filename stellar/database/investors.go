package database

// contains the WIP Investor struct which will be st ored in a separate bucket
import (
	"encoding/json"
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
	"github.com/stellar/go/keypair"
)

// the investor s truct contains all the investor details such as
// public key, seed (if account is created on the website) and ot her stuff which
// is yet to be decided

// ALl investors will be referenced by their public key, name is optional (maybe necessary?)
// we need to stil ldecide on identity and stuff and how much we want to track
// people who invest in the schools
type Investor struct {
	Index uint32
	// defauult index, gets us easy stats on how many people are there and stuff,
	// don't want to omit this
	Name string
	// display Name, different from UserName
	PublicKey string
	// the PublicKey used to identify you on the platform. We could still reference
	// people by name, but we needn't since we have the pk anyway.
	Seed string
	// optional, this is if the user created his account on our website
	// should be shown once and deleted permanently
	// add a notice like "WE DO NOT SAVE YOUR SEED" on the UI side
	AmountInvested float64
	// total amount, would be nice to track to contact them,
	// give them some kind of medals or something
	FirstSignedUp string
	// auto generated timestamp
	InvestedAssets []Order
	// array of asset codes this user has invested in
	// also I think we need a username + password for logging on to the platform itself
	// linking it here for now
	LoginUserName string
	// the thing you use to login to the platform
	LoginPassword string
	// LoginPassword is different from the seed you get if you choose
	// to open your account on the website. This is becasue even if you lose the
	// login password, you needn't worry too much about losing your funds, sicne you have
	// your seed and can send them to another address immediately.
}

var InvestorBucket = []byte("Investors")

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
		// generate a pk and seed pair and store it
		pair, err := keypair.Random()
		if err != nil {
			return a, err
		}
		a.Seed = pair.Seed()
		a.PublicKey = pair.Address()
		// display this seed but DON'T store this. Store this for now sicne we're just testing
		log.Println("This seed will be deleted from our servers. Note it down and please don't forget", a.Seed)
	}
	a.AmountInvested = float64(0)
	a.FirstSignedUp = utils.Timestamp()
	// don't set InvestedAssets
	a.LoginUserName = uname
	a.LoginPassword = pwhash
	// now we have a new investor, take this and then send this off to be stored in the database
	return a, nil
}

func CheckPassword(pwHash string) (bool, error) {
	// frontend should serve sha3(password) to us, we compare that and what's stored
	// in the database to s ee if they match
	return true, nil
}

func InsertInvestor(a Investor) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(InvestorBucket) // the orders bucket contains all our orders
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
		b, err := tx.CreateBucketIfNotExists(InvestorBucket)
		if err != nil {
			return err
		}
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

func RetrieveInvestor(key uint32) (Investor, error) {
	var inv Investor
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(InvestorBucket)
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

func SearchForInvestorPassword(pwhash string) (Investor, error) {
	var inv Investor
	// this is very ugly, but the only way it works right now (see TODO earlier)
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(InvestorBucket)
		if err != nil {
			return err
		}
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
			// we have the investor class, check password
			if rInvestor.LoginPassword == pwhash {
				log.Println("FOUDN INVESOTR")
				inv = rInvestor
			}
		}
		return fmt.Errorf("Not Found")
	})
	return inv, err
}
