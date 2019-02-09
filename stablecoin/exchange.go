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
		return err
	}

	if utils.StoF(balance) != utils.StoF(convAmount) {
		return fmt.Errorf("insufficient balance")
	}

	hash, err := assets.TrustAsset(Code, consts.StableCoinAddress, consts.StablecoinTrustLimit, recipientPK, recipientSeed)
	if err != nil {
		return err
	}
	log.Println("tx hash for trusting stableUSD: ", hash)
	// now send coins across and see if our tracker detects it
	_, hash, err = xlm.SendXLM(PublicKey, convAmount, recipientSeed, "Sending XLM to bootstrap")
	if err != nil {
		return err
	}
	log.Println("tx hash for sent xlm: ", hash, "pubkey: ", recipientPK)
	return nil
}

// offer to exchange user's xlm balance for stableusd if the user does not have enough
// stableUSD to compelte the payment
func OfferExchange(publicKey string, seed string, invAmount string) error {

	balance, err := xlm.GetAssetBalance(publicKey, Code)
	if err != nil {
		return err
	}

	balF := utils.StoF(balance)
	invF := utils.StoF(invAmount)
	if balF < invF {
		// user's stablecoin balance is less than the amount he wishes to invest, get stablecoin
		// equal to the amount he wishes to exchange
		diff := invF - balF
		// checking whether the user has enough xlm balance to cover for the exchange is done by Exchange()
		amount := oracle.ExchangeXLMforUSD(utils.FtoS(diff))
		err := Exchange(publicKey, seed, utils.FtoS(amount))
		if err != nil {
			return fmt.Errorf("Unable to exchange XLM for USD and automate payment. Please get more STABLEUSD to fulfil the payment")
		}
	}

	return nil
}
