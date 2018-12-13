package assets

// so we create two entities here = issuer and recipient
// the issuer is us, wewant to give out solar panel contracts and we issue assets
// there should be n investors, who will be investing in this contract and we can
// define a preset amount in dollars based on the value of the bond.
// the recipient which is the school has to first agree on a time bound period
// of payback, absed on which another asset (payBack Asset) is issued and this
// will automatically update based on the payments amde by the school.

// for eg, take ABC School to which we assign a bond of 14000 USD with 10 investors
// the first step would be to issue an asset which is wort 1:1 with the dollar
// (or use dollars itself, have to see how this works out) So now we have 14000
// INV_schoolname_capacity tokens which are possessed by investors as proof that they invested in this specific asset.
// now that we have 14000 INV tokens created, we need to create  the payback tokens
// for the case of this example, lets assume that there are 3 options:
// A. 3 YEARS = 36 mo
// B. 5 YEARS = 60 mo
// C. 7 YEARS = 84 mo
// now we need to create a peg for INV token based on the years the school chooses.
// Lets assume that the school chooses 5 years. We need to peg the PB (payback) token
// like 14000 USD : 60 PB which means 233.333 USD = 1 PB token
// now here, we round this up for ease of granularity, so this would be 234 USD a month.
// The schoole could now choose to payback 1 PB a month, which would mean it gains ownership in exactly
// 5 years, it could also pay faster, which would mean that they own the asset earlier
// years / months reamining is simply the balance in payback tokens (50.42 PB for eg)
// and we can use this to display users how much time they have remaining to own the asset

// In net, we have to create 2 assets:
// 1. Investor Token UNIQUE to each bond
// 2. Payback Token UNIQUE to each bond
// lets leave validation for later since that requires state validation stuff

import (
	"log"

	"github.com/boltdb/bolt"
	accounts "github.com/YaleOpenLab/smartPropertyMVP/stellar/accounts"
	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
)

func AssetID(inputString string) string {
	// so the assetID right now is a hash of the asset name, concatenated investor public keys and nonces
	x := utils.SHA3hash(inputString)
	return "YOL" + x[64:73] // max length of an asset in stellar is 12
	// log.Fatal(fmt.Errorf("All good"))
	// return nil
}

func CalculatePayback(balance string, noOfMonths string) error {
	// the idea is that we should be able ot pass an assetId to this function
	// and it must calculate how much time we have left for payback. For this example
	// until twe do the db stuff, lets pass a few params (although this could be done
	// separately as well).
	return nil
}

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
