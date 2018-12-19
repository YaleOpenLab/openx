package stablecoin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	assets "github.com/YaleOpenLab/smartPropertyMVP/stellar/assets"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
	oracle "github.com/YaleOpenLab/smartPropertyMVP/stellar/oracle"
	"github.com/boltdb/bolt"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

// the idea of this stablecoin package is to issue a stablecoin on stellar testnet
// so that we can test the function of something similar on mainnet. The stablecoin
// provider should be stored in a different database becuase we will not be migrating
// this.

// The idea is to issue a single USD token for every USD t hat we receive on our
// account, this should be automated and we must not have any kind of user interaction that is in
//  place here. We also need a stablecoin Code, which we shall call as "STABLEUSD"
// for easy reference. Most functions would be similar to the one in assets.go,
// but need to be tailored to suit our requirements

type StableIssuer struct {
	Index     uint32
	Seed      string
	PublicKey string
	// Fields are enough since this is a meta structure.
}

var StableUSD build.Asset
var Issuer StableIssuer

// STABLECOIN SEED IS: SDEG3MRXNFXSZVSPBVIT3TJXVXTEALMMWZMPNXHH4RFL2QGCALJVJSY2 and STABLECOIN PUBLICKEY IS GBAACP6UUXZAB5ZAYAHWEYLNKORWB36WVBZBXWNPFXQTDY2AIQFM6D7Y
var StableBucket = []byte("Stablecoins")

func CreateStableCoin() build.Asset {
	// need to set a couple flags here
	return build.CreditAsset(StableUSD.Code, Issuer.PublicKey)
}

func OpenDB() (*bolt.DB, error) {
	db, err := bolt.Open("sbc.db", 0600, nil)
	if err != nil {
		log.Println("Couldn't open database, exiting!")
		return db, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(StableBucket) // the orders bucket contains all our orders
		if err != nil {
			return err
		}
		return nil
	})
	return db, err
}

func InsertIssuer(a StableIssuer) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	// check whether another issuer already exists. if so, quit
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(StableBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.Uint32toB(a.Index)), encoded)
	})
	return err
}

func CheckStableIssuer() error {
	var rIssuer StableIssuer
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(StableBucket)
		i := uint32(1)
		x := b.Get(utils.Uint32toB(i))
		if x == nil {
			log.Println("Deteceted no other Stable issuer, returning")
			// this is where the key does not exist
			// and this is what we want
			return nil
		}
		err := json.Unmarshal(x, &rIssuer)
		if err != nil {
			return err
		}
		return nil
	})
	if rIssuer.Index != 0 {
		// this is the case that we want, so catch this and return
		fmt.Println("Found another stablecoin instance running, please remember the seed")
		return fmt.Errorf("Found another stablecoin instance running, please remember the seed")
	}
	return nil
}

func RetrieveStableIssuer() (StableIssuer, error) {
	// retrieves the platforms (more like the publickey)
	var rIssuer StableIssuer
	db, err := OpenDB()
	if err != nil {
		return rIssuer, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(StableBucket)
		i := uint32(1)
		x := b.Get(utils.Uint32toB(i))
		if x == nil {
			// this is where the key does not exist
			return nil
		}
		err := json.Unmarshal(x, &rIssuer)
		if err != nil {
			return nil
		}
		return nil
	})
	return rIssuer, err
}

func SetVals(PublicKey string, Seed string) {
	StableUSD.Code = "STABLEUSD"
	StableUSD.Issuer = PublicKey
	Issuer.PublicKey = PublicKey
	Issuer.Seed = Seed
}

func InitStableCoin() error {
	var x StableIssuer
	var err error

	if err = CheckStableIssuer(); err != nil {
		// there exits a stable issuer already,  retrieve and return publickey
		sI, err := RetrieveStableIssuer()
		if err != nil {
			return err
		}
		SetVals(sI.PublicKey, sI.Seed)
		return nil
	}
	// there is no instance of StableIssuer running, so instantiate one
	x.Seed, x.PublicKey, err = xlm.GetKeyPair()
	log.Printf("STABLECOIN SEED IS: %s and STABLECOIN PUBLICKEY IS %s", x.Seed, x.PublicKey)
	if err != nil {
		// don't return since we depend on this to work to continue further program
		// runs
		log.Fatal(err)
	}
	x.Index = 1 // only one stable instance
	err = InsertIssuer(x)
	if err != nil {
		// no way / ened to continue after this
		log.Fatal(err)
	}
	// set parameters for stablecoin ehre to avoid issues
	SetVals(x.PublicKey, x.Seed)
	err = xlm.GetXLM(x.PublicKey)
	if err != nil {
		// no way / ened to continue after this
		log.Fatal(err)
	}
	Issuer = x // set the local val to the global one
	_ = CreateStableCoin()
	return nil
}

// so now the above functions setup the stablecoin and create an asset. We ideally
// need a go routine that constantly listens for payments to this address (in XLM)
// and then hands out a fixed sum of 5USD per XLM deposited to the address which
// deposited the XLM
func ListenForPayments() {
	// this will be started as a goroutine
	// address := Issuer.PublicKey
	const address = "GBAACP6UUXZAB5ZAYAHWEYLNKORWB36WVBZBXWNPFXQTDY2AIQFM6D7Y"
	// this thing above has to be hardcoded because stellar's APi wants it like so
	// stupid stuff, but we need to go ahead with it
	ctx := context.Background() // start as a goroutine
	cursor := horizon.Cursor("now")
	fmt.Println("Waiting for a payment...")
	err := utils.DefaultTestNetClient.StreamPayments(ctx, address, &cursor, func(payment horizon.Payment) {
		/*
			Sample Response:
			Payment type payment
			Payment Paging Token 5424212982374401
			Payment From GC76MINOSNQUMDBNONBARFYFCCQA5QLNQSJOVANR5RNQVHQRB5B46B6I
			Payment To GBAACP6UUXZAB5ZAYAHWEYLNKORWB36WVBZBXWNPFXQTDY2AIQFM6D7Y
			Payment Asset Type native
			Payment Asset Code
			Payment Asset Issuer
			Payment Amount 10.0000000
			Payment Memo Type
			Payment Memo
		*/
		log.Println("Stablecoin payment to/from detected")
		log.Println("Payment type", payment.Type)
		log.Println("Payment Paging Token", payment.PagingToken)
		log.Println("Payment From", payment.From)
		log.Println("Payment To", payment.To)
		log.Println("Payment Asset Type", payment.AssetType)
		log.Println("Payment Asset Code", payment.AssetCode)
		log.Println("Payment Asset Issuer", payment.AssetIssuer)
		log.Println("Payment Amount", payment.Amount)
		log.Println("Payment Memo Type", payment.Memo.Type)
		log.Println("Payment Memo", payment.Memo.Value)
		if payment.Type == "payment" && payment.AssetType == "native" {
			// store the stuff that we want here
			payee := payment.From
			amount := payment.Amount
			log.Printf("Received request for stablecoin from %s worth %s", payee, amount)
			xlmWorth := oracle.ExchangeXLMforUSD(amount)
			log.Println("The deposited amount is worth: ", xlmWorth)
			// now send the stableusd asset over to this guy
			_, hash, err := assets.SendAssetFromIssuer(StableUSD.Code, payee, utils.FloatToString(xlmWorth), Issuer.Seed, Issuer.PublicKey)
			if err != nil {
				log.Println("Error while sending USD Assets back to payee: ", payee)
				//  don't skip here, there's technically nothing we can do
			}
			log.Println("Successful payment, hash: ", hash)
		}
	})

	if err != nil {
		panic(err)
	}
}
