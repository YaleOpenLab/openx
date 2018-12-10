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
	"encoding/hex"
	"golang.org/x/crypto/sha3"
	"log"
	"strconv"

	accounts "github.com/Varunram/smartPropertyMVP/stellar/accounts"
)

func AssetID(inputString string) string {
	// so the assetID right now is a hash of the asset name, concatenated investor public keys and nonces
	x := SHA3hash(inputString)
	log.Println("LGHTR", len(x), x[64:80])
	return "YOL" + x[64:73] // max length of an asset in stellar is 12
	// log.Fatal(fmt.Errorf("All good"))
	// return nil
}

func SHA3hash(inputString string) string {
	byteString := sha3.Sum512([]byte(inputString))
	hexString := hex.EncodeToString(byteString[:])
	// so now we have a SHA3hash that we can use to assign unique ids to our assets
	return hexString
}

func SetupAssets(issuer *accounts.Account, investor *accounts.Account, recipient *accounts.Account, investedAmount int, noOfYears int) error {
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
	log.Printf("Payback Token Ratio for asset class is 1PB: %d USD tokens and payback period is %d", convRatio, payBackPeriodLeft) // +1 to round up

	// so now we create payBack and investor tokens for this asset class
	// the issuer is the platform itself, so people have to trust us (maybe give proofs for this?)
	// what is the guarantee that I don't issue as many PBTokens as I want? none?
	INVAssetName := AssetID("INVTokens_" + assetName) // ie. sha3(INVTokens_School_PuertoRico_1)[64:76]
	PBAssetName := AssetID("PBTokens_" + assetName)
	iAmt := strconv.Itoa(investedAmount)
	pbAmt := strconv.Itoa(noOfYears * 12)
	log.Printf("Created asset names are %s and %s and amounts are %s and %s", INVAssetName, PBAssetName, iAmt, pbAmt)
	PBasset := issuer.CreateAsset(PBAssetName)
	INVasset := issuer.CreateAsset(INVAssetName)

	// so I have the assets for this school created
	// now the investors need to trust me only for investedAmount of INVTokens

	err := investor.TrustAsset(INVasset, string(iAmt))
	if err != nil {
		log.Fatal(err)
	}

	// and the school needs to trust me only for paybackTokens amount of PB tokens
	err = recipient.TrustAsset(PBasset, string(pbAmt))
	if err != nil {
		log.Fatal(err)
	}

	// so now the invesotr has his tokens, send paybackTokens to the school
	log.Println("Sending PBasset for: ", string(pbAmt))
	err = issuer.SendAsset(PBAssetName, recipient.PublicKey, pbAmt)
	if err != nil {
		log.Fatal(err)
	}

	// so now the investors trust the issued asset and the recipients trust the issued asset.
	// send the assets over to the investor
	log.Println("Sending INVasset for: ", string(iAmt))
	err = issuer.SendAsset(INVAssetName, investor.PublicKey, iAmt)
	if err != nil {
		log.Fatal(err)
	}
	// the assets are with the concerned parties
	return nil
}
