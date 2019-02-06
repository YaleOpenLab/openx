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
	oracle "github.com/OpenFinancing/openfinancing/oracle"
	stablecoin "github.com/OpenFinancing/openfinancing/stablecoin"
	utils "github.com/OpenFinancing/openfinancing/utils"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
)

// Payback is called when the receiver of the DebtAsset wants to pay a fixed amount
// of money back to the issuer of the DebtAssets. One way to imagine this would be
// like an electricity bill, something that people pay monthly but only that in this
// case, the electricity is free, so they pay directly towards the solar panels.
// The process of Payback roughly involves the followign steps:
// 1. Pay the issuer in DebtAssets with whatever amount desired.
// The oracle price of
// electricity cost is a lower bound (since the government would not like it if people
// default on their payments). (MW: Explain this lower bound and default issue more)
// Anything below the lower bound gets a warning in
// project for people to pay more, we could also have a threshold mechanism that says
// if a person constantly defaults for more than half the owed amount for three
// consecutive months, we sell power directly to the grid. This could also be used
// for a rating system, where the frontend UI can have a rating based on whether
// the recipient has defaulted or not in the past.
// 2. The receiver checks whether the amount is greater than Oracle Threshold and
// if so, sends back PaybackAssets, which stand for the month equivalent of payments.
// eg. the school has opted for a 5 year payback period, the school owes the issuer
// 60 PaybackAssets and the issuer sends back 1PaybackAsset every month if the school pays
// invested_amount/60 DebtAssets back to the issuer
// 3. The recipient checks whether the PaybackAssets received correlate to the amount
// that it sent and if not, raises the dispute since the forward DebtAsset payment
// is on chain and resolves the dispute itself using existing off chain legal frameworks
// (issued bonds, agreements, etc)
// TODO: evaluate whether we need PaybackAsset
func Payback(recpIndex int, projIndex int, assetName string, amount string, recipientSeed string,
	platformPubkey string) error {
	// in this flow, we exchange xlm for stableUSD and run the process in stableUSD but
	// in reality, we run using stableUSD directly wihtout th XLM part. Anonymous investors
	// or people who are a bit more  careful with keys could still use the XLM -> StableUSD
	// bridge, but the need for that and kyc regulations have to be evaluated.
	// Also, if a user can hold balance in btc / xlm, we could direct them to exchange it
	// using the DEX taht stellar provides and use that asset (or we could setup a payment
	// provider which accepts fiat + crypto and issue this asset ourselves)
	issuerPubkey, _, err := wallet.RetrieveSeed(issuer.CreatePath(projIndex), consts.IssuerSeedPwd)
	if err != nil {
		return err
	}

	recipient, err := database.RetrieveRecipient(recpIndex)
	if err != nil {
		return err
	}

	project, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}
	// once we have the stablecoin here, we can remove the assetName
	StableBalance, err := xlm.GetAssetBalance(recipient.U.PublicKey, "STABLEUSD")
	// checks for the stablecoin asset
	if err != nil || (utils.StoF(StableBalance) < utils.StoF(amount)) {
		log.Println("You do not have the required stablecoin balance, please refill")
		return err
	}
	// pay stableUSD back to platform
	_, stableUSDHash, err := assets.SendAsset(stablecoin.Code, consts.StableCoinAddress, platformPubkey, amount, recipientSeed, recipient.U.PublicKey, "Opensolar payback: "+utils.ItoS(projIndex))
	if err != nil {
		log.Println("SEND ASSET ERR:", err, platformPubkey, amount, recipientSeed, recipient.U.PublicKey)
		return err
	}
	log.Println("Paid back platform in  stableUSD, txhash: ", stableUSDHash)

	DEBAssetBalance, err := xlm.GetAssetBalance(recipient.U.PublicKey, assetName)
	if err != nil {
		fmt.Println("Don't have the debt asset in possession", err)
		return err
	}

	monthlyBill := oracle.MonthlyBill()
	if err != nil {
		log.Println("Unable to fetch oracle price, exiting")
		return err
	}

	log.Println("Retrieved average price from oracle: ", monthlyBill)
	confHeight, debtPaybackHash, err := assets.SendAssetToIssuer(assetName, issuerPubkey, amount, recipientSeed, recipient.U.PublicKey)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Paid debt amount: ", amount, " back to issuer, tx hash: ", debtPaybackHash, " ", confHeight)
	// log.Println("Checking balance to see if our account was debited")
	newBalance, err := xlm.GetAssetBalance(recipient.U.PublicKey, assetName)
	if err != nil {
		return err
	}

	newBalanceFloat := utils.StoF(newBalance)
	DEBAssetBalanceFloat := utils.StoF(DEBAssetBalance)
	mBillFloat := utils.StoF(monthlyBill)

	paidAmount := DEBAssetBalanceFloat - newBalanceFloat
	log.Println("Old Balance: ", DEBAssetBalanceFloat, " New Balance: ", newBalanceFloat, " Paid: ", paidAmount, " Bill Amount: ", mBillFloat)

	// would be nice to take some additional action like sending a notification or
	// something to investors or to the email address given so that everyone is made
	// aware of this and there's data transparency once the recipient pays back to the
	// platform - could be a service which people can subscriibe to in case they need
	// something

	if paidAmount < mBillFloat {
		log.Println("Amount paid is less than amount required, please make sure to cover this next time")
	} else if paidAmount > mBillFloat {
		log.Println("You've chosen to pay more than what is required for this month. Adjusting payback period accordingly")
	} else {
		log.Println("You've paid exactly what is required for this month. Payback period remains as usual")
	}
	// we need to update the database here
	// no need to retrieve this project again because we have it already
	project.Params.BalLeft -= paidAmount
	project.Params.DateLastPaid = utils.Unix()
	if project.Params.BalLeft == 0 {
		log.Println("YOU HAVE PAID OFF THIS ASSET, TRANSFERRING OWNERSHIP OF ASSET TO YOU")
		// don't delete the asset from the received assets list, we still need it so
		// that we can look back and find out hwo many assets this particular
		// enttiy has been invested in, have a leaderboard kind of thing, etc.
		project.Stage = 7
		// we should call neighbourly or some ohter partner here to transfer assets
		// using the bond they provide us with
		// the nice part here is that the recipient can not pay off more than what is
		// invested because the trustline will not allow such an incident to happen
	}
	// balLeft must be updated on the server side and can be challenged easily
	// if there's some discrepancy since the tx's are on the blockchain
	err = project.updateRecipient(recipient)
	if err != nil {
		return err
	}

	err = project.Save()
	if err != nil {
		return err
	}
	if recipient.U.Notification {
		notif.SendPaybackNotifToRecipient(projIndex, recipient.U.Email, stableUSDHash, debtPaybackHash)
	}
	for _, elem := range project.ProjectInvestors {
		if elem.U.Notification {
			notif.SendPaybackNotifToInvestor(projIndex, elem.U.Email, stableUSDHash, debtPaybackHash)
		}
	}
	return err
}

// CalculatePayback is a function that simply sums the PaybackAsset
// balance and returns them to the frontend UI for a nice display
func (project Project) CalculatePayback(amount string) string {
	// the idea is that we should be able to pass an assetId to this function
	// and it must calculate how much time we have left for payback. For this example
	// until we do the db stuff, lets pass a few params (although this could be done
	// separately as well).
	amountF := utils.StoF(amount)
	amountPB := (amountF / float64(project.Params.TotalValue)) * float64(project.Params.Years*12)
	amountPBString := utils.FtoS(amountPB)
	return amountPBString
}

func monitorPaybacks(recpIndex int, projIndex int, debtAssetCode string) {
	// monitor whether the user is paying back regularly towards the given project
	// this is a routine similar to the general notification routine but focused more on the
	// payback and tracking that. Also, if the other thread fails, nothing major happens except
	// notifications, but if this one fails, we can't track paybacks, so this one has to be
	// isolated. various thresholds that we will be using for notification services are defined as constants
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
		timeElapsed := nowTime - project.Params.DateLastPaid            // this would be in seconds (unix time)
		period := int64(project.PaybackPeriod * consts.OneWeekInSecond) // in seconds due to the const
		factor := timeElapsed / period

		if factor <= 1 {
			// don't do anything since the suer has been apying back regularly
			log.Println("User: ", recipient.U.Email, "is on track paying towards order: ", projIndex)
			// maybe even update reputation here on a fractional basis depending on a user's timely payments
		} else if factor > consts.NormalThreshold && factor < consts.AlertThreshold {
			// person has not paid back for one-two consecutive period, send gentle reminder
			notif.SendNicePaybackAlertEmail(projIndex, recipient.U.Email)
		} else if factor >= consts.SternAlertThreshold && factor < consts.DisconnectionThreshold {
			// person has not paid back for four consecutive cycles, send reminder
			notif.SendSternPaybackAlertEmail(projIndex, recipient.U.Email)
			for _, elem := range project.ProjectInvestors {
				// send an email to recipients to assure them that we're on the issue and will be acting
				// soon if the recipient fails to pay again.
				notif.SendSternPaybackAlertEmailI(projIndex, elem.U.Email)
			}
			notif.SendSternPaybackAlertEmailG(projIndex, project.Guarantor.U.Email)
			// send an email out to the guarantor
		} else if factor >= consts.DisconnectionThreshold {
			// send a disconnection notice to the recipient and let them know we have redirected
			// power towards the grid. Also maybe email ourselves in this case so that we can
			// contact them personally to resolve the issue as soon as possible.
			notif.SendDisconnectionEmail(projIndex, recipient.U.Email)
			for _, elem := range project.ProjectInvestors {
				// send an email out to each investor to let them know that the recipient
				// has defaulted on payments and that we have acted on the issue.
				notif.SendDisconnectionEmailI(projIndex, elem.U.Email)
			}
			notif.SendDisconnectionEmailG(projIndex, project.Guarantor.U.Email)
			// send an email out to the guarantor
		}

		time.Sleep(time.Duration(consts.OneWeekInSecond)) // poll every week to ch eck progress on payments
	}
}
