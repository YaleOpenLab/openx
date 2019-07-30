package opensolar

import (
	//"log"
	"encoding/json"
	"github.com/pkg/errors"

	edb "github.com/Varunram/essentials/database"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
)

// Save or Insert inserts a specific Project into the database
func (a *Project) Save() error {
	return edb.Save(consts.DbDir, database.ProjectsBucket, a, a.Index)
}

// RetrieveProject retrieves the project with the specified index from the database
func RetrieveProject(key int) (Project, error) {
	var inv Project
	x, err := edb.Retrieve(consts.DbDir, database.ProjectsBucket, key)
	if err != nil {
		return inv, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &inv)
	return inv, err
}

// RetrieveAllProjects retrieves all projects from the database
func RetrieveAllProjects() ([]Project, error) {
	var projects []Project
	x, err := edb.RetrieveAllKeys(consts.DbDir, database.ProjectsBucket)
	if err != nil {
		return projects, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range x {
		var temp Project
		err = json.Unmarshal(value, &temp)
		if err != nil {
			return projects, errors.New("could not unmarshal json")
		}
		projects = append(projects, temp)
	}

	return projects, nil
}

// RetrieveProjectsAtStage retrieves projects at a specific stage
func RetrieveProjectsAtStage(stage int) ([]Project, error) {
	var arr []Project
	if stage > 9 { // check for this and fail early instead of wasting compute time on this
		return arr, errors.Wrap(errors.New("stage can not be greater than 9, quitting!"), "stage can not be greater than 9, quitting!")
	}

	projects, err := RetrieveAllProjects()
	if err != nil {
		return arr, err
	}

	for _, project := range projects {
		if project.Stage == stage {
			arr = append(arr, project)
		}
	}

	return arr, nil
}

// RetrieveContractorProjects retrieves projects that are associated with a specific contractor
func RetrieveContractorProjects(stage int, index int) ([]Project, error) {
	var arr []Project
	if stage > 9 { // check for this and fail early instead of wasting compute time on this
		return arr, errors.Wrap(errors.New("stage can not be greater than 9, quitting!"), "stage can not be greater than 9, quitting!")
	}

	projects, err := RetrieveAllProjects()
	if err != nil {
		return arr, err
	}

	for _, project := range projects {
		if project.Stage == stage && project.ContractorIndex == index {
			arr = append(arr, project)
		}
	}

	return arr, nil
}

// RetrieveOriginatorProjects retrieves projects that are associated with a specific originator
func RetrieveOriginatorProjects(stage int, index int) ([]Project, error) {
	var arr []Project
	if stage > 9 { // check for this and fail early instead of wasting compute time on this
		return arr, errors.Wrap(errors.New("stage can not be greater than 9, quitting!"), "stage can not be greater than 9, quitting!")
	}

	projects, err := RetrieveAllProjects()
	if err != nil {
		return arr, err
	}

	for _, project := range projects {
		if project.Stage == stage && project.OriginatorIndex == index {
			arr = append(arr, project)
		}
	}

	return arr, nil
}

// RetrieveRecipientProjects retrieves projects that are associated with a specific originator
func RetrieveRecipientProjects(stage int, index int) ([]Project, error) {
	var arr []Project
	if stage > 9 { // check for this and fail early instead of wasting compute time on this
		return arr, errors.Wrap(errors.New("stage can not be greater than 9, quitting!"), "stage can not be greater than 9, quitting!")
	}

	projects, err := RetrieveAllProjects()
	if err != nil {
		return arr, err
	}

	for _, project := range projects {
		if project.Stage == stage && project.RecipientIndex == index {
			arr = append(arr, project)
		}
	}

	return arr, nil
}

// RetrieveLockedProjects retrieves all the projects that are locked and are waiting
// for the recipient to unlock them
func RetrieveLockedProjects() ([]Project, error) {
	var arr []Project

	projects, err := RetrieveAllProjects()
	if err != nil {
		return arr, err
	}

	for _, project := range projects {
		if project.Lock {
			arr = append(arr, project)
		}
	}

	return arr, nil
}

// SaveOriginatorMoU saves the MoU's hash in the platform's database
func SaveOriginatorMoU(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.StageData = append(a.StageData, hash)
	return a.Save()
}

// SaveContractHash saves a contract's hash in the platform's database
func SaveContractHash(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.StageData = append(a.StageData, hash)
	return a.Save()
}

// SaveInvPlatformContract saves the investor-platform contract's hash in the platform's database
func SaveInvPlatformContract(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.StageData = append(a.StageData, hash)
	return a.Save()
}

// SaveRecPlatformContract saves the recipient-platform contract's hash in the platform's database
func SaveRecPlatformContract(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.StageData = append(a.StageData, hash)
	return a.Save()
}
