package ozones

import (
	"encoding/json"
	"github.com/pkg/errors"

	database "github.com/YaleOpenLab/openx/database"
	utils "github.com/YaleOpenLab/openx/utils"
	"github.com/boltdb/bolt"
)

// NewLivingUnitCoop returns a new living coop and automatically saves it
func NewLivingUnitCoop(mdate string, mrights string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, totalAmount float64, typeOfUnit string, monthlyPayment float64,
	title string, location string, description string) (LivingUnitCoop, error) {
	var coop LivingUnitCoop
	coop.MaturationDate = mdate
	coop.MemberRights = mrights
	coop.SecurityType = stype
	coop.InterestRate = intrate
	coop.Rating = rating
	coop.BondIssuer = bIssuer
	coop.Underwriter = uWriter
	coop.Title = title
	coop.Location = location
	coop.Description = description
	coop.DateInitiated = utils.Timestamp()

	x, err := RetrieveAllLivingUnitCoops()
	if err != nil {
		return coop, errors.Wrap(err, "could not retrieve all living unit coops")
	}
	coop.Index = len(x) + 1
	coop.UnitsSold = 0
	coop.Amount = totalAmount
	coop.TypeOfUnit = typeOfUnit
	coop.MonthlyPayment = monthlyPayment
	err = coop.Save()
	return coop, err
}

// NewConstructionBond returns a New Construction Bond and automatically stores it in the db
func NewConstructionBond(mdate string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, unitCost float64, itype string, nUnits int, tax string, recIndex int,
	title string, location string, description string) (ConstructionBond, error) {
	var cBond ConstructionBond
	cBond.MaturationDate = mdate
	cBond.SecurityType = stype
	cBond.InterestRate = intrate
	cBond.Rating = rating
	cBond.BondIssuer = bIssuer
	cBond.Underwriter = uWriter
	cBond.Title = title
	cBond.Location = location
	cBond.Description = description
	cBond.DateInitiated = utils.Timestamp()

	x, err := RetrieveAllConstructionBonds()
	if err != nil {
		return cBond, errors.Wrap(err, "could not retrieve all living unit coops")
	}

	cBond.Index = len(x) + 1
	cBond.CostOfUnit = unitCost
	cBond.InstrumentType = itype
	cBond.NoOfUnits = nUnits
	cBond.Tax = tax
	cBond.RecipientIndex = recIndex
	err = cBond.Save()
	return cBond, err
}

func (a *LivingUnitCoop) Save() error {
	db, err := database.OpenDB()
	if err != nil {
		return errors.Wrap(err, "coild not open db")
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.CoopBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			return errors.Wrap(err, "Failed to marshal json")
		}
		return b.Put([]byte(utils.ItoB(a.Index)), encoded)
	})
	return err
}

func (a *ConstructionBond) Save() error {
	db, err := database.OpenDB()
	if err != nil {
		return errors.Wrap(err, "coild not open db")
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.BondBucket)
		encoded, err := json.Marshal(a)
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
		return arr, errors.Wrap(err, "coild not open db")
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
			err := json.Unmarshal(x, &rCoop)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal json")
			}
			arr = append(arr, rCoop)
		}
		return nil
	})
	return arr, err
}

// RetrieveAllConstructionBonds gets a list of all User in the database
func RetrieveAllConstructionBonds() ([]ConstructionBond, error) {
	var arr []ConstructionBond
	db, err := database.OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "coild not open db")
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
			err := json.Unmarshal(x, &rBond)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal json")
			}
			arr = append(arr, rBond)
		}
		return nil
	})
	return arr, err
}

// RetrieveCoop retrieves a specifi coop from the database
func RetrieveLivingUnitCoop(key int) (LivingUnitCoop, error) {
	var bond LivingUnitCoop
	db, err := database.OpenDB()
	if err != nil {
		return bond, errors.Wrap(err, "coild not open db")
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.CoopBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return errors.New("Retrieved LivingUnitCoop nil")
		}
		return json.Unmarshal(x, &bond)
	})
	return bond, err
}

func RetrieveConstructionBond(key int) (ConstructionBond, error) {
	var bond ConstructionBond
	db, err := database.OpenDB()
	if err != nil {
		return bond, errors.Wrap(err, "coild not open db")
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.BondBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return errors.New("Retreived Bond returns nil")
		}
		return json.Unmarshal(x, &bond)
	})
	return bond, err
}
