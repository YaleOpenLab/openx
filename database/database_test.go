// +build all travis

package database

import (
	"log"
	"os"
	"testing"
)

// TODO: rewrite how this works and split between platforms and database
// go test --tags="all" -coverprofile=test.txt .
func TestDb(t *testing.T) {
	var err error
	CreateHomeDir()
	os.Remove(os.Getenv("HOME") + "/.openfinancing/database/" + "/yol.db")
	err = os.MkdirAll(os.Getenv("HOME")+"/.openfinancing/database", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	db, err := OpenDB()
	if err != nil {
		t.Fatal(err)
	}
	db.Close() // close immmediately after check

	inv, err := NewInvestor("investor1", "blah", "blah", "Investor1")
	if err != nil {
		t.Fatal(err)
	}
	err = inv.Save()
	if err != nil {
		t.Fatal(err)
	}
	log.Println("INDEX IS: ", inv.U.Index)
	// func NewRecipient(uname string, pwd string, seedpwd string, Name string) (Recipient, error) {
	recp, err := NewRecipient("recipient1", "blah", "blah", "Recipient1")
	if err != nil {
		t.Fatal(err)
	}
	err = recp.Save()
	if err != nil {
		t.Fatal(err)
	}
	// func NewUser(uname string, pwd string, seedpwd string, Name string) (User, error) {
	user, err := NewUser("user1", "blah", "blah", "User1")
	if err != nil {
		t.Fatal(err)
	}
	err = user.Save()
	if err != nil {
		t.Fatal(err)
	}

	// try retrieving existing stuff
	inv1, err := RetrieveInvestor(inv.U.Index)
	if err != nil {
		t.Fatal(err)
	}
	if inv1.U.Name != "Investor1" {
		t.Fatalf("Usernames don't match. quitting!")
	}

	tmpinv, err := RetrieveInvestor(1000)
	if tmpinv.U.Index != 0 {
		t.Fatalf("Investor shouldn't exist, but does, quitting!")
	}

	rec1, err := RetrieveRecipient(recp.U.Index)
	if err != nil {
		t.Fatal(err)
	}
	if rec1.U.Name != "Recipient1" {
		t.Fatalf("Usernames don't match. quitting!")
	}

	tmprec, err := RetrieveRecipient(1000)
	if tmprec.U.Index != 0 {
		t.Fatalf("Investor shouldn't exist, but does, quitting!")
	}

	user1, err := RetrieveUser(user.Index)
	if err != nil {
		t.Fatal(err)
	}
	if user1.Name != "User1" {
		t.Fatalf("Usernames don't match. quitting!")
	}

	tmpuser, err := RetrieveUser(1000)
	if tmpuser.Index != 0 {
		t.Fatalf("Investor shouldn't exist, but does, quitting!")
	}

	// test length of users in bucket
	allInv, err :=  RetrieveAllInvestors()
	if err != nil {
		t.Fatal(err)
	}
	if len(allInv) != 1 {
		t.Fatalf("Unknown investors existing, quitting!")
	}

	allRec, err :=  RetrieveAllRecipients()
	if err != nil {
		t.Fatal(err)
	}
	if len(allRec) != 1 {
		t.Fatalf("Unknown recipients existing, quitting!")
	}

	allUser, err := RetrieveAllUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(allUser) != 3 {
		t.Fatalf("Unknown users existing, quitting!")
	}

	_, err = ValidateInvestor("investor1", "ed2df20bb16ecb0b4b149cf8e7d9819afd608b22999e707364196187fca0cf38544c9f3eb981ad81cef18562e4c818370eab068992639af7d70488945265197f")
	if err != nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateRecipient("recipient1", "ed2df20bb16ecb0b4b149cf8e7d9819afd608b22999e707364196187fca0cf38544c9f3eb981ad81cef18562e4c818370eab068992639af7d70488945265197f")
	if err != nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateUser("user1", "ed2df20bb16ecb0b4b149cf8e7d9819afd608b22999e707364196187fca0cf38544c9f3eb981ad81cef18562e4c818370eab068992639af7d70488945265197f")
	if err != nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateInvestor("blah", "ed2df20bb16ecb0b4b149cf8e7d9819afd608b22999e707364196187fca0cf38544c9f3eb981ad81cef18562e4c818370eab068992639af7d70488945265197f")
	if err == nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateRecipient("blah", "ed2df20bb16ecb0b4b149cf8e7d9819afd608b22999e707364196187fca0cf38544c9f3eb981ad81cef18562e4c818370eab068992639af7d70488945265197f")
	if err == nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateUser("blah", "ed2df20bb16ecb0b4b149cf8e7d9819afd608b22999e707364196187fca0cf38544c9f3eb981ad81cef18562e4c818370eab068992639af7d70488945265197f")
	if err == nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	voteBalance := inv.VotingBalance
	err = inv.AddVotingBalance(10000)
	if err != nil {
		t.Fatal(err)
	}
	if inv.VotingBalance - voteBalance != 10000 {
		t.Fatalf("Voting Balance not added, quitting!")
	}
	err = inv.DeductVotingBalance(10000)
	if err != nil {
		t.Fatal(err)
	}
	if inv.VotingBalance - voteBalance != 0 {
		t.Fatalf("Voting Balance not added, quitting!")
	}
	// func (a *Investor) CanInvest(balance string, targetBalance string) bool {
	if inv.CanInvest("100", "1000") {
		t.Fatalf("CanInvest Returns true!")
	}

	err = user.GenKeys("blah")
	if err != nil {
		t.Fatalf("Not able to generate keys, quitting!")
	}

	// func (a *User) GetSeed(seedpwd string) (string, error) {
	_, err = user.GetSeed("blah")
	if err != nil {
		t.Fatal(err)
	}

	err = DeleteKeyFromBucket(user.Index, UserBucket)
	if err != nil {
		t.Fatal(err)
	}
}
