package database

import (
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
)

func NewOriginator(uname string, pwd string, Name string, Address string, Description string) (ContractEntity, error) {
	newOriginator, err := NewContractEntity(uname, pwd, Name, Address, Description, "originator")
	if err != nil {
		return newOriginator, err
	}

	// insert the originator into the database
	err = InsertContractEntity(newOriginator)
	if err != nil {
		return newOriginator, err
	}

	return newOriginator, err
}

func (originator *ContractEntity) OriginContract(panelSize string, totalValue int, location string, years int, metadata string, recIndex int) (Contract, error) {
	//log.Println("NEW ORING: ", newOriginator)
	var pc Contract
	var err error
	// PanelSize, TotalValue, Location, Years, Metadata
	// skip Recipient now
	allOrders, err := RetrieveAllOrders()
	if err != nil {
		return pc, err
	}

	if len(allOrders) == 0 {
		pc.O.Index = 1
	} else {
		pc.O.Index = len(allOrders) + 1
	}

	pc.O.PanelSize = panelSize
	pc.O.TotalValue = totalValue
	pc.O.Location = location
	pc.O.Years = years
	pc.O.Metadata = metadata
	pc.O.DateInitiated = utils.Timestamp()
	pc.O.Origin = true
	oRecipient, err := RetrieveRecipient(recIndex)
	if err != nil {
		return pc, err
	}
	pc.O.OrderRecipient = oRecipient
	pc.O.Stage = 0
	originator.ProposedContracts = append(originator.ProposedContracts, pc)
	// need to update the database
	err = InsertContractEntity(*originator)
	if err != nil {
		return pc, err
	}

	// each origin contract is assumed to go directly into the orders which can be
	// tkaen up by contractors.
	// TODO: multiple originators
	// insert this into the db since this is the order the recipient needs to  choose
	// the best option for
	err = InsertOrder(pc.O) // assume this originated order is final
	if err != nil {
		return pc, err
	}

	return pc, err
}
