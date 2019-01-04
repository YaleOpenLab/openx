package main

import (
	"fmt"
	"log"

	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
)

func InsertDummyData() error {
	var err error
	// populate database with dumym data
	var project1 database.DBParams
	var contract1 database.Project
	var contract2 database.Project
	var contract3 database.Project

	var rec database.Recipient
	allRecs, err := database.RetrieveAllRecipients()
	if err != nil {
		log.Fatal(err)
	}
	if len(allRecs) == 0 {
		// there is no recipient right now, so create a dummy recipient
		var err error
		rec, err = database.NewRecipient("martin", "p", "Martin")
		if err != nil {
			log.Fatal(err)
		}
	}

	var inv database.Investor
	allInvs, err := database.RetrieveAllInvestors()
	if err != nil {
		log.Fatal(err)
	}
	if len(allInvs) == 0 {
		var err error
		inv, err = database.NewInvestor("john", "p", "John")
		if err != nil {
			log.Fatal(err)
		}
		err = inv.AddVotingBalance(100000)
		// this function saves as well, so there's no need to save again
		if err != nil {
			log.Fatal(err)
		}
	}
	// NewOriginator(uname string, pwd string, Name string, Address string, Description string)
	newOriginator, err := database.NewOriginator("john", "p", "John Doe", "14 ABC Street London", "This is a sample originator")
	if err != nil {
		log.Fatal(err)
	}

	c1, err := database.NewContractor("john", "p", "John Doe", "14 ABC Street London", "This is a sample contractor")
	if err != nil {
		log.Println(err)
	}

	project1.Index = 1
	project1.PanelSize = "100 1000 sq.ft homes each with their own private spaces for luxury"
	project1.TotalValue = 14000
	project1.Location = "India Basin, San Francisco"
	project1.MoneyRaised = 0
	project1.Metadata = "India Basin is an upcoming creative project based in San Francisco that seeks to invite innovators from all around to participate"
	project1.Funded = false
	project1.INVAssetCode = ""
	project1.DEBAssetCode = ""
	project1.PBAssetCode = ""
	project1.DateInitiated = utils.Timestamp()
	project1.Years = 3
	project1.ProjectRecipient = rec
	contract1.Params = project1
	contract1.Contractor = c1
	contract1.Originator = newOriginator
	contract1.Stage = 3
	err = contract1.Save()
	if err != nil {
		return fmt.Errorf("Error inserting project into db")
	}

	project1.Index = 2
	project1.PanelSize = "180 1200 sq.ft homes in a high rise building 0.1mi from Kendall Square"
	project1.TotalValue = 30000
	project1.Location = "Kendall Square, Boston"
	project1.MoneyRaised = 0
	project1.Metadata = "Kendall Square is set in the heart of Cambridge and is a popular startup IT hub"
	project1.Funded = false
	project1.INVAssetCode = ""
	project1.DEBAssetCode = ""
	project1.PBAssetCode = ""
	project1.DateInitiated = utils.Timestamp()
	project1.Years = 5
	project1.ProjectRecipient = rec
	contract2.Params = project1
	contract2.Contractor = c1
	contract2.Originator = newOriginator
	contract2.Stage = 3
	err = contract2.Save()
	if err != nil {
		return fmt.Errorf("Error inserting project into db")
	}

	project1.Index = 3
	project1.PanelSize = "260 1500 sq.ft homes set in a medieval cathedral style construction"
	project1.TotalValue = 40000
	project1.Location = "Trafalgar Square, London"
	project1.MoneyRaised = 0
	project1.Metadata = "Trafalgar Square is set in the heart of London's financial district, with big banks all over"
	project1.Funded = false
	project1.INVAssetCode = ""
	project1.DEBAssetCode = ""
	project1.PBAssetCode = ""
	project1.DateInitiated = utils.Timestamp()
	project1.Years = 7
	project1.ProjectRecipient = rec
	contract3.Params = project1
	contract3.Contractor = c1
	contract3.Originator = newOriginator
	contract3.Stage = 3
	err = contract3.Save()
	if err != nil {
		return fmt.Errorf("Error inserting project into db")
	}

	pc, err := newOriginator.OriginContract("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", 1) // 1 is the idnex for martin
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.RetrieveProject(pc.Params.Index)
	if err != nil {
		log.Fatal(err)
	}

	// Each contractor building off of this must reference the project index in their
	// proposed contract to enable searchability of the bucket. And each contractor
	// must build off of this in their proposed Contracts
	// Contractor stuff below
	/*
		_, err = contractor1.ProposeContract(pc.P.PanelSize, 28000, "Puerto Rico", 6, pc.P.Metadata+" we supply our own devs and provide insurance guarantee as well. Dual audit maintenance upto 1 year. Returns capped as per defaults", 1, biddingProject.Index) // 1 for retrieving martin as the recipient
		if err != nil {
			log.Fatal(err)
		}
	*/
	// competing contractor details follow
	_, err = database.NewContractor("sam", "p", "Samuel Jackson", "14 ABC Street London", "This is a competing contractor")
	if err != nil {
		log.Fatal(err)
	}

	/*
		_, err = contractor2.ProposeContract(pc.P.PanelSize, 35000, "Puerto Rico", 5, pc.P.Metadata+" free lifetime service, developers and insurance also provided", 1, biddingProject.Index) // 1 for retrieving martin as the recipient
		if err != nil {
			log.Fatal(err)
		}
	*/
	_, err = database.NewOriginator("samuel", "p", "Samuel L. Jackson", "ABC Street, London", "I am an originator")
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.RetrieveAllContractEntities("originator")
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.RetrieveAllContractEntities("contractor")
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
