package tickers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
)

// package tickers implements handlers for getting price from cryptocurrencyu markets

// we take the three largest (no wash trading) markets for XLM USD and return their weighted average
// to arrive at the price for XLM-USD. This price is indicative and not final since there will be latency
// involved between price display and trade finality.

// BinanceReq is the binance ticker from the API
var BinanceReq = "https://api.binance.com/api/v1/ticker/price?symbol=XLMUSDT"

// CoinbaseReq is the coinbase ticker from the API
var CoinbaseReq = "https://api.pro.coinbase.com/products/XLM-USD/ticker"

// KrakenReq is the kraken ticker from the API
var KrakenReq = "https://api.kraken.com/0/public/Ticker?pair=XLMUSD"

// BinanceVol is the binance ticker from the API
var BinanceVol = "https://api.binance.com/api/v1/ticker/24hr?symbol=XLMUSDT"

// BinanceTickerResponse defines the ticker API response from Binanace
type BinanceTickerResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// BinanceVolumeResponse defines the structure of binance's volume endpoint response
type BinanceVolumeResponse struct {
	// there are other fields as well, but we ignore them for now
	Symbol string `json:"symbol"`
	Volume string `json:"volume"`
}

// CoinbaseTickerResponse defines the structure of coinbase's ticker endpoitt response
type CoinbaseTickerResponse struct {
	TradeId int    `json:"trade_id"`
	Price   string `json:"price"`
	Volume  string `json:"volume"`
}

// KrakenTickerResponse defines the structure of kraken's ticker response
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

// BinanceTicker gets price data from Binance
func BinanceTicker() (float64, error) {
	data, err := erpc.GetRequest(BinanceReq)
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
	price, err := utils.ToFloat(response.Price)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return price, nil
}

// BinanceVolume gets volume data from Binance
func BinanceVolume() (float64, error) {
	data, err := erpc.GetRequest(BinanceVol)
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

	volume, err := utils.ToFloat(response.Volume)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return volume, nil // volume is in xlm and not usd
}

// CoinbaseTicker gets ticker data from coinbase
func CoinbaseTicker() (float64, error) {
	data, err := erpc.GetRequest(CoinbaseReq)
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
	price, err := utils.ToFloat(response.Price)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return price, nil
}

// CoinbaseVolume gets volume data from coinbase
func CoinbaseVolume() (float64, error) {
	data, err := erpc.GetRequest(CoinbaseReq)
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
	volume, err := utils.ToFloat(response.Volume)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return volume, nil
}

// KrakenTicker gets ticker data from kraken
func KrakenTicker() (float64, error) {
	data, err := erpc.GetRequest(KrakenReq)
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
	price, err := utils.ToFloat(response.Result.XXLMZUSD.C[0])
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return price, nil
}

// KrakenVolume gets volume data from kraken
func KrakenVolume() (float64, error) {
	data, err := erpc.GetRequest(KrakenReq)
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
	volume, err := utils.ToFloat(response.Result.XXLMZUSD.V[1]) // we want volume over the last 24 hours
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return volume, nil
}

// XLMUSD returns aggregated XLMUSD ticker data
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

// ExchangeXLMforUSD retrieves the current price of XLM/USD and then returns the USD amount
func ExchangeXLMforUSD(amount float64) float64 {
	// defines the rate for 1 usd = x XLM. Currently hardcoded to 10
	// exchangeRate := 0.1 // hardcode for now, can query cmc apis later
	exchangeRate := 10000000.0 // rig the exchange rate so that we can test some stuff
	return amount * exchangeRate
}
