package database

import (
	edb "github.com/Varunram/essentials/database"
	consts "github.com/YaleOpenLab/openx/consts"
	"github.com/boltdb/bolt"
)

// the database package maintains read / write operations to the orderbook
// we need an orderbook because there is no state on Stellar which makes it
// difficult for us to store this on the blockchain. We use boltdb now since we don't
// do that much relational mapping, but in the case we need that, we can modify
// this package to do that.
// package database contains base level stuff that will be required by all the
// sub platforms in the platform of platforms model. Currently contains Investors,
// Recipients and Users, but can be expanded to contain even stages, if that's deemed
// to be common across platforms
// define the name of the buckets that we interact with.

var ProjectsBucket = []byte("Projects")
var InvestorBucket = []byte("Investors")
var RecipientBucket = []byte("Recipients")
var ContractorBucket = []byte("Contractors")
var UserBucket = []byte("Users")
var BondBucket = []byte("Bonds")
var CoopBucket = []byte("Coop")
var InspectorBucket = []byte("Inspector")

// CreateHomeDir creates a home directory
func CreateHomeDir() {
	edb.CreateDirs(consts.HomeDir, consts.DbDir, consts.OpenSolarIssuerDir, consts.OpzonesIssuerDir)
}

// OpenDB opens the db
func OpenDB() (*bolt.DB, error) {
	return edb.CreateDB(consts.DbDir, ProjectsBucket, InvestorBucket, RecipientBucket, ContractorBucket, UserBucket, BondBucket, CoopBucket, InspectorBucket)
}

// DeleteKeyFromBucket deletes a given key from the bucket
func DeleteKeyFromBucket(key int, bucketName []byte) error {
	return edb.DeleteKeyFromBucket(consts.DbDir, key, bucketName)
}
