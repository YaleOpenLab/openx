package opensolar

import (
	"fmt"
	"github.com/pkg/errors"
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
	xlm "github.com/YaleOpenLab/openx/xlm"
)

// This script represents the smart contract that powers a project in this particular platform. Designed to be monolithic by design
// so that we can have everything we automate in one place for easy audits.

// —ACTOR REPUTATION SCHEMES—

// 1. Automatic Contants Based on Stage Completion and Project Success:
// These constants are the reputation weights associated with a project on the opensolar platform. For eg if
// a project's total worth is 10000 and everything in the project goes well and
// all entities are satisfied by the outcome, the originator gains 1000 points,
// the contractor gains 3000 points and so on. MW: These are allocated at what point in terms of the project stages? They will have to vary
// Thresholds relate to the payment cycles owed by the Recipient. MW: How are these executed, and how are points added or removes? Its unclear
const (
	InvestorWeight         = 0.1 // the percentage weight of the project's total reputation assigned to the investor
	OriginatorWeight       = 0.1 // the percentage weight of the project's total reputation assigned to the originator
	ContractorWeight       = 0.3 // the percentage weight of the project's total reputation assigned to the contractor
	DeveloperWeight        = 0.2 // the percentage weight of the project's total reputation assigned to the developer
	RecipientWeight        = 0.3 // the percentage weight of the project's total reputation assigned to the recipient
	NormalThreshold        = 1   // NormalThreshold is the normal payback interval of 1 payback period. Regular notifications are sent regardless of whether the user has paid back towards the project.
	AlertThreshold         = 2   // AlertThreshold is the threshold above which the user gets a nice email requesting a quick payback whenever possible
	SternAlertThreshold    = 4   // SternAlertThreshold is the threshold above when the user gets a warning that services will be disconnected if the user doesn't payback soon.
	DisconnectionThreshold = 6   // DisconnectionThreshold is the threshold above which the user gets a notification telling that services have been disconnected.
)

// TODO:
// 2. Peer-based star-rating
// This should be the normal 5 star system that users get from other users that are involved in the same transaction.


// TODO: Consider that in the family of Recipients or Investors, there are more than one actor, and sometimes signatory authorization is from only some of the actors.
// See the Project Stages document for reference of Beneficiary or Investor families. A clear example is a Recipient that is the actual issuer of the security,
// and another that is the actual offtaker.

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
		return errors.Wrap(err, "couldn't retrieve project")
	}
	if project.Stage != 0 {
		return errors.New("Project stage not zero")
	}
	if !VerifyBeforeAuthorizing(projIndex) {
		return errors.New("Originator not verified")
	}
	recipient, err := database.RetrieveRecipient(recpIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve recipient")
	}
	if project.RecipientIndex != recipient.U.Index {
		return errors.New("You can't authorize a project which is not assigned to you!")
	}

	err = project.SetStage(1) // set the project as originated
	if err != nil {
		return errors.Wrap(err, "Error while setting origin project")
	}

	err = RepOriginatedProject(project.Originator.U.Index, project.Index)
	if err != nil {
		return errors.Wrap(err, "error while increasing reputation of originator")
	}

	return nil
}

// —VOTING SCHEMES—
// MW: Lets design this together. Very cool to have votes (which are 'Likes'), but why only investors can vote? Why not projects at stage 1?
// What does it mean if a project has high votes?

// VoteTowardsProposedProject is a handler that an investor would use to vote towards a
// specific proposed project on the platform.
func VoteTowardsProposedProject(invIndex int, votes int, projectIndex int) error {
	inv, err := database.RetrieveInvestor(invIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve investor")
	}
	if votes > inv.VotingBalance {
		return errors.New("Can't vote with an amount greater than available balance")
	}

	project, err := RetrieveProject(projectIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	if project.Stage != 2 {
		return errors.New("You can't vote for a project with stage less than 2")
	}

	project.Votes += votes
	err = project.Save()
	if err != nil {
		return errors.Wrap(err, "couldn't save project")
	}

	err = inv.DeductVotingBalance(votes)
	if err != nil {
		return errors.Wrap(err, "error while deducitng voting balance of investor")
	}

	fmt.Println("CAST VOTE TOWARDS PROJECT SUCCESSFULLY")
	return nil
}

// -- INVESTMENT VERIFICATIONS--
// the preInvestmentChecks associated with the opensolar platform when an Investor bids an investment amount of a specific project
func preInvestmentCheck(projIndex int, invIndex int, invAmount string) (Project, error) {
	var project Project
	var investor database.Investor
	var err error

	project, err = RetrieveProject(projIndex)
	if err != nil {
		return project, errors.Wrap(err, "couldn't retrieve project")
	}

	investor, err = database.RetrieveInvestor(invIndex)
	if err != nil {
		return project, errors.Wrap(err, "couldn't retrieve investor")
	}

	// here we should check whether the investor has adequate STABLEUSD or XLM and not just the stablecoin
	// since we automate asset conversion in the MunibondInvest function
	if !investor.CanInvest(invAmount) {
		return project, errors.New("Investor has less balance than what is required to invest in this project")
	}

	// check if investment amount is greater than or equal to the project requirements
	if utils.StoF(invAmount) > project.TotalValue-project.MoneyRaised {
		return project, errors.New("Investment amount greater than what is required! Adjust your investment")
	}

	if project.SeedAssetCode == "" && project.InvestorAssetCode == "" {
		// this project does not have an asset issuer associated with it yet since there has been
		// no seed round nor investment round
		project.InvestorAssetCode = assets.AssetID(consts.InvestorAssetPrefix + project.Metadata) // you can retrieve asetCodes anywhere since metadata is assumed to be unique
		err = project.Save()
		if err != nil {
			return project, errors.Wrap(err, "couldn't save project")
		}
		err = issuer.InitIssuer(consts.OpenSolarIssuerDir, projIndex, consts.IssuerSeedPwd)
		if err != nil {
			return project, errors.Wrap(err, "error while initializing issuer")
		}
		err = issuer.FundIssuer(consts.OpenSolarIssuerDir, projIndex, consts.IssuerSeedPwd, consts.PlatformSeed)
		if err != nil {
			return project, errors.Wrap(err, "error while funding issuer")
		}
	}

	return project, nil
}

// --SEED INVESTMENT--
// SeedInvest is the seed investment function of the opensolar platform
func SeedInvest(projIndex int, invIndex int, recpIndex int, invAmount string,
	invSeed string, recpSeed string) error {

	project, err := preInvestmentCheck(projIndex, invIndex, invAmount)
	if err != nil {
		return errors.Wrap(err, "error while performing pre investment check")
	}

	// MW: Consider other seed investments in stages before the big raise of stage 4
	if project.Stage != 1 {
		return fmt.Errorf("project stage not at 1, you either have passed the seed stage or project is not at seed stage yet")
	}

// MW: Here it is using a specific model investment, eg. Muni Bond. If this is hard coded here, how can you set an opensolar project as equity crowdfunding or bond or debt?
	if project.InvestmentType != "munibond" {
		return fmt.Errorf("other investment models are not supported right now, quitting")
	}

	err = model.MunibondInvest(consts.OpenSolarIssuerDir, invIndex, invSeed, invAmount, projIndex,
		project.SeedAssetCode, project.TotalValue)
	if err != nil {
		return errors.Wrap(err, "error while investing")
	}

	err = project.updateProjectAfterInvestment(invAmount, invIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't update project after investment")
	}

	return err
}

// -- INVEST --
// Invest is the main invest function of the opensolar platform
func Invest(projIndex int, invIndex int, invAmount string, invSeed string) error {
	var err error

	// run preinvestment checks to make sure everything is okay
	project, err := preInvestmentCheck(projIndex, invIndex, invAmount)
	if err != nil {
		return errors.Wrap(err, "pre investment check failed")
	}

	if project.InvestmentType != "munibond" {
		return fmt.Errorf("other investment models are not supported right now, quitting")
	}

	if project.Stage != 4 {
		return fmt.Errorf("project not at stage where it can solicit investment, quitting")
	}
	// call the model and invest in the particular project
	err = model.MunibondInvest(consts.OpenSolarIssuerDir, invIndex, invSeed, invAmount, projIndex,
		project.InvestorAssetCode, project.TotalValue)
	if err != nil {
		log.Println("Error while seed investing", err)
		return errors.Wrap(err, "error while seed investing")
	}

	// once the investment is complete, update the project and store in the database
	err = project.updateProjectAfterInvestment(invAmount, invIndex)
	if err != nil {
		return errors.Wrap(err, "failed to update project after investment")
	}

	return err
}

// the updateProjectAfterInvestment of the opensolar platform
func (project *Project) updateProjectAfterInvestment(invAmount string, invIndex int) error {
	// MW: It seems that all your messages strings relate to errors, but not to confirmed transactions. It would be useful to add those
	var err error
	project.MoneyRaised += utils.StoF(invAmount)
	project.InvestorIndices = append(project.InvestorIndices, invIndex)
	err = project.Save()
	if err != nil {
		return errors.Wrap(err, "couldn't save project")
	}

	if project.MoneyRaised == project.TotalValue {
		project.Lock = true
		err = project.Save()
		if err != nil {
			return errors.Wrap(err, "couldn't save project")
		}

		err = project.sendRecipientNotification()
		if err != nil {
			return errors.Wrap(err, "error while sending notifications to recipient")
		}

		err = InitEscrow(consts.EscrowDir, project.Index, consts.EscrowPwd)
		if err != nil {
			return errors.Wrap(err, "error while initializing issuer")
		}

		err = TransferFundsToEscrow(project.TotalValue, project.Index)
		if err != nil {
			log.Println(err)
			return errors.Wrap(err, "could not transfer funds to the escrow, quitting!")
		}

		go sendRecipientAssets(project.Index)
	}

	// we need to udpate the project investment map here
	project.InvestorMap = make(map[string]float64) // make the map

	for _, elem := range project.InvestorIndices {
		investor, err := database.RetrieveInvestor(elem)
		if err != nil {
			return errors.Wrap(err, "error while retrieving investors, quitting")
		}

		balanceS, err := xlm.GetAssetBalance(investor.U.PublicKey, project.InvestorAssetCode)
		if err != nil {
			return errors.Wrap(err, "error while retrieving asset balance, quitting")
		}

		balanceF, err := utils.StoFWithCheck(balanceS)
		if err != nil {
			return errors.Wrap(err, "error while converting to float, quitting")
		}

		percentageInvestment := balanceF / project.TotalValue
		project.InvestorMap[investor.U.PublicKey] = percentageInvestment
	}

	err = project.Save()
	log.Println("INVESTOR MAP: ", project.InvestorMap)
	if err != nil {
		return errors.Wrap(err, "error while saving project, quitting")
	}
	return nil
}

// MW: Why does the recipient have to unlock the project? Why is the project locked in the first place?
// sendRecipientNotification sends the notification to the recipient requesting them
// to logon to the platform and unlock the project that has just been invested in
func (project *Project) sendRecipientNotification() error {
	recipient, err := database.RetrieveRecipient(project.RecipientIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve recipient")
	}
	notif.SendUnlockNotifToRecipient(project.Index, recipient.U.Email)
	return nil
}

// UnlockProject unlocks a specific project that has just been invested in
func UnlockProject(username string, pwhash string, projIndex int, seedpwd string) error {
	fmt.Println("UNLOCKING PROJECT")
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}

	recipient, err := database.ValidateRecipient(username, pwhash)
	if err != nil {
		return errors.Wrap(err, "couldn't validate recipient")
	}

	if recipient.U.Index != project.RecipientIndex {
		return errors.New("Recipient Indices don't match, quitting!")
	}

	recpSeed, err := wallet.DecryptSeed(recipient.U.EncryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "error while decrpyting seed")
	}

	checkPubkey, err := wallet.ReturnPubkey(recpSeed)
	if err != nil {
		return errors.Wrap(err, "couldn't get public key from seed")
	}

	if checkPubkey != recipient.U.PublicKey {
		log.Println("Invalid seed")
		return errors.New("Failed to unlock project")
	}

	if !project.Lock {
		return errors.New("Project not locked")
	}

	project.LockPwd = seedpwd
	project.Lock = false
	err = project.Save()
	if err != nil {
		return errors.Wrap(err, "couldn't save project")
	}
	return nil
}

// sendRecipientAssets sends a recipient the debt asset and the payback asset associated with
// the opensolar platform
func sendRecipientAssets(projIndex int) error {
	startTime := utils.Unix()
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "Couldn't retrieve project")
	}

	for utils.Unix()-startTime < consts.LockInterval {
		log.Printf("WAITING FOR PROJECT %d TO BE UNLOCKED", projIndex)
		project, err = RetrieveProject(projIndex)
		if err != nil {
			return errors.Wrap(err, "Couldn't retrieve project")
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
		return errors.Wrap(err, "Couldn't retrieve project")
	}

	recipient, err := database.RetrieveRecipient(project.RecipientIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve recipienrt")
	}

	recpSeed, err := wallet.DecryptSeed(recipient.U.EncryptedSeed, project.LockPwd)
	if err != nil {
		return errors.Wrap(err, "couldn't decrypt seed")
	}

	project.LockPwd = "" // set lockpwd to nil immediately after retrieving seed
	metadata := project.Metadata

	project.DebtAssetCode = assets.AssetID(consts.DebtAssetPrefix + metadata)
	project.PaybackAssetCode = assets.AssetID(consts.PaybackAssetPrefix + metadata)

	err = model.MunibondReceive(consts.OpenSolarIssuerDir, project.RecipientIndex, projIndex, project.DebtAssetCode,
		project.PaybackAssetCode, project.Years, recpSeed, project.TotalValue, project.PaybackPeriod)
	if err != nil {
		return errors.Wrap(err, "error while receiving assets from issuer on recipient's end")
	}

	err = project.updateProjectAfterAcceptance()
	if err != nil {
		return errors.Wrap(err, "failed to update project after acceptance of asset")
	}

	return nil
}

// - PROJECT INVESTMENT UPDATES THROUGHOUT 'THE RAISE' IN STAGE 7 --
// updateProjectAfterAcceptance updates the project after acceptance of investment
// by the recipient
func (project *Project) updateProjectAfterAcceptance() error {

	project.BalLeft = float64(project.TotalValue)
	project.Stage = Stage5.Number // set to stage 5 (after the raise is done, we need to wait for people to actually construct the solar panels)

	err := project.Save()
	if err != nil {
		return errors.Wrap(err, "couln't save project")
	}

	go monitorPaybacks(project.RecipientIndex, project.Index)
	return nil
}

// MW: Here, the project jumps from stage 5, the end to the raise, to stage 7, the payback period. What happens to everything in between?

// -- SOLAR OFFTAKING PAYMENTS IN STAGE 7 --
// Payback pays the platform back in STABLEUSD and DebtAsset and receives PaybackAssets
// in return. Price to be paid per month depends on the electricity consumed by the recipient
// in the particular time frame
// If we allow a user to hold balances in btc / xlm, we could direct them to exchange the coin for STABLEUSD
// (or we could setup a payment provider which accepts fiat + crypto and do this ourselves)
func Payback(recpIndex int, projIndex int, assetName string, amount string, recipientSeed string) error {

	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "Couldn't retrieve project")
	}

	escrowPath := CreatePath(consts.EscrowDir, projIndex)

	if project.InvestmentType != "munibond" {
		return fmt.Errorf("other investment models are not supported right now, quitting")
	}

	err = model.MunibondPayback(consts.OpenSolarIssuerDir, escrowPath, recpIndex, amount, recipientSeed, projIndex, assetName, project.InvestorIndices)
	if err != nil {
		return errors.Wrap(err, "Error while paying back the issuer")
	}

	// MW: Ownership of asset could shift as payments happen, or flip at the end.
	// 		Also, wouldnt it make sense to make the 'Ownership Flip or Handoff' as a separate function? Since this will have to trigger changes in a registry?
	project.BalLeft -= utils.StoF(amount) // can directly change this since we've checked for it in the MunibondPayback call
	project.DateLastPaid = utils.Unix()
	if project.BalLeft == 0 {
		log.Println("YOU HAVE PAID OFF THIS ASSET, TRANSFERRING OWNERSHIP OF ASSET TO YOU")
		project.Stage = 9 // stage 9 is the disposal stage, we don't wait for stage 9 to complete and hence leave it as is, just deleting the account and stuff associated with the project
		// we should call neighborly or some other partner here to transfer assets using the bond they provide us with
		// the nice part here is that the recipient can not pay off more than what is
		// invested because the trustline will not allow such an incident to happen
	}

	err = project.Save()
	if err != nil {
		return errors.Wrap(err, "coudln't save project")
	}

	escrowPubkey, escrowSeed, err := wallet.RetrieveSeed(escrowPath, consts.EscrowPwd)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve issuer seed")
	}

	err = DistributePayments(escrowSeed, escrowPubkey, projIndex, utils.StoI(amount))
	if err != nil {
		return errors.Wrap(err, "error while distributing payments")
	}

	return err
}

func DistributePayments(escrowSeed string, escrowPubkey string, projIndex int, amount int) error {
	// this should act as the service which redistributes payments received out to the parties involved
	// amount is the amount that we want to give back to the investors and other entities involved
	project, err := RetrieveProject(projIndex)
	if err != nil {
		errors.Wrap(err, "couldn't retrieve project, quitting!")
	}

	fixedRate := 0.05 // 5 % of the totla investment as return or somethign similar. Should not be hardcoded
	// TODO: return money to the developers and other people involved
	amountGivenBack := fixedRate * float64(amount)
	for pubkey, percentage := range project.InvestorMap {
		// send x to this pubkey
		txAmount := percentage * amountGivenBack
		_, _, err := xlm.SendXLM(pubkey, utils.FtoS(txAmount), escrowSeed, "returns")
		if err != nil {
			log.Println(err) // if there is an error with one payback, doesn't mean we should stop and wait for the others
			continue
		}
	}
	return nil
	// we have the projects, we need to find the percentages donated by investors
}

// CalculatePayback calculates the amount of payback assets that must be issued in relation
// to the total amount invested in the project
// MW: Why is this after payback?
func (project Project) CalculatePayback(amount string) string {
	amountF := utils.StoF(amount)
	amountPB := (amountF / float64(project.TotalValue)) * float64(project.Years*12)
	amountPBString := utils.FtoS(amountPB)
	return amountPBString
}

// monitorPaybacks monitors whether the user is paying back regularly towards the given project
// thread has to be isolated since if this fails, we stop tracking paybacks by the recipient.
// TODO: Add first loss guarantor here
func monitorPaybacks(recpIndex int, projIndex int) {
	for {
		project, err := RetrieveProject(projIndex)
		if err != nil {
			log.Println("Couldn't retrieve project")
			continue
		}

		recipient, err := database.RetrieveRecipient(recpIndex)
		if err != nil {
			log.Println("Couldn't retrieve recipient")
			continue
		}

		// this will be our payback period and we need to check if the user pays us back

		nowTime := utils.Unix()
		timeElapsed := nowTime - project.DateLastPaid                   // this would be in seconds (unix time)
		period := int64(project.PaybackPeriod * consts.OneWeekInSecond) // in seconds due to the const
		if period == 0 {
			period = 1 // for the test suite
		}
		factor := timeElapsed / period

		// Reputation adjustments based on payback history:
		if factor <= 1 {
			// don't do anything since the user has been paying back regularly
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

		time.Sleep(time.Duration(consts.OneWeekInSecond)) // poll every week to ch eckprogress on payments
	}
}
