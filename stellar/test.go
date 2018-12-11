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
	INVAssetName, PBAssetName, DEBAssetName, err := assets.SetupAssets(&issuer, &investor, &recipient, opts.InvAmount, opts.RecYears)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Investor: %s, Payback: %s,  Debt: %s", INVAssetName, PBAssetName, DEBAssetName)
	// At this point, we have created the assets according to params passed and
	// now we would want to simulate the situation where people pay the party
	// in question. This woul;d broadly involve the given steps:
	// 1. Pay the ISSUER in USD tokens (since we assume USDtoken is equivalent to USD in question)
	// the extended question though is if we omit the PBToken directly, but that would mean
	// we have no record of the agreed period on the blockchain. Kind of weird, but since there
	// is no state, we have  to adopt this.
	// so the user transfers x USD tokens back to the issuer and then once the transaction
	// is confirmed (which should be relatively fast in Stellar due to its quorum)
	// we call the balance API to see whether we've trnasferred the assets. If our balance
	// in the token increases by x amount, we convert the amount to payback tokens
	// and then transfer the payback tokens to the recipient OR we could transfer the payback
	// tokens from client to contract, but that would mean they could sign arbitrary amounts
	// since they hold the seed. Hence we should transfer the payback tokens to the recipient
	// to show progress in ownership.
	// In short,
	// the recipient pays in USD tokens and receives pyaback tokens as return

	// this checks for balance, would come into use later on to check if we sent
	// the right amomunt of money to the user

	// this gets the balance of all the coins belonging to a specific account
	// balances, err := recipient.GetAllBalances()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	debtAssetBalance, err := recipient.GetAssetBalance(DEBAssetName)
	if err != nil {
		log.Fatal(err)
	}

	pbAssetBalance, err := recipient.GetAssetBalance(PBAssetName)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Debt balance: %s, Payback Balance: %s", debtAssetBalance, pbAssetBalance)

	// now we need to simulate a situation where the recipient pays back a certain
	// portion of the funds
	// again, onboarding is omitted here, since that's a bigger problem that we hopefully
	// can delegate to other parties
	// an alternate idea is that they  can buy s tellar and repay, if we choose to
	// take that route, we must modify some stuff and make things easier.
/*
	confHeight, txHash, err := recipient.SendAsset(DEBAssetName, issuer.PublicKey , "200") // send some coins from the issuer to the recipient
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Sending Debt Token to issuer", confHeight, txHash)
	// now we need to check the balance to ensure that the amount was actually paid
	// THis has to be done on the server side, so this might be a good point of
	// distinction between server and client runnign software.

	newDebtBalance, err := recipient.GetAssetBalance(DEBAssetName)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Old balance: %s, New Balance: %s", debtAssetBalance, newDebtBalance)
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
