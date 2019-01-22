package xlm

import (
	"github.com/stellar/go/build"
	"github.com/stellar/go/network"
)

func SetAuthImmutable(seed string) (int32, string, error) {
	//  Create with Auth immutable since we don't want the asset to be revocable
	passphrase := network.TestNetworkPassphrase
	tx, err := build.Transaction(
		build.SourceAccount{seed},
		build.AutoSequence{TestNetClient},
		build.Network{passphrase},
		build.MemoText{"Set Auth Immutable"},
		build.SetOptions(
			build.SetAuthImmutable(),
		),
	)

	if err != nil {
		return -1, "", err
	}

	return SendTx(seed, tx)
}

func FreezeAccount(seed string) (int32, string, error) {
	passphrase := network.TestNetworkPassphrase
	tx, err := build.Transaction(
		build.SourceAccount{seed},
		build.AutoSequence{TestNetClient},
		build.Network{passphrase},
		build.MemoText{"Freezing Account"},
		build.SetOptions(
			build.MasterWeight(0),
			build.SetThresholds(0, 0, 0),
		),
	)

	if err != nil {
		return -1, "", err
	}

	return SendTx(seed, tx)
}
