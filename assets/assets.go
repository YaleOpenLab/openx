// package assets contains asset related functions like calculating AssetID and
// sets up DEBTokens, PBTokens and INVTokens for a specific project that it has been
// passed
// the entities in the system are described in the README file and this part
// will explain how the PBTokens, INVTokens and DEBTokens work.
// 1. INVToken - An INVToken is issued by the issuer for every USD that the investor
// has invested in the contract. This peg needs to be ensured maybe in protocol
// with stablecoins on Stellar or we need to provide an easy onboarding scheme
// for users into the crypto world using other means. The investor receives
// INVTokens as proof of investment but profit return mechanism is not taken into
// account here, since htat needs clear definition on how much investors get each
// period for investing in the project.
// TODO: INVTokens should be set with an
// immutable flag so that the isuser can't renege on issuing this assets at any
// future time
// 2. DEBToken - for each INVToken (and indirectly, USD invested in the project),
// we issue a DEBToken to the recipient of the assets so that they can pay us back.
// DEBTokens are also lunked with PBTokens and they should be immutable as well,
// so that the issuer can not change the amount of debt at any point in the future.
// 3. PBToken - each PBToken denoted a month of appropriate payback. A month's worth
// of payback is decided by the recipient, who decides the payback period of the
// given assets at the time of creation. PBTokens are non-fungible, it means
// that one project's payback token is not worth the same as the other project's PBToken.
// the other two tokens are fungible - each INVToken is worth +1USD and each DEBToken
// is worth -1 USD and can be transferred to other peers willing to take profit / debt
// on behalf of the above entities. SInce PBToken is not fungible, the flag
// authorization_required needs to be set and a party without a trustline with
// the issuer can not trade in this asset (and ideally, the issuer will not accept
// trustlines in this new asset)
// The hard part is ensuring that the assets are pegged to the USD in a stable way.
// we could ensure the peg ourselves by accepting USD off chain, but that's not provable
// on chain and the investor has to trust the issuer with that. Also, in this case,
// anonymous investors wouldn't be able to invest, which is something that would be
// nice to have
// TODO: Add flags to assets, onboarding
package assets

import (
	"fmt"
	"log"

	consts "github.com/OpenFinancing/openfinancing/consts"
	database "github.com/OpenFinancing/openfinancing/database"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	utils "github.com/OpenFinancing/openfinancing/utils"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
	"github.com/stellar/go/build"
)

// AssetID assigns a unique assetID to each asset. We assume that there won't be more
// than 68719476736 (16^9) assets that are created at any point, so we're good.
// the total AssetID must be less than 12 characters in length, so we take the first
// three for a human readable identifier and then the last 9 are random hex characaters
// passed through SHA3
func AssetID(inputString string) string {
	// so the assetID right now is a hash of the asset name, concatenated investor public keys and nonces
	x := utils.SHA3hash(inputString)
	return "YOL" + x[64:73] // max length of an asset in stellar is 12
	// log.Fatal(fmt.Errorf("All good"))
	// return nil
}

// CreateAsset creates a new asset belonging to the public key referenced above
func CreateAsset(assetName string, PublicKey string) build.Asset {
	// need to set a couple flags here
	return build.CreditAsset(assetName, PublicKey)
}

// TrustAsset trusts a specific asset issued by a particular public key and signs
// a transaction with a preset limit on how much it is willing to trsut that issuer's
// asset for
func TrustAsset(asset build.Asset, limit string, PublicKey string, Seed string) (string, error) {
	// TRUST is FROM recipient TO issuer
	trustTx, err := build.Transaction(
		build.SourceAccount{PublicKey},
		build.AutoSequence{SequenceProvider: xlm.TestNetClient},
		build.TestNetwork,
		build.Trust(asset.Code, asset.Issuer, build.Limit(limit)),
	)

	_, txHash, err := xlm.SendTx(Seed, trustTx)
	return txHash, err
}

// SendAsset transfers _amount_ number of assets from the caller to the destination
// and returns an error if the destination doesn't have a trustline with the issuer
// This method is called by the issuer of the asset
func SendAssetFromIssuer(assetName string, destination string, amount string, Seed string, PublicKey string) (int32, string, error) {
	// this transaction is FROM issuer TO recipient
	paymentTx, err := build.Transaction(
		build.SourceAccount{PublicKey},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: xlm.TestNetClient},
		build.MemoText{"Sending Asset: " + assetName},
		build.Payment(
			build.Destination{AddressOrSeed: destination},
			build.CreditAmount{assetName, PublicKey, amount},
			// CreditAmount identifies the asset by asset Code and issuer pubkey
		),
	)

	if err != nil {
		return -1, "", err
	}
	return xlm.SendTx(Seed, paymentTx)
}

// InvestInProject invests in a particular oder issued by _issuer_ with seed _issuerSeed_
// the _investor_ decides to invest _investmentAmountS_ amount of USD Tokens in
// a particular _uContract.Params_. If the invested amount makes the money raised equal to
// the total value of the _uContract.Params_, we issue the PBTokens and DEBTokens to the
// _recipient_
func InvestInProject(issuerPublicKey string, issuerSeed string, investor *database.Investor, recipient *database.Recipient, investmentAmountS string, uContract solar.SolarProject, investorSeed string, recipientSeed string) (solar.SolarProject, error) {
	var err error

	// invest only in integer values as of now, TODO: change to float
	investmentAmount := utils.StoI(investmentAmountS)
	// check if investment amount is greater than or equal to the project requirements
	amtLeft := uContract.Params.TotalValue - uContract.Params.MoneyRaised
	if investmentAmount > amtLeft {
		fmt.Println("User is trying to invest more than what is needed, print and exit")
		return uContract, fmt.Errorf("User is trying to invest more than what is needed, print and exit")
	}

	// user has decided to invest in a part of the project (don't know if full yet)
	// so if there has been no token codes assigned yet, we need to create them and
	// assign them here
	// you can retrieve these anywhere since the metadata will most likely be unique
	assetName := AssetID(uContract.Params.Metadata)
	if uContract.Params.INVAssetCode == "" {
		// this person is the first investor, set the investor token name
		INVAssetCode := AssetID(consts.INVAssetPrefix + assetName)
		uContract.Params.INVAssetCode = INVAssetCode   // set the investeor code
		_ = CreateAsset(INVAssetCode, issuerPublicKey) // create the asset itself, since it would not have bene created earlier
	}
	// we should check here whether the investor has enough USDTokens in project to be
	// able to ivnest in the asset
	if !investor.CanInvest(investor.U.PublicKey, investmentAmountS) {
		log.Println("Investor has less balance than what is required to ivnest in this asset")
		return uContract, err
	}
	var INVAsset build.Asset
	INVAsset.Code = uContract.Params.INVAssetCode
	INVAsset.Issuer = issuerPublicKey
	// INVAsset is not a native token, so don't set that
	// now we need to send the investor the INVAssets as proof of investment
	txHash, err := TrustAsset(INVAsset, utils.ItoS(uContract.Params.TotalValue), investor.U.PublicKey, investorSeed)
	// trust upto the total value of the asset
	if err != nil {
		return uContract, err
	}
	log.Println("Investor trusted asset: ", INVAsset.Code, " tx hash: ", txHash)
	log.Println("Sending INVAsset: ", INVAsset.Code, "for: ", investmentAmount)
	_, txHash, err = SendAssetFromIssuer(INVAsset.Code, investor.U.PublicKey, investmentAmountS, issuerSeed, issuerPublicKey)
	if err != nil {
		return uContract, err
	}
	log.Printf("Sent INVAsset %s to investor %s with txhash %s", INVAsset.Code, investor.U.PublicKey, txHash)
	// investor asset sent, update uContract.Params's BalLeft
	uContract.Params.MoneyRaised += investmentAmount
	fmt.Println("Updating investor to handle invested amounts and assets")
	investor.AmountInvested += float64(investmentAmount)
	investor.InvestedAssets = append(investor.InvestedAssets, uContract.Params.DEBAssetCode)
	err = investor.Save() // save investor creds now that we're done
	if err != nil {
		return uContract, err
	}
	fmt.Println("Updated investor database")
	// append the investor class to the list of project investors
	// if the same investor has invested twice, he will appear twice
	// can be resolved on the UI side by requiring unique, so not doing that here
	uContract.Params.ProjectInvestors = append(uContract.Params.ProjectInvestors, *investor)
	if uContract.Params.MoneyRaised == uContract.Params.TotalValue {
		// this project covers up the amount nedeed for the project, so set the DEBAssetCode
		// and PBAssetCodes, generate them and give to the recipient
		uContract.Params.DEBAssetCode = AssetID(consts.DEBAssetPrefix + assetName)
		uContract.Params.PBAssetCode = AssetID(consts.PBAssetPrefix + assetName)
		DEBasset := CreateAsset(uContract.Params.DEBAssetCode, issuerPublicKey)
		PBasset := CreateAsset(uContract.Params.PBAssetCode, issuerPublicKey)
		// and the school needs to trust me only for paybackTokens amount of PB tokens
		pbAmtTrust := utils.ItoS(uContract.Params.Years * 12 * 2) // two way exchange possible, to account for errors
		txHash, err = TrustAsset(PBasset, pbAmtTrust, recipient.U.PublicKey, recipientSeed)
		if err != nil {
			return uContract, err
		}
		log.Println("Recipient Trusted Payback asset: ", PBasset.Code, " tx hash: ", txHash)

		txHash, err = TrustAsset(DEBasset, utils.ItoS(uContract.Params.TotalValue*2), recipient.U.PublicKey, recipientSeed) // since debt = invested amount
		// *2 is for sending the amount back
		if err != nil {
			return uContract, err
		}
		log.Println("Recipient Trusted Debt asset: ", DEBasset.Code, " tx hash: ", txHash)
		log.Println("Sending DEBasset: ", uContract.Params.DEBAssetCode)
		_, txHash, err = SendAssetFromIssuer(uContract.Params.DEBAssetCode, recipient.U.PublicKey, utils.ItoS(uContract.Params.TotalValue), issuerSeed, issuerPublicKey) // same amount as debt
		if err != nil {
			return uContract, err
		}
		log.Printf("Sent DEBasset to recipient %s with txhash %s", recipient.U.PublicKey, txHash)
		uContract.Params.BalLeft = float64(uContract.Params.TotalValue)
		recipient.ReceivedSolarProjects = append(recipient.ReceivedSolarProjects, uContract.Params.DEBAssetCode)
		uContract.Params.ProjectRecipient = *recipient // need to udpate uContract.Params each time recipient is mutated
		// only here does the recipient part change, so update it only here
		// TODO: keep note of who all invested in this asset (even though it should be
		// easy to get that from the blockchain)
		if uContract.Params.DEBAssetCode == "" {
			log.Fatal("Empty debt asset code")
		}
		err = recipient.Save()
		if err != nil {
			return uContract, err
		}
		uContract.Stage = 4 // 4 is the funded stage
		err = uContract.Save()
		if err != nil {
			log.Println("Couldn't insert project")
			return uContract, err
		}
		fmt.Println("Updated recipient bucket")
		return uContract, nil
	}
	// update the project finally now that we have updated other databases
	err = uContract.Save()
	return uContract, err
}

func SendPBAsset(project solar.SolarProject, destination string, amount string, Seed string, PublicKey string) error {
	// need to calculate how much PBAsset we need to send back.
	amountS := project.CalculatePayback(amount)
	_, txHash, err := SendAssetFromIssuer(project.Params.PBAssetCode, destination, amountS, Seed, PublicKey)
	log.Println("TXHASH for payback is: ", txHash)
	return err
}
