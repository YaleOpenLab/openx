// +build all travis

package xlm

import (
	"testing"
	"time"

	"github.com/stellar/go/build"
)

func TestXLM(t *testing.T) {
	seed, address, err := GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	err = GetXLM(address)
	if err != nil {
		t.Fatal(err)
	}
	err = GetXLM("blah")
	if err == nil {
		t.Fatal("Invalid Address, shouldn't work!")
	}
	// now we have coins, so we can send to another account
	// need to create an account through
	_, destPubKey, err := GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLMCreateAccount(destPubKey, "2", "wrong seed")
	// create the destiantion account by sendignb some coins to bootstrap
	if err == nil {
		t.Fatalf("Wrong seed, shouldn't work!")
	}
	_, _, err = SendXLMCreateAccount("wrong pubkey", "2", seed)
	// create the destiantion account by sendignb some coins to bootstrap
	if err == nil {
		t.Fatalf("Wrong pubkey, shouldn't work!")
	}
	_, _, err = SendXLMCreateAccount(destPubKey, "2", seed)
	// create the destiantion account by sendignb some coins to bootstrap
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLM(destPubKey, "1", seed, "")
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLM(destPubKey, "1", "wrong seed", "")
	if err == nil {
		t.Fatalf("Wrong seed, shouldn't work!")
	}
	_, _, err = SendXLM("wrong pubkey", "1", seed, "")
	if err == nil {
		t.Fatalf("Wrong pubkey, shouldn't work!")
	}
	if AccountExists("blah") {
		t.Fatalf("Dummy account exists, shouldn't!")
	}
	if !AccountExists(destPubKey) {
		t.Fatalf("Account which should exist doesn't, quitting!")
	}
	_, testPubKey, err := GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	err = RefillAccount(testPubKey, seed)
	if err != nil {
		t.Fatal(err)
	}
	err = RefillAccount("blah", seed)
	if err == nil {
		t.Fatal("Not catching wrong pubkey error, quitting!")
	}
	// don't test the reverse becuase apparently there's some problem in catching the next block
	// or something
	testseed, pk, err := GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	err = GetXLM(pk)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLM(destPubKey, "9998", testseed, "")
	if err != nil {
		t.Fatal(err)
	}
	err = RefillAccount(pk, seed)
	if err != nil {
		t.Fatal(err)
	}
	// testseed, pk ; seed, address
	_ = build.CreditAsset("BLAH1", pk)
	_, err = trustAsset("BLAH1", pk, "10", address, seed)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = sendAssetFromIssuer("BLAH1", address, "8", testseed, pk)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	_, err = GetAssetTrustLimit(address, "BLAH1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetAssetTrustLimit(address, "BLAH2")
	if err == nil {
		t.Fatalf("not able to catch invalid asset error")
	}
	_, _, err = SetAuthImmutable(seed)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = FreezeAccount(seed)
	if err != nil {
		t.Fatal(err)
	}
	oldTc := TestNetClient.URL
	TestNetClient.URL = "blah"
	_, _, err = SetAuthImmutable(seed)
	if err == nil {
		t.Fatalf("can send tx with an invalid client url, quitting!")
	}
	_, _, err = FreezeAccount(seed)
	if err == nil {
		t.Fatalf("can send tx with an invalid client url, quitting!")
	}
	TestNetClient.URL = oldTc
}

// BELOW FUNCTIONS BORROWED FROM ASSETS DUE TO IMPORT LIMITATIONS
func trustAsset(assetCode string, assetIssuer string, limit string, PublicKey string, Seed string) (string, error) {
	// TRUST is FROM PublicKey TO Seed
	trustTx, err := build.Transaction(
		build.SourceAccount{PublicKey},
		build.AutoSequence{SequenceProvider: TestNetClient},
		build.TestNetwork,
		build.Trust(assetCode, assetIssuer, build.Limit(limit)),
	)

	if err != nil {
		return "", err
	}

	_, txHash, err := SendTx(Seed, trustTx)
	return txHash, err
}

func sendAssetFromIssuer(assetName string, destination string, amount string,
	issuerSeed string, issuerPubkey string) (int32, string, error) {
	// this transaction is FROM issuer TO recipient
	paymentTx, err := build.Transaction(
		build.SourceAccount{issuerPubkey},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: TestNetClient},
		build.MemoText{"Sending Asset: " + assetName},
		build.Payment(
			build.Destination{AddressOrSeed: destination},
			build.CreditAmount{assetName, issuerPubkey, amount},
			// CreditAmount identifies the asset by asset Code and issuer pubkey
		),
	)

	if err != nil {
		return -1, "", err
	}
	return SendTx(issuerSeed, paymentTx)
}
