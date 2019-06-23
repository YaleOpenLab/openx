package xlm

import (
	"github.com/pkg/errors"
	"net/http"
)

// GetXLM makes an API call to friendbot, which gives us 10000 testnet XLM
func GetXLM(PublicKey string) error {
	// get some coins from the stellar robot for testing
	// gives only a constant amount of stellar, so no need to pass it a coin param
	// send stabelcoin from the platform instead of friendbot since it seems to be unstable

	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + PublicKey)
	if err != nil || resp.Status != "200 OK" {
		return errors.New("API Request did not succeed") // need this separately
	}

	// _, txhash, err := SendXLMCreateAccount(PublicKey, "6", consts.PlatformSeed)

	// log.Println("txhash for setting up account: ", txhash)
	// if err != nil {
	// return errors.New("Platform could not send XLM to account, quitting!")
	// }
	return nil
}
