package main

import (
	"github.com/pkg/errors"
	"log"

	database "github.com/YaleOpenLab/openx/database"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	opzones "github.com/YaleOpenLab/openx/platforms/ozones"
	utils "github.com/YaleOpenLab/openx/utils"
)

func newSolarProject(index int, panelsize string, totalValue float64, location string, moneyRaised float64,
	metadata string, invAssetCode string, debtAssetCode string, pbAssetCode string, years int, recpIndex int,
	contractor opensolar.Entity, originator opensolar.Entity, stage int, pbperiod int, auctionType string) (opensolar.Project, error) {

	var project opensolar.Project
	project.Index = index
	project.PanelSize = panelsize
	project.TotalValue = totalValue
	project.State = location
	project.MoneyRaised = moneyRaised
	project.Metadata = metadata
	project.InvestorAssetCode = invAssetCode
	project.DebtAssetCode = debtAssetCode
	project.PaybackAssetCode = pbAssetCode
	project.DateInitiated = utils.Timestamp()
	project.ETA = years
	project.RecipientIndex = recpIndex
	project.Contractor = contractor
	project.Originator = originator
	project.Stage = stage
	project.PaybackPeriod = pbperiod
	project.AuctionType = auctionType
	project.InvestmentType = "munibond"

	var x1 opensolar.TermsHelper
	x1.Variable = "Security Type"
	x1.Value = "Municipal Bond"
	x1.RelevantParty = "PR DofEd"
	x1.Note = "Promoted by PR governor's office"
	x1.Status = "Demo"
	x1.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var x2 opensolar.TermsHelper
	x2.Variable = "PPA Tariff"
	x2.Value = "0.24 ct/KWh"
	x2.RelevantParty = "oracle X / PREPA"
	x2.Note = "Variable anchored to local tariff"
	x2.Status = "Signed"
	x2.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var x3 opensolar.TermsHelper
	x3.Variable = "Return (TEY)"
	x3.Value = "3.1%"
	x3.RelevantParty = "Broker Dealer"
	x3.Note = "Variable tied to tariff"
	x3.Status = "Signed"
	x3.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var x4 opensolar.TermsHelper
	x4.Variable = "Maturity"
	x4.Value = "+/- 2025"
	x4.RelevantParty = "Broker Dealer"
	x4.Note = "Tax adjusted Yield"
	x4.Status = "Signed"
	x4.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var x5 opensolar.TermsHelper
	x5.Variable = "Guarantee"
	x5.Value = "50%"
	x5.RelevantParty = "Foundation X"
	x5.Note = "First-loss upon breach"
	x5.Status = "Started"
	x5.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var x6 opensolar.TermsHelper
	x6.Variable = "Insurance"
	x6.Value = "Premium"
	x6.RelevantParty = "Allianz CS"
	x6.Note = "Hurricane Coverage"
	x6.Status = "Started"
	x6.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	project.Terms = append(project.Terms, x1, x2, x3, x4, x5, x6)
	err := project.Save()
	if err != nil {
		return project, errors.New("Error inserting project into db")
	}
	return project, nil
}

// newLivingUnitCoop creates a new living unit coop
func newLivingUnitCoop(mdate string, mrights string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, totalAmount float64, typeOfUnit string, monthlyPayment float64,
	title string, location string, description string) (opzones.LivingUnitCoop, error) {
	var coop opzones.LivingUnitCoop
	coop.MaturationDate = mdate
	coop.MemberRights = mrights
	coop.SecurityType = stype
	coop.InterestRate = intrate
	coop.Rating = rating
	coop.BondIssuer = bIssuer
	coop.Underwriter = uWriter
	coop.Title = title
	coop.Location = location
	coop.Description = description
	coop.DateInitiated = utils.Timestamp()

	x, err := opzones.RetrieveAllLivingUnitCoops()
	if err != nil {
		return coop, errors.Wrap(err, "could not retrieve all living unit coops")
	}
	coop.Index = len(x) + 1
	coop.UnitsSold = 0
	coop.Amount = totalAmount
	coop.TypeOfUnit = typeOfUnit
	coop.MonthlyPayment = monthlyPayment
	err = coop.Save()
	return coop, err
}

// newConstructionBond returns a New Construction Bond and automatically stores it in the db
func newConstructionBond(mdate string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, unitCost float64, itype string, nUnits int, tax string, recIndex int,
	title string, location string, description string) (opzones.ConstructionBond, error) {
	var cBond opzones.ConstructionBond
	cBond.MaturationDate = mdate
	cBond.SecurityType = stype
	cBond.InterestRate = intrate
	cBond.Rating = rating
	cBond.BondIssuer = bIssuer
	cBond.Underwriter = uWriter
	cBond.Title = title
	cBond.Location = location
	cBond.Description = description
	cBond.DateInitiated = utils.Timestamp()

	x, err := opzones.RetrieveAllConstructionBonds()
	if err != nil {
		return cBond, errors.Wrap(err, "could not retrieve all living unit coops")
	}

	cBond.Index = len(x) + 1
	cBond.CostOfUnit = unitCost
	cBond.InstrumentType = itype
	cBond.NoOfUnits = nUnits
	cBond.Tax = tax
	cBond.RecipientIndex = recIndex
	err = cBond.Save()
	return cBond, err
}

// InsertDummyData inserts sample data
func InsertDummyData() error {
	var err error
	// populate database with dumym data
	var recp database.Recipient
	allRecs, err := database.RetrieveAllRecipients()
	if err != nil {
		log.Fatal(err)
	}
	if len(allRecs) == 0 {
		// there is no recipient right now, so create a dummy recipient
		var err error
		recp, err = database.NewRecipient("martin", "p", "x", "Martin")
		if err != nil {
			log.Fatal(err)
		}
		recp.U.Notification = true
		err = recp.AddEmail("varunramganesh@gmail.com")
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

	// MW: Are these users that engage with demo projects?

	originator, err := opensolar.NewOriginator("samuel", "p", "x", "Samuel L. Jackson", "ABC Street, London", "I am an originator")
	if err != nil {
		log.Fatal(err)
	}

	contractor, err := opensolar.NewContractor("sam", "p", "x", "Samuel Jackson", "14 ABC Street London", "This is a competing contractor")
	if err != nil {
		log.Fatal(err)
	}

	_, err = newConstructionBond("Dec 21 2021", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Opportunity Zone Construction", 200, "5% tax for 10 years", 1, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		log.Fatal(err)
	}

	_, err = newConstructionBond("Apr 2 2025", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Opportunity Zone Construction", 400, "No tax for 20 years", 1, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		log.Fatal(err)
	}

	_, err = newConstructionBond("Jul 9 2029", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Opportunity Zone Construction", 100, "3% Tax for 5 Years", 1, "Osaka Development Project", "Osaka", "This Project seeks to develop cutting edge technologies in Osaka")
	if err != nil {
		log.Fatal(err)
	}

	_, err = newLivingUnitCoop("Dec 21 2021", "Member Rights Link", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Coop Model", 4000, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		log.Fatal(err)
	}

	_, err = newLivingUnitCoop("Apr 2 2025", "Member Rights Link", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Coop Model", 1000, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		log.Fatal(err)
	}

	_, err = newLivingUnitCoop("Jul 9 2029", "Member Rights Link", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Coop Model", 2000, "Osaka Development Project", "Osaka", "ODP seeks to develop cutting edge technologies in Osaka and invites investors all around the world to be a part of this new age")
	if err != nil {
		log.Fatal(err)
	}

	_, err = newSolarProject(1, "100 1000 sq.ft homes each with their own private spaces for luxury", 14000, "India Basin, San Francisco",
		0, "India Basin is an upcoming creative project based in San Francisco that seeks to invite innovators from all around to participate", "", "", "",
		3, recp.U.Index, contractor, originator, 4, 2, "blind")

	if err != nil {
		log.Fatal(err)
	}

	_, err = newSolarProject(2, "180 1200 sq.ft homes in a high rise building 0.1mi from Kendall Square", 30000, "Kendall Square, Boston",
		0, "Kendall Square is set in the heart of Cambridge and is a popular startup IT hub", "", "", "",
		5, recp.U.Index, contractor, originator, 4, 2, "blind")

	if err != nil {
		log.Fatal(err)
	}

	_, err = newSolarProject(3, "260 1500 sq.ft homes set in a medieval cathedral style construction", 40000, "Trafalgar Square, London",
		0, "Trafalgar Square is set in the heart of London's financial district, with big banks all over", "", "", "",
		7, recp.U.Index, contractor, originator, 4, 2, "blind")

	if err != nil {
		log.Fatal(err)
	}

	_, err = originator.Originate("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", 1, "blind") // 1 is the idnex for martin
	if err != nil {
		log.Fatal(err)
	}

	// MWTODO: get comments on various fields in this file
	// MW: Please explain the "p" and "x" and "shortname" uses
	demoInv, err := database.NewInvestor("openlab", "p", "x", "Yale OpenLab")
	if err != nil {
		log.Fatal(err)
	}

	demoRec, err := database.NewRecipient("supasto", "p", "x", "S.U. Pasto School")
	if err != nil {
		log.Fatal(err)
	}

	demoOrig, err := opensolar.NewOriginator("dci", "p", "x", "MIT DCI", "MIT Building E14-15", "The MIT Media Lab's Digital Currency Initiative")
	if err != nil {
		log.Fatal(err)
	}

	demoCont, err := opensolar.NewContractor("mw", "p", "x", "Martin Wainstein", "254 Elm Street, New Haven, CT", "Martin Wainstein from the Yale OpenLab")
	if err != nil {
		log.Fatal(err)
	}

	demoDevel, err := opensolar.NewDeveloper("gs", "p", "x", "Genmoji Solar", "Genmoji, San Juan, Puerto Rico", "Genmoji Solar")
	if err != nil {
		log.Fatal(err)
	}

	demoGuar, err := opensolar.NewGuarantor("ml", "p", "x", "MIT Media Lab", "MIT Building E14-15", "The MIT Media Lab is an interdisciplinary lab with innovators from all around the globe")
	if err != nil {
		log.Fatal(err)
	}


	// Demo Project
	// STAGE 7 - Puerto Rico
	var demoProject opensolar.Project

	indexHelp, err := opensolar.RetrieveAllProjects()
	if err != nil {
		log.Fatal(err)
	}

	demoProject.Index = len(indexHelp) + 1
	demoProject.PanelSize = "1000W" 
	demoProject.TotalValue = 8000 + 2000
	demoProject.State = "S.U. Pasto School, Aibonito"
	demoProject.MoneyRaised = 10000
	demoProject.ETA = 5
	demoProject.PaybackPeriod = 2	// In number of weeks in which payments are triggered
	demoProject.InterestRate = 0.029
	demoProject.Metadata = "This is a pilot initiative of the MIT-Yale effort to integrate solar platforms with IoT data and blockchain based payment systems to help develop community shelters in Puerto Rico"
	demoProject.Inverter = "Schneider Conext SW 230V 2024"
	demoProject.ChargeRegulator = "Schneider MPPT60"
	demoProject.ControlPanel = "Schneider XW SCP"
	demoProject.CommBox = "Schneider Conext Insight"
	demoProject.ACTransfer = "Eaton Manual throw switches between grid and solar+grid setups"
	demoProject.SolarCombiner = "MidNite"
	demoProject.Batteries = "Advance Autoparts Deep cycle 600A"
	demoProject.IoTHub = "Yale Open Powermeter w/ RaspberryPi3"
	demoProject.DateInitiated = "01/23/2018"
	demoProject.DateFunded = "06/19/2018"
	demoProject.BalLeft = 10000 // assume recipient has not paid anything back to us yet

	// MW: Who are these originators and contractors? Can we give names to them?
	demoProject.Originator = demoOrig
	demoProject.Contractor = demoCont
	demoProject.Developer = demoDevel
	demoProject.Guarantor = demoGuar
	demoProject.ContractorFee = 2000
	demoProject.DeveloperFee = 6000
	demoProject.RecipientIndex = demoRec.U.Index
	demoProject.Stage = 7
	demoProject.AuctionType = "private"
	demoProject.StageData = append(demoProject.StageData, "ipfshash") // TODO: replace this with the real ipfs hash for the demo
	demoProject.Reputation = 10000                                    // fix this equal to total value
	demoProject.InvestorIndices = append(demoProject.InvestorIndices, demoInv.U.Index)
	demoProject.InvestmentType = "Municipal Bond"

	err = demoProject.Save()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Martin To Revise
	// omwp: One Mega Watt Project
	// STAGE 4 - New Hampshire
	var omwp opensolar.Project
	indexHelp, err = opensolar.RetrieveAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	omwp.Index = len(indexHelp) + 1

	nd1, err := opensolar.NewDeveloper("solardev", "p", "x", "First Solar", "First Solar, Earth", "First Solar")
	if err != nil {
		log.Fatal(err)
	}

	nd2, err := opensolar.NewDeveloper("Host", "p", "x", "Lancaster Town", "Host", "Host")
	if err != nil {
		log.Fatal(err)
	}

	nd3, err := opensolar.NewDeveloper("Lancaster Solar Engineer Solutions", "p", "x", "Independent RFP Engineer", "Independent RFP Engineer", "Independent RFP Engineer")
	if err != nil {
		log.Fatal(err)
	}

	nd4, err := opensolar.NewDeveloper("Simple Service Provider", "p", "x", "Simple Service Provider", "Simple Service Provider", "Simple Service Provider")
	if err != nil {
		log.Fatal(err)
	}

	nd5, err := opensolar.NewDeveloper("Vendor", "p", "x", "Vendor", "Vendor", "Vendor")
	if err != nil {
		log.Fatal(err)
	}

	nd6, err := opensolar.NewDeveloper("Auditors", "p", "x", "Auditors", "Auditors", "Auditors")
	if err != nil {
		log.Fatal(err)
	}

	nd7, err := opensolar.NewGuarantor("Insurance Agent", "p", "x", "Insurance Agent", "Insurance Agent", "Insurance Agent")
	if err != nil {
		log.Fatal(err)
	}

	nd8, err := opensolar.NewDeveloper("Utility", "p", "x", "Utility", "Utility", "Utility")
	if err != nil {
		log.Fatal(err)
	}

	omwp.DeveloperIndices = append(omwp.DeveloperIndices, nd1.U.Index, nd2.U.Index, nd3.U.Index, nd4.U.Index, nd5.U.Index, nd6.U.Index, nd7.U.Index, nd8.U.Index)
	omwp.MainDeveloper = nd1

	g1, err := opensolar.NewGuarantor("Green Bank", "p", "x", "Green Bank", "Green Bank", "Green Bank")
	if err != nil {
		log.Fatal(err)
	}
	omwp.Guarantor = g1

	i1, err := database.NewInvestor("Green Bank", "p", "x", "Green Bank")
	if err != nil {
		log.Fatal(err)
	}

	i2, err := database.NewInvestor("OZ Fund", "p", "x", "OZ Fund")
	if err != nil {
		log.Fatal(err)
	}

	i3, err := database.NewInvestor("Tax Equity Business", "p", "x", "Tax Equity Business")
	if err != nil {
		log.Fatal(err)
	}

	omwp.InvestorIndices = append(omwp.InvestorIndices, i1.U.Index, i2.U.Index, i3.U.Index)

	r1, err := database.NewRecipient("city", "p", "x", "city")
	if err != nil {
		log.Fatal(err)
	}

	r2, err := database.NewRecipient("shelter", "p", "x", "shelter")
	if err != nil {
		log.Fatal(err)
	}

	omwp.RecipientIndices = append(omwp.RecipientIndices, r1.U.Index, r2.U.Index)

	omwp.TotalValue = 2000000
	omwp.MoneyRaised = 150000
	omwp.ETA = 20
	omwp.DebtInvestor1 = "OZ Fund"
	omwp.DebtInvestor2 = "Green Bank"
	omwp.TaxEquityInvestor = "Lancaster Bank"
	omwp.State = "NH"
	omwp.Country = "USA"
	omwp.InterestRate = 0.05
	omwp.Tax = "Free for x years"
	omwp.PanelSize = "1MW"
	omwp.Batteries = "210 kWh 1x Tesla Powerpack"
	omwp.Metadata = "Neighborhood 1MW solar array on the field next to Lancaster Elementary High School. The project was originated by the head of the community organization, Ben Southworth, who is also active in the parent teacher association (PTA). The city of Lancaster has agreed to give a 20 year lease of the land to the project if the school gets to own the solar array after the lease expires. The school is located in an opportunity zone"
	omwp.BlendedCapitalInvestorIndex = i2.U.Index
	omwp.Stage = 4

	err = omwp.Save()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Martin To Revise
	// tkwp: Ten Kilowatt Project
	// STAGE 8 - Connecticut Homeless Shelter
	var tkwp opensolar.Project
	indexHelp, err = opensolar.RetrieveAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	tkwp.Index = len(indexHelp) + 1
	tkwp.TotalValue = 30000
	tkwp.Stage = 8
	tkwp.MoneyRaised = 30000
	tkwp.ETA = 7
	tkwp.State = "CT"
	tkwp.Country = "US"
	tkwp.InterestRate = 0.05
	//MW: The string doesn't like % to be included. Also Tax could be: 'TaxCredit' parameter of getting funds back, and 'TaxAmount' or 'TaxDebit' which is the percent of tax taken from the project revenue. Both can be specific % value, and also a string (eventually a drop down) describing the structure. 
	tkwp.Tax = "0.3 Tax Credit"
	tkwp.PanelSize = "15KW"
	tkwp.Metadata = "Residential solar array for a homeless shelter. The project was originated by a member of the board of the homeless shelter who gets the shelter to purchase all the electricity at a discounted rate. The shelter chooses to lease the roof for free over the lifetime of the project. The originator knows the solar developer who set up the project company"

	i1, err = database.NewInvestor("Matthew Moroney", "p", "x", "Matthew Moroney")
	if err != nil {
		log.Fatal(err)
	}

	i2, err = database.NewInvestor("Franz Hochstrasser", "p", "x", "Franz Hochstrasser")
	if err != nil {
		log.Fatal(err)
	}

	i3, err = database.NewInvestor("Connecticut Green Bank", "p", "x", "Connecticut Green Bank")
	if err != nil {
		log.Fatal(err)
	}

	i4, err := database.NewInvestor("Yale University Community Fund", "p", "x", "Yale University Community Fund")
	if err != nil {
		log.Fatal(err)
	}

	i5, err := database.NewInvestor("Jerome Green", "p", "x", "Jerome Green")
	if err != nil {
		log.Fatal(err)
	}

	i6, err := database.NewInvestor("Open Solar Revolving Fund", "p", "x", "Open Solar Revolving Fund")
	if err != nil {
		log.Fatal(err)
	}

	nd1, err = opensolar.NewDeveloper("Yale Architecture", "p", "x", "Yale Architecture", "Yale Architecture", "Yale Architecture")
	if err != nil {
		log.Fatal(err)
	}

	nd2, err = opensolar.NewDeveloper("CT Solar", "p", "x", "CT Solar", "CT Solar, CT", "CT Solar")
	if err != nil {
		log.Fatal(err)
	}

	nd3, err = opensolar.NewDeveloper("Host", "p", "x", "Columbus House", "Columbus House", "Columbus House")
	if err != nil {
		log.Fatal(err)
	}

	nd4, err = opensolar.NewGuarantor("RGreen Blend Fund", "p", "x", "RGreen Blend Fund", "RGreen Blend Fund", "RGreen Blend Fund")
	if err != nil {
		log.Fatal(err)
	}

	nd5, err = opensolar.NewDeveloper("Avangrid RECs", "p", "x", "Avangrid RECs", "Avangrid RECs", "Avangrid RECs")
	if err != nil {
		log.Fatal(err)
	}

	no1, err := opensolar.NewOriginator("Raise Green", "p", "x", "Raise Green", "Raise Green", "Raise Green")
	if err != nil {
		log.Fatal(err)
	}

	tkwp.Tax = "self sheltering"
	tkwp.InvestorIndices = append(tkwp.InvestorIndices, i1.U.Index, i2.U.Index, i3.U.Index, i4.U.Index, i5.U.Index, i6.U.Index)
	tkwp.DeveloperIndices = append(tkwp.DeveloperIndices, nd1.U.Index, nd2.U.Index, nd3.U.Index, nd4.U.Index, nd5.U.Index)
	tkwp.Originator = no1
	tkwp.InvestmentType = "Regulation Crowdfunding"

	r1, err = database.NewRecipient("Shelter1 Community Solar", "p", "x", "Shelter1 Community Solar")
	if err != nil {
		log.Fatal(err)
	}

	r2, err = database.NewRecipient("Columbus House", "p", "x", "Columbus House")
	if err != nil {
		log.Fatal(err)
	}

	tkwp.RecipientIndices = append(tkwp.RecipientIndices, r1.U.Index, r2.U.Index)

	err = tkwp.Save()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
