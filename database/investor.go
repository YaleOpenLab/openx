package database

import (
	"github.com/pkg/errors"

	tickers "github.com/Varunram/essentials/crypto/exchangetickers"
	xlm "github.com/Varunram/essentials/crypto/xlm"
	edb "github.com/Varunram/essentials/database"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
)

// Investor defines the investor structure
type Investor struct {
	VotingBalance float64 // this will be equal to the amount of stablecoins that the investor possesses,
	// should update this every once in a while to ensure voting consistency.
	// These are votes to show opinions about bids done by contractors on the specific projects that investors invested in.
	AmountInvested float64
	// total amount, would be nice to track to contact them,
	// give them some kind of medals or something
	InvestedSolarProjects        []string
	InvestedSolarProjectsIndices []int
	InvestedBonds                []string
	InvestedCoops                []string
	// array of asset codes this user has invested in
	U           *User
	WeightedROI string
	// the weightedROI for all the projects under the investor's umbrella
	AllTimeReturns []float64
	// the all time returns accumulated by the investor during his time on the platform indexed by project index
	ReceivedRECs string
	// The renewable enrgy  certificated received by the investor as part o
	Prorata string
	// the pro rata in all the projects that the in vestor has invested in
}

// NewInvestor creates a new investor object when passed the username, password hash,
// name and an option to generate the seed and publicKey.
func NewInvestor(uname string, pwd string, seedpwd string, Name string) (Investor, error) {
	var a Investor
	var err error
	user, err := NewUser(uname, pwd, seedpwd, Name)
	if err != nil {
		return a, errors.Wrap(err, "error while creating a new user")
	}
	a.U = &user
	a.AmountInvested = float64(0)
	err = a.Save()
	return a, err
}

// Save inserts a passed Investor object into the database
func (a *Investor) Save() error {
	return edb.Save(consts.DbDir, InvestorBucket, a, a.U.Index)
}

// RetrieveInvestorHelper is a helper function associated with the RetrieveInvestor function
func RetrieveInvestorHelper(key int) (Investor, error) {

	var inv Investor
	x, err := edb.Retrieve(consts.DbDir, InvestorBucket, key)
	if err != nil {
		return inv, errors.Wrap(err, "error while retrieving key from bucket")
	}

	return x.(Investor), nil
}

// RetrieveInvestor retrieves a particular investor indexed by key from the database
func RetrieveInvestor(key int) (Investor, error) {
	var inv Investor
	user, err := RetrieveUser(key)
	if err != nil {
		return inv, err
	}
	inv, err = RetrieveInvestorHelper(key)
	if err != nil {
		return inv, err
	}
	inv.U = &user
	return inv, inv.Save()
}

// RetrieveAllUsers gets a list of all User in the database
func RetrieveAllInvestors() ([]Investor, error) {
	var investors []Investor
	x, err := edb.RetrieveAllKeys(consts.DbDir, InvestorBucket)
	if err != nil {
		return investors, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range x {
		investors = append(investors, value.(Investor))
	}

	return investors, nil
}

// ValidateInvestor is a function to validate the investors username and password to log them into the platform, and find the details related to the investor
// This is separate from the publicKey/seed pair (which are stored encrypted in the database); since we can help users change their password, but we can't help them retrieve their seed.
func ValidateInvestor(name string, pwhash string) (Investor, error) {
	var rec Investor
	user, err := ValidateUser(name, pwhash)
	if err != nil {
		return rec, errors.Wrap(err, "failed to validate user")
	}
	return RetrieveInvestor(user.Index)
}

// AddVotingBalance adds / subtracts voting balance
func (a *Investor) ChangeVotingBalance(votes float64) error {
	// this function is caled when we want to refund the user with the votes once
	// an order has been finalized.
	a.VotingBalance += votes
	if a.VotingBalance < 0 {
		a.VotingBalance = 0 // to ensure no one has negative votes or something
	}
	return a.Save()
}

// CanInvest checks whether an investor has the required balance to invest in a project
func (a *Investor) CanInvest(targetBalance string) bool {
	usdBalance, err := xlm.GetAssetBalance(a.U.StellarWallet.PublicKey, "STABLEUSD")
	if err != nil {
		usdBalance = "0"
	}

	xlmBalance, err := xlm.GetNativeBalance(a.U.StellarWallet.PublicKey)
	if err != nil {
		xlmBalance = "0"
	}

	// need to fetch the oracle price here for the order
	oraclePrice := tickers.ExchangeXLMforUSD(xlmBalance)
	if (utils.StoF(usdBalance) > utils.StoF(targetBalance)) || oraclePrice > utils.StoF(targetBalance) {
		// return true since the user has enough USD balance to pay for the order
		return true
	}
	return false
}

// TopReputationInvestors gets a list of all the investors with top reputation
func TopReputationInvestors() ([]Investor, error) {
	allInvestors, err := RetrieveAllInvestors()
	if err != nil {
		return allInvestors, errors.Wrap(err, "failed to retrieve all investors")
	}
	for i := range allInvestors {
		for j := range allInvestors {
			if allInvestors[i].U.Reputation > allInvestors[j].U.Reputation {
				tmp := allInvestors[i]
				allInvestors[i] = allInvestors[j]
				allInvestors[j] = tmp
			}
		}
	}
	return allInvestors, nil
}
