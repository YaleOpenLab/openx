package bonds

// this param is shared by both bonds and coops, could add other stuff in here as well
type BondCoopParams struct {
	Index          int
	MaturationDate string
	MemberRights   string
	SecurityType   string
	InterestRate   float64
	Rating         string
	BondIssuer     string
	Underwriter    string
	DateInitiated  string // date the project was created
	INVAssetCode   string
	Title          string
	Description    string
	Location       string
}
