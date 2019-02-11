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

// the smart contract that powers this particular platform. Designed to be monolithic by design
// so that we can have everything we automate in one place for easy audits.

// these are the reputation weights associated with a project on the opensolar platform. For eg if
// a project's total worth is 10000 and everything in the project goes well and
// all entities are satisfied by the outcome, the originator gains 1000 points,
// the contractor gains 3000 points and so on
// MWTODO: get comments on the weights and tweak them if needed
const (
	InvestorWeight         = 0.1 // the percentage weight of the project's total reputation assigned to the investor
	OriginatorWeight       = 0.1 // the percentage weight of the project's total reputation assigned to the originator
	ContractorWeight       = 0.3 // the percentage weight of the project's total reputation assigned to the contractor
	DeveloperWeight        = 0.2 // the percentage weight of the project's total reputation assigned to the developer
	RecipientWeight        = 0.3 // the percentage weight of the project's total reputation assigned to the recipient
	NormalThreshold        = 1   // NormalThreshold is the normal payback interval of 1 payback period. Regular notifications are sent regardless of whether the user has paid back towards the project.
	AlertThreshold         = 2   // AlertThreshold is the threshold above which the user gets a nice email requesting a quick payback whenever possible
	SternAlertThreshold    = 4   // SternAlertThreshold is the threshold above the user gets a warning that services will be disconnected if the user doesn't payback soon
	DisconnectionThreshold = 6   // DisconnectionThreshold is the threshold above which the user gets a notification telling that services have been disconnected.
)

// TODO: Consider that for this authorization to happen, there could be a
// verification requirement (eg. that the project is relatively feasible),
// and that it may need several approvals for it (eg. Recipient can be two
// figures here — the school entity (more visible) and the department of
// education (more admin) who is the actual issuer) along with a validation
// requirement

// VerifyBeforeAuthorizing verifies some information on the originator before upgrading
// the project stage
func VerifyBeforeAuthorizing(projIndex int) bool {
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return false
	}
	// TODO: In the future, this would involve the kyc operator to check the originator's credentials
	fmt.Printf("ORIGINATOR'S NAME IS: %s and PROJECT's METADATA IS: %s", project.Originator.U.Name, project.Metadata)
	return true
}

// RecipientAuthorize allows a recipient to authorize a specific project
func RecipientAuthorize(projIndex int, recpIndex int) error {
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}
	if project.Stage != 0 {
		return fmt.Errorf("Project stage not zero")
	}
	if !VerifyBeforeAuthorizing(projIndex) {
		return fmt.Errorf("Originator not verified")
	}
	recipient, err := database.RetrieveRecipient(recpIndex)
	if err != nil {
		return err
	}
	if project.RecipientIndex != recipient.U.Index {
		return fmt.Errorf("You can't authorize a project which is not assigned to you!")
	}

	err = project.SetOriginProject() // set the project as originated
	if err != nil {
		return err
	}

	err = RepOriginatedProject(project.Originator.U.Index, project.Index)
	if err != nil {
		return err
	}

	/* set the open for money stage if we choose to have it in the end
	err = project.SetOpenForMoneyStage()
	if err != nil {
		return err
	}
	*/
	return nil
}

// VoteTowardsProposedProject is a handler that an investor would use to vote towards a
// specific proposed project on the platform
func VoteTowardsProposedProject(invIndex int, votes int, projectIndex int) error {
	inv, err := database.RetrieveInvestor(invIndex)
	if err != nil {
		return err
	}
	if votes > inv.VotingBalance {
		return fmt.Errorf("Can't vote with an amount greater than available balance")
	}

	project, err := RetrieveProject(projectIndex)
	if err != nil {
		return err
	}
	if project.Stage != 2 {
		return fmt.Errorf("You can't vote for a project with stage less than 2")
	}

	project.Votes += votes
	err = project.Save()
	if err != nil {
		return err
	}

	err = inv.DeductVotingBalance(votes)
	if err != nil {
		return err
	}

	fmt.Println("CAST VOTE TOWARDS PROJECT SUCCESSFULLY")
	return nil
}

// the preInvestmentChecks associated with the opensolar platform
func preInvestmentCheck(projIndex int, invIndex int, invAmount string) (Project, error) {
	var project Project
	var investor database.Investor
	var err error

	project, err = RetrieveProject(projIndex)
	if err != nil {
		return project, err
	}

	investor, err = database.RetrieveInvestor(invIndex)
	if err != nil {
		return project, err
	}

	if !investor.CanInvest(invAmount) {
		log.Println("Investor has less balance than what is required to ivnest in this asset")
		return project, fmt.Errorf("Investor has less balance than what is required to ivnest in this asset")
	}

	// check if investment amount is greater than or equal to the project requirements
	if utils.StoF(invAmount) > project.TotalValue-project.MoneyRaised {
		return project, err
	}

	// user has decided to invest in a part of the project (don't know if full yet)
	// no asset codes assigned yet, we need to create them
	if project.SeedAssetCode == "" && project.InvestorAssetCode == "" {
		// this project does not have an issuer associated with it yet since there has been
		// no seed round and an investment round
		project.InvestorAssetCode = assets.AssetID(consts.InvestorAssetPrefix + project.Metadata) // you can retrieve asetCodes anywhere since metadata is assumed to be unique
		err = project.Save()
		if err != nil {
			return project, err
		}
		err = issuer.InitIssuer(projIndex, consts.IssuerSeedPwd)
		if err != nil {
			return project, err
		}
		err = issuer.FundIssuer(projIndex, consts.IssuerSeedPwd, consts.PlatformSeed)
		if err != nil {
			return project, err
		}
	}

	return project, nil
}

// the SeedInvest function of the opensolar platform
func SeedInvest(projIndex int, invIndex int, recpIndex int, invAmount string,
	invSeed string, recpSeed string) error {

	project, err := preInvestmentCheck(projIndex, invIndex, invAmount)
	if err != nil {
		return err
	}

	err = model.MunibondInvest(invIndex, invSeed, invAmount, projIndex,
		project.SeedAssetCode, project.TotalValue)
	if err != nil {
		return err
	}

	err = project.updateProjectAfterInvestment(invAmount, invIndex)
	if err != nil {
		return err
	}

	return err
}

// the main invest function of the opensolar platform
func Invest(projIndex int, invIndex int, invAmount string, invSeed string) error {
	var err error

	project, err := preInvestmentCheck(projIndex, invIndex, invAmount)
	if err != nil {
		return err
	}

	err = model.MunibondInvest(invIndex, invSeed, invAmount, projIndex,
		project.InvestorAssetCode, project.TotalValue)
	if err != nil {
		return err
	}

	err = project.updateProjectAfterInvestment(invAmount, invIndex)
	if err != nil {
		return err
	}

	return err
}

// the updateProjectAfterInvestment of the opensolar platform
func (project *Project) updateProjectAfterInvestment(invAmount string, invIndex int) error {

	var err error
	project.MoneyRaised += utils.StoF(invAmount)
	project.InvestorIndices = append(project.InvestorIndices, invIndex)
	err = project.Save()
	if err != nil {
		return err
	}

	if project.MoneyRaised == project.TotalValue {
		project.Lock = true
		err = project.Save()
		if err != nil {
			return err
		}

		err = project.sendRecipientNotification()
		if err != nil {
			return err
		}

		go sendRecipientAssets(project.Index)
	}

	return nil
}

// sendRecipientNotification sends the notification to the recipient requesting them
// to logon to the platform and unlock the project that has just been invested in
func (project *Project) sendRecipientNotification() error {
	recipient, err := database.RetrieveRecipient(project.RecipientIndex)
	if err != nil {
		return err
	}
	notif.SendUnlockNotifToRecipient(project.Index, recipient.U.Email)
	return nil
}

// UnlockProject unlocks a specific project that has just been invested in
func UnlockProject(username string, pwhash string, projIndex int, seedpwd string) error {
	fmt.Println("UNLOCKING PROJECT")
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}

	recipient, err := database.ValidateRecipient(username, pwhash)
	if err != nil {
		return err
	}

	if recipient.U.Index != project.RecipientIndex {
		log.Println("CHECKINDEICE", recipient.U.Index, project.RecipientIndex)
		return fmt.Errorf("Seeds don't match, quitting!")
	}

	recpSeed, err := wallet.DecryptSeed(recipient.U.EncryptedSeed, seedpwd)
	if err != nil {
		log.Println(err)
		return err
	}

	checkPubkey, err := wallet.ReturnPubkey(recpSeed)
	if err != nil {
		log.Println(err)
		return err
	}

	if checkPubkey != recipient.U.PublicKey {
		log.Println("Invalid seed")
		return fmt.Errorf("Failed to unlock project")
	}

	if !project.Lock {
		return fmt.Errorf("Project not locked")
	}

	project.LockPwd = seedpwd
	project.Lock = false
	err = project.Save()
	if err != nil {
		return err
	}
	return nil
}

// sendRecipientAssets sends a recipient the debt asset and the payback asset associated with
// the opensolar platform
func sendRecipientAssets(projIndex int) error {
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

	// lock is open, retrieve project and transfer assets
	project, err = RetrieveProject(projIndex)
	if err != nil {
		log.Println(err)
		return err
	}

	recipient, err := database.RetrieveRecipient(project.RecipientIndex)
	if err != nil {
		log.Println(err)
		return err
	}

	recpSeed, err := wallet.DecryptSeed(recipient.U.EncryptedSeed, project.LockPwd)
	if err != nil {
		log.Println(err)
		return err
	}

	project.LockPwd = "" // set lockpwd to nil immediately after retrieving seed
	metadata := project.Metadata

	project.DebtAssetCode = assets.AssetID(consts.DebtAssetPrefix + metadata)
	project.PaybackAssetCode = assets.AssetID(consts.PaybackAssetPrefix + metadata)

	err = model.MunibondReceive(project.RecipientIndex, projIndex, project.DebtAssetCode,
		project.PaybackAssetCode, project.Years, recpSeed, project.TotalValue, project.PaybackPeriod)
	if err != nil {
		return err
	}

	err = project.updateProjectAfterAcceptance()
	if err != nil {
		return err
	}

	return nil
}

// updateProjectAfterAcceptance updates the project after acceptance of investment
// by the recipient
func (project *Project) updateProjectAfterAcceptance() error {

	project.BalLeft = float64(project.TotalValue)
	project.Stage = FundedProject // set funded project stage

	err := project.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	go monitorPaybacks(project.RecipientIndex, project.Index)
	return nil
}

// Payback pays the platform back in STABLEUSD and DebtAsset and receives PaybackAssets
// in return. Price to be paid per month depends on the electricity consumed by the recipient
// in the particular time frame
// If we allow a user to hold balances in btc / xlm, we could direct them to exchange the coin for STABLEUSD
// (or we could setup a payment provider which accepts fiat + crypto and do this ourselves)
func Payback(recpIndex int, projIndex int, assetName string, amount string, recipientSeed string) error {

	project, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}

	err = model.MunibondPayback(recpIndex, amount, recipientSeed, projIndex, assetName, project.InvestorIndices)
	if err != nil {
		return err
	}

	project.BalLeft -= utils.StoF(amount) // can directly change this since we've checked for it in the MunibondPayback call
	project.DateLastPaid = utils.Unix()
	if project.BalLeft == 0 {
		log.Println("YOU HAVE PAID OFF THIS ASSET, TRANSFERRING OWNERSHIP OF ASSET TO YOU")
		project.Stage = 7
		// we should call neighborly or some other partner here to transfer assets using the bond they provide us with
		// the nice part here is that the recipient can not pay off more than what is
		// invested because the trustline will not allow such an incident to happen
	}

	err = project.Save()
	if err != nil {
		return err
	}

	return err
}

// CalculatePayback calculates the amount of payback assets that must be issued in relation
// to the total amount invested in the project
func (project Project) CalculatePayback(amount string) string {
	amountF := utils.StoF(amount)
	amountPB := (amountF / float64(project.TotalValue)) * float64(project.Years*12)
	amountPBString := utils.FtoS(amountPB)
	return amountPBString
}

// monitorPaybacks monitors whether the user is paying back regularly towards the given project
// thread has to be isolated since if this fails, we stop tracking paybacks by the recipient.
func monitorPaybacks(recpIndex int, projIndex int) {
	for {
		project, err := RetrieveProject(projIndex)
		if err != nil {
			log.Println(err)
		}

		recipient, err := database.RetrieveRecipient(recpIndex)
		if err != nil {
			log.Println(err)
		}

		// this will be our payback period and we need to check if the user pays us back

		nowTime := utils.Unix()
		timeElapsed := nowTime - project.DateLastPaid                   // this would be in seconds (unix time)
		period := int64(project.PaybackPeriod * consts.OneWeekInSecond) // in seconds due to the const
		if period == 0 {
			period = 1 // for the test suite
		}
		factor := timeElapsed / period

		if factor <= 1 {
			// don't do anything since the suer has been paying back regularly
			log.Println("User: ", recipient.U.Email, "is on track paying towards order: ", projIndex)
			// maybe even update reputation here on a fractional basis depending on a user's timely payments
		} else if factor > NormalThreshold && factor < AlertThreshold {
			// person has not paid back for one-two consecutive period, send gentle reminder
			notif.SendNicePaybackAlertEmail(projIndex, recipient.U.Email)
		} else if factor >= SternAlertThreshold && factor < DisconnectionThreshold {
			// person has not paid back for four consecutive cycles, send reminder
			notif.SendSternPaybackAlertEmail(projIndex, recipient.U.Email)
			for _, i := range project.InvestorIndices {
				// send an email to recipients to assure them that we're on the issue and will be acting
				// soon if the recipient fails to pay again.
				investor, err := database.RetrieveInvestor(i)
				if err != nil {
					log.Println(err)
					continue
				}
				if investor.U.Notification {
					notif.SendSternPaybackAlertEmailI(projIndex, investor.U.Email)
				}
			}
			notif.SendSternPaybackAlertEmailG(projIndex, project.Guarantor.U.Email)
			// send an email out to the guarantor
		} else if factor >= DisconnectionThreshold {
			// send a disconnection notice to the recipient and let them know we have redirected
			// power towards the grid. Also maybe email ourselves in this case so that we can
			// contact them personally to resolve the issue as soon as possible.
			notif.SendDisconnectionEmail(projIndex, recipient.U.Email)
			for _, i := range project.InvestorIndices {
				// send an email to recipients to assure them that we're on the issue and will be acting
				// soon if the recipient fails to pay again.
				investor, err := database.RetrieveInvestor(i)
				if err != nil {
					log.Println(err)
					continue
				}
				if investor.U.Notification {
					notif.SendDisconnectionEmailI(projIndex, investor.U.Email)
				}
			}
			notif.SendDisconnectionEmailG(projIndex, project.Guarantor.U.Email)
			// send an email out to the guarantor
		}

		time.Sleep(time.Duration(consts.OneWeekInSecond)) // poll every week to ch eck progress on payments
	}
}
