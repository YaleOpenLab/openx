package xlm

import (
	"github.com/pkg/errors"
	//	"log"

	wallet "github.com/YaleOpenLab/openx/xlm/wallet"
	horizon "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	build "github.com/stellar/go/txnbuild"
)

func NewBuyOrder(encryptedSeed []byte, seedpwd string, assetName string,
	destination string, amount string, price string) (int32, string, error) {
	seed, err := wallet.DecryptSeed(encryptedSeed, seedpwd)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not decrypt seed, quitting")
	}

	mykp, err := keypair.Parse(seed)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not parse keypair, quitting")
	}

	client := horizon.DefaultTestNetClient
	ar := horizon.AccountRequest{AccountID: mykp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
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

func NewSellOrder(encryptedSeed []byte, seedpwd string, assetName string,
	destination string, amount string, price string) (int32, string, error) {
	seed, err := wallet.DecryptSeed(encryptedSeed, seedpwd)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not decrypt seed, quitting")
	}

	mykp, err := keypair.Parse(seed)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not parse keypair, quitting")
	}

	client := horizon.DefaultTestNetClient
	ar := horizon.AccountRequest{AccountID: mykp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
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
