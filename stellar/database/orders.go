package database

// the database package maintains read / write operations to the orderbook
// we need an orderbook because there is no state on Stellar which makes it
// difficult for us to store this on the blockchain. We use boltdb no sicne we don't
// do that much relational mapping, but in the case we need that, we can modify
// this package to do that.

// what do we need to store?
import (
	"encoding/json"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
)

// NewOrder creates a new order struct with the order parameters pased to the function
// quite ugly with all the parameters passed, would be nice if we could rewrite this
// in a nicer way
func NewOrder(panelSize string, totalValue int, location string, moneyRaised int, metadata string, INVAssetCode string, DEBAssetCode string, PBAssetCode string) (Order, error) {
	var a Order
	// need to get a new index since we have a small bug on that
	allOrders, err := RetrieveAllOrders()
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
	err = InsertOrder(a)
	if err != nil {
		return a, err
	}
	return a, nil
}

// DeleteOrder deltes a given value corresponding to the ky from the database
// DeleteOrder should be used only in cases where something is wrong from our side
// while creating an order. For other cases, we should set Live to False and edit
// the order
// TODO: make delete not mess up with indices, which it currently does
func DeleteOrder(key uint32) error {
	// deleting order might be dangerous since that would mess with the RetrieveAllOrders
	// function, have it in here for now, don't do too much with it / fiox retrieve all
	// to handle this case
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(OrdersBucket)
		err := b.Delete(utils.Uint32toB(key))
		if err != nil {
			return err
		}
		log.Println("Deleted order with key: ", key)
		return nil
	})
	return err
}

// RetrieveAllOrders retrieves all orders from the default database (for use only
// by frontend RPCs which want to query us)
func RetrieveAllOrders() ([]Order, error) {
	var arr []Order
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(OrdersBucket)
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

// RetrieveOrder retrievs a single value corresponding to a given key from
// the default database. For use only by RPC calls
func RetrieveOrder(key uint32) (Order, error) {
	var rOrder Order
	db, err := OpenDB()
	if err != nil {
		return rOrder, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(OrdersBucket)
		x := b.Get(utils.Uint32toB(key))
		err = json.Unmarshal(x, &rOrder)
		if err != nil {
			return err
		}
		return nil
	})
	return rOrder, err
}

// InsertOrder inserts a passed order into the given database
// TODO: need locks over insert and retrieve operations since BOLTdb supports only
// one operation at a time.
func InsertOrder(order Order) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(OrdersBucket)
		encoded, err := json.Marshal(order)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.Uint32toB(order.Index)), encoded)
	})
	return err
}
