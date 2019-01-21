package solar

import (
	"fmt"
)

// Contract auctions are specific for public infrastructure and for projects with multiple stakeholders
// where transparency is important. The main point is to avoid corruption and produce positive competition
// among providers. When funding solar project using international investor money, transparency and auditability
// to them is  also a crucial aspect since they are often not close to the project, and want to make sure there
// is a tobust due diligence prior to unlocking the funds.

// Different Auctions or Tenders are designed based on the nature of the project.
// In general, the criteria for selection is price, technical quality (eg. hardware), engineering model, development time
// and other perks offered by developers (eg. extra guarantees).

// you need to have a lock in period beyond which contractors can not post what
// stuff they want. now, how do you choose which contractor wins? Ideally,
// the school would want the most stuff but you need to vet which contracts are good
// and not.

type ContractAuction struct {
	// TODO: this struct isn't used yet as it needs handlers and stuff, but when
	// we move off main.go for testing, this must be used in order to make stuff
	// easier for us.
	AllContracts    []Project
	AllContractors  []Entity
	WinningContract Project // do we need this?
}

// auctions contains stuff related to choosing the best contract and potentially
// future auction logic that might need to be housed here
func SelectContractBlind(arr []Project) (Project, error) {
	// in a blind auction, the bid with the highest price wins
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

func SelectContractVickrey(arr []Project) (Project, error) {
	var winningContract Project
	if len(arr) == 0 {
		return winningContract, fmt.Errorf("Empty array passed!")
	}
	// array is not empty, min 1 elem
	winningContract = arr[0]
	var pos int
	for i, elem := range arr {
		if elem.Params.TotalValue < winningContract.Params.TotalValue {
			winningContract = elem
			pos = i
			continue
		}
	}
	// here we have the highest bidder. Now we need to delete this guy from the array
	// and get the second highest bidder
	// delete a[pos] from arr
	arr = append(arr[:pos], arr[pos+1:]...)
	if len(arr) == 0 {
		// means only one contract was proposed for this project, so fall back to blind auction
		return winningContract, nil
	}
	vickreyPrice := arr[0].Params.TotalValue
	for _, elem := range arr {
		if elem.Params.TotalValue < vickreyPrice {
			vickreyPrice = elem.Params.TotalValue
		}
	}
	// we have the winner, who's elem and we have the price which is vickreyPrice
	// voerwrite the winning contractor's contract
	winningContract.Params.TotalValue = vickreyPrice
	return winningContract, winningContract.Save()
}

func SelectContractTime(arr []Project) (Project, error) {
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

func (project *Project) SetAuctionType(auctionType string) error {
	// see https://en.wikipedia.org/wiki/Auction for primary auction types
	switch auctionType {
	case "blind":
		project.AuctionType = "blind"
	case "vickrey":
		project.AuctionType = "vickrey"
	case "english":
		project.AuctionType = "english"
	case "dutch":
		project.AuctionType = "dutch"
	default:
		project.AuctionType = "blind"
	}
	return project.Save()
}
