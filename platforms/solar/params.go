package solar

import (
	database "github.com/OpenFinancing/openfinancing/database"
)

type SolarParams struct {
	Index int // an Index to keep quick track of how many projects exist

	PanelSize   string // size of the given panel, for diplsaying to the user who wants to bid stuff
	TotalValue  int    // the total money that we need from investors
	Location    string // where this specific solar panel is located
	MoneyRaised int    // total money that has been raised until now
	Years       int    // number of years the recipient has chosen to opt for
	Metadata    string // any other metadata can be stored here

	Votes int // the number of votes towards a proposed contract by investors

	// once all funds have been raised, we need to set assetCodes
	INVAssetCode string
	DEBAssetCode string
	PBAssetCode  string

	BalLeft float64 // denotes the balance left to pay by the party

	DateInitiated string // date the project was created
	DateFunded    string // date when the project was funded
	DateLastPaid  string // date the project was last paid

	ProjectRecipient database.Recipient
	ProjectInvestors []database.Investor // TODO: get feedback on whether this is in its right place and whether this should be moved to contract
	// Percentage raised is not stored in the database since that can be calculated by the UI
}
