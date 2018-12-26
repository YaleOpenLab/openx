package database

import (
	"fmt"
)

type ContractAuction struct {
	// TODO: this struct isn't used yet as it needs handlers and stuff, but when
	// we move off main.go for testinge, this must be  used in order to make stuff
	// easier for us.
	// this is called when there is an originated order live and when there are
	// contractors who want to get this price. This is a blind auction and the
	// choosing criteria is just price for now.
	// TODO: decide this criteria
	AllContracts    []Contract
	AllContractors  []ContractEntity
	WinningContract Contract // do we need this?
}

// auctions contains stuff related to choosing the best contract and potentially
// future auction logic that might need to be housed here
func SelectContractByPrice(arr []Contract) (Contract, error) {
	var a Contract
	if len(arr) == 0 {
		return a, fmt.Errorf("Empty array passed!")
	}
	// array is not empty, min 1 elem
	a = arr[0]
	for _, elem := range arr {
		if elem.O.TotalValue < a.O.TotalValue {
			a = elem
			continue
		}
	}
	return a, nil
}

func SelectContractByTime(arr []Contract) (Contract, error) {
	var a Contract
	if len(arr) == 0 {
		return a, fmt.Errorf("Empty array passed!")
	}
	// array is not empty, min 1 elem
	a = arr[0]
	for _, elem := range arr {
		if elem.O.Years < a.O.Years {
			a = elem
			continue
		}
	}
	return a, nil
}
