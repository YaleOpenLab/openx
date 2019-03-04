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

	TotalValue   float64 // the total money that we need from investors
	MoneyRaised  float64 // total money that has been raised until now
	Years        int     // number of years the recipient is expected to repay the initial investment amount by
	InterestRate float64 // the interest rate provided to potential investors
	BalLeft      float64 // denotes the balance left to pay by the party, percentage raised is not stored in the database since that can be calculated
	Votes        int     // the number of votes towards a proposed contract by investors

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
	Location        string // where this specific solar panel is located
	PanelSize       string // size of the given panel, for diplsaying to the user who wants to bid stuff
	Inverter        string
	ChargeRegulator string
	ControlPanel    string
	CommBox         string
	ACTransfer      string
	SolarCombiner   string
	Batteries       string
	IoTHub          string
	Metadata        string // other metadata which does not have an explicit name can be stored here. Used to derive assetIDs

	// List of entities other than the contractor
	Originator    Entity
	OriginatorFee float64 // fee paid to the originator included in the total value of the project
	Developer     Entity  // the developer who is responsible for installing the solar panels and the IoT hubs
	Guarantor     Entity  // the person guaranteeing the specific project in question

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

	// Various ipfs hashes that we need to store
	OriginatorMoUHash       string // the contract between the originator and the recipient at stage LegalContractStage
	ContractorContractHash  string // the contract between the contractor and the platform at stage ProposeProject
	InvPlatformContractHash string // the contract between the investor and the platform at stage FundedProject
	RecPlatformContractHash string // the contract between the recipient and the platform at stage FundedProject
	SpecSheetHash           string // the ipfs hash of the specification document containing installation details

	Reputation float64 // the positive reputation associated with a given project
	Lock       bool    // lock investment in order to wait for recipient's confirmation
	LockPwd    string  // the recipient's seedpwd. Will be set to null as soon as we use it.

	InvestmentType string // the type of investment - equity crowdfunding, municipal bond, normal crowdfunding, etc defined in models

	PaybackPeriod int // the frequency in number of weeks that the recipient has to pay the platform.
	// this has to be set to atleast a week since the payback monitoring thread runs every week. Ideally, we could
	/// provide users with a predefined list of payback periods periods
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
