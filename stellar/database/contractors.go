package database

import (
	"encoding/json"
	"fmt"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
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

func (contractor *ContractEntity) ProposeContract(panelSize string, totalValue int, location string, years int, metadata string, recipient Recipient, orderIndex int) (Contract, error) {
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

// we go through each contract entity and retrieve orders specific to the boIndex
// which is stored in their proposed contracts slice
func RetrieveAllProposedContracts(boIndex int) ([]ContractEntity, []Contract, error) {
	// boindex is the bidding order index which we should search for in all
	// contractors' proposed contracts
	var contractorsArr []ContractEntity
	var contractsArr []Contract
	temp, err := RetrieveAllUsers()
	if err != nil {
		return contractorsArr, contractsArr, err
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return contractorsArr, contractsArr, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ContractorBucket)
		for i := 1; i < limit; i++ {
			var rContractor ContractEntity
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// might be some other user like an investor or recipient
				continue
			}
			err := json.Unmarshal(x, &rContractor)
			if err != nil {
				return nil
			}
			if !rContractor.Contractor {
				continue
			}
			// is a contractor, search for the index of his proposed contracts
			contract1, err := FindInKey(boIndex, rContractor.ProposedContracts)
			if err != nil {
				// doesnt have a proposed contract for the specific recipient
				continue
			}
			// contract1 is the specific contract which has a bid towards this order
			// now we need to store the contractor and the contract for the bidding process
			contractorsArr = append(contractorsArr, rContractor)
			contractsArr = append(contractsArr, contract1)
			// default is to add all contractentities to the array
		}
		return nil
	})
	return contractorsArr, contractsArr, err
}

func FindInKey(key int, arr []Contract) (Contract, error) {
	var dummy Contract
	for _, elem := range arr {
		if elem.O.Index == key {
			return elem, nil
		}
	}
	return dummy, fmt.Errorf("Not found")
}
