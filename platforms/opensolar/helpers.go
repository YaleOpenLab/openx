package opensolar

import (
	"github.com/pkg/errors"
)

// findInKey finds a project within an array of projects, given a key or index
func findInKey(key int, arr []Project) (Project, error) {
	var dummy Project
	for _, elem := range arr {
		if elem.Index == key {
			return elem, nil
		}
	}
	return dummy, errors.New("Not found")
}
