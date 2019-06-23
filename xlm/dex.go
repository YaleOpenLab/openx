package xlm

import (
	"github.com/pkg/errors"
	//	"log"

	"github.com/stellar/go/network"
	build "github.com/stellar/go/txnbuild"
)

// package dex contains functions for interfacing with the stellar dex

// NewBuyOrder creates a new buy order on the stellar dex
func NewBuyOrder(seed string, assetName string, destination string,
	amount string, price string) (int32, string, error) {

	sourceAccount, mykp, err := ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not load client details, quitting")
	}

	buyOffer := build.ManageBuyOffer{
		Selling: build.NativeAsset{},
		Buying:  build.CreditAsset{assetName, destination},
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

	return SendTx(mykp, tx)
}

// NewSellOrder creates a new sell order on the stellar dex
func NewSellOrder(seed string, assetName string, destination string,
	amount string, price string) (int32, string, error) {

	sourceAccount, mykp, err := ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not load client details, quitting")
	}

	sellOffer := build.ManageBuyOffer{
		Selling: build.CreditAsset{assetName, destination},
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
