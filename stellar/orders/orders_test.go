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
  dummy.Index = 68
  dummy.PanelSize = "16 inches long, 36	inches wide"
  dummy.TotalValue = 14000
  dummy.Location = "Puerto Rico"
  dummy.MoneyRaised = 0
  dummy.Metadata = "This is our first project and would be nice if you could contribute"

  err = InsertOrder(dummy, db)
  if err != nil {
    t.Errorf("Inserting an order into the database failed")
    // shouldn't realyl fatal here, but htis is in main, so we can't return
  }
  order, err := RetrieveOrder(dummy.Index, db)
  if err != nil {
    t.Errorf("Retrieving order from the database failed")
    // again, shouldn't really fat a here, but we're in main
  }
  log.Println("Retrieved order: ", order)
  err = DeleteOrder(dummy.Index, db)
  if err != nil {
    t.Errorf("Deleting an  roder from the db failed")
  }
  log.Println("Deleted order")
  order, err = RetrieveOrder(dummy.Index, db)
  if err == nil {
    // this should fail because we're trying to read an empty key value pair
    t.Errorf("Found deleted entry, quitting!")
  }
}
