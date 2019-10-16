// +build all travis

package database

import (
	"log"
	"os"
	"testing"
	"time"

	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/Varunram/essentials/xlm"
	assets "github.com/Varunram/essentials/xlm/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	build "github.com/stellar/go/txnbuild"
)

// go test --tags="all" -coverprofile=test.txt .
func TestDb(t *testing.T) {
	var err error
	consts.SetConsts(false)
	os.Remove("blahopenx.db")
	consts.DbDir = "blah"   // set to a false db so that we can test errors arising from OpenDB()
	CreateHomeDir()         // create home directory if it doesn't exist yet
	os.Remove(consts.DbDir) // remove the test database file, if it exists
	db, err := OpenDB()
	if err != nil {
		t.Fatal(err)
	}
	db.Close() // close immmediately after check

	user, err := NewUser("user1", "blah", "blah", "User1")
	if err != nil {
		t.Fatal(err)
	}
	user.AccessToken = "ACCESSTOKEN"
	user.AccessTokenTimeout = utils.Unix() + 1000000
	err = user.Save()
	if err != nil {
		t.Fatal(err)
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
	if len(allUsers) != 1 {
		t.Fatalf("Unknown users existing, quitting!")
	}

	_, err = ValidateAccessToken("user1", "ACCESSTOKEN")
	if err != nil {
		log.Println(err)
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateAccessToken("blah", "ACCESSTOKEN")
	if err == nil {
		log.Println(err)
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	err = user.GenKeys("blah")
	if err != nil {
		t.Fatalf("Not able to generate keys, quitting!")
	}

	_, err = ValidateSeedpwd(user.Username, user.Pwhash, "blah")
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateSeedpwd("shouldfail", user.Pwhash, "blah")
	if err == nil {
		t.Fatalf("can't catch wrong username")
	}

	_, err = ValidateSeedpwd(user.Username, user.Pwhash, "shouldfail")
	if err == nil {
		t.Fatalf("can't catch wrong seedpwd")
	}

	err = user.GenKeys("algorand", "algorand")
	if err == nil {
		t.Fatalf("able to generate algorand keys even when daemon is not running")
	}

	err = xlm.GetXLM(user.StellarWallet.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	err = xlm.GetXLM(user.SecondaryWallet.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)
	err = user.IncreaseTrustLimit("blah", 10)
	if err != nil {
		t.Fatal(err)
	}
	_ = build.CreditAsset{"blah", user.StellarWallet.PublicKey}
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
	user.Admin = true
	err = user.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = user.SetBan(1000)
	if err != nil {
		t.Fatalf("Not able to set ban on legitimate user, quitting")
	}
	err = user.SetBan(1000)
	if err != nil {
		t.Fatalf("Able to set ban on user even if ban is already set, quitting")
	}
	err = user.Authorize(user.Index)
	if err != nil {
		t.Fatalf("Not able to set kyc flag, exiting!")
	}
	_, err = RetrieveAllUsersWithKyc()
	if err != nil {
		t.Fatal(err)
	}
	testuser, _ := SearchWithEmailId("blahx@blah.com")
	if testuser.Index != 0 {
		t.Fatalf("user with invalid email exists")
	}
	err = user.MoveFundsFromSecondaryWallet(10, "blah")
	if err != nil {
		t.Fatal(err)
	}
	err = user.MoveFundsFromSecondaryWallet(10, "shouldfail")
	if err == nil {
		t.Fatalf("decryption succeeds with invalid seedpwd for secondary account")
	}
	err = user.MoveFundsFromSecondaryWallet(100000, "shouldfail")
	if err == nil {
		t.Fatalf("can transfer more amount than possessed")
	}
	err = user.MoveFundsFromSecondaryWallet(-1, "blah")
	if err == nil {
		t.Fatalf("not able to catch invalid amount error")
	}
	err = user.SweepSecondaryWallet("blah")
	if err != nil {
		t.Fatal(err)
	}
	err = user.SweepSecondaryWallet("invalidseedpwd")
	if err == nil {
		t.Fatalf("no able to catch invalid seedpwd")
	}
	_, err = TopReputationUsers()
	if err != nil {
		t.Fatal(err)
	}
	token, err := user.GenAccessToken()
	if err != nil {
		t.Fatal(err)
	}
	_, err = ValidateSeedpwdAuthToken(user.Username, token, "blah")
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateSeedpwdAuthToken("fakeusername", token, "blah")
	if err == nil {
		t.Fatalf("not able to detect fake username")
	}

	_, err = ValidateSeedpwdAuthToken(user.Username, "fakeaccesstoken", "blah")
	if err == nil {
		t.Fatalf("not able to detect fake access token")
	}

	_, err = ValidateSeedpwdAuthToken(user.Username, token, "fakeseedpwd")
	if err == nil {
		t.Fatalf("not able to detect fake seedpwd")
	}

	err = user.AddtoMailbox("test", "test")
	if err != nil {
		t.Fatal(err)
	}
	var temp []byte
	err = user.ImportSeed(temp, "", "")
	if err == nil {
		t.Fatalf("able to decrypt empty byte array")
	}
	err = user.ImportSeed(user.StellarWallet.EncryptedSeed, user.StellarWallet.PublicKey, "blah")
	if err != nil {
		t.Fatal(err)
	}
	_, err = user.Authenticate2FA("dummypassword")
	if err == nil {
		t.Fatalf("able to authenticate empty 2fa password")
	}
	_, err = user.Generate2FA()
	if err != nil {
		t.Fatal(err)
	}
	err = user.GiveFeedback(556, 5)
	if err == nil {
		t.Fatalf("able to give feedback to a non existent user")
	}
	err = user.GiveFeedback(user.Index, 10)
	if err == nil {
		t.Fatalf("able to give more feedback than 5")
	}

	err = user.GiveFeedback(user.Index, 5)
	if err != nil {
		t.Fatal(err)
	}
	err = user.AddEmail("ghost@ghosts.com")
	if err != nil {
		t.Fatal(err)
	}
	_, err = CheckUsernameCollision(user.Username)
	if err == nil {
		t.Fatalf("can't catch username collision")
	}
	user.Admin = true
	err = user.Save()
	if err != nil {
		t.Fatal(err)
	}
	_, err = RetrieveAllAdmins()
	if err != nil {
		t.Fatal(err)
	}
	err = DeleteKeyFromBucket(user.Index, UserBucket)
	if err != nil {
		t.Fatal(err)
	}

	// end of user related tests

	err = NewPlatform("platform", "CODE", false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = RetrievePlatform(1) // GUESS?
	if err != nil {
		t.Fatal(err)
	}

	_, err = RetrieveAllPlatforms()
	if err != nil {
		t.Fatal(err)
	}

	_, err = RetrieveAllPfLim()
	if err != nil {
		t.Fatal(err)
	}

	os.Remove(consts.DbDir + "/openx.db")
}
