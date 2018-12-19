package database

import (
	"log"

	oracle "github.com/YaleOpenLab/smartPropertyMVP/stellar/oracle"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
	"github.com/stellar/go/build"
)

// handle all methods on entities here isntead of  having them over at their specific
// files in order to avoid searching where the methods are
// TrustAsset creates a trustline from the caller towards the specific asset
// and asset issuer with a _limit_ set on the maximum amount of tokens that can be sent
// through the trust channel. Each trustline costs 0.5XLM.
func (a *Investor) TrustAsset(asset build.Asset, limit string) (string, error) {
	// TRUST is FROM recipient TO issuer
	trustTx, err := build.Transaction(
		build.SourceAccount{a.PublicKey},
		build.AutoSequence{SequenceProvider: utils.DefaultTestNetClient},
		build.TestNetwork,
		build.Trust(asset.Code, asset.Issuer, build.Limit(limit)),
	)

	if err != nil {
		return "", err
	}

	trustTxe, err := trustTx.Sign(a.Seed)
	if err != nil {
		return "", err
	}

	trustTxeB64, err := trustTxe.Base64()
	if err != nil {
		return "", err
	}

	tx, err := utils.DefaultTestNetClient.SubmitTransaction(trustTxeB64)
	if err != nil {
		return "", err
	}

	log.Println("Trusted asset tx: ", tx.Hash)
	return tx.Hash, nil
}

// SendAssetToIssuer sends back assets fromn an asset holder to the issuer of the asset.
// This method is called by the receiver of assets.
func (a *Recipient) SendAssetToIssuer(assetName string, issuerPubkey string, amount string) (int32, string, error) {
	// SendAssetToIssuer is FROM recipient / investor to issuer
	paymentTx, err := build.Transaction(
		build.SourceAccount{a.PublicKey},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: utils.DefaultTestNetClient},
		build.Payment(
			build.Destination{AddressOrSeed: issuerPubkey},
			build.CreditAmount{assetName, issuerPubkey, amount},
		),
	)

	if err != nil {
		return -11, "", err
	}

	paymentTxe, err := paymentTx.Sign(a.Seed)
	if err != nil {
		return -11, "", err
	}

	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {
		return -11, "", err
	}

	tx, err := utils.DefaultTestNetClient.SubmitTransaction(paymentTxeB64)
	if err != nil {
		return -11, "", err
	}

	return tx.Ledger, tx.Hash, nil
}

// Payback is called when the receiver of the DEBToken wants to pay a fixed amount
// of money back to the issuer of the DEBTokens. One way to imagine this would be
// like an electricity bill, something that people pay monthly but only that in this
// case, the electricity is free, so they pay directly towards the solar panels.
// The process of Payback roughly involves the followign steps:
// 1. Pay the issuer in DEBTokens with whatever amount desired.
// The oracle price of
// electricity cost is a lower bound (since the government would not like it if people
// default on their payments). Anything below the lower bound gets a warning in
// order for people to pay more, we could also have a threshold mechanism that says
// if a person constantly defaults for more than half the owed amount for three
// consecutive months, we sell power directly to the grid. THis could also be used
// for a rating system, where the frontend UI can have a rating based on whether
// the recipient has defaulted or not in the past.
// 2. The receiver checks whether the amount is greater than Oracle Threshold and
// if so, sends back PBTokens, which stand for the month equivalent of payments.
// eg. the school has opted for a 5 year payback period, the school owes the issuer
// 60 PBTokens and the issuer sends back 1PBToken every month if the school pays
// invested_amount/60 DEBTokens back to the issuer
// 3. The recipient checks whether the PBTokens received correlate to the amount
// that it sent and if not, raises the dispute since the forward DEBToken payment
// is on chain and resolves the dispute itself using existing off chain legal frameworks
// (issued bonds, agreements, etc)
func (a *Recipient) Payback(index uint32, assetName string, issuerPubkey string, amount string) error {

	oldBalance, err := xlm.GetAssetBalance(assetName, a.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	PBAmount, err := oracle.PriceOracle()
	if err != nil {
		log.Println("Unable to fetch oracle price, exiting")
		return err
	}

	log.Println("Retrieved average price from oracle: ", PBAmount)
	// the oracke needs to know the assetName so that it can find the other details
	// about this asset from the db. This should run on the server side and must
	// be split when we do run client side stuff.
	// hardcode for now, need to add the oracle here so that we
	// can do this dynamically
	confHeight, txHash, err := a.SendAssetToIssuer(assetName, issuerPubkey, amount)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Paid debt amount: ", amount, " back to issuer, tx hash: ", txHash, " ", confHeight)
	log.Println("Checking balance to see if our account was debited")
	newBalance, err := xlm.GetAssetBalance(assetName, a.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	newBalanceFloat := utils.FtoS(newBalance)
	oldBalanceFloat := utils.FtoS(oldBalance)
	amountFloat := utils.FtoS(PBAmount)

	paidAmount := oldBalanceFloat - newBalanceFloat
	log.Println("Old Balance: ", oldBalanceFloat, "New Balance: ", newBalanceFloat, "Paid: ", paidAmount, "Req Amount: ", amountFloat)

	// would be nice to take some additional action like sending a notification or
	// something to investors or to the email address given so that everyone is made
	// aware of this and there's data transparency

	if paidAmount < amountFloat {
		log.Println("Amount paid is less than amount required, balance not updating, please amke sure to cover this next time")
	} else if paidAmount > amountFloat {
		log.Println("You've chosen to pay more than what is required for this month. Adjusting payback period accordingly")
	} else {
		log.Println("You've paid exactly what is required for this month. Payback period remains as usual")
	}
	// we need to update the database here
	givenOrder, err := RetrieveOrder(index)
	if err != nil {
		log.Println("Given order not found in the database")
		return err
	}
	givenOrder.BalLeft = float64(givenOrder.TotalValue) - paidAmount
	givenOrder.DateLastPaid = utils.Timestamp()
	// balLeft must be updated on the server side and can be challenged easily
	// if there's some discrepancy since the tx's are on the blockchain
	return InsertOrder(givenOrder)
}
