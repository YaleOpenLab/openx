package database

// recipient.go defines all recipient related functions that are not defined on
// the struct itself.
import (
	"encoding/json"
	"fmt"
	"log"

	oracle "github.com/YaleOpenLab/smartPropertyMVP/stellar/oracle"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
	"github.com/boltdb/bolt"
	"github.com/stellar/go/build"
)

type Recipient struct {
	ReceivedProjects []DBParams
	// ReceivedProjects denotes the projects that have been received by the recipient
	// instead of storing the PaybackAssets and the DebtAssets, we store this
	U User
	// user related functions are called as an instance directly
	// TODO: better name? idk
}

func NewRecipient(uname string, pwd string, seedpwd string, Name string) (Recipient, error) {
	var a Recipient
	var err error
	a.U, err = NewUser(uname, pwd, seedpwd, Name)
	if err != nil {
		return a, err
	}
	err = a.Save()
	return a, err
}

// all operations are mostly similar to that of the Recipient class
// TODO: merge where possible by adding an extra bucket param
func (a *Recipient) Save() error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(RecipientBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.ItoB(a.U.Index)), encoded)
	})
	return err
}

// RetrieveAllRecipients gets a list of all Recipient in the database
func RetrieveAllRecipients() ([]Recipient, error) {
	var arr []Recipient
	temp, err := RetrieveAllUsers()
	if err != nil {
		return arr, err
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(RecipientBucket)
		i := 1
		for ; i < limit; i++ {
			var rRecipient Recipient
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// this is where the key does not exist
				continue
			}
			err := json.Unmarshal(x, &rRecipient)
			//if err != nil && rRecipient.Live == false {
			if err != nil {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			arr = append(arr, rRecipient)
		}
		return nil
	})
	return arr, err
}

func RetrieveRecipient(key int) (Recipient, error) {
	var inv Recipient
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(RecipientBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return nil
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, err
}

func ValidateRecipient(name string, pwhash string) (Recipient, error) {
	var rec Recipient
	user, err := ValidateUser(name, pwhash)
	if err != nil {
		return rec, err
	}
	return RetrieveRecipient(user.Index)
}

// SendAssetToIssuer sends back assets fromn an asset holder to the issuer of the asset.
func (a *Recipient) SendAssetToIssuer(assetName string, issuerPubkey string, amount string, seed string) (int32, string, error) {
	// SendAssetToIssuer is FROM recipient / investor to issuer
	// TODO: the platform / issuer doesn't send back the PBToken since PBTOkens are
	// disabled as of now, can add back in later if needed.
	paymentTx, err := build.Transaction(
		build.SourceAccount{a.U.PublicKey},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: xlm.TestNetClient},
		build.Payment(
			build.Destination{AddressOrSeed: issuerPubkey},
			build.CreditAmount{assetName, issuerPubkey, amount},
		),
	)

	if err != nil {
		return -11, "", err
	}

	return xlm.SendTx(seed, paymentTx)
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
// project for people to pay more, we could also have a threshold mechanism that says
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
func (a *Recipient) Payback(uContract Project, assetName string, issuerPubkey string, amount string, seed string) error {
	// once we have the stablecoin here, we can remove the assetName
	StableBalance, err := xlm.GetAssetBalance(a.U.PublicKey, "STABLEUSD")
	// checks for the stablecoin asset
	if err != nil {
		log.Println("YOU HAVE NO STABLECOIN BALANCE, PLEASE REFILL ACCOUNT")
		return fmt.Errorf("YOU HAVE NO STABLECOIN BALANCE, PLEASE REFILL ACCOUNT")
	}

	DEBAssetBalance, err := xlm.GetAssetBalance(a.U.PublicKey, assetName)
	if err != nil {
		log.Println("Don't have the debt asset in possession")
		log.Fatal(err)
	}

	if utils.StoF(amount) > utils.StoF(StableBalance) {
		// check whether the recipient has enough StableUSD tokens in project to make
		// this happen
		log.Println("YOU CAN'T SEND AN AMOUNT MORE THAN WHAT YOU HAVE")
		return fmt.Errorf("YOU CAN'T SEND AN AMOUNT MORE THAN WHAT YOU HAVE")
	}
	// check balance in DEBAssetCode anmd
	monthlyBill, err := oracle.MonthlyBill()
	if err != nil {
		log.Println("Unable to fetch oracle price, exiting")
		return err
	}

	log.Println("Retrieved average price from oracle: ", monthlyBill)
	// the oracle needs to know the assetName so that it can find the other details
	// about this asset from the db. This should run on the server side and must
	// be split when we do run client side stuff.
	// hardcode for now, need to add the oracle here so that we
	// can do this dynamically
	// send amount worth DEBTokens back to issuer
	confHeight, txHash, err := a.SendAssetToIssuer(assetName, issuerPubkey, amount, seed)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Paid debt amount: ", amount, " back to issuer, tx hash: ", txHash, " ", confHeight)
	log.Println("Checking balance to see if our account was debited")
	newBalance, err := xlm.GetAssetBalance(a.U.PublicKey, assetName)
	if err != nil {
		log.Fatal(err)
	}

	newBalanceFloat := utils.StoF(newBalance)
	DEBAssetBalanceFloat := utils.StoF(DEBAssetBalance)
	mBillFloat := utils.StoF(monthlyBill)

	paidAmount := DEBAssetBalanceFloat - newBalanceFloat
	log.Println("Old Balance: ", DEBAssetBalanceFloat, "New Balance: ", newBalanceFloat, "Paid: ", paidAmount, "Bill Amount: ", mBillFloat)

	// would be nice to take some additional action like sending a notification or
	// something to investors or to the email address given so that everyone is made
	// aware of this and there's data transparency

	if paidAmount < mBillFloat {
		log.Println("Amount paid is less than amount required, balance not updating, please make sure to cover this next time")
	} else if paidAmount > mBillFloat {
		log.Println("You've chosen to pay more than what is required for this month. Adjusting payback period accordingly")
	} else {
		log.Println("You've paid exactly what is required for this month. Payback period remains as usual")
	}
	// we need to update the database here
	// no need to retrieve this project again because we have it already
	uContract.Params.BalLeft -= paidAmount
	uContract.Params.DateLastPaid = utils.Timestamp()
	if uContract.Params.BalLeft == 0 {
		log.Println("YOU HAVE PAID OFF THIS ASSET, TRANSFERRING OWNERSHIP OF ASSET TO YOU")
		// don't delete the asset from the received assets list, we still need it so
		// that we can look back and find out hwo many assets this particular
		// enttiy has been invested in, have a leaderboard kind of thing, etc.
		uContract.Stage = 7
		// we should call neighbourly or some ohter partner here to transfer assets
		// using the bond they provide us with
		// the nice part here is that the recipient can not pay off more than what is
		// invested because the trustline will not allow such an incident to happen
	}
	// balLeft must be updated on the server side and can be challenged easily
	// if there's some discrepancy since the tx's are on the blockchain
	err = a.UpdateProjectSlice(uContract.Params)
	if err != nil {
		return err
	}
	fmt.Println("UPDATED ORDER: ", uContract.Params)
	err = uContract.Save()
	if err != nil {
		return err
	}
	return err
}

func (a *Recipient) UpdateProjectSlice(project DBParams) error {
	pos := -1
	for i, mem := range a.ReceivedProjects {
		if mem.DEBAssetCode == project.DEBAssetCode {
			log.Println("Rewriting the thing in our copy")
			// rewrite the thing in memory that we have
			pos = i
			break
		}
	}
	if pos != -1 {
		// rewrite the thing in memory
		a.ReceivedProjects[pos] = project
		err := a.Save()
		return err
	}
	return fmt.Errorf("Not found")
}
