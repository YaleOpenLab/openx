package main

import (
	"fmt"
	"log"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/algod/models"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/types"
)

// These constants represent the algod REST endpoint and the corresponding
// API token. You can retrieve these from the `algod.net` and `algod.token`
// files in the algod data directory.
var algodAddress = "http://localhost:49809"
var algodToken = "57724c8fd1146e26d8f9805734414c4374f3528fc1201796feb701a2358bdd55"
var Client algod.Client

const kmdAddress = "http://localhost:7833"
const kmdToken = "a91d47703ce61823872df82d072470d2ceb203f5c27cea1012dea3d8d7eacaf7"

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

func GetLatestBlock(status models.NodeStatus) (models.Block, error) {
	return Client.Block(status.LastRound)
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
	Client, err = algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		fmt.Printf("failed to make algod client: %s\n", err)
		return
	}

	// Print algod status
	nodeStatus, err := GetStatus(Client)
	if err != nil {
		fmt.Printf("error getting algod status: %s\n", err)
		return
	}

	// Fetch block information
	block, err := GetLatestBlock(nodeStatus)
	if err != nil {
		fmt.Printf("error getting last block: %s\n", err)
		return
	}

	log.Println("BLOCK: ", block)

	address, err := CreateNewWallet("blah", "x")
	if err != nil {
		log.Println("error while creating a new walelt, exiting")
		return
	}

	log.Println("New address generated: ", address)
}
