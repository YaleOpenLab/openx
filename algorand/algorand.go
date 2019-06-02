package algorand

import (
	"fmt"
	"github.com/pkg/errors"
	"log"

	consts "github.com/YaleOpenLab/openx/consts"
	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/algod/models"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

// These constants represent the algod REST endpoint and the corresponding
// API token. You can retrieve these from the `algod.net` and `algod.token`
// files in the algod data directory.
var AlgodClient algod.Client
var KmdClient kmd.Client

func GetStatus(Client algod.Client) (models.NodeStatus, error) {
	var status models.NodeStatus
	status, err := AlgodClient.Status()
	if err != nil {
		fmt.Printf("error getting algod status: %s\n", err)
		return status, err
	}

	return status, nil
}

func InitClient() (algod.Client, error) {
	var err error
	AlgodClient, err = algod.MakeClient(consts.AlgodAddress, consts.AlgodToken)
	if err != nil {
		fmt.Printf("failed to make algod client: %s\n", err)
		return AlgodClient, nil
	}

	return AlgodClient, nil
}

func InitKmdClient() (kmd.Client, error) {
	return kmd.MakeClient(consts.KmdAddress, consts.KmdToken)
}

func GetLatestBlock(status models.NodeStatus) (models.Block, error) {
	return AlgodClient.Block(status.LastRound)
}

func GetBlock(blockNumber uint64) (models.Block, error) {
	return AlgodClient.Block(blockNumber)
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

func SendTransactionToSelf(walletName string, password string, fromAddr string, amount uint64) (string, error) {
	// Get the list of wallets
	listResponse, err := KmdClient.ListWallets()
	if err != nil {
		return "", errors.Wrap(err, "error listing wallets")
	}

	// Find our walletID in the list
	var ourWalletId string
	fmt.Printf("Got %d wallet(s):\n", len(listResponse.Wallets))
	for _, wallet := range listResponse.Wallets {
		fmt.Printf("ID: %s\tName: %s\n", wallet.ID, wallet.Name)
		if wallet.Name == walletName {
			fmt.Printf("found wallet '%s' with ID: %s\n", wallet.Name, wallet.ID)
			ourWalletId = wallet.ID
		}
	}

	walletHandleToken, err := CreateNewWalletHandle(ourWalletId, password)
	if err != nil {
		return "", errors.Wrap(err, "error initializing wallet handle")
	}

	// Generate a new address from the wallet handle
	toAddr, err := GenerateAddress(walletHandleToken)
	if err != nil {
		return "", errors.Wrap(err, "error generating key")
	}

	log.Println("Generated address 2: ", toAddr)

	// Get the suggested transaction parameters
	txParams, err := AlgodClient.SuggestedParams()
	if err != nil {
		return "", errors.Wrap(err, "error getting suggested tx params")
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
		return "", errors.Wrap(err, "error creating transaction")
	}

	// Sign the transaction
	signResponse, err := KmdClient.SignTransaction(walletHandleToken, password, tx)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign transaction with kmd")
	}

	// Broadcast the transaction to the network
	// Note that this transaction will get rejected because the accounts do not have any tokens
	sendResponse, err := AlgodClient.SendRawTransaction(signResponse.SignedTransaction)
	if err != nil {
		return "", errors.Wrap(err, "failed to send transaction")
	}

	log.Println("Transaction ID: ", sendResponse.TxID)
	return sendResponse.TxID, nil
}

func SendTransaction(walletName string, password string, amount uint64, fromAddr string, toAddr string) (string, error) {
	// Get the list of wallets
	listResponse, err := KmdClient.ListWallets()
	if err != nil {
		return "", errors.Wrap(err, "error listing wallets")
	}

	// Find our walletID in the list
	var ourWalletId string
	for _, wallet := range listResponse.Wallets {
		if wallet.Name == walletName {
			fmt.Printf("found wallet '%s' with ID: %s\n", wallet.Name, wallet.ID)
			ourWalletId = wallet.ID
		}
	}

	// get a wallet handle token to sign the transaction with
	walletHandleToken, err := CreateNewWalletHandle(ourWalletId, password)
	if err != nil {
		return "", errors.Wrap(err, "error initializing wallet handle")
	}

	// Get suggested parameters
	txParams, err := AlgodClient.SuggestedParams()
	if err != nil {
		return "", errors.Wrap(err, "error getting suggested tx params")
	}

	// params for estimation
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	fee := txParams.Fee
	lastRound := txParams.LastRound
	// (from, to string, fee, amount, firstRound, lastRound uint64, note []byte,
	// closeRemainderTo, genesisID string, genesisHash []byte) (types.Transaction, error)
	tx, err := transaction.MakePaymentTxn(fromAddr, toAddr, fee, amount, lastRound-50, lastRound+50, nil, "", genID, genHash)
	if err != nil {
		return "", errors.Wrap(err, "error creating transaction")
	}

	// Sign the tx
	signResponse, err := KmdClient.SignTransaction(walletHandleToken, password, tx)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign transaction with kmd")
	}

	// Send the tx
	sendResponse, err := AlgodClient.SendRawTransaction(signResponse.SignedTransaction)
	if err != nil {
		return "", errors.Wrap(err, "failed to send transaction")
	}

	log.Println("Transaction ID: ", sendResponse.TxID)
	return sendResponse.TxID, nil
}

func CreateNewWallet(name string, password string) (string, error) {
	var err error
	KmdClient, err = kmd.MakeClient(consts.KmdAddress, consts.KmdToken)
	if err != nil {
		return "", errors.Wrap(err, "failed to make kmd client")
	}

	// Create the example wallet, if it doesn't already exist
	createWalletResponse, err := KmdClient.CreateWallet(name, password, kmd.DefaultWalletDriver, types.MasterDerivationKey{})
	if err != nil {
		return "", errors.Wrap(err, "error creating wallet")
	}

	// We need the wallet ID in order to get a wallet handle, so we can add accounts
	walletID := createWalletResponse.Wallet.ID
	fmt.Printf("Created wallet '%s' with ID: %s\n", createWalletResponse.Wallet.Name, walletID)

	// Get a wallet handle. The wallet handle is used for things like signing transactions
	// and creating accounts. Wallet handles do expire, but they can be renewed
	walletHandleToken, err := CreateNewWalletHandle(walletID, password)
	if err != nil {
		return "", errors.Wrap(err, "failed to create new wallet handler")
	}

	// Generate a new address from the wallet handle
	address, err := GenerateAddress(walletHandleToken)
	if err != nil {
		return "", errors.Wrap(err, "failed to gneerate new address")
	}

	return address, nil
}
