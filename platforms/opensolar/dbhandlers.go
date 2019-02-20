package opensolar

import (
	"encoding/json"
	"github.com/pkg/errors"

	database "github.com/YaleOpenLab/openx/database"
	utils "github.com/YaleOpenLab/openx/utils"
	"github.com/boltdb/bolt"
)

// Save or Insert inserts a specific Project into the database
func (a *Project) Save() error {
	db, err := database.OpenDB()
	if err != nil {
		return errors.Wrap(err, "couldn't open db")
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ProjectsBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			return errors.Wrap(err, "couldn't marshal json")
		}
		return b.Put([]byte(utils.ItoB(a.Index)), encoded)
	})
	return err
}

// RetrieveProject retrieves the project with the specified index from the database
func RetrieveProject(key int) (Project, error) {
	var inv Project
	db, err := database.OpenDB()
	if err != nil {
		return inv, errors.Wrap(err, "couldn't open db")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ProjectsBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return errors.New("Retrieved project nil")
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, err
}

// RetrieveAllProjects retrieves all projects from the database
func RetrieveAllProjects() ([]Project, error) {
	var arr []Project
	db, err := database.OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "couldn't open db")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				break
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				return errors.Wrap(err, "couldn't marshal json")
			}
			arr = append(arr, rProject)
		}
		return nil
	})
	return arr, err
}

// RetrieveProjectsAtStage retrieves projects at a specific stage
func RetrieveProjectsAtStage(stage float64) ([]Project, error) {
	var arr []Project
	db, err := database.OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "couldn't open db")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				break
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				return errors.Wrap(err, "couldn't marshal json")
			}
			if rProject.Stage == stage {
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

// RetrieveContractorProjects retrieves projects that are associated with a specific contractor
func RetrieveContractorProjects(stage float64, index int) ([]Project, error) {
	var arr []Project
	db, err := database.OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "couldn't open db")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil // key does not exist
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				return errors.Wrap(err, "couldn't marshal json")
			}
			if rProject.Stage == stage && rProject.Contractor.U.Index == index {
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

// RetrieveOriginatorProjects retrieves projects that are associated with a specific originator
func RetrieveOriginatorProjects(stage float64, index int) ([]Project, error) {
	var arr []Project
	db, err := database.OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "couldn't open db")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				return errors.Wrap(err, "couldn't marshal json")
			}
			if rProject.Stage == stage && rProject.Originator.U.Index == index {
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

// RetrieveRecipientProjects retrieves projects that are associated with a specific originator
func RetrieveRecipientProjects(stage float64, index int) ([]Project, error) {
	var arr []Project
	db, err := database.OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "couldn't open db")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				return errors.Wrap(err, "couldn't marshal json")
			}
			if rProject.Stage == stage && rProject.RecipientIndex == index {
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

// RetrieveLockedProjects retrieves all the projects that are locked and are waiting
// for the recipient to unlock them
func RetrieveLockedProjects() ([]Project, error) {
	var arr []Project
	db, err := database.OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "couldn't open db")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				break
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				return errors.Wrap(err, "couldn't marshal json")
			}
			if rProject.Lock {
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

// SaveOriginatorMoU saves the MoU's hash in the platform's database
func SaveOriginatorMoU(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.OriginatorMoUHash = hash
	return a.Save()
}

// SaveContractHash saves a contract's hash in the platform's database
func SaveContractHash(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.ContractorContractHash = hash
	return a.Save()
}

// SaveInvPlatformContract saves the investor-platform contract's hash in the platform's database
func SaveInvPlatformContract(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.InvPlatformContractHash = hash
	return a.Save()
}

// SaveRecPlatformContract saves the recipient-platform contract's hash in the platform's database
func SaveRecPlatformContract(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.RecPlatformContractHash = hash
	return a.Save()
}
