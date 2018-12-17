package database

// Order is what is advertised on the frontend where investors can choose what
// projects to invest in.
type Order struct {
	Index         uint32  // an Index to keep quick track of how many orders exist
	PanelSize     string  // size of the given panel, for diplsaying to the user who wants to bid stuff
	TotalValue    int     // the total money that we need from investors
	Location      string  // where this specific solar panel is located
	MoneyRaised   int     // total money that has been raised until now
	Years         int     // number of years the recipient has chosen to opt for
	Metadata      string  // any other metadata can be stored here
	Live          bool    // check to see whether the current order is live or not
	INVAssetCode  string  // once all funds have been raised, we need to set assetCodes
	DEBAssetCode  string  // once all funds have been raised, we need to set assetCodes
	PBAssetCode   string  // once all funds have been raised, we need to set assetCodes
	BalLeft       float64 // denotes the balance left to pay by the party
	DateInitiated string  // date the order was created
	DateLastPaid  string  // date the order was last paid
	RecipientName string  // name of the recipient in order to assign the given assets
	// TODO: have an investor and recipient relation here
	// Percentage raised is not stored in the database since that can be calculated by the UI
}

type Recipient struct {
	Index uint32
	// defauult index, gets us easy stats on how many people are there and stuff,
	// don't want to omit this
	Name string
	// Name of the primary stakeholder involved (principal trustee of school, for eg.)
	PublicKey string
	// PublicKey denotes the public key of the recipient
	Seed string
	// Seed is the equivalent of a private key in stellar (stellar doesn't expose private keys)
	// do we make seed optional like that for the Recipient? Couple things to consider
	// here: if the recipient loses the publickey, it can nver send DEBTokens back
	// to the issuer, so it would be as if it reneged on the deal. Do we count on
	// technically less sound people to hold their public keys safely? I suggest
	// this would be  difficult in practice, so maybe enforce that they need to hold|
	// their account on the platform?
	FirstSignedUp string
	// auto generated timestamp
	ReceivedOrders []Order
	// ReceivedOrders denotes the orders that have been received by the recipient
	// instead of storing the PaybackAssets and the DebtAssets, we store this
	LoginUserName string
	// the thing you use to login to the platform
	LoginPassword string
	// password, which is separate from the generated seed.
}

// the investor s truct contains all the investor details such as
// public key, seed (if account is created on the website) and ot her stuff which
// is yet to be decided

// ALl investors will be referenced by their public key, name is optional (maybe necessary?)
// we need to stil ldecide on identity and stuff and how much we want to track
// people who invest in the schools
type Investor struct {
	Index uint32
	// defauult index, gets us easy stats on how many people are there and stuff,
	// don't want to omit this
	Name string
	// display Name, different from UserName
	PublicKey string
	// the PublicKey used to identify you on the platform. We could still reference
	// people by name, but we needn't since we have the pk anyway.
	Seed string
	// optional, this is if the user created his account on our website
	// should be shown once and deleted permanently
	// add a notice like "WE DO NOT SAVE YOUR SEED" on the UI side
	AmountInvested float64
	// total amount, would be nice to track to contact them,
	// give them some kind of medals or something
	FirstSignedUp string
	// auto generated timestamp
	InvestedAssets []Order
	// array of asset codes this user has invested in
	// also I think we need a username + password for logging on to the platform itself
	// linking it here for now
	LoginUserName string
	// the thing you use to login to the platform
	LoginPassword string
	// LoginPassword is different from the seed you get if you choose
	// to open your account on the website. This is becasue even if you lose the
	// login password, you needn't worry too much about losing your funds, sicne you have
	// your seed and can send them to another address immediately.
}
