package oracle

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"

	utils "github.com/YaleOpenLab/openx/utils"
)

// All tickers are requested from Binance due to it being the largest exchange by volume
var XLMUSDReq = "https://api.binance.com/api/v1/ticker/price?symbol=XLMUSDT"

// BinanceTickerResponse defines the ticker API response from Binanace
type BinanceTickerResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// GetRequest is a clone of the original GetRequest from rpc
func GetRequest(url string) ([]byte, error) {
	var dummy []byte
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("did not create new GET request", err)
		return dummy, err
	}
	req.Header.Set("Origin", "localhost")
	res, err := client.Do(req)
	if err != nil {
		log.Println("did not make request", err)
		return dummy, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// XLMUSD returns the XLMUSD ticker from Binance
func XLMUSD() (float64, error) {
	data, err := GetRequest(XLMUSDReq)
	if err != nil {
		log.Println("did not get response", err)
		return -1, errors.Wrap(err, "did not get response from Binance API")
	}

	var response BinanceTickerResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal response")
	}

	if response.Symbol != "XLMUSDT" {
		return -1, fmt.Errorf("ticker symbols don't match with API response")
	}
	// response.Price is in string, need to convert it to float
	price, err := utils.StoFWithCheck(response.Price)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return price, nil
}

// https://api.binance.com/api/v1/ticker/price?symbol=XLMUSDT
