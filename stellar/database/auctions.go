package database

import (
	"fmt"
)

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
