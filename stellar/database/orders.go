package database

// the database package maintains read / write operations to the orderbook
// we need an orderbook because there is no state on Stellar which makes it
// difficult for us to store this on the blockchain. We use boltdb no since we don't
// do that much relational mapping, but in the case we need that, we can modify
// this package to do that.

// what do we need to store?
import (
	"encoding/json"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"github.com/boltdb/bolt"
)

// Order is what is advertised on the frontend where investors can choose what
// projects to invest in.
// the idea of an Order should be constructed by ourselves based on some preset
// parameters. It should hold a contractor, developers, guarantor, originator
// along with the fields already defined. The other fields should be displayed
// on the frontend.
// Order is the structure that is actually stored in the database
// Order is a subset of the larger struct Contractor, which is still a TODO
type Order struct {
	Index         int     // an Index to keep quick track of how many orders exist
	PanelSize     string  // size of the given panel, for diplsaying to the user who wants to bid stuff
	TotalValue    int     // the total money that we need from investors
	Location      string  // where this specific solar panel is located
	MoneyRaised   int     // total money that has been raised until now
	Years         int     // number of years the recipient has chosen to opt for
	Metadata      string  // any other metadata can be stored here
	Live          bool    // check to see whether the current order is live or not
	Origin        bool    // if this order is an originated order
	Votes         int     // the number of votes of a proposed contract
	Stage         int     // the stage at which the specific order is at, can be used to see the progess of order
	PaidOff       bool    // whether the asset has been paidoff by the recipient
	INVAssetCode  string  // once all funds have been raised, we need to set assetCodes
	DEBAssetCode  string  // once all funds have been raised, we need to set assetCodes
	PBAssetCode   string  // once all funds have been raised, we need to set assetCodes
	BalLeft       float64 // denotes the balance left to pay by the party
	DateInitiated string  // date the order was created
	DateFunded    string  // date when the order was funded
	DateLastPaid  string  // date the order was last paid
	// instead of just holding the recipient's name here, we can hold the recipient
	OrderRecipient Recipient
	// also have an array of investors to keep track of who invested in these projects
	OrderInvestors []Investor
	// TODO: have an investor and recipient relation here
	// Percentage raised is not stored in the database since that can be calculated by the UI
}

// NewOrder creates a new order struct with the order parameters passed to the function
// quite ugly with all the parameters passed, would be nice if we could rewrite this
// in a nicer way
func NewOrder(panelSize string, totalValue int, location string, moneyRaised int, metadata string) (Order, error) {
	var a Order
	// need to get a new index since we have a small bug on that
	allOrders, err := RetrieveAllOrders()
	if err != nil {
		return a, err
	}

	if len(allOrders) == 0 {
		a.Index = 1
	} else {
		a.Index = len(allOrders) + 1
	}
	a.PanelSize = panelSize
	a.TotalValue = totalValue
	a.Location = location
	a.MoneyRaised = moneyRaised
	a.Metadata = metadata
	a.Live = true
	a.BalLeft = float64(totalValue)
	a.DateInitiated = utils.Timestamp()
	// need to insert this into the database
	err = InsertOrder(a)
	if err != nil {
		return a, err
	}
	return a, nil
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
		return b.Put([]byte(utils.ItoB(order.Index)), encoded)
	})
	return err
}

// RetrieveOrder retrievs a single value corresponding to a given key from
// the default database. For use only by RPC calls
func RetrieveOrder(key int) (Order, error) {
	var rOrder Order
	db, err := OpenDB()
	if err != nil {
		return rOrder, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(OrdersBucket)
		x := b.Get(utils.ItoB(key))
		err = json.Unmarshal(x, &rOrder)
		if err != nil {
			return err
		}
		return nil
	})
	return rOrder, err
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
		for i := 1; ; i++ {
			var rOrder Order
			x := b.Get(utils.ItoB(i))
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
			if !rOrder.Origin {
				// only return final orders
				arr = append(arr, rOrder)
			}
		}
		return nil
	})
	return arr, err
}

// RetrieveAllProposedOrders retrieves proposed orders from the default database
func RetrieveAllOriginatedOrders() ([]Order, error) {
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
		for i := 1; ; i++ {
			var rOrder Order
			x := b.Get(utils.ItoB(i))
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
			if rOrder.Origin {
				// only return originated orders, so that the function calling this can
				// display the originated orders
				arr = append(arr, rOrder)
			}
		}
		return nil
	})
	return arr, err
}

// DeleteOrder deltes a given value corresponding to the ky from the database
// DeleteOrder should be used only in cases where something is wrong from our side
// while creating an order. For other cases, we should set Live to False and edit
// the order
// TODO: make delete not mess up with indices, which it currently does
func DeleteOrder(key int) error {
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
		// TODO: we must first retrieve the order to see if it exists before trying to delete it right away
		err := b.Delete(utils.ItoB(key))
		if err != nil {
			return err
		}
		log.Println("Deleted order with key: ", key)
		return nil
	})
	return err
}
