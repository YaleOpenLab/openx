/*
package ozones

/*
import (
	platform "github.com/YaleOpenLab/openx/platforms"
)

// An Opportunity Zone  has multiple forms of investment opportunities within it. SOme examples are
// Construction Bonds, Living Unit Coops, Utility Investments, etc. Ozones aims to start off with a construction
// bond and living unit coops facilitated by CityScope and then build more complex constructions like Utility
// Investments and DAOs for governane mechanisms on top of the construction we have right now.

// ConstructionBond defines the various parameters that constitute a Construction Bond
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
	InterestRate   float64 // the interest rateoffered for this particular bond
	Rating         string  // the moody's / finch's rating for this particular bond
	BondIssuer     string  // the issuing bank of this particular bond
	BondHolders    string
	Underwriter    string // the underwriter that will provide guarantee against defaults

	DateInitiated     string // date the project was created
	InvestorAssetCode string
	DebtAssetCode     string
	InvestorIndices   []int // the array of investors who have invested in this particular construction bond
	RecipientIndex    int   // the index of the recipient who ideally would be the person constructing this particular space
	LockPwd           string
	Lock              bool
}

// LivingUnitCoop defines the various parameters that constitute a Living Unit Coop
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
	MonthlyPayment float64 // monthly payment that must be m  ade towards this investment
	MemberRights   string  // the rights that the member of this coop is entitled to
	InterestRate   float64 // the interest rateoffered for this particular bond
	Rating         string  // the moody's / finch rating for this particular bond
	BondIssuer     string  // the issuing bank of this particular bond
	Underwriter    string  // the underwriter that will provide guarantee against defaults

	DateInitiated     string // date the project was created
	InvestorAssetCode string // the main receipt that the investor receives on investing in this living coop unit
	ResidentIndices   []int  // the indices of all residents (i nthis case investors as well) in this living unit coop

	RecipientIndex int
	LockPwd        string
	Lock           bool
}

// InitializePlatform borrows the init platform method from the common platform handler
func InitializePlatform() error {
	return platform.InitializePlatform()
}

// RefillPlatform checks whether the publicKey passed has any xlm and if its balance
// is less than 21 XLM, it proceeds to ask the friendbot for more test xlm
func RefillPlatform(publicKey string) error {
	return platform.RefillPlatform(publicKey)
}
*/

// testdata
/*
// newLivingUnitCoop creates a new living unit coop
func newLivingUnitCoop(mdate string, mrights string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, totalAmount float64, typeOfUnit string, monthlyPayment float64,
	title string, location string, description string) (opzones.LivingUnitCoop, error) {
	var coop opzones.LivingUnitCoop
	coop.MaturationDate = mdate
	coop.MemberRights = mrights
	coop.SecurityType = stype
	coop.InterestRate = intrate
	coop.Rating = rating
	coop.BondIssuer = bIssuer
	coop.Underwriter = uWriter
	coop.Title = title
	coop.Location = location
	coop.Description = description
	coop.DateInitiated = utils.Timestamp()

	x, err := opzones.RetrieveAllLivingUnitCoops()
	if err != nil {
		return coop, errors.Wrap(err, "could not retrieve all living unit coops")
	}
	coop.Index = len(x) + 1
	coop.UnitsSold = 0
	coop.Amount = totalAmount
	coop.TypeOfUnit = typeOfUnit
	coop.MonthlyPayment = monthlyPayment
	err = coop.Save()
	return coop, err
}

// newConstructionBond returns a New Construction Bond and automatically stores it in the db
func newConstructionBond(mdate string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, unitCost float64, itype string, nUnits int, tax string, recIndex int,
	title string, location string, description string) (opzones.ConstructionBond, error) {
	var cBond opzones.ConstructionBond
	cBond.MaturationDate = mdate
	cBond.SecurityType = stype
	cBond.InterestRate = intrate
	cBond.Rating = rating
	cBond.BondIssuer = bIssuer
	cBond.Underwriter = uWriter
	cBond.Title = title
	cBond.Location = location
	cBond.Description = description
	cBond.DateInitiated = utils.Timestamp()

	x, err := opzones.RetrieveAllConstructionBonds()
	if err != nil {
		return cBond, errors.Wrap(err, "could not retrieve all living unit coops")
	}

	cBond.Index = len(x) + 1
	cBond.CostOfUnit = unitCost
	cBond.InstrumentType = itype
	cBond.NoOfUnits = nUnits
	cBond.Tax = tax
	cBond.RecipientIndex = recIndex
	err = cBond.Save()
	return cBond, err
}

/*
	_, err = newConstructionBond("Dec 21 2021", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Opportunity Zone Construction", 200, "5% tax for 10 years", 1, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		return err
	}

	_, err = newConstructionBond("Apr 2 2025", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Opportunity Zone Construction", 400, "No tax for 20 years", 1, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		return err
	}

	_, err = newConstructionBond("Jul 9 2029", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Opportunity Zone Construction", 100, "3% Tax for 5 Years", 1, "Osaka Development Project", "Osaka", "This Project seeks to develop cutting edge technologies in Osaka")
	if err != nil {
		return err
	}

	_, err = newLivingUnitCoop("Dec 21 2021", "Member Rights Link", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Coop Model", 4000, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		return err
	}

	_, err = newLivingUnitCoop("Apr 2 2025", "Member Rights Link", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Coop Model", 1000, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		return err
	}

	_, err = newLivingUnitCoop("Jul 9 2029", "Member Rights Link", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Coop Model", 2000, "Osaka Development Project", "Osaka", "ODP seeks to develop cutting edge technologies in Osaka and invites investors all around the world to be a part of this new age")
	if err != nil {
		return err
	}
*/
