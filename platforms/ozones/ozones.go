package ozones

import (
	database "github.com/YaleOpenLab/openx/database"
	platform "github.com/YaleOpenLab/openx/platforms"
)

// An Opportunity Zone  has multiple forms of investment opportunities within it. SOme examples are
// Construction Bonds, Living Unit Coops, Utility Investments, etc. Ozones aims to start off with a construction
// bond and living unit coops facilitated by CityScope and then build more complex constructions like Utility
// INvestments and DAOs for governane mechanisms on top of the construction we have right now.

// TODO: add description of what the construction bond model does in ozones below with inputs from John

// define the various parameters that constitute a Construction Bond
type ConstructionBond struct {
	Index int

	Title          string
	Location       string
	Description    string
	AmountRaised   float64 // amount raised is what has been raised until now
	InstrumentType string  // Opportunity Zone Construction, 10 year
	Amount         string  // amount is something like $10 million units upto $200 million
	CostOfUnit     float64 // $10 million
	NoOfUnits      int     // 20 (since investment required is acpped at $200 million)
	SecurityType   string  // The class of security that this security falls under
	Tax            string
	MaturationDate string  // date at which the bond expires
	InterestRate   float64 // teh interest rateoffered for this particular bond
	Rating         string  // the moody's / finch's rating for this particular bond
	BondIssuer     string  // the issuing bank of this particular bond
	BondHolders    string
	Underwriter    string // the underwriter that will provide guarantee against defaults

	DateInitiated     string // date the project was created
	InvestorAssetCode string
	DebtAssetCode     string
	InvestorIndices   []int // the array of investors who have invested in this particular construction bond
	RecipientIndex    int   // the index of the recipient who ideally would be the person constructing this particular space

	// TODO: add more parameters based on discussions and feedback from Martin and John
}

// define the various parameters that constitute a Living Unit Coop
type LivingUnitCoop struct {
	Index int

	Title          string
	Location       string
	Description    string
	UnitsSold      int
	TypeOfUnit     string  // 2 bedroom transformable coop unit: 600 feet, see link
	Amount         float64 // amount that is required to be invested in this living unit coop
	SecurityType   string  // The class of security that this security falls under
	MaturationDate string  // date at which the bond expires
	MonthlyPayment float64
	MemberRights   string  // the rights that the member of this coop is entitled to
	InterestRate   float64 // teh interest rateoffered for this particular bond
	Rating         string  // the moody's / finch rating for this particular bond
	BondIssuer     string  // the issuing bank of this particular bond
	Underwriter    string  // the underwriter that will provide guarantee against defaults

	DateInitiated     string // date the project was created
	InvestorAssetCode string // the main receipt that the investor receives on investing in this living coop unit
	Residents         []database.Investor

	// TODO: add more parameters based on discussions and feedback from Martin and John
}

func InitializePlatform() error {
	return platform.InitializePlatform()
}

// RefillPlatform checks whether the publicKey passed has any xlm and if its balance
// is less than 21 XLM, it proceeds to ask the friendbot for more test xlm
func RefillPlatform(publicKey string) error {
	return platform.RefillPlatform(publicKey)
}
