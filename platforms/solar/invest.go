package solar

import (
	"fmt"
	"log"

	assets "github.com/OpenFinancing/openfinancing/assets"
	consts "github.com/OpenFinancing/openfinancing/consts"
	database "github.com/OpenFinancing/openfinancing/database"
	issuer "github.com/OpenFinancing/openfinancing/issuer"
	utils "github.com/OpenFinancing/openfinancing/utils"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	"github.com/stellar/go/build"
)

func RetrieveValues(projIndex int, invIndex int, recpIndex int) (Project, database.Investor, database.Recipient, error) {
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
	return project, investor, recipient, nil
}

// this file does not contain any tests associated with it right now. In the future,
// once we have a robust frontend, we can modify the CLI interface to act as a test
// for this file

// InvestInProject invests in a particular solar project given required parameters
func InvestInProject(projIndex int, invIndex int, recpIndex int, invAmountS string,
	invSeed string, recpSeed string, platformSeed string) (Project, error) {
	var err error

	project, investor, recipient, err := RetrieveValues(projIndex, invIndex, recpIndex)
	if err != nil {
		return project, err
	}

	invAmount := utils.StoF(invAmountS)
	// check if investment amount is greater than or equal to the project requirements
	if invAmount > project.Params.TotalValue-project.Params.MoneyRaised {
		return project, fmt.Errorf("User is trying to invest more than what is needed")
	}

	var InvestorAsset build.Asset
	var PaybackAsset build.Asset
	var DebtAsset build.Asset
	// user has decided to invest in a part of the project (don't know if full yet)
	// so if there has been no asset codes assigned yet, we need to create them and
	// assign them here
	// you can retrieve these anywhere since the metadata will most likely be unique
	if project.Params.InvestorAssetCode == "" {
		// this person is the first investor, set the investor asset name and create the
		// issuer that will be created for this particular project
		project.Params.InvestorAssetCode = assets.AssetID(consts.InvestorAssetPrefix + project.Params.Metadata) // set the investor asset code
		err = issuer.InitIssuer(project.Params.Index, consts.IssuerSeedPwd)
		if err != nil {
			log.Fatal(err)
		}
		err = issuer.FundIssuer(project.Params.Index, consts.IssuerSeedPwd, platformSeed)
		if err != nil {
			log.Fatal(err)
		}
	}

	issuerPubkey, issuerSeed, err := wallet.RetrieveSeed(issuer.CreatePath(project.Params.Index), consts.IssuerSeedPwd)
	if err != nil {
		return project, err
	}

	// InvAsset is not a native asset, so don't set the native flag
	InvestorAsset = assets.CreateAsset(project.Params.InvestorAssetCode, issuerPubkey)
	// make investor trust the asset, trustlimit is upto the value of the project
	txHash, err := assets.TrustAsset(InvestorAsset.Code, issuerPubkey, utils.FtoS(project.Params.TotalValue), investor.U.PublicKey, invSeed)
	if err != nil {
		return project, err
	}
	log.Println("Investor trusted asset: ", InvestorAsset.Code, " tx hash: ", txHash)
	_, txHash, err = assets.SendAssetFromIssuer(InvestorAsset.Code, investor.U.PublicKey, invAmountS, issuerSeed, issuerPubkey)
	if err != nil {
		return project, err
	}
	log.Printf("Sent InvAsset %s to investor %s with txhash %s", InvestorAsset.Code, investor.U.PublicKey, txHash)
	// investor asset sent, update project.Params's BalLeft
	fmt.Println("Updating investor to handle invested amounts and assets")
	project.Params.MoneyRaised += invAmount
	investor.AmountInvested += float64(invAmount)
	investor.InvestedSolarProjects = append(investor.InvestedSolarProjects, InvestorAsset.Code)
	// keep note of who all invested in this asset (even though it should be easy
	// to get that from the blockchain)
	err = investor.Save() // save investor creds now that we're done
	if err != nil {
		return project, err
	}
	fmt.Println("Updated investor database")
	// append the investor class to the list of project investors
	// if the same investor has invested twice, he will appear twice
	// can be resolved on the UI side by requiring unique, so not doing that here
	project.ProjectInvestors = append(project.ProjectInvestors, investor)
	if project.Params.MoneyRaised == project.Params.TotalValue {
		// this project covers up the amount nedeed for the project, so set the DebtAssetCode
		// and PaybackAssetCodes, generate them and give to the recipient
		project.Params.DebtAssetCode = assets.AssetID(consts.DebtAssetPrefix + project.Params.Metadata)
		project.Params.PaybackAssetCode = assets.AssetID(consts.PaybackAssetPrefix + project.Params.Metadata)
		DebtAsset = assets.CreateAsset(project.Params.DebtAssetCode, issuerPubkey)
		PaybackAsset = assets.CreateAsset(project.Params.PaybackAssetCode, issuerPubkey)
		pbAmtTrust := utils.ItoS(project.Params.Years * 12 * 2) // two way exchange possible, to account for errors

		txHash, err = assets.TrustAsset(PaybackAsset.Code, issuerPubkey, pbAmtTrust, recipient.U.PublicKey, recpSeed)
		if err != nil {
			return project, err
		}
		log.Println("Recipient Trusts Debt asset: ", DebtAsset.Code, " tx hash: ", txHash)
		_, txHash, err = assets.SendAssetFromIssuer(PaybackAsset.Code, recipient.U.PublicKey, pbAmtTrust, issuerSeed, issuerPubkey) // same amount as debt
		if err != nil {
			return project, err
		}
		log.Printf("Sent PaybackAsset to recipient %s with txhash %s", recipient.U.PublicKey, txHash)
		txHash, err = assets.TrustAsset(DebtAsset.Code, issuerPubkey, utils.FtoS(project.Params.TotalValue*2), recipient.U.PublicKey, recpSeed)
		if err != nil {
			return project, err
		}
		log.Println("Recipient Trusts Payback asset: ", PaybackAsset.Code, " tx hash: ", txHash)
		_, txHash, err = assets.SendAssetFromIssuer(DebtAsset.Code, recipient.U.PublicKey, utils.FtoS(project.Params.TotalValue), issuerSeed, issuerPubkey) // same amount as debt
		if err != nil {
			return project, err
		}
		log.Printf("Sent PaybackAsset to recipient %s with txhash %s\n", recipient.U.PublicKey, txHash)
		project.Params.BalLeft = float64(project.Params.TotalValue)
		recipient.ReceivedSolarProjects = append(recipient.ReceivedSolarProjects, DebtAsset.Code)
		project.ProjectRecipient = recipient // need to udpate project.Params each time recipient is mutated
		// only here does the recipient part change, so update it only here
		err = recipient.Save()
		if err != nil {
			return project, err
		}
		project.Stage = FundedProject // set funded project stage
		err = project.Save()
		if err != nil {
			log.Println("Couldn't insert project")
			return project, err
		}
		fmt.Println("Updated recipient bucket")
		txhash, err := issuer.FreezeIssuer(project.Params.Index, "blah")
		if err != nil {
			return project, err
		}
		log.Printf("Tx hash for freezing issuer is: %s", txhash)
	}
	// update the project finally now that we have updated other databases
	err = project.Save()
	return project, err
}
