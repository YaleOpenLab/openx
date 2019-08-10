package database

import (
	edb "github.com/Varunram/essentials/database"
	consts "github.com/YaleOpenLab/openx/consts"
	"github.com/boltdb/bolt"
)

// the database package contains the handlers necesssary for openx to interact with the
// underlying boltdb database

var UserBucket = []byte("Users")

// CreateHomeDir creates the home and database directories
func CreateHomeDir() {
	edb.CreateDirs(consts.HomeDir, consts.DbDir)
	db, _ := edb.CreateDB(consts.DbDir+consts.DbName, UserBucket)
	db.Close()
}

// OpenDB opens the db and returns a pointer to the database
func OpenDB() (*bolt.DB, error) {
	return edb.OpenDB(consts.DbDir + consts.DbName)
}

// DeleteKeyFromBucket deletes an object from the passed bucket
func DeleteKeyFromBucket(key int, bucketName []byte) error {
	return edb.DeleteKeyFromBucket(consts.DbDir+consts.DbName, key, bucketName)
}
