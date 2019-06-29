package ozones

import (
	"github.com/pkg/errors"

	utils "github.com/Varunram/essentials/utils"
	database "github.com/YaleOpenLab/openx/database"
	"github.com/boltdb/bolt"
)

// Save saves the changes in a living unit coop
func (a *LivingUnitCoop) Save() error {
	db, err := database.OpenDB()
	if err != nil {
		return errors.Wrap(err, "could not open db")
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.CoopBucket)
		encoded, err := a.MarshalJSON()
		if err != nil {
			return errors.Wrap(err, "failed to marshal json")
		}
		return b.Put([]byte(utils.ItoB(a.Index)), encoded)
	})
	return err
}

// Save saves the changes in a construction bond
func (a *ConstructionBond) Save() error {
	db, err := database.OpenDB()
	if err != nil {
		return errors.Wrap(err, "could not open db")
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.BondBucket)
		encoded, err := a.MarshalJSON()
		if err != nil {
			return errors.Wrap(err, "failed to marshal json")
		}
		return b.Put([]byte(utils.ItoB(a.Index)), encoded)
	})
	return err
}

// RetrieveAllLivingUnitCoops gets a list of all User in the database
func RetrieveAllLivingUnitCoops() ([]LivingUnitCoop, error) {
	var arr []LivingUnitCoop
	db, err := database.OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "could not open db")
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.CoopBucket)
		for i := 1; ; i++ {
			var rCoop LivingUnitCoop
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := rCoop.UnmarshalJSON(x)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal json")
			}
			arr = append(arr, rCoop)
		}
	})
	return arr, err
}

// RetrieveAllConstructionBonds gets a list of all User in the database
func RetrieveAllConstructionBonds() ([]ConstructionBond, error) {
	var arr []ConstructionBond
	db, err := database.OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "could not open db")
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.BondBucket)
		for i := 1; ; i++ {
			var rBond ConstructionBond
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := rBond.UnmarshalJSON(x)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal json")
			}
			arr = append(arr, rBond)
		}
	})
	return arr, err
}

// RetrieveLivingUnitCoop retrieves a specifi coop from the database
func RetrieveLivingUnitCoop(key int) (LivingUnitCoop, error) {
	var bond LivingUnitCoop
	db, err := database.OpenDB()
	if err != nil {
		return bond, errors.Wrap(err, "could not open db")
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.CoopBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return errors.New("Retrieved LivingUnitCoop nil")
		}
		return bond.UnmarshalJSON(x)
	})
	return bond, err
}

// RetrieveConstructionBond retrieves the construction bond from memory
func RetrieveConstructionBond(key int) (ConstructionBond, error) {
	var bond ConstructionBond
	db, err := database.OpenDB()
	if err != nil {
		return bond, errors.Wrap(err, "could not open db")
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.BondBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return errors.New("Retrieved Bond returns nil")
		}
		return bond.UnmarshalJSON(x)
	})
	return bond, err
}
