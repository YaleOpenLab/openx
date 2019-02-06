package solar

import (
	"fmt"
	"log"
	"time"

	assets "github.com/OpenFinancing/openfinancing/assets"
	consts "github.com/OpenFinancing/openfinancing/consts"
	database "github.com/OpenFinancing/openfinancing/database"
	issuer "github.com/OpenFinancing/openfinancing/issuer"
	notif "github.com/OpenFinancing/openfinancing/notif"
	stablecoin "github.com/OpenFinancing/openfinancing/stablecoin"
	utils "github.com/OpenFinancing/openfinancing/utils"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
)

// this file does not contain any tests associated with it right now. In the future,
// once we have a robust frontend, we can modify the CLI interface to act as a test
// for this file

func PreInvestmentCheck(projIndex int, invIndex int, recpIndex int, invAmount string) (Project, database.Investor, database.Recipient, error) {
	var project Project
	var investor database.Investor
	var recipient database.Recipient
	var err error

	project, err = RetrieveProject(projIndex)
	if err != nil {
		return project, investor, recipient, err
	}

	investor, err = database.RetrieveInvestor(invIndex)
	if err != nil {
		return project, investor, recipient, err
	}

	recipient, err = database.RetrieveRecipient(recpIndex)
	if err != nil {
		return project, investor, recipient, err
	}

	if !investor.CanInvest(invAmount) {
		log.Println("Investor has less balance than what is required to ivnest in this asset")
		return project, investor, recipient, fmt.Errorf("Investor has less balance than what is required to ivnest in this asset")
	}

	// check if investment amount is greater than or equal to the project requirements
	if utils.StoF(invAmount) > project.Params.TotalValue-project.Params.MoneyRaised {
		return project, investor, recipient, err
	}

	return project, investor, recipient, nil
}

func SendUSDToPlatform(invSeed string, invAmount string, projIndex int) (string, error) {
	// send stableusd to the platform (not the issuer) since the issuer will be locked
	// and we can't use the funds. We also need ot be able to redeem the stablecoin for fiat
	// so we can't burn them
	platformPubkey, err := wallet.ReturnPubkey(consts.PlatformSeed)
	if err != nil {
		return "", err
	}

	invPubkey, err := wallet.ReturnPubkey(invSeed)
	if err != nil {
		return "", err
	}

	oldPlatformBalance, err := xlm.GetAssetBalance(platformPubkey, stablecoin.Code)
	if err != nil {
		return "", err
	}

	_, txhash, err := assets.SendAsset(stablecoin.Code, stablecoin.PublicKey, platformPubkey, invAmount, invSeed, invPubkey, "Opensolar investment: "+utils.ItoS(projIndex))
	if err != nil {
		log.Println("Sending stableusd to platform failed", platformPubkey, invAmount, invSeed, invPubkey)
		return txhash, err
	}

	log.Println("Sent STABLEUSD to platform, confirmation: ", txhash)
	time.Sleep(5 * time.Second) // wait for a block

	newPlatformBalance, err := xlm.GetAssetBalance(platformPubkey, stablecoin.Code)
	if err != nil {
		return txhash, err
	}

	if utils.StoF(newPlatformBalance)-utils.StoF(oldPlatformBalance) < utils.StoF(invAmount)-1 {
		return txhash, fmt.Errorf("Sent amount doesn't match with investment amount")
	}
	return txhash, nil
}

// InvestInProject invests in a particular solar project given required parameters
func InvestInProject(projIndex int, invIndex int, invAmount string, invSeed string) (Project, error) {
	var err error

	var project Project
	var investor database.Investor

	project, err = RetrieveProject(projIndex)
	if err != nil {
		return project, err
	}

	investor, err = database.RetrieveInvestor(invIndex)
	if err != nil {
		return project, err
	}

	// user has decided to invest in a part of the project (don't know if full yet)
	// no asset codes assigned yet, we need to create them
	// you can retrieve asetCodes anywhere since metadata is assumed to be unique
	if project.Params.SeedAssetCode == "" && project.Params.InvestorAssetCode == "" {
		// this project does not have an issuer associated with it yet since there has been
		// no seed round and an investment round
		project.Params.InvestorAssetCode = assets.AssetID(consts.InvestorAssetPrefix + project.Params.Metadata) // set the investor asset code
		err = issuer.InitIssuer(project.Params.Index, consts.IssuerSeedPwd)
		if err != nil {
			return project, err
		}
		err = issuer.FundIssuer(project.Params.Index, consts.IssuerSeedPwd, consts.PlatformSeed)
		if err != nil {
			return project, err
		}
	}

	// offer user to exchange xlm for stableusd and invest directly if the user does not have stableusd
	// this should be a menu on the Frontend but here we do this automatically
	err = stablecoin.OfferExchange(investor.U.PublicKey, invSeed, invAmount)
	if err != nil {
		log.Println(err)
		return project, err
	}

	stableTxHash, err := SendUSDToPlatform(invSeed, invAmount, project.Params.Index)
	if err != nil {
		log.Println(err)
		return project, err
	}

	issuerPubkey, issuerSeed, err := wallet.RetrieveSeed(issuer.CreatePath(project.Params.Index), consts.IssuerSeedPwd)
	if err != nil {
		log.Println(err)
		return project, err
	}

	InvestorAsset := assets.CreateAsset(project.Params.InvestorAssetCode, issuerPubkey)
	invTrustTxHash, err := assets.TrustAsset(InvestorAsset.Code, issuerPubkey, utils.FtoS(project.Params.TotalValue), investor.U.PublicKey, invSeed)
	if err != nil {
		return project, err
	}

	log.Println("Investor trusted asset: ", InvestorAsset.Code, " tx hash: ", invTrustTxHash)
	_, invAssetTxHash, err := assets.SendAssetFromIssuer(InvestorAsset.Code, investor.U.PublicKey, invAmount, issuerSeed, issuerPubkey)
	if err != nil {
		return project, err
	}

	log.Printf("Sent InvAsset %s to investor %s with txhash %s", InvestorAsset.Code, investor.U.PublicKey, invAssetTxHash)
	// investor asset sent, update project.Params's BalLeft
	fmt.Println("Updating investor to handle invested amounts and assets")
	project.Params.MoneyRaised += utils.StoF(invAmount)
	project.ProjectInvestors = append(project.ProjectInvestors, investor)
	investor.AmountInvested += utils.StoF(invAmount)
	investor.InvestedSolarProjects = append(investor.InvestedSolarProjects, InvestorAsset.Code)
	// keep note of who all invested in this asset (even though it should be easy
	// to get that from the blockchain)
	err = investor.Save()
	if err != nil {
		return project, err
	}

	if investor.U.Notification {
		notif.SendInvestmentNotifToInvestor(projIndex, investor.U.Email, stableTxHash, invTrustTxHash, invAssetTxHash)
	}

	/*
		// The main difference between the RPC and non RPC version is that we don't send
		// any assets to the recipient in this case. We need to have handlers which will
		// take care of follow ups to send assets to the recipient later
			err = project.sendRecipientAssets(recipient, issuerPubkey, issuerSeed, recpSeed)
			if err != nil {
				return project, err
			}
	*/
	err = project.Save()
	if err != nil {
		return project, err
	}
	if project.Params.MoneyRaised == project.Params.TotalValue {
		project.Lock = true
		err = project.Save()
		if err != nil {
			return project, err
		}
		project.sendRecipientNotification()
		go sendRecipientAssets(project.Params.Index, issuerPubkey, issuerSeed)
	}
	return project, err
}

// SeedInvestInProject is similar to InvestInProject differing only in that it distributes
// seed assets instead of investor assets
func SeedInvestInProject(projIndex int, invIndex int, recpIndex int, invAmount string,
	invSeed string, recpSeed string) (Project, error) {

	project, investor, _, err := PreInvestmentCheck(projIndex, invIndex, recpIndex, invAmount)
	if err != nil {
		return project, err
	}

	// limit seed investing to one round only (as per traditional standards)
	// so we need not detect if the user has invested already because if he has,
	// he should not be able to invest in the project again
	if project.Params.SeedAssetCode == "" {
		// this person is the first investor, set the investor asset name and create the
		// issuer that will be created for this particular project
		project.Params.SeedAssetCode = assets.AssetID(consts.InvestorAssetPrefix + project.Params.Metadata) // set the investor asset code
		err = issuer.InitIssuer(project.Params.Index, consts.IssuerSeedPwd)
		if err != nil {
			return project, err
		}
		err = issuer.FundIssuer(project.Params.Index, consts.IssuerSeedPwd, consts.PlatformSeed)
		if err != nil {
			return project, err
		}
	}

	// we now have the seed asset and the issuer setup
	stableTxHash, err := SendUSDToPlatform(invSeed, invAmount, project.Params.Index)
	if err != nil {
		return project, err
	}

	issuerPubkey, issuerSeed, err := wallet.RetrieveSeed(issuer.CreatePath(project.Params.Index), consts.IssuerSeedPwd)
	if err != nil {
		return project, err
	}

	project.Params.SeedAssetCode = assets.AssetID(consts.SeedAssetPrefix + project.Params.Metadata)
	SeedAsset := assets.CreateAsset(project.Params.SeedAssetCode, issuerPubkey)

	invTrustTxHash, err := assets.TrustAsset(SeedAsset.Code, issuerPubkey, utils.FtoS(project.Params.TotalValue), investor.U.PublicKey, invSeed)
	if err != nil {
		return project, err
	}

	log.Println("Investor trusted asset: ", SeedAsset.Code, " tx hash: ", invTrustTxHash)
	_, invAssetTxHash, err := assets.SendAssetFromIssuer(SeedAsset.Code, investor.U.PublicKey, invAmount, issuerSeed, issuerPubkey)
	if err != nil {
		return project, err
	}

	log.Printf("Sent SeedAsset %s to investor %s with txhash %s", SeedAsset.Code, investor.U.PublicKey, invAssetTxHash)

	project.Params.MoneyRaised += utils.StoF(invAmount)
	project.SeedInvestors = append(project.SeedInvestors, investor)
	investor.AmountInvested += utils.StoF(invAmount)
	investor.InvestedSolarProjects = append(investor.InvestedSolarProjects, SeedAsset.Code)

	err = investor.Save()
	if err != nil {
		return project, err
	}

	if investor.U.Notification {
		notif.SendSeedInvestmentNotifToInvestor(projIndex, investor.U.Email, stableTxHash, invTrustTxHash, invAssetTxHash)
	}

	err = project.Save()
	if err != nil {
		return project, err
	}
	if project.Params.MoneyRaised == project.Params.TotalValue {
		project.Lock = true
		err = project.Save()
		if err != nil {
			return project, err
		}
		project.sendRecipientNotification()
	}
	return project, err
}

func (project *Project) sendRecipientNotification() {
	// this project covers up the amount nedeed for the project, so send the recipient
	// a notification that their project has been invested in and that they need
	// to logon to the platform in order to accept the investment
	// notif.SendUnlockNotifToRecipient(project.Params.Index, project.ProjectRecipient.U.Email)
	notif.SendUnlockNotifToRecipient(project.Params.Index, "varunramganesh@gmail.com")
}

func sendRecipientAssets(projIndex int, issuerPubkey string, issuerSeed string) error {
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
	recipient := project.ProjectRecipient
	// now send the assets as normal
	project.Params.DebtAssetCode = assets.AssetID(consts.DebtAssetPrefix + project.Params.Metadata)
	project.Params.PaybackAssetCode = assets.AssetID(consts.PaybackAssetPrefix + project.Params.Metadata)

	DebtAsset := assets.CreateAsset(project.Params.DebtAssetCode, issuerPubkey)
	PaybackAsset := assets.CreateAsset(project.Params.PaybackAssetCode, issuerPubkey)

	pbAmtTrust := utils.ItoS(project.Params.Years * 12 * 2) // two way exchange possible, to account for errors

	recpPbTrustHash, err := assets.TrustAsset(PaybackAsset.Code, issuerPubkey, pbAmtTrust, recipient.U.PublicKey, recpSeed)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Recipient Trusts Debt asset: ", DebtAsset.Code, " tx hash: ", recpPbTrustHash)
	_, recpAssetHash, err := assets.SendAssetFromIssuer(PaybackAsset.Code, recipient.U.PublicKey, pbAmtTrust, issuerSeed, issuerPubkey) // same amount as debt
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Sent PaybackAsset to recipient %s with txhash %s", recipient.U.PublicKey, recpAssetHash)
	recpDebtTrustHash, err := assets.TrustAsset(DebtAsset.Code, issuerPubkey, utils.FtoS(project.Params.TotalValue*2), recipient.U.PublicKey, recpSeed)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Recipient Trusts Payback asset: ", PaybackAsset.Code, " tx hash: ", recpDebtTrustHash)
	_, recpDebtAssetHash, err := assets.SendAssetFromIssuer(DebtAsset.Code, recipient.U.PublicKey, utils.FtoS(project.Params.TotalValue), issuerSeed, issuerPubkey) // same amount as debt
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Sent PaybackAsset to recipient %s with txhash %s\n", recipient.U.PublicKey, recpDebtAssetHash)
	project.Params.BalLeft = float64(project.Params.TotalValue)
	project.ProjectRecipient = recipient // need to udpate project.Params each time recipient is mutated
	project.Stage = FundedProject        // set funded project stage
	recipient.ReceivedSolarProjects = append(recipient.ReceivedSolarProjects, DebtAsset.Code)

	err = recipient.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	err = project.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	txhash, err := issuer.FreezeIssuer(project.Params.Index, "blah")
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Tx hash for freezing issuer is: %s", txhash)
	if recipient.U.Notification {
		notif.SendInvestmentNotifToRecipient(project.Params.Index, recipient.U.Email, recpPbTrustHash, recpAssetHash, recpDebtTrustHash, recpDebtAssetHash)
	}
	fmt.Printf("PROJECT %d's INVESTMENT CONFIRMED!", project.Params.Index)
	err = project.Save()
	if err != nil {
		return err
	}
	// need to run a separate goroutine for payback
	go sendPaymentNotif(recipient.U.Index, project.Params.Index, recipient.U.Email)
	go monitorPaybacks(recipient.U.Index, project.Params.Index, project.Params.DebtAssetCode)
	return nil
}

// sendPaymentNotif sends a notification every payback period to the recipient to
// kindly remind him to payback towards the project
func sendPaymentNotif(recpIndex int, projIndex int, email string) {
	// setup a payback monitoring routine for monitoring if the recipient pays us back on time
	// the recipient must give his email to receive updates
	paybackTimes := 0
	for {
		project, err := RetrieveProject(projIndex)

		if err != nil {
			log.Println(err)
			message := "Error in payback routine, please contact help as soon as you receive this message " + err.Error()
			notif.SendAlertEmail(message, email) // don't catch the error here
			time.Sleep(time.Duration(project.PaybackPeriod * consts.OneWeekInSecond))
		}

		_, err = database.RetrieveRecipient(recpIndex) // need to retireve to make sure nothing goes awry
		if err != nil {
			log.Println(err)
			message := "Error while retrieving your account details, please contact help as soon as you receive this message " + err.Error()
			notif.SendAlertEmail(message, email) // don't catch the error here
			time.Sleep(time.Duration(project.PaybackPeriod * consts.OneWeekInSecond))
		}

		if paybackTimes == 0 {
			// sleep and bother during the next cycle
			time.Sleep(time.Duration(project.PaybackPeriod * consts.OneWeekInSecond))
		}

		// PAYBACK TIME!!
		// we don't know if the user has paid, but we send an email anyway
		notif.SendPaybackAlertEmail(project.Params.Index, email)
		// sleep until the next payment is due
		paybackTimes += 1
		log.Println("Sent: ", email, "a notification on payments for payment cycle: ", paybackTimes)
		time.Sleep(time.Duration(project.PaybackPeriod * consts.OneWeekInSecond))
	}
}
