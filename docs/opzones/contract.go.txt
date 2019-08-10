/*
package ozones

import (
	"github.com/pkg/errors"
	"log"
	"math"
	"time"

	assets "github.com/YaleOpenLab/openx/chains/xlm/assets"
	issuer "github.com/YaleOpenLab/openx/chains/xlm/issuer"
	wallet "github.com/YaleOpenLab/openx/chains/xlm/wallet"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	notif "github.com/YaleOpenLab/openx/notif"
)

// preInvestmentConstructionBonds defines the pre investment conditions pertaining to construction bonds
func preInvestmentConstructionBonds(projIndex int, invIndex int, invAmount float64) (ConstructionBond, error) {

	project, err := RetrieveConstructionBond(projIndex)
	if err != nil {
		return project, errors.Wrap(err, "couldn't retrieve construction bond from db")
	}

	investor, err := database.RetrieveInvestor(invIndex)
	if err != nil {
		return project, errors.Wrap(err, "couldn't retrieve investor from db")
	}

	rem := invAmount / project.CostOfUnit
	if math.Floor(rem) == 0 {
		return project, errors.New("You are trying to invest more than a unit's cost, do you want to invest in two units?")
	}

	assetName := assets.AssetID(project.MaturationDate + project.SecurityType + project.Rating + project.BondIssuer) // get a unique assetID

	if len(project.InvestorIndices) == 0 {
		// initialize issuer
		err = issuer.InitIssuer(consts.OpzonesIssuerDir, project.Index, consts.IssuerSeedPwd)
		if err != nil {
			return project, errors.Wrap(err, "error while initializing issuer")
		}
		err = issuer.FundIssuer(consts.OpzonesIssuerDir, project.Index, consts.IssuerSeedPwd, consts.PlatformSeed)
		if err != nil {
			return project, errors.Wrap(err, "error while funding issuer")
		}

		project.InvestorAssetCode = assets.AssetID(consts.BondAssetPrefix + assetName) // set the investor code
	}

	if !investor.CanInvest(invAmount) {
		return project, errors.Wrap(err, "Investor has less balance than what is required to ivnest in this asset")
	}

	return project, nil
}

// preInvestmentConstructionBonds defines the pre investment conditions pertaining to living unit coops
func preInvestmentLivingCoop(projIndex int, invIndex int, invAmount float64) (LivingUnitCoop, error) {

	project, err := RetrieveLivingUnitCoop(projIndex)
	if err != nil {
		return project, errors.Wrap(err, "couldn't retrieve living unit coop from db")
	}

	investor, err := database.RetrieveInvestor(invIndex)
	if err != nil {
		return project, errors.Wrap(err, "couldn't retrieve investor from db")
	}
	// check if investment amount is greater than the cost of a unit
	if invAmount != project.MonthlyPayment {
		return project, errors.New("You are trying to invest more than a unit's cost, do you want to invest in two units?")
	}

	assetName := assets.AssetID(project.Description)

	if len(project.ResidentIndices) == 0 {
		// initialize issuer
		err = issuer.InitIssuer(consts.OpzonesIssuerDir, project.Index, consts.IssuerSeedPwd)
		if err != nil {
			return project, errors.Wrap(err, "error while initializing issuer")
		}
		err = issuer.FundIssuer(consts.OpzonesIssuerDir, project.Index, consts.IssuerSeedPwd, consts.PlatformSeed)
		if err != nil {
			return project, errors.Wrap(err, "error while funding issuer")
		}

		project.InvestorAssetCode = assets.AssetID(consts.BondAssetPrefix + assetName) // set the investor code
	}

	if !investor.CanInvest(invAmount) {
		log.Println("Investor has less balance than what is required to ivnest in this asset")
		return project, err
	}

	return project, nil
}

// InvestInLivingUnitCoop invests in a particular living coop
func InvestInLivingUnitCoop(projIndex int, invIndex int, invAmount float64, invSeed string,
	recpSeed string) error {
	// we want to invest in this specific bond
	var err error

	project, err := preInvestmentLivingCoop(projIndex, invIndex, invAmount)
	if err != nil {
		return errors.Wrap(err, "could not check pre investment conditions in living unit coop")
	}

	err = Invest(projIndex, invIndex, project.InvestorAssetCode, invSeed,
		invAmount, project.Amount, project.ResidentIndices, "livingunitcoop")
	if err != nil {
		return errors.Wrap(err, "could not invest in living unit coop")
	}

	err = project.updateLivingUnitCoopAfterInvestment(invAmount, invIndex)
	if err != nil {
		return errors.Wrap(err, "Failed to update project after investment")
	}

	return nil
}

// InvestInConstructionBond invests in a particular construction bonm
func InvestInConstructionBond(projIndex int, invIndex int, invAmount float64, invSeed string) error {
	// we want to invest in this specific bond
	var err error

	project, err := preInvestmentConstructionBonds(projIndex, invIndex, invAmount)
	if err != nil {
		return errors.Wrap(err, "could not check pre investment conditions in construction bond")
	}

	trustLimit := project.CostOfUnit * float64(project.NoOfUnits)

	err = Invest(projIndex, invIndex, project.InvestorAssetCode, invSeed,
		invAmount, trustLimit, project.InvestorIndices, "constructionbond")
	if err != nil {
		return errors.Wrap(err, "could not invest in construction bond")
	}

	err = project.updateConstructionBondAfterInvestment(invAmount, invIndex)
	if err != nil {
		return errors.Wrap(err, "failed to update project after investment")
	}

	totalValue := float64(project.CostOfUnit * float64(project.NoOfUnits))
	project.DebtAssetCode = assets.AssetID(consts.DebtAssetPrefix + project.Description)

	if totalValue == project.AmountRaised {
		// send the recipient a notification to unlock the specific project and accept the investment
		err = project.sendRecipientNotification()
		if err != nil {
			return errors.Wrap(err, "error while sending notifications to recipient")
		}
		go sendRecipientAssets(projIndex, totalValue)
	}
	return nil
}

// sendRecipientAssets sends the recipient assets pertaining to an order
func sendRecipientAssets(projIndex int, totalValue float64) error {
	// send the recipient relevant debt asset
	startTime := utils.Unix()
	project, err := RetrieveConstructionBond(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve construction bond from db")
	}

	for utils.Unix()-startTime < consts.LockInterval {
		log.Printf("WAITING FOR PROJECT %d TO BE UNLOCKED", projIndex)
		project, err = RetrieveConstructionBond(projIndex)
		if err != nil {
			return errors.Wrap(err, "couldn't retrieve construction bond from db")
		}
		if !project.Lock {
			log.Println("Project UNLOCKED IN LOOP")
			break
		}
		time.Sleep(10 * time.Second)
	}

	project, err = RetrieveConstructionBond(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve construction bond from db")
	}

	recipient, err := database.RetrieveRecipient(project.RecipientIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve recipient from db")
	}

	recpSeed, err := wallet.DecryptSeed(recipient.U.StellarWallet.EncryptedSeed, project.LockPwd)
	if err != nil {
		return errors.Wrap(err, "couldn't decrypt seed")
	}

	err = ReceiveBond(consts.OpzonesIssuerDir, project.RecipientIndex, projIndex, project.DebtAssetCode, recpSeed, totalValue)
	if err != nil {
		return errors.Wrap(err, "failed to send assets to recipient project after investment")
	}

	project.LockPwd = ""
	return project.Save()
}

// sendRecipientNotification sends the notification to the recipient requesting them
// to logon to the platform and unlock the project that has just been invested in
func (project *ConstructionBond) sendRecipientNotification() error {
	recipient, err := database.RetrieveRecipient(project.RecipientIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve recipient from db")
	}
	notif.SendUnlockNotifToRecipient(project.Index, recipient.U.Email)
	project.Lock = true
	return project.Save()
}

// UnlockProject unlocks a specific project that has just been invested in
func UnlockProject(username string, pwhash string, projIndex int, seedpwd string, application string) error {
	log.Println("UNLOCKING PROJECT")
	recipient, err := database.ValidateRecipient(username, pwhash)
	if err != nil {
		return errors.Wrap(err, "couldn't validate recipient")
	}

	recpSeed, err := wallet.DecryptSeed(recipient.U.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "Error while decrpyting seed")
	}

	checkPubkey, err := wallet.ReturnPubkey(recpSeed)
	if err != nil {
		return errors.Wrap(err, "Couldn't get pubkey from seed")
	}

	if checkPubkey != recipient.U.StellarWallet.PublicKey {
		return errors.New("Failed to unlock project, public keys don't match")
	}

	if application == "constructionbond" {
		project, err := RetrieveConstructionBond(projIndex)
		if err != nil || !project.Lock {
			return errors.Wrap(err, "lock not set on project")
		}

		if recipient.U.Index != project.RecipientIndex {
			return errors.New("Recipient Indices don't match, quitting!")
		}

		project.LockPwd = seedpwd
		project.Lock = false
		err = project.Save()
		if err != nil {
			return errors.Wrap(err, "couldn't save project")
		}
	} else if application == "livingunitcoop" {
		project, err := RetrieveLivingUnitCoop(projIndex)
		if err != nil || !project.Lock {
			return errors.Wrap(err, "couldn't retrieve living unit coop")
		}

		if recipient.U.Index != project.RecipientIndex {
			return errors.New("Recipient Indices don't match, quitting!")
		}

		project.LockPwd = seedpwd
		project.Lock = false
		err = project.Save()
		if err != nil {
			return errors.Wrap(err, "couldn't save project")
		}
	}

	return nil
}

func (project *ConstructionBond) updateConstructionBondAfterInvestment(invAmount float64, invIndex int) error {
	project.InvestorIndices = append(project.InvestorIndices, invIndex)
	project.AmountRaised += invAmount
	return project.Save()
}

func (project *LivingUnitCoop) updateLivingUnitCoopAfterInvestment(invAmount float64, invIndex int) error {
	project.ResidentIndices = append(project.ResidentIndices, invIndex)
	project.UnitsSold += 1
	return project.Save()
}
*/
