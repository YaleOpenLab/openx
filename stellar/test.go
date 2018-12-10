package main

import (
	"fmt"
	"log"
	"os"

	accounts "github.com/Varunram/smartPropertyMVP/stellar/accounts"
	assets "github.com/Varunram/smartPropertyMVP/stellar/assets"
	orders "github.com/Varunram/smartPropertyMVP/stellar/orders"
	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose   []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	InvAmount int    `short:"i" description:"Desired investment" required:"true"`
	RecYears  int    `short:"r" description:"Number of years the recipient wants to repay in" required:"true"`
}

func ValidateInputs() {
	if !(opts.RecYears == 3 || opts.RecYears == 5 || opts.RecYears == 7) {
		// right now payoff periods are limited, I guess they don't need to be,
		// but in this case jsut are
		log.Fatal(fmt.Errorf("Number of years not supported"))
	}
}

func main() {
	db, err := orders.OpenDB()
	if err != nil {
		log.Fatal(err)
		// this means that we couldn't open the dtabase and we need to do something else
	}
	defer db.Close()
	log.Fatal("DB works")
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

	// everyone should have coins to setup trustlines. Testnet doe not allow for sending
	// coins to account that have no balance - BUG?
	// anyways, stellar has a fat testnet wallet, so no worries

	err = issuer.GetCoins() // get coins for issuer
	if err != nil {
		log.Fatal(err)
	}

	// the problem with this is we generally accept donations in crypto and then
	// people have to trust this that we don't print stuff out of thin air
	// instead of using our own coin, we could use stronghold coin (stablecoin on Stellar)
	// Stellar also has an immediate DEX, but do we use it? ethical stuff while dealing with
	// funds remiain
	err = investor.GetCoins() // get coins for recipient
	if err != nil {
		log.Fatal(err)
	}

	err = recipient.GetCoins() // get coins for recipient
	if err != nil {
		log.Fatal(err)
	}

	// before setting up the assets, we need to refer to the orderbook in order to
	// get the list of available offers and funding things. For this purpose, we could
	// build a hash table / a simple dictionary, but I think investors in general
	// would like more info, so a simple map should be enough.
	// And this needs to be st ored in a database somewhere so that we don't lose this
	// data. Also need cryptographic proofs that this data is what it is, because
	// there is no concept of state in stellar. Is there a better way?
	err = assets.SetupAssets(&issuer, &investor, &recipient, opts.InvAmount, opts.RecYears)
	if err != nil {
		log.Fatal(err)
	}

	// log.Fatal(fmt.Errorf("All good"))

	// so now at this point, we assume that the investor wants to buy some solar shares
	// we need to validate the amount that he gives us, but we'll take care of that later.
	// assume that a single investor has invested 14000 USD for this project
	// so now issue the INV asset

	// this checks for balance, would come into use later on to check if we sent
	// the right amomunt of money to the user
	// err = issuer.Balance()
	// if err != nil {
	// 	log.Fatal(err)
	// }

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
