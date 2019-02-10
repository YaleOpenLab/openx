package opensolar

// General Concept:
// This can be applied as a feedback system to actors that are part of projects
// (i.e. relate to Entities.go, such as contractors). It allows investors and
// recipients to give comments on the services provided by these entities in the project.
// the idea is that after a stage ends, we need to ask feedback on the entity whose
// work ends at that stage. Examples follow:
// 1. Originator - the work of an originator ends at stage 1 (after the contract is
// originated), so we need to ask feedback on the originator at this stage. Based
// on this feedback, have some parameters, grade them on a scale of 1-10 and then
// have bands which reward accordingly.
// 2. Contractor - the contractor's work ends at stage 5 depeding on his involvement
// in the project, whether he spnsors developers, etc and we need to ask for feedback
// once a contract enters stage 5.
// 3. Investor - an investor's work (maybe contribution?) ends at stage 7 and we need
// to ask the recipient on his feedback of the investor - whether he was courteous,
// friendly, paid visits to the property and similar
// 4. Recipient - a recipient's work is finished after stage 7 and we can ask the other
// entities on their feedback at this stage or ask them before this stage and take a
// weighted average of their reputation scores so we may have a reputation without
// waiting for years to give feedback.
// we need to have some kind of table which dictates who can give feedback on whom and
// how much it would be worth but we can define that later.

// TODO: add additional fields here based on what feedback we collect, would depend
// on the frontend design and implementation as well
// TODO: complete this area once we have a rudimentary frontend so that we can test stuff
// more easily and link it with the reputation system
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
