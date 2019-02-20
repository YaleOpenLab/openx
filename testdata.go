package main

import (
	"errors"
	"log"

	database "github.com/YaleOpenLab/openx/database"
	solar "github.com/YaleOpenLab/openx/platforms/opensolar"
	opzones "github.com/YaleOpenLab/openx/platforms/ozones"
	utils "github.com/YaleOpenLab/openx/utils"
)

func InsertDummyData() error {
	var err error
	// populate database with dumym data
	var project1 solar.Project
	var project2 solar.Project
	var project3 solar.Project
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

	_, err = opzones.NewConstructionBond("Dec 21 2021", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Opportunity Zone Construction", 200, "5% tax for 10 years", 1, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		log.Fatal(err)
	}

	_, err = opzones.NewConstructionBond("Apr 2 2025", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Opportunity Zone Construction", 400, "No tax for 20 years", 1, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		log.Fatal(err)
	}

	_, err = opzones.NewConstructionBond("Jul 9 2029", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Opportunity Zone Construction", 100, "3% Tax for 5 Years", 1, "Osaka Development Project", "Osaka", "This Project seeks to develop cutting edge technologies in Osaka")
	if err != nil {
		log.Fatal(err)
	}
	// newParams(mdate string, mrights string, stype string, intrate float64, rating string, bIssuer string, uWriter string
	// unitCost float64, itype string, nUnits int, tax string
	coop, err := opzones.NewLivingUnitCoop("Dec 21 2021", "Member Rights Link", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Coop Model", 4000, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		log.Fatal(err)
	}

	_, err = opzones.NewLivingUnitCoop("Apr 2 2025", "Member Rights Link", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Coop Model", 1000, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		log.Fatal(err)
	}

	_, err = opzones.NewLivingUnitCoop("Jul 9 2029", "Member Rights Link", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Coop Model", 2000, "Osaka Development Project", "Osaka", "ODP seeks to develop cutting edge technologies in Osaka and invites investors all around the world to be a part of this new age")
	if err != nil {
		log.Fatal(err)
	}
	_, err = opzones.RetrieveLivingUnitCoop(coop.Index)
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
		log.Fatal(err)
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
	project1.RecipientIndex = rec.U.Index
	project1.Contractor = c1
	project1.Originator = newOriginator
	project1.Stage = 3
	project1.PaybackPeriod = 2
	project1.AuctionType = "blind"
	err = project1.Save()
	if err != nil {
		return errors.New("Error inserting project into db")
	}

	project2.Index = 2
	project2.PanelSize = "180 1200 sq.ft homes in a high rise building 0.1mi from Kendall Square"
	project2.TotalValue = 30000
	project2.Location = "Kendall Square, Boston"
	project2.MoneyRaised = 0
	project2.Metadata = "Kendall Square is set in the heart of Cambridge and is a popular startup IT hub"
	project2.InvestorAssetCode = ""
	project2.DebtAssetCode = ""
	project2.PaybackAssetCode = ""
	project2.DateInitiated = utils.Timestamp()
	project2.Years = 5
	project2.RecipientIndex = rec.U.Index
	project2.Contractor = c1
	project2.Originator = newOriginator
	project2.Stage = 3
	project2.PaybackPeriod = 2
	project2.AuctionType = "blind"
	err = project2.Save()
	if err != nil {
		return errors.New("Error inserting project into db")
	}

	project3.Index = 3
	project3.PanelSize = "260 1500 sq.ft homes set in a medieval cathedral style construction"
	project3.TotalValue = 40000
	project3.Location = "Trafalgar Square, London"
	project3.MoneyRaised = 0
	project3.Metadata = "Trafalgar Square is set in the heart of London's financial district, with big banks all over"
	project3.InvestorAssetCode = ""
	project3.DebtAssetCode = ""
	project3.PaybackAssetCode = ""
	project3.DateInitiated = utils.Timestamp()
	project3.Years = 7
	project3.RecipientIndex = rec.U.Index
	project3.Contractor = c1
	project3.Originator = newOriginator
	project3.Stage = 3
	project3.PaybackPeriod = 2
	project3.AuctionType = "blind"
	err = project3.Save()
	if err != nil {
		return errors.New("Error inserting project into db")
	}

	pc, err := newOriginator.Originate("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", 1, "blind") // 1 is the idnex for martin
	if err != nil {
		log.Fatal(err)
	}

	_, err = solar.RetrieveProject(pc.Index)
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

	demoProject.Index = len(indexHelp) + 1
	demoProject.PanelSize = "10x 100W Komaes Solar Panels"
	demoProject.TotalValue = 8000 + 2000
	demoProject.Location = "S.U. Pasto School, Puerto Rico"
	demoProject.MoneyRaised = 10000
	demoProject.Years = 5
	demoProject.PaybackPeriod = 2
	demoProject.InterestRate = 0.029
	demoProject.Metadata = "This is a pilot initiative of the MIT-Yale effort to integrate solar platforms with IoT data and blockchain based payment systems to help develop community shelters in Puerto Rico"
	demoProject.Inverter = "Schneider Conext SW 230V 2024"
	demoProject.ChargeRegulator = "Schneider MPPT60"
	demoProject.ControlPanel = "Schneider XW SCP"
	demoProject.CommBox = "Schneider Conext Insight"
	demoProject.ACTransfer = "Eaton Manual throw switches between grid and solar+grid setups"
	demoProject.SolarCombiner = "MidNite"
	demoProject.Batteries = "Advance Autoparts Deep cycle 600A"
	demoProject.IoTHub = "Raspberry Pi 3"
	demoProject.DateInitiated = "01/23/2018"
	demoProject.DateFunded = "06/19/2018"
	demoProject.BalLeft = 10000 // assume recipient has not paid anything back to us yet

	demoProject.Originator = demoOrig
	demoProject.Contractor = demoCont
	demoProject.Developer = demoDevel
	demoProject.Guarantor = demoGuar
	demoProject.ContractorFee = 2000
	demoProject.DeveloperFee = 6000
	demoProject.RecipientIndex = demoRec.U.Index
	demoProject.Stage = 6
	demoProject.AuctionType = "private"
	demoProject.SpecSheetHash = "ipfshash" // TODO: replace this with the real ipfs hash for the demo
	demoProject.Reputation = 10000         // fix this equal to total value
	demoProject.InvestorIndices = append(demoProject.InvestorIndices, demoInv.U.Index)
	demoProject.InvestmentType = "Municipal Bond"

	err = demoProject.Save()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
