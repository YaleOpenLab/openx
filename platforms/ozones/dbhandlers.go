package ozones

import (
	"encoding/json"
	"fmt"
	"log"

	database "github.com/YaleOpenLab/openx/database"
	utils "github.com/YaleOpenLab/openx/utils"
	"github.com/boltdb/bolt"
)

// newParams defiens a common function for all the sub parts of the open housing platform. Can be thoguht
// of more like a common subset on which paramters for different models are defined on
func newParams(mdate string, mrights string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, title string, location string, description string) BondCoopParams {
	var rParams BondCoopParams
	rParams.MaturationDate = mdate
	rParams.MemberRights = mrights
	rParams.SecurityType = stype
	rParams.InterestRate = intrate
	rParams.Rating = rating
	rParams.BondIssuer = bIssuer
	rParams.Underwriter = uWriter
	rParams.Title = title
	rParams.Location = location
	rParams.Description = description
	rParams.DateInitiated = utils.Timestamp()
	return rParams
}

// NewCoop returns a new living coop and automatically saves it
func NewCoop(mdate string, mrights string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, totalAmount float64, typeOfUnit string, monthlyPayment float64,
	title string, location string, description string) (Coop, error) {
	var cCoop Coop
	cCoop.Params = newParams(mdate, mrights, stype, intrate, rating, bIssuer, uWriter, title, location, description)
	x, err := RetrieveAllCoops()
	if err != nil {
		return cCoop, err
	}

	cCoop.Params.Index = len(x) + 1
	cCoop.UnitsSold = 0
	cCoop.TotalAmount = totalAmount
	cCoop.TypeOfUnit = typeOfUnit
	cCoop.MonthlyPayment = monthlyPayment
	err = cCoop.Save()
	return cCoop, err
}

// NewBond returns a New Construction Bond and automatically stores it in the db
func NewBond(mdate string, mrights string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, unitCost float64, itype string, nUnits int, tax string, recIndex int,
	title string, location string, description string) (ConstructionBond, error) {
	var cBond ConstructionBond
	cBond.Params = newParams(mdate, mrights, stype, intrate, rating, bIssuer, uWriter, title, location, description)
	x, err := RetrieveAllBonds()
	if err != nil {
		return cBond, err
	}

	cBond.Params.Index = len(x) + 1
	cBond.CostOfUnit = unitCost
	cBond.InstrumentType = itype
	cBond.NoOfUnits = nUnits
	cBond.Tax = tax
	cBond.RecipientIndex = recIndex
	err = cBond.Save()
	return cBond, err
}

func (a *Coop) Save() error {
	db, err := database.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.CoopBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.ItoB(a.Params.Index)), encoded)
	})
	return err
}

func (a *ConstructionBond) Save() error {
	db, err := database.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.BondBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.ItoB(a.Params.Index)), encoded)
	})
	return err
}

// RetrieveAllCoops gets a list of all User in the database
func RetrieveAllCoops() ([]Coop, error) {
	var arr []Coop
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.CoopBucket)
		for i := 1; ; i++ {
			var rCoop Coop
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rCoop)
			if err != nil {
				return err
			}
			arr = append(arr, rCoop)
		}
		return nil
	})
	return arr, err
}

// RetrieveAllBonds gets a list of all User in the database
func RetrieveAllBonds() ([]ConstructionBond, error) {
	var arr []ConstructionBond
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
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
				return err
			}
			arr = append(arr, rBond)
		}
		return nil
	})
	return arr, err
}

// RetrieveCoop retrieves a specifi coop from the database
func RetrieveCoop(key int) (Coop, error) {
	var bond Coop
	db, err := database.OpenDB()
	if err != nil {
		return bond, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.CoopBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return fmt.Errorf("Retrieved Coop nil")
		}
		return json.Unmarshal(x, &bond)
	})
	return bond, err
}

func RetrieveBond(key int) (ConstructionBond, error) {
	var bond ConstructionBond
	db, err := database.OpenDB()
	if err != nil {
		return bond, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.BondBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return fmt.Errorf("Retreived Bond returns nil")
		}
		return json.Unmarshal(x, &bond)
	})
	return bond, err
}
