package database

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	aes "github.com/OpenFinancing/openfinancing/aes"
	notif "github.com/OpenFinancing/openfinancing/notif"
	utils "github.com/OpenFinancing/openfinancing/utils"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
	"github.com/boltdb/bolt"
)

// the user structure houses all entities that are of type "User". This contains
// commonly used functions so that we need not repeat the ssame thing for every instance.
type User struct {
	Index int
	// default index, gets us easy stats on how many people are there
	EncryptedSeed []byte
	// EncryptedSeed stores the AES-256 encrypted seed of the user. This way, even
	// if the platform is hacked, the user's funds are still safe
	Name string
	// Name of the primary stakeholder involved (principal trustee of school, for eg.)
	PublicKey string
	// PublicKey denotes the public key of the recipient
	Username string
	// the username you use to login to the platform
	Pwhash string
	// the password hash, which you use to authenticate on the platform
	Address string
	// the registered address of the above company
	Description string
	// Does the contractor need to have a seed and a publickey?
	// we assume that it does in this case and proceed.
	// information on company credentials, their experience
	Image string
	// image can be company logo, founder selfie
	FirstSignedUp string
	// auto generated timestamp
	Kyc bool
	// false if kyc is not accepted / reviewed, true if user has been verified.
	// TODO: evaluate kyc providers and get a trusted partner who can do this for us (see kyc-services.md)
	Inspector bool
	// inspector is a kyc inspector who valdiates the data of people who would like
	// to signup on the platform
	Email string
	// user email to send out notifications
	Notification bool
	// GDPR, if user wants to opt in, set this to true. Default is false
	Reputation float64
	// Reputation contains the max reputation that can be gained by a user. Reputation increases
	// for each completed bond and decreases for each bond cancelled. The frontend
	// could have a table based on reputation scores and use the appropriate scores for
	// awarding badges or something to users with high reputation
	LocalAssets []string
	// a collection of assets that the user can own and trade locally using the emulator
}

// User is a metastrucutre that contains commonly used keys within a single umbrella
// so that we can import it wherever needed.
func NewUser(uname string, pwd string, seedpwd string, Name string) (User, error) {
	// call this after the user has failled in username and password.
	// Store hashed password in the database
	var a User

	allUsers, err := RetrieveAllUsers()
	if err != nil {
		return a, err
	}

	// the ugly indexing thing again, need to think of something better here
	if len(allUsers) == 0 {
		a.Index = 1
	} else {
		a.Index = len(allUsers) + 1
	}

	a.Name = Name
	err = a.GenKeys(seedpwd)
	if err != nil {
		return a, err
	}
	a.Username = uname
	a.Pwhash = utils.SHA3hash(pwd) // store tha sha3 hash
	// now we have a new User, take this and then send this struct off to be stored in the database
	a.FirstSignedUp = utils.Timestamp()
	a.Kyc = false
	a.Notification = false
	err = a.Save()
	return a, err // since user is a meta structure, insert it and then return the function
}

// InsertUser inserts a passed User object into the database
func (a *User) Save() error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			return err
		}
		return b.Put([]byte(utils.ItoB(a.Index)), encoded)
	})
	return err
}

func RetrieveAllUsersWithoutKyc() ([]User, error) {
	var arr []User
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; ; i++ {
			var rUser User
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rUser)
			if err != nil {
				return err
			}
			if !rUser.Kyc {
				arr = append(arr, rUser)
			}
		}
		return nil
	})
	return arr, err
}

func RetrieveAllUsersWithKyc() ([]User, error) {
	var arr []User
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; ; i++ {
			var rUser User
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rUser)
			if err != nil {
				return err
			}
			if rUser.Kyc {
				arr = append(arr, rUser)
			}
		}
		return nil
	})
	return arr, err
}

// RetrieveAllUsers gets a list of all User in the database
func RetrieveAllUsers() ([]User, error) {
	var arr []User
	db, err := OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; ; i++ {
			var rUser User
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rUser)
			if err != nil {
				return err
			}
			arr = append(arr, rUser)
		}
		return nil
	})
	return arr, err
}

// RetrieveUser retrieves a particular User indexed by key from the database
func RetrieveUser(key int) (User, error) {
	var inv User
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return fmt.Errorf("Retrieved user nil, quitting!")
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, err
}

func ValidateUser(name string, pwhash string) (User, error) {
	var inv User
	temp, err := RetrieveAllUsers()
	if err != nil {
		return inv, err
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; i < limit; i++ {
			var rUser User
			x := b.Get(utils.ItoB(i))
			err := json.Unmarshal(x, &rUser)
			if err != nil {
				return err
			}
			// check names
			if rUser.Username == name && rUser.Pwhash == pwhash {
				inv = rUser
				return nil
			}
		}
		return fmt.Errorf("Not Found")
	})
	return inv, err
}

func (a *User) GenKeys(seedpwd string) error {
	var err error
	var seed string
	seed, a.PublicKey, err = xlm.GetKeyPair()
	if err != nil {
		return err
	}
	// don't store the seed in the database
	a.EncryptedSeed, err = aes.Encrypt([]byte(seed), seedpwd)
	err = a.Save()
	return err
}

func (a *User) GetSeed(seedpwd string) (string, error) {
	return wallet.DecryptSeed(a.EncryptedSeed, seedpwd)
}

// CheckUsernameCollision checks if a username is available to a new user who
// wants to signup on the platform
func CheckUsernameCollision(uname string) error {
	temp, err := RetrieveAllUsers()
	if err != nil {
		return err
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; i < limit; i++ {
			var rUser User
			x := b.Get(utils.ItoB(i))
			err := json.Unmarshal(x, &rUser)
			if err != nil {
				return err
			}
			// check names
			if rUser.Username == uname {
				return fmt.Errorf("Username collision")
			}
		}
		return nil
	})
	return err
}

// Everything above this is exactly the same as the investor class. Need to replicate
// because of bucket issues, hopefully there's a nicer way
// package kyc is designed to emulate the working of a kyc entity in the system
// when someone signs up on the system, they appear in they KYC reviewer's panel
// and the inspector has to approve them for them to  be able to go on the platform.
// iF rejected, they can choose to submit additional information that they are indeed
// compliant and in that case we can allow them unto the platform
// Roughly this should involve a new bool in the user bucket which says kyc and only
// the inspector should have the power to set it to true.
// the inspector itself requires kyc though, so we shall have an admin account which can
// kickoff the kyc process.
// MWTODO: what do we do with these KYC powers? what features are open and what can be
// viewed only by going through KYC?
func (a *User) Authorize(userIndex int) error {
	// we don't really mind who this user is since all we need to verify is his identity
	if !a.Inspector {
		return fmt.Errorf("You don't have the required permissions to kyc a person")
	}
	user, err := RetrieveUser(userIndex)
	// we want to retrieve only users who have not gone through KYC before
	if err != nil {
		return err
	}
	if user.Kyc {
		return fmt.Errorf("user already KYC'd")
	}
	user.Kyc = true
	return user.Save()
}

func AddInspector(userIndex int) error {
	// this should only be called by the platform itself and not open to others
	user, err := RetrieveUser(userIndex)
	if err != nil {
		return err
	}
	user.Inspector = true
	return user.Save()
}

// these two functions can be used as internal hnadlers and hte RPC can save reputation directly
func (a *User) IncreaseReputation(reputation float64) error {
	a.Reputation += reputation
	return a.Save()
}

func (a *User) DecreaseReputation(reputation float64) error {
	a.Reputation -= reputation
	return a.Save()
}

func TopReputationUsers() ([]User, error) {
	// these reputation functions should mostly be used by the frontend through the
	// RPC to display to other users what other users' reputation is.
	allUsers, err := RetrieveAllUsers()
	if err != nil {
		return allUsers, err
	}
	for i, _ := range allUsers {
		for j, _ := range allUsers {
			if allUsers[i].Reputation > allUsers[j].Reputation {
				tmp := allUsers[i]
				allUsers[i] = allUsers[j]
				allUsers[j] = tmp
			}
		}
	}
	return allUsers, nil
}

func AgreeToContractConditions(contractHash string, projIndex string,
	debtAssetCode string, userIndex int, seedpwd string) error {
	// we need to display this on the frontend and once the user presses agree, commit
	// a tx to the blockchain with the outcome
	message := "I agree to the terms and conditions specified in contract " + contractHash +
		"and by signing this message to the blockchain agree that I accept the investment in project " + projIndex +
		"whose debt asset is: " + debtAssetCode
	// hash the message and transmit the message in 5 parts
	// eg.
	// CONTRACTHASH9a768ace36ff3d17
	// 71d5c145a544de3d68343b2e7609
	// 3cb7b2a8ea89ac7f1a20c852e6fc
	// 1d71275b43abffefac381c5b906f
	// 55c3bcff4225353d02f1d3498758

	user, err := RetrieveUser(userIndex)
	if err != nil {
		log.Println(err)
		return err
	}

	seed, err := wallet.DecryptSeed(user.EncryptedSeed, seedpwd)
	if err != nil {
		log.Println(err)
		return err
	}

	messageHash := "CONTRACTHASH" + strings.ToUpper(utils.SHA3hash(message))
	firstPart := messageHash[:28] // higher limit is not included in the slice
	secondPart := messageHash[28:56]
	thirdPart := messageHash[56:84]
	fourthPart := messageHash[84:112]
	fifthPart := messageHash[112:140]

	timeStamp := utils.I64toS(utils.Unix())
	_, firstHash, err := xlm.SendXLM(user.PublicKey, timeStamp, seed, firstPart)
	if err != nil {
		log.Println(err)
		return err
	}

	_, secondHash, err := xlm.SendXLM(user.PublicKey, timeStamp, seed, secondPart)
	if err != nil {
		log.Println(err)
		return err
	}

	_, thirdHash, err := xlm.SendXLM(user.PublicKey, timeStamp, seed, thirdPart)
	if err != nil {
		log.Println(err)
		return err
	}

	_, fourthHash, err := xlm.SendXLM(user.PublicKey, timeStamp, seed, fourthPart)
	if err != nil {
		log.Println(err)
		return err
	}

	_, fifthHash, err := xlm.SendXLM(user.PublicKey, timeStamp, seed, fifthPart)
	if err != nil {
		log.Println(err)
		return err
	}

	//if user.Notification {
	notif.SendContractNotification(firstHash, secondHash, thirdHash, fourthHash, fifthHash, "varunramganesh@gmail.com")
	//}

	return nil
}
