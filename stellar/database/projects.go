package database

import (
	"encoding/json"
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
)

// the database package maintains read / write operations to the orderbook
// we need an orderbook because there is no state on Stellar which makes it
// difficult for us to store this on the blockchain. We use boltdb no since we don't
// do that much relational mapping, but in the case we need that, we can modify
// this package to do that.

// DBParam is a backend meta structure used by the backend Project, which encompasses
// more information than this structure but all that information would nto be
// needed for transacting in assets and interfacing with other elements in the system
type DBParams struct {
	// Data regarding the location of the project&
	Index int // an Index to keep quick track of how many projects exist

	PanelSize   string // size of the given panel, for diplsaying to the user who wants to bid stuff
	TotalValue  int    // the total money that we need from investors
	Location    string // where this specific solar panel is located
	MoneyRaised int    // total money that has been raised until now
	Years       int    // number of years the recipient has chosen to opt for
	Metadata    string // any other metadata can be stored here

	Funded  bool // convenience parameter to see if the project has been funded
	PaidOff bool // whether the asset has been paidoff by the recipient

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
	ProjectInvestors []Investor
	// Percentage raised is not stored in the database since that can be calculated by the UI
}

// A contract has six Stages (right now an order has 6 stages and later both will be merged)
// TODO: implement 0.5 stages
// seed funding and seeda ssets are also TODOs, thoguh investors can see the assets
// now and can transfer funds if they really want to
// look into state commitments and committing state in the memo field of transactions
// and then having to propagate one transaction for ever major state change

// A legal contract should ideally be sotred on ipfs and we must keep track of the
// ipfs hash so that we can retrieve it later when required
// this is a metastructure now since we dont store this directly in the database.
// TODO: store this seaprately
type Project struct {
	Params DBParams

	Originator    Entity // a specific contract must hold the person who originated it
	Contractor    Entity // the person with the proposed contract
	Guarantor     Entity // the person guaranteeing the specific project in question
	OriginatorFee int    // fee paid to the originator from the total fee of the project
	ContractorFee int    // fee paid to the contractor from the total fee of the project

	Stage float64
	// this could also have votes associated with it, but we aren't doing that right away
	// TODO: add vote stuff here
}

// TODO: get comments on the various stages involved here
var (
	OriginProposedContractStage = 0 // Stage 0: Originator approaches the recipient to originate an order
	// LegalContractStage          = 0.5 // Stage 0.5: Legal contract between the originator and the recipient, out of blockchain
	OriginContractStage         = 1   // Stage 1: Originator proposes a contract on behalf of the recipient
	OpenForMoneyStage           = 1.5 // Stage 1.5: The contract, even though not final, is now open to investors' money
	ProposedContractStage       = 2   // Stage 2: Contractors propose their contracts and investors can vote on them if they want to
	RecipientFinalContractStage = 3   // Stage 3: Recipient chooses a particular contract for finalization
	FundedContractStage         = 4   // Stage 4: Review the legal contract and finalize a particular contractor
	InstalledProjectStage       = 5   // Stage 5: Installation of the panels / houses by the developer and contractor
	PowerGenerationStage        = 6   // Stage 6: Power generation and trigerring automatic payments, cover breach, etc.
	// Stage 5 is the boundary between when a contract becomes a project
)

// things that are stored in the contracts database are of minimum stage 1. The contracts database is only for final
// contracts. All other contracts are stored in their respective entities' slices and NOT in the contracts database.
// so a contract's rough workflow is like
// origincontract (0) -> approval by recipient (1) -> OpenForMoneyStage (1.5) -> ...
func NewOriginProject(project DBParams, originator Entity) (Project, error) {
	// need variadic params to store optional stuff
	var proposedContract Project
	proposedContract.Params = project
	proposedContract.Originator = originator
	err := proposedContract.SetOriginContractStage()
	return proposedContract, err
}

func (contract *Project) RecipientAuthorizeContract(recipient Recipient) error {
	if contract.Params.ProjectRecipient.U.Name != recipient.U.Name {
		return fmt.Errorf("You can't authorize a contract which is not yours")
	}
	// set the contract as both originated and ready for investors' money
	err := contract.SetOriginContractStage()
	if err != nil {
		return err
	}
	err = contract.SetOpenForMoneyStage()
	if err != nil {
		return err
	}
	return nil
}

// the proposed contracts are stored inside a contractor's proposecontracts slice
// and are not stored in the contracts db. So one we have a winning contract and contractor,
// we must check the originator for that project and then upgrade the contract associated
// with it

func FinalizeProject(finalizedContract Project) error {
	// now we need to search using the project's location and size field since a contractor
	// can not change that while proposing a contract
	// retrieve all contracts and check
	allContractsDB, err := RetrieveAllProjects()
	if err != nil {
		return err
	}
	for _, dbContracts := range allContractsDB {
		if dbContracts.Params.Location == finalizedContract.Params.Location && dbContracts.Params.PanelSize == finalizedContract.Params.PanelSize {
			// this is the contract whose stage we need ot upgrade and whose thing we must add to the contract
			// TODO: weak check, should have something better here
			dbContracts.Params = finalizedContract.Params         // overwrite price related details
			dbContracts.Contractor = finalizedContract.Contractor // store the contractor for the given order
			dbContracts.Guarantor = finalizedContract.Guarantor   // add guarantor
			dbContracts.SetRecipientFinalContractStage()          // set the stage to be open for investors
			dbContracts.Save()                                    // save in db
			return nil
		}
	}
	return fmt.Errorf("Finalized Project not found in db")
}

func (a *Project) Save() error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ContractBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.ItoB(a.Params.Index)), encoded)
	})
	return err
}

func RetrieveProject(key int) (Project, error) {
	var inv Project
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ContractBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return nil
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, err
}

func RetrieveAllProjects() ([]Project, error) {
	var arr []Project
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ContractBucket)
		for i := 1; ; i++ {
			var rContract Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				break
			}
			err := json.Unmarshal(x, &rContract)
			if err != nil {
				return err
			}
			// append only contracts which are open for funding and below
			arr = append(arr, rContract)
		}
		return nil
	})
	return arr, err
}

func RetrieveAllProposedProjects(recpIndex int) ([]Project, error) {
	var arr []Project
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ContractBucket)
		for i := 1; ; i++ {
			var rContract Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				break
			}
			err := json.Unmarshal(x, &rContract)
			if err != nil {
				return err
			}
			if rContract.Stage == 2 && rContract.Params.ProjectRecipient.U.Index == recpIndex {
				// append only those contracts with stage 2
				arr = append(arr, rContract)
			}
		}
		return nil
	})
	return arr, err
}

func (a *Project) SetOriginProposedContractStage() error {
	a.Stage = 0
	return a.Save()
}

func (a *Project) SetLegalContractStage() error {
	a.Stage = 0.5
	return a.Save()
}

func (a *Project) SetOriginContractStage() error {
	a.Stage = 1
	return a.Save()
}

func (a *Project) SetOpenForMoneyStage() error {
	a.Stage = 1.5
	return a.Save()
}

func (a *Project) SetProposedContractStage() error {
	a.Stage = 2
	return a.Save()
}

func (a *Project) SetRecipientFinalContractStage() error {
	a.Stage = 3
	return a.Save()
}

func (a *Project) SetFundedContractStage() error {
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

// this function retrieves all originated contracts, we don't bother about the other
// projects (projects may be final / non final as well)
func RetrieveOriginatedProjects() ([]Project, error) {
	var arr []Project
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(ContractBucket)
		for i := 1; ; i++ {
			var rContract Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// this is where the key does not exist
				return nil
			}
			err := json.Unmarshal(x, &rContract)
			if err != nil {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			if rContract.Stage == 1 {
				// return contracts which have been originated and are not final yet
				arr = append(arr, rContract)
			}
		}
		return nil
	})
	return arr, err
}

func RetrieveStage3Projects() ([]Project, error) {
	var arr []Project
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(ContractBucket)
		for i := 1; ; i++ {
			var rContract Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// this is where the key does not exist
				return nil
			}
			err := json.Unmarshal(x, &rContract)
			if err != nil {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			if rContract.Stage == 3 {
				// return contracts which have been originated and are not final yet
				arr = append(arr, rContract)
			}
		}
		return nil
	})
	return arr, err
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
