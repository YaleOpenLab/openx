package database

import (
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
)

// When contractors are proposing a contract towards something,
// we need to be sure that they are not following the price and are instead giving
// their best quote possible. In this case, a blind auction method is the best
// and that's what we have right now. If we want this to be an auction as well, we
// need to have a specific date of sorts where all the contractors can propose
// contracts immmediately, without latency.
// Also, have some kind of deposit for Contractors (5% or something) so that they
// don't go back on their investment and slash their ivnestment by 10% if this happens
// and distribute that amount to the recipient directly and reduce everyone's bids
// by that amount to account for the change in underlying Project
// also, a given Contractor right now is allowed only for one final bid for blind
// auction advantages (no price disvocery, etc). If we want to change this, we must
// have an auction handler that will take care of this.

func NewContractor(uname string, pwd string, Name string, Address string, Description string) (Entity, error) {
	newContractor, err := NewEntity(uname, pwd, Name, Address, Description, "contractor")
	if err != nil {
		return newContractor, err
	}

	// insert the contractor into the database
	err = InsertEntity(newContractor)
	if err != nil {
		return newContractor, err
	}

	return newContractor, err
}

func (contractor *Entity) ProposeContract(panelSize string, totalValue int, location string, years int, metadata string, recIndex int, projectIndex int) (Project, error) {
	var pc Project
	var err error

	// for this, create a new contract and store in the contracts db. Wea re sorting
	// by stage, so it shouldn't matter a whole lot
	indexCheck, err := RetrieveAllProjects()
	if err != nil {
		return pc, err
	}
	pc.Params.Index = len(indexCheck) + 1
	pc.Params.PanelSize = panelSize
	pc.Params.TotalValue = totalValue
	pc.Params.Location = location
	pc.Params.Years = years
	pc.Params.Metadata = metadata
	pc.Params.DateInitiated = utils.Timestamp()
	iRecipient, err := RetrieveRecipient(recIndex)
	if err != nil {
		return pc, err
	}
	pc.Params.ProjectRecipient = iRecipient
	pc.Stage = 2 // 2 since we need to filter this out while retrieving the propsoed contracts
	pc.Contractor = *contractor
	// instead of storing in this proposedcontracts slice, store it as a project, but not a contract and retrieve by stage
	err = pc.Save()
	// don't insert the project since the contractor's projects are not final
	return pc, err
}
