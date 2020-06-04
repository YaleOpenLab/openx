// +build all travis

package database

import (
	"log"
	"os"
	"sync"
	"testing"

	"github.com/Varunram/essentials/utils"
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
	consts.DbDir = "blah" // set to a false db so that we can test errors arising from OpenDB()
	xlm.SetConsts(10, false)
	username := "testusername"
	userpwhash := utils.SHA3hash("testpass")
	seedpwd := "x"
	email := "User1"

	PlatformBucket = []byte("FakePlatforms")
	UserBucket = []byte("FakeUsers")
	err = NewPlatform("platform", "CODE", false)
	if err == nil {
		t.Fatalf("unable to catch wrong platform bucket name")
	}
	_, err = RetrievePlatform(1) // GUESS?
	if err == nil {
		t.Fatalf("unable to catch wrong platform bucket name")
	}
	_, err = RetrieveAllPlatforms()
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}
	user, err := NewUser(username, userpwhash, seedpwd, email)
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}
	_, err = CheckUsernameCollision(user.Username)
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}
	_, err = SearchWithEmailID("fakeemail")
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}
	_, err = RetrieveAllUsers()
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}
	err = AddInspector(1)
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}
	err = user.SetBan(user.Index)
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}
	_, err = RetrieveAllUsersWithoutKyc()
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}
	_, err = RetrieveAllUsersWithKyc()
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}
	_, err = TopReputationUsers()
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}
	_, err = RetrieveAllAdmins()
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}
	_, err = ValidatePwhash("fakename", "fakepwhash")
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}
	_, err = ValidatePwhashReg("username", "pwhash")
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}
	_, err = ValidateAccessToken("fakeusername", "fakeaccessToken")
	if err == nil {
		t.Fatalf("failed access token length check")
	}
	_, err = ValidateAccessToken("fakeusername", utils.GetRandomString((32)))
	if err == nil {
		t.Fatalf("unable to catch wrong user bucket name")
	}

	PlatformBucket = []byte("Platforms")
	UserBucket = []byte("Users")

	consts.StablecoinPublicKey = "GAVEVWKMXVQ2WSCBTR7M5UKRVFFWIA52VP7ISDKZSEJKQS2VYG4D6C6P"
	consts.PlatformPublicKey = "GAJJMQAP5KG7GVCOVY2NUUJCVFX72GXZKMUQUCWUGN55EKFS3MXFAMEZ"
	CreateHomeDir()         // create home directory if it doesn't exist yet
	os.Remove(consts.DbDir) // remove the test database file, if it exists
	db, err := OpenDB()
	if err != nil {
		t.Fatal(err)
	}
	db.Close() // close immmediately after check

	user, err = NewUser(username, userpwhash, seedpwd, email)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("UINDEX: ", user.Index)
	err = user.VerReq()
	if err != nil {
		t.Fatal(err)
	}

	err = user.UnverReq()
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewUser(username, userpwhash, seedpwd, email)
	if err == nil {
		t.Fatalf("can create user with same username")
	}

	_, err = ValidatePwhashReg(username, userpwhash)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidatePwhashReg("fakeusername", userpwhash)
	if err == nil {
		t.Fatalf("unable to catch fake username in validatepwhashreg")
	}

	accessToken, err := user.GenAccessToken()
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateAccessToken(username, accessToken)
	if err == nil {
		t.Fatalf("can validate access token without user being confirmed")
	}

	_, err = ValidatePwhash(username, accessToken)
	if err == nil {
		t.Fatalf("can validate pwhash without user being confirmed")
	}

	_, err = CheckUsernameCollision(user.Username)
	if err == nil {
		t.Fatalf("can't catch username collision")
	}

	user.Conf = true
	err = user.Save()
	if err != nil {
		t.Fatal(err)
	}

	user, err = RetrieveUser(user.Index)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateAccessToken(username, accessToken)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateAccessToken(username, "faketoken")
	if err == nil {
		log.Println(err)
		t.Fatalf("didn't fail on fake token")
	}

	_, err = ValidateAccessToken(username, "faketokenfaketokenfaketokenfaketokenfaketoken")
	if err == nil {
		log.Println(err)
		t.Fatalf("didn't fail on fake token")
	}

	_, err = ValidateAccessToken("fakeusername", accessToken)
	if err == nil {
		log.Println(err)
		t.Fatalf("didn't fail on fake username")
	}

	_, err = ValidateAccessToken("fakeusername", "faketoken")
	if err == nil {
		log.Println(err)
		t.Fatalf("didn't fail on fake username and token")
	}

	user.AccessToken[accessToken] = 0
	err = user.Save()
	if err != nil {
		t.Fatal(err)
	}

	user, err = RetrieveUser(user.Index)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateAccessToken(username, accessToken)
	if err == nil {
		log.Println(err)
		t.Fatalf("didn't error out on token timeout")
	}

	user1000, _ := RetrieveUser(1000)
	if user1000.Index != 0 {
		t.Fatalf("User shouldn't exist, but does, quitting!")
	}

	allUsers, err := RetrieveAllUsers()
	if err != nil {
		t.Fatal(err)
	}

	if len(allUsers) != 1 {
		log.Println(len(allUsers))
		t.Fatalf("Unknown users exist, quitting!")
	}

	err = user.GenKeys(seedpwd)
	if err != nil {
		t.Fatalf("unable to generate keys, quitting!")
	}

	_, err = ValidateSeedpwd(username, userpwhash, seedpwd)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateSeedpwd("fakeusername", userpwhash, seedpwd)
	if err == nil {
		t.Fatalf("can't catch fake username")
	}

	_, err = ValidateSeedpwd(user.Username, user.Pwhash, "fakeseedpwd")
	if err == nil {
		t.Fatalf("can't catch fake seedpwd")
	}

	fakeEmailUser, _ := SearchWithEmailID("fakeemail")
	if fakeEmailUser.Index != 0 {
		t.Fatalf("user with invalid email exists")
	}

	err = user.AddEmail(email)
	if err != nil {
		t.Fatal(err)
	}

	_, err = SearchWithEmailID(email)
	if err != nil {
		t.Fatal(err)
	}

	accessToken, err = user.GenAccessToken()
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateSeedpwdAuthToken(user.Username, accessToken, seedpwd)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateSeedpwdAuthToken("fakeusername", accessToken, seedpwd)
	if err == nil {
		t.Fatalf("unable to detect fake username")
	}

	_, err = ValidateSeedpwdAuthToken(user.Username, "fakeaccesstoken", seedpwd)
	if err == nil {
		t.Fatalf("unable to detect fake access token")
	}

	_, err = ValidateSeedpwdAuthToken(user.Username, accessToken, "fakeseedpwd")
	if err == nil {
		t.Fatalf("unable to detect fake seedpwd")
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
		t.Fatalf("unable to catch inspector permission error")
	}
	err = user.SetBan(100)
	if err == nil {
		t.Fatalf("able to ban a user even though person is not an inspector, quitting")
	}
	err = user.SetBan(user.Index)
	if err == nil {
		t.Fatalf("able to ban own self")
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
	err = user.Authorize(1000)
	if err == nil {
		t.Fatalf("able to authorize non existent user")
	}
	err = user.Authorize(user.Index)
	if err == nil {
		t.Fatalf("unable to catch error when user has already KYC'd")
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
		t.Fatalf("unable to save user for banning, quitting")
	}
	user.Admin = true
	err = user.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = user.SetBan(1000)
	if err != nil {
		t.Fatalf("unable to set ban on legitimate user, quitting")
	}
	err = user.SetBan(1000)
	if err != nil {
		t.Fatalf("Able to set ban on user even if ban is already set, quitting")
	}
	err = user.Authorize(user.Index)
	if err != nil {
		t.Fatalf("unable to set kyc flag, exiting!")
	}
	_, err = RetrieveAllUsersWithKyc()
	if err != nil {
		t.Fatal(err)
	}

	_, err = TopReputationUsers()
	if err != nil {
		t.Fatal(err)
	}
	err = user.AddtoMailbox("subject", "message")
	if err != nil {
		t.Fatal(err)
	}
	var temp []byte
	err = user.ImportSeed(temp, "", "")
	if err == nil {
		t.Fatalf("able to decrypt empty byte array")
	}
	err = user.ImportSeed(user.StellarWallet.EncryptedSeed, user.StellarWallet.PublicKey, seedpwd)
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

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err = xlm.GetXLM(user.StellarWallet.PublicKey)
		if err != nil {
			log.Fatal(err)
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err = xlm.GetXLM(user.SecondaryWallet.PublicKey)
		if err != nil {
			log.Fatal(err)
		}
	}(&wg)

	wg.Wait()

	var wg2 sync.WaitGroup

	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err = user.IncreaseTrustLimit(seedpwd, 10)
		if err != nil {
			log.Fatal(err)
		}
	}(&wg2)

	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err = user.MoveFundsFromSecondaryWallet(10, seedpwd)
		if err != nil {
			log.Fatal(err)
		}
		err = user.SweepSecondaryWallet(seedpwd)
		if err != nil {
			log.Fatal(err)
		}
	}(&wg2)

	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		pkSeed, pk, err := xlm.GetKeyPair()
		if err != nil {
			log.Fatal(err)
		}
		assetCode := "assetcode"
		_ = build.CreditAsset{Code: assetCode, Issuer: pkSeed} // this account doesn't exist yet, so this should fail

		_, err = assets.TrustAsset(assetCode, pk, -1, seedpwd)
		if err == nil {
			log.Fatalf("can trust invalid asset")
		}
	}(&wg2)

	wg2.Wait()

	err = user.IncreaseTrustLimit("fakeseedpwd", 10)
	if err == nil {
		t.Fatalf("unable to catch fake seedpwd error")
	}

	consts.StablecoinPublicKey = "fake"

	err = user.IncreaseTrustLimit(seedpwd, 10)
	if err == nil {
		t.Fatalf("unable to catch incorrect stablecoin pubkey error")
	}

	consts.StablecoinPublicKey = "GAVEVWKMXVQ2WSCBTR7M5UKRVFFWIA52VP7ISDKZSEJKQS2VYG4D6C6P"

	err = user.MoveFundsFromSecondaryWallet(10, "shouldfail")
	if err == nil {
		t.Fatalf("decryption succeeds with invalid seedpwd for secondary account")
	}
	err = user.MoveFundsFromSecondaryWallet(100000, "shouldfail")
	if err == nil {
		t.Fatalf("unable to catch incorrect seedpwd")
	}
	err = user.MoveFundsFromSecondaryWallet(100000, seedpwd)
	if err == nil {
		t.Fatalf("can transfer more amount than possessed")
	}
	err = user.MoveFundsFromSecondaryWallet(-1, seedpwd)
	if err == nil {
		t.Fatalf("unable to catch invalid amount error")
	}

	err = user.SweepSecondaryWallet("invalidseedpwd")
	if err == nil {
		t.Fatalf("no able to catch invalid seedpwd")
	}

	for i := 0; i < 5; i++ {
		accessToken, err = user.GenAccessToken()
		if err != nil {
			t.Fatal(err)
		}
	}

	err = user.AllLogout()
	if err != nil {
		t.Fatal(err)
	}

	os.Remove(consts.DbDir + "/openx.db")
}
