package main

import (
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

func populateStaticData1kw() (int, error) {

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

func create1kwProject() error {
	// setup all the entities that will be involved with the project here
	projIndex, err := populateStaticData1kw()
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

	project, err := opensolar.RetrieveProject(6)
	if err != nil {
		return -1, err
	}
	// Define things that will be displayed on the frontend
	project.AutoReloadInterval = -1
	project.ResilienceRating = 0.8
	project.ContractList = nil
	project.Architecture.SolarOutputImage = ""
	project.Architecture.SolarArray = "1000 kW"
	project.Architecture.DailyAvgGeneration = "4000 kWh"
	project.Architecture.System = "Battery + Grid"
	project.Architecture.InverterSize = "1.25MW"
	project.Architecture.DesignDescription = ""
	project.Context = ""
	project.SummaryImage = ""
	project.ExplorePageSummary.Solar = project.PanelSize
	project.ExplorePageSummary.Storage = project.ExecutiveSummary.ProjectSize["Storage"]
	project.ExplorePageSummary.Tariff = project.ExecutiveSummary.Financials["Tariff (Variable)"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = project.ExecutiveSummary.Financials["Return (TEY)"]
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

func create1mwProject() error {
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

	project, err := opensolar.RetrieveProject(7)
	if err != nil {
		return -1, err
	}

	project.ExplorePageSummary.Solar = project.PanelSize
	project.ExplorePageSummary.Storage = project.ExecutiveSummary.ProjectSize["Storage"]
	project.ExplorePageSummary.Tariff = project.ExecutiveSummary.Financials["Tariff (Fixed)"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = project.ExecutiveSummary.Financials["Return (TEY)"]
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

func create10kwProject() error {

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

func create10mwProject() error {
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

	project, err := opensolar.RetrieveProject(8)
	if err != nil {
		return err
	}
	project.MoneyRaised = 0

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

	project.ResilienceRating = 0.75
	project.ExplorePageSummary.Solar = project.PanelSize
	project.ExplorePageSummary.Storage = project.ExecutiveSummary.ProjectSize["Storage"]
	project.ExplorePageSummary.Tariff = project.ExecutiveSummary.Financials["Tariff (Variable)"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = project.ExecutiveSummary.Financials["Expected Return (Non TEY)"]
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

	project, err := opensolar.RetrieveProject(9)
	if err != nil {
		return err
	}
	project.MoneyRaised = 0
	// Define the various entities that are associated with a specific project
	project.RecipientIndex = recipient1.U.Index
	project.OriginatorIndex = originator1.U.Index
	// project.GuarantorIndex = guarantor1.U.Index
	// project.ContractorIndex = contractor1.U.Index
	project.MainDeveloperIndex = developer1.U.Index
	project.BlendedCapitalInvestorIndex = -1
	project.InvestorIndices = append(project.InvestorIndices, investor1.U.Index, investor2.U.Index, investor3.U.Index)
	project.RecipientIndices = append(project.RecipientIndices, recipient1.U.Index, recipient2.U.Index, recipient3.U.Index, recipient4.U.Index, recipient5.U.Index)
	project.DeveloperIndices = append(project.DeveloperIndices, developer1.U.Index, developer2.U.Index)
	// Define things that will be displayed on the frontend
	project.ResilienceRating = 0.8
	project.ExplorePageSummary.Solar = project.PanelSize
	project.ExplorePageSummary.Storage = project.ExecutiveSummary.ProjectSize["Storage"]
	project.ExplorePageSummary.Tariff = project.ExecutiveSummary.Financials["Tariff"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = project.ExecutiveSummary.Financials["Return"]
	project.ExplorePageSummary.Rating = project.Rating
	project.ExplorePageSummary.Tax = "N/A"
	project.ExplorePageSummary.ETA = project.EstimatedAcquisition

	project.DPIntroImage = "https://images.openx.solar/OpenSolarProjects/1_Rwanda/9.jpg"
	project.OImages = append(project.OImages, "https://images.openx.solar/OpenSolarProjects/1_Rwanda/2.jpg", "https://images.openx.solar/OpenSolarProjects/1_Rwanda/10.jpg")
	project.AImages = append(project.AImages, "https://images.openx.solar/OpenSolarProjects/1_Rwanda/5.jpg", "https://images.openx.solar/OpenSolarProjects/1_Rwanda/normal.png")
	project.EImages = append(project.EImages, "https://images.openx.solar/OpenSolarProjects/1_Rwanda/6.jpg")
	project.CEImages = append(project.CEImages, "https://images.openx.solar/OpenSolarProjects/1_Rwanda/8.jpg", "https://images.openx.solar/OpenSolarProjects/1_Rwanda/1.jpg", "https://images.openx.solar/OpenSolarProjects/1_Rwanda/11.jpg", "https://images.openx.solar/OpenSolarProjects/1_Rwanda/7.jpg")
	err = project.Save()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
