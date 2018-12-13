// package assets contains asset related functions like calculating AssetID and
// sets up DEBTokens, PBTokens and INVTokens for a specific order that it has been
// passed
// the entities in the system are described in the README file and this part
// will explain how the PBTokens, INVTokens and DEBTokens work.
// 1. INVToken - An INVToken is issued by the issuer for every USD that the investor
// has invested in the contract. This peg needs to be ensured maybe in protocol
// with stablecoins on Stellar or we need to provide an easy onboarding scheme
// for users into the crypto worls using other means. The inevestor receives
// INVTokens as proof of investment but profit return mechanism is not taken into
// account here, since htat needs clear definition on how much investors get each
// period for inevesting in the project. TODO: INVTokens should be set with an
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
// TODO: Add flags to assets, onboarding, multiple investors and more
package assets

import (
	"log"

	"github.com/boltdb/bolt"
	accounts "github.com/YaleOpenLab/smartPropertyMVP/stellar/accounts"
	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
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

// CalculatePayback is a TODO function that should simply some up the PBToken
// balance and then return them to the frontend UI for a nice display
func CalculatePayback(balance string, noOfMonths string) error {
	// the idea is that we should be able ot pass an assetId to this function
	// and it must calculate how much time we have left for payback. For this example
	// until twe do the db stuff, lets pass a few params (although this could be done
	// separately as well).
	return nil
}

// SetupAsset sets up assets based on a given order
func SetupAsset(db *bolt.DB, issuer *accounts.Account, investor *accounts.Account, recipient *accounts.Account, investedAmount int, noOfYears int) (database.Order, error) {
	var newOrder database.Order
	assetName := AssetID("School_PuertoRico_1")
	// the reason why we have an int here is to avoid parsing
	// issues like dealing with random user strings "abc" could also be a valid input
	// if we decide to accept strings as our user input
	convRatio := float64(investedAmount/(noOfYears*12) + 1) // x usd = 1 PB token
	// the +1 is to offset the ratio to a whole number and make paybacks slightly less
	// which would mean the investors get paid ~months*1 more, which can be offset in another place
	paybackTokens := noOfYears * 12 // float to have granularity of sorts
	// the school would pay us back in USD tokens however, we use the conversion ratio of usd/234 to calculate payback period
	// assume the school pays 230 usd this year
	payBackPeriodLeft := float64(paybackTokens) - 200.0/convRatio
	log.Println("Payback Token Ratio for asset class is 1PB: ", convRatio, " USD tokens and payback period is ", payBackPeriodLeft) // +1 to round up

	// so now we create payBack and investor tokens for this asset class
	// the issuer is the platform itself, so people have to trust us (maybe give proofs for this?)
	// what is the guarantee that I don't issue as many PBTokens as I want? none?
	INVAssetName := AssetID("INVTokens_" + assetName) // ie. sha3(INVTokens_School_PuertoRico_1)[64:76]
	DEBAssetName := AssetID("DEBTokens_" + assetName) // ie. sha3(INVTokens_School_PuertoRico_1)[64:76]
	PBAssetName := AssetID("PBTokens_" + assetName)
	iAmt := utils.IntToString(investedAmount)
	dAmt := utils.IntToString(investedAmount*2)
	pbAmt := utils.IntToString(noOfYears * 12)
	log.Printf("Created asset names are %s and %s and amounts are %s and %s", INVAssetName, PBAssetName, iAmt, pbAmt)
	PBasset := issuer.CreateAsset(PBAssetName)
	INVasset := issuer.CreateAsset(INVAssetName)
	DEBasset := issuer.CreateAsset(DEBAssetName)

	// so I have the assets for this school created
	// now the investors need to trust me only for investedAmount of INVTokens

	txHash, err := investor.TrustAsset(INVasset, string(iAmt))
	if err != nil {
		return newOrder, err
	}
	log.Println("Investor trusted asset: ", INVasset.Code, " tx hash: ", txHash)

	// and the school needs to trust me only for paybackTokens amount of PB tokens
	txHash, err = recipient.TrustAsset(PBasset, string(pbAmt))
	if err != nil {
		return newOrder, err
	}
	log.Println("Recipient Trusted Payback asset: ", PBasset.Code, " tx hash: ", txHash)

	txHash, err = recipient.TrustAsset(DEBasset, string(dAmt)) // since debt = invested amount
	if err != nil {
		return newOrder, err
	}
	log.Println("Recipient Trusted Debt asset: ", DEBasset.Code, " tx hash: ", txHash)

	// so now the investor has his tokens, send paybackTokens to the school
	// log.Println("Sending PBasset for: ", string(pbAmt))
	// err = issuer.SendAsset(PBAssetName, recipient.PublicKey, pbAmt)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// so now the investors trust the issued asset and the recipients trust the issued asset.
	// send the assets over to the investor
	log.Println("Sending INVasset: ", INVAssetName, "for: ", string(iAmt))
	_, _, err = issuer.SendAsset(INVAssetName, investor.PublicKey, iAmt)
	if err != nil {
		return newOrder, err
	}
	log.Println("Sending DEBasset: ", DEBAssetName, "for: ", string(iAmt))
	_, _, err = issuer.SendAsset(DEBAssetName, recipient.PublicKey, iAmt) // same amount as debt
	if err != nil {
		return newOrder, err
	}

	newOrder, err = database.NewOrder(db, "16x20 panels", investedAmount, "Puerto Rico", investedAmount, "This is test data", INVAssetName, DEBAssetName, PBAssetName)
	if err != nil {
		log.Println("Error creating a new order. Quitting!")
		log.Fatal(err)
	}
	log.Println("Created new order: ", newOrder)

	// return asset names since we need to track this stuff
	return newOrder, nil
}
