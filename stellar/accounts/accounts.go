// the accoutns package is a meta package that interacts with the stellar testnet
// APi and fetches coins initially for the user.
// the accounts package deals with connecting to the stellar testnet API to fetch
// coins for the issuer,
package accounts

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	"github.com/boltdb/bolt"
	"github.com/stellar/go/build"
	clients "github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon" // using this since hte client/horizon package has some deprecated fields
)

type Account struct {
	Seed      string
	PublicKey string
}

var DefaultTestNetClient = &clients.Client{
	URL:  "https://horizon-testnet.stellar.org",
	HTTP: http.DefaultClient,
}

// SetupAccount is a handler to setup a new account (issuer / investor / school)
func SetupAccount() Account {
	a, err := New()
	if err != nil {
		log.Fatal(err)
	}
	return a
}

// New() generates a new ed25519 keypair, assigns them to the meta strucutre Account
// and Returns it
func New() (Account, error) {
	var a Account
	pair, err := keypair.Random()
	// so key value pairs over here are ed25519 key pairs instead of bitcoin style key pairs
	// they also seem to sue al lcaps, which I don't know why
	// friendbot creates the account for us, on mainnet, there is no friendbot, so
	// we need to fil an address and then create accounts from that
	if err != nil {
		return a, err
	}
	// log.Println("MY SEED IS: ", pair.Seed())
	a.Seed = pair.Seed()
	a.PublicKey = pair.Address()
	return a, nil
}

// Account.SetupAccount() is a method on the structure Account that
// creates a new account using the stellar build.CreateAccount function and
// sends _amount_ number of stellar lumens to the newly created account.
// Note that the destination must alreayd have a keypair generated for this to work
// or else we'd be burning the coins since we wouldn't have the public key
// associated with it,
func (issuer *Account) SetupAccount(recipientPubKey string, amount string) error {
	passphrase := network.TestNetworkPassphrase
	// we need to set a couple flags here to make sure that the issuer can't
	// withdraw the asset, the whole "trustless" thing. Also, this is verifiable
	// and that's nice (you could just read the flags)
	// weird part is that I can't seem to figure out how to set these flags,
	// so leaving this until the end
	tx, err := build.Transaction(
		build.SourceAccount{issuer.Seed},
		build.AutoSequence{DefaultTestNetClient},
		build.Network{passphrase},
		build.CreateAccount(
			build.Destination{recipientPubKey},
			build.NativeAmount{amount},
			// build.SetAuthRequired(),
			// build.SetAuthImmutable(),
		),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	txe, err := tx.Sign(issuer.Seed)
	if err != nil {
		fmt.Println(err)
		return err
	}

	txeB64, err := txe.Base64()

	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)
	log.Println("LEDGER: ", resp.Ledger, "Hash: ", resp.Hash)
	return nil
}

// GetCoins makes an API call to the friendbot on stellar testnet, which gives
// us 10000 XLM for use. We don't need 10000XLM (we need only ~3 XLM for setting up
// various trustlines), but there's no option to receive less, so we're having to call
// this. On mainnet, we'd be refilling the accoutns manually, so this function
// wouldn't exist.
func (a *Account) GetCoins() error {
	// get some coins from the stellar robot for testing
	// gives only a constant amount of stellar, so no need to pass it a coin param
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + a.PublicKey)
	if err != nil || resp == nil {
		log.Println("ERRORED OUT while calling friendbot, no coins for us")
		return err
	}
	return nil
}

// GetAssetBalance calls the stellar testnet API to get all balances
// and then runs through the balances to get the balance of a specific account
func (a *Account) GetAssetBalance(assetCode string) (string, error) {

	account, err := DefaultTestNetClient.LoadAccount(a.PublicKey)
	if err != nil {
		return "", nil
	}

	for _, balance := range account.Balances {
		if balance.Asset.Code == assetCode {
			return balance.Balance, nil
		}
	}

	return "", nil
}

// GetAllBalances calls  the stellar testnet API to get all the balances associated
// with a certain account.
func (a *Account) GetAllBalances() ([]horizon.Balance, error) {

	account, err := DefaultTestNetClient.LoadAccount(a.PublicKey)
	if err != nil {
		return nil, nil
	}

	return account.Balances, nil
}

// SendCoins sends _amount_ number of native tokens (XLM) to the specified destination
// address using the stellar testnet API
func (a *Account) SendCoins(destination string, amount string) (int32, string, error) {

	if _, err := DefaultTestNetClient.LoadAccount(destination); err != nil {
		// if destination doesn't exist, do nothing
		// returning -11 since -1 maybe returned for unconfirmed tx or something like that
		return -11, "", err
	}

	passphrase := network.TestNetworkPassphrase

	tx, err := build.Transaction(
		build.Network{passphrase},
		build.SourceAccount{a.Seed},
		build.AutoSequence{DefaultTestNetClient},
		build.Payment(
			build.Destination{destination},
			build.NativeAmount{amount},
		),
	)

	if err != nil {
		return -11, "", err
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(a.Seed)
	if err != nil {
		return -11, "", err
	}

	txeB64, err := txe.Base64()
	if err != nil {
		return -11, "", err
	}
	// And finally, send it off to Stellar!
	resp, err := DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		return -11, "", err
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)
	return resp.Ledger, resp.Hash, nil
}

// CreateAsset is a method on Account that creates a new asset with code
// _assetName_ belonging to the caller
func (a *Account) CreateAsset(assetName string) build.Asset {
	// need to set a couple flags here
	return build.CreditAsset(assetName, a.PublicKey)
}

// TrustAsset creates a trustline from the caller towards the specific asset
// and asset issuer with a _limit_ set on the maximum amount of tokens that can be sent
// through the trust channel. Each trustline costs 0.5XLM.
func (a *Account) TrustAsset(asset build.Asset, limit string) (string, error) {
	// TRUST is FROM recipient TO issuer
	trustTx, err := build.Transaction(
		build.SourceAccount{a.PublicKey},
		build.AutoSequence{SequenceProvider: DefaultTestNetClient},
		build.TestNetwork,
		build.Trust(asset.Code, asset.Issuer, build.Limit(limit)),
	)

	if err != nil {
		return "", err
	}

	trustTxe, err := trustTx.Sign(a.Seed)
	if err != nil {
		return "", err
	}

	trustTxeB64, err := trustTxe.Base64()
	if err != nil {
		return "", err
	}

	tx, err := DefaultTestNetClient.SubmitTransaction(trustTxeB64)
	if err != nil {
		return "", err
	}

	log.Println("Trusted asset tx: ", tx.Hash)
	return tx.Hash, nil
}

// SendAsset transfers _amount_ number of assets from the caller to the destination
// and returns an error if the destination doesn't have a trustline with the issuer
// This method is called by the issuer of the asset
func (a *Account) SendAsset(assetName string, destination string, amount string) (int32, string, error) {
	// this transaction is FROM issuer TO recipient
	paymentTx, err := build.Transaction(
		build.SourceAccount{a.PublicKey},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: DefaultTestNetClient},
		build.Payment(
			build.Destination{AddressOrSeed: destination},
			build.CreditAmount{assetName, a.PublicKey, amount},
			// CreditAmount identifies the asset by asset Code and issuer pubkey
		),
	)

	if err != nil {
		return -11, "", err
	}

	paymentTxe, err := paymentTx.Sign(a.Seed)
	if err != nil {
		return -11, "", err
	}

	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {
		return -11, "", err
	}

	tx, err := DefaultTestNetClient.SubmitTransaction(paymentTxeB64)
	if err != nil {
		return -11, "", err
	}

	return tx.Ledger, tx.Hash, nil
}

// SendAssetToIssuer sends back assets fromn an asset holder to the issuer of the asset.
// This method is called by the receiver of assets.
func (a *Account) SendAssetToIssuer(assetName string, issuerPubkey string, amount string) (int32, string, error) {
	// SendAssetToIssuer is FROM recipient / investor to issuer
	paymentTx, err := build.Transaction(
		build.SourceAccount{a.PublicKey},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: DefaultTestNetClient},
		build.Payment(
			build.Destination{AddressOrSeed: issuerPubkey},
			build.CreditAmount{assetName, issuerPubkey, amount},
		),
	)

	if err != nil {
		return -11, "", err
	}

	paymentTxe, err := paymentTx.Sign(a.Seed)
	if err != nil {
		return -11, "", err
	}

	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {
		return -11, "", err
	}

	tx, err := DefaultTestNetClient.SubmitTransaction(paymentTxeB64)
	if err != nil {
		return -11, "", err
	}

	return tx.Ledger, tx.Hash, nil
}

// PriceOracle returns the power tariffs and any data that we need to  certify
// that is in the real world. Right now, this is hardcoded since we need to come up
// with a construct to get the price data in a reliable way - this could be a website
// were poeple erport this or certified authorities can timestamp this on chain
// or similar. Web s craping governemnt websites might work, but that seems too
// overkill for what we're doing now.
func PriceOracle() (string, error) {
	// right now, community consensus look like the price of electricity is
	// $0.2 per kWH in Puerto Rico, so hardcoding that here.
	priceOfElectricity := 0.2
	// since solar is free, they just need to pay this and then in some x time (chosen
	// when the order is created / confirmed on the school side), they
	// can own the panel.
	// the average energy consumption in puerto rico seems to be 5,657 kWh or about
	// 471 kWH per household. lets take 600 accounting for a 20% error margin.
	averageConsumption := float64(600)
	avgString := utils.FloatToString(priceOfElectricity*averageConsumption)
	return avgString, nil
}

// PriceOracleInFloat does the same thing as PriceOracle, but returns the data
// as a float for use in appropriate places
func PriceOracleInFloat() (float64) {
	priceOfElectricity := 0.2
	averageConsumption := float64(600)
	return priceOfElectricity*averageConsumption
}

// Payback is called when the receiver of the DEBToken wants to pay a fixed amount
// of money back to the issuer of the DEBTokens. One way to imagine this would be
// like an electricity bill, something that people pay monthly but only that in this
// case, the electricity is free, so they pay directly towards the solar panels.
// The process of Payback roughly involves the followign steps:
// 1. Pay the issuer in DEBTokens with whatever amount desired.
// The oracle price of
// electricity cost is a lower bound (since the government would not like it if people
// default on their payments). Anything below the lower bound gets a warning in
// order for people to pay more, we could also have a threshold mechanism that says
// if a person constantly defaults for more than half the owed amount for three
// consecutive months, we sell power directly to the grid. THis could also be used
// for a rating system, where the frontend UI can have a rating based on whether
// the recipient has defaulted or not in the past.
// 2. The receiver checks whether the amount is greater than Oracle Threshold and
// if so, sends back PBTokens, which stand for the month equivalent of payments.
// eg. the school has opted for a 5 year payback period, the school owes the issuer
// 60 PBTokens and the issuer sends back 1PBToken every month if the school pays
// invested_amount/60 DEBTokens back to the issuer
// 3. The recipient checks whether the PBTokens received correlate to the amount
// that it sent and if not, raises the dispute since the forward DEBToken payment
// is on chain and resolves the dispute itself using existing off chain legal frameworks
// (issued bonds, agreements, etc)
func (a *Account) Payback(db *bolt.DB, index uint32, assetName string, issuerPubkey string, amount string) error {

	oldBalance, err := a.GetAssetBalance(assetName)
	if err != nil {
		log.Fatal(err)
	}

	PBAmount, err := PriceOracle()
	if err != nil {
		log.Println("Unable to fetch oracle price, exiting")
		return err
	}

	log.Println("Retrieved average price from oracle: ", PBAmount)
	// the oracke needs to know the assetName so that it can find the other details
	// about this asset from the db. This should run on the server side and must
	// be split when we do run client side stuff.
	// hardcode for now, need to add the oracle here so that we
	// can do this dynamically
	confHeight, txHash, err := a.SendAssetToIssuer(assetName, issuerPubkey, amount)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Paid debt amount: ", amount, " back to issuer, tx hash: ", txHash, " ", confHeight)
	log.Println("Checking balance to see if our account was debited")
	newBalance, err := a.GetAssetBalance(assetName)
	if err != nil {
		log.Fatal(err)
	}

	newBalanceFloat, err := strconv.ParseFloat(newBalance, 32) // 32 bit floats
	if err != nil {
		log.Println(err)
		return err
	}
	oldBalanceFloat, err := strconv.ParseFloat(oldBalance, 32)
	if err != nil {
		log.Println(err)
		return err
	}
	amountFloat, err := strconv.ParseFloat(PBAmount, 32)
	if err != nil {
		log.Println(err)
		return err
	}

	paidAmount := oldBalanceFloat - newBalanceFloat
	log.Println("Old Balance: ", oldBalanceFloat, "New Balance: ", newBalanceFloat, "Paid: ", paidAmount, "Req Amount: ", amountFloat)

	// would be nice to take some additional action like sending a notification or
	// something to investors or to the email address given so that everyone is made
	// aware of this and there's data transparency

	if paidAmount < amountFloat {
		log.Println("Amount paid is less than amount required, balance not updating, please amke sure to cover this next time")
	} else if paidAmount > amountFloat {
		log.Println("You've chosen to pay more than what is required for this month. Adjusting payback period accordingly")
	} else {
		log.Println("You've paid exactly what is required for this month. Payback period remains as usual")
	}
	// we need to update the database here
	givenOrder, err := database.RetrieveOrder(index, db)
	if err != nil {
		log.Println("Given order not found in the database")
		return err
	}
	givenOrder.BalLeft = float64(givenOrder.TotalValue) - paidAmount
	givenOrder.DateLastPaid = time.Now().Format(time.RFC850)
	// balLeft must be updated on the server side and can be challenged easily
	// if there's some discrepancy since the tx's are on the blockchain
	return database.InsertOrder(givenOrder, db)
}
