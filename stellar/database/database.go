package database

import (
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
)

// defines the name of the buckets that we interact with.
var ProjectsBucket = []byte("Projects")
var InvestorBucket = []byte("Investors")
var RecipientBucket = []byte("Recipients")
var PlatformBucket = []byte("Platforms")
var ContractorBucket = []byte("Contractors")
var UserBucket = []byte("Users")

// TODO: need locks over this to ensure no one's using the db while we are
func OpenDB() (*bolt.DB, error) {
	db, err := bolt.Open("yol.db", 0600, nil)
	if err != nil {
		log.Println("Couldn't open database, exiting!")
		return db, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(ProjectsBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(InvestorBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(RecipientBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(PlatformBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(ContractorBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(UserBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		return nil
	})
	return db, err
}

// DeleteProject deltes a given value corresponding to the ky from the database
// DeleteProject should be used only in cases where something is wrong from our side
// while creating an project. For other cases, we should set Live to False and edit
// the project
// TODO: make delete not mess up with indices, which it currently does
// DeleteKeyFromBucket deletes a given key from the bucket _bucketName
func DeleteKeyFromBucket(key int, bucketName []byte) error {
	// deleting project might be dangerous since that would mess with the other
	// functions, have it in here for now, don't do too much with it / fiox retrieve all
	// to handle this case
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		err := b.Delete(utils.ItoB(key))
		if err != nil {
			return err
		}
		log.Printf("Deleted element with key: %d in bucket %s", key, string(bucketName))
		return nil
	})
	return err
}
