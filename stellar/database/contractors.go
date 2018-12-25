package database

import (
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
)

func NewContractor(uname string, pwd string, Name string, Address string, Description string) (ContractEntity, error) {
	newContractor, err := NewContractEntity(uname, pwd, Name, Address, Description, "contractor")
	if err != nil {
		return newContractor, err
	}

	// insert the contractor into the database
	err = InsertContractEntity(newContractor)
	if err != nil {
		return newContractor, err
	}

	return newContractor, err
}

func (contractor *ContractEntity) ProposeContract(panelSize string, totalValue int, location string, years int, metadata string, recipient Recipient, orderIndex uint32) (Contract, error) {
	var pc Contract
	var err error

	pc.O.Index = orderIndex
	pc.O.PanelSize = panelSize
	pc.O.TotalValue = totalValue
	pc.O.Location = location
	pc.O.Years = years
	pc.O.Metadata = metadata
	pc.O.OrderRecipient = recipient
	pc.O.DateInitiated = utils.Timestamp()
	contractor.ProposedContracts = append(contractor.ProposedContracts, pc)

	err = InsertContractEntity(*contractor)
	if err != nil {
		return pc, err
	}

	// don't insert the order since the contractor's orders are not final
	return pc, err
}
