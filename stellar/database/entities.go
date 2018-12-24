package database

// entities defines a list of all the entities involved in the system like
// investor, recipient, order, etc
// Order is what is advertised on the frontend where investors can choose what
// projects to invest in.
// the idea of an Order should be constructed by ourselves based on some preset
// parameters. It should hold a contractor, developers, guarantor, originator
// along with the fields already defined. The other fields should be displayed
// on the frontend.
// Order is the structure that is actually stored in the database
// Order is a subset of the larger struct Contractor, which is still a TODO
type Order struct {
	Index         uint32  // an Index to keep quick track of how many orders exist
	PanelSize     string  // size of the given panel, for diplsaying to the user who wants to bid stuff
	TotalValue    int     // the total money that we need from investors
	Location      string  // where this specific solar panel is located
	MoneyRaised   int     // total money that has been raised until now
	Years         int     // number of years the recipient has chosen to opt for
	Metadata      string  // any other metadata can be stored here
	Live          bool    // check to see whether the current order is live or not
	PaidOff       bool    // whether the asset has been paidoff by the recipient
	INVAssetCode  string  // once all funds have been raised, we need to set assetCodes
	DEBAssetCode  string  // once all funds have been raised, we need to set assetCodes
	PBAssetCode   string  // once all funds have been raised, we need to set assetCodes
	BalLeft       float64 // denotes the balance left to pay by the party
	DateInitiated string  // date the order was created
	DateLastPaid  string  // date the order was last paid
	// instead of just holding the recipient's name here, we can hold the recipient
	OrderRecipient Recipient
	// also have an array of investors to keep track of who invested in these projects
	OrderInvestors []Investor
	// TODO: have an investor and recipient relation here
	// Percentage raised is not stored in the database since that can be calculated by the UI
}

// the user structure  houses all entities that are of type "User". This contains
// commonly used functions so that we need not repeat the ssame thing for every instance.
type User struct {
	Index uint32
	// default index, gets us easy stats on how many people are there and stuff,
	// don't want to omit this
	Seed string
	// Seed is the equivalent of a private key in stellar (stellar doesn't expose private keys)
	// do we make seed optional like that for the Recipient? Couple things to consider
	// here: if the recipient loses the publickey, it can nver send DEBTokens back
	// to the issuer, so it would be as if it reneged on the deal. Do we count on
	// technically less sound people to hold their public keys safely? I suggest
	// this would be  difficult in practice, so maybe enforce that they need to hold|
	// their account on the platform?
	Name string
	// Name of the primary stakeholder involved (principal trustee of school, for eg.)
	PublicKey string
	// PublicKey denotes the public key of the recipient
	LoginUserName string
	// the username you use to login to the platform
	LoginPassword string
	// the password, which you use to authenticate on the platform
	FirstSignedUp string
	// auto generated timestamp
}

type Recipient struct {
	ReceivedOrders []Order
	// ReceivedOrders denotes the orders that have been received by the recipient
	// instead of storing the PaybackAssets and the DebtAssets, we store this
	U User
	// user related functions are called as an instance directly
	// TODO: better name? idk
}

// the investor struct contains all the investor details such as
// public key, seed (if account is created on the website) and ot her stuff which
// is yet to be decided

// All investors will be referenced by their public key, name is optional (maybe necessary?)
// we need to stil ldecide on identity and stuff and how much we want to track
// people who invest in the schools
type Investor struct {
	AmountInvested float64
	// total amount, would be nice to track to contact them,
	// give them some kind of medals or something
	InvestedAssets []Order
	// array of asset codes this user has invested in
	// also I think we need a username + password for logging on to the platform itself
	// linking it here for now
	U User
	// user related functions are called as an instance directly
}

// the contractor super struct comprises of various entities within it. Its a
// super class because combining them results in less duplication of code
// TODO: in some ways, the Name, LoginUserName and LoginPassword fields can be
// devolved into a separate User struct, that would result in less duplication as
// well
type Contractor struct {
	// User defines common params such as name, seed, publickey
	U User
	// the name of the contractor / company that is contracting
	Address string
	// the registered address of the above company
	Description string
	// Does the contractor need to have a seed and a publickey?
	// we assume that it does in this case and proceed.
	// information on company credentials, their experience
	Image string
	// image can be company logo, founder selfie
	// hash of the password in reality
	IsContractor bool
	// A contractor is party who proposes a specific some of money towards a
	// particular project. This is the actual amount that the investors invest in.
	// This ideally must include the developer fee within it, so that investors
	// don't have to invest in two things. It would also make sense because the contractors
	// sometimes would hire developers themselves.
	IsGuarantor bool
	// A guarantor is somebody who can assure investors that the school will get paid
	// on time. This authority should be trusted and either should be vetted by the law
	// or have a multisig paying out to the investors beyond a certain timeline if they
	// don't get paid by the school. This way, the guarantor can be anonymous, like the
	// nice Pineapple Fund guy. THis can also be an insurance company, who is willing to
	// guarantee for specific school and the school can pay him out of chain / have
	// that as fee within the contract the originator
	IsDeveloper bool
	// A developer is someone who installs the required equipment (Raspberry Pi,
	// network adapters, anti tamper installations and similar) In the initial
	// orders, this will be us, since we'd be installign the pi ourselves, but in
	// the future, we expect third party developers / companies to do this for us
	// and act in a decentralized fashion. This money can either be paid out of chain
	// in fiat or can be a portion of the funds the investors chooses to invest in.
	// a contractor may also employ developers by himself, so this entity is not
	// strictly necessary.
	IsOriginator bool
	// An Originator is an entity that will start a project and get a fixed fee for
	// rendering its service. An Originator's role is not restricted, the originator
	// can also be the developer, contractor or guarantor. The originator should take
	// the responsibility of auditing the requirements of the project - panel size,
	// location, number of panels needed, etc. He then should ideally be able to fill
	// out some kind of form on the website so that the originator's proposal is live
	// and shown to potential investors. The originators get paid only when the order
	// is live, else they can just spam, without any actual investment
	PastContracts []Contract
	// list of all the contracts that the contractor has won in the past
	PresentContracts []Contract
	// list of all contracts that the contractor is presently undertaking1
	PastFeedback []Feedback
	// feedback received on the contractor from parites involved in the past
	// What kind of proof do we want from the company? KYC?
	// maybe we could have a photo op like exchanges do these days, with the owner
	// holding up his drivers' license or similar
}

// how does a contract evolve into an order? or do we make contracts orders?
// but we want people to be able to bid on contracts, so is it better having both
// as a single entity? ask during call and confirm so that we can do stuff. Maybe
// the "Order" struct that we use now can be a child struct of the Contract struct

type Feedback struct {
	Content string
	// the content of the feedback, good / bad
	// maybe we could have a  rating system baked in? a star based rating system?
	// would be nice, idk
	From Contractor
	// who gave the feedback?
	To Contractor
	// regarding whom is this feedback about
	Date string
	// time at which this feedback was written
	RelatedContract []Contract
	// the contract regarding which this feedback is directed at
}

// the proposal part desribed above is a collection of Contracts from different persons.
// we aren't defining the Contract part for each entity because that would casue repetition.
// Instead, a single Contract struct can be used as an engineering proposal, as a quote,
// etc.
// A contract is a superset of an order and is used to display to people what is needed.
// it will have an order class inside it, similar to how the User class exsits inside
// the investor and recipient classes.
type Contract struct {
	// a contract belongs to a Contractor, so there is no need for a reverse mapping
	// from the Contract to the Contractor
	// TODO: What stuff goes in here?
	// since there is no state trie or similar in stellar, we need to hash the contract
	// parameters and then reference it. This will be immutable and can't be changed
	// this could also be a role taken up by the legal oracle. Also need to ask
	// neighbourly and swytch regarding this. Maybe we need a bridge from stellar to
	// ethereum to interface with the ERC721 or maybe we coul have an oracle that
	// does this for us.
	O Order
}
