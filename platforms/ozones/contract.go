package ozones

import (
	"fmt"
	"log"
	"time"
	"math"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	issuer "github.com/YaleOpenLab/openx/issuer"
	model "github.com/YaleOpenLab/openx/models/debtcrowdfunding"
	notif "github.com/YaleOpenLab/openx/notif"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
)

func preInvestmentConstructionBonds(projIndex int, invIndex int, invAmount string) (ConstructionBond, error) {

	project, err := RetrieveConstructionBond(projIndex)
	if err != nil {
		return project, err
	}

	investor, err := database.RetrieveInvestor(invIndex)
	if err != nil {
		return project, err
	}
	// check if investment amount is greater than the cost of a unit
	rem := float64(utils.StoF(invAmount)) / project.CostOfUnit
	if math.Floor(rem) == 0 {
		return project, fmt.Errorf("You are trying to invest more than a unit's cost, do you want to invest in two units?")
	}

	assetName := assets.AssetID(project.MaturationDate + project.SecurityType + project.Rating + project.BondIssuer) // get a unique assetID

	if len(project.InvestorIndices) == 0 {
		// initialize issuer
		err = issuer.InitIssuer(consts.OpzonesIssuerDir, project.Index, consts.IssuerSeedPwd)
		if err != nil {
			log.Println("Error while initializing issuer", err)
			return project, err
		}
		err = issuer.FundIssuer(consts.OpzonesIssuerDir, project.Index, consts.IssuerSeedPwd, consts.PlatformSeed)
		if err != nil {
			log.Println("Error while funding issuer", err)
			return project, err
		}

		project.InvestorAssetCode = assets.AssetID(consts.BondAssetPrefix + assetName) // set the investor code
	}

	if !investor.CanInvest(invAmount) {
		log.Println("Investor has less balance than what is required to ivnest in this asset")
		return project, err
	}

	return project, nil
}

func preInvestmentLivingCoop(projIndex int, invIndex int, invAmount string) (LivingUnitCoop, error) {

	project, err := RetrieveLivingUnitCoop(projIndex)
	if err != nil {
		return project, err
	}

	investor, err := database.RetrieveInvestor(invIndex)
	if err != nil {
		return project, err
	}
	// check if investment amount is greater than the cost of a unit
	if float64(utils.StoF(invAmount)) != project.MonthlyPayment {
		return project, fmt.Errorf("You are trying to invest more than a unit's cost, do you want to invest in two units?")
	}

	assetName := assets.AssetID(project.Description)

	if len(project.ResidentIndices) == 0 {
		// initialize issuer
		err = issuer.InitIssuer(consts.OpzonesIssuerDir, project.Index, consts.IssuerSeedPwd)
		if err != nil {
			log.Println("Error while initializing issuer", err)
			return project, err
		}
		err = issuer.FundIssuer(consts.OpzonesIssuerDir, project.Index, consts.IssuerSeedPwd, consts.PlatformSeed)
		if err != nil {
			log.Println("Error while funding issuer", err)
			return project, err
		}

		project.InvestorAssetCode = assets.AssetID(consts.BondAssetPrefix + assetName) // set the investor code
	}

	if !investor.CanInvest(invAmount) {
		log.Println("Investor has less balance than what is required to ivnest in this asset")
		return project, err
	}

	return project, nil
}

// Invest in a particular living coop
func InvestInLivingUnitCoop(projIndex int, invIndex int, invAmount string, invSeed string,
	recpSeed string) error {
	// we want to invest in this specific bond
	var err error

	project, err := preInvestmentLivingCoop(projIndex, invIndex, invAmount)
	if err != nil {
		log.Println(err)
		return err
	}

	err = model.Invest(projIndex, invIndex, project.InvestorAssetCode, invSeed,
		invAmount, utils.FtoS(project.Amount), project.ResidentIndices, "livingunitcoop")
	if err != nil {
		log.Println(err)
		return err
	}

	err = project.updateLivingUnitCoopAfterInvestment(invAmount, invIndex)
	if err != nil {
		log.Println("Failed to update project after investment", err)
		return err
	}

	return nil
}

// Invest in a particular construction bonm
func InvestInConstructionBond(projIndex int, invIndex int, invAmount string, invSeed string) error {
	// we want to invest in this specific bond
	var err error

	project, err := preInvestmentConstructionBonds(projIndex, invIndex, invAmount)
	if err != nil {
		log.Println(err)
		return err
	}

	trustLimit := utils.FtoS(project.CostOfUnit * float64(project.NoOfUnits))

	err = model.Invest(projIndex, invIndex, project.InvestorAssetCode, invSeed,
		invAmount, trustLimit, project.InvestorIndices, "constructionbond")
	if err != nil {
		log.Println(err)
		return err
	}

	err = project.updateConstructionBondAfterInvestment(invAmount, invIndex)
	if err != nil {
		log.Println("Failed to update project after investment", err)
		return err
	}

	totalValue := float64(project.CostOfUnit * float64(project.NoOfUnits))
	project.DebtAssetCode = assets.AssetID(consts.DebtAssetPrefix + project.Description)

	if totalValue == project.AmountRaised {
		// send the recipient a notification to unlock the specific project and accept the investment
		err = project.sendRecipientNotification()
		if err != nil {
			log.Println("Error while sending notifications to recipient", err)
			return err
		}
		go sendRecipientAssets(projIndex, totalValue)
	}
	return nil
}

func sendRecipientAssets(projIndex int, totalValue float64) error {
	// send the recipient relevant debt asset
	startTime := utils.Unix()
	project, err := RetrieveConstructionBond(projIndex)
	if err != nil {
		log.Println(err)
		return err
	}

	for utils.Unix()-startTime < consts.LockInterval {
		log.Printf("WAITING FOR PROJECT %d TO BE UNLOCKED", projIndex)
		project, err = RetrieveConstructionBond(projIndex)
		if err != nil {
			log.Println(err)
			return err
		}
		if !project.Lock {
			log.Println("Project UNLOCKED IN LOOP")
			break
		}
		time.Sleep(10 * time.Second)
	}

	project, err = RetrieveConstructionBond(projIndex)
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
		log.Println("Couldn't decrypt seed", err)
		return err
	}

	err = model.ReceiveBond(consts.OpzonesIssuerDir, project.RecipientIndex, projIndex, project.DebtAssetCode, recpSeed, totalValue)
	if err != nil {
		log.Println("Failed to send assets to recipient project after investment", err)
		return err
	}

	project.LockPwd = ""
	return project.Save()
}

// sendRecipientNotification sends the notification to the recipient requesting them
// to logon to the platform and unlock the project that has just been invested in
func (project *ConstructionBond) sendRecipientNotification() error {
	recipient, err := database.RetrieveRecipient(project.RecipientIndex)
	if err != nil {
		log.Println(err)
		return err
	}
	notif.SendUnlockNotifToRecipient(project.Index, recipient.U.Email)
	project.Lock = true
	return project.Save()
}

// sendRecipientNotification sends the notification to the recipient requesting them
// to logon to the platform and unlock the project that has just been invested in
func (project *LivingUnitCoop) sendDeveloperNotification() error {
	recipient, err := database.RetrieveRecipient(project.RecipientIndex)
	if err != nil {
		log.Println(err)
		return err
	}
	notif.SendUnlockNotifToRecipient(project.Index, recipient.U.Email)
	return nil
}

// UnlockProject unlocks a specific project that has just been invested in
func UnlockProject(username string, pwhash string, projIndex int, seedpwd string, application string) error {
	fmt.Println("UNLOCKING PROJECT")
	recipient, err := database.ValidateRecipient(username, pwhash)
	if err != nil {
		log.Println(err)
		return err
	}

	recpSeed, err := wallet.DecryptSeed(recipient.U.EncryptedSeed, seedpwd)
	if err != nil {
		log.Println("Error while decrpyting seed", err)
		return err
	}

	checkPubkey, err := wallet.ReturnPubkey(recpSeed)
	if err != nil {
		log.Println("Couldn't get pubkey from seed", err)
		return err
	}

	if checkPubkey != recipient.U.PublicKey {
		log.Println("Invalid seed")
		return fmt.Errorf("Failed to unlock project")
	}

	if application == "constructionbond" {
		project, err := RetrieveConstructionBond(projIndex)
		if err != nil || !project.Lock {
			log.Println(err)
			return err
		}

		if recipient.U.Index != project.RecipientIndex {
			return fmt.Errorf("Recipient Indices don't match, quitting!")
		}

		project.LockPwd = seedpwd
		project.Lock = false
		err = project.Save()
		if err != nil {
			log.Println(err)
			return err
		}
	} else if application == "livingunitcoop" {
		project, err := RetrieveLivingUnitCoop(projIndex)
		if err != nil || !project.Lock {
			log.Println(err)
			return err
		}

		if recipient.U.Index != project.RecipientIndex {
			return fmt.Errorf("Recipient Indices don't match, quitting!")
		}

		project.LockPwd = seedpwd
		project.Lock = false
		err = project.Save()
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (project *ConstructionBond) updateConstructionBondAfterInvestment(invAmount string, invIndex int) error {
	project.InvestorIndices = append(project.InvestorIndices, invIndex)
	// TODO: have the amount in escrow or something
	project.AmountRaised += utils.StoF(invAmount)
	return project.Save()
}

func (project *LivingUnitCoop) updateLivingUnitCoopAfterInvestment(invAmount string, invIndex int) error {
	project.ResidentIndices = append(project.ResidentIndices, invIndex)
	project.UnitsSold += 1
	return project.Save()
}
