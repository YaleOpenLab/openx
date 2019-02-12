package stablecoin

import (
	"fmt"
	"log"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	oracle "github.com/YaleOpenLab/openx/oracle"
	utils "github.com/YaleOpenLab/openx/utils"
	xlm "github.com/YaleOpenLab/openx/xlm"
)

func Exchange(recipientPK string, recipientSeed string, convAmount string) error {

	if !xlm.AccountExists(recipientPK) {
		return fmt.Errorf("Account does not exist, quitting!")
	}

	// check whether user has enough xlm to pay. If not, quit
	balance, err := xlm.GetNativeBalance(recipientPK)
	if err != nil {
		log.Println(err)
		return err
	}

	if utils.StoF(balance) <= utils.StoF(convAmount) {
		return fmt.Errorf("insufficient balance")
	}

	hash, err := assets.TrustAsset(consts.Code, consts.StableCoinAddress, consts.StablecoinTrustLimit, recipientPK, recipientSeed)
	if err != nil {
		return err
	}
	log.Println("tx hash for trusting stableUSD: ", hash)
	// now send coins across and see if our tracker detects it
	_, hash, err = xlm.SendXLM(consts.StablecoinPublicKey, convAmount, recipientSeed, "Exchange XLM for stablecoin")
	if err != nil {
		return err
	}
	log.Println("tx hash for sent xlm: ", hash, "pubkey: ", recipientPK)
	return nil
}

// offer to exchange user's xlm balance for stableusd if the user does not have enough
// stableUSD to compelte the payment
func OfferExchange(publicKey string, seed string, invAmount string) error {

	log.Println("OFFERING EXCHANGE TO INVESTOR")
	balance, err := xlm.GetAssetBalance(publicKey, consts.Code)
	if err != nil {
		// the user does not have a balance in STABLEUSD
		balance = "0"
	}

	balF := utils.StoF(balance)
	invF := utils.StoF(invAmount)
	if balF < invF {
		// user's stablecoin balance is less than the amount he wishes to invest, get stablecoin
		// equal to the amount he wishes to exchange
		diff := invF - balF + 1 // the extra 1 is to cover for fees
		// checking whether the user has enough xlm balance to cover for the exchange is done by Exchange()
		xlmBalance, err := xlm.GetNativeBalance(publicKey)
		if err != nil {
			return err
		}

		totalUSD := oracle.ExchangeXLMforUSD(xlmBalance) // amount in stablecoin that the user would receive for diff

		if totalUSD < diff {
			return fmt.Errorf("User does not have enough funds to compelte this transaction")
		}

		// now we need to exchange XLM equal to diff in stablecoin
		exchangeRate := oracle.ExchangeXLMforUSD("1")
		// 1 xlm can fetch exchangeRate USD, how much xlm does diff USD need?
		amountToExchange := diff / exchangeRate

		err = Exchange(publicKey, seed, utils.FtoS(amountToExchange))
		if err != nil {
			return fmt.Errorf("Unable to exchange XLM for USD and automate payment. Please get more STABLEUSD to fulfil the payment")
		}
	}

	return nil
}
