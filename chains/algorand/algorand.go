package algorand

import (
	"github.com/pkg/errors"
	"log"

	utils "github.com/Varunram/essentials/utils"
	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/algod/models"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

// Algorand's model is similar to that of ethereum and stellar (account based model)
// these accounts have a name ("blah") and each one of them can have multiple addresses
// associated with them. Each address should have a minimum balance of 0.1 Algo (100,000 microAlgos)
// Algorand doesn't have a horizone API (PHEW!) and needs one to directly interact with the blockchain
// or create a private test network

// These constants represent the algod REST endpoint and the corresponding
// API token. You can retrieve these from the `algod.net` and `algod.token`
// files in the algod data directory.

// AlgodClient is a package-level gloabal variable
var AlgodClient algod.Client

// KmdClient is a package-level gloabal variable
var KmdClient kmd.Client

// InitAlgodClient initializes a new algorand daemon client
func InitAlgodClient() (algod.Client, error) {
	var err error
	AlgodClient, err = algod.MakeClient(AlgodAddress, AlgodToken)
	if err != nil {
		log.Printf("failed to make algod client: %s\n", err)
		return AlgodClient, nil
	}

	return AlgodClient, nil
}

// InitKmdClient initializes a new key management daemon client
func InitKmdClient() (kmd.Client, error) {
	return kmd.MakeClient(KmdAddress, KmdToken)
}

// Init initializes the algod client and the kmd client
func Init() error {
	var err error
	AlgodClient, err = InitAlgodClient()
	if err != nil {
		return err
	}

	KmdClient, err = InitKmdClient()
	return err
}

// GetLatestBlock getst the latest block from the blockchain
func GetLatestBlock(status models.NodeStatus) (models.Block, error) {
	return AlgodClient.Block(status.LastRound)
}

// GetBlock gets the details ofa  given block from the algorand blockchain
func GetBlock(blockNumber uint64) (models.Block, error) {
	return AlgodClient.Block(blockNumber)
}

// GetStatus gets the status of a given algod client
func GetStatus(Client algod.Client) (models.NodeStatus, error) {
	var status models.NodeStatus
	status, err := AlgodClient.Status()
	if err != nil {
		log.Printf("error getting algod status: %s\n", err)
		return status, err
	}

	return status, nil
}

// CreateNewWallet creates a new wallet
func CreateNewWallet(name string, password string) (string, error) {
	response, err := KmdClient.CreateWallet(name, password, kmd.DefaultWalletDriver, types.MasterDerivationKey{})
	if err != nil {
		return "", errors.Wrap(err, "error creating wallet")
	}

	walletID := response.Wallet.ID
	return walletID, nil
}

// generateWalletToken creates a wallet handle and is used for things like signing transactions
// and creating accounts. Wallet handles do expire, but they can be renewed
func generateWalletToken(walletID string, password string) (string, error) {
	// Get a wallet handle. The wallet handle is used for things like signing transactions
	// and creating accounts. Wallet handles do expire, but they can be renewed
	initResponse, err := KmdClient.InitWalletHandle(walletID, password)
	if err != nil {
		log.Printf("Error initializing wallet handle: %s\n", err)
		return "", err
	}

	walletHandleToken := initResponse.WalletHandleToken
	return walletHandleToken, nil
}

// generateAddress generates an address from a given wallet handle
func generateAddress(walletHandleToken string) (string, error) {
	// Generate a new address from the wallet handle
	genResponse, err := KmdClient.GenerateKey(walletHandleToken)
	if err != nil {
		log.Printf("Error generating key: %s\n", err)
		return "", err
	}
	log.Printf("Generated address %s\n", genResponse.Address)
	return genResponse.Address, nil
}

func sendTx(fromAddr string, toAddr string, amount uint64, note []byte,
	walletHandleToken string, password string) (string, error) {
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

// GetAlgo seeds a given wallet with a specific number of algos similar to what friendbot
// does on stellar
func GetAlgo(walletName string, password string) (string, error) {
	ourWalletId, err := getWalletId(walletName)
	if err != nil {
		return "", errors.Wrap(err, "couldn't get wallet id of wallet, quitting")
	}

	// get a wallet handle token to sign the transaction with
	walletHandleToken, err := generateWalletToken(ourWalletId, password)
	if err != nil {
		return "", errors.Wrap(err, "error initializing wallet handle")
	}

	// Generate a new address from the wallet handle
	toAddr, err := generateAddress(walletHandleToken)
	if err != nil {
		return "", errors.Wrap(err, "error generating key")
	}

	log.Println("Generated to address: ", toAddr)
	fromAddr := "YXU3MTTKV74UAGED6ROTHVVPEY5646WI3N5FLLQZWFV66AFKVQ5PMMYDZE"
	amount := uint64(150000) // 100000 is the minimum balance an account can hold

	note := []byte("cool")

	txid, err := sendTx(fromAddr, toAddr, amount, note, walletHandleToken, password)
	if err != nil {
		return "", errors.Wrap(err, "couldn't send tx")
	}
	log.Println("txid: ", txid)
	return txid, nil
}

// SendAlgoToSelf sends algos to another address owned by the same user
func SendAlgoToSelf(walletName string, password string, fromAddr string, amount uint64) (string, error) {
	// Get the list of wallets
	ourWalletId, err := getWalletId(walletName)
	if err != nil {
		return "", errors.Wrap(err, "couldn't get wallet id of wallet, quitting")
	}

	// get a wallet handle token to sign the transaction with
	walletHandleToken, err := generateWalletToken(ourWalletId, password)
	if err != nil {
		return "", errors.Wrap(err, "error initializing wallet handle")
	}

	// Generate a new address from the wallet handle
	toAddr, err := generateAddress(walletHandleToken)
	if err != nil {
		return "", errors.Wrap(err, "error generating key")
	}

	log.Println("Generated to address: ", toAddr)

	note := []byte("cool")
	txid, err := sendTx(fromAddr, toAddr, amount, note, walletHandleToken, password)
	if err != nil {
		return "", errors.Wrap(err, "couldn't send tx")
	}
	log.Println("txid: ", txid)
	return txid, nil
}

func getWalletId(walletName string) (string, error) {
	listResponse, err := KmdClient.ListWallets()
	if err != nil {
		return "", errors.Wrap(err, "error listing wallets")
	}

	// Find our walletID in the list
	var ourWalletId string
	for _, wallet := range listResponse.Wallets {
		if wallet.Name == walletName {
			log.Printf("found wallet '%s' with ID: %s\n", wallet.Name, wallet.ID)
			ourWalletId = wallet.ID
		}
	}

	return ourWalletId, nil
}

// SendAlgo sends algos to another address from a source account
func SendAlgo(walletName string, password string, amount uint64, fromAddr string, toAddr string) (string, error) {
	// Get the list of wallets
	ourWalletId, err := getWalletId(walletName)
	if err != nil {
		return "", errors.Wrap(err, "couldn't get wallet id of wallet, quitting")
	}

	// get a wallet handle token to sign the transaction with
	walletHandleToken, err := generateWalletToken(ourWalletId, password)
	if err != nil {
		return "", errors.Wrap(err, "error initializing wallet handle")
	}

	note := []byte("cool")
	txid, err := sendTx(fromAddr, toAddr, amount, note, walletHandleToken, password)
	if err != nil {
		return "", errors.Wrap(err, "couldn't send tx")
	}
	log.Println("txid: ", txid)
	return txid, nil
}

// CreateNewWalletAndAddress creates a new wallet and an address
func CreateNewWalletAndAddress(name string, password string) (string, error) {
	var err error

	walletID, err := CreateNewWallet(name, password)
	if err != nil {
		return "", errors.Wrap(err, "couldn't create new wallet id, quitting")
	}

	return GenerateNewAddress(walletID, name, password)
}

// GenerateNewAddress generates a new address associated with the given wallet
func GenerateNewAddress(walletID string, name string, password string) (string, error) {
	var err error

	walletHandleToken, err := generateWalletToken(walletID, password)
	if err != nil {
		return "", errors.Wrap(err, "failed to create new wallet handler")
	}

	// Generate a new address from the wallet handle
	address, err := generateAddress(walletHandleToken)
	if err != nil {
		return "", errors.Wrap(err, "failed to gneerate new address")
	}

	return address, nil
}

// GenerateBackup gets the seedphrase from the walletName for backup
func GenerateBackup(walletName string, password string) (string, error) {
	// Get the list of wallets
	ourWalletId, err := getWalletId(walletName)
	if err != nil {
		return "", errors.Wrap(err, "couldn't get wallet id of wallet, quitting")
	}

	// get a wallet handle token to sign the transaction with
	walletHandleToken, err := generateWalletToken(ourWalletId, password)
	if err != nil {
		return "", errors.Wrap(err, "error initializing wallet handle")
	}

	// Get the backup phrase
	resp, err := KmdClient.ExportMasterDerivationKey(walletHandleToken, password)
	if err != nil {
		return "", errors.Wrap(err, "error exporting backup phrase")
	}

	// This string should be kept in a safe place and not shared
	backupPhrase, err := mnemonic.FromKey(resp.MasterDerivationKey[:])
	if err != nil {
		return "", errors.Wrap(err, "error getting backup phrase")
	}

	return backupPhrase, nil
}

// AlgorandWallet defines the algorand wallet strcuture
type AlgorandWallet struct {
	WalletName string
	WalletID   string
}

// GenNewWallet generates a new algorand wallet
func GenNewWallet(walletName string, password string) (AlgorandWallet, error) {

	var x AlgorandWallet
	var err error
	if len(walletName) > 16 {
		return x, errors.New("wallet name too long, quitting")
	}

	x.WalletName = "algowl" + utils.GetRandomString(16-len(walletName))
	x.WalletID, err = CreateNewWallet(x.WalletName, password)
	if err != nil {
		return x, errors.Wrap(err, "couldn't create new wallet id, quitting")
	}

	return x, nil
}
