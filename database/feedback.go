package database

// entities defines a list of all the entities involved in the system like
// investor, recipient, order, etc
// how does a contract evolve into an order? or do we make contracts orders?
// but we want people to be able to bid on contracts, so is it better having both
// as a single entity? ask during call and confirm so that we can do stuff. Maybe
// the "project" struct that we use now can be a child struct of the Project struct

type Feedback struct {
	Content string
	// the content of the feedback, good / bad
	// maybe we could have a  rating system baked in? a star based rating system?
	// would be nice, idk
	From Entity
	// who gave the feedback?
	To Entity
	// regarding whom is this feedback about
	Date string
	// time at which this feedback was written
	RelatedContract []Project
	// the contract regarding which this feedback is directed at
}
