package main

import (
	"fmt"
	"log"

	database "github.com/YaleOpenLab/openx/database"
	bonds "github.com/YaleOpenLab/openx/platforms/bonds"
	solar "github.com/YaleOpenLab/openx/platforms/solar"
	utils "github.com/YaleOpenLab/openx/utils"
)

func InsertDummyData() error {
	var err error
	// populate database with dumym data
	var project1 solar.SolarParams
	var contract1 solar.Project
	var contract2 solar.Project
	var contract3 solar.Project
	var rec database.Recipient
	allRecs, err := database.RetrieveAllRecipients()
	if err != nil {
		log.Fatal(err)
	}
	if len(allRecs) == 0 {
		// there is no recipient right now, so create a dummy recipient
		var err error
		rec, err = database.NewRecipient("martin", "p", "x", "Martin")
		if err != nil {
			log.Fatal(err)
		}
		rec.U.Notification = true
		err = rec.AddEmail("varunramganesh@gmail.com")
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
		inv, err = database.NewInvestor("john", "p", "x", "John")
		if err != nil {
			log.Fatal(err)
		}
		err = inv.AddVotingBalance(100000)
		// this function saves as well, so there's no need to save again
		if err != nil {
			log.Fatal(err)
		}
		err = database.AddInspector(inv.U.Index)
		if err != nil {
			log.Fatal(err)
		}
		x, err := database.RetrieveUser(inv.U.Index)
		if err != nil {
			log.Fatal(err)
		}
		inv.U = x
		err = inv.Save()
		if err != nil {
			log.Fatal(err)
		}
		err = x.Authorize(inv.U.Index)
		if err != nil {
			log.Fatal(err)
		}
		inv.U.Notification = true
		err = inv.AddEmail("varunramganesh@gmail.com")
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = bonds.NewBond("Dec 21 2021", "Member Rights Link", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Opportunity Zone Construction", 200, "5% tax for 10 years", 1, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		log.Fatal(err)
	}

	_, err = bonds.NewBond("Apr 2 2025", "Member Rights Link", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Opportunity Zone Construction", 400, "No tax for 20 years", 1, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		log.Fatal(err)
	}

	_, err = bonds.NewBond("Jul 9 2029", "Member Rights Link", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Opportunity Zone Construction", 100, "3% Tax for 5 Years", 1, "Osaka Development Project", "Osaka", "This Project seeks to develop cutting edge technologies in Osaka")
	if err != nil {
		log.Fatal(err)
	}
	// newParams(mdate string, mrights string, stype string, intrate float64, rating string, bIssuer string, uWriter string
	// unitCost float64, itype string, nUnits int, tax string
	coop, err := bonds.NewCoop("Dec 21 2021", "Member Rights Link", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Coop Model", 4000, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		log.Fatal(err)
	}

	_, err = bonds.NewCoop("Apr 2 2025", "Member Rights Link", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Coop Model", 1000, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		log.Fatal(err)
	}

	_, err = bonds.NewCoop("Jul 9 2029", "Member Rights Link", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Coop Model", 2000, "Osaka Development Project", "Osaka", "ODP seeks to develop cutting edge technologies in Osaka and invites investors all around the world to be a part of this new age")
	if err != nil {
		log.Fatal(err)
	}
	_, err = bonds.RetrieveCoop(coop.Params.Index)
	if err != nil {
		log.Fatal(err)
	}
	// NewOriginator(uname string, pwd string, Name string, Address string, Description string)
	newOriginator, err := solar.NewOriginator("john", "p", "x", "John Doe", "14 ABC Street London", "This is a sample originator")
	if err != nil {
		log.Fatal(err)
	}

	c1, err := solar.NewContractor("john", "p", "x", "John Doe", "14 ABC Street London", "This is a sample contractor")
	if err != nil {
		log.Println(err)
	}

	project1.Index = 1
	project1.PanelSize = "100 1000 sq.ft homes each with their own private spaces for luxury"
	project1.TotalValue = 14000
	project1.Location = "India Basin, San Francisco"
	project1.MoneyRaised = 0
	project1.Metadata = "India Basin is an upcoming creative project based in San Francisco that seeks to invite innovators from all around to participate"
	project1.InvestorAssetCode = ""
	project1.DebtAssetCode = ""
	project1.PaybackAssetCode = ""
	project1.DateInitiated = utils.Timestamp()
	project1.Years = 3
	contract1.Params = project1
	contract1.ProjectRecipient = rec
	contract1.Contractor = c1
	contract1.Originator = newOriginator
	contract1.Stage = 3
	contract1.AuctionType = "blind"
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
	project1.InvestorAssetCode = ""
	project1.DebtAssetCode = ""
	project1.PaybackAssetCode = ""
	project1.DateInitiated = utils.Timestamp()
	project1.Years = 5
	contract2.ProjectRecipient = rec
	contract2.Params = project1
	contract2.Contractor = c1
	contract2.Originator = newOriginator
	contract2.Stage = 3
	contract2.AuctionType = "blind"
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
	project1.InvestorAssetCode = ""
	project1.DebtAssetCode = ""
	project1.PaybackAssetCode = ""
	project1.DateInitiated = utils.Timestamp()
	project1.Years = 7
	contract3.ProjectRecipient = rec
	contract3.Params = project1
	contract3.Contractor = c1
	contract3.Originator = newOriginator
	contract3.Stage = 3
	contract3.AuctionType = "blind"
	err = contract3.Save()
	if err != nil {
		return fmt.Errorf("Error inserting project into db")
	}

	pc, err := newOriginator.Originate("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", 1, "blind") // 1 is the idnex for martin
	if err != nil {
		log.Fatal(err)
	}

	_, err = solar.RetrieveProject(pc.Params.Index)
	if err != nil {
		log.Fatal(err)
	}

	// Each contractor building off of this must reference the project index in their
	// proposed contract to enable searchability of the bucket. And each contractor
	// must build off of this in their proposed Contracts
	// Contractor stuff below, competing contractor details follow
	_, err = solar.NewContractor("sam", "p", "x", "Samuel Jackson", "14 ABC Street London", "This is a competing contractor")
	if err != nil {
		log.Fatal(err)
	}

	_, err = solar.NewOriginator("samuel", "p", "x", "Samuel L. Jackson", "ABC Street, London", "I am an originator")
	if err != nil {
		log.Fatal(err)
	}

	_, err = solar.RetrieveAllEntities("originator")
	if err != nil {
		log.Fatal(err)
	}

	_, err = solar.RetrieveAllEntities("contractor")
	if err != nil {
		log.Fatal(err)
	}

	// MWTODO: get comments on various fileds in this file
	demoInv, err := database.NewInvestor("Yale OpenLab", "p", "x", "Yale OpenLab")
	if err != nil {
		log.Fatal(err)
	}

	demoRec, err := database.NewRecipient("S.U. Pasto School, Puerto Rico", "p", "x", "S.U. Pasto School")
	if err != nil {
		log.Fatal(err)
	}

	demoOrig, err := solar.NewOriginator("MIT Digital Curreny Initiative", "p", "x", "MIT DCI", "MIT Building E14-15", "The MIT Media Lab's Digital Currency Initiative")
	if err != nil {
		log.Fatal(err)
	}

	demoCont, err := solar.NewContractor("Martin Wainstein", "p", "x", "Martin Wainstein", "254 Elm Street, New Haven, CT", "Martin Wainstein from the Yale OpenLab")
	if err != nil {
		log.Fatal(err)
	}

	demoDevel, err := solar.NewDeveloper("Genmoji Solar", "p", "x", "Genmoji Solar", "Genmoji, San Juan, Puerto Rico", "Genmoji Solar")
	if err != nil {
		log.Fatal(err)
	}

	demoGuar, err := solar.NewGuarantor("MIT Media Lab", "p", "x", "MIT Media Lab", "MIT Building E14-15", "The MIT Media Lab is an interdisciplinary lab with innovators from all around the globe")
	if err != nil {
		log.Fatal(err)
	}

	var demoProject solar.Project

	indexHelp, err := solar.RetrieveAllProjects()
	if err != nil {
		log.Fatal(err)
	}

	demoProject.Params.Index = len(indexHelp) + 1
	demoProject.Params.PanelSize = "10x 100W Komaes Solar Panels"
	demoProject.Params.TotalValue = 8000 + 2000
	demoProject.Params.Location = "S.U. Pasto School, Puerto Rico"
	demoProject.Params.MoneyRaised = 10000
	demoProject.Params.Years = 5
	demoProject.Params.InterestRate = 0.029
	demoProject.Params.Metadata = "This is a pilot initiative of the MIT-Yale effort to integrate solar platforms with IoT data and blockchain based payment systems to help develop community shelters in Puerto Rico"
	demoProject.Params.Inverter = "Schneider Conext SW 230V 2024"
	demoProject.Params.ChargeRegulator = "Schneider MPPT60"
	demoProject.Params.ControlPanel = "Schneider XW SCP"
	demoProject.Params.CommBox = "Schneider Conext Insight"
	demoProject.Params.ACTransfer = "Eaton Manual throw switches between grid and solar+grid setups"
	demoProject.Params.SolarCombiner = "MidNite"
	demoProject.Params.Batteries = "Advance Autoparts Deep cycle 600A"
	demoProject.Params.IoTHub = "Raspberry Pi 3"
	demoProject.Params.DateInitiated = "01/23/2018"
	demoProject.Params.DateFunded = "06/19/2018"
	demoProject.Params.BalLeft = 10000 // assume recipient has not paid anything back to us yet

	demoProject.Originator = demoOrig
	demoProject.Contractor = demoCont
	demoProject.Developer = demoDevel
	demoProject.Guarantor = demoGuar
	demoProject.ContractorFee = 2000
	demoProject.DeveloperFee = 6000
	demoProject.ProjectRecipient = demoRec
	demoProject.Stage = 6
	demoProject.AuctionType = "private"
	demoProject.SpecSheetHash = "ipfshash" // TODO: replace this with the real ipfs hash for the demo
	demoProject.Reputation = 10000         // fix this equal to total value
	demoProject.ProjectInvestors = append(demoProject.ProjectInvestors, demoInv)
	demoProject.InvestmentType = "Municipal Bond"

	err = demoProject.Save()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
