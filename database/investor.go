package database

import (
	"github.com/pkg/errors"

	oracle "github.com/YaleOpenLab/openx/oracle"
	utils "github.com/YaleOpenLab/openx/utils"
	xlm "github.com/YaleOpenLab/openx/xlm"
	"github.com/boltdb/bolt"
)

// the investor struct contains all the investor details such as
// public key, seed (if account is created on the website) and other stuff which
// is yet to be decided

// All investors will be referenced by their public key, name is optional (maybe necessary?)
// we need to still decide on identity and stuff and how much we want to track
// people who invest in projects

// Investor defines the investor structure
type Investor struct {
	VotingBalance int // this will be equal to the amount of stablecoins that the investor possesses,
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
// name and an option to generate the seed and publicKey. This is done because if
// we decide to allow anonymous investors to invest on our platform, we can easily
// insert their publickey into the system and then have hanlders for them signing
// transactions
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
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(InvestorBucket)
		encoded, err := a.MarshalJSON()
		if err != nil {
			return errors.Wrap(err, "error while marshaling json struct")
		}
		return b.Put([]byte(utils.ItoB(a.U.Index)), encoded)
	})
	return err
}

// RetrieveInvestorHelper is a helper function associated with the RetrieveInvestor function
func RetrieveInvestorHelper(key int) (Investor, error) {
	var inv Investor
	db, err := OpenDB()
	if err != nil {
		return inv, errors.Wrap(err, "failed to open db")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(InvestorBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			// no investor with the specific details
			return errors.New("No investor found with required credentials")
		}
		return inv.UnmarshalJSON(x)
	})
	return inv, err
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

// RetrieveAllInvestors gets a list of all investors in the database
func RetrieveAllInvestors() ([]Investor, error) {
	// this route is broken because it reads through keys sequentially
	// need to see keys until the length of the users database
	var arr []Investor
	temp, err := RetrieveAllUsers()
	if err != nil {
		return arr, errors.Wrap(err, "failed to retrieve all users")
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "failed to open db")
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(InvestorBucket)
		for i := 1; i < limit; i++ {
			var rInvestor Investor
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// this is where the key does not exist, we search until limit, so don't error out
				continue
			}
			err := rInvestor.UnmarshalJSON(x)
			if err != nil {
				// error in unmarshalling this struct, error out
				return errors.Wrap(err, "failed to unmarshal json")
			}
			arr = append(arr, rInvestor)
		}
		return nil
	})
	return arr, err
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

// DeductVotingBalance deducts voting balance
func (a *Investor) DeductVotingBalance(votes int) error {
	// maybe once the user can presses the vote button manually, we can fetch balance
	// and show him available votes onthe frontend
	a.VotingBalance -= votes
	return a.Save()
}

// AddVotingBalance adds voting balance
func (a *Investor) AddVotingBalance(votes int) error {
	// this function is caled when we want to refund the user with the votes once
	// an order has been finalized.
	a.VotingBalance += votes
	return a.Save()
}

// CanInvest checks whether an investor has the required balance to invest in a project
func (a *Investor) CanInvest(targetBalance string) bool {
	usdBalance, err := xlm.GetAssetBalance(a.U.PublicKey, "STABLEUSD")
	if err != nil {
		usdBalance = "0"
	}

	xlmBalance, err := xlm.GetNativeBalance(a.U.PublicKey)
	if err != nil {
		xlmBalance = "0"
	}

	// need to fetch the oracle price here for the order
	oraclePrice := oracle.ExchangeXLMforUSD(xlmBalance)
	if (utils.StoF(usdBalance) > utils.StoF(targetBalance)) || oraclePrice > utils.StoF(targetBalance) {
		// return true since the user has enough USD balance to pay for the order
		return true
	}
	return false
}

// the following two functions on reputation are repeated for recipients and entities
// but are necessary for th RPC which woukd call these functions in various scenarios
// eg. when negative feedback is approved  by multiple parties and they decide to
// reduce the reputation of the user

// ChangeInvReputation changes the investor's reputation
func ChangeInvReputation(invIndex int, reputation float64) error {
	a, err := RetrieveInvestor(invIndex)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve investor")
	}
	if reputation > 0 {
		err = a.U.IncreaseReputation(reputation)
	} else {
		err = a.U.DecreaseReputation(reputation)
	}
	if err != nil {
		return errors.Wrap(err, "Error while changing reputation")
	}
	return a.Save()
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
