package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"time"

	database "github.com/YaleOpenLab/openx/database"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
	xlm "github.com/YaleOpenLab/openx/xlm"
)

func createAllStaticEntities() error {
	// PR static entities
	var err error
	_, err = opensolar.NewOriginator("DCI", "p", "x", "MIT DCI", "MIT Building E14-15", "The MIT Media Lab's Digital Currency Initiative")
	if err != nil {
		log.Fatal(err)
	}

	_, err = opensolar.NewContractor("MartinWainstein", "p", "x", "Martin Wainstein", "254 Elm Street, New Haven, CT", "Martin Wainstein from the Yale OpenLab")
	if err != nil {
		log.Fatal(err)
	}

	_, err = opensolar.NewDeveloper("gs", "p", "x", "Genmoji Solar", "Genmoji, San Juan, Puerto Rico", "Genmoji Solar")
	if err != nil {
		log.Fatal(err)
	}

	_, err = opensolar.NewDeveloper("nbly", "p", "x", "Neighborly Securities", "San Francisco, CA", "Broker Dealer")
	if err != nil {
		log.Fatal(err)
	}

	_, err = opensolar.NewGuarantor("mitml", "p", "x", "MIT Media Lab", "MIT Building E14-15", "The MIT Media Lab is an interdisciplinary lab with innovators from all around the globe")
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func populateStaticDataPR() (int, error) {

	project, err := opensolar.RetrieveProject(5)

	project.ExplorePageSummary.Solar = project.PanelSize
	project.ExplorePageSummary.Storage = project.ExecutiveSummary.ProjectSize["Storage"]
	project.ExplorePageSummary.Tariff = project.ExecutiveSummary.Financials["Tariff (Variable)"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = project.ExecutiveSummary.Financials["Return (TEY)"]
	project.ExplorePageSummary.Rating = project.Rating
	project.ExplorePageSummary.Tax = "N/A"
	project.ExplorePageSummary.ETA = project.EstimatedAcquisition

	project.DPIntroImage = "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/1.jpg"
	project.OHeroImage = "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/3.jpg"
	project.OImages = append(project.OImages, "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/2.jpg", "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/5.jpg")
	project.AImages = append(project.AImages, "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/7.png", "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/9.png", "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/normal.png")
	project.EImages = append(project.EImages, "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/6.jpg", "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/8.png", "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/10.png", "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/11.png")
	project.CEImages = append(project.CEImages, "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/12.png", "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/13.png", "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/14.png", "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/15.png")
	project.BNImages = append(project.PSImages, "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/16.png", "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/17.png", "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/18.png", "https://images.openx.solar/OpenSolarProjects/7_YaleMIT/19.png")

	project.OriginatorIndex = 1
	project.GuarantorIndex = 5
	project.ContractorIndex = 2
	project.MainDeveloperIndex = 3
	project.DeveloperIndices = append(project.DeveloperIndices, 3, 4)
	project.ContractorFee = 2000
	project.OriginatorFee = 0
	project.DeveloperFee = append(project.DeveloperFee, 6000)

	err = project.Save()
	if err != nil {
		return -1, err
	}

	return project.Index, nil
}

func invHelper(invName, invDescription string) (database.Investor, string, error) {
	// setup investor account
	passwd := "p"
	seedpwd := "x"
	investor1, err := database.NewInvestor(invName, passwd, seedpwd, invDescription)
	if err != nil {
		return investor1, "", err
	}
	invSeed, err := wallet.DecryptSeed(investor1.U.EncryptedSeed, seedpwd)
	if err != nil {
		return investor1, "", err
	}
	err = xlm.GetXLM(investor1.U.PublicKey)
	if err != nil {
		return investor1, "", err
	}
	return investor1, invSeed, nil
}

func recpHelper(recpName, recpDescription string) (database.Recipient, string, error) {
	// setup recipient account
	passwd := "p"
	seedpwd := "x"
	recipient, err := database.NewRecipient(recpName, passwd, seedpwd, recpDescription)
	if err != nil {
		return recipient, "", err
	}
	recpSeed, err := wallet.DecryptSeed(recipient.U.EncryptedSeed, seedpwd)
	if err != nil {
		return recipient, "", err
	}
	err = xlm.GetXLM(recipient.U.PublicKey)
	if err != nil {
		return recipient, "", err
	}
	return recipient, recpSeed, nil

}

// this file contains the data that we need to display on the frontend

func I1R1(projIndex int, invName string, invDescription string, recpName string, recpDescription string) error {
	project, err := opensolar.RetrieveProject(projIndex)
	if err != nil {
		log.Fatal(err)
	}

	oldStage := project.Stage
	project.Stage = 4 // to enable investments on this particular project
	err = project.Save()
	if err != nil {
		return err
	}

	// passwd := "p"
	seedpwd := "x"

	investor1, invSeed, err := invHelper(invName, invDescription)
	if err != nil {
		log.Fatal(err)
	}

	recipient1, _, err := recpHelper(recpName, recpDescription)
	if err != nil {
		log.Fatal(err)
	}

	project.RecipientIndex = recipient1.U.Index
	err = project.Save()
	if err != nil {
		return err
	}

	err = opensolar.Invest(projIndex, investor1.U.Index, utils.FtoS(project.TotalValue), invSeed)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Second)
	err = opensolar.UnlockProject(recipient1.U.Username, recipient1.U.Pwhash, projIndex, seedpwd)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(100 * time.Second)
	project.Stage = oldStage
	err = project.Save()
	if err != nil {
		return err
	}
	return nil
}

func createPuertoRicoProject() error {
	// setup all the entities that will be involved with the project here
	projIndex, err := populateStaticDataPR()
	if err != nil {
		log.Fatal(err)
	}

	err = I1R1(projIndex, "OpenLab", "Yale OpenLab", "SUpasto", "S.U. Pasto School")
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func populateStaticData1MW() (int, error) {
	originator1, err := opensolar.NewOriginator("testorig", "p", "x", "testorig", "testorig", "testorig")
	if err != nil {
		log.Fatal(err)
	}

	contractor1, err := opensolar.NewContractor("testcont", "p", "x", "testcont", "testcont", "testcont")
	if err != nil {
		log.Fatal(err)
	}

	developer1, err := opensolar.NewDeveloper("SolarDev", "p", "x", "First Solar", "Solar Rd, San Diego, California", "Main contractor for full solar development")
	if err != nil {
		log.Fatal(err)
	}

	developer2, err := opensolar.NewDeveloper("LancasterSolar", "p", "x", "Town of Lancaste NH", "Lancaster, New Hampshire", "Host")
	if err != nil {
		log.Fatal(err)
	}

	developer3, err := opensolar.NewDeveloper("LancasterRFP", "p", "x", "Lancaster Solar Engineer Solutions", "25 Lancaster Rd, New Hampshire", "Independent RFP Engineer")
	if err != nil {
		log.Fatal(err)
	}

	developer4, err := opensolar.NewDeveloper("Simple Service Provider", "p", "x", "Simple Service Provider", "Simple Service Provider", "Simple Service Provider")
	if err != nil {
		log.Fatal(err)
	}

	developer5, err := opensolar.NewDeveloper("VendorX", "p", "x", "Solar Racking Systems Inc", "34 Crack St, Boston", "Retail Vendor")
	if err != nil {
		log.Fatal(err)
	}

	developer6, err := opensolar.NewDeveloper("NEpool", "p", "x", "New England Pool Registered Auditor", "56 Hamden Ave, Stamford, CT", "REC Auditors for New England")
	if err != nil {
		log.Fatal(err)
	}

	developer7, err := opensolar.NewGuarantor("AllianzCS", "p", "x", "Allianz Climate Solutions", "34 5th, New York, NY", "Insurance Agent")
	if err != nil {
		log.Fatal(err)
	}

	developer8, err := opensolar.NewDeveloper("UIavangrid", "p", "x", "Avangrid Networks", "100 Marsh Hill Rd, New Haven, CT", "Utility")
	if err != nil {
		log.Fatal(err)
	}

	guarantor1, err := opensolar.NewGuarantor("GreenBank", "p", "x", "NH Green Bank", "67 Washington Rd, New Hampshire", "Impact-first escrow provider")
	if err != nil {
		log.Fatal(err)
	}

	var project opensolar.Project
	indexHelp, err := opensolar.RetrieveAllProjects()
	if err != nil {
		log.Fatal(err)
	}

	// This is to populate the table of Terms and Conditions in the front end.
	var terms1 opensolar.TermsHelper
	terms1.Variable = "Security Type"
	terms1.Value = "Equity Notes"
	terms1.RelevantParty = "Lancaster Mutual Solar"
	terms1.Note = "Co-owned by the town of Lancaster"
	terms1.Status = "Issued"
	terms1.SupportDoc = "https://openlab.yale.edu"

	var terms2 opensolar.TermsHelper
	terms2.Variable = "PPA Avg. Tariff"
	terms2.Value = "0.18 ct/KWh"
	terms2.RelevantParty = "Multiple Parties"
	terms2.Note = "Local schools ans business offtaking"
	terms2.Status = "Signed"
	terms2.SupportDoc = "https://openlab.yale.edu"

	var terms3 opensolar.TermsHelper
	terms3.Variable = "Return (TEY)"
	terms3.Value = "4.8%"
	terms3.RelevantParty = "Broker Dealer"
	terms3.Note = "Tax equivalent yield, with capital gains"
	terms3.Status = "Approv"
	terms3.SupportDoc = "https://openlab.yale.edu"

	var terms4 opensolar.TermsHelper
	terms4.Variable = "Maturity"
	terms4.Value = "2026"
	terms4.RelevantParty = "Broker Dealer"
	terms4.Note = "By convertible notes"
	terms4.Status = "Signed"
	terms4.SupportDoc = "https://openlab.yale.edu"

	var terms5 opensolar.TermsHelper
	terms5.Variable = "Guarantee"
	terms5.Value = "50000"
	terms5.RelevantParty = "NH Green Bank"
	terms5.Note = "First-loss escrow upon breach"
	terms5.Status = "Signed"
	terms5.SupportDoc = "https://openlab.yale.edu"

	var terms6 opensolar.TermsHelper
	terms6.Variable = "Insurance"
	terms6.Value = "Premium"
	terms6.RelevantParty = "Allstate"
	terms6.Note = "Force Majeur"
	terms6.Status = "Signed"
	terms6.SupportDoc = "https://openlab.yale.edu"

	var esHelper opensolar.ExecutiveSummaryHelper

	esHelper.Investment = make(map[string]string)
	esHelper.Financials = make(map[string]string)
	esHelper.ProjectSize = make(map[string]string)
	esHelper.SustainabilityMetrics = make(map[string]string)

	esHelper.Investment["Capex"] = "2000000"
	esHelper.Investment["Hardware Ratio"] = "62"
	esHelper.Investment["First-Loss-Escrow"] = "50000"
	esHelper.Investment["Maturity(Fixed)"] = "2029"

	esHelper.Financials["Return (TEY)"] = "4.8"
	esHelper.Financials["Insurance"] = "Premium"
	esHelper.Financials["Tariff (Variable)"] = "0.18 ct/kWh"
	esHelper.Financials["REC Value"] = "$154"

	esHelper.ProjectSize["PV Solar"] = "1 MW"
	esHelper.ProjectSize["Storage"] = "210kWh"
	esHelper.ProjectSize["Array Style"] = "Land & Roof"
	esHelper.ProjectSize["Inverter Capacity"] = "1.25 MW"

	esHelper.SustainabilityMetrics["Carbon Drawdown"] = "0.1 t/kWh"
	esHelper.SustainabilityMetrics["Community Value"] = "6/7"
	esHelper.SustainabilityMetrics["LCA"] = "7/7"

	var communityblock1 opensolar.CommunityEngagementHelper
	communityblock1.Width = 12
	communityblock1.Title = "Climate Education, Awareness & Governance"
	communityblock1.ImageURL = ""
	communityblock1.Content = ""
	communityblock1.Link = ""

	project.Index = len(indexHelp) + 1
	project.Name = "Lancaster Mutual Solar"
	project.State = "NH"
	project.Country = "US"
	project.TotalValue = 2000000
	project.PanelSize = "1MW"
	project.PanelTechnicalDescription = ""
	project.Inverter = ""
	project.ChargeRegulator = ""
	project.ControlPanel = ""
	project.CommBox = ""
	project.ACTransfer = ""
	project.SolarCombiner = ""
	project.Batteries = "210 kWh 1x Tesla Powerpack"
	project.IoTHub = ""
	project.Rating = "AAA"
	project.Metadata = "Neighborhood 1MW solar array on the field next to Lancaster Elementary High School. The project was originated by the head of the community organization, Ben Southworth, who is also active in the parent teacher association (PTA). The city of Lancaster has agreed to give a 20 year lease of the land to the project if the school gets to own the solar array after the lease expires. The school is located in an opportunity zone"

	// Define parameters related to finance
	project.EstimatedAcquisition = 20
	project.BalLeft = -1
	project.InterestRate = 0.05
	project.Tax = "Tax free Opportunity Zone"

	// Define dates of creation and funding
	project.DateInitiated = ""
	project.DateFunded = ""
	project.DateLastPaid = -1

	// Define technical parameters
	project.AuctionType = "blind"
	project.InvestmentType = "munibond"
	project.PaybackPeriod = 4
	project.Stage = 4
	project.SeedInvestmentFactor = 1.1
	project.SeedInvestmentCap = 500
	project.ProposedInvetmentCap = 15000
	project.SelfFund = 0

	// Describe issuer of security and the broker dealer
	project.SecurityIssuer = "Lancaster Mutual Fund"
	project.BrokerDealer = "Neighborly Securities"

	// Define things that will be displayed on the frontend
	project.Terms = append(project.Terms, terms1, terms2, terms3, terms4, terms5, terms6)
	project.ExecutiveSummary = esHelper
	project.AutoReloadInterval = -1
	project.ResilienceRating = 0.8
	project.ActionsRequired = ""
	project.Bullets.Bullet1 = "Eligible for Opportunity Zone Investments"
	project.Bullets.Bullet2 = "Community owned mutual funds as bond issuer"
	project.Bullets.Bullet3 = "1 MW grid with full offtake agreements"
	var hashHelper opensolar.HashHelper
	project.Hashes = hashHelper
	project.ContractList = nil
	project.CommunityEngagement = append(project.CommunityEngagement, communityblock1)
	project.Architecture.SolarOutputImage = ""
	project.Architecture.SolarArray = "1000 kW"
	project.Architecture.DailyAvgGeneration = "4000 kWh"
	project.Architecture.System = "Battery + Grid"
	project.Architecture.InverterSize = "1.25MW"
	project.Architecture.DesignDescription = ""
	project.Context = ""
	project.SummaryImage = ""
	project.ExplorePageSummary.Solar = project.PanelSize
	project.ExplorePageSummary.Storage = esHelper.ProjectSize["Storage"]
	project.ExplorePageSummary.Tariff = esHelper.Financials["Tariff (Variable)"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = esHelper.Financials["Return (TEY)"]
	project.ExplorePageSummary.Rating = project.Rating
	project.ExplorePageSummary.Tax = "N/A"
	project.ExplorePageSummary.ETA = project.EstimatedAcquisition

	project.DPIntroImage = "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/1.png"
	project.OHeroImage = "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/1.png"
	project.OImages = append(project.OImages, "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/6.jpg", "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/7.jpg", "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/4.jpg", "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/3.jpg")
	project.AImages = append(project.AImages, "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/2.jpg", "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/normal.png")
	project.EImages = append(project.EImages, "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/9.png", "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/8.png")
	project.CEImages = append(project.CEImages, "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/5.jpg")
	project.BNImages = append(project.BNImages, "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/10.png", "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/11.png", "https://images.openx.solar/OpenSolarProjects/4_NH_Lancaster/1.png")

	project.FEText = make(map[string]interface{})
	project.FEText, err = parseJsonText("data-sandbox/newhampshire.json")
	if err != nil {
		log.Fatal(err)
	}
	project.EngineeringLayoutType = "complex"
	project.OriginatorIndex = originator1.U.Index
	project.GuarantorIndex = guarantor1.U.Index
	project.ContractorIndex = contractor1.U.Index
	project.MainDeveloperIndex = developer1.U.Index
	project.DeveloperIndices = append(project.DeveloperIndices, developer1.U.Index, developer2.U.Index, developer3.U.Index, developer4.U.Index, developer5.U.Index, developer6.U.Index, developer7.U.Index, developer8.U.Index)
	project.DebtInvestor1 = "OZFunds"
	project.DebtInvestor2 = "GreenBank"
	project.TaxEquityInvestor = "TaxEquity"
	err = project.Save()
	if err != nil {
		log.Fatal(err)
	}
	return project.Index, nil
}

func I3R1(projIndex int, invName1 string, invDescription1 string, invName2 string, invDescription2 string,
	invName3 string, invDescription3 string, invAmount1 string, invAmount2 string, invAmount3 string,
	recpName string, recpDescription string) error {

	project, err := opensolar.RetrieveProject(projIndex)
	if err != nil {
		log.Fatal(err)
	}

	oldStage := project.Stage
	project.Stage = 4 // to enable investments on this particular project
	err = project.Save()
	if err != nil {
		return err
	}

	// passwd := "p"
	// seedpwd := "x"

	investor1, invSeed1, err := invHelper(invName1, invDescription1)
	if err != nil {
		log.Fatal(err)
	}

	investor2, invSeed2, err := invHelper(invName2, invDescription2)
	if err != nil {
		log.Fatal(err)
	}

	investor3, invSeed3, err := invHelper(invName3, invDescription3)
	if err != nil {
		log.Fatal(err)
	}

	recipient1, _, err := recpHelper(recpName, recpDescription)
	if err != nil {
		log.Fatal(err)
	}

	project.RecipientIndex = recipient1.U.Index
	err = project.Save()
	if err != nil {
		return err
	}

	err = opensolar.Invest(projIndex, investor1.U.Index, invAmount1, invSeed1)
	if err != nil {
		log.Fatal(err)
	}

	err = opensolar.Invest(projIndex, investor2.U.Index, invAmount2, invSeed2)
	if err != nil {
		log.Fatal(err)
	}

	err = opensolar.Invest(projIndex, investor3.U.Index, invAmount3, invSeed3)
	if err != nil {
		log.Fatal(err)
	}

	// update local project with changes from storage
	project, err = opensolar.RetrieveProject(projIndex)
	if err != nil {
		log.Fatal(err)
	}

	project.BlendedCapitalInvestorIndex = investor1.U.Index
	project.Stage = oldStage
	err = project.Save()
	if err != nil {
		return err
	}
	return nil
}

func createOneMegaWattProject() error {
	// setup all the entities involved with the project here
	projIndex, err := populateStaticData1MW()
	if err != nil {
		log.Fatal(err)
	}

	err = I3R1(projIndex, "OZFunds", "OZ FundCo", "GreenBank", "NH Green Bank", "TaxEquity", "Lancaster Lumber Mill Coop",
		"1000000", "400000", "100000", "LancasterHigh", "Lancaster Elementary School")
	if err != nil {
		log.Fatal(err)
	}

	recipient2, err := database.NewRecipient("LancasterT", "p", "x", "Town of Lancaster NH")
	if err != nil {
		log.Fatal(err)
	}

	project, err := opensolar.RetrieveProject(projIndex)
	if err != nil {
		log.Fatal(err)
	}

	// Define the various entities that are associated with a specific project
	project.RecipientIndices = append(project.RecipientIndices, recipient2.U.Index)
	err = project.Save()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func populateStaticData10KW() (int, error) {
	developer1, err := opensolar.NewDeveloper("YaleArchitecture", "p", "x", "Yale School of Architecture", "45 York St, New Haven, CT", "System and layout designer")
	if err != nil {
		log.Fatal(err)
	}

	developer2, err := opensolar.NewDeveloper("CTSolar", "p", "x", "Connecticut Solar", "45 Sun Street, Stamford, CT", "Solar system installer")
	if err != nil {
		log.Fatal(err)
	}

	developer3, err := opensolar.NewDeveloper("ColumbusHouse", "p", "x", "Columbus House", "21 Hagrid Ave, New Haven, CT", "Project Host")
	if err != nil {
		log.Fatal(err)
	}

	// We have in these examples one user that is covering different roles. And this is something good for the demo eventually. The example is Raise Green (both originator and guarantor) and Avangrid (REC developer here, Utility in another project)
	// How should we create these users so that they have these different entity properties?
	developer4, err := opensolar.NewGuarantor("RGreenFund", "p", "x", "RaiseGreen Blend Fund", "21 orange st, New Haven, CT", "Impact-first blended capital provider")
	if err != nil {
		log.Fatal(err)
	}

	developer5, err := opensolar.NewDeveloper("Avangrid", "p", "x", "Avangrid RECs", "100 Marsh Hill Rd, New Haven, CT", "Certifier of RECs and provider of REC meter")
	if err != nil {
		log.Fatal(err)
	}

	originator1, err := opensolar.NewOriginator("RaiseGreen", "p", "x", "Raise Green", "21 orange st, New Haven, CT", "Project originator")
	if err != nil {
		log.Fatal(err)
	}

	contractor1, err := opensolar.NewContractor("testcont", "p", "x", "testcont", "testcont", "testcont")
	if err != nil {
		log.Fatal(err)
	}

	guarantor1, err := opensolar.NewGuarantor("testguarantor", "p", "x", "testguarantor", "testguarantor", "testguarantor")
	if err != nil {
		log.Fatal(err)
	}

	var project opensolar.Project
	indexHelp, err := opensolar.RetrieveAllProjects()
	if err != nil {
		log.Fatal(err)
	}

	var esHelper opensolar.ExecutiveSummaryHelper

	esHelper.Investment = make(map[string]string)
	esHelper.Financials = make(map[string]string)
	esHelper.ProjectSize = make(map[string]string)
	esHelper.SustainabilityMetrics = make(map[string]string)

	esHelper.Investment["Capex"] = "30000"
	esHelper.Investment["Hardware"] = "$3 Non Voting"
	esHelper.Investment["Raise Type"] = "Reg CF"
	esHelper.Investment["Certification"] = "N/A"

	esHelper.Financials["Equity Value"] = "130%"
	esHelper.Financials["Insurance"] = "Basic"
	esHelper.Financials["Tariff (Fixed)"] = "0.12 ct/kWh"
	esHelper.Financials["Maturity (Fixed)"] = "2019"

	esHelper.ProjectSize["PV Solar"] = "10 kW"
	esHelper.ProjectSize["Storage"] = "N/A Grid Tied"
	esHelper.ProjectSize["% Critical"] = "100"
	esHelper.ProjectSize["Inverter Capacity"] = "15 kW"

	esHelper.SustainabilityMetrics["Carbon Drawdown"] = "0.1 t/kWh"
	esHelper.SustainabilityMetrics["Community Value"] = "7/7"
	esHelper.SustainabilityMetrics["LCA"] = ""

	var communityblock1 opensolar.CommunityEngagementHelper
	communityblock1.Width = 12
	communityblock1.Title = "Consultation"
	communityblock1.ImageURL = ""
	communityblock1.Content = ""
	communityblock1.Link = ""

	// This is to populate the table of Terms and Conditions in the front end.
	var terms1 opensolar.TermsHelper
	terms1.Variable = "Security Type"
	terms1.Value = "Reg CF"
	terms1.RelevantParty = "NH Community Solar"
	terms1.Note = "Special Purpose Vehicle"
	terms1.Status = "Flipped"
	terms1.SupportDoc = "https://openlab.yale.edu"

	var terms2 opensolar.TermsHelper
	terms2.Variable = "PPA Tariff"
	terms2.Value = "0.12 ct/KWh"
	terms2.RelevantParty = "NH Homeless Shelter"
	terms2.Note = "Fixed PPA determined by offtaker"
	terms2.Status = "Signed"
	terms2.SupportDoc = "https://openlab.yale.edu"

	var terms3 opensolar.TermsHelper
	terms3.Variable = "Return"
	terms3.Value = "130%"
	terms3.RelevantParty = "Equity Value"
	terms3.Note = "Growth in value. No tax incentives"
	terms3.Status = "Open"
	terms3.SupportDoc = "https://openlab.yale.edu"

	var terms4 opensolar.TermsHelper
	terms4.Variable = "Ownership Flip"
	terms4.Value = "2019"
	terms4.RelevantParty = "Convertible Note"
	terms4.Note = "Crowd investors sell stock"
	terms4.Status = "Flipped"
	terms4.SupportDoc = "https://openlab.yale.edu"

	var terms5 opensolar.TermsHelper
	terms5.Variable = "Guarantee"
	terms5.Value = "N/A"
	terms5.RelevantParty = "N/A"
	terms5.Note = "No guarantees of breach"
	terms5.Status = "None"
	terms5.SupportDoc = "https://openlab.yale.edu"

	var terms6 opensolar.TermsHelper
	terms6.Variable = "Insurance"
	terms6.Value = "Basic"
	terms6.RelevantParty = "CT Insurers"
	terms6.Note = "Force Majeur"
	terms6.Status = "Signed"
	terms6.SupportDoc = "https://openlab.yale.edu"

	project.Index = len(indexHelp) + 1
	project.Name = "New Haven Shelter Solar 2"
	project.State = "CT"
	project.Country = "US"
	project.TotalValue = 30000
	project.PanelSize = "15kW"
	project.PanelTechnicalDescription = ""
	project.Inverter = ""
	project.ChargeRegulator = ""
	project.ControlPanel = ""
	project.CommBox = ""
	project.ACTransfer = ""
	project.SolarCombiner = ""
	project.Batteries = ""
	project.IoTHub = ""
	project.Rating = "Premium"
	project.Metadata = "Residential solar array for a homeless shelter. The project was originated by a member of the board of the homeless shelter who gets the shelter to purchase all the electricity at a discounted rate. The shelter chooses to lease the roof for free over the lifetime of the project. The originator knows the solar developer who set up the project company"

	// Define parameters related to finance
	project.EstimatedAcquisition = 0 // this project already flipped ownership
	project.BalLeft = 0
	project.InterestRate = 0.05
	project.Tax = "0.3 Tax Credit"

	// Define dates of creation and funding
	project.DateInitiated = ""
	project.DateFunded = ""
	project.DateLastPaid = -1

	// Define technical parameters
	project.AuctionType = "blind"
	project.InvestmentType = "munibond"
	project.PaybackPeriod = 4
	project.Stage = 8
	project.SeedInvestmentFactor = 1.1
	project.SeedInvestmentCap = 500
	project.ProposedInvetmentCap = 15000
	project.SelfFund = 0

	// Describe issuer of security and the broker dealer
	project.SecurityIssuer = ""
	project.BrokerDealer = ""

	// Define things that will be displayed on the frontend
	project.Terms = append(project.Terms, terms1, terms2, terms3, terms4, terms5, terms6)
	project.ExecutiveSummary = esHelper
	project.AutoReloadInterval = -1
	project.ResilienceRating = 0.6
	project.ActionsRequired = ""
	project.Bullets.Bullet1 = "Community owned solar in homeless shelter"
	project.Bullets.Bullet2 = "Siginificantly alleviates financial pressure due to high CT power cost"
	project.Bullets.Bullet3 = "Grid-tied with REC offtaking"
	var hashHelper opensolar.HashHelper
	project.Hashes = hashHelper
	project.ContractList = nil
	project.CommunityEngagement = append(project.CommunityEngagement, communityblock1)
	project.Architecture.SolarOutputImage = ""
	project.Architecture.SolarArray = "50 x 200W"
	project.Architecture.DailyAvgGeneration = "20000 Wh"
	project.Architecture.System = "Grid Tied"
	project.Architecture.InverterSize = "15 kW"
	project.Architecture.DesignDescription = ""
	project.Context = ""
	project.SummaryImage = ""
	project.ExplorePageSummary.Solar = project.PanelSize
	project.ExplorePageSummary.Storage = esHelper.ProjectSize["Storage"]
	project.ExplorePageSummary.Tariff = esHelper.Financials["Tariff (Fixed)"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = esHelper.Financials["Return (TEY)"]
	project.ExplorePageSummary.Rating = project.Rating
	project.ExplorePageSummary.Tax = "N/A"
	project.ExplorePageSummary.ETA = project.EstimatedAcquisition

	project.DPIntroImage = "https://images.openx.solar/OpenSolarProjects/8_NewHaven/3.png"
	project.OHeroImage = ""
	project.OImages = append(project.OImages, "https://images.openx.solar/OpenSolarProjects/8_NewHaven/4.jpg", "https://images.openx.solar/OpenSolarProjects/8_NewHaven/12.jpg")
	project.OOImages = append(project.OOImages, "https://images.openx.solar/OpenSolarProjects/8_NewHaven/6.jpg")
	project.AImages = append(project.AImages, "https://images.openx.solar/OpenSolarProjects/8_NewHaven/7.jpg", "https://images.openx.solar/OpenSolarProjects/8_NewHaven/normal.png")
	project.EImages = append(project.EImages, "https://images.openx.solar/OpenSolarProjects/8_NewHaven/9.jpg", "https://images.openx.solar/OpenSolarProjects/8_NewHaven/10.jpg")
	project.CEImages = append(project.CEImages, "https://images.openx.solar/OpenSolarProjects/8_NewHaven/12.jpg")
	project.PSImages = append(project.PSImages, "https://images.openx.solar/OpenSolarProjects/8_NewHaven/11.jpg")
	project.BNImages = append(project.BNImages, "")

	project.EngineeringLayoutType = "basic"

	project.FEText = make(map[string]interface{})
	project.FEText, err = parseJsonText("data-sandbox/newhaven.json")
	if err != nil {
		log.Fatal(err)
	}
	project.OriginatorIndex = originator1.U.Index
	project.GuarantorIndex = guarantor1.U.Index
	project.ContractorIndex = contractor1.U.Index
	project.MainDeveloperIndex = developer1.U.Index
	project.DeveloperIndices = append(project.DeveloperIndices, developer1.U.Index, developer2.U.Index, developer3.U.Index, developer4.U.Index, developer5.U.Index)

	err = project.Save()
	if err != nil {
		log.Fatal(err)
	}

	return project.Index, nil
}

func I6R1(projIndex int, invName1 string, invDescription1 string, invName2 string, invDescription2 string,
	invName3 string, invDescription3 string, invName4 string, invDescription4 string, invName5 string, invDescription5 string,
	invName6 string, invDescription6 string, invAmount1 string, invAmount2 string, invAmount3 string, invAmount4 string,
	invAmount5 string, invAmount6 string, recpName string, recpDescription string) error {

	project, err := opensolar.RetrieveProject(projIndex)
	if err != nil {
		log.Fatal(err)
	}

	recipient1, _, err := recpHelper(recpName, recpDescription)
	if err != nil {
		log.Fatal(err)
	}

	oldStage := project.Stage
	project.RecipientIndex = recipient1.U.Index
	project.Stage = 4 // to enable investments on this particular project
	err = project.Save()
	if err != nil {
		return err
	}

	// passwd := "p"
	seedpwd := "x"
	investor1, invSeed1, err := invHelper(invName1, invDescription1)
	if err != nil {
		log.Fatal(err)
	}

	investor2, invSeed2, err := invHelper(invName2, invDescription2)
	if err != nil {
		log.Fatal(err)
	}

	investor3, invSeed3, err := invHelper(invName3, invDescription3)
	if err != nil {
		log.Fatal(err)
	}

	investor4, invSeed4, err := invHelper(invName4, invDescription4)
	if err != nil {
		log.Fatal(err)
	}

	investor5, invSeed5, err := invHelper(invName5, invDescription5)
	if err != nil {
		log.Fatal(err)
	}

	investor6, invSeed6, err := invHelper(invName6, invDescription6)
	if err != nil {
		log.Fatal(err)
	}

	err = project.Save()
	if err != nil {
		return err
	}

	err = opensolar.Invest(projIndex, investor1.U.Index, invAmount1, invSeed1)
	if err != nil {
		log.Fatal(err)
	}

	err = opensolar.Invest(projIndex, investor2.U.Index, invAmount2, invSeed2)
	if err != nil {
		log.Fatal(err)
	}

	err = opensolar.Invest(projIndex, investor3.U.Index, invAmount3, invSeed3)
	if err != nil {
		log.Fatal(err)
	}

	err = opensolar.Invest(projIndex, investor4.U.Index, invAmount4, invSeed4)
	if err != nil {
		log.Fatal(err)
	}

	err = opensolar.Invest(projIndex, investor5.U.Index, invAmount5, invSeed5)
	if err != nil {
		log.Fatal(err)
	}

	err = opensolar.Invest(projIndex, investor6.U.Index, invAmount6, invSeed6)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Second)
	err = opensolar.UnlockProject(recipient1.U.Username, recipient1.U.Pwhash, projIndex, seedpwd)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(100 * time.Second)
	project, err = opensolar.RetrieveProject(projIndex)
	if err != nil {
		log.Fatal(err)
	}

	project.BlendedCapitalInvestorIndex = investor1.U.Index
	project.Stage = oldStage
	err = project.Save()
	if err != nil {
		return err
	}
	return nil
}

func createTenKiloWattProject() error {

	projIndex, err := populateStaticData10KW()
	if err != nil {
		log.Fatal(err)
	}

	err = I6R1(projIndex, "MatthewMoroney", "Matthew Moroney", "FranzHochstrasser", "Franz Hochstrasser", "CTGreenBank", "Connecticut Green Bank",
		"YaleUniversity", "Yale University Community Fund", "JeromeGreen", "Jerome Green", "OpenSolarFund", "Open Solar Revolving Fund",
		"4000", "4000", "4000", "4000", "4000", "10000", "colhouse", "Columbus House Foundation")

	if err != nil {
		log.Fatal(err)
	}

	recipient2, err := database.NewRecipient("ColumbusHouse", "p", "x", "Columbus House Foundation")
	if err != nil {
		log.Fatal(err)
	}

	project, err := opensolar.RetrieveProject(projIndex)
	if err != nil {
		log.Fatal(err)
	}
	// Define the various entities that are associated with a specific project
	project.RecipientIndices = append(project.RecipientIndices, recipient2.U.Index)
	err = project.Save()
	if err != nil {
		return err
	}
	return nil
}

func createTenMegaWattProject() error {
	// create the required entities that we need over here
	investor1, err := database.NewInvestor("emcoll", "p", "x", "Emerson Collective")
	if err != nil {
		log.Fatal(err)
	}

	investor2, err := database.NewInvestor("prqozfund", "p", "x", "Puerto Rico QOZ Fund")
	if err != nil {
		log.Fatal(err)
	}

	recipient1, err := database.NewRecipient("prgov", "p", "x", "PR Government")
	if err != nil {
		log.Fatal(err)
	}

	recipient2, err := database.NewRecipient("prschools", "p", "x", "Puerto Rico Solar Schools Limited")
	if err != nil {
		log.Fatal(err)
	}

	recipient3, err := database.NewRecipient("prdoe", "p", "x", "Puerto Rico Department of Education")
	if err != nil {
		log.Fatal(err)
	}

	originator1, err := opensolar.NewOriginator("MartinWainstein1", "p", "x", "Martin Wainstein", "254 Elm Street, New Haven, CT", "Martin Wainstein from the Yale OpenLab")
	if err != nil {
		log.Fatal(err)
	}

	developer1, err := opensolar.NewDeveloper("hst", "p", "x", "HST Solar", "25 Hewlett St, San Francisco, CA", "Preliminary finance and engineering assessment")
	if err != nil {
		log.Fatal(err)
	}

	developer2, err := opensolar.NewDeveloper("FemaRoofs", "p", "x", "FEMA Puerto Rico", "“45 Old Town Rd, Puerto Rico", "Civil engineering assessment of school roofs")
	if err != nil {
		log.Fatal(err)
	}

	contractor1, err := opensolar.NewContractor("testcont", "p", "x", "testcont", "testcont", "testcont")
	if err != nil {
		log.Fatal(err)
	}

	guarantor1, err := opensolar.NewGuarantor("testguarantor", "p", "x", "testguarantor", "testguarantor", "testguarantor")
	if err != nil {
		log.Fatal(err)
	}

	var project opensolar.Project
	indexHelp, err := opensolar.RetrieveAllProjects()
	if err != nil {
		log.Fatal(err)
	}

	// This is to populate the table of Terms and Conditions in the front end.
	var terms1 opensolar.TermsHelper
	terms1.Variable = "Security Type"
	terms1.Value = "Municipal Bond"
	terms1.RelevantParty = "PR DofEd"
	terms1.Note = "Not yet issued. See informal Agreements"
	terms1.Status = "Open"
	terms1.SupportDoc = "https://openlab.yale.edu"

	var terms2 opensolar.TermsHelper
	terms2.Variable = "PPA Tariff"
	terms2.Value = "0.24 ct/KWh"
	terms2.RelevantParty = "Oracle X / PREPA"
	terms2.Note = "Not signed. Expected as variable tariff"
	terms2.Status = "Open"
	terms2.SupportDoc = "https://openlab.yale.edu"

	var terms3 opensolar.TermsHelper
	terms3.Variable = "Return (TEY)"
	terms3.Value = "3.5%"
	terms3.RelevantParty = "See Broker Dealer"
	terms3.Note = "Tax equivalent yield with capital gains"
	terms3.Status = "Open"
	terms3.SupportDoc = "https://openlab.yale.edu"

	var terms4 opensolar.TermsHelper
	terms4.Variable = "Maturity"
	terms4.Value = "+/- 2025"
	terms4.RelevantParty = "Broker Dealer"
	terms4.Note = "Variable tied to tariff"
	terms4.Status = "Open"
	terms4.SupportDoc = "https://openlab.yale.edu"

	var terms5 opensolar.TermsHelper
	terms5.Variable = "Guarantee"
	terms5.Value = "15%"
	terms5.RelevantParty = "FEMA"
	terms5.Note = "First-loss upon breach"
	terms5.Status = "Started"
	terms5.SupportDoc = "https://openlab.yale.edu"

	var terms6 opensolar.TermsHelper
	terms6.Variable = "Insurance"
	terms6.Value = "Premium"
	terms6.RelevantParty = "Allianz CS"
	terms6.Note = "Hurricane Coverage"
	terms6.Status = "Started"
	terms6.SupportDoc = "https://openlab.yale.edu"

	var esHelper opensolar.ExecutiveSummaryHelper

	esHelper.Investment = make(map[string]string)
	esHelper.Financials = make(map[string]string)
	esHelper.ProjectSize = make(map[string]string)
	esHelper.SustainabilityMetrics = make(map[string]string)

	esHelper.Investment["Capex"] = "19000000"
	esHelper.Investment["Hardware"] = "70"
	esHelper.Investment["First-Loss-Escrow"] = "15%"
	esHelper.Investment["Maturity"] = ""

	esHelper.Financials["Expected Return (Non TEY)"] = "2.5%"
	esHelper.Financials["Insurance"] = "Basic Force Majeur"
	esHelper.Financials["Tariff (Variable)"] = "0.24 ct/kWh"
	esHelper.Financials["REC Value"] = "$234/MWh"

	esHelper.ProjectSize["PV Solar"] = "300 x 30kW"
	esHelper.ProjectSize["Storage"] = "350 x 2.5 kWh"
	esHelper.ProjectSize["% Critical"] = "20"
	esHelper.ProjectSize["Inverter Capacity"] = "300 x 35 kW"

	esHelper.SustainabilityMetrics["Carbon Drawdown"] = "0.1t/kWh"
	esHelper.SustainabilityMetrics["Community Value"] = "6/7"
	esHelper.SustainabilityMetrics["LCA"] = "N/A"

	var communityblock1 opensolar.CommunityEngagementHelper
	communityblock1.Width = 3
	communityblock1.Title = "Consultation"
	communityblock1.ImageURL = ""
	communityblock1.Content = ""
	communityblock1.Link = ""

	var communityblock2 opensolar.CommunityEngagementHelper
	communityblock2.Width = 3
	communityblock2.Title = "Participation"
	communityblock2.ImageURL = ""
	communityblock2.Content = ""
	communityblock2.Link = ""

	var communityblock3 opensolar.CommunityEngagementHelper
	communityblock3.Width = 3
	communityblock3.Title = "Outreach"
	communityblock3.ImageURL = ""
	communityblock3.Content = ""
	communityblock3.Link = ""

	var communityblock4 opensolar.CommunityEngagementHelper
	communityblock4.Width = 3
	communityblock4.Title = "Governance"
	communityblock4.ImageURL = ""
	communityblock4.Content = ""
	communityblock4.Link = ""

	project.Index = len(indexHelp) + 1
	project.Name = "Puerto Rico Solar School Bond 1"
	project.State = "Puerto Rico"
	project.Country = "US"
	project.TotalValue = 19000000
	project.PanelSize = "10MW"
	project.PanelTechnicalDescription = ""
	project.Inverter = ""
	project.ChargeRegulator = ""
	project.ControlPanel = ""
	project.CommBox = ""
	project.ACTransfer = ""
	project.SolarCombiner = ""
	project.Batteries = "900 Ah"
	project.IoTHub = ""
	project.Rating = "AA+"
	project.Metadata = "Transformation of 300 Puerto Rican public schools into solar powered emergency shelters. Each school will have around 30kW solar and 2kWh battery bank to cover critical loads including refrigeration of food and medicine, and an emergency telecommunication system with first responders. Backed by the Office of the Governor. 10 MW aggregate solar capacity. Nodes for community microgrids"

	// Define parameters related to finance
	project.MoneyRaised = 0
	project.EstimatedAcquisition = 8
	project.BalLeft = -1
	project.InterestRate = 0.029
	project.Tax = ""

	// Define dates of creation and funding
	project.DateInitiated = "01/23/2019"
	project.DateFunded = ""
	project.DateLastPaid = -1

	// Define technical parameters
	project.AuctionType = "private"
	project.InvestmentType = "munibond"
	project.PaybackPeriod = 4
	project.Stage = 2
	project.SeedInvestmentFactor = 1.1
	project.SeedInvestmentCap = 500
	project.ProposedInvetmentCap = 15000
	project.SelfFund = 0

	// Describe issuer of security and the broker dealer
	project.SecurityIssuer = "Neighborly Securities"
	project.BrokerDealer = "Broker Dealer"

	// Define the various entities that are associated with a specific project
	project.RecipientIndex = recipient1.U.Index
	project.OriginatorIndex = originator1.U.Index
	project.GuarantorIndex = guarantor1.U.Index
	project.ContractorIndex = contractor1.U.Index
	project.MainDeveloperIndex = developer1.U.Index
	project.BlendedCapitalInvestorIndex = -1
	project.InvestorIndices = append(project.InvestorIndices, investor1.U.Index, investor2.U.Index)
	project.SeedInvestorIndices = nil
	project.RecipientIndices = append(project.RecipientIndices, recipient1.U.Index, recipient2.U.Index, recipient3.U.Index)
	project.DeveloperIndices = append(project.DeveloperIndices, developer1.U.Index, developer2.U.Index)
	project.ContractorFee = 2000
	project.OriginatorFee = 0
	project.DeveloperFee = append(project.DeveloperFee, 6000)
	project.DebtInvestor1 = ""
	project.DebtInvestor2 = ""
	project.TaxEquityInvestor = ""

	// Define things that will be displayed on the frontend
	project.Terms = append(project.Terms, terms1, terms2, terms3, terms4, terms5, terms6)
	project.ExecutiveSummary = esHelper
	project.AutoReloadInterval = -1
	project.ResilienceRating = 0.75
	project.ActionsRequired = ""
	project.Bullets.Bullet1 = "Backed by the Governor's office"
	project.Bullets.Bullet2 = "Critical Loads covered (telecom and refrigeration)"
	project.Bullets.Bullet3 = "Certification of social impact"
	var hashHelper opensolar.HashHelper
	project.Hashes = hashHelper
	project.ContractList = nil
	project.CommunityEngagement = append(project.CommunityEngagement, communityblock1, communityblock2, communityblock3, communityblock4)
	project.Architecture.SolarOutputImage = ""
	project.Architecture.SolarArray = "300 x 30kW"
	project.Architecture.DailyAvgGeneration = "400 MWh"
	project.Architecture.System = "350 Tesla Powerwells"
	project.Architecture.InverterSize = "300 x 35 kW"
	project.Architecture.DesignDescription = ""
	project.Context = ""
	project.SummaryImage = ""
	project.ExplorePageSummary.Solar = project.PanelSize
	project.ExplorePageSummary.Storage = esHelper.ProjectSize["Storage"]
	project.ExplorePageSummary.Tariff = esHelper.Financials["Tariff (Variable)"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = esHelper.Financials["Expected Return (Non TEY)"]
	project.ExplorePageSummary.Rating = project.Rating
	project.ExplorePageSummary.Tax = "N/A"
	project.ExplorePageSummary.ETA = project.EstimatedAcquisition

	project.DPIntroImage = "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/1.jpg"
	project.OHeroImage = "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/10.png"
	project.OImages = append(project.OImages, "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/3.jpg", "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/8.jpg")
	project.OOImages = append(project.OOImages, "")
	project.AImages = append(project.AImages, "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/9.png", "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/normal.png", "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/2.jpg")
	project.EImages = append(project.EImages, "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/11.png", "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/12.png", "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/13.png")
	project.CEImages = append(project.CEImages, "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/7.jpg", "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/14.jpg", "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/15.jpg", "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/16.jpg")
	project.PSImages = append(project.PSImages, "https://images.openx.solar/OpenSolarProjects/2_PR_Bonds/17.jpg")
	project.BNImages = append(project.BNImages, "")

	project.EngineeringLayoutType = "simple"
	project.FEText = make(map[string]interface{})
	project.FEText, err = parseJsonText("data-sandbox/prbonds.json")
	if err != nil {
		log.Fatal(err)
	}

	err = project.Save()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func createOneHundredKiloWattProject() error {

	investor1, err := database.NewInvestor("jjackson", "p", "x", "Jerome Jackson")
	if err != nil {
		log.Fatal(err)
	}

	investor2, err := database.NewInvestor("esare", "p", "x", "Eliah Sare")
	if err != nil {
		log.Fatal(err)
	}

	investor3, err := database.NewInvestor("yaleuf", "p", "x", "Yale University Fund")
	if err != nil {
		log.Fatal(err)
	}

	recipient1, err := database.NewRecipient("ubaduef", "p", "x", "Ubadu Energy Collective")
	if err != nil {
		log.Fatal(err)
	}

	recipient2, err := database.NewRecipient("sunshinegschool", "p", "x", "Sunshine Garden School")
	if err != nil {
		log.Fatal(err)
	}

	recipient3, err := database.NewRecipient("ubaduth", "p", "x", "Ubadu Town Hall")
	if err != nil {
		log.Fatal(err)
	}

	recipient4, err := database.NewRecipient("dwbrf", "p", "x", " Doctors without borders, Rwanda chapter")
	if err != nil {
		log.Fatal(err)
	}

	recipient5, err := database.NewRecipient("largerof", "p", "x", "Large Residential offtakers")
	if err != nil {
		log.Fatal(err)
	}

	originator1, err := opensolar.NewOriginator("DjiembeMbeba", "p", "x", "Djiembe Mbeba", "Ubadu village, Rwanda", "Project originator")
	if err != nil {
		log.Fatal(err)
	}

	developer1, err := opensolar.NewDeveloper("SolarPartners", "p", "x", "Solar Partners", "34 Hiete st, Somaliland", "MiniGrid game developer")
	if err != nil {
		log.Fatal(err)
	}

	developer2, err := opensolar.NewDeveloper("hst2", "p", "x", "HST Solar", "25 Hewlett St, San Francisco, CA", "Preliminary finance and engineering assessment")
	if err != nil {
		log.Fatal(err)
	}

	var project opensolar.Project
	indexHelp, err := opensolar.RetrieveAllProjects()
	if err != nil {
		log.Fatal(err)
	}

	// This is to populate the table of Terms and Conditions in the front end
	var terms1 opensolar.TermsHelper
	terms1.Variable = "Security Type"
	terms1.Value = "Equity Crowdfunding"
	terms1.RelevantParty = "Ubadu Collective"
	terms1.Note = "Coop is not incorporated yet"
	terms1.Status = "Open"
	terms1.SupportDoc = "https://openlab.yale.edu"

	var terms2 opensolar.TermsHelper
	terms2.Variable = "PPA Tariff"
	terms2.Value = "0.12 ct/KWh"
	terms2.RelevantParty = "Ubadu Collective"
	terms2.Note = "Average PPA, from tiered offtakers"
	terms2.Status = "Open"
	terms2.SupportDoc = "https://openlab.yale.edu"

	var terms3 opensolar.TermsHelper
	terms3.Variable = "Exp. Return"
	terms3.Value = "2.3%"
	terms3.RelevantParty = "Equity Value"
	terms3.Note = "Growth value. No tax incentives"
	terms3.Status = "Open"
	terms3.SupportDoc = "https://openlab.yale.edu"

	var terms4 opensolar.TermsHelper
	terms4.Variable = "Ownership Flip"
	terms4.Value = "2027"
	terms4.RelevantParty = "By convertible notes"
	terms4.Note = "Crowd investors sell stock"
	terms4.Status = "Open"
	terms4.SupportDoc = "https://openlab.yale.edu"

	var terms5 opensolar.TermsHelper
	terms5.Variable = "Guarantee"
	terms5.Value = "20%"
	terms5.RelevantParty = "Africa Fund"
	terms5.Note = "Agreed but pending"
	terms5.Status = "Open"
	terms5.SupportDoc = "https://openlab.yale.edu"

	var terms6 opensolar.TermsHelper
	terms6.Variable = "Insurance"
	terms6.Value = "N/A"
	terms6.RelevantParty = "N/A"
	terms6.Note = "Defining insurance parties"
	terms6.Status = "Open"
	terms6.SupportDoc = "https://openlab.yale.edu"

	var esHelper opensolar.ExecutiveSummaryHelper

	esHelper.Investment = make(map[string]string)
	esHelper.Financials = make(map[string]string)
	esHelper.ProjectSize = make(map[string]string)
	esHelper.SustainabilityMetrics = make(map[string]string)

	esHelper.Investment["Capex"] = "230000"
	esHelper.Investment["Hardware"] = "75"
	esHelper.Investment["FirstLossEscrow"] = "Equity w/Notes" // cahnge this to raise type
	esHelper.Investment["CertificationCosts"] = "2028"        // change to maturity

	esHelper.Financials["Return"] = "2.3%"
	esHelper.Financials["First-Loss Escrow"] = "20%"
	esHelper.Financials["Tariff"] = "0.24 ct/kWh"
	esHelper.Financials["REC Value"] = "In Process"

	esHelper.ProjectSize["PVSolar"] = "4 x 25 kW"
	esHelper.ProjectSize["Storage"] = "25 kWh"
	esHelper.ProjectSize["Critical"] = "100"
	esHelper.ProjectSize["InverterCapacity"] = "4 x 30 kW"

	esHelper.SustainabilityMetrics["CarbonDrawdown"] = "N/A"
	esHelper.SustainabilityMetrics["CommunityValue"] = "7/7"
	esHelper.SustainabilityMetrics["LCA"] = "N/A"

	var communityblock1 opensolar.CommunityEngagementHelper
	communityblock1.Title = "Consultation"
	communityblock1.ImageURL = ""
	communityblock1.Content = ""
	communityblock1.Link = ""

	var communityblock2 opensolar.CommunityEngagementHelper
	communityblock2.Title = "Participation"
	communityblock2.ImageURL = ""
	communityblock2.Content = ""
	communityblock2.Link = ""

	var communityblock3 opensolar.CommunityEngagementHelper
	communityblock3.Title = "Outreach"
	communityblock3.ImageURL = ""
	communityblock3.Content = ""
	communityblock3.Link = ""

	var communityblock4 opensolar.CommunityEngagementHelper
	communityblock4.Title = "Governance"
	communityblock4.ImageURL = ""
	communityblock4.Content = ""
	communityblock4.Link = ""

	project.Index = len(indexHelp) + 1
	project.Name = "Ubadu Village Microgrid, Rwanda"
	project.State = "Kigali"
	project.Country = "Rwanda"
	project.TotalValue = 10000
	project.PanelSize = "1kW"
	project.PanelTechnicalDescription = ""
	project.Inverter = ""
	project.ChargeRegulator = ""
	project.ControlPanel = ""
	project.CommBox = ""
	project.ACTransfer = ""
	project.SolarCombiner = ""
	project.Batteries = ""
	project.IoTHub = "Yale Open Powermeter w/ RaspberryPi3"
	project.Rating = "N/A"
	project.Metadata = "The community of Ubadu, Rwanda has no access to grid electricity yet shows a  growing local economy. This microgrid project, will serve 250 homes, a school, an infirmary and the town hall"

	// Define parameters related to finance
	project.MoneyRaised = 0
	project.EstimatedAcquisition = 5
	project.BalLeft = 10000
	project.InterestRate = 0.029
	project.Tax = "Insert tax scheme here"

	// Define dates of creation and funding
	project.DateInitiated = ""
	project.DateFunded = ""
	project.DateLastPaid = -1

	// Define technical parameters
	project.AuctionType = "blind"
	project.InvestmentType = "munibond"
	project.PaybackPeriod = 4
	project.Stage = 1
	project.SeedInvestmentFactor = 1.1
	project.SeedInvestmentCap = 500
	project.ProposedInvetmentCap = 15000
	project.SelfFund = 0

	// Describe issuer of security and the broker dealer
	project.SecurityIssuer = ""
	project.BrokerDealer = ""

	// Define the various entities that are associated with a specific project
	project.RecipientIndex = recipient1.U.Index
	project.OriginatorIndex = originator1.U.Index
	// project.GuarantorIndex = guarantor1.U.Index
	// project.ContractorIndex = contractor1.U.Index
	project.MainDeveloperIndex = developer1.U.Index
	project.BlendedCapitalInvestorIndex = -1
	project.InvestorIndices = append(project.InvestorIndices, investor1.U.Index, investor2.U.Index, investor3.U.Index)
	project.SeedInvestorIndices = nil
	project.RecipientIndices = append(project.RecipientIndices, recipient1.U.Index, recipient2.U.Index, recipient3.U.Index, recipient4.U.Index, recipient5.U.Index)
	project.DeveloperIndices = append(project.DeveloperIndices, developer1.U.Index, developer2.U.Index)
	project.ContractorFee = 2000
	project.OriginatorFee = 0
	project.DeveloperFee = append(project.DeveloperFee, 6000)
	project.DebtInvestor1 = ""
	project.DebtInvestor2 = ""
	project.TaxEquityInvestor = ""
	// Define things that will be displayed on the frontend
	project.Terms = append(project.Terms, terms1, terms2, terms3, terms4, terms5, terms6)
	project.ExecutiveSummary = esHelper
	project.AutoReloadInterval = -1
	project.ResilienceRating = 0.8
	project.ActionsRequired = ""
	project.Bullets.Bullet1 = "Community owned cooperative mircrogrid"
	project.Bullets.Bullet2 = "Transactive energy capability"
	project.Bullets.Bullet3 = "Certified high social impact"
	var hashHelper opensolar.HashHelper
	project.Hashes = hashHelper
	project.ContractList = nil
	project.CommunityEngagement = append(project.CommunityEngagement, communityblock1, communityblock2, communityblock3, communityblock4)
	project.Architecture.SolarOutputImage = ""
	project.Architecture.SolarArray = "10x 100 W"
	project.Architecture.DailyAvgGeneration = "4000 kWh"
	project.Architecture.System = "600A Deep Cycle"
	project.Architecture.InverterSize = "2024W 230V"
	project.Architecture.DesignDescription = ""
	project.Context = ""
	project.SummaryImage = ""
	project.ExplorePageSummary.Solar = project.PanelSize
	project.ExplorePageSummary.Storage = esHelper.ProjectSize["Storage"]
	project.ExplorePageSummary.Tariff = esHelper.Financials["Tariff"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = esHelper.Financials["Return"]
	project.ExplorePageSummary.Rating = project.Rating
	project.ExplorePageSummary.Tax = "N/A"
	project.ExplorePageSummary.ETA = project.EstimatedAcquisition

	project.DPIntroImage = "https://images.openx.solar/OpenSolarProjects/1_Rwanda/9.jpg"
	project.OHeroImage = ""
	project.OImages = append(project.OImages, "https://images.openx.solar/OpenSolarProjects/1_Rwanda/2.jpg", "https://images.openx.solar/OpenSolarProjects/1_Rwanda/10.jpg")
	project.OOImages = append(project.OOImages, "")
	project.AImages = append(project.AImages, "https://images.openx.solar/OpenSolarProjects/1_Rwanda/5.jpg", "https://images.openx.solar/OpenSolarProjects/1_Rwanda/normal.png")
	project.EImages = append(project.EImages, "https://images.openx.solar/OpenSolarProjects/1_Rwanda/6.jpg")
	project.CEImages = append(project.CEImages, "https://images.openx.solar/OpenSolarProjects/1_Rwanda/8.jpg", "https://images.openx.solar/OpenSolarProjects/1_Rwanda/1.jpg", "https://images.openx.solar/OpenSolarProjects/1_Rwanda/11.jpg", "https://images.openx.solar/OpenSolarProjects/1_Rwanda/7.jpg")
	project.PSImages = append(project.PSImages, "")
	project.BNImages = append(project.BNImages, "")
	project.EngineeringLayoutType = "basic"
	project.FEText = make(map[string]interface{})
	project.FEText, err = parseJsonText("data-sandbox/ubadu.json")
	if err != nil {
		log.Fatal(err)
	}
	err = project.Save()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func testSolarProject(index int, panelsize string, totalValue float64, location string, moneyRaised float64,
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
	project.EstimatedAcquisition = years
	project.RecipientIndex = recpIndex
	project.ContractorIndex = contractor.U.Index
	project.OriginatorIndex = originator.U.Index
	project.Stage = stage
	project.PaybackPeriod = pbperiod
	project.AuctionType = auctionType
	project.InvestmentType = "munibond"

	err := project.Save()
	if err != nil {
		return project, errors.New("Error inserting project into db")
	}
	return project, nil
}

func parseJsonText(fileName string) (map[string]interface{}, error) {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}
