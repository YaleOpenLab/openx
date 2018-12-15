package database

// the database package maintains read / write operations to the orderbook
// we need an orderbook because there is no state on Stellar which makes it
// difficult for us to store this on the blockchain. We use boltdb no sicne we don't
// do that much relational mapping, but in the case we need that, we can modify
// this package to do that.

// what do we need to store?
import (
	"encoding/json"
	"fmt"
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
	Years         int     // number of years the recipient has chosen to opt for
	Metadata      string  // any other metadata can be stored here
	Live          bool    // check to see whether the current order is live or not
	INVAssetCode  string  // once all funds have been raised, we need to set assetCodes
	DEBAssetCode  string  // once all funds have been raised, we need to set assetCodes
	PBAssetCode   string  // once all funds have been raised, we need to set assetCodes
	BalLeft       float64 // denotes the balance left to pay by the party
	DateInitiated string  // date the order was created
	DateLastPaid  string  // date the order was last paid
	RecipientName string  // name of the recipient in order to assign the given assets
	// TODO: have an investor and recipient relation here
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

// InsertOrder inserts a passed order into the given database
// TODO: need locks over insert and retrieve operations since BOLTdb supports only
// one operation at a time.
func InsertOrderRPC(order Order) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
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

// PrettyPrintOrder pretty prints orders
func PrettyPrintOrders(orders []Order) {
	for _, order := range orders {
		fmt.Println("    ORDER NUMBER: ", order.Index)
		fmt.Println("          Panel Size: ", order.PanelSize)
		fmt.Println("          Total Value: ", order.TotalValue)
		fmt.Println("          Location: ", order.Location)
		fmt.Println("          Money Raised: ", order.MoneyRaised)
		fmt.Println("          Metadata: ", order.Metadata)
		fmt.Println("          Years: ", order.Years)
		if order.Live {
			fmt.Println("          Investor Asset Code: ", order.INVAssetCode)
			fmt.Println("          Debt Asset Code: ", order.DEBAssetCode)
			fmt.Println("          Payback Asset Code: ", order.PBAssetCode)
			fmt.Println("          Balance Left: ", order.BalLeft)
		}
		fmt.Println("          Date Initiated: ", order.DateInitiated)
		if order.Live {
			fmt.Println("          Date Last Paid: ", order.DateLastPaid)
		}
	}
}

// PrettyPrintOrder pretty prints orders
func PrettyPrintOrder(order Order) {
	fmt.Println("    ORDER NUMBER: ", order.Index)
	fmt.Println("          Panel Size: ", order.PanelSize)
	fmt.Println("          Total Value: ", order.TotalValue)
	fmt.Println("          Location: ", order.Location)
	fmt.Println("          Money Raised: ", order.MoneyRaised)
	fmt.Println("          Metadata: ", order.Metadata)
	fmt.Println("          Years: ", order.Years)
	if order.Live {
		fmt.Println("          Investor Asset Code: ", order.INVAssetCode)
		fmt.Println("          Debt Asset Code: ", order.DEBAssetCode)
		fmt.Println("          Payback Asset Code: ", order.PBAssetCode)
		fmt.Println("          Balance Left: ", order.BalLeft)
	}
	fmt.Println("          Date Initiated: ", order.DateInitiated)
	if order.Live {
		fmt.Println("          Date Last Paid: ", order.DateLastPaid)
	}
}

func InsertDummyData() error {
	var err error
	// populate database with dumym data
	var order1 Order

	order1.Index = 1
	order1.PanelSize = "100 1000 sq.ft homes each with their own private spaces for luxury"
	order1.TotalValue = 14000
	order1.Location = "India Basin, San Francisco"
	order1.MoneyRaised = 0
	order1.Metadata = "India Basin is an upcoming creative project based in San Francisco that seeks to invite innovators from all around to participate"
	order1.Live = false
	order1.INVAssetCode = ""
	order1.DEBAssetCode = ""
	order1.PBAssetCode = ""
	order1.DateInitiated = utils.Timestamp()
	order1.Years = 3
	order1.RecipientName = "Martin" // this is not the username of the recipient
	err = InsertOrderRPC(order1)
	if err != nil {
		return fmt.Errorf("Error inserting order into db")
	}

	order1.Index = 2
	order1.PanelSize = "180 1200 sq.ft homes in a high rise building 0.1mi from Kendall Square"
	order1.TotalValue = 30000
	order1.Location = "Kendall Square, Boston"
	order1.MoneyRaised = 0
	order1.Metadata = "Kendall Square is set in the heart of Cambridge and is a popular startup IT hub"
	order1.Live = false
	order1.INVAssetCode = ""
	order1.DEBAssetCode = ""
	order1.PBAssetCode = ""
	order1.DateInitiated = utils.Timestamp()
	order1.Years = 5
	order1.RecipientName = "Martin" // this is not the username of the recipient

	err = InsertOrderRPC(order1)
	if err != nil {
		return fmt.Errorf("Error inserting order into db")
	}

	order1.Index = 3
	order1.PanelSize = "260 1500 sq.ft homes set in a medieval cathedral style construction"
	order1.TotalValue = 40000
	order1.Location = "Trafalgar Square, London"
	order1.MoneyRaised = 0
	order1.Metadata = "Trafalgar Square is set in the heart of London's financial district, with big banks all over"
	order1.Live = false
	order1.INVAssetCode = ""
	order1.DEBAssetCode = ""
	order1.PBAssetCode = ""
	order1.DateInitiated = utils.Timestamp()
	order1.Years = 7
	order1.RecipientName = "Martin" // this is not the username of the recipient

	err = InsertOrderRPC(order1)
	if err != nil {
		return fmt.Errorf("Error inserting order into db")
	}

	var inv Investor
	allInvs, err := RetrieveAllInvestors()
	if err != nil {
		log.Fatal(err)
	}
	if len(allInvs) == 0 {
		inv.Index = 1
		inv.LoginUserName = "john"
		inv.LoginPassword = "e9a75486736a550af4fea861e2378305c4a555a05094dee1dca2f68afea49cc3a50e8de6ea131ea521311f4d6fb054a146e8282f8e35ff2e6368c1a62e909716"
		inv.Name = "John"
		err = InsertInvestor(inv)
		if err != nil {
			log.Fatal(err)
		}
	} else if len(allInvs) == 1 {
		// don't do anything
	}

	allRecs, err := RetrieveAllRecipients()
	if err != nil {
		log.Fatal(err)
	}
	if len(allRecs) == 0 {
		var rec Recipient
		rec.Index = 1
		rec.LoginUserName = "martin"
		rec.LoginPassword = "e9a75486736a550af4fea861e2378305c4a555a05094dee1dca2f68afea49cc3a50e8de6ea131ea521311f4d6fb054a146e8282f8e35ff2e6368c1a62e909716"
		rec.Name = "Martin"
		err = InsertRecipient(rec)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
