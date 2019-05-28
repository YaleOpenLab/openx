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
var BinanceReq = "https://api.binance.com/api/v1/ticker/price?symbol=XLMUSDT"
var CoinbaseReq = "https://api.coinbase.com/v2/prices/XLM-USD/sell"
var KrakenReq = "https://api.kraken.com/0/public/Ticker?pair=XLMUSD"

// BinanceTickerResponse defines the ticker API response from Binanace
type BinanceTickerResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type CoinbaseTickerResponse struct {
	Data struct {
		Base     string `json:"base"`
		Currency string `json:"currency"`
		Amount   string `json:""amount`
	} `json:"data"`
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

func BinanceTicker() (float64, error) {
	data, err := GetRequest(BinanceReq)
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

func CoinbaseTicker() (float64, error) {
	data, err := GetRequest(CoinbaseReq)
	if err != nil {
		log.Println("did not get response", err)
		return -1, errors.Wrap(err, "did not get response from Coinbase API")
	}

	var response CoinbaseTickerResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal response")
	}

	if response.Data.Base != "XLM" {
		return -1, fmt.Errorf("ticker symbols don't match with API response")
	}
	// response.Price is in string, need to convert it to float
	price, err := utils.StoFWithCheck(response.Data.Amount)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return price, nil
}

type KrakenTickerResponse struct {
	Error  []string `json:"error"`
	Result struct {
		XXLMZUSD struct {
			// there's some additional info here but we don't require that
			C []string // c = last trade closed array(<price>, <lot volume>),
		}
	}
}

func KrakenTicker() (float64, error) {
	data, err := GetRequest(KrakenReq)
	if err != nil {
		log.Println("did not get response", err)
		return -1, errors.Wrap(err, "did not get response from Kraken API")
	}

	var response KrakenTickerResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal response")
	}

	// response.Price is in string, need to convert it to float
	price, err := utils.StoFWithCheck(response.Result.XXLMZUSD.C[0])
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return price, nil
}

// XLMUSD returns the XLMUSD ticker from Binance
func XLMUSD() (float64, error) {
	return BinanceTicker()
	// we could pre assign weightages here or read that dynamically be analysing volume data from each.
	// but  is that needed? we might as well show users all three, tell them these are the prices and that
	// actual prices might vary or something like that
}

// https://api.binance.com/api/v1/ticker/price?symbol=XLMUSDT
