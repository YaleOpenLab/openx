package stablecoin

// the idea of this stablecoin package is to issue a stablecoin on stellar testnet
// so that we can test the function of something similar on mainnet. The stablecoin
// provider should be stored in a different database because we will not be migrating
// this.

// The idea is to issue a single USD token for every USD t hat we receive on our
// account, this should be automated and we must not have any kind of user interaction that is in
// place here. We also need a stablecoin Code, which we shall call as "STABLEUSD"
// for easy reference. Most functions would be similar to the one in assets.go,
// but need to be tailored to suit our requirements

// the USD token defined here is what is issued by the speciifc bank. Ideally, we
// could accept a tx hash and check it as well, but since we can query balances,
// much easier to do it this way.
// or can be something like a stablecoin or token
import (
	"context"
	"fmt"
	"log"
	"os"

	assets "github.com/OpenFinancing/openfinancing/assets"
	consts "github.com/OpenFinancing/openfinancing/consts"
	oracle "github.com/OpenFinancing/openfinancing/oracle"
	scan "github.com/OpenFinancing/openfinancing/scan"
	utils "github.com/OpenFinancing/openfinancing/utils"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

var StableUSD build.Asset
var PublicKey string
var Seed string
var Code = "STABLEUSD"

// InitStableCoin returns the platform structure and the seed
func InitStableCoin() error {
	var publicKey string
	var seed string
	var err error
	// now we can be sure we have the directory, check for seed
	if _, err := os.Stat(consts.StableCoinSeedFile); !os.IsNotExist(err) {
		// the seed exists
		fmt.Println("ENTER YOUR PASSWORD TO DECRYPT THE STABLECOIN SEED FILE")
		password, err := scan.ScanRawPassword()
		if err != nil {
			log.Println(err)
			return err
		}
		publicKey, seed, err = wallet.RetrieveSeed(consts.StableCoinSeedFile, password)
	} else {
		fmt.Println("Enter a password to encrypt your stablecoin's master seed. Please store this in a very safe place. This prompt will not ask to confirm your password")
		password, err := scan.ScanRawPassword()
		if err != nil {
			return err
		}
		publicKey, seed, err = wallet.NewSeed(consts.StableCoinSeedFile, password)
		err = xlm.GetXLM(publicKey)
	}
	// the user doesn't have seed, so create a new platform
	if err != nil {
		return err
	}
	PublicKey = publicKey
	Seed = seed
	StableUSD.Issuer = PublicKey
	StableUSD.Code = "STABLEUSD"
	log.Printf("STABLECOIN PUBLICKEY IS: %s, SEED is: %s", PublicKey, Seed)
	go ListenForPayments()
	return err
}

// ListenForPayments listens for payments to the stablecoin account and once it
// gets the transaction hash from the rmeote API, calculates how much USD it owes
// for the amount deposited and then transfers the StableUSD asset to the payee
// Prices are retrieved from an oracle.
func ListenForPayments() {
	// the publicKey above has to be hardcoded as a constant because stellar's API wants it like so
	// stupid stuff, but we need to go ahead with it. In reality, this shouldn't
	// be much of a problem since we expect that the platform's seed will be
	// constant
	ctx := context.Background() // start in the background context
	cursor := horizon.Cursor("now")
	fmt.Println("Waiting for a payment...")
	err := xlm.TestNetClient.StreamPayments(ctx, consts.StableCoinAddress, &cursor, func(payment horizon.Payment) {
		/*
			Sample Response:
			Payment type payment
			Payment Paging Token 5424212982374401
			Payment From GC76MINOSNQUMDBNONBARFYFCCQA5QLNQSJOVANR5RNQVHQRB5B46B6I
			Payment To GBAACP6UUXZAB5ZAYAHWEYLNKORWB36WVBZBXWNPFXQTDY2AIQFM6D7Y
			Payment Asset Type native
			Payment Asset Code
			Payment Asset Issuer
			Payment Amount 10.0000000
			Payment Memo Type
			Payment Memo
		*/
		log.Println("Stablecoin payment to/from detected")
		log.Println("Payment type", payment.Type)
		log.Println("Payment Paging Token", payment.PagingToken)
		log.Println("Payment From", payment.From)
		log.Println("Payment To", payment.To)
		log.Println("Payment Asset Type", payment.AssetType)
		log.Println("Payment Asset Code", payment.AssetCode)
		log.Println("Payment Asset Issuer", payment.AssetIssuer)
		log.Println("Payment Amount", payment.Amount)
		log.Println("Payment Memo Type", payment.Memo.Type)
		log.Println("Payment Memo", payment.Memo.Value)
		if payment.Type == "payment" && payment.AssetType == "native" {
			// store the stuff that we want here
			payee := payment.From
			amount := payment.Amount
			log.Printf("Received request for stablecoin from %s worth %s", payee, amount)
			xlmWorth := oracle.ExchangeXLMforUSD(amount)
			log.Println("The deposited amount is worth: ", xlmWorth)
			// now send the stableusd asset over to this guy
			_, hash, err := assets.SendAssetFromIssuer(Code, payee, utils.FtoS(xlmWorth), Seed, PublicKey)
			if err != nil {
				log.Println("Error while sending USD Assets back to payee: ", payee)
				log.Println(err)
				//  don't skip here, there's technically nothing we can do
			}
			log.Println("Successful payment, hash: ", hash)
		}
	})

	if err != nil {
		// we shouldn't ideally fatal here, but do since we're testing out stuff
		log.Fatal(err)
	}
}
