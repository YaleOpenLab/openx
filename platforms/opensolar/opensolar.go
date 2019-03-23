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
	Index int // an Index to keep quick track of how many projects exist

	Name           string // the name of the project / the identifier by which its referred to
	TotalValue     float64 // the total money that we need from investors
	MoneyRaised    float64 // total money that has been raised until now
	ETA            int     // the year in which the recipient is expected to repay the initial investment amount by
	BalLeft        float64 // denotes the balance left to pay by the party, percentage raised is not stored in the database since that can be calculated
	Votes          int     // the number of votes towards a proposed contract by investors
	SecurityIssuer string  // the issuer of the security
	BrokerDealer   string  // the broker dealer associated with the project

	// different assets associated with the platform
	InvestorAssetCode string
	DebtAssetCode     string
	PaybackAssetCode  string
	SeedAssetCode     string

	// dates on which the project was initiated, funded an last paid towards
	DateInitiated string
	DateFunded    string
	DateLastPaid  int64 // int64 ie unix time since we need comparisons on this one

	// params that help define the specifications of the installation
	// MW: What does state really mean here? Is this the project name?
	State        string  // the state in which the project has been installed in
	Country      string  // the country in which the project has been installed in
	InterestRate float64 // the rate of return for investors
	Tax          string  // the specifications of the tax system associated with this particular project
	// TODO: I see we have 'Panel Size' which should just be the denominal value only (eg. 1000W), but there should also be a 'Panel technical description'
	// This should talk about '10x 100W Komaes etc'
	PanelSize       string // size of the given panel, for diplsaying to the user who wants to bid stuff
	Inverter        string
	ChargeRegulator string
	ControlPanel    string
	CommBox         string
	ACTransfer      string
	SolarCombiner   string
	//TODO: Batteries should also have a fixed nominal value of capacity, as well as one describing what setup it is.
	Batteries string
	IoTHub    string
	Metadata  string // other metadata which does not have an explicit name can be stored here. Used to derive assetIDs

	// List of entities other than the contractor
	Originator           Entity
	OriginatorFee        float64 // fee paid to the originator included in the total value of the project
	Developer            Entity  // the developer who is responsible for installing the solar panels and the IoT hubs
	Guarantor            Entity  // the person guaranteeing the specific project in question
	SeedInvestmentFactor float64 // the factor that a seed investor's investment is multiplied by in case he does invest at the seed stage
	SeedInvestmentCap    float64 // the max amount that a seed investor can put in a project when it is in its seed stages
	ProposedInvetmentCap float64 // the max amount that an investor can invest in when the project is in its proposed stage (stage 2)
	SelfFund             float64 // the amount that a beneficiary / recipient puts in a project wihtout asking from other investors. This is not included as a seed investment because this would mean the recipient pays his own investment back in the project

	// List of contractor entities
	Contractor             Entity  // the person with the proposed contract
	ContractorFee          float64 // fee paid to the contractor from the total fee of the project
	SecondaryContractor    Entity  // this is the secondary contractor involved in the project
	SecondaryContractorFee float64 // the fee to be paid towards the secondary contractor
	TertiaryContractor     Entity  // tertiary contractor if any can be added to the system
	TertiaryContractorFee  float64 // the fee to be paid towards the tertiary contractor
	DeveloperFee           float64 // the fee charged by the developer

	RecipientIndex      int   // The index of the project's recipient
	InvestorIndices     []int // The various investors who have invested in the project
	SeedInvestorIndices []int // Investors who took part before the contract was at stage 3

	Stage       int    // the stage at which the contract is at, float due to potential support of 0.5 state changes in the future
	AuctionType string // the type of the auction in question. Default is blind auction unless explicitly mentioned

	Reputation float64 // the positive reputation associated with a given project
	Lock       bool    // lock investment in order to wait for recipient's confirmation
	LockPwd    string  // the recipient's seedpwd. Will be set to null as soon as we use it.

	InvestmentType string // the type of investment - equity crowdfunding, municipal bond, normal crowdfunding, etc defined in models

	PaybackPeriod int // the frequency in number of weeks that the recipient has to pay the platform.
	// this has to be set to atleast a week since the payback monitoring thread runs every week. Ideally, we could
	// provide users with a predefined list of payback periods periods

	// List of checklists that the user can go and check in the past whether they have been fulfilled or not
	// this is a string-strign map since I can add any arbitrary data that I want to without checking for stuff.
	// this is of length 9 since there are nine stages defined for the opensolar platform
	StageChecklist []map[string]bool

	// List of data associated with each checkpoint in order for someone who comes in later to verify
	// that we indeed have the right project. THe various hashes and stuff are stored here instead of
	// having separate fields for each contract
	StageData   []string
	InvestorMap map[string]float64 // publicKey: percentage donation

	Terms              []TermsHelper // the terms of the project
	AutoReloadInterval float64       // the interval in which the user's funds reach zero
	ActionsRequired    string        // the action(s) required by the user

	// these are bullet points that would be displayed along with project decription on the main screen
	Bullet1 string
	Bullet2 string
	Bullet3 string

	Pictures []string // an array of the pictures in base64 that are stored on the backend
	// TOOD: see if we can handle this in a simpler way

	ResilienceRating     float64
	InvestmentMetrics    InvestmentHelper
	FinancialMetrics     FinancialHelper
	ProjectSizeMetric    ProjectSizeHelper
	SustainabilityMetric SustainabilityHelper

	LegalProjectOverviewHash string
	LegalPPAHash             string
	LegalRECAgreementHash    string
	GuarantorAgreementHash   string
	ContractorAgreementHash  string
	StakeholderAgreementHash string
	CommunityEnergyHash      string
	FinancialReportingHash   string

	// list of smart contracts that we must link to on the project page
	Contract1 string
	Contract2 string
	Contract3 string
	Contract4 string
	Contract5 string

	DeveloperIndices            []int
	MainDeveloper               Entity
	MainOriginator              Entity
	BlendedCapitalInvestorIndex int
	RecipientIndices            []int
	DebtInvestor1               string
	DebtInvestor2               string
	TaxEquityInvestor           string
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
