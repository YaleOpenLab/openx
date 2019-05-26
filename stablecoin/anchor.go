package stablecoin

import (
	"github.com/pkg/errors"
	"log"
	// "time"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	oracle "github.com/YaleOpenLab/openx/oracle"
	utils "github.com/YaleOpenLab/openx/utils"
	xlm "github.com/YaleOpenLab/openx/xlm"
)

func GetAnchorUSD(recpSeed string, amountUSDs string) (string, error) {
	txhash, err := assets.TrustAsset(consts.AnchorUSDCode, consts.AnchorUSDAddress, consts.AnchorUSDTrustLimit, recpSeed)
	// txhash, err := assets.TrustAsset(consts.Code, consts.StableCoinAddress, consts.StablecoinTrustLimit, recpSeed)
	if err != nil {
		return txhash, errors.Wrap(err, "couldn't trust anchorUSD")
	}
	log.Println("tx hash for trusting stableUSD: ", txhash)
	// now send coins across and see if our tracker detects it
	// the given amount is in USD, we need to convert it into XLM since we're sending XLM
	amountUSD, err := utils.StoFWithCheck(amountUSDs)
	if err != nil {
		return txhash, err
	}

	exchangeRate, err := oracle.XLMUSD()
	if err != nil {
		return txhash, errors.Wrap(err, "error in fetching price from oracle")
	}
	amountXLM := exchangeRate * amountUSD

	log.Println("Exchanging: ", amountXLM, " XLM for anchorUSD")
	_, txhash, err = xlm.SendXLM(consts.AnchorUSDAddress, utils.FtoS(amountXLM), recpSeed, "Exchange XLM for anchorUSD")
	if err != nil {
		return txhash, errors.Wrap(err, "couldn't send xlm")
	}
	log.Println("tx hash for sent xlm: ", txhash)
	return txhash, nil
}
