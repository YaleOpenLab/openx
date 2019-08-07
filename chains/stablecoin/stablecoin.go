package stablecoin

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tickers "github.com/Varunram/essentials/crypto/exchangetickers"
	xlm "github.com/Varunram/essentials/crypto/xlm"
	assets "github.com/Varunram/essentials/crypto/xlm/assets"
	wallet "github.com/Varunram/essentials/crypto/xlm/wallet"
	scan "github.com/Varunram/essentials/scan"
	utils "github.com/Varunram/essentials/utils"
	"github.com/pkg/errors"
	horizon "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon/operations"
)

// Package stablecoin implements a stablecoin with code STABLEUSD built on Stellar.

// InitStableCoin starts the stablecoin daemon
func InitStableCoin(mainnet bool) (string, string, error) {
	if Mainnet {
		return "", "", errors.New("Stablecoin in mainnet defaults to AnchorUSD")
	}
	var publicKey string
	var seed string
	// now we can be sure we have the directory, check for seed
	if _, err := os.Stat(StableCoinSeedFile); !os.IsNotExist(err) {
		// the seed exists
		fmt.Println("ENTER YOUR PASSWORD TO DECRYPT THE STABLECOIN SEED FILE")
		password, err := scan.ScanRawPassword()
		if err != nil {
			return "", "", errors.Wrap(err, "couldn't scan raw password")
		}
		publicKey, seed, err = wallet.RetrieveSeed(StableCoinSeedFile, password)
		if err != nil {
			return "", "", err
		}
	} else {
		// stablecoin doesn't exist yet
		fmt.Println("Enter a password to encrypt your stablecoin's master seed. Please store this in a very safe place. This prompt will not ask to confirm your password")
		password, err := scan.ScanRawPassword()
		if err != nil {
			return "", "", err
		}
		publicKey, seed, err = wallet.NewSeedStore(StableCoinSeedFile, password)
		if err != nil {
			return "", "", err
		}
		if !mainnet {
			err = xlm.GetXLM(publicKey)
			if err != nil {
				return "", "", err
			}
		}
	}

	if Mainnet {
		if !xlm.AccountExists(publicKey) {
			return "", "", errors.New("please refill your account: " + publicKey + " with funds to setup a stellar account. Your seed is: " + seed)
		} else {
			err := xlm.GetXLM(publicKey)
			if err != nil {
				return "", "", err
			}
		}
	}

	// the user doesn't have seed, so create a new platform
	StablecoinPublicKey = publicKey
	StablecoinSeed = seed

	go ListenForPayments()
	return StablecoinPublicKey, StablecoinSeed, nil
}

// ListenForPayments listens to all payments to/from the stablecoin address
func ListenForPayments() {
	client := xlm.TestNetClient
	opRequest := horizon.OperationRequest{ForAccount: StablecoinPublicKey}

	ctx, cancel := context.WithCancel(context.Background()) // cancel
	defer cancel()
	go func() {
		// Stop streaming after 60 seconds.
		log.Println("monitoring payments made towards address")
		time.Sleep(5 * time.Second) // refresh the thread every 5 seconds to check for payments
		// cancel() don't cancel the handler, let it run indefinitely
	}()

	printHandler := func(op operations.Operation) {
		/*
			log.Println("stablecoin operation: ", op)
			log.Println("PAGING TOKEN: ", op.PagingToken())
			log.Println("GETTYPE TOKEN: ", op.GetType())
			log.Println("GETID TOKEN: ", op.GetID())
			log.Println("GetTransactionHash TOKEN: ", op.GetTransactionHash())
			log.Println("IsTransactionSuccessful TOKEN: ", op.IsTransactionSuccessful())
			log.Println("IsTransactionSuccessful TOKEN: ", op)
		*/
		if op.IsTransactionSuccessful() {
			switch payment := op.(type) {
			case operations.Payment:
				log.Println("sending stablecoin to counterparty")
				// log.Println("CHECK THIS OUT: ", payment.Asset.Type)
				if payment.Asset.Type == "native" { // native asset
					payee := payment.From
					amount, _ := utils.ToFloat(payment.Amount)
					log.Println("Received request for stablecoin from", payee, ",worth", amount)
					xlmWorth := tickers.ExchangeXLMforUSD(amount)
					log.Println("The deposited amount is worth: ", xlmWorth)
					// now send the stableusd asset over to this guy
					_, hash, err := assets.SendAssetFromIssuer(StablecoinCode, payee, xlmWorth, StablecoinSeed, StablecoinPublicKey)
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
