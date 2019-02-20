package investmentcrowdfunding

import (
	"fmt"
	"log"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	issuer "github.com/YaleOpenLab/openx/issuer"
	models "github.com/YaleOpenLab/openx/models"
	notif "github.com/YaleOpenLab/openx/notif"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
)

// Debt Crowdfunding is a model where the investor loans out some initial capital and receives interest on that investment

// Invest invests in a particular project
func Invest(projIndex int, invIndex int, invAssetCode string, invSeed string,
	invAmount string, trustLimit string, investorIndices []int, application string) error {

	issuerPath := consts.OpzonesIssuerDir
	issuerPubkey, issuerSeed, err := wallet.RetrieveSeed(issuer.CreatePath(issuerPath, projIndex), consts.IssuerSeedPwd)
	if err != nil {
		log.Println("Unable to retrieve issuer seed", err)
		log.Println(err)
		return err
	}

	if len(investorIndices) == 0 {
		_ = assets.CreateAsset(invAssetCode, issuerPubkey) // create the asset, since it would not have been created earlier
	}

	investor, err := database.RetrieveInvestor(invIndex)
	if err != nil {
		return err
	}

	stableTxHash, err := SendUSDToPlatform(invSeed, invAmount, "Opzones investment: "+utils.ItoS(projIndex))
	if err != nil {
		log.Println("Unable to send STABELUSD to platform", err)
		return err
	}

	log.Println("STABLEHASH: ", stableTxHash)

	// make investor trust the asset that we provide
	investorAssetTrustHash, err := assets.TrustAsset(invAssetCode, issuerPubkey, trustLimit, investor.U.PublicKey, invSeed)
	// trust upto the total value of the asset
	if err != nil {
		return err
	}
	log.Println("Investor trusted asset: ", invAssetCode, " tx hash: ", investorAssetTrustHash)
	log.Println("Sending InvestorAsset: ", invAssetCode, "for: ", invAmount)
	_, investorAssetHash, err := assets.SendAssetFromIssuer(invAssetCode, investor.U.PublicKey, invAmount, issuerSeed, issuerPubkey)
	if err != nil {
		return err
	}
	log.Printf("Sent InvestorAsset %s to investor %s with txhash %s", invAssetCode, investor.U.PublicKey, investorAssetHash)
	// investor asset sent, update a.Params's BalLeft
	investor.AmountInvested += utils.StoF(invAmount)
	if application == "constructionBond" {
		investor.InvestedBonds = append(investor.InvestedBonds, invAssetCode)
	} else if application == "livingunitcoop" {
		investor.InvestedCoops = append(investor.InvestedCoops, invAssetCode)
	}
	err = investor.Save() // save investor creds now that we're done
	if err != nil {
		return err
	}

	if investor.U.Notification {
		go notif.SendInvestmentNotifToInvestor(projIndex, investor.U.Email, stableTxHash, investorAssetTrustHash, investorAssetHash)
	}
	// need to check whether we need to send assets to investor here
	return nil
}

// ReceiveBond sends out assets to the recipient
func ReceiveBond(issuerPath string, recpIndex int, projIndex int, debtAssetCode string,
	recpSeed string, totalValue float64) error {

	recipient, err := database.RetrieveRecipient(recpIndex)
	if err != nil {
		log.Println("Unable to retrieve recipient from database", err)
		return err
	}

	issuerPubkey, issuerSeed, err := wallet.RetrieveSeed(issuer.CreatePath(issuerPath, projIndex), consts.IssuerSeedPwd)
	if err != nil {
		log.Println("Unable to retrieve issuer seed", err)
		return err
	}

	// this investment asset is a one way asset ie the recipient need not pay this back to the platform since these funds
	// would be used in the unit's construction
	debtTrustHash, err := assets.TrustAsset(debtAssetCode, issuerPubkey, utils.FtoS(totalValue), recipient.U.PublicKey, recpSeed)
	if err != nil {
		log.Println("Error while trusting investment asset", err)
		return err
	}
	log.Printf("Recipient Trusts Investment asset %s with txhash %s", debtAssetCode, debtTrustHash)

	_, recpDebtAssetHash, err := assets.SendAssetFromIssuer(debtAssetCode, recipient.U.PublicKey, utils.FtoS(totalValue), issuerSeed, issuerPubkey) // same amount as investment
	if err != nil {
		log.Println("Error while sending investment asset", err)
		return err
	}

	log.Printf("Sent DebtAsset to recipient %s with txhash %s\n", recipient.U.PublicKey, recpDebtAssetHash)
	recipient.ReceivedConstructionBonds = append(recipient.ReceivedConstructionBonds, debtAssetCode)
	err = recipient.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	txhash, err := issuer.FreezeIssuer(issuerPath, projIndex, "blah")
	if err != nil {
		log.Println("Error while freezing issuer", err)
		return err
	}

	log.Printf("Tx hash for freezing issuer is: %s", txhash)
	fmt.Printf("PROJECT %d's INVESTMENT CONFIRMED!", projIndex)

	if recipient.U.Notification {
		go notif.SendInvestmentNotifToRecipientOZ(projIndex, recipient.U.Email, debtTrustHash, recpDebtAssetHash)
	}

	return nil
}

// SendUSDToPlatform can be used to send USD back to the platform
func SendUSDToPlatform(invSeed string, invAmount string, memo string) (string, error) {
	return models.SendUSDToPlatform(invSeed, invAmount, memo)
}
