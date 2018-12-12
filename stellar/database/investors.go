package database

// this should actually be in the databse package, but since this needs the
import (
	"log"
	"encoding/json"

	"github.com/stellar/go/keypair"
	utils "github.com/Varunram/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
)

// the investor s truct contains all the investor details such as
// public key, seed (if account is created on the website) and ot her stuff which
// is yet to be decided

// ALl investors will be referenced by their public key, name is optional (maybe necessary?)
// we need to stil ldecide on identity and stuff and how much we want to track
// people who invest in the schools
type Investor struct {
	Index uint32
	// defauult index, gets us easy s tats on how many  people are there and stuff,
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

func NewInvestor(uname string, pwhash string, Name string, pkgen bool) error {
	// call this after the user has failled in username and password. Store hashed password
	// in the database
	var a Investor
	a.Name = Name
	if pkgen {
		// generate a pk and seed pair and store it
		pair, err := keypair.Random()
		if err != nil {
			return err
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
	// now we have a new investore, take this and then send this off to be stored in the database
	return nil
}

func CheckPassword(pwHash string) (bool, error) {
	// frontend should serve sha3(password) to us, we compare that and what's stored
	// in the database to s ee if they match
	return true, nil
}

func Testfn() string {
	log.Println("Testing this package")
	return "hello"
}

func InsertInvestor(a Investor, db *bolt.DB) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}

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
	})
	return err
}
