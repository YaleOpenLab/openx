package ozones

import (
	database "github.com/YaleOpenLab/openx/database"
	platform "github.com/YaleOpenLab/openx/platforms"
)

type BondCoopParams struct {
	Index             int
	MaturationDate    string
	MemberRights      string
	SecurityType      string
	InterestRate      float64
	Rating            string
	BondIssuer        string
	Underwriter       string
	DateInitiated     string // date the project was created
	InvestorAssetCode string
	Title             string
	Description       string
	Location          string
}

// TODO: change name of bonds to something better. Add description of the bond platform below
// ConstructionBond contains the paramters for the Construction Bond model of the housing platform
// paramters defined here are not exhaustive and more can be added if desired
type ConstructionBond struct {
	Params BondCoopParams
	// common set of params that we need for openfinancing
	AmountRaised   float64
	CostOfUnit     float64
	InstrumentType string
	NoOfUnits      int
	Tax            string
	Investors      []database.Investor
	RecipientIndex int
}

// the coop struct uses the same base params as the bond model
type Coop struct {
	Params         BondCoopParams
	UnitsSold      int
	TotalAmount    float64
	TypeOfUnit     string
	MonthlyPayment float64
	Residents      []database.Investor
}

func InitializePlatform() error {
	return platform.InitializePlatform()
}

// RefillPlatform checks whether the publicKey passed has any xlm and if its balance
// is less than 21 XLM, it proceeds to ask the friendbot for more test xlm
func RefillPlatform(publicKey string) error {
	return platform.RefillPlatform(publicKey)
}
