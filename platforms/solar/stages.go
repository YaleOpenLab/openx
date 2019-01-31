package solar

import (
	database "github.com/OpenFinancing/openfinancing/database"
)

var (
	PreOriginProject      = float64(0) // Stage 0: Originator approaches the recipient to originate an order. This is a project proposal.
	LegalContractStage    = 0.5        // Stage 0.5: Legal agreement (eg. MOU or letter of intent) between the originator and the recipient, out of blockchain. Can use a 2 of 2 multisig.
	OriginProject         = float64(1) // Stage 1: Originator/s proposes a contract on behalf of the recipient.
	OpenForMoneyStage     = 1.5        // Stage 1.5: The project, even though not final, is now open to investors' money
	ProposedProject       = float64(2) // Stage 2: Contractors propose their contracts and investors can vote on them if they want to
	FinalizedProject      = float64(3) // Stage 3: Recipient chooses a particular contract for finalization. This can be arbitraty or following a specific tender process
	FundedProject         = float64(4) // Stage 4: Extend and Review legal contracts, (re)open for investment and finalize a particular contractor
	InstalledProjectStage = float64(5) // Stage 5: Installation of the panels / houses by the developer and contractor
	PowerGenerationStage  = float64(6) // Stage 6: Power generation and trigerring automatic payments, cover breach, etc.
	DebtPaidOffStage      = float64(7) // Stage 7: The stage at which the recipient pays back for his solar panels
)

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
	a.Reputation = a.Params.TotalValue
	// upgrade reputation since totalValue would have changed originated contract
	err := a.Save()
	if err != nil {
		return err
	}
	// modify originator reputation now that the final price is fixed
	return RepOriginatedProject(a.Originator.U.Index, a.Params.Index)
}

func (a *Project) SetFundedProject() error {
	a.Stage = 4
	return a.Save()
}

func (a *Project) SetInstalledProjectStage() error {
	a.Stage = 5
	err := a.Save()
	if err != nil {
		return err
	}
	// modify contractor Reputation now that a project has been installed
	err = a.Contractor.U.IncreaseReputation(a.Params.TotalValue * ContractorWeight)
	if err != nil {
		return err
	}
	for _, elem := range a.ProjectInvestors {
		err := database.ChangeInvReputation(elem.U.Index, a.Params.TotalValue*InvestorWeight)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Project) SetPowerGenerationStage() error {
	a.Stage = 6
	err := a.Save()
	if err != nil {
		return err
	}
	// set the reputation for the recipient here
	return database.ChangeRecpReputation(a.ProjectRecipient.U.Index, a.Params.TotalValue*RecipientWeight)
}
