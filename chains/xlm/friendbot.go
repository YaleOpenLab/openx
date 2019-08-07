package xlm

import (
	"github.com/pkg/errors"
	"net/http"
)

// GetXLM makes an API call to the stellar friendbot, which gives 10000 testnet XLM
func GetXLM(PublicKey string) error {
	if Mainnet {
		return errors.New("no friendbot on mainnet, quitting")
	}
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + PublicKey)
	if err != nil || resp.Status != "200 OK" {
		return errors.New("API Request did not succeed")
	}
	return nil
}
