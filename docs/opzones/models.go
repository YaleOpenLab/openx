/*
package ozones

/*
import (
  "log"
	"time"

	xlm "github.com/YaleOpenLab/openx/chains/xlm"
  assets "github.com/YaleOpenLab/openx/chains/xlm/assets"
  issuer "github.com/YaleOpenLab/openx/chains/xlm/issuer"
  wallet "github.com/YaleOpenLab/openx/chains/xlm/wallet"
  utils "github.com/Varunram/essentials/utils"
  consts "github.com/YaleOpenLab/openx/consts"
  database "github.com/YaleOpenLab/openx/database"
  notif "github.com/YaleOpenLab/openx/notif"
  "github.com/pkg/errors"
)
// Invest invests in a particular project
func Invest(projIndex int, invIndex int, invAssetCode string, invSeed string,
	invAmount float64, trustLimit float64, investorIndices []int, application string) error {

	issuerPath := consts.OpzonesIssuerDir
	issuerPubkey, issuerSeed, err := wallet.RetrieveSeed(issuer.GetPath(issuerPath, projIndex), consts.IssuerSeedPwd)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve issuer seed")
	}

	if len(investorIndices) == 0 {
		_ = assets.CreateAsset(invAssetCode, issuerPubkey) // create the asset, since it would not have been created earlier
	}

	investor, err := database.RetrieveInvestor(invIndex)
	if err != nil {
		return err
	}

	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	stableTxHash, err := SendUSDToPlatform(invSeed, invAmount, "Opzones investment: "+projIndexString)
	if err != nil {
		log.Println("Unable to send STABELUSD to platform", err)
		return err
	}

	log.Println("STABLEHASH: ", stableTxHash)

	// make investor trust the asset that we provide
	investorAssetTrustHash, err := assets.TrustAsset(invAssetCode, issuerPubkey, trustLimit, invSeed)
	// trust upto the total value of the asset
	if err != nil {
		return err
	}
	log.Println("Investor trusted asset: ", invAssetCode, " tx hash: ", investorAssetTrustHash)
	log.Println("Sending InvestorAsset: ", invAssetCode, "for: ", invAmount)
	_, investorAssetHash, err := assets.SendAssetFromIssuer(invAssetCode, investor.U.StellarWallet.PublicKey, invAmount, issuerSeed, issuerPubkey)
	if err != nil {
		return err
	}
	log.Printf("Sent InvestorAsset %s to investor %s with txhash %s", invAssetCode, investor.U.StellarWallet.PublicKey, investorAssetHash)
	// investor asset sent, update a.Params's BalLeft

	investor.AmountInvested += invAmount
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

	issuerPubkey, issuerSeed, err := wallet.RetrieveSeed(issuer.GetPath(issuerPath, projIndex), consts.IssuerSeedPwd)
	if err != nil {
		log.Println("Unable to retrieve issuer seed", err)
		return err
	}

	// this investment asset is a one way asset ie the recipient need not pay this back to the platform since these funds
	// would be used in the unit's construction

	debtTrustHash, err := assets.TrustAsset(debtAssetCode, issuerPubkey, totalValue, recpSeed)
	if err != nil {
		log.Println("Error while trusting investment asset", err)
		return err
	}
	log.Printf("Recipient Trusts Investment asset %s with txhash %s", debtAssetCode, debtTrustHash)

	_, recpDebtAssetHash, err := assets.SendAssetFromIssuer(debtAssetCode, recipient.U.StellarWallet.PublicKey, totalValue, issuerSeed, issuerPubkey) // same amount as investment
	if err != nil {
		log.Println("Error while sending investment asset", err)
		return err
	}

	log.Printf("Sent DebtAsset to recipient %s with txhash %s\n", recipient.U.StellarWallet.PublicKey, recpDebtAssetHash)
	recipient.ReceivedConstructionBonds = append(recipient.ReceivedConstructionBonds, debtAssetCode)
	err = recipient.Save()
	if err != nil {
		return errors.Wrap(err, "can't save recipient")
	}

	txhash, err := issuer.FreezeIssuer(issuerPath, projIndex, "blah")
	if err != nil {
		log.Println("Error while freezing issuer", err)
		return err
	}

	log.Printf("Tx hash for freezing issuer is: %s", txhash)
	log.Printf("PROJECT %d's INVESTMENT CONFIRMED!", projIndex)

	if recipient.U.Notification {
		go notif.SendInvestmentNotifToRecipientOZ(projIndex, recipient.U.Email, debtTrustHash, recpDebtAssetHash)
	}

	return nil
}


// SendUSDToPlatform sends STABLEUSD back to the platform for investment
func SendUSDToPlatform(invSeed string, invAmount float64, memo string) (string, error) {
	// send stableusd to the platform (not the issuer) since the issuer will be locked
	// and we can't use the funds. We also need ot be able to redeem the stablecoin for fiat
	// so we can't burn them
	var oldPlatformBalance float64
	var err error
	oldPlatformBalance, err = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.StablecoinCode)
	if err != nil {
		// platform does not have stablecoin, shouldn't arrive here ideally
		oldPlatformBalance = 0
	}

	var txhash string
	if !consts.Mainnet {
		_, txhash, err = assets.SendAsset(consts.StablecoinCode, consts.StablecoinPublicKey, consts.PlatformPublicKey, invAmount, invSeed, memo)
		if err != nil {
			return txhash, errors.Wrap(err, "sending stableusd to platform failed")
		}
	} else {
		_, txhash, err = assets.SendAsset(consts.AnchorUSDCode, consts.AnchorUSDAddress, consts.PlatformPublicKey, invAmount, invSeed, memo)
		if err != nil {
			return txhash, errors.Wrap(err, "sending stableusd to platform failed")
		}
	}

	log.Println("Sent STABLEUSD to platform, confirmation: ", txhash)
	time.Sleep(5 * time.Second) // wait for a block

	var newPlatformBalance float64
	if !consts.Mainnet {
		newPlatformBalance, err = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.StablecoinCode)
		if err != nil {
			return txhash, errors.Wrap(err, "error while getting asset balance")
		}
	} else {
		newPlatformBalance, err = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.AnchorUSDCode)
		if err != nil {
			return txhash, errors.Wrap(err, "error while getting asset balance")
		}
	}

	if newPlatformBalance-oldPlatformBalance < invAmount-1 {
		return txhash, errors.New("Sent amount doesn't match with investment amount")
	}
	return txhash, nil
}
*/
