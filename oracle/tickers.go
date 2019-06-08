package oracle

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"

	utils "github.com/YaleOpenLab/openx/utils"
)

// we take the three largest (no wash trading) markets for XLM USD and return their weighted average
// to arrive at the price for XLM-USD. This price is indicative and not final since there will be latency
// involved between price display and trade finality.
var BinanceReq = "https://api.binance.com/api/v1/ticker/price?symbol=XLMUSDT"
var CoinbaseReq = "https://api.pro.coinbase.com/products/XLM-USD/ticker"
var KrakenReq = "https://api.kraken.com/0/public/Ticker?pair=XLMUSD"

var BinanceVol = "https://api.binance.com/api/v1/ticker/24hr?symbol=XLMUSDT"

// BinanceTickerResponse defines the ticker API response from Binanace
type BinanceTickerResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type BinanceVolumeResponse struct {
	// there are other fields as well, but we ignore them for now
	Symbol string `json:"symbol"`
	Volume string `json:"volume"`
}

type CoinbaseTickerResponse struct {
	TradeId int    `json:"trade_id"`
	Price   string `json:"price"`
	Volume  string `json:"volume"`
}

type KrakenTickerResponse struct {
	Error  []string `json:"error"`
	Result struct {
		XXLMZUSD struct {
			// there's some additional info here but we don't require that
			C []string // c = last trade closed array(<price>, <lot volume>),
			V []string // volume array(<today>, <last 24 hours>)
		}
	}
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
		return -1, errors.New("ticker symbols don't match with API response")
	}
	// response.Price is in string, need to convert it to float
	price, err := utils.StoFWithCheck(response.Price)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return price, nil
}

func BinanceVolume() (float64, error) {
	data, err := GetRequest(BinanceVol)
	if err != nil {
		log.Println("did not get response", err)
		return -1, errors.Wrap(err, "did not get response from Binance API")
	}

	var response BinanceVolumeResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal response")
	}

	if response.Symbol != "XLMUSDT" {
		return -1, errors.New("ticker symbols don't match with API response")
	}

	volume, err := utils.StoFWithCheck(response.Volume)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return volume, nil // volume is in xlm and not usd
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

	// response.Price is in string, need to convert it to float
	price, err := utils.StoFWithCheck(response.Price)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return price, nil
}

func CoinbaseVolume() (float64, error) {
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

	// response.Price is in string, need to convert it to float
	volume, err := utils.StoFWithCheck(response.Volume)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return volume, nil
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

func KrakenVolume() (float64, error) {
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
	volume, err := utils.StoFWithCheck(response.Result.XXLMZUSD.V[1]) // we want volume over the last 24 hours
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return volume, nil
}

// XLMUSD returns the XLMUSD ticker from Binance
func XLMUSD() (float64, error) {
	binanceVolume, err := CoinbaseVolume()
	if err != nil {
		return -1, err
	}

	cbVolume, err := CoinbaseVolume()
	if err != nil {
		return -1, err
	}

	krakenVolume, err := KrakenVolume()
	if err != nil {
		return -1, err
	}

	binanceTicker, err := CoinbaseTicker()
	if err != nil {
		return -1, err
	}

	cbTicker, err := CoinbaseTicker()
	if err != nil {
		return -1, err
	}

	krakenTicker, err := KrakenTicker()
	if err != nil {
		return -1, err
	}

	netVolume := binanceVolume + cbVolume + krakenVolume

	// return weighted average of all the prices
	return binanceTicker*(binanceVolume/netVolume) + cbTicker*(cbVolume/netVolume) + krakenTicker*(krakenVolume/netVolume), nil
}
