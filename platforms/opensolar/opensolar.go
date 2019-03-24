package opensolar

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	consts "github.com/YaleOpenLab/openx/consts"
	notif "github.com/YaleOpenLab/openx/notif"
	platform "github.com/YaleOpenLab/openx/platforms"
)

// A Project is the investment structure that will be invested in by people. In the case
// of the opensolar platform, this is a solar system.

// TODO: add more parameters here that would help identify a given solar project

// Project defines the project struct
type Project struct {
	// Describe the project
	Index                     int     // an Index to keep track of how many projects exist
	Name                      string  // the name of the project / the identifier by which its referred to
	State                     string  // the state in which the project has been installed in
	Country                   string  // the country in which the project has been installed in
	TotalValue                float64 // the total money that we need from investors
	PanelSize                 string  // size of the given panel, for diplsaying to the user who wants to bid stuff
	PanelTechnicalDescription string  // This should talk about '10x 100W Komaes etc'
	Inverter                  string  // the inverter of the installed project
	ChargeRegulator           string  // the charge regulator of the installed project
	ControlPanel              string  // the control panel of the installed project
	CommBox                   string  // the comm box of the installed project
	ACTransfer                string  // the AC transfer of the installed project
	SolarCombiner             string  // the solar combiner of the installed project
	Batteries                 string  // the batteries of the installed project. TODO: Batteries should also have a fixed nominal value of capacity, as well as one describing what setup it is.
	IoTHub                    string  // the IoT Hub installed as part of the project
	Metadata                  string  // other metadata which does not have an explicit name can be stored here. Used to derive assetIDs

	// Define parameters related to finance
	MoneyRaised          float64 // total money that has been raised until now
	EstimatedAcquisition int     // the year in which the recipient is expected to repay the initial investment amount by
	BalLeft              float64 // denotes the balance left to pay by the party, percentage raised is not stored in the database since that can be calculated
	InterestRate         float64 // the rate of return for investors
	Tax                  string  // the specifications of the tax system associated with this particular project

	// Define dates of creation and funding
	DateInitiated string // date the project was created on the platform
	DateFunded    string // date that the project completed the stage 4-5 migration
	DateLastPaid  int64  // int64 ie unix time since we need comparisons on this one

	// Define technical paramters
	AuctionType          string  // the type of the auction in question. Default is blind auction unless explicitly mentioned
	InvestmentType       string  // the type of investment - equity crowdfunding, municipal bond, normal crowdfunding, etc defined in models
	PaybackPeriod        int     // the frequency in number of weeks that the recipient has to pay the platform.
	Stage                int     // the stage at which the contract is at, float due to potential support of 0.5 state changes in the future
	InvestorAssetCode    string  // the code of the asset given to investors on investment in the project
	DebtAssetCode        string  // the code of the asset given to recipients on receiving a project
	PaybackAssetCode     string  // the code of the asset given to recipients on receiving a project
	SeedAssetCode        string  // the code of the asset given to seed investors on seed investment in the project
	SeedInvestmentFactor float64 // the factor that a seed investor's investment is multiplied by in case he does invest at the seed stage
	SeedInvestmentCap    float64 // the max amount that a seed investor can put in a project when it is in its seed stages
	ProposedInvetmentCap float64 // the max amount that an investor can invest in when the project is in its proposed stage (stage 2)
	SelfFund             float64 // the amount that a beneficiary / recipient puts in a project wihtout asking from other investors. This is not included as a seed investment because this would mean the recipient pays his own investment back in the project

	// Describe issuer of security and the broker dealer
	SecurityIssuer string // the issuer of the security
	BrokerDealer   string // the broker dealer associated with the project

	// Define the various entities that are associated with a specific project
	RecipientIndex              int       // The index of the project's recipient
	OriginatorIndex             int       // the originator of the project
	GuarantorIndex              int       // the person guaranteeing the specific project in question
	ContractorIndex             int       // the person with the proposed contract
	MainDeveloperIndex          int       // the main developer of the project
	BlendedCapitalInvestorIndex int       // the index of the blended capital investor
	InvestorIndices             []int     // The various investors who have invested in the project
	SeedInvestorIndices         []int     // Investors who took part before the contract was at stage 3
	RecipientIndices            []int     // the indices of the recipient family (offtakers, beneficiaries, etc)
	DeveloperIndices            []int     // the indices of the developers involved in the project`
	ContractorFee               float64   // fee paid to the contractor from the total fee of the project
	OriginatorFee               float64   // fee paid to the originator included in the total value of the project
	DeveloperFee                []float64 // the fees charged by the developers
	DebtInvestor1               string    // debt investor index, if any
	DebtInvestor2               string    // debt investor index, if any
	TaxEquityInvestor           string    // tax equity investor if any

	// Define parameters that will not be defined directly but will be used for the backend flow
	Lock           bool               // lock investment in order to wait for recipient's confirmation
	LockPwd        string             // the recipient's seedpwd. Will be set to null as soon as we use it.
	Votes          int                // the number of votes towards a proposed contract by investors
	AmountOwed     float64            // the amoutn owed to investors as a cumulative sum. Used in case of a breach
	Reputation     float64            // the positive reputation associated with a given project
	StageData      []string           // the data associated with stage migrations
	StageChecklist []map[string]bool  // the checklist that has to be completed before moving on to the next stage
	InvestorMap    map[string]float64 // publicKey: percentage donation
	WaterfallMap   map[string]float64 // publickey:amount map ni order to pay multiple accounts. A bit ugly, but should work fine. Make map before using

	// Define things that will be displayed on the frontend
	Terms                    []TermsHelper        // the terms of the project
	InvestmentMetrics        InvestmentHelper     // investment metrics that might be useful to an investor
	FinancialMetrics         FinancialHelper      // financial metrics that might be useful to an investor
	ProjectSizeMetric        ProjectSizeHelper    // a metric which shows the size of the project
	SustainabilityMetric     SustainabilityHelper // a metric which shows the sustainability index of the project
	AutoReloadInterval       float64              // the interval in which the user's funds reach zero
	ResilienceRating         float64              // resilience of the project
	ActionsRequired          string               // the action(s) required by the user
	Bullet1                  string               // bullet points to be displayed on the project summary page
	Bullet2                  string               // bullet points to be displayed on the project summary page
	Bullet3                  string               // bullet points to be displayed on the project summary page
	LegalProjectOverviewHash string               // hash to be displayed on the project details page
	LegalPPAHash             string               // hash to be displayed on the project details page
	LegalRECAgreementHash    string               // hash to be displayed on the project details page
	GuarantorAgreementHash   string               // hash to be displayed on the project details page
	ContractorAgreementHash  string               // hash to be displayed on the project details page
	StakeholderAgreementHash string               // hash to be displayed on the project details page
	CommunityEnergyHash      string               // hash to be displayed on the project details page
	FinancialReportingHash   string               // hash to be displayed on the project details page
	Contract1                string               // contracts which will be linked to on the project details page
	Contract2                string               // contracts which will be linked to on the project details page
	Contract3                string               // contracts which will be linked to on the project details page
	Contract4                string               // contracts which will be linked to on the project details page
	Contract5                string               // contracts which will be linked to on the project details page
}

// Terms a terms and conditions struct. WIll be used as an array in the main project
type TermsHelper struct {
	Variable      string
	Value         string
	RelevantParty string
	Note          string
	Status        string
	SupportDoc    string
}

type InvestmentHelper struct {
	Capex              string
	Hardware           float64
	FirstLossEscrow    string
	CertificationCosts string
}

type FinancialHelper struct {
	Return    float64
	Insurance string
	Tariff    string
	Maturity  string
}

type ProjectSizeHelper struct {
	PVSolar          string
	Storage          string
	Critical         float64
	InverterCapacity string
}

type SustainabilityHelper struct {
	CarbonDrawdown  string
	CommnunityValue string
	LCA             string
}

//easyjson:json
type SolarProjectArray []Project

// InitializePlatform imports handlers from the main platform struct that are necessary for starting the platform
func InitializePlatform() error {
	return platform.InitializePlatform()
}

// RefillPlatform checks whether the publicKey passed has any xlm and if its balance
// is less than 21 XLM, it proceeds to ask the friendbot for more test xlm
func RefillPlatform(publicKey string) error {
	return platform.RefillPlatform(publicKey)
}

const tellerUrl = "https://localhost"

type statusResponse struct {
	Code   int
	Status string
}

// MonitorTeller monitors a teller and checks whether its live. If not, send an email to platform admins
func MonitorTeller(projIndex int) {
	// call this function only after a specific order has been accepted by the recipient
	for {
		project, err := RetrieveProject(projIndex)
		if err != nil {
			log.Println(err)
			continue
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		req, err := http.NewRequest("GET", tellerUrl+"/ping", nil)
		if err != nil {
			log.Println("did not create new GET request", err)
			notif.SendTellerDownEmail(project.Index, project.RecipientIndex)
			time.Sleep(consts.TellerPollInterval * time.Second)
			continue
		}

		req.Header.Set("Origin", "localhost")
		res, err := client.Do(req)
		if err != nil {
			log.Println("did not make request", err)
			notif.SendTellerDownEmail(project.Index, project.RecipientIndex)
			time.Sleep(consts.TellerPollInterval * time.Second)
			continue
		}
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("error while reading response body", err)
			notif.SendTellerDownEmail(project.Index, project.RecipientIndex)
			time.Sleep(consts.TellerPollInterval * time.Second)
			continue
		}

		var x statusResponse

		err = x.UnmarshalJSON(data)
		if err != nil {
			log.Println("error while unmarshalling data", err)
			notif.SendTellerDownEmail(project.Index, project.RecipientIndex)
			time.Sleep(consts.TellerPollInterval * time.Second)
			continue
		}

		if x.Code != 200 || x.Status != "HEALTH OK" {
			notif.SendTellerDownEmail(project.Index, project.RecipientIndex)
		}

		res.Body.Close()
		time.Sleep(consts.TellerPollInterval * time.Second)
	}
}
