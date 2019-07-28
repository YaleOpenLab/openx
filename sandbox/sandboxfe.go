package sandbox

// fe sandbox contains data related to the sandbox that is displayed on the frontend. Can be updated to
// reflect changes on the frontend. File last updated: May 2019
import (
	"log"
	"time"

	xlm "github.com/Varunram/essentials/crypto/xlm"
	assets "github.com/Varunram/essentials/crypto/xlm/assets"
	wallet "github.com/Varunram/essentials/crypto/xlm/wallet"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
)

func createAllStaticEntities() error {
	// PR static entities
	var err error
	_, err = opensolar.NewOriginator("dci@test.com", "password", "x", "MIT DCI", "MIT Building E14-15", "The MIT Media Lab's Digital Currency Initiative")
	if err != nil {
		return err
	}

	_, err = opensolar.NewContractor("martinwainstein@test.com", "password", "x", "Martin Wainstein", "254 Elm Street, New Haven, CT", "Martin Wainstein from the Yale OpenLab")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("gs@test.com", "password", "x", "Genmoji Solar", "Genmoji, San Juan, Puerto Rico", "Genmoji Solar")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("nbly@test.com", "password", "x", "Neighborly Securities", "San Francisco, CA", "Broker Dealer")
	if err != nil {
		return err
	}

	_, err = opensolar.NewGuarantor("mitml@test.com", "password", "x", "MIT Media Lab", "MIT Building E14-15", "The MIT Media Lab is an interdisciplinary lab with innovators from all around the globe")
	if err != nil {
		return err
	}

	// 1MW project (5)
	_, err = opensolar.NewContractor("testcont@test.com", "password", "x", "testcont", "testcont", "testcont")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("solardev@test.com", "password", "x", "First Solar", "Solar Rd, San Diego, California", "Main contractor for full solar development")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("LancasterSolar@test.com", "password", "x", "Town of Lancaste NH", "Lancaster, New Hampshire", "Host")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("LancasterRFP@test.com", "password", "x", "Lancaster Solar Engineer Solutions", "25 Lancaster Rd, New Hampshire", "Independent RFP Engineer")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("SimpleServiceProvider@test.com", "password", "x", "Simple Service Provider", "Simple Service Provider", "Simple Service Provider")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("VendorX@test.com", "password", "x", "Solar Racking Systems Inc", "34 Crack St, Boston", "Retail Vendor")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("NEpool@test.com", "password", "x", "New England Pool Registered Auditor", "56 Hamden Ave, Stamford, CT", "REC Auditors for New England")
	if err != nil {
		return err
	}

	_, err = opensolar.NewGuarantor("AllianzCS@test.com", "password", "x", "Allianz Climate Solutions", "34 5th, New York, NY", "Insurance Agent")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("UIavangrid@test.com", "password", "x", "Avangrid Networks", "100 Marsh Hill Rd, New Haven, CT", "Utility")
	if err != nil {
		return err
	}

	_, err = opensolar.NewGuarantor("GreenBank@test.com", "password", "x", "NH Green Bank", "67 Washington Rd, New Hampshire", "Impact-first escrow provider")
	if err != nil {
		return err
	}

	_, err = opensolar.NewOriginator("ben@test.com", "password", "x", "Ben Southworth", "Lancaster, NH", "Originator of the Lancaster oz fund community")
	if err != nil {
		return err
	}

	// 10kW project (16)
	_, err = opensolar.NewDeveloper("YaleArchitecture@test.com", "password", "x", "Yale School of Architecture", "45 York St, New Haven, CT", "System and layout designer")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("CTSolar@test.com", "password", "x", "Connecticut Solar", "45 Sun Street, Stamford, CT", "Solar system installer")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("ColumbusHouse@test.com", "password", "x", "Columbus House", "21 Hagrid Ave, New Haven, CT", "Project Host")
	if err != nil {
		return err
	}

	_, err = opensolar.NewGuarantor("RGreenFund@test.com", "password", "x", "RaiseGreen Blend Fund", "21 orange st, New Haven, CT", "Impact-first blended capital provider")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("Avangrid@test.com", "password", "x", "Avangrid RECs", "100 Marsh Hill Rd, New Haven, CT", "Certifier of RECs and provider of REC meter")
	if err != nil {
		return err
	}

	_, err = opensolar.NewOriginator("RaiseGreen@test.com", "password", "x", "Raise Green", "21 orange st, New Haven, CT", "Project originator")
	if err != nil {
		return err
	}

	_, err = opensolar.NewContractor("testcont@test.com", "password", "x", "testcont", "testcont", "testcont")
	if err != nil {
		return err
	}

	_, err = opensolar.NewGuarantor("testguarantor@test.com", "password", "x", "testguarantor", "testguarantor", "testguarantor")
	if err != nil {
		return err
	}

	// 10MW project (24)
	_, err = database.NewInvestor("emcoll@test.com", "password", "x", "Emerson Collective")
	if err != nil {
		return err
	}

	_, err = database.NewInvestor("prqozfund@test.com", "password", "x", "Puerto Rico QOZ Fund")
	if err != nil {
		return err
	}

	_, err = database.NewRecipient("prgov@test.com", "password", "x", "PR Government")
	if err != nil {
		return err
	}

	_, err = database.NewRecipient("prschools@test.com", "password", "x", "Puerto Rico Solar Schools Limited")
	if err != nil {
		return err
	}

	_, err = database.NewRecipient("prdoe@test.com", "password", "x", "Puerto Rico Department of Education")
	if err != nil {
		return err
	}

	_, err = opensolar.NewOriginator("MartinWainstein1@test.com", "password", "x", "Martin Wainstein", "254 Elm Street, New Haven, CT", "Martin Wainstein from the Yale OpenLab")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("hst@test.com", "password", "x", "HST Solar", "25 Hewlett St, San Francisco, CA", "Preliminary finance and engineering assessment")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("FemaRoofs@test.com", "password", "x", "FEMA Puerto Rico", "“45 Old Town Rd, Puerto Rico", "Civil engineering assessment of school roofs")
	if err != nil {
		return err
	}

	_, err = opensolar.NewContractor("testcont@test.com", "password", "x", "testcont", "testcont", "testcont")
	if err != nil {
		return err
	}

	_, err = opensolar.NewGuarantor("testguarantor@test.com", "password", "x", "testguarantor", "testguarantor", "testguarantor")
	if err != nil {
		return err
	}

	// 100 kw project (34)
	_, err = database.NewInvestor("jjackson@test.com", "password", "x", "Jerome Jackson")
	if err != nil {
		return err
	}

	_, err = database.NewInvestor("esare@test.com", "password", "x", "Eliah Sare")
	if err != nil {
		return err
	}

	_, err = database.NewInvestor("yaleuf@test.com", "password", "x", "Yale University Fund")
	if err != nil {
		return err
	}

	_, err = database.NewRecipient("villageec@test.com", "password", "x", "Village Energy Collective")
	if err != nil {
		return err
	}

	_, err = database.NewRecipient("sunshinegschool@test.com", "password", "x", "Sunshine Garden School")
	if err != nil {
		return err
	}

	_, err = database.NewRecipient("ubaduth@test.com", "password", "x", "Ubadu Town Hall")
	if err != nil {
		return err
	}

	_, err = database.NewRecipient("dwbrf@test.com", "password", "x", " Doctors without borders, Rwanda chapter")
	if err != nil {
		return err
	}

	_, err = database.NewRecipient("largerof@test.com", "password", "x", "Large Residential offtakers")
	if err != nil {
		return err
	}

	// insert small residential offtakers as well
	_, err = opensolar.NewOriginator("embeba@test.com", "password", "x", "School Principa", "Village, Rwanda", "Project originator")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("solarpartners@test.com", "password", "x", "Solar Partners", "KG 10, House 25 Gasabo district, kamatamuUrugwiro, Kacyiru, Kigali", "MiniGrid game developer")
	if err != nil {
		return err
	}

	_, err = opensolar.NewDeveloper("hst2@test.com", "password", "x", "HST Solar", "25 Hewlett St, San Francisco, CA", "Preliminary finance and engineering assessment")
	if err != nil {
		return err
	}

	return nil
}

func populateStaticData1kw() error {

	project, err := opensolar.RetrieveProject(4)
	if err != nil {
		return err
	}

	project.ExplorePageSummary.Solar = project.PanelSize
	project.ExplorePageSummary.Storage = project.ExecutiveSummary.ProjectSize["storage"]
	project.ExplorePageSummary.Tariff = project.ExecutiveSummary.Financials["tariff (variable)"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = project.ExecutiveSummary.Financials["return(tey)"]
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
		return err
	}

	return nil
}

func populateStaticData1mw() error {
	project, err := opensolar.RetrieveProject(5)
	if err != nil {
		log.Println("bug?")
		return err
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
	project.ExplorePageSummary.Storage = project.ExecutiveSummary.ProjectSize["storage"]
	project.ExplorePageSummary.Tariff = project.ExecutiveSummary.Financials["tariff (variable)"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = project.ExecutiveSummary.Financials["return (tey)"]
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

	project.OriginatorIndex = 16
	project.GuarantorIndex = 15
	project.ContractorIndex = 6
	project.MainDeveloperIndex = 7
	project.DeveloperIndices = append(project.DeveloperIndices, 7, 8, 9, 10, 11, 12, 13, 14)
	project.DebtInvestor1 = "OZFunds"
	project.DebtInvestor2 = "GreenBank"
	project.TaxEquityInvestor = "TaxEquity"
	err = project.Save()
	if err != nil {
		return err
	}
	return nil
}

func populateStaticData10kw() error {
	project, err := opensolar.RetrieveProject(6)
	if err != nil {
		return err
	}

	project.ExplorePageSummary.Solar = project.PanelSize
	project.ExplorePageSummary.Storage = project.ExecutiveSummary.ProjectSize["storage"]
	project.ExplorePageSummary.Tariff = project.ExecutiveSummary.Financials["tariff(fixed)"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = "3%" // since we don't have that in the exec page summary, we hardcode that here
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

	project.OriginatorIndex = 22
	project.GuarantorIndex = 24
	project.ContractorIndex = 23
	project.MainDeveloperIndex = 17
	project.DeveloperIndices = append(project.DeveloperIndices, 17, 18, 19, 20, 21)

	err = project.Save()
	if err != nil {
		return err
	}

	return nil
}

func populateStaticData10MW() error {
	// create the required entities that we need over here
	project, err := opensolar.RetrieveProject(7)
	if err != nil {
		return err
	}
	project.MoneyRaised = 0

	// Define the various entities that are associated with a specific project
	project.RecipientIndex = 27
	project.OriginatorIndex = 30
	project.GuarantorIndex = 34
	project.ContractorIndex = 33
	project.MainDeveloperIndex = 31
	project.BlendedCapitalInvestorIndex = -1
	project.InvestorIndices = append(project.InvestorIndices, 25, 26)
	project.SeedInvestorIndices = nil
	project.RecipientIndices = append(project.RecipientIndices, 27, 28, 29)
	project.DeveloperIndices = append(project.DeveloperIndices, 31, 32)

	project.ResilienceRating = 0.75
	project.ExplorePageSummary.Solar = project.PanelSize
	project.ExplorePageSummary.Storage = project.ExecutiveSummary.ProjectSize["storage"]
	project.ExplorePageSummary.Tariff = project.ExecutiveSummary.Financials["tariff (variable)"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = project.ExecutiveSummary.Financials["expected return (non tey)"]
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
		return err
	}

	return nil
}

func populateStaticData100KW() error {

	project, err := opensolar.RetrieveProject(8)
	if err != nil {
		return err
	}
	project.MoneyRaised = 0
	// Define the various entities that are associated with a specific project
	project.RecipientIndex = 38
	project.OriginatorIndex = 43
	// project.GuarantorIndex = guarantor1.U.Index
	// project.ContractorIndex = contractor1.U.Index
	project.MainDeveloperIndex = 44
	project.BlendedCapitalInvestorIndex = -1
	project.InvestorIndices = append(project.InvestorIndices, 35, 36, 37)
	project.RecipientIndices = append(project.RecipientIndices, 38, 39, 40, 41, 42)
	project.DeveloperIndices = append(project.DeveloperIndices, 44, 45)
	// Define things that will be displayed on the frontend
	project.ResilienceRating = 0.8
	project.ExplorePageSummary.Solar = project.PanelSize
	project.ExplorePageSummary.Storage = project.ExecutiveSummary.ProjectSize["storage"]
	project.ExplorePageSummary.Tariff = project.ExecutiveSummary.Financials["tariff (fixed)"]
	project.ExplorePageSummary.Stage = project.Stage
	project.ExplorePageSummary.Return = project.ExecutiveSummary.Financials["expected return"]
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
		return err
	}

	return nil
}

func bootstrapInvestor(invName, invDescription string) (database.Investor, string, error) {
	// setup investor account
	log.Println(consts.StablecoinSeed, consts.StablecoinPublicKey, consts.StablecoinCode)
	passwd := "password"
	seedpwd := "x"
	investor1, err := database.NewInvestor(invName, passwd, seedpwd, invDescription)
	if err != nil {
		return investor1, "", err
	}
	invSeed, err := wallet.DecryptSeed(investor1.U.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return investor1, "", err
	}
	err = xlm.GetXLM(investor1.U.StellarWallet.PublicKey)
	if err != nil {
		return investor1, "", err
	}
	// trust the stablecoin issuer and give the investor a fixed number of stableusd to invest
	// this helps prevent calling the exchange function that is implicitly called in the payment function
	_, err = assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, "10000000000", invSeed)
	if err != nil {
		return investor1, "", err
	}
	_, _, err = assets.SendAssetFromIssuer(consts.StablecoinCode, investor1.U.StellarWallet.PublicKey, "1000000", consts.StablecoinSeed, consts.StablecoinPublicKey)
	if err != nil {
		return investor1, "", err
	}
	return investor1, invSeed, nil
}

func bootstrapRecipient(recpName, recpDescription string) (database.Recipient, string, error) {
	// setup recipient account
	passwd := "password"
	seedpwd := "x"
	recipient, err := database.NewRecipient(recpName, passwd, seedpwd, recpDescription)
	if err != nil {
		return recipient, "", err
	}
	recpSeed, err := wallet.DecryptSeed(recipient.U.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return recipient, "", err
	}
	err = xlm.GetXLM(recipient.U.StellarWallet.PublicKey)
	if err != nil {
		return recipient, "", err
	}
	return recipient, recpSeed, nil

}

func oneInvestor(projIndex int, invName string, invDescription string, recpName string,
	recpDescription string) error {
	project, err := opensolar.RetrieveProject(projIndex)
	if err != nil {
		return err
	}

	oldStage := project.Stage
	project.Stage = 4 // to enable investments on this particular project
	err = project.Save()
	if err != nil {
		return err
	}

	// passwd := "password"
	seedpwd := "x"

	investor1, invSeed, err := bootstrapInvestor(invName, invDescription)
	if err != nil {
		return err
	}

	recipient1, _, err := bootstrapRecipient(recpName, recpDescription)
	if err != nil {
		return err
	}

	project.RecipientIndex = recipient1.U.Index
	err = project.Save()
	if err != nil {
		return err
	}

	totalValueString, err := utils.ToString(project.TotalValue)
	if err != nil {
		return err
	}
	err = opensolar.Invest(projIndex, investor1.U.Index, totalValueString, invSeed)
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	err = opensolar.UnlockProject(recipient1.U.Username, recipient1.U.Pwhash, projIndex, seedpwd)
	if err != nil {
		return err
	}

	time.Sleep(100 * time.Second)
	project, err = opensolar.RetrieveProject(projIndex)
	if err != nil {
		return err
	}
	project.Stage = oldStage
	err = project.Save()
	if err != nil {
		return err
	}
	return nil
}

func threeInvestor(projIndex int, invName1 string, invDescription1 string, invName2 string, invDescription2 string,
	invName3 string, invDescription3 string, invAmount1 string, invAmount2 string, invAmount3 string,
	recpName string, recpDescription string) error {

	project, err := opensolar.RetrieveProject(projIndex)
	if err != nil {
		return err
	}

	oldStage := project.Stage
	project.Stage = 4 // to enable investments on this particular project
	err = project.Save()
	if err != nil {
		return err
	}

	// passwd := "password"
	// seedpwd := "x"

	investor1, invSeed1, err := bootstrapInvestor(invName1, invDescription1)
	if err != nil {
		return err
	}

	investor2, invSeed2, err := bootstrapInvestor(invName2, invDescription2)
	if err != nil {
		return err
	}

	investor3, invSeed3, err := bootstrapInvestor(invName3, invDescription3)
	if err != nil {
		return err
	}

	recipient1, _, err := bootstrapRecipient(recpName, recpDescription)
	if err != nil {
		return err
	}

	project.RecipientIndex = recipient1.U.Index
	err = project.Save()
	if err != nil {
		return err
	}

	err = opensolar.Invest(projIndex, investor1.U.Index, invAmount1, invSeed1)
	if err != nil {
		return err
	}

	err = opensolar.Invest(projIndex, investor2.U.Index, invAmount2, invSeed2)
	if err != nil {
		return err
	}

	err = opensolar.Invest(projIndex, investor3.U.Index, invAmount3, invSeed3)
	if err != nil {
		return err
	}

	// update local project with changes from storage
	project, err = opensolar.RetrieveProject(projIndex)
	if err != nil {
		return err
	}

	project.BlendedCapitalInvestorIndex = investor1.U.Index
	project.Stage = oldStage
	err = project.Save()
	if err != nil {
		return err
	}
	return nil
}

func sixInvestor(projIndex int, invName1 string, invDescription1 string, invName2 string, invDescription2 string,
	invName3 string, invDescription3 string, invName4 string, invDescription4 string, invName5 string, invDescription5 string,
	invName6 string, invDescription6 string, invAmount1 string, invAmount2 string, invAmount3 string, invAmount4 string,
	invAmount5 string, invAmount6 string, recpName string, recpDescription string) error {

	project, err := opensolar.RetrieveProject(projIndex)
	if err != nil {
		return err
	}

	recipient1, _, err := bootstrapRecipient(recpName, recpDescription)
	if err != nil {
		return err
	}

	oldStage := project.Stage
	project.RecipientIndex = recipient1.U.Index
	project.Stage = 4 // to enable investments on this particular project
	err = project.Save()
	if err != nil {
		return err
	}

	// passwd := "password"
	seedpwd := "x"
	investor1, invSeed1, err := bootstrapInvestor(invName1, invDescription1)
	if err != nil {
		return err
	}

	investor2, invSeed2, err := bootstrapInvestor(invName2, invDescription2)
	if err != nil {
		return err
	}

	investor3, invSeed3, err := bootstrapInvestor(invName3, invDescription3)
	if err != nil {
		return err
	}

	investor4, invSeed4, err := bootstrapInvestor(invName4, invDescription4)
	if err != nil {
		return err
	}

	investor5, invSeed5, err := bootstrapInvestor(invName5, invDescription5)
	if err != nil {
		return err
	}

	investor6, invSeed6, err := bootstrapInvestor(invName6, invDescription6)
	if err != nil {
		return err
	}

	err = project.Save()
	if err != nil {
		return err
	}

	err = opensolar.Invest(projIndex, investor1.U.Index, invAmount1, invSeed1)
	if err != nil {
		return err
	}

	err = opensolar.Invest(projIndex, investor2.U.Index, invAmount2, invSeed2)
	if err != nil {
		return err
	}

	err = opensolar.Invest(projIndex, investor3.U.Index, invAmount3, invSeed3)
	if err != nil {
		return err
	}

	err = opensolar.Invest(projIndex, investor4.U.Index, invAmount4, invSeed4)
	if err != nil {
		return err
	}

	err = opensolar.Invest(projIndex, investor5.U.Index, invAmount5, invSeed5)
	if err != nil {
		return err
	}

	err = opensolar.Invest(projIndex, investor6.U.Index, invAmount6, invSeed6)
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	err = opensolar.UnlockProject(recipient1.U.Username, recipient1.U.Pwhash, projIndex, seedpwd)
	if err != nil {
		return err
	}

	time.Sleep(100 * time.Second)
	project, err = opensolar.RetrieveProject(projIndex)
	if err != nil {
		return err
	}

	project.BlendedCapitalInvestorIndex = investor1.U.Index
	project.Stage = oldStage
	err = project.Save()
	if err != nil {
		return err
	}
	return nil
}

func populateDynamicData1kw() error {
	// setup all the entities that will be involved with the project here
	err := oneInvestor(4, "OpenLab@test.com", "Yale OpenLab", "supasto@test.com", "S.U. Pasto School")
	if err != nil {
		return err
	}
	return nil
}

func populateDynamicData1mw() error {
	// setup all the entities involved with the project here
	err := threeInvestor(5, "OZFunds@test.com", "OZ FundCo", "GreenBank@test.com", "NH Green Bank", "TaxEquity@test.com", "Lancaster Lumber Mill Coop",
		"1000000", "400000", "100000", "LancasterHigh@test.com", "Lancaster Elementary School")
	if err != nil {
		return err
	}

	recipient2, err := database.NewRecipient("Lancastert@test.com", "password", "x", "Town of Lancaster NH")
	if err != nil {
		return err
	}

	project, err := opensolar.RetrieveProject(5)
	if err != nil {
		return err
	}

	// Define the various entities that are associated with a specific project
	project.RecipientIndices = append(project.RecipientIndices, recipient2.U.Index)
	err = project.Save()
	if err != nil {
		return err
	}
	return nil
}

func populateDynamicData10kw() error {

	err := sixInvestor(6, "MatthewMoroney@test.com", "Matthew Moroney", "FranzHochstrasser@test.com", "Franz Hochstrasser", "CTGreenBank@test.com", "Connecticut Green Bank",
		"YaleUniversity@test.com", "Yale University Community Fund", "JeromeGreen@test.com", "Jerome Green", "OpenSolarFund@test.com", "Open Solar Revolving Fund",
		"4000", "4000", "4000", "4000", "4000", "10000", "colhouse@test.com", "Columbus House Foundation")

	if err != nil {
		return err
	}

	recipient2, err := database.NewRecipient("ColumbusHouse@test.com", "password", "x", "Columbus House Foundation")
	if err != nil {
		return err
	}

	project, err := opensolar.RetrieveProject(6)
	if err != nil {
		return err
	}
	// Define the various entities that are associated with a specific project
	project.RecipientIndices = append(project.RecipientIndices, recipient2.U.Index)
	err = project.Save()
	if err != nil {
		return err
	}
	return nil
}
