package xlm

// the xlm package is a package that interacts with the stellar testnet
// API and fetches testnet coins for the user
// right now, there are multiple fucntions which are not used through the program,
// would be nice to use them when we require so.
import (
	"fmt"
	"log"
	"net/http"

	consts "github.com/YaleOpenLab/smartPropertyMVP/stellar/consts"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
)

func GetKeyPair() (string, string, error) {
	pair, err := keypair.Random()
	return pair.Seed(), pair.Address(), err
}

// GetCoins makes an API call to the friendbot on stellar testnet, which gives
// us 10000 XLM for use. We don't need 10000XLM (we need only ~3 XLM for setting up
// various trustlines), but there's no option to receive less, so we're having to call
// this. On mainnet, we'd be refilling the accoutns manually, so this function
// wouldn't exist.
func GetXLM(PublicKey string) error {
	// get some coins from the stellar robot for testing
	// gives only a constant amount of stellar, so no need to pass it a coin param
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + PublicKey)
	if err != nil {
		log.Println("ERRORED OUT while calling friendbot, no coins for us")
		return err
	}
	if resp.Status != "200 OK" {
		return fmt.Errorf("API Request did not succeed")
	}
	return nil
}

// check whether an accoutn exists, not needed now since we do the check ourselves
// in multiple places
func AccountExists(address string) bool {
	_, err := TestNetClient.LoadAccount(address)
	if err != nil {
		return false
	}
	return true
}

func SendTx(Seed string, tx *build.TransactionBuilder) (int32, string, error) {
	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(Seed)
	if err != nil {
		return -1, "", err
	}

	txeB64, err := txe.Base64()
	if err != nil {
		return -1, "", err
	}
	// And finally, send it off to Stellar
	resp, err := TestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		return -1, "", err
	}

	fmt.Printf("Propagated Transaction: %s, sequence: %d", resp.Hash, resp.Ledger)
	return resp.Ledger, resp.Hash, nil
}

// Generating a keypair on stellar doesn't mean that you can send funds to it
// you need to call the CreateAccount method in project to be able to send funds
// to it
func SendXLMCreateAccount(destination string, amount string, Seed string) (int32, string, error) {
	// destination will not exist yet, so don't check
	passphrase := network.TestNetworkPassphrase
	tx, err := build.Transaction(
		build.SourceAccount{Seed},
		build.AutoSequence{TestNetClient},
		build.Network{passphrase},
		build.MemoText{"Sending Boootstrap Money"},
		build.CreateAccount(
			build.Destination{destination},
			build.NativeAmount{amount},
		),
	)

	if err != nil {
		return -1, "", err
	}

	return SendTx(Seed, tx)
}

// SendXLM sends _amount_ number of native tokens (XLM) to the specified destination
// address using the stellar testnet API
func SendXLM(destination string, amount string, Seed string, memo string) (int32, string, error) {
	// don't check if the account exists or not, hopefully it does
	passphrase := network.TestNetworkPassphrase
	tx, err := build.Transaction(
		build.SourceAccount{Seed},
		build.AutoSequence{TestNetClient},
		build.Network{passphrase},
		build.MemoText{memo},
		build.Payment(
			build.Destination{destination},
			build.NativeAmount{amount},
		),
	)

	if err != nil {
		return -1, "", err
	}

	return SendTx(Seed, tx)
}

func RefillAccount(publicKey string, platformSeed string) error {
	var err error
	if !AccountExists(publicKey) {
		// there is no account under the user's name
		// means we need to setup an account first
		log.Println("Account does not exist, creating: ", publicKey)
		_, _, err = SendXLMCreateAccount(publicKey, consts.DonateBalance, platformSeed)
		if err != nil {
			log.Println("Account Could not be created")
			return err
		}
	}
	// balance is in string, convert to float
	balance, err := GetNativeBalance(publicKey)
	if err != nil {
		return err
	}
	balanceI := utils.StoF(balance)
	if balanceI < 3 { // to setup trustlines
		_, _, err = SendXLM(publicKey, consts.DonateBalance, platformSeed, "Sending XLM to refill")
		if err != nil {
			log.Println("Account doesn't have funds")
			return err
		}
	}
	return nil
}
