package opensolar

// Feedback defines a structure that can be used for providing feedback about entities
type Feedback struct {
	Content string
	// the content of the feedback, good / bad
	// maybe we could have a rating system baked in? a star based rating system?
	// would be nice, idk
	From Entity
	// who gave the feedback?
	To Entity
	// regarding whom is this feedback about
	Date string
	// time at which this feedback was written
	Contract []Project
	// the contract regarding which this feedback is directed at
}
