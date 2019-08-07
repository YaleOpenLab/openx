// +build all travis

package xlm

import (
	"testing"
	"time"

	"github.com/stellar/go/network"
	build "github.com/stellar/go/txnbuild"
)

func TestXLM(t *testing.T) {
	RefillAmount = 10
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
	_, _, err = SendXLMCreateAccount(destPubKey, 2, "wrong seed")
	// create the destinations account by sendignb some coins to bootstrap
	if err == nil {
		t.Fatalf("Wrong seed, shouldn't work!")
	}
	_, _, err = SendXLMCreateAccount("wrong pubkey", 2, seed)
	// create the destinations account by sendignb some coins to bootstrap
	if err == nil {
		t.Fatalf("Wrong pubkey, shouldn't work!")
	}
	_, _, err = SendXLMCreateAccount(destPubKey, 2, seed)
	// create the destinations account by sendignb some coins to bootstrap
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLM(destPubKey, 1, seed, "")
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLM(destPubKey, 1, "wrong seed", "")
	if err == nil {
		t.Fatalf("Wrong seed, shouldn't work!")
	}
	_, _, err = SendXLM("wrong pubkey", 1, seed, "")
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
	// don't test the reverse because apparently there's some problem in catching the next block
	// or something
	testseed, pk, err := GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	err = GetXLM(pk)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLM(destPubKey, 9998, testseed, "")
	if err != nil {
		t.Fatal(err)
	}
	err = RefillAccount(pk, seed)
	if err != nil {
		t.Fatal(err)
	}
	// testseed, pk ; seed, address
	_ = build.CreditAsset{"BLAH1", pk}
	_, err = trustAsset("BLAH1", pk, "10", seed)
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
	oldTc := TestNetClient.HorizonURL
	TestNetClient.HorizonURL = "blah"
	_, _, err = SetAuthImmutable(seed)
	if err == nil {
		t.Fatalf("can send tx with an invalid client url, quitting!")
	}
	_, _, err = FreezeAccount(seed)
	if err == nil {
		t.Fatalf("can send tx with an invalid client url, quitting!")
	}
	TestNetClient.HorizonURL = oldTc
}

func trustAsset(assetCode string, assetIssuer string, limit string, seed string) (string, error) {
	// TRUST is FROM Seed TO assetIssuer
	passphrase := network.TestNetworkPassphrase
	sourceAccount, mykp, err := ReturnSourceAccount(seed)
	if err != nil {
		return "", err
	}

	op := build.ChangeTrust{
		Line:  build.CreditAsset{assetCode, assetIssuer},
		Limit: limit,
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       passphrase,
	}

	_, txHash, err := SendTx(mykp, tx)
	if err != nil {
		return "", err
	}

	return txHash, err
}

func sendAssetFromIssuer(assetCode string, destination string, amount string,
	seed string, issuerPubkey string) (int32, string, error) {

	passphrase := network.TestNetworkPassphrase
	sourceAccount, mykp, err := ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", err
	}

	op := build.Payment{
		Destination: destination,
		Amount:      amount,
		Asset:       build.CreditAsset{assetCode, issuerPubkey},
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       passphrase,
	}

	return SendTx(mykp, tx)
}
