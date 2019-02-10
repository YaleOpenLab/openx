package ozones

import (
	"fmt"
	"log"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	utils "github.com/YaleOpenLab/openx/utils"
)

// TODO: add the recipient's role here, whether to give him an asset or do nothing

// Invest in a particular living coop
func (a *Coop) Invest(issuerPublicKey string, issuerSeed string, investor *database.Investor,
	investmentAmountS string, investorSeed string) error {
	// we want to invest in this specific bond
	var err error
	investmentAmount := utils.StoI(investmentAmountS)
	// check if investment amount is greater than the cost of a unit
	if float64(investmentAmount) > a.MonthlyPayment || float64(investmentAmount) < a.MonthlyPayment {
		fmt.Println("You are trying to invest more or less than a month's payment")
		return fmt.Errorf("You are trying to invest more or less than a month's payment")
	}
	assetName := assets.AssetID(a.Params.MaturationDate + a.Params.SecurityType + a.Params.Rating + a.Params.BondIssuer) // get a unique assetID

	if a.Params.InvestorAssetCode == "" {
		// this person is the first investor, set the investor token name
		InvestorAssetCode := assets.AssetID(consts.CoopAssetPrefix + assetName)
		a.Params.InvestorAssetCode = InvestorAssetCode             // set the investeor code
		_ = assets.CreateAsset(InvestorAssetCode, issuerPublicKey) // create the asset itself, since it would not have bene created earlier
	}
	/*
		dont check stableUSD balance for now
		if !investor.CanInvest(investor.U.PublicKey, investmentAmountS) {
			log.Println("Investor has less balance than what is required to ivnest in this asset")
			return a, err
		}
	*/
	// make investor trust the asset that we provide
	txHash, err := assets.TrustAsset(a.Params.InvestorAssetCode, issuerPublicKey, utils.FtoS(a.TotalAmount), investor.U.PublicKey, investorSeed)
	// trust upto the total value of the asset
	if err != nil {
		return err
	}
	log.Println("Investor trusted asset: ", a.Params.InvestorAssetCode, " tx hash: ", txHash)
	log.Println("Sending INVAsset: ", a.Params.InvestorAssetCode, "for: ", investmentAmount)
	_, txHash, err = assets.SendAssetFromIssuer(a.Params.InvestorAssetCode, investor.U.PublicKey, investmentAmountS, issuerSeed, issuerPublicKey)
	if err != nil {
		return err
	}
	log.Printf("Sent INVAsset %s to investor %s with txhash %s", a.Params.InvestorAssetCode, investor.U.PublicKey, txHash)
	// investor asset sent, update a.Params's BalLeft
	a.UnitsSold += 1
	investor.AmountInvested += float64(investmentAmount)
	investor.InvestedCoops = append(investor.InvestedCoops, a.Params.InvestorAssetCode)
	err = investor.Save() // save investor creds now that we're done
	if err != nil {
		return err
	}
	a.Residents = append(a.Residents, *investor)
	err = a.Save()
	return err
}

// Invest in a particular construction bonm
func (a *ConstructionBond) Invest(issuerPublicKey string, issuerSeed string, investor *database.Investor,
	recipient *database.Recipient, investmentAmountS string, investorSeed string, recipientSeed string) error {
	// we want to invest in this specific bond
	var err error
	investmentAmount := utils.StoI(investmentAmountS)
	// check if investment amount is greater than the cost of a unit
	if float64(investmentAmount) > a.CostOfUnit {
		return fmt.Errorf("You are trying to invest more than a unit's cost, do you want to invest in two units?")
	}
	assetName := assets.AssetID(a.Params.MaturationDate + a.Params.SecurityType + a.Params.Rating + a.Params.BondIssuer) // get a unique assetID

	if a.Params.InvestorAssetCode == "" {
		// this person is the first investor, set the investor token name
		InvestorAssetCode := assets.AssetID(consts.BondAssetPrefix + assetName)
		a.Params.InvestorAssetCode = InvestorAssetCode             // set the investeor code
		_ = assets.CreateAsset(InvestorAssetCode, issuerPublicKey) // create the asset itself, since it would not have bene created earlier
	}
	/*
		dont check stableUSD balance for now
		if !investor.CanInvest(investor.U.PublicKey, investmentAmountS) {
			log.Println("Investor has less balance than what is required to ivnest in this asset")
			return a, err
		}
	*/
	// make investor trust the asset that we provide
	txHash, err := assets.TrustAsset(a.Params.InvestorAssetCode, issuerPublicKey, utils.FtoS(a.CostOfUnit*float64(a.NoOfUnits)), investor.U.PublicKey, investorSeed)
	// trust upto the total value of the asset
	if err != nil {
		return err
	}
	log.Println("Investor trusted asset: ", a.Params.InvestorAssetCode, " tx hash: ", txHash)
	log.Println("Sending INVAsset: ", a.Params.InvestorAssetCode, "for: ", investmentAmount)
	_, txHash, err = assets.SendAssetFromIssuer(a.Params.InvestorAssetCode, investor.U.PublicKey, investmentAmountS, issuerSeed, issuerPublicKey)
	if err != nil {
		return err
	}
	log.Printf("Sent INVAsset %s to investor %s with txhash %s", a.Params.InvestorAssetCode, investor.U.PublicKey, txHash)
	// investor asset sent, update a.Params's BalLeft
	a.AmountRaised += float64(investmentAmount)
	investor.AmountInvested += float64(investmentAmount)
	investor.InvestedBonds = append(investor.InvestedBonds, a.Params.InvestorAssetCode)
	err = investor.Save() // save investor creds now that we're done
	if err != nil {
		return err
	}
	a.Investors = append(a.Investors, *investor)
	err = a.Save()
	return err
}
