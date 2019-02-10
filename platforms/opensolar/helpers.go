package opensolar

import (
	"fmt"

	database "github.com/YaleOpenLab/openx/database"
)

// findInKey finds a project within an array of projects, given a key or index
func findInKey(key int, arr []Project) (Project, error) {
	var dummy Project
	for _, elem := range arr {
		if elem.Index == key {
			return elem, nil
		}
	}
	return dummy, fmt.Errorf("Not found")
}

// updateRecipient updstes the project's ProjectRecipient field
func (project *Project) updateRecipient(a database.Recipient) error {
	pos := -1
	for i, mem := range a.ReceivedSolarProjects {
		if mem == project.DebtAssetCode {
			// rewrite the thing in memory that we have
			pos = i
			break
		}
	}
	if pos != -1 {
		// rewrite the thing in memory
		a.ReceivedSolarProjects[pos] = project.DebtAssetCode
		err := a.Save()
		return err
	}
	return nil
}
