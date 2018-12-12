package main

import (
	"fmt"
	"log"
	"os"

	accounts "github.com/Varunram/smartPropertyMVP/stellar/accounts"
	assets "github.com/Varunram/smartPropertyMVP/stellar/assets"
	orders "github.com/Varunram/smartPropertyMVP/stellar/orders"
	utils "github.com/Varunram/smartPropertyMVP/stellar/utils"
	server "github.com/Varunram/smartPropertyMVP/stellar/server"
	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose   []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	InvAmount int    `short:"i" description:"Desired investment" required:"true"`
	RecYears  int    `short:"r" description:"Number of years the recipient wants to repay in. Can be 3, 5 or 7 years." required:"true"`
}

func ValidateInputs() {
	if !(opts.RecYears == 3 || opts.RecYears == 5 || opts.RecYears == 7) {
		// right now payoff periods are limited, I guess they don't need to be,
		// but in this case jsut are
		log.Fatal(fmt.Errorf("Number of years not supported"))
	}
}

func main() {
	var err error
	server.SetupServer() // this must be towards the end
	db, err := orders.OpenDB()
	if err != nil {
		log.Fatal(err)
		// this means that we couldn't open the database and we need to do something else
	}
	defer db.Close()
	_, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("InvAmount: %d USD, RecYears: %d years, Verbose: %t", opts.InvAmount, opts.RecYears, opts.Verbose)
	ValidateInputs()

	issuer := accounts.SetupAccount()
	investor := accounts.SetupAccount()
	recipient := accounts.SetupAccount()

	log.Println("The issuer's public key and private key are: ", issuer.PublicKey, " ", issuer.Seed)
	log.Println("The investor's public key and private key are: ", investor.PublicKey, " ", investor.Seed)
	log.Println("The recipient's public key and private key are: ", recipient.PublicKey, " ", recipient.Seed)

	// everyone should have coins to setup trustlines.
	// anyways, stellar has a fat testnet wallet, so no worry that this might
	// get depleted

	err = issuer.GetCoins() // get coins for issuer
	if err != nil {
		log.Fatal(err)
	}

	err = issuer.SetupAccount(recipient.PublicKey, "10")
	if err != nil {
		log.Println("Recipient Account not setup")
		log.Fatal(err)
	}

	err = issuer.SetupAccount(investor.PublicKey, "10")
	if err != nil {
		log.Println("Investor Account not setup")
		log.Fatal(err)
	}

	// the problem with this is we generally accept donations in crypto and then
	// people have to trust this that we don't print stuff out of thin air
	// instead of using our own coin, we could use stronghold coin (stablecoin on Stellar)
	// Stellar also has an immediate DEX, but do we use it? ethical stuff while dealing with
	// funds remiain
	// before setting up the assets, we need to refer to the orderbook in order to
	// get the list of available offers and funding things. For this purpose, we could
	// build a hash table / a simple dictionary, but I think investors in general
	// would like more info, so a simple map should be enough.
	// And this needs to be stored in a database somewhere so that we don't lose this
	// data. Also need cryptographic proofs that this data is what it is, because
	// there is no concept of state in stellar. Is there a better way?
	a, err := assets.SetupAsset(db, &issuer, &investor, &recipient, opts.InvAmount, opts.RecYears)
	if err != nil {
		log.Fatal(err)
	}
	// In short, the recipient pays in DEBtokens and receives PBtokens in return

	// this checks for balance, would come into use later on to check if we sent
	// the right amomunt of money to the user
	// balances, err := recipient.GetAllBalances()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// now we need to simulate a situation where the recipient pays back a certain
	// portion of the funds
	// onboarding is omitted here, that's a bigger problem that we hopefully
	// can delegate to other parties like Neighborly
	// an alternate idea is that they can buy stellar and repay, if we choose to
	// take that route, we must use a coin on stellar as an anchor to receive this token.
	// in this way, we need to check native balance and then use the anchor
	// right now don't do that, but should do in future to solicit donations from
	// the community, who would be generally dealing in XLM (and not DEBtoken)

	// another idea is that you could speculate on DEBtoken by having a market
	// for it, that would reuqire to relax the flags a bit. Right now, we don't
	// use an authorization flag, but we should since we don't want alternate markets
	// to develop. If we do, don't set the flag
	paybackAmount := "210"
	err = recipient.Payback(db, a.Index, a.DEBAssetCode, issuer.PublicKey, paybackAmount)
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}
	// after this ,we must update the steuff on the server side and send a payback token
	// to let the user know that he has paid x amoutn of money.
	// this however, would be the money paid / money that has to be paid per month
	// in total, this should be payBackPeriod * 12

	paybackAmountF := utils.StringToFloat(paybackAmount)
	refundS := utils.FloatToString(paybackAmountF / accounts.PriceOracleInFloat())
	// weird conversion stuff, but have to since the amount should be in a string
	blockHeight, txHash, err := issuer.SendAsset(a.PBAssetCode, recipient.PublicKey, refundS)
	if err != nil {
		log.Println("Error while sending a payback token, notify help immediately")
		log.Fatal(err)
	}
	log.Println("Sent payback token to recipient", blockHeight, txHash)
	tOrder, err := orders.RetrieveOrder(a.Index, db)
	if err != nil {
		log.Println("Error retrieving from db")
		log.Fatal(err)
	}
	log.Println("Test whether this was updated: ", tOrder)

	debtAssetBalance, err := recipient.GetAssetBalance(a.DEBAssetCode)
	if err != nil {
		log.Fatal(err)
	}

	pbAssetBalance, err := recipient.GetAssetBalance(a.PBAssetCode)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Debt balance: %s, Payback Balance: %s", debtAssetBalance, pbAssetBalance)

	/*
		confHeight, txHash, err := issuer.SendCoins(recipient.PublicKey, "3.34") // send some coins from the issuer to the recipient
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Confirmation height is: ", confHeight, " and txHash is: ", txHash)

		asset := issuer.CreateAsset(assetName) // create the asset that we want

		trustLimit := "100" // trust only 100 barrels of oil from Petro
		err = recipient.TrustAsset(asset, trustLimit)
		if err != nil {
			log.Println("Trust limit is in the wrong format")
			log.Fatal(err)
		}

		err = issuer.SendAsset(assetName, recipient.PublicKey, "3.34")
		if err != nil {
			log.Fatal(err)
		}
	*/
}
