// +build all travis

package database

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	xlm "github.com/Varunram/essentials/crypto/xlm"
	assets "github.com/Varunram/essentials/crypto/xlm/assets"
	wallet "github.com/Varunram/essentials/crypto/xlm/wallet"
	consts "github.com/YaleOpenLab/openx/consts"
	build "github.com/stellar/go/txnbuild"
)

// go test --tags="all" -coverprofile=test.txt .
func TestDb(t *testing.T) {
	var err error
	consts.SetConsts()
	CreateHomeDir()         // create home directory if it doesn't exist yet
	consts.DbDir = "blah"   // set to a false db so that we can test errors arising from OpenDB()
	os.Remove(consts.DbDir) // remove the test database file, if it exists
	db, err := OpenDB()
	if err != nil {
		t.Fatal(err)
	}
	db.Close() // close immmediately after check
	inv, err := NewInvestor("investor1", "blah", "blah", "Investor1")
	if err != nil {
		t.Fatal(err)
	}
	xc, err := RetrieveAllInvestors()
	if len(xc) != 1 {
		t.Fatal(err)
	}

	// try retrieving existing stuff
	inv1, err := RetrieveInvestor(1)
	if err != nil {
		t.Fatal(err)
	}
	if inv1.U.Name != "Investor1" {
		t.Fatalf("Investor names don't match, quitting!")
	}
	// func NewRecipient(uname string, pwd string, seedpwd string, Name string) (Recipient, error) {
	recp, err := NewRecipient("recipient1", "blah", "blah", "Recipient1")
	if err != nil {
		t.Fatal(err)
	}

	rec1, err := RetrieveRecipient(recp.U.Index)
	if err != nil {
		t.Fatal(err)
	}

	if rec1.U.Name != "Recipient1" {
		t.Fatalf("Recipient usernames don't match. quitting!")
	}

	allRec, err := RetrieveAllRecipients()
	if err != nil || len(allRec) != 1 {
		log.Println("length of all recipients not 1")
		t.Fatal(err)
	}

	user, err := NewUser("user1", "blah", "blah", "User1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = RetrieveInvestor(1000)
	if err == nil {
		t.Fatalf("Investor shouldn't exist, but does, quitting!")
	}

	_, err = RetrieveRecipient(1000)
	if err == nil {
		t.Fatalf("Recipient shouldn't exist, but does, quitting!")
	}

	user1, err := RetrieveUser(user.Index)
	if err != nil {
		t.Fatal(err)
	}
	if user1.Name != "User1" {
		t.Fatalf("Usernames don't match. quitting!")
	}

	tmpuser, _ := RetrieveUser(1000)
	if tmpuser.Index != 0 {
		t.Fatalf("Investor shouldn't exist, but does, quitting!")
	}

	allUsers, err := RetrieveAllUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(allUsers) != 3 {
		t.Fatalf("Unknown users existing, quitting!")
	}

	// check if each of the validate functions work
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

	// check voting balance routes
	voteBalance := inv.VotingBalance
	err = inv.ChangeVotingBalance(10000)
	if err != nil {
		t.Fatal(err)
	}
	if inv.VotingBalance-voteBalance != 10000 {
		t.Fatalf("Voting Balance not added, quitting!")
	}
	err = inv.ChangeVotingBalance(-10000)
	if err != nil {
		t.Fatal(err)
	}
	if inv.VotingBalance-voteBalance != 0 {
		t.Fatalf("Voting Balance not added, quitting!")
	}

	// check CanInvest Route
	if inv.CanInvest(1000) {
		t.Fatalf("CanInvest Returns true!")
	}

	err = user.GenKeys("blah")
	if err != nil {
		t.Fatalf("Not able to generate keys, quitting!")
	}

	err = DeleteKeyFromBucket(user.Index, UserBucket)
	if err != nil {
		t.Fatal(err)
	}

	// check the asset functions below. For some weird reason, placing these tests
	// above confuses the other routes, so placing everything here so that we can
	// isolate them from the other routes.
	err = xlm.GetXLM(recp.U.StellarWallet.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	err = xlm.GetXLM(inv.U.StellarWallet.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	err = inv.U.IncreaseTrustLimit("blah", 10)
	if err != nil {
		t.Fatal(err)
	}
	_ = build.CreditAsset{"blah", recp.U.StellarWallet.PublicKey}
	invSeed, err := wallet.DecryptSeed(inv.U.StellarWallet.EncryptedSeed, "blah")
	if err != nil {
		t.Fatal(err)
	}
	hash, err := assets.TrustAsset("blah", recp.U.StellarWallet.PublicKey, 100, invSeed)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("HASH IS: ", hash)
	_, err = assets.TrustAsset("blah", recp.U.StellarWallet.PublicKey, -1, "blah")
	if err == nil {
		t.Fatalf("can trust asset with invalid s eed!")
	}
	pkSeed, _, err := xlm.GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	_ = build.CreditAsset{"blah2", pkSeed} // this account doesn't exist yet, so this should fail
	_, err = assets.TrustAsset("blah2", "", -1, "blah")
	if err == nil {
		t.Fatalf("can trust invalid asset")
	}
	_, err = RetrieveAllUsersWithoutKyc()
	if err != nil {
		t.Fatal(err)
	}
	_, err = RetrieveAllUsersWithKyc()
	if err != nil {
		t.Fatal(err)
	}
	err = user.ChangeReputation(1.0)
	if err != nil {
		t.Fatal(err)
	}
	err = user.ChangeReputation(1.0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = TopReputationUsers()
	if err != nil {
		t.Fatal(err)
	}
	err = user.Authorize(user.Index)
	if err == nil {
		t.Fatalf("Not able to catch inspector permission error")
	}
	err = user.SetBan(100)
	if err == nil {
		t.Fatalf("able to ban a user even though person is not an inspector, quitting")
	}
	err = AddInspector(user.Index)
	if err != nil {
		t.Fatal(err)
	}
	user, err = RetrieveUser(user.Index)
	if err != nil {
		t.Fatal(err)
	}
	err = user.Authorize(user.Index)
	if err != nil {
		t.Fatal(err)
	}
	err = user.SetBan(user.Index)
	if err == nil {
		t.Fatalf("able to  set a ban on self, quitting")
	}
	err = user.SetBan(-1)
	if err == nil {
		t.Fatalf("able to set ban on user who doesn't exist, quitting")
	}
	var banTest User
	banTest.Index = 1000
	err = banTest.Save()
	if err != nil {
		t.Fatalf("not able to save user for banning, quitting")
	}
	err = user.SetBan(1000)
	if err != nil {
		log.Println(err)
		t.Fatalf("Not able to set ban on legitimate user, quitting")
	}
	err = user.SetBan(1000)
	if err != nil {
		t.Fatalf("Able to set ban on user even if ban is already set, quitting")
	}
	err = user.Authorize(user.Index)
	if err == nil {
		t.Fatalf("Able to authorize KYC'd user, exiting!")
	}
	_, err = RetrieveAllUsersWithKyc()
	if err != nil {
		t.Fatal(err)
	}
	_, err = TopReputationRecipients()
	if err != nil {
		t.Fatal(err)
	}
	err = recp.U.ChangeReputation(1.0)
	if err != nil {
		t.Fatal(err)
	}
	err = recp.U.ChangeReputation(-1.0)
	if err != nil {
		t.Fatal(err)
	}
	err = recp.U.AddEmail("blah@blah.com")
	if err != nil {
		t.Fatal(err)
	}
	_, err = SearchWithEmailId("blah@blah.com")
	if err != nil {
		t.Fatal(err)
	}
	testuser, _ := SearchWithEmailId("blahx@blah.com")
	if testuser.StellarWallet.PublicKey != "" {
		t.Fatalf("user with invalid email exists")
	}
	err = xlm.GetXLM(inv.U.SecondaryWallet.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	err = inv.U.MoveFundsFromSecondaryWallet(10, "blah")
	if err != nil {
		t.Fatal(err)
	}
	err = inv.U.MoveFundsFromSecondaryWallet(-1, "blah")
	if err == nil {
		t.Fatalf("not able to catch invalid amount error")
	}
	err = inv.U.SweepSecondaryWallet("blah")
	if err != nil {
		t.Fatal(err)
	}
	err = inv.U.SweepSecondaryWallet("invalidseedpwd")
	if err == nil {
		t.Fatalf("no able to catch invalid seedpwd")
	}
	_, err = TopReputationInvestors()
	if err != nil {
		t.Fatal(err)
	}
	err = inv.U.ChangeReputation(1.0)
	if err != nil {
		t.Fatal(err)
	}
	err = inv.U.ChangeReputation(-1.0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = TopReputationRecipients()
	if err != nil {
		t.Fatal(err)
	}
	_, err = TopReputationInvestors()
	if err != nil {
		t.Fatal(err)
	}
	_, err = TopReputationUsers()
	if err != nil {
		t.Fatal(err)
	}
	var blah Investor
	var blahuser User
	blah.U = &blahuser
	blah.U.Name = "Cool"
	blahBytes, err := json.Marshal(blah)
	if err != nil {
		t.Fatal(err)
	}

	var uBlah Investor
	err = json.Unmarshal(blahBytes, &uBlah)
	if err != nil {
		t.Fatal(err)
	}

	os.Remove(consts.DbDir + "/openx.db")
}
