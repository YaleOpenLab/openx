package database

import (
	"log"
	"testing"
)

func TestDb(t *testing.T) {
	db, err := OpenDB()
	if err != nil {
		log.Fatal(err)
		// this means that we couldn't open the dtabase and we need to do something else
	}
	var dummy Order
	dummy.Index = 1
	dummy.PanelSize = "16 inches long, 36	inches wide"
	dummy.TotalValue = 14000
	dummy.Location = "Puerto Rico"
	dummy.MoneyRaised = 0
	dummy.Metadata = "This is a test entry and if present in the database, should be deleted"
	dummy.Live = true

	err = InsertOrder(dummy, db)
	if err != nil {
		t.Errorf("Inserting an order into the database failed")
		// shouldn't really fatal here, but this is in main, so we can't return
	}
	order, err := RetrieveOrder(dummy.Index, db)
	if err != nil {
		log.Println(err)
		t.Errorf("Retrieving order from the database failed")
		// again, shouldn't really fat a here, but we're in main
	}
	log.Println("Retrieved order: ", order)
	dummy.Index = 2
	err = InsertOrder(dummy, db)
	if err != nil {
		log.Println(err)
		t.Errorf("Inserting an order into the database failed")
		// shouldn't really fatal here, but this is in main, so we can't return
	}
	orders, err := RetrieveAllOrders(db)
	if err != nil {
		log.Println("Retrieve all error: ", err)
		t.Errorf("Failed in retrieving all orders")
	}
	log.Println("Retrieved orders: ", orders)
	err = DeleteOrder(dummy.Index, db)
	if err != nil {
		log.Println(err)
		t.Errorf("Deleting an  roder from the db failed")
	}
	log.Println("Deleted order")
	_, err = RetrieveOrder(dummy.Index, db)
	if err == nil {
		log.Println(err)
		// this should fail because we're trying to read an empty key value pair
		t.Errorf("Found deleted entry, quitting!")
	}
	db.Close() // needed since only order has two sets of functions, one for RPC
	// connections and the other for non RPC connections
	inv, err := NewInvestor("user1", "blah", "cool", true)
	if err != nil {
		log.Fatal(err)
	}
	err = InsertInvestor(inv)
	if err != nil {
		log.Fatal(err)
	}
	allInvestors, err := RetrieveAllInvestors()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Retrieved all investors: ", allInvestors)
	inv, err = NewInvestor("user1", "b921f75437050f0f7d2caba6303d165309614d524e3d7e6bccf313f39d113468d30e1e2ac01f91f6c9b66c083d393f49b3177345311849edb026bb86ee624be0", "cool", true)
	if err != nil {
		log.Fatal(err)
	}
	err = InsertInvestor(inv)
	if err != nil {
		log.Fatal(err)
	}
	_, err = SearchForPassword("b921f75437050f0f7d2caba6303d165309614d524e3d7e6bccf313f39d113468d30e1e2ac01f91f6c9b66c083d393f49b3177345311849edb026bb86ee624be0")
	if err != nil {
		log.Fatal(err)
	}
}
