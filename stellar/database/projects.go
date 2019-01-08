package database

// the database package maintains read / write operations to the orderbook
// we need an orderbook because there is no state on Stellar which makes it
// difficult for us to store this on the blockchain. We use boltdb no since we don't
// do that much relational mapping, but in the case we need that, we can modify
// this package to do that.
import (
	"encoding/json"
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
)

// DBParam is a backend meta structure used by the backend Project, which encompasses
// more information than this structure but all that information would nto be
// needed for transacting in assets and interfacing with other elements in the system

type DBParams struct {
	Index int // an Index to keep quick track of how many projects exist

	PanelSize   string // size of the given panel, for diplsaying to the user who wants to bid stuff
	TotalValue  int    // the total money that we need from investors
	Location    string // where this specific solar panel is located
	MoneyRaised int    // total money that has been raised until now
	Years       int    // number of years the recipient has chosen to opt for
	Metadata    string // any other metadata can be stored here

	Votes int // the number of votes towards a proposed contract by investors

	// once all funds have been raised, we need to set assetCodes
	INVAssetCode string
	DEBAssetCode string
	PBAssetCode  string

	BalLeft float64 // denotes the balance left to pay by the party

	DateInitiated string // date the project was created
	DateFunded    string // date when the project was funded
	DateLastPaid  string // date the project was last paid

	ProjectRecipient Recipient
	ProjectInvestors []Investor // TODO: get feedback on whether this is in its right place and whether this should be moved to contract
	// Percentage raised is not stored in the database since that can be calculated by the UI
}

// Contracts and Projects are used interchangeably below
// A contract has six Stages (right now an order has 6 stages and later both will be merged)
// seed funding and seed assets are also TODOs, though investors can see the assets
// now and can transfer funds if they really want to
// TODO: propagate one transaction for ever major state change

// A legal contract should ideally be stored on ipfs and we must keep track of the
// ipfs hash so that we can retrieve it later when required

// A Project is what is stored in the database and what is used by other packages
// Project imports DBParams since having everythin inside one struct is tedious
// and DBParams already has lots of keys. Also, this doesn't affect the way its
// actually stored in the database, so its a nice way to do it.
// DBParams is also what's needed by the assets and other stuff whereas the other fields
// are needed in other parts, another nice distinction
type Project struct {
	Params DBParams // Params is the former Order struct improted into the new Project structure

	Originator    Entity // a specific contract must hold the person who originated it
	Contractor    Entity // the person with the proposed contract
	Guarantor     Entity // the person guaranteeing the specific project in question
	OriginatorFee int    // fee paid to the originator from the total fee of the project
	ContractorFee int    // fee paid to the contractor from the total fee of the project

	Stage float64 // the stage at which the contract is at, float due to potential support of 0.5 state changes in the future
}

// TODO: get comments on the various stages involved here
var (
	PreOriginProject      = float64(0) // Stage 0: Originator approaches the recipient to originate an order
	LegalContractStage    = 0.5        // Stage 0.5: Legal contract between the originator and the recipient, out of blockchain
	OriginProject         = float64(1) // Stage 1: Originator proposes a contract on behalf of the recipient
	OpenForMoneyStage     = 1.5        // Stage 1.5: The contract, even though not final, is now open to investors' money
	ProposedProject       = float64(2) // Stage 2: Contractors propose their contracts and investors can vote on them if they want to
	FinalizedProject      = float64(3) // Stage 3: Recipient chooses a particular contract for finalization
	FundedProject         = float64(4) // Stage 4: Review the legal contract and finalize a particular contractor
	InstalledProjectStage = float64(5) // Stage 5: Installation of the panels / houses by the developer and contractor
	PowerGenerationStage  = float64(6) // Stage 6: Power generation and trigerring automatic payments, cover breach, etc.
	DebtPaidOffStage      = float64(7) // Stage 7: The stage at which the recipient pays back for his solar panels
)

// so a contract's rough workflow is like
// origincontract (0) -> approval by recipient (1) -> OpenForMoneyStage (1.5) -> ...
// NewOriginProject returns a new project passed a project and originator to assign to
// stage is set automatically to 1 by the call to SetOriginProject
func NewOriginProject(project DBParams, originator Entity) (Project, error) {
	// need variadic params to store optional stuff
	var proposedProject Project
	proposedProject.Params = project
	proposedProject.Originator = originator
	err := proposedProject.SetOriginProject()
	return proposedProject, err
}

// RecipientAuthorizeContract authorizes a specific project from the recipients side
// if you already have the project and the recipient struct ready to pass.
// this function is not used right now, but would be when we finalize the various
// stages of a project
func (project *Project) RecipientAuthorizeContract(recipient Recipient) error {
	if project.Params.ProjectRecipient.U.Name != recipient.U.Name {
		return fmt.Errorf("You can't authorize a project which is not assigned to you!")
	}
	// set the project as both originated and ready for investors' money
	err := project.SetOriginProject()
	if err != nil {
		return err
	}
	err = project.SetOpenForMoneyStage()
	if err != nil {
		return err
	}
	return nil
}

// FinalizeProject finalizes a specific project proposed by contractors and sets the
// stage to three and allows investors to f ormally invest. Investors can technically
// invest after stage 1.5 with something like seed funding, but this is the main funding part
// we are looking towards
// select a contract here directly so that we avoid the extra db lookup if we accepted index
func FinalizeProject(project Project) error {
	// now we need to search using the project's location and size field since a contractor
	// can not change that while proposing a contract
	// retrieve all contracts and check
	allProjectsDB, err := RetrieveAllProjects()
	if err != nil {
		return err
	}
	for _, dbProjects := range allProjectsDB {
		if dbProjects.Params.Location == project.Params.Location && dbProjects.Params.PanelSize == project.Params.PanelSize {
			// this is the contract whose stage we need ot upgrade and whose thing we must add to the contract
			// TODO: weak check, should have something better here
			dbProjects.Params = project.Params         // overwrite price related details
			dbProjects.Contractor = project.Contractor // store the contractor for the given order
			dbProjects.Guarantor = project.Guarantor   // add guarantor
			dbProjects.SetFinalizedProject()           // set the stage to be open for investors
			dbProjects.Save()                          // save in db
			return nil
		}
	}
	return fmt.Errorf("Finalized Project not found in db")
}

// Save or Insert inserts a specific Project into the database
func (a *Project) Save() error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ProjectsBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.ItoB(a.Params.Index)), encoded)
	})
	return err
}

// RetrieveProject retrieves the project with the specified index from the database
func RetrieveProject(key int) (Project, error) {
	var inv Project
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ProjectsBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return nil
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, err
}

// RetrieveAllProjects retrieves all projects from the database
func RetrieveAllProjects() ([]Project, error) {
	var arr []Project
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				break
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				return err
			}
			// append only contracts which are open for funding and below
			arr = append(arr, rProject)
		}
		return nil
	})
	return arr, err
}

func RetrieveProjects(stage float64) ([]Project, error) {
	var arr []Project
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// this is where the key does not exist
				return nil
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			if rProject.Stage == stage {
				// return contracts which have been originated and are not final yet
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

func RetrieveProjectsC(stage float64, index int) ([]Project, error) {
	var arr []Project
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// this is where the key does not exist
				return nil
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			if rProject.Stage == stage && rProject.Contractor.U.Index == index {
				// return contracts which have been originated and are not final yet
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

// RetrieveOriginProjectsIO is used when we want to display the list of originated
// contracts to the originator
func RetrieveProjectsO(stage float64, index int) ([]Project, error) {
	var arr []Project
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// this is where the key does not exist
				return nil
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			if rProject.Stage == stage && rProject.Originator.U.Index == index {
				// return contracts which have been originated and are not final yet
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

func RetrieveProjectsR(stage float64, index int) ([]Project, error) {
	var arr []Project
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// this is where the key does not exist
				return nil
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			if rProject.Stage == stage && rProject.Params.ProjectRecipient.U.Index == index {
				// return contracts which have been originated and are not final yet
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

func PromoteStage0To1Project(projects []Project, index int) error {
	// we need to upgrade the contract's whose index is contractIndex to stage 1
	for _, elem := range projects {
		if elem.Params.Index == index {
			// increase this contract's stage
			log.Println("UPGRADING PROJECT INDEX", elem.Params.Index)
			err := elem.SetOriginProject()
			return err
		}
	}
	return fmt.Errorf("Project not found, erroring!")
}

// the following functions are helper functions to set the stage for a specific
// project
// we could also alternately define contract states and then read the state from
// our side and then compress this into a single function
func (a *Project) SetPreOriginProject() error {
	a.Stage = 0
	return a.Save()
}

func (a *Project) SetLegalContractStage() error {
	a.Stage = 0.5
	return a.Save()
}

func (a *Project) SetOriginProject() error {
	a.Stage = 1
	return a.Save()
}

func (a *Project) SetOpenForMoneyStage() error {
	a.Stage = 1.5
	return a.Save()
}

func (a *Project) SetProposedProject() error {
	a.Stage = 2
	return a.Save()
}

func (a *Project) SetFinalizedProject() error {
	a.Stage = 3
	return a.Save()
}

func (a *Project) SetFundedProject() error {
	a.Stage = 4
	return a.Save()
}

func (a *Project) SetInstalledProjectStage() error {
	a.Stage = 5
	return a.Save()
}

func (a *Project) SetPowerGenerationStage() error {
	a.Stage = 6
	return a.Save()
}

func FindInKey(key int, arr []Project) (Project, error) {
	var dummy Project
	for _, elem := range arr {
		if elem.Params.Index == key {
			return elem, nil
		}
	}
	return dummy, fmt.Errorf("Not found")
}

// CalculatePayback is a TODO function that should simply sum the PBToken
// balance and then return them to the frontend UI for a nice display
func (project Project) CalculatePayback(amount string) string {
	// the idea is that we should be able ot pass an assetId to this function
	// and it must calculate how much time we have left for payback. For this example
	// until twe do the db stuff, lets pass a few params (although this could be done
	// separately as well).
	// TODO: this functon needs to be the payback function
	amountF := utils.StoF(amount)
	amountPB := (amountF / float64(project.Params.TotalValue)) * float64(project.Params.Years*12)
	amountPBString := utils.FtoS(amountPB)
	return amountPBString
}
