package opensolar

import (
	"fmt"
	"log"
	"time"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	issuer "github.com/YaleOpenLab/openx/issuer"
	model "github.com/YaleOpenLab/openx/models/munibond"
	notif "github.com/YaleOpenLab/openx/notif"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
)

// this file does not contain any tests associated with it right now. In the future,
// once we have a robust frontend, we can modify the CLI interface to act as a test
// for this file

// PreInvestmentCheck is exclusive to a particular platform and should be defined for each platform
// performing its own set of checks before allowing an investor to invest in the project.
func PreInvestmentCheck(projIndex int, invIndex int, invAmount string) (Project, database.Investor, error) {
	var project Project
	var investor database.Investor
	var err error

	project, err = RetrieveProject(projIndex)
	if err != nil {
		return project, investor, err
	}

	investor, err = database.RetrieveInvestor(invIndex)
	if err != nil {
		return project, investor, err
	}

	if !investor.CanInvest(invAmount) {
		log.Println("Investor has less balance than what is required to ivnest in this asset")
		return project, investor, fmt.Errorf("Investor has less balance than what is required to ivnest in this asset")
	}

	// check if investment amount is greater than or equal to the project requirements
	if utils.StoF(invAmount) > project.Params.TotalValue-project.Params.MoneyRaised {
		return project, investor, err
	}

	// user has decided to invest in a part of the project (don't know if full yet)
	// no asset codes assigned yet, we need to create them
	// you can retrieve asetCodes anywhere since metadata is assumed to be unique
	if project.Params.SeedAssetCode == "" && project.Params.InvestorAssetCode == "" {
		// this project does not have an issuer associated with it yet since there has been
		// no seed round and an investment round
		project.Params.InvestorAssetCode = assets.AssetID(consts.InvestorAssetPrefix + project.Params.Metadata) // set the investor asset code
		err = project.Save()
		if err != nil {
			return project, investor, err
		}
		err = issuer.InitIssuer(projIndex, consts.IssuerSeedPwd)
		if err != nil {
			return project, investor, err
		}
		err = issuer.FundIssuer(projIndex, consts.IssuerSeedPwd, consts.PlatformSeed)
		if err != nil {
			return project, investor, err
		}
	}

	return project, investor, nil
}

// UpdateProjectAfterInvestment is a required function to be defined on all platforms which
// requires the invested asset to be updated shortly after the investment
func (project *Project) UpdateProjectAfterInvestment(invAmount string, investor database.Investor) error {

	var err error
	project.Params.MoneyRaised += utils.StoF(invAmount)
	project.ProjectInvestors = append(project.ProjectInvestors, investor)
	// keep note of who invested in this asset (even though it should be easy
	// to get that from the blockchain)

	err = project.Save()
	if err != nil {
		return err
	}

	if project.Params.MoneyRaised == project.Params.TotalValue {
		project.Lock = true
		err = project.Save()
		if err != nil {
			return err
		}
		project.sendRecipientNotification()
		go sendRecipientAssets(project.Params.Index)
	}

	return nil
}

// InvestInProject invests in a particular solar project given required parameters
func InvestInProject(projIndex int, invIndex int, invAmount string, invSeed string) error {
	var err error

	project, investor, err := PreInvestmentCheck(projIndex, invIndex, invAmount)
	if err != nil {
		return err
	}

	err = model.SingleInvestmentModelInv(investor, invSeed, invAmount, projIndex,
		project.Params.InvestorAssetCode, project.Params.TotalValue)
	if err != nil {
		return err
	}

	err = project.UpdateProjectAfterInvestment(invAmount, investor)
	if err != nil {
		return err
	}

	return err
}

// SeedInvestInProject is similar to InvestInProject differing only in that it distributes
// seed assets instead of investor assets
func SeedInvestInProject(projIndex int, invIndex int, recpIndex int, invAmount string,
	invSeed string, recpSeed string) error {

	project, investor, err := PreInvestmentCheck(projIndex, invIndex, invAmount)
	if err != nil {
		return err
	}

	err = model.SingleInvestmentModelInv(investor, invSeed, invAmount, projIndex,
		project.Params.SeedAssetCode, project.Params.TotalValue)
	if err != nil {
		return err
	}

	err = project.UpdateProjectAfterInvestment(invAmount, investor)
	if err != nil {
		return err
	}

	return err
}

func (project *Project) sendRecipientNotification() {
	// this project covers up the amount nedeed for the project, so send the recipient
	// a notification that their project has been invested in and that they need
	// to logon to the platform in order to accept the investment
	// notif.SendUnlockNotifToRecipient(projIndex, project.ProjectRecipient.U.Email)
	notif.SendUnlockNotifToRecipient(project.Params.Index, project.ProjectRecipient.U.Email)
}

func (project *Project) UpdateProjectAfterAcceptance() error {

	recipient, err := database.RetrieveRecipient(project.ProjectRecipient.U.Index)
	if err != nil {
		return err
	}

	project.Params.BalLeft = float64(project.Params.TotalValue)
	project.ProjectRecipient = recipient // need to udpate project.Params each time recipient is mutated
	project.Stage = FundedProject        // set funded project stage

	err = project.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	// need to run a separate goroutine for payback
	go monitorPaybacks(recipient.U.Index, project.Params.Index, project.Params.DebtAssetCode)
	return nil
}

func sendRecipientAssets(projIndex int) error {
	// this project covers up the amount nedeed for the project, so set the DebtAssetCode
	// and PaybackAssetCodes, generate them and give to the recipient
	// we need the recipient's seed here, so we need to wait on the frontend and require
	// confirmation from the recipient or something
	// we need the recipient's seed before we can proceed further
	startTime := utils.Unix()
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}

	for utils.Unix()-startTime < consts.LockInterval {
		log.Printf("WAITING FOR PROJECT %d TO BE UNLOCKED", projIndex)
		project, err = RetrieveProject(projIndex)
		if err != nil {
			return err
		}
		if !project.Lock {
			log.Println("Project UNLOCKED IN LOOP")
			break
		}
		time.Sleep(10 * time.Second)
	}

	// here, we hope that the recipient's account is setup already
	// by the time the function reaches here, the lock would have been opened
	// update our copy of the project
	project, err = RetrieveProject(projIndex)
	if err != nil {
		log.Println(err)
		return err
	}

	recpSeed, err := wallet.DecryptSeed(project.ProjectRecipient.U.EncryptedSeed, project.LockPwd)
	if err != nil {
		log.Println(err)
		return err
	}

	metadata := project.Params.Metadata

	project.Params.DebtAssetCode = assets.AssetID(consts.DebtAssetPrefix + metadata)
	project.Params.PaybackAssetCode = assets.AssetID(consts.PaybackAssetPrefix + metadata)

	err = model.SingleInvestmentModelRecp(project.ProjectRecipient, projIndex, project.Params.DebtAssetCode,
		project.Params.PaybackAssetCode, project.Params.Years, recpSeed, project.Params.TotalValue, project.PaybackPeriod)
	if err != nil {
		return err
	}

	err = project.UpdateProjectAfterAcceptance()
	if err != nil {
		return err
	}

	return nil
}
