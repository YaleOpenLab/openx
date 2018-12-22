// package assets contains asset related functions like calculating AssetID and
// sets up DEBTokens, PBTokens and INVTokens for a specific order that it has been
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
// that one order's payback token is not worth the same as the other order's PBToken.
// the other two tokens are fungible - each INVToken is worth +1USD and each DEBToken
// is worth -1 USD and can be trnasferred to other peers willing to take profit / debt
// on behalf of the above entities. SInce PBToken is not fungible, the flag
// authorization_required needs to be set and a party without a trustline with
// the issuer can not trade in this asset (and ideally, the issuer will not accept
// trustlines in this new asset)
// Supported payback periods right now are
// A. 3 YEARS = 36 PBTokens
// B. 5 YEARS = 60 PBTokens
// C. 7 YEARS = 84 PBTokens
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
	"strconv"

	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
	"github.com/stellar/go/build"
)

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
		build.AutoSequence{SequenceProvider: utils.DefaultTestNetClient},
		build.TestNetwork,
		build.Trust(asset.Code, asset.Issuer, build.Limit(limit)),
	)

	if err != nil {
		return "", err
	}

	trustTxe, err := trustTx.Sign(Seed)
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

// SendAsset transfers _amount_ number of assets from the caller to the destination
// and returns an error if the destination doesn't have a trustline with the issuer
// This method is called by the issuer of the asset
func SendAssetFromIssuer(assetName string, destination string, amount string, Seed string, PublicKey string) (int32, string, error) {
	// this transaction is FROM issuer TO recipient
	paymentTx, err := build.Transaction(
		build.SourceAccount{PublicKey},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: utils.DefaultTestNetClient},
		build.Payment(
			build.Destination{AddressOrSeed: destination},
			build.CreditAmount{assetName, PublicKey, amount},
			// build.MemoText{"Sending Solar Asset"}, // apparently we
			// can put whatever we want here, but it doesn't work
			// CreditAmount identifies the asset by asset Code and issuer pubkey
		),
	)

	if err != nil {
		return -1, "", err
	}

	paymentTxe, err := paymentTx.Sign(Seed)
	if err != nil {
		return -1, "", err
	}

	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {
		return -1, "", err
	}

	tx, err := utils.DefaultTestNetClient.SubmitTransaction(paymentTxeB64)
	if err != nil {
		return -1, "", err
	}

	return tx.Ledger, tx.Hash, nil
}

 // InvestInOrder invests in a particular oder issued by _issuer_ wtih seed _issuerSeed_
 // the _investor_ decides to invest _investmentAmountS_ amount of USD Tokens in
 // a particular _uOrder_. If the invested amount makes the money raised equal to
 // the total value of the _uOrder_, we issue the PBTokens and DEBTokens to the
 // _recipient_
func InvestInOrder(issuer *database.Platform, issuerSeed string, investor *database.Investor, recipient *database.Recipient, investmentAmountS string, uOrder database.Order) (database.Order, error) {
	var partOrder database.Order
	var err error

	// invest only in integer values as of now, TODO: change to float
	investmentAmount := utils.StoI(investmentAmountS)
	//  check if investment amount is greater than or equal to the order requirements
	amtLeft := uOrder.TotalValue - uOrder.MoneyRaised
	if investmentAmount > amtLeft {
		fmt.Println("User is trying to invest more thna what is needed, print and exit")
		return partOrder, fmt.Errorf("User is trying to invest more thna what is needed, print and exit")
	}

	// user has decided to invest in a part of the order (don't know if full yet)
	// so if there has been no token codes assigned yet, we need to create them and
	// assign them here
	// you can retrieve these anywhere since the metadata will mostt likely be unique
	assetName := AssetID(uOrder.Metadata)
	if uOrder.INVAssetCode == "" {
		// this person is the first investor, set the investor token name
		INVAssetCode := AssetID("INVTokens_" + assetName)
		uOrder.INVAssetCode = INVAssetCode              // set the investeor code
		_ = CreateAsset(INVAssetCode, issuer.PublicKey) // create the asset itself, since it would not have bene created earlier
	}
	// we should check here whether the investor has enough USDTokens in order to be
	// able to ivnest in the asset
	err = xlm.GetUSDTokenBalance(investor.U.PublicKey, investmentAmountS)
	if err != nil {
		log.Println("Investor has less balance than what is required to ivnest in this asset")
		return uOrder, err
	}
	var INVAsset build.Asset
	INVAsset.Code = uOrder.INVAssetCode
	INVAsset.Issuer = issuer.PublicKey
	// INVAsset is not a native token, so don't set that
	// now we need to send the investor the INVAssets as proof of investment
	txHash, err := TrustAsset(INVAsset, utils.IntToString(uOrder.TotalValue), investor.U.PublicKey, investor.U.Seed)
	// trust upto the total value of the asset
	if err != nil {
		return uOrder, err
	}
	log.Println("Investor trusted asset: ", INVAsset.Code, " tx hash: ", txHash)
	log.Println("Sending INVAsset: ", INVAsset.Code, "for: ", investmentAmount)
	_, txHash, err = SendAssetFromIssuer(INVAsset.Code, investor.U.PublicKey, strconv.Itoa(investmentAmount), issuerSeed, issuer.PublicKey)
	if err != nil {
		return uOrder, err
	}
	log.Printf("Sent INVAsset %s to investor %s with txhash %s", INVAsset.Code, investor.U.PublicKey, txHash)
	// investor asset sent, update uOrder's BalLeft
	uOrder.MoneyRaised += investmentAmount
	fmt.Println("Updating investor to handle invested amounts and assets")
	investor.AmountInvested += float64(investmentAmount)
	investor.InvestedAssets = append(investor.InvestedAssets, uOrder)
	err = database.InsertInvestor(*investor) // save investor creds now that we're done
	if err != nil {
		return uOrder, err
	}
	fmt.Println("Updated investor database")
	if uOrder.MoneyRaised == uOrder.TotalValue {
		// this order covers up the amount nedeed for the order, so set the DEBAssetCode
		// and PBAssetCodes, generate them and give to the recipient
		DEBAssetCode := AssetID("DEBTokens_" + assetName)
		PBAssetCode := AssetID("PBTokens_" + assetName)
		DEBasset := CreateAsset(DEBAssetCode, issuer.PublicKey)
		PBasset := CreateAsset(PBAssetCode, issuer.PublicKey)
		// and the school needs to trust me only for paybackTokens amount of PB tokens
		pbAmtTrust := utils.IntToString(uOrder.Years * 24)
		txHash, err = TrustAsset(PBasset, pbAmtTrust, recipient.U.PublicKey, recipient.U.Seed)
		if err != nil {
			return uOrder, err
		}
		log.Println("Recipient Trusted Payback asset: ", PBasset.Code, " tx hash: ", txHash)

		txHash, err = TrustAsset(DEBasset, strconv.Itoa(uOrder.TotalValue * 2), recipient.U.PublicKey, recipient.U.Seed) // since debt = invested amount
		// *2 is for sending the amount back
		if err != nil {
			return uOrder, err
		}
		log.Println("Recipient Trusted Debt asset: ", DEBasset.Code, " tx hash: ", txHash)
		log.Println("Sending DEBasset: ", DEBAssetCode)
		_, txHash, err = SendAssetFromIssuer(DEBAssetCode, recipient.U.PublicKey, strconv.Itoa(uOrder.TotalValue), issuerSeed, issuer.PublicKey) // same amount as debt
		if err != nil {
			return uOrder, err
		}
		log.Printf("Sent DEBasset to recipient %s with txhash %s", recipient.U.PublicKey, txHash)
		uOrder.Live = true
		uOrder.DEBAssetCode = DEBAssetCode
		uOrder.PBAssetCode = PBAssetCode
		uOrder.BalLeft = float64(uOrder.TotalValue)
		recipient.ReceivedOrders = append(recipient.ReceivedOrders, uOrder)
		uOrder.OrderRecipient = *recipient // need to udpate uOrder each time recipient is mutated
		// only here does the recipient part change, so update it only here
		// TODO: keep note of who all invested in this asset (even though it should be
		// easy to get that from the blockchain)
		log.Println("UORDER BEFORE DB: ", uOrder.DEBAssetCode)
		if uOrder.DEBAssetCode == "" {
			log.Fatal("DOnt work")
		}
		err = database.InsertRecipient(*recipient)
		if err != nil {
			return uOrder, err
		}
		err = database.InsertOrder(uOrder)
		if err != nil {
			log.Println("Couldn't insert order")
			return uOrder, err
		}
		fmt.Println("Updated recipient bucket")
		return uOrder, nil
	}
	// update the order finally now that we have updated other databases
	err = database.InsertOrder(uOrder)
	return uOrder, err
}
