package oracle

//  oracle defines oracle related functions
// right now, most are just placehlders, but in the future, they should call
// remote APIs and act as oracles
import (
	utils "github.com/YaleOpenLab/openx/utils"
)

// MonthlyBill returns the power tariffs and any data that we need to certify
// that is in the real world. Right now, this is hardcoded since we need to come up
// with a construct to get the price data in a reliable way - this could be a website
// were poeple erport this or certified authorities can timestamp this on chain
// or similar. Web scraping government websites might work, but that seems too
// overkill for what we're doing now.
func MonthlyBill() string {
	// right now, community consensus look like the price of electricity is
	// $0.2 per kWH in Puerto Rico, so hardcoding that here.
	priceOfElectricity := 0.2
	// since solar is free, they just need to pay this and then in some x time (chosen
	// when the order is created / confirmed on the school side), they
	// can own the panel.
	// the average energy consumption in puerto rico seems to be 5,657 kWh or about
	// 471 kWH per household. lets take 600 accounting for a 20% error margin.
	averageConsumption := float64(600)
	avgString := utils.FtoS(priceOfElectricity * averageConsumption)
	return avgString
}

// MonthlyBillInFloat does the same thing as PriceOracle, but returns the data
// as a float for use in appropriate places
func MonthlyBillInFloat() float64 {
	priceOfElectricity := 0.2
	averageConsumption := float64(600)
	return priceOfElectricity * averageConsumption
}

// ExchangeXLMforUSD retrieves the current price of XLM/USD and then returns the USD amount
// that the XLM deposited is worth and takes a percentage premium that emulates
// how real world exchanges would behave. This fee is 0.01% for now
func ExchangeXLMforUSD(amount string) float64 {
	// defines the rate for 1 usd = x XLM. Currently hardcoded to 10
	amountF := utils.StoF(amount)
	// exchangeRate := 0.1 // hardcode for now, can query cmc apis later
	exchangeRate := 10000000.0 // rig the exchange rate so that we can test some stuff
	return amountF * exchangeRate
}
