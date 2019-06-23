package xlm

import (
	"github.com/pkg/errors"
	//"log"

	consts "github.com/YaleOpenLab/openx/consts"
	oracle "github.com/YaleOpenLab/openx/oracle"
	utils "github.com/YaleOpenLab/openx/utils"
	"github.com/stellar/go/network"
	build "github.com/stellar/go/txnbuild"
)

// package dex contains functions for interfacing with the stellar dex

// NewBuyOrder creates a new buy order on the stellar dex
func NewBuyOrder(seed string, assetName string, issuer string,
	amount string, price string) (int32, string, error) {

	sourceAccount, mykp, err := ReturnSourceAccount(seed)
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
	return SendTx(mykp, tx)
}

// NewSellOrder creates a new sell order on the stellar dex
func NewSellOrder(seed string, assetName string, issuer string, amount string,
	price string) (int32, string, error) {

	sourceAccount, mykp, err := ReturnSourceAccount(seed)
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

	return SendTx(mykp, tx)
}

// DexStableCoinBuy gets the price from an oracle and places an order on the DEX to buy AnchorUSD
func DexStableCoinBuy(seed string, amount string) (int32, string, error) {
	assetName := "USD"
	issuer := consts.AnchorUSDAddress
	price, err := oracle.BinanceTicker()
	if err != nil {
		return -1, "", errors.New("could not fetch price form binance, quitting")
	}
	price = price * 1.02 // a small premium to get the order fulfilled immediately
	return NewBuyOrder(seed, assetName, issuer, amount, utils.FtoS(price))
}

// DexStableCoinBuy gets the price from an oracle and places an order on the DEX to sell AnchorUSD
func DexStableCoinSell(seed string, amount string) (int32, string, error) {
	assetName := "USD"
	issuer := consts.AnchorUSDAddress
	price, err := oracle.BinanceTicker()
	if err != nil {
		return -1, "", errors.New("could not fetch price form binance, quitting")
	}
	price = price * 1.02 // a small premium to get the order fulfilled immediately
	return NewSellOrder(seed, assetName, issuer, amount, utils.FtoS(price))
}
