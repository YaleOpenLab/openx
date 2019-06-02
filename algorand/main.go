package main

import (
	"fmt"
	"log"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/algod/models"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

// These constants represent the algod REST endpoint and the corresponding
// API token. You can retrieve these from the `algod.net` and `algod.token`
// files in the algod data directory.
const algodAddress = "http://localhost:50435"
const algodToken = "df6740f7618f699b0417f764b6447fa7e690f9514c73cd60184314ae16141030"
var Client algod.Client

const kmdAddress = "http://localhost:51976"
const kmdToken = "755071c9616f4ebac31512e4db7993dc056f12790d94d634e978a66dfc44ce9b"

var KmdClient kmd.Client

func GetStatus(Client algod.Client) (models.NodeStatus, error) {
	var status models.NodeStatus
	status, err := Client.Status()
	if err != nil {
		fmt.Printf("error getting algod status: %s\n", err)
		return status, err
	}

	return status, nil
}

func InitClient() (algod.Client, error) {
	var err error
	Client, err = algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		fmt.Printf("failed to make algod client: %s\n", err)
		return Client, nil
	}

	return Client, nil
}

func InitKmdClient() (kmd.Client, error) {
	return kmd.MakeClient(kmdAddress, kmdToken)
}

func GetLatestBlock(status models.NodeStatus) (models.Block, error) {
	return Client.Block(status.LastRound)
}

func GetBlock(blockNumber uint64) (models.Block, error) {
	return Client.Block(blockNumber)
}

func CreateNewWalletHandle(walletID string, password string) (string, error) {
	// Get a wallet handle. The wallet handle is used for things like signing transactions
	// and creating accounts. Wallet handles do expire, but they can be renewed
	initResponse, err := KmdClient.InitWalletHandle(walletID, password)
	if err != nil {
		fmt.Printf("Error initializing wallet handle: %s\n", err)
		return "", err
	}

	walletHandleToken := initResponse.WalletHandleToken
	return walletHandleToken, nil
}

func GenerateAddress(walletHandleToken string) (string, error) {
	// Generate a new address from the wallet handle
	genResponse, err := KmdClient.GenerateKey(walletHandleToken)
	if err != nil {
		fmt.Printf("Error generating key: %s\n", err)
		return "", err
	}
	fmt.Printf("Generated address %s\n", genResponse.Address)
	return genResponse.Address, nil
}

func SignTransaction(walletName string, password string, amount uint64) {
	// Get the list of wallets
	listResponse, err := KmdClient.ListWallets()
	if err != nil {
		fmt.Printf("error listing wallets: %s\n", err)
		return
	}
	log.Println("response: ", listResponse)
	// Find our wallet name in the list
	var ourWalletId string
	fmt.Printf("Got %d wallet(s):\n", len(listResponse.Wallets))
	for _, wallet := range listResponse.Wallets {
		fmt.Printf("ID: %s\tName: %s\n", wallet.ID, wallet.Name)
		if wallet.Name == walletName {
			fmt.Printf("found wallet '%s' with ID: %s\n", wallet.Name, wallet.ID)
			ourWalletId = wallet.ID
		}
	}

	exampleWalletHandleToken, err := CreateNewWalletHandle(ourWalletId, password)
	if err != nil {
		fmt.Printf("Error initializing wallet handle: %s\n", err)
		return
	}

	fromAddr := "YXU3MTTKV74UAGED6ROTHVVPEY5646WI3N5FLLQZWFV66AFKVQ5PMMYDZE"

	// Generate a new address from the wallet handle
	toAddr, err := GenerateAddress(exampleWalletHandleToken)
	if err != nil {
		fmt.Printf("Error generating key: %s\n", err)
		return
	}
	fmt.Printf("Generated address 2 %s\n", toAddr)

	// Get the suggested transaction parameters
	txParams, err := Client.SuggestedParams()
	if err != nil {
		fmt.Printf("error getting suggested tx params: %s\n", err)
		return
	}

	// Make transaction
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	fee := txParams.Fee
	lastRound := txParams.LastRound
	// (from, to string, fee, amount, firstRound, lastRound uint64, note []byte,
	// closeRemainderTo, genesisID string, genesisHash []byte) (types.Transaction, error)
	tx, err := transaction.MakePaymentTxn(fromAddr, toAddr, fee, amount, lastRound-50, lastRound+50, nil, "", genID, genHash)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
		return
	}

	// Sign the transaction
	signResponse, err := KmdClient.SignTransaction(exampleWalletHandleToken, password, tx)
	if err != nil {
		fmt.Printf("Failed to sign transaction with kmd: %s\n", err)
		return
	}

	fmt.Printf("kmd made signed transaction with bytes: %x\n", signResponse.SignedTransaction)

	// Broadcast the transaction to the network
	// Note that this transaction will get rejected because the accounts do not have any tokens
	sendResponse, err := Client.SendRawTransaction(signResponse.SignedTransaction)
	if err != nil {
		fmt.Printf("failed to send transaction: %s\n", err)
		return
	}

	fmt.Printf("Transaction ID: %s\n", sendResponse.TxID)
}

func CreateNewWallet(name string, password string) (string, error) {
	var err error
	KmdClient, err = kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		fmt.Printf("failed to make kmd client: %s\n", err)
		return "", err
	}
	fmt.Println("Made a kmd client")

	// Create the example wallet, if it doesn't already exist
	createWalletResponse, err := KmdClient.CreateWallet(name, password, kmd.DefaultWalletDriver, types.MasterDerivationKey{})
	if err != nil {
		fmt.Printf("error creating wallet: %s\n", err)
		return "", err
	}

	// We need the wallet ID in order to get a wallet handle, so we can add accounts
	walletID := createWalletResponse.Wallet.ID
	fmt.Printf("Created wallet '%s' with ID: %s\n", createWalletResponse.Wallet.Name, walletID)

	// Get a wallet handle. The wallet handle is used for things like signing transactions
	// and creating accounts. Wallet handles do expire, but they can be renewed
	walletHandleToken, err := CreateNewWalletHandle(walletID, password)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Generate a new address from the wallet handle
	address, err := GenerateAddress(walletHandleToken)
	if err != nil {
		return "", nil
	}
	return address, nil
}

func main() {
	// Create an algod client
	var err error
	Client, err = InitClient()
	if err != nil {
		fmt.Printf("failed to make algod client: %s\n", err)
		return
	}

	KmdClient, err = InitKmdClient()
	if err != nil {
		log.Println(err)
		return
	}

	/*
	// Print algod status
	nodeStatus, err := GetStatus(Client)
	if err != nil {
		fmt.Printf("error getting algod status: %s\n", err)
		return
	}

	log.Println("NODE STATUS: ", nodeStatus)
	// Fetch block information
	block, err := GetLatestBlock(nodeStatus)
	if err != nil {
		fmt.Printf("error getting last block: %s\n", err)
		return
	}

	log.Println("BLOCK: ", block.Hash)

	address, err := CreateNewWallet("blah", "x")
	if err != nil {
		log.Println("error while creating a new wallet, exiting")
		return
	}

	log.Println("New address generated: ", address)
	*/
	SignTransaction("blah", "x", 100000)
}
