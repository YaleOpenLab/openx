package main

import (
	"github.com/pkg/errors"
	"log"

	database "github.com/YaleOpenLab/openx/database"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	opzones "github.com/YaleOpenLab/openx/platforms/ozones"
	utils "github.com/YaleOpenLab/openx/utils"
)

/*
TEMPLATE FOR A FRONTEND TEST PROJECT
project.Index = //                     int     // an Index to keep track of how many projects exist
project.Name = //                      string  // the name of the project / the identifier by which its referred to
project.State = //                     string  // the state in which the project has been installed in
project.Country = //                   string  // the country in which the project has been installed in
project.TotalValue = //                float64 // the total money that we need from investors
project.PanelSize = //                 string  // size of the given panel, for diplsaying to the user who wants to bid stuff
project.PanelTechnicalDescription = // string  // This should talk about '10x 100W Komaes etc'
project.Inverter = //                  string  // the inverter of the installed project
project.ChargeRegulator = //           string  // the charge regulator of the installed project
project.ControlPanel = //              string  // the control panel of the installed project
project.CommBox = //                   string  // the comm box of the installed project
project.ACTransfer = //                string  // the AC transfer of the installed project
project.SolarCombiner = //             string  // the solar combiner of the installed project
project.Batteries = //                 string  // the batteries of the installed project. TODO: Batteries should also have a fixed nominal value of capacity, as well as one describing what setup it is.
project.IoTHub = //                    string  // the IoT Hub installed as part of the project
project.Metadata = //                  string  // other metadata which does not have an explicit name can be stored here. Used to derive assetIDs

// Define parameters related to finance
project.MoneyRaised = //          float64 // total money that has been raised until now
project.EstimatedAcquisition = // int     // the year in which the recipient is expected to repay the initial investment amount by
project.BalLeft = //              float64 // denotes the balance left to pay by the party, percentage raised is not stored in the database since that can be calculated
project.InterestRate = //         float64 // the rate of return for investors
project.Tax = //                  string  // the specifications of the tax system associated with this particular project

// Define dates of creation and funding
project.DateInitiated = // string // date the project was created on the platform
project.DateFunded = //    string // date that the project completed the stage 4-5 migration
project.DateLastPaid = //  int64  // int64 ie unix time since we need comparisons on this one

// Define technical paramters
project.AuctionType = //          string  // the type of the auction in question. Default is blind auction unless explicitly mentioned
project.InvestmentType = //       string  // the type of investment - equity crowdfunding, municipal bond, normal crowdfunding, etc defined in models
project.PaybackPeriod = //        int     // the frequency in number of weeks that the recipient has to pay the platform.
project.Stage = //                int     // the stage at which the contract is at, float due to potential support of 0.5 state changes in the future
project.SeedInvestmentFactor = // float64 // the factor that a seed investor's investment is multiplied by in case he does invest at the seed stage
project.SeedInvestmentCap = //    float64 // the max amount that a seed investor can put in a project when it is in its seed stages
project.ProposedInvetmentCap = // float64 // the max amount that an investor can invest in when the project is in its proposed stage (stage 2)
project.SelfFund = //             float64 // the amount that a beneficiary / recipient puts in a project wihtout asking from other investors. This is not included as a seed investment because this would mean the recipient pays his own investment back in the project

// Describe issuer of security and the broker dealer
SecurityIssuer string // the issuer of the security
BrokerDealer   string // the broker dealer associated with the project

// Define the various entities that are associated with a specific project
project.RecipientIndex = //              int       // The index of the project's recipient
project.OriginatorIndex = //             int       // the originator of the project
project.GuarantorIndex = //              int       // the person guaranteeing the specific project in question
project.ContractorIndex = //             int       // the person with the proposed contract
project.MainDeveloperIndex = //          int       // the main developer of the project
project.BlendedCapitalInvestorIndex = // int       // the index of the blended capital investor
project.InvestorIndices = //             []int     // The various investors who have invested in the project
project.SeedInvestorIndices = //         []int     // Investors who took part before the contract was at stage 3
project.RecipientIndices = //            []int     // the indices of the recipient family (offtakers, beneficiaries, etc)
project.DeveloperIndices = //            []int     // the indices of the developers involved in the project`
project.ContractorFee = //               float64   // fee paid to the contractor from the total fee of the project
project.OriginatorFee = //               float64   // fee paid to the originator included in the total value of the project
project.DeveloperFee = //                []float64 // the fees charged by the developers
project.DebtInvestor1 = //               string    // debt investor index, if any
project.DebtInvestor2 = //               string    // debt investor index, if any
project.TaxEquityInvestor = //           string    // tax equity investor if any

// Define things that will be displayed on the frontend
project.Terms = //                    []TermsHelper        // the terms of the project
project.InvestmentMetrics = //        InvestmentHelper     // investment metrics that might be useful to an investor
project.FinancialMetrics = //         FinancialHelper      // financial metrics that might be useful to an investor
project.ProjectSizeMetric = //        ProjectSizeHelper    // a metric which shows the size of the project
project.SustainabilityMetric = //     SustainabilityHelper // a metric which shows the sustainability index of the project
project.AutoReloadInterval = //       float64              // the interval in which the user's funds reach zero
project.ResilienceRating = //         float64              // resilience of the project
project.ActionsRequired = //          string               // the action(s) required by the user
project.Bullet1 = //                  string               // bullet points to be displayed on the project summary page
project.Bullet2 = //                  string               // bullet points to be displayed on the project summary page
project.Bullet3 = //                  string               // bullet points to be displayed on the project summary page
project.LegalProjectOverviewHash = // string               // hash to be displayed on the project details page
project.LegalPPAHash = //             string               // hash to be displayed on the project details page
project.LegalRECAgreementHash = //    string               // hash to be displayed on the project details page
project.GuarantorAgreementHash = //   string               // hash to be displayed on the project details page
project.ContractorAgreementHash = //  string               // hash to be displayed on the project details page
project.StakeholderAgreementHash = // string               // hash to be displayed on the project details page
project.CommunityEnergyHash = //      string               // hash to be displayed on the project details page
project.FinancialReportingHash = //   string               // hash to be displayed on the project details page
project.Contract1 = //                string               // contracts which will be linked to on the project details page
project.Contract2 = //                string               // contracts which will be linked to on the project details page
project.Contract3 = //                string               // contracts which will be linked to on the project details page
project.Contract4 = //                string               // contracts which will be linked to on the project details page
project.Contract5 = //                string               // contracts which will be linked to on the project details page

// follow this by the terms and conditions of the specific project
*/
func createPuertoRicoProject() error {
	// setup all the entities that will be involved with the project here
	investor1, err := database.NewInvestor("OpenLab", "p", "x", "Yale OpenLab")
	if err != nil {
		log.Fatal(err)
	}

	recipient1, err := database.NewRecipient("SUpasto", "p", "x", "S.U. Pasto School")
	if err != nil {
		log.Fatal(err)
	}

	originator1, err := opensolar.NewOriginator("DCI", "p", "x", "MIT DCI", "MIT Building E14-15", "The MIT Media Lab's Digital Currency Initiative")
	if err != nil {
		log.Fatal(err)
	}

	contractor1, err := opensolar.NewContractor("MartinWainstein", "p", "x", "Martin Wainstein", "254 Elm Street, New Haven, CT", "Martin Wainstein from the Yale OpenLab")
	if err != nil {
		log.Fatal(err)
	}

	developer1, err := opensolar.NewDeveloper("gs", "p", "x", "Genmoji Solar", "Genmoji, San Juan, Puerto Rico", "Genmoji Solar")
	if err != nil {
		log.Fatal(err)
	}

	developer2, err := opensolar.NewDeveloper("nbly", "p", "x", "Neighborly Securities", "San Francisco, CA", "Broker Dealer")
	if err != nil {
		log.Fatal(err)
	}

	guarantor1, err := opensolar.NewGuarantor("mitml", "p", "x", "MIT Media Lab", "MIT Building E14-15", "The MIT Media Lab is an interdisciplinary lab with innovators from all around the globe")
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
	terms1.Value = "Municipal Bond"
	terms1.RelevantParty = "PR DofEd"
	terms1.Note = "Promoted by PR governor's office"
	terms1.Status = "Demo"
	terms1.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms2 opensolar.TermsHelper
	terms2.Variable = "PPA Tariff"
	terms2.Value = "0.24 ct/KWh"
	terms2.RelevantParty = "oracle X / PREPA"
	terms2.Note = "Variable anchored to local tariff"
	terms2.Status = "Signed"
	terms2.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms3 opensolar.TermsHelper
	terms3.Variable = "Return (TEY)"
	terms3.Value = "3.1%"
	terms3.RelevantParty = "Broker Dealer"
	terms3.Note = "Variable tied to tariff"
	terms3.Status = "Signed"
	terms3.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms4 opensolar.TermsHelper
	terms4.Variable = "Maturity"
	terms4.Value = "+/- 2025"
	terms4.RelevantParty = "Broker Dealer"
	terms4.Note = "Tax adjusted Yield"
	terms4.Status = "Signed"
	terms4.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms5 opensolar.TermsHelper
	terms5.Variable = "Guarantee"
	terms5.Value = "50%"
	terms5.RelevantParty = "Foundation X"
	terms5.Note = "First-loss upon breach"
	terms5.Status = "Started"
	terms5.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms6 opensolar.TermsHelper
	terms6.Variable = "Insurance"
	terms6.Value = "Premium"
	terms6.RelevantParty = "Allianz CS"
	terms6.Note = "Hurricane Coverage"
	terms6.Status = "Started"
	terms6.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var esHelper opensolar.ExecutiveSummaryHelper

	esHelper.Investment.Capex = "5000"
	esHelper.Investment.Hardware = 60
	esHelper.Investment.FirstLossEscrow = "30%"
	esHelper.Investment.CertificationCosts = "N/A"

	esHelper.Financials.Return = 3.1
	esHelper.Financials.Insurance = "Premium"
	esHelper.Financials.Tariff = "0.24 ct/kWh"
	esHelper.Financials.Maturity = "2028 variable"

	esHelper.ProjectSize.PVSolar = "1 kW"
	esHelper.ProjectSize.Storage = "200 Wh"
	esHelper.ProjectSize.Critical = 2
	esHelper.ProjectSize.InverterCapacity = "2.5 kW"

	esHelper.SustainabilityMetrics.CarbonDrawdown = "0.1 t/kWh"
	esHelper.SustainabilityMetrics.CommunityValue = "5/7"
	esHelper.SustainabilityMetrics.LCA = "N/A"

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
	project.Name = "Puerto Rico Pilot Project"
	project.State = "Puerto Rico"
	project.Country = "US"
	project.TotalValue = 10000
	project.PanelSize = "1kW"
	project.PanelTechnicalDescription = "10x 100W Komaes"
	project.Inverter = "Schneider Conext SW 230V 2024"
	project.ChargeRegulator = "Schneider MPPT60"
	project.ControlPanel = "Schneider XW SCP"
	project.CommBox = "Schneider Conext Insight"
	project.ACTransfer = "Eaton Manual throw switches between grid and solar+grid setups"
	project.SolarCombiner = "MidNite"
	project.Batteries = "Advance Autoparts Deep cycle 600A"
	project.IoTHub = "Yale Open Powermeter w/ RaspberryPi3"
	project.Metadata = "This project is a pilot initiative from MIT MediaLab's DCI & the Yale Openlab at Tsai CITY, as to integrate the opensolar platforms with IoT data and blockchain based payment systems to help finance community energy in Puerto Rico"

	// Define parameters related to finance
	project.MoneyRaised = 10000
	project.EstimatedAcquisition = 5
	project.BalLeft = 10000
	project.InterestRate = 0.029
	project.Tax = "Insert tax scheme here"

	// Define dates of creation and funding
	project.DateInitiated = "01/23/2018"
	project.DateFunded = "06/19/2018"
	project.DateLastPaid = -1

	// Define technical paramters
	project.AuctionType = "private"
	project.InvestmentType = "munibond"
	project.PaybackPeriod = 4
	project.Stage = 7
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
	project.InvestorIndices = append(project.InvestorIndices, investor1.U.Index)
	project.SeedInvestorIndices = nil
	project.RecipientIndices = append(project.RecipientIndices, recipient1.U.Index)
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
	project.Bullet1 = "Research proejct on smart financing"
	project.Bullet2 = "Critical loads to the Admin Building"
	project.Bullet3 = "Grid-tied and storage"
	project.LegalProjectOverviewHash = ""
	project.LegalPPAHash = ""
	project.LegalRECAgreementHash = ""
	project.GuarantorAgreementHash = ""
	project.ContractorAgreementHash = ""
	project.StakeholderAgreementHash = ""
	project.CommunityEnergyHash = ""
	project.FinancialReportingHash = ""
	project.Contract1 = ""
	project.Contract2 = ""
	project.Contract3 = ""
	project.Contract4 = ""
	project.Contract5 = ""
	project.CommunityEngagement = append(project.CommunityEngagement, communityblock1, communityblock2, communityblock3, communityblock4)
	project.Architecture.SolarOutputImage = ""
	project.Architecture.SolarArray = "10x 100 W"
	project.Architecture.DailyAvgGeneration = "4000 kWh"
	project.Architecture.System = "600A Deep Cycle"
	project.Architecture.InverterSize = "2024W 230V"
	project.Architecture.DesignDescription = ""
	project.Context = ""
	project.SummaryImage= ""

	err = project.Save()
	if err != nil {
		return err
	}

	return nil
}

func createOneMegaWattProject() error {
	// setup all the entities involved with the project here
	investor1, err := database.NewInvestor("GreenBank", "p", "x", "NH Green Bank")
	if err != nil {
		log.Fatal(err)
	}

	investor2, err := database.NewInvestor("OZFunds", "p", "x", "OZ FundCo")
	if err != nil {
		log.Fatal(err)
	}

	investor3, err := database.NewInvestor("TaxEquity", "p", "x", "Lancaster Lumber Mill Coop")
	if err != nil {
		log.Fatal(err)
	}

	recipient1, err := database.NewRecipient("LancasterHigh", "p", "x", "Lancaster Elementary School")
	if err != nil {
		log.Fatal(err)
	}

	recipient2, err := database.NewRecipient("LancasterT", "p", "x", "Town of Lancaster NH")
	if err != nil {
		log.Fatal(err)
	}

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

	// This is to populate the table of Terms and Conditions in the front end. TODO: change this inline with the FE
	var terms1 opensolar.TermsHelper
	terms1.Variable = "Security Type"
	terms1.Value = "Equity Notes"
	terms1.RelevantParty = "Lancaster Mutual Solar"
	terms1.Note = "Co-owned by the town of Lancaster"
	terms1.Status = "Issued"
	terms1.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms2 opensolar.TermsHelper
	terms2.Variable = "PPA Avg. Tariff"
	terms2.Value = "0.18 ct/KWh"
	terms2.RelevantParty = "Multiple Parties"
	terms2.Note = "Local schools ans business offtaking"
	terms2.Status = "Signed"
	terms2.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms3 opensolar.TermsHelper
	terms3.Variable = "Return (TEY)"
	terms3.Value = "4.8%"
	terms3.RelevantParty = "Broker Dealer"
	terms3.Note = "Tax equivalent yield, with capital gains"
	terms3.Status = "Approv"
	terms3.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms4 opensolar.TermsHelper
	terms4.Variable = "Maturity"
	terms4.Value = "2026"
	terms4.RelevantParty = "Broker Dealer"
	terms4.Note = "By convertible notes"
	terms4.Status = "Signed"
	terms4.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms5 opensolar.TermsHelper
	terms5.Variable = "Guarantee"
	terms5.Value = "50000"
	terms5.RelevantParty = "NH Green Bank"
	terms5.Note = "First-loss escrow upon breach"
	terms5.Status = "Signed"
	terms5.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms6 opensolar.TermsHelper
	terms6.Variable = "Insurance"
	terms6.Value = "Premium"
	terms6.RelevantParty = "Allstate"
	terms6.Note = "Force Majeur"
	terms6.Status = "Signed"
	terms6.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var esHelper opensolar.ExecutiveSummaryHelper

	esHelper.Investment.Capex = "2000000"
	esHelper.Investment.Hardware = 62
	esHelper.Investment.FirstLossEscrow = "50000"
	esHelper.Investment.CertificationCosts = ""

	esHelper.Financials.Return = 4.8
	esHelper.Financials.Insurance = "Premium"
	esHelper.Financials.Tariff = "0.18 ct/kWh"
	esHelper.Financials.Maturity = ""

	esHelper.ProjectSize.PVSolar = "1MW"
	esHelper.ProjectSize.Storage = "210kWh"
	esHelper.ProjectSize.Critical = -1
	esHelper.ProjectSize.InverterCapacity = "1.25MW"

	esHelper.SustainabilityMetrics.CarbonDrawdown = "0.1 t/kWh"
	esHelper.SustainabilityMetrics.CommunityValue = "6/7"
	esHelper.SustainabilityMetrics.LCA = "7/7"

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
	project.Metadata = "Neighborhood 1MW solar array on the field next to Lancaster Elementary High School. The project was originated by the head of the community organization, Ben Southworth, who is also active in the parent teacher association (PTA). The city of Lancaster has agreed to give a 20 year lease of the land to the project if the school gets to own the solar array after the lease expires. The school is located in an opportunity zone"

	// Define parameters related to finance
	project.MoneyRaised = 1500000
	project.EstimatedAcquisition = 20
	project.BalLeft = -1
	project.InterestRate = 0.05
	project.Tax = "Tax free Opportunity Zone"

	// Define dates of creation and funding
	project.DateInitiated = ""
	project.DateFunded = ""
	project.DateLastPaid = -1

	// Define technical paramters
	project.AuctionType = "blind"
	project.InvestmentType = "equity"
	project.PaybackPeriod = 4
	project.Stage = 4
	project.SeedInvestmentFactor = 1.1
	project.SeedInvestmentCap = 500
	project.ProposedInvetmentCap = 15000
	project.SelfFund = 0

	// Describe issuer of security and the broker dealer
	project.SecurityIssuer = "Lancaster Mutual Fund"
	project.BrokerDealer = "Neighborly Securities"

	// Define the various entities that are associated with a specific project
	project.RecipientIndex = recipient1.U.Index
	project.OriginatorIndex = originator1.U.Index
	project.GuarantorIndex = guarantor1.U.Index
	project.ContractorIndex = contractor1.U.Index
	project.MainDeveloperIndex = developer1.U.Index
	project.BlendedCapitalInvestorIndex = investor2.U.Index
	project.InvestorIndices = append(project.InvestorIndices, investor1.U.Index, investor3.U.Index)
	project.SeedInvestorIndices = nil
	project.RecipientIndices = append(project.RecipientIndices, recipient1.U.Index, recipient2.U.Index)
	project.DeveloperIndices = append(project.DeveloperIndices, developer1.U.Index, developer2.U.Index, developer3.U.Index, developer4.U.Index, developer5.U.Index, developer6.U.Index, developer7.U.Index, developer8.U.Index)
	project.ContractorFee = 2000
	project.OriginatorFee = 0
	project.DeveloperFee = append(project.DeveloperFee, 6000)
	project.DebtInvestor1 = "OZFunds"
	project.DebtInvestor2 = "GreenBank"
	project.TaxEquityInvestor = "TaxEquity"

	// Define things that will be displayed on the frontend
	project.Terms = append(project.Terms, terms1, terms2, terms3, terms4, terms5, terms6)
	project.ExecutiveSummary = esHelper
	project.AutoReloadInterval = -1
	project.ResilienceRating = 0.8
	project.ActionsRequired = ""
	project.Bullet1 = "Eligible for Opportunity Zone Investments"
	project.Bullet2 = "Community owned mutual funds as bond issuer"
	project.Bullet3 = "1 MW grid with full offtake agreements"
	project.LegalProjectOverviewHash = ""
	project.LegalPPAHash = ""
	project.LegalRECAgreementHash = ""
	project.GuarantorAgreementHash = ""
	project.ContractorAgreementHash = ""
	project.StakeholderAgreementHash = ""
	project.CommunityEnergyHash = ""
	project.FinancialReportingHash = ""
	project.Contract1 = ""
	project.Contract2 = ""
	project.Contract3 = ""
	project.Contract4 = ""
	project.Contract5 = ""
	project.CommunityEngagement = append(project.CommunityEngagement, communityblock1)
	project.Architecture.SolarOutputImage = ""
	project.Architecture.SolarArray = "1000 kW"
	project.Architecture.DailyAvgGeneration = "4000 kWh"
	project.Architecture.System = "Battery + Grid"
	project.Architecture.InverterSize = "1.25MW"
	project.Architecture.DesignDescription = ""
	project.Context = ""
	project.SummaryImage= ""

	err = project.Save()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func createTenKiloWattProject() error {

	investor1, err := database.NewInvestor("MatthewMoroney", "p", "x", "Matthew Moroney")
	if err != nil {
		log.Fatal(err)
	}

	investor2, err := database.NewInvestor("FranzHochstrasser", "p", "x", "Franz Hochstrasser")
	if err != nil {
		log.Fatal(err)
	}

	investor3, err := database.NewInvestor("CTGreenBank", "p", "x", "Connecticut Green Bank")
	if err != nil {
		log.Fatal(err)
	}

	investor4, err := database.NewInvestor("YaleUniversity", "p", "x", "Yale University Community Fund")
	if err != nil {
		log.Fatal(err)
	}

	investor5, err := database.NewInvestor("JeromeGreen", "p", "x", "Jerome Green")
	if err != nil {
		log.Fatal(err)
	}

	investor6, err := database.NewInvestor("OpenSolarFund", "p", "x", "Open Solar Revolving Fund")
	if err != nil {
		log.Fatal(err)
	}

	recipient1, err := database.NewRecipient("Shelter1", "p", "x", "Shelter1 Community Solar")
	if err != nil {
		log.Fatal(err)
	}

	recipient2, err := database.NewRecipient("ColumbusHouse", "p", "x", "Columbus House Foundation")
	if err != nil {
		log.Fatal(err)
	}

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

	esHelper.Investment.Capex = "30000"
	esHelper.Investment.Hardware = 3
	esHelper.Investment.FirstLossEscrow = "N/A"
	esHelper.Investment.CertificationCosts = "N/A"

	esHelper.Financials.Return = -1
	esHelper.Financials.Insurance = "Basic"
	esHelper.Financials.Tariff = "0.12 ct/kWh"
	esHelper.Financials.Maturity = "2019"

	esHelper.ProjectSize.PVSolar = "10 kW"
	esHelper.ProjectSize.Storage = "N/A Grid Tied"
	esHelper.ProjectSize.Critical = 100
	esHelper.ProjectSize.InverterCapacity = "15 kW"

	esHelper.SustainabilityMetrics.CarbonDrawdown = "0.1 t/kWh"
	esHelper.SustainabilityMetrics.CommunityValue = "7/7"
	esHelper.SustainabilityMetrics.LCA = ""

	var communityblock1 opensolar.CommunityEngagementHelper
	communityblock1.Width = 12
	communityblock1.Title = "Consultation"
	communityblock1.ImageURL = ""
	communityblock1.Content = ""
	communityblock1.Link = ""

	// This is to populate the table of Terms and Conditions in the front end. TODO: change this inline with the FE
	var terms1 opensolar.TermsHelper
	terms1.Variable = "Security Type"
	terms1.Value = "Reg CF"
	terms1.RelevantParty = "NH Community Solar"
	terms1.Note = "Special Purpose Vehicle"
	terms1.Status = "Flipped"
	terms1.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms2 opensolar.TermsHelper
	terms2.Variable = "PPA Tariff"
	terms2.Value = "0.12 ct/KWh"
	terms2.RelevantParty = "NH Homeless Shelter"
	terms2.Note = "Fixed PPA determined by offtaker"
	terms2.Status = "Signed"
	terms2.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms3 opensolar.TermsHelper
	terms3.Variable = "Return"
	terms3.Value = "130%"
	terms3.RelevantParty = "Equity Value"
	terms3.Note = "Growth in value. No tax incentives"
	terms3.Status = "Open"
	terms3.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms4 opensolar.TermsHelper
	terms4.Variable = "Ownership Flip"
	terms4.Value = "2019"
	terms4.RelevantParty = "Convertible Note"
	terms4.Note = "Crowd investors sell stock"
	terms4.Status = "Flipped"
	terms4.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms5 opensolar.TermsHelper
	terms5.Variable = "Guarantee"
	terms5.Value = "N/A"
	terms5.RelevantParty = "N/A"
	terms5.Note = "No guarantees of breach"
	terms5.Status = "None"
	terms5.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms6 opensolar.TermsHelper
	terms6.Variable = "Insurance"
	terms6.Value = "Basic"
	terms6.RelevantParty = "CT Insurers"
	terms6.Note = "Force Majeur"
	terms6.Status = "Signed"
	terms6.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	project.Index = len(indexHelp) + 1
	project.Name = "Homeless Shelter Construction"
	project.State = "Puerto Rico"
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
	project.Metadata = "Residential solar array for a homeless shelter. The project was originated by a member of the board of the homeless shelter who gets the shelter to purchase all the electricity at a discounted rate. The shelter chooses to lease the roof for free over the lifetime of the project. The originator knows the solar developer who set up the project company"

	// Define parameters related to finance
	project.MoneyRaised = 30000
	project.EstimatedAcquisition = 0 // this project already flipped ownership
	project.BalLeft = 0
	project.InterestRate = 0.05
	project.Tax = "0.3 Tax Credit"

	// Define dates of creation and funding
	project.DateInitiated = ""
	project.DateFunded = ""
	project.DateLastPaid = -1

	// Define technical paramters
	project.AuctionType = "blind"
	project.InvestmentType = "regcf"
	project.PaybackPeriod = 4
	project.Stage = 8
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
	project.GuarantorIndex = guarantor1.U.Index
	project.ContractorIndex = contractor1.U.Index
	project.MainDeveloperIndex = developer1.U.Index
	project.BlendedCapitalInvestorIndex = -1
	project.InvestorIndices = append(project.InvestorIndices, investor1.U.Index, investor2.U.Index, investor3.U.Index, investor4.U.Index, investor5.U.Index, investor6.U.Index)
	project.SeedInvestorIndices = nil
	project.RecipientIndices = append(project.RecipientIndices, recipient1.U.Index, recipient2.U.Index)
	project.DeveloperIndices = append(project.DeveloperIndices, developer1.U.Index, developer2.U.Index, developer3.U.Index, developer4.U.Index, developer5.U.Index)
	project.ContractorFee = 0
	project.OriginatorFee = 0
	project.DeveloperFee = append(project.DeveloperFee, 0)
	project.DebtInvestor1 = ""
	project.DebtInvestor2 = ""
	project.TaxEquityInvestor = ""

	// Define things that will be displayed on the frontend
	project.Terms = append(project.Terms, terms1, terms2, terms3, terms4, terms5, terms6)
	project.ExecutiveSummary = esHelper
	project.AutoReloadInterval = -1
	project.ResilienceRating = 0.6
	project.ActionsRequired = ""
	project.Bullet1 = "Community owned solar in homeless shelter"
	project.Bullet2 = "Siginificantly alleviates financial pressure due to high CT power cost"
	project.Bullet3 = "Grid-tied with REC offtaking"
	project.LegalProjectOverviewHash = ""
	project.LegalPPAHash = ""
	project.LegalRECAgreementHash = ""
	project.GuarantorAgreementHash = ""
	project.ContractorAgreementHash = ""
	project.StakeholderAgreementHash = ""
	project.CommunityEnergyHash = ""
	project.FinancialReportingHash = ""
	project.Contract1 = ""
	project.Contract2 = ""
	project.Contract3 = ""
	project.Contract4 = ""
	project.Contract5 = ""
	project.CommunityEngagement = append(project.CommunityEngagement, communityblock1)
	project.Architecture.SolarOutputImage = ""
	project.Architecture.SolarArray = "50 x 200W"
	project.Architecture.DailyAvgGeneration = "20000 Wh"
	project.Architecture.System = "Grid Tied"
	project.Architecture.InverterSize = "15 kW"
	project.Architecture.DesignDescription = ""
	project.Context = ""
	project.SummaryImage= ""

	err = project.Save()
	if err != nil {
		log.Fatal(err)
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

	// This is to populate the table of Terms and Conditions in the front end. TODO: change this inline with the FE
	var terms1 opensolar.TermsHelper
	terms1.Variable = "Security Type"
	terms1.Value = "Municipal Bond"
	terms1.RelevantParty = "PR DofEd"
	terms1.Note = "Not yet issued. See informal Agreements"
	terms1.Status = "Open"
	terms1.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms2 opensolar.TermsHelper
	terms2.Variable = "PPA Tariff"
	terms2.Value = "0.24 ct/KWh"
	terms2.RelevantParty = "Oracle X / PREPA"
	terms2.Note = "Not signed. Expected as variable tariff"
	terms2.Status = "Open"
	terms2.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms3 opensolar.TermsHelper
	terms3.Variable = "Return (TEY)"
	terms3.Value = "3.5%"
	terms3.RelevantParty = "See Broker Dealer"
	terms3.Note = "Tax equivalent yield with capital gains"
	terms3.Status = "Open"
	terms3.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms4 opensolar.TermsHelper
	terms4.Variable = "Maturity"
	terms4.Value = "+/- 2025"
	terms4.RelevantParty = "Broker Dealer"
	terms4.Note = "Variable tied to tariff"
	terms4.Status = "Open"
	terms4.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms5 opensolar.TermsHelper
	terms5.Variable = "Guarantee"
	terms5.Value = "15%"
	terms5.RelevantParty = "FEMA"
	terms5.Note = "First-loss upon breach"
	terms5.Status = "Started"
	terms5.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms6 opensolar.TermsHelper
	terms6.Variable = "Insurance"
	terms6.Value = "Premium"
	terms6.RelevantParty = "Allianz CS"
	terms6.Note = "Hurricane Coverage"
	terms6.Status = "Started"
	terms6.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var esHelper opensolar.ExecutiveSummaryHelper

	esHelper.Investment.Capex = "19000000"
	esHelper.Investment.Hardware = 70
	esHelper.Investment.FirstLossEscrow = "15%"
	esHelper.Investment.CertificationCosts = ""

	esHelper.Financials.Return = 2.5
	esHelper.Financials.Insurance = "Basic Force Majeur"
	esHelper.Financials.Tariff = "0.24 ct/kWh"
	esHelper.Financials.Maturity = "2028"

	esHelper.ProjectSize.PVSolar = "300 x 30kW"
	esHelper.ProjectSize.Storage = "350 x 2.5 kWh"
	esHelper.ProjectSize.Critical = 20
	esHelper.ProjectSize.InverterCapacity = "300 x 35 kW"

	esHelper.SustainabilityMetrics.CarbonDrawdown = "0.1t/kWh"
	esHelper.SustainabilityMetrics.CommunityValue = "6/7"
	esHelper.SustainabilityMetrics.LCA = "N/A"

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
	project.Name = "Puerto Rico Public School Bond 10"
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

	// Define technical paramters
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
	project.Bullet1 = "Backed by the Governor's office"
	project.Bullet2 = "Critical Loads covered (telecom and refrigeration)"
	project.Bullet3 = "Certification of social impact"
	project.LegalProjectOverviewHash = ""
	project.LegalPPAHash = ""
	project.LegalRECAgreementHash = ""
	project.GuarantorAgreementHash = ""
	project.ContractorAgreementHash = ""
	project.StakeholderAgreementHash = ""
	project.CommunityEnergyHash = ""
	project.FinancialReportingHash = ""
	project.Contract1 = ""
	project.Contract2 = ""
	project.Contract3 = ""
	project.Contract4 = ""
	project.Contract5 = ""
	project.CommunityEngagement = append(project.CommunityEngagement, communityblock1, communityblock2, communityblock3, communityblock4)
	project.Architecture.SolarOutputImage = ""
	project.Architecture.SolarArray = "300 x 30kW"
	project.Architecture.DailyAvgGeneration = "400 MWh"
	project.Architecture.System = "350 Tesla Powerwells"
	project.Architecture.InverterSize = "300 x 35 kW"
	project.Architecture.DesignDescription = ""
	project.Context = ""
	project.SummaryImage= ""

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

	// This is to populate the table of Terms and Conditions in the front end. TODO: change this inline with the FE
	var terms1 opensolar.TermsHelper
	terms1.Variable = "Security Type"
	terms1.Value = "Municipal Bond"
	terms1.RelevantParty = "PR DofEd"
	terms1.Note = "Promoted by PR governor's office"
	terms1.Status = "Demo"
	terms1.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms2 opensolar.TermsHelper
	terms2.Variable = "PPA Tariff"
	terms2.Value = "0.24 ct/KWh"
	terms2.RelevantParty = "oracle X / PREPA"
	terms2.Note = "Variable anchored to local tariff"
	terms2.Status = "Signed"
	terms2.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms3 opensolar.TermsHelper
	terms3.Variable = "Return (TEY)"
	terms3.Value = "3.1%"
	terms3.RelevantParty = "Broker Dealer"
	terms3.Note = "Variable tied to tariff"
	terms3.Status = "Signed"
	terms3.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms4 opensolar.TermsHelper
	terms4.Variable = "Maturity"
	terms4.Value = "+/- 2025"
	terms4.RelevantParty = "Broker Dealer"
	terms4.Note = "Tax adjusted Yield"
	terms4.Status = "Signed"
	terms4.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms5 opensolar.TermsHelper
	terms5.Variable = "Guarantee"
	terms5.Value = "50%"
	terms5.RelevantParty = "Foundation X"
	terms5.Note = "First-loss upon breach"
	terms5.Status = "Started"
	terms5.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var terms6 opensolar.TermsHelper
	terms6.Variable = "Insurance"
	terms6.Value = "Premium"
	terms6.RelevantParty = "Allianz CS"
	terms6.Note = "Hurricane Coverage"
	terms6.Status = "Started"
	terms6.SupportDoc = "https://openlab.yale.edu" // replace this with the relevant doc

	var investmentHelper opensolar.InvestmentHelper
	var financialHelper opensolar.FinancialHelper
	var projectSizeHelper opensolar.ProjectSizeHelper
	var sustainabilityHelper opensolar.SustainabilityHelper

	investmentHelper.Capex = ""
	investmentHelper.Hardware = -1
	investmentHelper.FirstLossEscrow = ""
	investmentHelper.CertificationCosts = ""

	financialHelper.Return = -1
	financialHelper.Insurance = ""
	financialHelper.Tariff = ""
	financialHelper.Maturity = ""

	projectSizeHelper.PVSolar = ""
	projectSizeHelper.Storage = ""
	projectSizeHelper.Critical = -1
	projectSizeHelper.InverterCapacity = ""

	sustainabilityHelper.CarbonDrawdown = ""
	sustainabilityHelper.CommunityValue = ""
	sustainabilityHelper.LCA = ""

	project.Index = len(indexHelp) + 1
	project.Name = "Rwanda Community Microgrid"
	project.PanelSize = "100kW"
	project.State = "Khigali"
	project.Country = "Rwanda"
	project.Batteries = "25 kWh"
	project.Stage = 1
	project.DateInitiated = "03/25/2019"
	project.Metadata = "The community of Ubadu, Rwanda has no access to electricity yet shows a growing local economy. This microgrid project, developed a collaboration with Yale and MIT, aims to serve 250 homes, including its only school, ‘Sunshine Garden,’ the town infirmary led by a team of doctors without borders, and the town hall. Community cooperative with international backing. 20% first loss fund secured. Currently doing engineering due diligence for development quotes"
	project.TotalValue = 230000
	project.SeedInvestmentCap = 5000
	project.MoneyRaised = 1250
	project.InterestRate = 0.023
	project.EstimatedAcquisition = 7
	project.PaybackPeriod = 4
	project.InvestmentType = "equity"
	project.InvestorIndices = append(project.InvestorIndices, investor1.U.Index, investor2.U.Index, investor3.U.Index)
	project.RecipientIndices = append(project.RecipientIndices, recipient1.U.Index, recipient2.U.Index, recipient3.U.Index, recipient4.U.Index, recipient5.U.Index)
	project.DeveloperIndices = append(project.DeveloperIndices, developer1.U.Index, developer2.U.Index)
	project.OriginatorIndex = originator1.U.Index
	project.Terms = append(project.Terms, terms1, terms2, terms3, terms4, terms5, terms6)

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

// ALL 5 PROJECT DATA WILL BE ADDED HERE FOR THE DEMO
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

	_, err = testSolarProject(1, "100 1000 sq.ft homes each with their own private spaces for luxury", 14000, "India Basin, San Francisco",
		0, "India Basin is an upcoming creative project based in San Francisco that seeks to invite innovators from all around to participate", "", "", "",
		3, recp.U.Index, contractor, originator, 4, 2, "blind")

	if err != nil {
		log.Fatal(err)
	}

	_, err = testSolarProject(2, "180 1200 sq.ft homes in a high rise building 0.1mi from Kendall Square", 30000, "Kendall Square, Boston",
		0, "Kendall Square is set in the heart of Cambridge and is a popular startup IT hub", "", "", "",
		5, recp.U.Index, contractor, originator, 4, 2, "blind")

	if err != nil {
		log.Fatal(err)
	}

	_, err = testSolarProject(3, "260 1500 sq.ft homes set in a medieval cathedral style construction", 40000, "Trafalgar Square, London",
		0, "Trafalgar Square is set in the heart of London's financial district, with big banks all over", "", "", "",
		7, recp.U.Index, contractor, originator, 4, 2, "blind")

	if err != nil {
		log.Fatal(err)
	}

	_, err = originator.Originate("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", 1, "blind") // 1 is the idnex for martin
	if err != nil {
		log.Fatal(err)
	}

	// project: Puerto Rico Project
	// STAGE 7 - Puerto Rico
	err = createPuertoRicoProject()
	if err != nil {
		log.Fatal(err)
	}
	// project: One Mega Watt Project
	// STAGE 4 - New Hampshire
	err = createOneMegaWattProject()
	if err != nil {
		log.Fatal(err)
	}
	// project: Ten Kilowatt Project
	// STAGE 8 - Connecticut Homeless Shelter
	err = createTenKiloWattProject()
	if err != nil {
		log.Fatal(err)
	}
	// project: Ten Mega Watt Project
	// STAGE 2 - Puerto Rico Public School Bond
	err = createTenMegaWattProject()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
