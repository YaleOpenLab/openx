package database

import (
	"fmt"
)

// you need to have a lock in period beyond which contractors can not post what
// stuff they want. now, how do you choose which contractor wins? Ideally,
// the school would want the most stuff but you need to vet which contracts are good
// and not.

type ContractAuction struct {
	// TODO: this struct isn't used yet as it needs handlers and stuff, but when
	// we move off main.go for testinge, this must be  used in order to make stuff
	// easier for us.
	// this is called when there is an originated order live and when there are
	// contractors who want to get this price. This is a blind auction and the
	// choosing criteria is just price for now.
	// TODO: decide this criteria
	AllContracts    []Project
	AllContractors  []Entity
	WinningContract Project // do we need this?
}

// auctions contains stuff related to choosing the best contract and potentially
// future auction logic that might need to be housed here
func SelectContractByPrice(arr []Project) (Project, error) {
	var a Project
	if len(arr) == 0 {
		return a, fmt.Errorf("Empty array passed!")
	}
	// array is not empty, min 1 elem
	a = arr[0]
	for _, elem := range arr {
		if elem.Params.TotalValue < a.Params.TotalValue {
			a = elem
			continue
		}
	}
	return a, nil
}

func SelectContractByTime(arr []Project) (Project, error) {
	var a Project
	if len(arr) == 0 {
		return a, fmt.Errorf("Empty array passed!")
	}
	// array is not empty, min 1 elem
	a = arr[0]
	for _, elem := range arr {
		if elem.Params.Years < a.Params.Years {
			a = elem
			continue
		}
	}
	return a, nil
}
