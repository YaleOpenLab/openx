package stablecoin

import (
	"github.com/pkg/errors"
	"log"
	// "time"

	tickers "github.com/Varunram/essentials/crypto/exchangetickers"
	xlm "github.com/Varunram/essentials/crypto/xlm"
	assets "github.com/Varunram/essentials/crypto/xlm/assets"
)

// anchor implements stuff which is needed to interact with the anchor stablecoin

// GetAnchorUSDe exchanges XLM for AnchorUSD
func GetAnchorUSD(recpSeed string, amountUSD float64) (string, error) {

	txhash, err := assets.TrustAsset(AnchorUSDCode, AnchorUSDAddress, AnchorUSDTrustLimit, recpSeed)
	if err != nil {
		return txhash, errors.Wrap(err, "couldn't trust anchorUSD")
	}
	log.Println("tx hash for trusting stableUSD: ", txhash)
	// now send coins across and see if our tracker detects it
	// the given amount is in USD, we need to convert it into XLM since we're sending XLM
	exchangeRate, err := tickers.XLMUSD()
	if err != nil {
		return txhash, errors.Wrap(err, "error in fetching price from oracle")
	}
	amountXLM := exchangeRate * amountUSD

	log.Println("Exchanging: ", amountXLM, " XLM for anchorUSD")
	_, txhash, err = xlm.SendXLM(AnchorUSDAddress, amountXLM, recpSeed, "Exchange XLM for anchorUSD")
	if err != nil {
		return txhash, errors.Wrap(err, "couldn't send xlm")
	}
	log.Println("tx hash for sent xlm: ", txhash)
	return txhash, nil
}
