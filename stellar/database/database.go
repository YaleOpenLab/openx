package database

import (
	"log"

	"github.com/boltdb/bolt"
)

var OrdersBucket = []byte("Orders")
var InvestorBucket = []byte("Investors")
var RecipientBucket = []byte("Recipients")
var PlatformBucket = []byte("Platforms")
var ContractorBucket = []byte("Contractors")

// TODO: need locks over this to ensure no one's using the db while we are
func OpenDB() (*bolt.DB, error) {
	db, err := bolt.Open("yol.db", 0600, nil)
	if err != nil {
		log.Println("Couldn't open database, exiting!")
		return db, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(OrdersBucket) // the orders bucket contains all our orders
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(InvestorBucket) // the orders bucket contains all our orders
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(RecipientBucket) // the orders bucket contains all our orders
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(PlatformBucket) // the orders bucket contains all our orders
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(ContractorBucket) // the orders bucket contains all our orders
		if err != nil {
			return err
		}
		return nil
	})
	return db, err
}
