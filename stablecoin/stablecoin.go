package stablecoin

// the idea of this stablecoin package is to issue a stablecoin on stellar testnet
// so that we can test the function of something similar on mainnet. The stablecoin
// provider should be stored in a different database because we will not be migrating
// this.

// The idea is to issue a single USD asset for every USD that we receive on our
// account, this should be automated and we must not have any kind of user interaction that is in
// place here. We also need a stablecoin Code, which we shall call as "STABLEUSD"
// for easy reference. Most functions would be similar to the one in assets.go,
// but need to be tailored to suit our requirements

// the USD asset defined here is what is issued by the speciifc bank. Ideally, we
// could accept a tx hash and check it as well, but since we can query balances,
// much easier to do it this way.
// or can be something like a stablecoin or asset
import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	consts "github.com/YaleOpenLab/openx/consts"
	oracle "github.com/YaleOpenLab/openx/oracle"
	scan "github.com/YaleOpenLab/openx/scan"
	utils "github.com/YaleOpenLab/openx/utils"
	xlm "github.com/YaleOpenLab/openx/xlm"
	assets "github.com/YaleOpenLab/openx/xlm/assets"
	wallet "github.com/YaleOpenLab/openx/xlm/wallet"
	"github.com/pkg/errors"
	horizon "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon/operations"
)

// InitStableCoin returns the platform structure and the seed
func InitStableCoin() error {
	var publicKey string
	var seed string
	// now we can be sure we have the directory, check for seed
	if _, err := os.Stat(consts.StableCoinSeedFile); !os.IsNotExist(err) {
		// the seed exists
		fmt.Println("ENTER YOUR PASSWORD TO DECRYPT THE STABLECOIN SEED FILE")
		password, err := scan.ScanRawPassword()
		if err != nil {
			return errors.Wrap(err, "couldn't scan raw password")
		}
		publicKey, seed, err = wallet.RetrieveSeed(consts.StableCoinSeedFile, password)
		// catch error here due to scope sharing
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Enter a password to encrypt your stablecoin's master seed. Please store this in a very safe place. This prompt will not ask to confirm your password")
		password, err := scan.ScanRawPassword()
		if err != nil {
			return err
		}
		publicKey, seed, err = wallet.NewSeed(consts.StableCoinSeedFile, password)
		if err != nil {
			return err
		}
		err = xlm.GetXLM(publicKey)
		if err != nil {
			return err
		}
	}
	// the user doesn't have seed, so create a new platform
	consts.StablecoinPublicKey = publicKey
	consts.StablecoinSeed = seed

	go ListenForPayments()
	return nil
}

func ListenForPayments() {
	client := xlm.TestNetClient
	// all payments
	opRequest := horizon.OperationRequest{ForAccount: consts.StableCoinAddress}

	ctx, _ := context.WithCancel(context.Background()) // cancel
	go func() {
		// Stop streaming after 60 seconds.
		log.Println("monitoring payments made towards address")
		time.Sleep(5 * time.Second) // refresh the thread every 5 seconds to check for payments
		// cancel() don't cancel the handler, let it run indefinitely
	}()

	printHandler := func(op operations.Operation) {
		log.Println("stablecoin operation: ", op)
		log.Println("PAGING TOKEN: ", op.PagingToken())
		log.Println("GETTYPE TOKEN: ", op.GetType())
		log.Println("GETID TOKEN: ", op.GetID())
		log.Println("GetTransactionHash TOKEN: ", op.GetTransactionHash())
		log.Println("IsTransactionSuccessful TOKEN: ", op.IsTransactionSuccessful())
		log.Println("IsTransactionSuccessful TOKEN: ", op)

		if op.IsTransactionSuccessful() {
			switch payment := op.(type) {
			case operations.Payment:
				log.Println("sending stablecoin to counterparty")
				log.Println("CHECK THIS OUT: ", payment.Asset.Type)
				if payment.Asset.Type == "native" { // native asset
					payee := payment.From
					amount := payment.Amount
					log.Printf("Received request for stablecoin from %s worth %s", payee, amount)
					xlmWorth := oracle.ExchangeXLMforUSD(amount)
					log.Println("The deposited amount is worth: ", xlmWorth)
					// now send the stableusd asset over to this guy
					_, hash, err := assets.SendAssetFromIssuer(consts.StablecoinCode, payee, utils.FtoS(xlmWorth), consts.StablecoinSeed, consts.StablecoinPublicKey)
					if err != nil {
						log.Println("Error while sending USD Assets back to payee: ", payee, err)
						//  don't skip here, there's technically nothing we can do
					}
					log.Println("Successful payment, hash: ", hash)
				}
			}
		}
	}

	err := client.StreamPayments(ctx, opRequest, printHandler)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
