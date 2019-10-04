package database

import (
	"encoding/json"
	"github.com/pkg/errors"

	edb "github.com/Varunram/essentials/database"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
)

// Platform is a struct which holds all platform related info
type Platform struct {
	Index   int
	Name    string
	Code    string
	Timeout int64
}

// NewPlatform creates a new platform and stores it in the database
func NewPlatform(name string, code string, timeout bool) error {
	index, err := RetrieveAllPfLim()
	if err != nil {
		return errors.Wrap(err, "could not retrieve all keys from the database")
	}

	var x Platform
	x.Index = index + 1
	x.Name = name
	x.Code = code
	x.Timeout = utils.Unix() + 2600000 // 10 months
	// timeout is set to true by default
	return x.Save()
}

// Save inserts a Platform object into the database
func (a *Platform) Save() error {
	return edb.Save(consts.DbDir+consts.DbName, PlatformBucket, a, a.Index)
}

// RetrievePlatform retrieves a Platform from the database
func RetrievePlatform(key int) (Platform, error) {
	var pf Platform
	x, err := edb.Retrieve(consts.DbDir+consts.DbName, PlatformBucket, key)
	if err != nil {
		return pf, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &pf)
	return pf, err
}

// RetrieveAllPlatforms retrieves all platforms from the database
func RetrieveAllPlatforms() ([]Platform, error) {
	var arr []Platform
	x, err := edb.RetrieveAllKeys(consts.DbDir+consts.DbName, PlatformBucket)
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all platforms")
	}
	for _, value := range x {
		var temp Platform
		err := json.Unmarshal(value, &temp)
		if err != nil {
			return arr, errors.New("error while unmarshalling json, quitting")
		}
		arr = append(arr, temp)
	}

	return arr, nil
}

// RetrieveAllPfLim gets the number of platforms in the platform bucket
func RetrieveAllPfLim() (int, error) {
	return edb.RetrieveAllKeysLim(consts.DbDir+consts.DbName, PlatformBucket)
}
