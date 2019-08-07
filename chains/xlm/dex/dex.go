package dex

import (
	"github.com/pkg/errors"
	//"log"

	ticker "github.com/Varunram/essentials/crypto/exchangetickers"
	stablecoin "github.com/Varunram/essentials/crypto/stablecoin"
	xlm "github.com/Varunram/essentials/crypto/xlm"
	utils "github.com/Varunram/essentials/utils"
	"github.com/stellar/go/network"
	build "github.com/stellar/go/txnbuild"
)

// package dex contains functions for interfacing with the stellar dex

// NewBuyOrder creates a new buy order on the stellar dex
func NewBuyOrder(seed string, assetName string, issuer string,
	amount string, price string) (int32, string, error) {

	sourceAccount, mykp, err := xlm.ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not load client details, quitting")
	}

	buyOffer := build.ManageBuyOffer{
		Selling: build.NativeAsset{},
		Buying:  build.CreditAsset{assetName, issuer},
		Amount:  amount,
		Price:   price,
		OfferID: 0,
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&buyOffer},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       network.TestNetworkPassphrase,
	}

	// once the offer is completed, we need to send a follow up tx to send funds to the requested address
	return xlm.SendTx(mykp, tx)
}

// NewSellOrder creates a new sell order on the stellar dex
func NewSellOrder(seed string, assetName string, issuer string, amount string,
	price string) (int32, string, error) {

	sourceAccount, mykp, err := xlm.ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not load client details, quitting")
	}

	sellOffer := build.ManageBuyOffer{
		Selling: build.CreditAsset{assetName, issuer},
		Buying:  build.NativeAsset{},
		Amount:  amount,
		Price:   price,
		OfferID: 0,
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&sellOffer},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       network.TestNetworkPassphrase,
	}

	return xlm.SendTx(mykp, tx)
}

// DexStableCoinBuy gets the price from an oracle and places an order on the DEX to buy AnchorUSD
func DexStableCoinBuy(seed string, amount string) (int32, string, error) {
	assetName := "USD"
	issuer := stablecoin.AnchorUSDAddress
	price, err := ticker.BinanceTicker()
	if err != nil {
		return -1, "", errors.New("could not fetch price form binance, quitting")
	}
	price = price * 1.02 // a small premium to get the order fulfilled immediately
	ftss, _ := utils.ToString(price)
	return NewBuyOrder(seed, assetName, issuer, amount, ftss)
}

// DexStableCoinSell places a sell order for STABLEUSD on the Stellar dex
func DexStableCoinSell(seed string, amount string) (int32, string, error) {
	assetName := "USD"
	issuer := stablecoin.AnchorUSDAddress
	price, err := ticker.BinanceTicker()
	if err != nil {
		return -1, "", errors.New("could not fetch price form binance, quitting")
	}
	price = price * 1.02 // a small premium to get the order fulfilled immediately
	ftss, _ := utils.ToString(price)
	return NewSellOrder(seed, assetName, issuer, amount, ftss)
}
