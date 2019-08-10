/*
package ozones

/*
import (
	"encoding/json"
	"github.com/pkg/errors"

	edb "github.com/Varunram/essentials/database"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
)

// Save saves the changes in a living unit coop
func (a *LivingUnitCoop) Save() error {
	return edb.Save(consts.DbDir+consts.DbName, database.CoopBucket, a, a.Index)
}

// Save saves the changes in a construction bond
func (a *ConstructionBond) Save() error {
	return edb.Save(consts.DbDir+consts.DbName, database.BondBucket, a, a.Index)
}

// RetrieveAllLivingUnitCoops gets a list of all User in the database
func RetrieveAllLivingUnitCoops() ([]LivingUnitCoop, error) {
	var arr []LivingUnitCoop
	x, err := edb.RetrieveAllKeys(consts.DbDir+consts.DbName, database.CoopBucket)
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range x {
		var temp LivingUnitCoop
		err := json.Unmarshal(value, &temp)
		if err != nil {
			return arr, errors.New("error while unmarshalling json, quitting")
		}
		arr = append(arr, temp)
	}

	return arr, nil
}

// RetrieveAllConstructionBonds gets a list of all User in the database
func RetrieveAllConstructionBonds() ([]ConstructionBond, error) {
	var arr []ConstructionBond
	x, err := edb.RetrieveAllKeys(consts.DbDir+consts.DbName, database.BondBucket)
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range x {
		var temp ConstructionBond
		err := json.Unmarshal(value, &temp)
		if err != nil {
			return arr, errors.New("error while unmarshalling json, quitting")
		}
		arr = append(arr, temp)
	}

	return arr, nil
}

// RetrieveLivingUnitCoop retrieves a specifi coop from the database
func RetrieveLivingUnitCoop(key int) (LivingUnitCoop, error) {
	var elem LivingUnitCoop
	x, err := edb.Retrieve(consts.DbDir+consts.DbName, database.CoopBucket, key)
	if err != nil {
		return elem, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &elem)
	return elem, err
}

// RetrieveConstructionBond retrieves the construction bond from memory
func RetrieveConstructionBond(key int) (ConstructionBond, error) {
	var elem ConstructionBond
	x, err := edb.Retrieve(consts.DbDir+consts.DbName, database.BondBucket, key)
	if err != nil {
		return elem, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &elem)
	return elem, err
}
*/
