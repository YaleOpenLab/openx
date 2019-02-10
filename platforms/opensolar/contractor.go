package opensolar

import (
	"fmt"

	database "github.com/YaleOpenLab/openx/database"
	utils "github.com/YaleOpenLab/openx/utils"
)

// When contractors are proposing a contract towards something,
// we need to be sure that they are not following the price and are instead giving
// their best quote possible. In this case, a blind auction method is the best
// and that's what we have right now. If we want this to be an auction as well, we
// need to have a specific date of sorts where all the contractors can propose
// contracts immmediately, without latency.
// A given Contractor right now is allowed only for one final bid for blind
// auction advantages (i.e. no price discovery, etc). If we want to change this, we must
// Also, have some kind of deposit for Contractors (5% or something) so that they
// don't go back on their investment and slash their investment by 10% if this happens
// and distribute that amount to the recipient directly and reduce everyone's bids
// by that amount to account for the change in underlying Project

// Contractors are created here inheriting properties from Users/Entities
func NewContractor(uname string, pwd string, seedpwd string, Name string, Address string, Description string) (Entity, error) {
	// Create a new entity with the boolean of 'contractor' set to 'true.' This is
	// done just by passing the string "contractor"
	// TODO: Consider other specific information needed for contractors, other than
	// the ones set for users and entities. It can go here, or set as a separate struct.
	return newEntity(uname, pwd, seedpwd, Name, Address, Description, "contractor")
}

func (contractor *Entity) Propose(panelSize string, totalValue float64, location string,
	years int, metadata string, recIndex int, projectIndex int, auctionType string) (Project, error) {
	var pc Project
	var err error

	// for this, create a new contract and store in the contracts db. We are sorting
	// by stage, so it shouldn't matter a whole lot
	indexCheck, err := RetrieveAllProjects()
	if err != nil {
		return pc, fmt.Errorf("Projects could not be retrieved!")
	}
	pc.Params.Index = len(indexCheck) + 1
	pc.Params.PanelSize = panelSize
	pc.Params.TotalValue = totalValue
	pc.Params.Location = location
	pc.Params.Years = years
	pc.Params.Metadata = metadata
	pc.Params.DateInitiated = utils.Timestamp()
	iRecipient, err := database.RetrieveRecipient(recIndex)
	if err != nil {
		return pc, err
	}
	pc.ProjectRecipient = iRecipient
	pc.Stage = 2 // 2 since we need to filter this out while retrieving the propsoed contracts
	pc.AuctionType = auctionType
	pc.Contractor = *contractor
	// instead of storing in this proposedcontracts slice, store it as a project, but not a contract and retrieve by stage
	err = pc.Save()
	// don't insert the project since the contractor's projects are not final
	return pc, err
}

func (contractor *Entity) AddCollateral(amount float64, data string) error {
	// sets the collaterla ofr a particular contractor. Since this is not required
	// by other entities, this is declared here
	contractor.Collateral += amount
	contractor.CollateralData = append(contractor.CollateralData, data)
	return contractor.Save()
}

func (contractor *Entity) Slash(contractValue float64) error {
	// slash an entity's reputation score if it reneges back on an agreed contract
	contractor.U.Reputation -= contractValue * 0.1
	return contractor.Save()
}

func RepInstalledProject(contrIndex int, projIndex int) error {
	contractor, err := RetrieveEntity(contrIndex)
	if err != nil {
		return err
	}
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}
	err = project.SetInstalledProjectStage()
	if err != nil {
		return err
	}
	// by default, add reputation to the contractor for installed projects. If the
	// recipient disputes this claim, then review and reduce the contractor's reputation
	contractor.U.Reputation += project.Params.TotalValue * ContractorWeight
	return contractor.Save()
}
