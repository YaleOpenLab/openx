package main

import (
	"encoding/json"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
)

func parseYaml(fileName string, feJson string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName(fileName)
	viper.AddConfigPath("./data-sandbox")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "error while reading values from config file")
	}

	var project opensolar.Project
	terms := make([]opensolar.TermsHelper, 6)
	termsHelper := viper.Get("Terms").(map[string]interface{})

	i := 0
	for _, elem := range termsHelper {
		// elem inside here is a map of "variable": values.
		newMap := elem.(map[string]interface{})
		terms[i].Variable = newMap["variable"].(string)
		terms[i].Value = newMap["value"].(string)
		terms[i].RelevantParty = newMap["relevantparty"].(string)
		terms[i].Note = newMap["note"].(string)
		terms[i].Status = newMap["status"].(string)
		terms[i].SupportDoc = newMap["supportdoc"].(string)
		i += 1
	}

	project.Terms = terms
	var executiveSummary opensolar.ExecutiveSummaryHelper

	execSummaryReader := viper.Get("ExecutiveSummary.Investment").(map[string]interface{})
	execSummaryWriter := make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.Investment = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.Financials").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.Financials = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.ProjectSize").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.ProjectSize = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.SustainabilityMetrics").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.SustainabilityMetrics = execSummaryWriter

	project.ExecutiveSummary = executiveSummary

	var bullets opensolar.BulletHelper
	bullets.Bullet1 = viper.Get("Bullets.Bullet1").(string)
	bullets.Bullet2 = viper.Get("Bullets.Bullet2").(string)
	bullets.Bullet3 = viper.Get("Bullets.Bullet3").(string)

	project.Bullets = bullets

	var architecture opensolar.ArchitectureHelper

	architecture.SolarArray = viper.Get("Architecture.SolarArray").(string)
	architecture.DailyAvgGeneration = viper.Get("Architecture.DailyAvgGeneration").(string)
	architecture.System = viper.Get("Architecture.System").(string)
	architecture.InverterSize = viper.Get("Architecture.InverterSize").(string)

	project.Architecture = architecture

	project.Index = viper.Get("Index").(int)
	project.Name = viper.Get("Name").(string)
	project.State = viper.Get("State").(string)
	project.Country = viper.Get("Country").(string)
	project.TotalValue = viper.Get("TotalValue").(float64)
	project.Metadata = viper.Get("Metadata").(string)
	project.PanelSize = viper.Get("PanelSize").(string)
	project.PanelTechnicalDescription = viper.Get("PanelTechnicalDescription").(string)
	project.Inverter = viper.Get("Inverter").(string)
	project.ChargeRegulator = viper.Get("ChargeRegulator").(string)
	project.ControlPanel = viper.Get("ControlPanel").(string)
	project.CommBox = viper.Get("CommBox").(string)
	project.ACTransfer = viper.Get("ACTransfer").(string)
	project.SolarCombiner = viper.Get("SolarCombiner").(string)
	project.Batteries = viper.Get("Batteries").(string)
	project.IoTHub = viper.Get("IoTHub").(string)
	project.Rating = viper.Get("Rating").(string)
	project.EstimatedAcquisition = viper.Get("EstimatedAcquisition").(int)
	project.BalLeft = viper.Get("BalLeft").(float64)
	project.InterestRate = viper.Get("InterestRate").(float64)
	project.Tax = viper.Get("Tax").(string)
	project.DateInitiated = viper.Get("DateInitiated").(string)
	project.DateFunded = viper.Get("DateFunded").(string)
	project.AuctionType = viper.Get("AuctionType").(string)
	project.InvestmentType = viper.Get("InvestmentType").(string)
	project.PaybackPeriod = viper.Get("PaybackPeriod").(int)
	project.Stage = viper.Get("Stage").(int)
	project.SeedInvestmentFactor = viper.Get("SeedInvestmentFactor").(float64)
	project.SeedInvestmentCap = viper.Get("SeedInvestmentCap").(float64)
	project.ProposedInvetmentCap = viper.Get("ProposedInvetmentCap").(float64)
	project.SelfFund = viper.Get("SelfFund").(float64)
	project.SecurityIssuer = viper.Get("SecurityIssuer").(string)
	project.BrokerDealer = viper.Get("BrokerDealer").(string)
	project.EngineeringLayoutType = viper.Get("EngineeringLayoutType").(string)

	project.FEText, err = parseJsonText(feJson)
	if err != nil {
		log.Fatal(err)
	}

	return project.Save()
}

func populateStaticData() error {
	var err error
	err = createAllStaticEntities()
	if err != nil {
		return err
	}
	err = parseYaml("1kwy", "data-sandbox/1kw.json")
	if err != nil {
		return err
	}
	// project: One Kilowatt Project
	// STAGE 7 - Puerto Rico
	err = populateStaticData1kw()
	if err != nil {
		return err
	}
	err = parseYaml("1mwy", "data-sandbox/1mw.json")
	if err != nil {
		return err
	}
	// project: One Megawatt Project
	// STAGE 4 - New Hampshire
	err = populateStaticData1mw()
	if err != nil {
		return err
	}
	err = parseYaml("10kwy", "data-sandbox/10kw.json")
	if err != nil {
		return err
	}
	// project: Ten Kilowatt Project
	// STAGE 8 - Connecticut Homeless Shelter
	err = populateStaticData10kw()
	if err != nil {
		return err
	}
	err = parseYaml("10mwy", "data-sandbox/10mw.json")
	if err != nil {
		return err
	}
	// project: Ten Megawatt Project
	// STAGE 2 - Puerto Rico Public School Bond
	err = populateStaticData10MW()
	if err != nil {
		return err
	}
	err = parseYaml("100kwy", "data-sandbox/100kw.json")
	if err != nil {
		return err
	}
	// project: One HUndred Kilowatt Project
	// STAGE 1 - Rwanda Project
	err = populateStaticData100KW()
	if err != nil {
		return err
	}
	return nil
}

func populateDynamicData() error {
	var err error
	err = populateDynamicData1kw()
	if err != nil {
		return err
	}
	err = populateDynamicData1mw()
	if err != nil {
		return err
	}
	err = populateDynamicData10kw()
	if err != nil {
		return err
	}
	return nil
}

// CreateSandbox is the main function that controls data insertion as part of the sandbox environment
func CreateSandbox() error {
	// project: Puerto Rico Project
	// STAGE 7 - Puerto Rico
	var err error
	err = populateStaticData()
	if err != nil {
		return err
	}
	err = populateDynamicData()
	if err != nil {
		return err
	}
	return nil
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