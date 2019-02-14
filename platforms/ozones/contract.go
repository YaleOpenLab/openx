package ozones

import (
	"fmt"
	"log"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	issuer "github.com/YaleOpenLab/openx/issuer"
	model "github.com/YaleOpenLab/openx/models/debtcrowdfunding"
	utils "github.com/YaleOpenLab/openx/utils"
)

// TODO: add the recipient's role here, whether to give him an asset or do nothing

// Invest in a particular living coop
func (a *LivingUnitCoop) Invest(issuerPubkey string, issuerSeed string, investor *database.Investor,
	invAmountS string, invSeed string) error {
	// we want to invest in this specific bond
	var err error
	invAmount := utils.StoI(invAmountS)
	// check if investment amount is greater than the cost of a unit
	if float64(invAmount) > a.MonthlyPayment || float64(invAmount) < a.MonthlyPayment {
		fmt.Println("You are trying to invest more or less than a month's payment")
		return fmt.Errorf("You are trying to invest more or less than a month's payment")
	}
	assetName := assets.AssetID(a.MaturationDate + a.SecurityType + a.Rating + a.BondIssuer) // get a unique assetID

	if a.InvestorAssetCode == "" {
		// this person is the first investor, set the investor token name
		InvestorAssetCode := assets.AssetID(consts.CoopAssetPrefix + assetName)
		a.InvestorAssetCode = InvestorAssetCode                 // set the investeor code
		_ = assets.CreateAsset(InvestorAssetCode, issuerPubkey) // create the asset itself, since it would not have bene created earlier
	}

	if !investor.CanInvest(invAmountS) {
		log.Println("Investor has less balance than what is required to ivnest in this asset")
		return err
	}

	// investor can invest in this project, send stablecoin to the platform
	txHash, err := assets.TrustAsset(a.InvestorAssetCode, issuerPubkey, utils.FtoS(a.Amount), investor.U.PublicKey, invSeed)
	if err != nil {
		return err
	}

	log.Println("Investor trusted asset: ", a.InvestorAssetCode, " tx hash: ", txHash)
	log.Println("Sending INVAsset: ", a.InvestorAssetCode, "for: ", invAmount)
	_, txHash, err = assets.SendAssetFromIssuer(a.InvestorAssetCode, investor.U.PublicKey, invAmountS, issuerSeed, issuerPubkey)
	if err != nil {
		return err
	}
	log.Printf("Sent INVAsset %s to investor %s with txhash %s", a.InvestorAssetCode, investor.U.PublicKey, txHash)
	// investor asset sent, update a.Params's BalLeft
	a.UnitsSold += 1
	investor.AmountInvested += float64(invAmount)
	investor.InvestedCoops = append(investor.InvestedCoops, a.InvestorAssetCode)
	err = investor.Save() // save investor creds now that we're done
	if err != nil {
		return err
	}
	a.Residents = append(a.Residents, *investor)
	err = a.Save()
	return err
}

func preInvestmentCheck(invIndex int, projIndex int, invAmount string) (ConstructionBond, error) {

	project, err := RetrieveConstructionBond(projIndex)
	if err != nil {
		return project, err
	}

	investor, err := database.RetrieveInvestor(invIndex)
	if err != nil {
		return project, err
	}
	// check if investment amount is greater than the cost of a unit
	if float64(utils.StoF(invAmount)) != project.CostOfUnit {
		return project, fmt.Errorf("You are trying to invest more than a unit's cost, do you want to invest in two units?")
	}

	assetName := assets.AssetID(project.MaturationDate + project.SecurityType + project.Rating + project.BondIssuer) // get a unique assetID

	if len(project.InvestorIndices) == 0 {
		// initialize issuer
		err = issuer.InitIssuer(consts.OpzonesIsuserDir, project.Index, consts.IssuerSeedPwd)
		if err != nil {
			log.Println("Error while initializing issuer", err)
			return project, err
		}
		err = issuer.FundIssuer(consts.OpzonesIsuserDir, project.Index, consts.IssuerSeedPwd, consts.PlatformSeed)
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

// Invest in a particular construction bonm
func InvestInConstructionBond(projIndex int, invIndex int, recpIndex int, invAmount string,
	invSeed string, recpSeed string) error {
	// we want to invest in this specific bond
	var err error

	project, err := preInvestmentCheck(projIndex, invIndex, invAmount)
	if err != nil {
		log.Println(err)
		return err
	}

	trustLimit := utils.FtoS(project.CostOfUnit * float64(project.NoOfUnits))

	err = model.Invest(projIndex, invIndex, project.InvestorAssetCode, invSeed,
		invAmount, trustLimit, project.InvestorIndices)
	if err != nil {
		log.Println(err)
		return err
	}

	err = project.updateConstructionBondAfterInvestment(invAmount, invIndex, recpIndex)
	if err != nil {
		log.Println("Failed to update project after investment", err)
		return err
	}

	totalValue := float64(project.CostOfUnit * float64(project.NoOfUnits))
	project.DebtAssetCode = assets.AssetID(consts.DebtAssetPrefix + project.Description)

	if totalValue == project.AmountRaised {
		// send the recipient relevant debt asset
		err = model.Receive(consts.OpzonesIsuserDir, recpIndex, projIndex, project.DebtAssetCode, recpSeed, totalValue)
		if err != nil {
			log.Println("Failed to send assets to recipient project after investment", err)
			return err
		}
	}
	return err
}

func (project *ConstructionBond) updateConstructionBondAfterInvestment(invAmount string, invIndex int, recpIndex int) error {
	project.InvestorIndices = append(project.InvestorIndices, invIndex)
	// TODO: have the amount in escrow or something
	project.AmountRaised += utils.StoF(invAmount)
	project.RecipientIndex = recpIndex
	return project.Save()
}
