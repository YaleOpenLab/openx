package orders

import(
  "testing"
  "log"
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
    // shouldn't realyl fatal here, but htis is in main, so we can't return
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
    // shouldn't realyl fatal here, but htis is in main, so we can't return
  }
  orders, err := RetrieveAll(db)
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
  order, err = RetrieveOrder(dummy.Index, db)
  if err == nil {
    log.Println(err)
    // this should fail because we're trying to read an empty key value pair
    t.Errorf("Found deleted entry, quitting!")
  }
}
