// +build all algo

package algorand

import (
	"log"
	"testing"

	utils "github.com/Varunram/essentials/utils"
)

func TestAlgorandFunctions(t *testing.T) {

	var err error
	AlgodClient, err = InitAlgodClient()
	if err != nil {
		t.Fatalf("failed to make algod client\n")
	}

	KmdClient, err = InitKmdClient()
	if err != nil {
		t.Fatalf("failed to init kmd client")
	}

	nodeStatus, err := GetStatus(AlgodClient)
	if err != nil {
		t.Fatalf("error getting algod status\n")
	}

	log.Println("NODE STATUS: ", nodeStatus)
	// Fetch block information
	block, err := GetLatestBlock(nodeStatus)
	if err != nil {
		t.Fatalf("error getting last block\n")
	}
	log.Println("BLOCK: ", block.Hash)

	walletName := utils.GetRandomString(6)
	address, err := CreateNewWalletAndAddress(walletName, "x")
	if err != nil {
		t.Fatalf("error while creating a new wallet, exiting")
	}
	log.Println("New address generated: ", address)

	fromAddr := "YXU3MTTKV74UAGED6ROTHVVPEY5646WI3N5FLLQZWFV66AFKVQ5PMMYDZE" // the account that has funds
	txhash, err := SendAlgoToSelf("blah", "x", fromAddr, 100000)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Sent tx hash: ", txhash)

	txhash, err = SendAlgo("blah", "x", 100000, fromAddr, address)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Sent tx hash: ", txhash)
	txhash, err = GetAlgo("blah", "x")
	if err != nil {
		t.Fatal()
	}
	log.Println("Sent tx hash: ", txhash)
}
