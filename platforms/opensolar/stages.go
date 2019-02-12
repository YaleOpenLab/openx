package opensolar

import (
	database "github.com/YaleOpenLab/openx/database"
	"log"
)

// this file contains the different stages associated with an opensolar project and the handlers
// that are used to modify them

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
	a.Reputation = a.TotalValue // upgrade reputation since totalValue might have changed from the originated contract
	err := a.Save()
	if err != nil {
		return err
	}
	return RepOriginatedProject(a.Originator.U.Index, a.Index) // modify originator reputation now that the final price is fixed
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

	err = a.Contractor.U.IncreaseReputation(a.TotalValue * ContractorWeight) // modify contractor Reputation now that a project has been installed
	if err != nil {
		return err
	}

	for _, i := range a.InvestorIndices {
		elem, err := database.RetrieveInvestor(i)
		if err != nil {
			log.Println(err)
			return err
		}
		err = database.ChangeInvReputation(elem.U.Index, a.TotalValue*InvestorWeight)
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

	return database.ChangeRecpReputation(a.RecipientIndex, a.TotalValue*RecipientWeight) // modify recipient reputation now that the system had begun power generation
}
