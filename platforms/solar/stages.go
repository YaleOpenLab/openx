package solar

import (
	"fmt"
)

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

// the following functions are helper functions to set the stage for a specific
// project
// we could also alternately define contract states and then read the state from
// our side and then compress this into a single function
func (a *SolarProject) SetPreOriginProject() error {
	a.Stage = 0
	return a.Save()
}

func (a *SolarProject) SetLegalContractStage() error {
	a.Stage = 0.5
	return a.Save()
}

func (a *SolarProject) SetOriginProject() error {
	a.Stage = 1
	return a.Save()
}

func (a *SolarProject) SetOpenForMoneyStage() error {
	a.Stage = 1.5
	return a.Save()
}

func (a *SolarProject) SetProposedProject() error {
	a.Stage = 2
	return a.Save()
}

func (a *SolarProject) SetFinalizedProject() error {
	a.Stage = 3
	return a.Save()
}

func (a *SolarProject) SetFundedProject() error {
	a.Stage = 4
	return a.Save()
}

func (a *SolarProject) SetInstalledProjectStage() error {
	a.Stage = 5
	return a.Save()
}

func (a *SolarProject) SetPowerGenerationStage() error {
	a.Stage = 6
	return a.Save()
}

func PromoteStage0To1Project(projects []SolarProject, index int) error {
	// we need to upgrade the contract's whose index is contractIndex to stage 1
	for _, elem := range projects {
		if elem.Params.Index == index {
			return elem.SetOriginProject() // upgrade stage of this project
		}
	}
	return fmt.Errorf("SolarProject not found, erroring!")
}
