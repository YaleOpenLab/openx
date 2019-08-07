package stablecoin

import (
	"github.com/pkg/errors"
	"log"
	"time"

	tickers "github.com/Varunram/essentials/crypto/exchangetickers"
	xlm "github.com/Varunram/essentials/crypto/xlm"
	assets "github.com/Varunram/essentials/crypto/xlm/assets"
	// utils "github.com/Varunram/essentials/utils"
)

// Exchange exchanges xlm for STABLEUSD
func Exchange(recipientPK string, recipientSeed string, convAmount float64) error {

	if Mainnet {
		return errors.New("Exchange in mainent needs to be done through dex")
	}

	if !xlm.AccountExists(recipientPK) {
		return errors.New("Account does not exist, quitting!")
	}

	// check whether user has enough xlm to pay. If not, quit
	balance, err := xlm.GetNativeBalance(recipientPK)
	if err != nil {
		return errors.Wrap(err, "couldn't get native balance from api")
	}

	if balance < convAmount {
		return errors.New("balance is less than amount requested")
	}

	var trustLimit float64
	trustLimit, err = xlm.GetAssetTrustLimit(recipientPK, StablecoinCode)
	if err != nil {
		// asset doesn't exist
		trustLimit = 0
	}

	if trustLimit < convAmount && trustLimit != 0 {
		return errors.Wrap(err, "trust limit doesn't warrant investment, please contact platform admin")
	}

	hash, err := assets.TrustAsset(StablecoinCode, StablecoinPublicKey, StablecoinTrustLimit, recipientSeed)
	if err != nil {
		return errors.Wrap(err, "couldn't trust asset")
	}
	log.Println("tx hash for trusting stableUSD: ", hash)
	// now send coins across and see if our tracker detects it
	log.Println(StablecoinPublicKey, convAmount, recipientSeed, "Exchange XLM for stablecoin")
	_, hash, err = xlm.SendXLM(StablecoinPublicKey, convAmount, recipientSeed, "Exchange XLM for stablecoin")
	if err != nil {
		return errors.Wrap(err, "couldn't send xlm")
	}
	log.Println("tx hash for sent xlm: ", hash, "pubkey: ", recipientPK)
	return nil
}

// OfferExchange offers to exchange user's xlm balance for stableusd if the user does not have enough
// stableUSD to complete the payment
func OfferExchange(publicKey string, seed string, invAmount float64) error {

	if Mainnet {
		return errors.New("Exchange offers in mainnet need to be done through dex")
	}

	balance, err := xlm.GetAssetBalance(publicKey, StablecoinCode)
	if err != nil {
		// the user does not have a balance in STABLEUSD
		balance = 0
	}

	if balance < invAmount {
		log.Println("Offering xlm to stableusd exchange to investor")
		// user's stablecoin balance is less than the amount he wishes to invest, get stablecoin
		// equal to the amount he wishes to exchange
		diff := invAmount - balance + 10 // the extra 1 is to cover for fees
		// checking whether the user has enough xlm balance to cover for the exchange is done by Exchange()
		xlmBalance, err := xlm.GetNativeBalance(publicKey)
		if err != nil {
			return errors.Wrap(err, "couldn't get native balance from api")
		}

		totalUSD := tickers.ExchangeXLMforUSD(xlmBalance) // amount in stablecoin that the user would receive for diff

		if totalUSD < diff {
			return errors.New("User does not have enough funds to complete this transaction")
		}

		// now we need to exchange XLM equal to diff in stablecoin
		exchangeRate := tickers.ExchangeXLMforUSD(1)
		// 1 xlm can fetch exchangeRate USD, how much xlm does diff USD need?
		amountToExchange := diff / exchangeRate
		log.Println(diff, exchangeRate, amountToExchange)
		err = Exchange(publicKey, seed, amountToExchange)
		if err != nil {
			return errors.Wrap(err, "Unable to exchange XLM for USD and automate payment. Please get more STABLEUSD to fulfil the payment")
		}
		time.Sleep(10 * time.Second) // 5 seconds for issuing stalbeusd to the person who's requested for it
	} else {
		log.Println("User has sufficient stablecoin balance, not exchanging xlm for usd")
	}

	return nil
}
