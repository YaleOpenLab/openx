// +build all travis

package database

import (
	"log"
	"os"
	"testing"
)

// go test --tags="all" -coverprofile=test.txt .
func TestDb(t *testing.T) {
	var err error
	db, err := OpenDB()
	if err != nil {
		t.Fatal(err)
	}
	db.Close() // close immmediately after check
	dummy, err := NewOrder("16 inches long, 36	inches wide", 14000, "Puerto Rico", 0, "This is a test entry and if present in the database, should be deleted")
	if err != nil {
		t.Fatal(err)
	}

	err = InsertOrder(dummy)
	if err != nil {
		t.Errorf("Inserting an order into the database failed")
		// shouldn't really fatal here, but this is in main, so we can't return
	}
	order, err := RetrieveOrder(dummy.Index)
	if err != nil {
		log.Println(err)
		t.Errorf("Retrieving order from the database failed")
		// again, shouldn't really fat a here, but we're in main
	}
	klx1, err := RetrieveOrder(1000)
	if klx1.Index != 0 {
		t.Fatalf("REtrieved order which does not exist, quitting!")
	}
	log.Println("Retrieved order: ", order)
	if order.Index != dummy.Index {
		t.Fatalf("Indices don't match, quitting!")
	}
	dummy.Index = 2 // change index and try inserting another order
	err = InsertOrder(dummy)
	if err != nil {
		log.Println(err)
		t.Errorf("Inserting an order into the database failed")
		// shouldn't really fatal here, but this is in main, so we can't return
	}
	orders, err := RetrieveAllOrders()
	if err != nil {
		log.Println("Retrieve all error: ", err)
		t.Errorf("Failed in retrieving all orders")
	}
	if orders[0].Index != 1 {
		t.Fatalf("Index of first element doesn't match, quitting!")
	}
	log.Println("Retrieved orders: ", orders)
	oOrders, err := RetrieveAllOrders()
	if err != nil {
		log.Println("Retrieve all error: ", err)
		t.Errorf("Failed in retrieving all originated orders")
	}
	if len(oOrders) != 2 {
		t.Fatalf("Originated orders present!")
	}
	err = DeleteOrder(dummy.Index)
	if err != nil {
		log.Println(err)
		t.Errorf("Deleting an  roder from the db failed")
	}
	log.Println("Deleted order")
	// err = DeleteOrder(1000) this would work becuase the key will not be found and hence would not return an error
	_, err = RetrieveOrder(dummy.Index)
	if err == nil {
		log.Println(err)
		// this should fail because we're trying to read an empty key value pair
		t.Errorf("Found deleted entry, quitting!")
	}
	// now we repeat the same tests for all other entities
	// connections and the other for non RPC connections
	inv, err := NewInvestor("inv1", "blah", "cool")
	if err != nil {
		t.Fatal(err)
	}
	err = InsertInvestor(inv)
	if err != nil {
		t.Fatal(err)
	}
	rInv, err := RetrieveInvestor(1)
	if err != nil {
		t.Fatal(err)
	}
	ix1, err := RetrieveInvestor(1000)
	if ix1.U.Index != 0 {
		t.Fatalf("Invalid Investor exists")
	}
	if rInv.U.Name != inv.U.Name || rInv.U.LoginUserName != inv.U.LoginUserName || rInv.U.LoginPassword != inv.U.LoginPassword {
		t.Fatalf("Usernames don't match, quitting!")
	}
	inv, err = NewInvestor("inv2", "b921f75437050f0f7d2caba6303d165309614d524e3d7e6bccf313f39d113468d30e1e2ac01f91f6c9b66c083d393f49b3177345311849edb026bb86ee624be0", "cool")
	if err != nil {
		t.Fatal(err)
	}
	err = InsertInvestor(inv)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ValidateInvestor("inv2",
		"f28f9cc1f8c415d5c43c5bef02afb9493c0ed4236b876bff0c3d98d31e134b5505a1401604f301077b8a8f7c2a482afec428fe04a4b8ada6ad1337e10c6ebb99")
	if err != nil {
		t.Fatal(err)
	}
	allInvestors, err := RetrieveAllInvestors()
	if err != nil {
		t.Fatal(err)
	}
	if len(allInvestors) != 2 {
		t.Fatalf("Lengths of invesotrs don't match, quitting!")
	}
	if allInvestors[0].U.Name != "cool" || allInvestors[0].U.LoginUserName != "inv1" {
		t.Fatalf("UserNames don't match, quitting!")
	}
	// investors entity testing over, test recipients below in the same way
	// now we repeat the same tests for all other entities
	// connections and the other for non RPC connections
	recp, err := NewRecipient("user1", "blah", "cool")
	if err != nil {
		t.Fatal(err)
	}
	err = InsertRecipient(recp)
	if err != nil {
		t.Fatal(err)
	}
	rRecp, err := RetrieveRecipient(recp.U.Index)
	if err != nil {
		t.Fatal(err)
	}
	rx1, err := RetrieveRecipient(1000)
	if rx1.U.Index != 0 {
		t.Fatalf("Invalid Recipient exists")
	}
	if rRecp.U.Name != recp.U.Name || rRecp.U.LoginUserName != recp.U.LoginUserName || rRecp.U.LoginPassword != recp.U.LoginPassword {
		t.Fatalf("Usernames don't match, quitting!")
	}
	recp, err = NewRecipient("user1", "b921f75437050f0f7d2caba6303d165309614d524e3d7e6bccf313f39d113468d30e1e2ac01f91f6c9b66c083d393f49b3177345311849edb026bb86ee624be0", "cool")
	if err != nil {
		t.Fatal(err)
	}
	err = InsertRecipient(recp)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ValidateRecipient("user1",
		"f28f9cc1f8c415d5c43c5bef02afb9493c0ed4236b876bff0c3d98d31e134b5505a1401604f301077b8a8f7c2a482afec428fe04a4b8ada6ad1337e10c6ebb99")
	if err != nil {
		t.Fatal(err)
	}
	allRecipients, err := RetrieveAllRecipients()
	if err != nil {
		t.Fatal(err)
	}
	if len(allRecipients) != 2 {
		t.Fatalf("Lengths of recipients don't match, quitting!")
	}
	if allRecipients[0].U.Name != "cool" || allRecipients[0].U.LoginUserName != "user1" {
		t.Fatalf("UserNames don't match, quitting!")
	}
	// tests for originators
	newCE, err := NewContractEntity("OrigTest", "pwd", "NameOrigTest", "123 ABC Street", "OrigDescription", "originator")
	if err != nil {
		t.Fatal(err)
	}
	err = InsertContractEntity(newCE)
	if err != nil {
		t.Fatal(err)
	}
	allOrigs, err := RetrieveAllContractEntities("originator")
	if err != nil {
		t.Fatal(err)
	}
	acz1, err := RetrieveAllContractEntities("random")
	if len(acz1) != 0 {
		log.Println(acz1)
		t.Fatalf("Category which does not exist exists?")
	}
	if len(allOrigs) != 1 || allOrigs[0].U.Name != "NameOrigTest" {
		t.Fatal("Names don't match, quitting!")
	}
	// tests for contractors
	newCE, err = NewContractor("ConTest", "pwd", "NameConTest", "123 ABC Street", "ConDescription") // use and test this as well
	if err != nil {
		t.Fatal(err)
	}
	err = InsertContractEntity(newCE)
	if err != nil {
		t.Fatal(err)
	}
	allConts, err := RetrieveAllContractEntities("contractor")
	if err != nil {
		t.Fatal(err)
	}
	if len(allConts) != 1 || allConts[0].U.Name != "NameConTest" {
		t.Fatal("Names don't match, quitting!")
	}
	_, err = newCE.ProposeContract("100 16x32 panels", 28000, "Puerto Rico", 6, "LEED+ Gold rated panels and this is random data out of nowhere and we supply our own devs and provide insurance guarantee as well. Dual audit maintenance upto 1 year. Returns capped as per defaults", 1, 1)
	// 1 for retrieving martin as the recipient and 1 is the order Index
	if err != nil {
		log.Fatal(err)
	}
	// now try retrieving this proposed contract
	rO, err := RetrieveOrder(1)
	if err != nil {
		t.Fatal(err)
	}
	_, allPCs, err := RetrieveAllProposedContracts(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(allPCs) != 1 { // add check for stuff here
		log.Println("LEN ASDASD", len(allPCs))
	}
	rPC, err := FindInKey(1, allPCs)
	if err != nil {
		t.Fatal(err)
	}
	if rPC.O.Index != rO.Index {
		t.Fatal("Indices don't match")
	}

	// now come the failure cases which should fail and we shall catch the case when they don't
	_, allPCs, err = RetrieveAllProposedContracts(2)
	if len(allPCs) != 0 {
		log.Println("LEBNGRG: ", len(allPCs))
		t.Fatalf("Retrieving a missing contract succeeds, quitting!")
	}
	rPC, err = FindInKey(2, allPCs)
	if err == nil {
		t.Fatalf("Entity which should be missing exists!")
	}

	// Checks for users
	u1, err := NewUser("uname1", "pwd1", "NewUser")
	if err != nil {
		t.Fatal(err)
	}
	err = u1.GenKeys()
	if err != nil {
		t.Fatal(err)
	}
	rU1, err := RetrieveUser(u1.Index)
	if err != nil {
		t.Fatal(err)
	}
	ux1, err := RetrieveUser(1000)
	if ux1.Index != 0 {
		t.Fatalf("Retrieving invalid user succeeds, quitting!")
	}
	if rU1.Name != "NewUser" {
		t.Fatalf("User Names don't match")
	}
	vU, err := ValidateUser("uname1", "01e2b5c87c9ecbad25df63f47098bd1b47a9ffbdce46967a0cca8922924fc009e8f7a41c45109b47982ee7b9b36f617cab3f1d029b17ff110fffa91ed4ab4d27")
	if err != nil {
		t.Fatal(err)
	}
	if vU.Name != "NewUser" {
		t.Fatalf("User Names don't match")
	}
	allUs, err := RetrieveAllUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(allUs) != 7 {
		t.Fatal("Length of all users doesn't match")
	}
	// newCE, err = NewContractEntity("ConTest", "pwd", "NameConTest", "123 ABC Street", "ConDescription", "contractor")
	rCE, err := SearchForContractEntity("ConTest", "9f88a8d40b90616715f868ed195d24e5df994f56bce34eddb022c213484eb0f220d8907e4ecd8f64ddd364cb30bb5758b32ee26541f340b930f7e5bf756907a4")
	if err != nil {
		t.Fatal(err)
	}
	if rCE.Contractor != true {
		log.Println("THe reole: ", rCE.Contractor, rCE)
		t.Fatal("Roles don't match!")
	}
	trC1, err := RetrieveContractEntity(6)
	if err != nil || trC1.U.Index == 0 {
		t.Fatal("Contract Entity lookup failed")
	}
	newCE2, err := NewOriginator("OrigTest2", "pwd", "NameOrigTest2", "123 ABC Street", "OrigDescription2")
	if err != nil {
		t.Fatal(err)
	}
	err = InsertContractEntity(newCE2)
	if err != nil {
		t.Fatal(err)
	}
	_, err = newCE2.OriginContract("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", 1) // 1 is the index for martin
	if err != nil {
		log.Fatal(err)
	}
	allOOs, err := RetrieveAllOriginatedOrders()
	if err != nil {
		t.Fatal(err)
	}
	if len(allOOs) != 1 {
		log.Println("Length of all Originated Orders: ", len(allOOs))
		t.Fatalf("Length of all orignated orders doesn't match")
	}
	err = DeleteKeyFromBucket(recp.U.Index, RecipientBucket)
	if err != nil {
		t.Fatal(err)
	}
	// can't test the payback stuff since we need the recipient to have funds and stuff
	// maybe test the investor sending funds and stuff as well? but that's imported
	// by main and we could have problems here related to that
	os.Remove("yol.db")
}
