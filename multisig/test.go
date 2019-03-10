package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stellar/go/network"
	"github.com/stellar/go/build"
	clients "github.com/stellar/go/clients/horizon"
)

/*
InflationDest("GCT7S5BA6ZC7SV7GGEMEYJTWOBYTBOA7SC4JEYP7IAEDG7HQNIWKRJ4G"),
	SetAuthRequired(),
	SetAuthRevocable(),
	SetAuthImmutable(),
	ClearAuthRequired(),
	ClearAuthRevocable(),
	ClearAuthImmutable(),
	MasterWeight(1),
	SetThresholds(2, 3, 4),
	HomeDomain("stellar.org"),
	AddSigner("GC6DDGPXVWXD5V6XOWJ7VUTDYI7VKPV2RAJWBVBHR47OPV5NASUNHTJW", 5),
*/
var TestNetClient = &clients.Client{
	// URL: "http://35.192.122.229:8080",
	URL:  "https://horizon-testnet.stellar.org",
	HTTP: http.DefaultClient,
}

func Multisig2of2() {
	// don't check if the account exists or not, hopefully it does
	party1Pubkey := "GBQXTMPZ6NP7ND4ZGV7N7X6J6RBW2SVH7HB4I3I6L3SRTM6GJ4WDDNAN"
	party1Seed := "SA55EOJKRERBF5YQG5DXYYGIR5NKUZ2ENQECBYSWNKSA4WGRIHKMMLIW"

	party3Pubkey := "GB67ZCS5CJWRIQCU6VYGMT2VPQUWOJDYXZULI2W7FXXHYED4HTMXSJAE"
	party3Seed := "SCV2OGWZ6ZUCI2XBNTJMEEHB7RKVO7CAZM6UHMLOVCH57EGKTWUM27IU"

	party2Pubkey := "GB77NPDNGWCGKWYKFLYBNZSEN5OVHR2IVMYTV2HJPK7SAQ4MBWARR6KO"
	party2Seed := "SBI2Y7Q7GUL7PZCSDDDWRKKB6UPMBTHV2ESK2JS3NHMGK5LQN7MVP3AA"

	log.Println(party1Seed, party1Pubkey, party3Seed, party3Pubkey, party2Seed, party2Pubkey)

	memo := "testsign"
	amount := "1"

	tx, err := build.Transaction(
		build.SourceAccount{party2Pubkey},
		build.AutoSequence{TestNetClient},
		build.Network{network.TestNetworkPassphrase},
		build.MemoText{memo},
		build.Payment(
			build.Destination{party2Pubkey},
			build.NativeAmount{amount},
		),
		build.SetOptions(
			build.MasterWeight(1),
			build.AddSigner(party1Pubkey, 1),
			build.SetThresholds(2, 2, 2),
		),
	)

	if err != nil {
		log.Println("error while constructing tx", err)
		return
	}

	txe, err := tx.Sign(party2Seed) // sign using party 2's seed
	if err != nil {
		log.Println("second party couldn't sign tx", err)
		return
	}

	txeB64, err := txe.Base64()
	if err != nil {
		log.Println(err)
		return
	}
	// And finally, send it off to Stellar
	resp, err := TestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		log.Println("error while submitting tx", err)
		return
	}

	fmt.Printf("Propagated Transaction: %s, sequence: %d\n", resp.Hash, resp.Ledger)
	return
}

func TestMutlisig2of2() {
	party1Pubkey := "GBQXTMPZ6NP7ND4ZGV7N7X6J6RBW2SVH7HB4I3I6L3SRTM6GJ4WDDNAN"
	party1Seed := "SA55EOJKRERBF5YQG5DXYYGIR5NKUZ2ENQECBYSWNKSA4WGRIHKMMLIW"

	party3Pubkey := "GB67ZCS5CJWRIQCU6VYGMT2VPQUWOJDYXZULI2W7FXXHYED4HTMXSJAE"
	party3Seed := "SCV2OGWZ6ZUCI2XBNTJMEEHB7RKVO7CAZM6UHMLOVCH57EGKTWUM27IU"

	party2Pubkey := "GB77NPDNGWCGKWYKFLYBNZSEN5OVHR2IVMYTV2HJPK7SAQ4MBWARR6KO"
	party2Seed := "SBI2Y7Q7GUL7PZCSDDDWRKKB6UPMBTHV2ESK2JS3NHMGK5LQN7MVP3AA"

	memo := "works"

	log.Println(party1Seed, party1Pubkey, party3Seed, party3Pubkey, party2Seed, party2Pubkey)

	tx, err := build.Transaction(
		build.SourceAccount{party2Pubkey},
		build.AutoSequence{TestNetClient},
		build.Network{network.TestNetworkPassphrase},
		build.MemoText{memo},
		build.Payment(
			build.Destination{party2Pubkey},
			build.NativeAmount{"1"},
		),
	)

	txe, err := tx.Sign(party2Seed, party1Seed) // sign using party 2's seed
	if err != nil {
		log.Println("second party couldn't sign tx", err)
		return
	}

	txeB64, err := txe.Base64()
	if err != nil {
		log.Println(err)
		return
	}

	// now we have the base64. WE NEED TO SIGN IT USING THE OTHER SEED
	// And finally, send it off to Stellar
	resp, err := TestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		log.Println("error while submitting tx", err)
		return
	}

	fmt.Printf("Propagated Transaction: %s, sequence: %d\n", resp.Hash, resp.Ledger)
	return
}
