package database

// the database package maintains read / write operations to the orderbook
// we need an orderbook because there is no state on Stellar which makes it
// difficult for us to store this on the blockchain. We use boltdb no sicne we don't
// do that much relational mapping, but in the case we need that, we can modify
// this package to do that.

// what do we need to store?
import (
	// "fmt"
	"encoding/json"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
)

// Order is what is advertised on the frontend where investors can choose what
// projects to invest in.
type Order struct {
	Index         uint32  // an Index to keep quick track of how many orders exist
	PanelSize     string  // size of the given panel, for diplsaying to the user who wants to bid stuff
	TotalValue    int     // the total money that we need from investors
	Location      string  // where this specific solar panel is located
	MoneyRaised   int     // total money that has been raised until now
	Metadata      string  // any other metadata can be stored here
	Live          bool    // check to see whether the current order is live or not
	INVAssetCode  string  // once all funds have been raised, we need to set assetCodes
	DEBAssetCode  string  // once all funds have been raised, we need to set assetCodes
	PBAssetCode   string  // once all funds have been raised, we need to set assetCodes
	BalLeft       float64 // denotes the balance left to pay by the party
	DateInitiated string  // date the order was created
	DateLastPaid  string  // date the order was last paid
	// Percentage raised is not stored in the database since that can be calculated by the UI
}

var OrdersBucket = []byte("Orders")

// TODO: need locks over this to ensure no one's using the db while we are
func OpenDB() (*bolt.DB, error) {
	db, err := bolt.Open("yol.db", 0600, nil)
	if err != nil {
		log.Println("Couldn't open database, exiting!")
		return nil, err
	}
	return db, nil
}

// NewOrder creates a new order struct with the order parameters pased to the function
// quite ugly with all the parameters passed, would be nice if we could rewrite this
// in a nicer way
func NewOrder(db *bolt.DB, panelSize string, totalValue int, location string, moneyRaised int, metadata string, INVAssetCode string, DEBAssetCode string, PBAssetCode string) (Order, error) {
	var a Order
	// need to get a new index since we have a small bug on that
	allOrders, err := RetrieveAllOrders(db)
	if err != nil {
		return a, err
	}

	if len(allOrders) == 0 {
		a.Index = 1
	} else {
		a.Index = uint32(len(allOrders) + 1)
	}
	a.PanelSize = panelSize
	a.TotalValue = totalValue
	a.Location = location
	a.MoneyRaised = moneyRaised
	a.Metadata = metadata
	a.Live = true
	a.INVAssetCode = INVAssetCode
	a.DEBAssetCode = DEBAssetCode
	a.PBAssetCode = PBAssetCode
	a.BalLeft = float64(totalValue)
	a.DateInitiated = utils.Timestamp()
	// need to insert this into the database
	err = InsertOrder(a, db)
	if err != nil {
		return a, err
	}
	return a, nil
}

// InsertOrder inserts a passed order into the given database
// TODO: need locks over insert and retrieve operations since BOLTdb supports only
// one operation at a time.
func InsertOrder(order Order, db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(OrdersBucket) // the orders bucket contains all our orders
		if err != nil {
			log.Fatal(err)
			return err
		}
		encoded, err := json.Marshal(order)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.Uint32toB(order.Index)), encoded)
	})
	return err
}

// RetrieveOrder retrieves the given value from the database corresponding to the
// given key
func RetrieveOrder(key uint32, db *bolt.DB) (Order, error) {
	var rOrder Order
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(OrdersBucket)
		if err != nil {
			return err
		}
		x := b.Get(utils.Uint32toB(key))
		err = json.Unmarshal(x, &rOrder)
		if err != nil {
			return err
		}
		return nil
	})
	return rOrder, err
}

// DeleteOrder deltes a given value corresponding to the ky from the database
// DeleteOrder should be used only in cases where something is wrong from our side
// while creating an order. For other cases, we should set Live to False and edit
// the order
// TODO: make delete not mess up with indices, which it currently does
func DeleteOrder(key uint32, db *bolt.DB) error {
	// deleting order might be dangerous since that would mess with the RetrieveAllOrders
	// function, have it in here for now, don't do too much with it / fiox retrieve all
	// to handle this case
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(OrdersBucket)
		if err != nil {
			return err
		}
		err = b.Delete(utils.Uint32toB(key))
		if err != nil {
			return err
		}
		log.Println("Deleted order with key: ", key)
		return nil
	})
	return err
}

// RetrieveAllOrders retrieves all orders from the given database
func RetrieveAllOrders(db *bolt.DB) ([]Order, error) {
	var arr []Order
	err := db.Update(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b, err := tx.CreateBucketIfNotExists(OrdersBucket)
		if err != nil {
			return err
		}
		i := uint32(1)
		for ; ; i++ {
			var rOrder Order
			x := b.Get(utils.Uint32toB(i))
			if x == nil {
				// this is where the key does not exist
				return nil
			}
			err := json.Unmarshal(x, &rOrder)
			if err != nil && rOrder.Live == false {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			arr = append(arr, rOrder)
		}
		return nil
	})
	return arr, err
}

// RetrieveAllOrdersWithoutDB retrieves all orders from the default database (for use only
// by frontend RPCs which want to query us)
func RetrieveAllOrdersWithoutDB() ([]Order, error) {
	var arr []Order
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b, err := tx.CreateBucketIfNotExists(OrdersBucket)
		if err != nil {
			return err
		}
		i := uint32(1)
		for ; ; i++ {
			var rOrder Order
			x := b.Get(utils.Uint32toB(i))
			if x == nil {
				// this is where the key does not exist
				return nil
			}
			err := json.Unmarshal(x, &rOrder)
			if err != nil && rOrder.Live == false {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			arr = append(arr, rOrder)
		}
		return nil
	})
	return arr, err
}

// RetrieveOrderRPC retrievs a single value corresponding to a given key from
// the default database. For use only by RPC calls
func RetrieveOrderRPC(key uint32) (Order, error) {
	var rOrder Order
	db, err := OpenDB()
	if err != nil {
		return rOrder, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(OrdersBucket)
		if err != nil {
			return err
		}
		x := b.Get(utils.Uint32toB(key))
		err = json.Unmarshal(x, &rOrder)
		if err != nil {
			return err
		}
		return nil
	})
	return rOrder, err
}
