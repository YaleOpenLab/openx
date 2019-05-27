package database

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/pkg/errors"
	"log"

	aes "github.com/YaleOpenLab/openx/aes"
	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	recovery "github.com/YaleOpenLab/openx/recovery"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
	xlm "github.com/YaleOpenLab/openx/xlm"
	"github.com/boltdb/bolt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	crypto "github.com/ethereum/go-ethereum/crypto"
)

// User is a metastrucutre that contains commonly used keys within a single umbrella
// so that we can import it wherever needed.
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
	City string
	// the city of residence of the resident
	ZipCode string
	// the zipcode of hte particular city
	Country string
	// the coutnry of residence of the resident
	RecoveryPhone string
	// the phone number where we need to send recovery codes to
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
	Inspector bool
	// inspector is a kyc inspector who valdiates the data of people who would like
	// to signup on the platform
	Banned bool
	// a field which can be used to set a ban on a user. Can be only used by inspectors in the event someone
	// who has KYC is known to behave in a suspicious way.
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
	RecoveryShares []string
	// RecoveryShares are shares that you could hare out to a party and one could reconstruct the
	// seed from 2 out of 3 parts. Based on Shamir's Secret Sharing Scheme.
	PwdResetCode string

	SecondaryWallet Wallet
	// SecondaryWallet defines a higher level wallet which can be imagined to be similar to a savings account

	EthereumWallet EthWallet
	// EthereumWallet defines a separate wallet for ethereum which people can use to control their ERC721 RECs

	PendingDocuments map[string]string
	// a Pending documents map to keep track of documents that the user in question has to keep track of
	// related to a specific project. The key is the same as the value of the project and the value is a description
	// of what exactly needs to be submitted.
	KYC KycStruct

	StarRating map[int]int // peer bases tarr rating that users can give of each other. Can be gamed, but this is complemented by
	// the automated feedback system, so we should be good.

	GivenStarRating map[int]int // to keep track of users for whom you've given feedback
}

type KycStruct struct {
	PassportPhoto  string // should be a base64 string or similar according to what the API provider wants
	IDCardPhoto    string
	DriversLicense string
	PersonalPhoto  string // a selfie to verify that  the person registering on the platform is the same person whose documents have been uploaded
}

// EthWallet contains the structures needed for an ethereum wallet
type EthWallet struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

// Wallet contains the stuff that we need for a wallet.
type Wallet struct {
	EncryptedSeed []byte // the seedpwd for this would be the same as the one for the primary wallet
	// since we don't want the user to remember like 10 passwords
	PublicKey string
}

// NewUser creates a new user
func NewUser(uname string, pwd string, seedpwd string, Name string) (User, error) {
	// call this after the user has failled in username and password.
	// Store hashed password in the database
	var a User

	allUsers, err := RetrieveAllUsers()
	if err != nil {
		return a, errors.Wrap(err, "Error while retrieving all users from database")
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
		return a, errors.Wrap(err, "Error while generating public and private keys")
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

// Save inserts a passed User object into the database
func (a *User) Save() error {
	db, err := OpenDB()
	if err != nil {
		return errors.Wrap(err, "Error while opening database")
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		encoded, err := a.MarshalJSON()
		if err != nil {
			return errors.Wrap(err, "Error while marshaling json")
		}
		return b.Put([]byte(utils.ItoB(a.Index)), encoded)
	})
	return err
}

// RetrieveAllUsersWithoutKyc retrieves all users without kyc
func RetrieveAllUsersWithoutKyc() ([]User, error) {
	var arr []User
	db, err := OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "Error while opening database")
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
			err := rUser.UnmarshalJSON(x)
			if err != nil {
				return errors.Wrap(err, "Error while unmarshalling json")
			}
			if !rUser.Kyc {
				arr = append(arr, rUser)
			}
		}
	})
	return arr, err
}

// RetrieveAllUsersWithKyc retrieves all users with kyc
func RetrieveAllUsersWithKyc() ([]User, error) {
	var arr []User
	db, err := OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "Error while opening database")
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
			err := rUser.UnmarshalJSON(x)
			if err != nil {
				return errors.Wrap(err, "Error while unmarshalling json")
			}
			if rUser.Kyc {
				arr = append(arr, rUser)
			}
		}
	})
	return arr, err
}

// RetrieveAllUsers gets a list of all User in the database
func RetrieveAllUsers() ([]User, error) {
	var arr []User
	db, err := OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "Error while opening database")
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
			err := rUser.UnmarshalJSON(x)
			if err != nil {
				return errors.Wrap(err, "Error while unmarshalling json")
			}
			arr = append(arr, rUser)
		}
	})
	return arr, err
}

// RetrieveUser retrieves a particular User indexed by key from the database
func RetrieveUser(key int) (User, error) {
	var inv User
	db, err := OpenDB()
	if err != nil {
		return inv, errors.Wrap(err, "error while opening database")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return errors.New("retrieved user nil, quitting!")
		}
		return inv.UnmarshalJSON(x)
	})
	return inv, err
}

// ValidateUser validates a particular user
func ValidateUser(name string, pwhash string) (User, error) {
	var inv User
	temp, err := RetrieveAllUsers()
	if err != nil {
		return inv, errors.Wrap(err, "error while retrieving all users from database")
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return inv, errors.Wrap(err, "could not open db, quitting!")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; i < limit; i++ {
			var rUser User
			x := b.Get(utils.ItoB(i))
			err := rUser.UnmarshalJSON(x)
			if err != nil {
				return errors.Wrap(err, "could not unmarshal json, quitting!")
			}
			// check names
			if rUser.Username == name && rUser.Pwhash == pwhash {
				inv = rUser
				return nil
			}
		}
		return errors.New("Not Found")
	})
	return inv, err
}

// GenKeys generates a keypair for the user
func (a *User) GenKeys(seedpwd string) error {
	var err error
	var seed string
	seed, a.PublicKey, err = xlm.GetKeyPair()
	if err != nil {
		return errors.Wrap(err, "error while generating public and private key pair")
	}
	// don't store the seed in the database
	a.EncryptedSeed, err = aes.Encrypt([]byte(seed), seedpwd)
	if err != nil {
		return errors.Wrap(err, "error while encrypting seed")
	}

	tmp, err := recovery.Create(2, 3, seed)
	if err != nil {
		return errors.Wrap(err, "error while storing recovery shares")
	}

	a.RecoveryShares = append(a.RecoveryShares, tmp...) // this is for the primary account

	secSeed, secPubkey, err := xlm.GetKeyPair()
	if err != nil {
		return errors.Wrap(err, "could not generate secondary keypair")
	}

	a.SecondaryWallet.PublicKey = secPubkey
	a.SecondaryWallet.EncryptedSeed, err = aes.Encrypt([]byte(secSeed), seedpwd)
	if err != nil {
		return errors.Wrap(err, "error while encrypting seed")
	}

	ecdsaPrivkey, err := crypto.GenerateKey()
	if err != nil {
		return errors.Wrap(err, "could not generate an ethereum keypair, quitting!")
	}

	privateKeyBytes := crypto.FromECDSA(ecdsaPrivkey)
	a.EthereumWallet.PrivateKey = hexutil.Encode(privateKeyBytes)[2:]
	a.EthereumWallet.Address = crypto.PubkeyToAddress(ecdsaPrivkey.PublicKey).Hex()

	publicKeyECDSA, ok := ecdsaPrivkey.Public().(*ecdsa.PublicKey)
	if !ok {
		return errors.Wrap(err, "error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	a.EthereumWallet.PublicKey = hexutil.Encode(publicKeyBytes)[4:] // an ethereum address is 65 bytes long and hte first byte is 0x04 for DER encoding, so we omit that

	if crypto.PubkeyToAddress(*publicKeyECDSA).Hex() != a.EthereumWallet.Address {
		return errors.Wrap(err, "addresses don't match, quitting!")
	}

	err = a.Save()
	return err
}

// CheckUsernameCollision checks if a username is available to a new user who
// wants to signup on the platform
func CheckUsernameCollision(uname string) (User, error) {
	var dummy User
	temp, err := RetrieveAllUsers()
	if err != nil {
		return dummy, errors.Wrap(err, "error while retrieving all users from database")
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return dummy, errors.Wrap(err, "error while opening database")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; i < limit; i++ {
			var rUser User
			x := b.Get(utils.ItoB(i))
			err := rUser.UnmarshalJSON(x)
			if err != nil {
				return errors.Wrap(err, "error while unmarshalling json")
			}
			// check names
			if rUser.Username == uname {
				dummy = rUser
				return errors.New("Username collision")
			}
		}
		return nil
	})
	return dummy, err
}

// Authorize authorizes a user
func (a *User) Authorize(userIndex int) error {
	// we don't really mind who this user is since all we need to verify is his identity
	if !a.Inspector {
		return errors.New("You don't have the required permissions to kyc a person")
	}
	user, err := RetrieveUser(userIndex)
	// we want to retrieve only users who have not gone through KYC before
	if err != nil {
		return errors.Wrap(err, "error while retrieving user from database")
	}
	if user.Kyc {
		return errors.New("user already KYC'd")
	}
	user.Kyc = true
	return user.Save()
}

// AddInspector adds a kyc inspector
func AddInspector(userIndex int) error {
	// this should only be called by the platform itself and not open to others
	user, err := RetrieveUser(userIndex)
	if err != nil {
		return errors.Wrap(err, "error while retrieving user from database")
	}
	user.Inspector = true
	return user.Save()
}

// these two functions can be used as internal hnadlers and hte RPC can save reputation directly

// IncreaseReputation increases reputation
func (a *User) IncreaseReputation(reputation float64) error {
	a.Reputation += reputation
	return a.Save()
}

// DecreaseReputation decreases reputation
func (a *User) DecreaseReputation(reputation float64) error {
	a.Reputation -= reputation
	return a.Save()
}

// TopReputationUsers gets the users with top reputation
func TopReputationUsers() ([]User, error) {
	// these reputation functions should mostly be used by the frontend through the
	// RPC to display to other users what other users' reputation is.
	allUsers, err := RetrieveAllUsers()
	if err != nil {
		return allUsers, errors.Wrap(err, "error while retrieving all users from database")
	}
	for i := range allUsers {
		for j := range allUsers {
			if allUsers[i].Reputation > allUsers[j].Reputation {
				tmp := allUsers[i]
				allUsers[i] = allUsers[j]
				allUsers[j] = tmp
			}
		}
	}
	return allUsers, nil
}

func IncreaseTrustLimit(userIndex int, seedpwd string, trust string) error {

	user, err := RetrieveUser(userIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve user from database, quitting!")
	}

	seed, err := wallet.DecryptSeed(user.EncryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "couldn't decrypt seed, quitting!")
	}

	// we now have the seed, so we should upgrade the trustlimit by the margin requested. The margin passed here
	// must not include the old trustlimit

	trustLimit := utils.StoF(trust) + utils.StoF(consts.StablecoinTrustLimit)

	_, err = assets.TrustAsset(consts.StablecoinCode, consts.StableCoinAddress, utils.FtoS(trustLimit), seed)
	if err != nil {
		return errors.Wrap(err, "couldn't trust asset, quitting!")
	}

	return nil
}

// SearchWithEmailId searches for a given user who has the given email id
func SearchWithEmailId(email string) (User, error) {
	var foundUser User
	db, err := OpenDB()
	if err != nil {
		return foundUser, errors.Wrap(err, "Error while opening database")
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
			err := rUser.UnmarshalJSON(x)
			if err != nil {
				return errors.Wrap(err, "Error while unmarshalling json")
			}
			if rUser.Email == email {
				foundUser = rUser
			}
		}
	})
	return foundUser, err
}

// MoveFundsFromSecondaryWallet moves funds from the secondary wallet to the primary wallet
func MoveFundsFromSecondaryWallet(userIndex int, pwhash string, amount string, seedpwd string) error {
	user, err := RetrieveUser(userIndex)
	if err != nil {
		return errors.Wrap(err, "could not retrieve user, quitting")
	}

	if user.Pwhash != pwhash {
		return fmt.Errorf("pw hashes don't match, quitting")
	}
	amountI, err := utils.StoFWithCheck(amount)
	if err != nil {
		return errors.Wrap(err, "amount not float, quitting")
	}
	// unlock secondary account
	secSeed, err := wallet.DecryptSeed(user.SecondaryWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "could not unlock secondary seed, quitting")
	}

	// get secondary balance
	secFunds, err := xlm.GetNativeBalance(user.SecondaryWallet.PublicKey)
	if err != nil {
		return errors.Wrap(err, "could not get xlm balance of secondary account")
	}

	if amountI > utils.StoF(secFunds) {
		return fmt.Errorf("amount to be transferred is greater than the funds available in the secondary account, quitting")
	}

	// send the tx over
	_, txhash, err := xlm.SendXLM(user.PublicKey, amount, secSeed, "fund transfer to secondary")
	if err != nil {
		return errors.Wrap(err, "error while transferring funds to secondary account, quitting")
	}

	log.Println("transfer sec-prim tx hash: ", txhash)
	return nil
}

// SweepSecondaryWallet sweeps fudsd from the secondary account to the primary account
func SweepSecondaryWallet(userIndex int, pwhash string, seedpwd string) error {
	// unlock secondary account

	user, err := RetrieveUser(userIndex)
	if err != nil {
		return errors.Wrap(err, "could not retrieve user, quitting")
	}

	if user.Pwhash != pwhash {
		return fmt.Errorf("pw hashes don't match, quitting")
	}

	secSeed, err := wallet.DecryptSeed(user.SecondaryWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "could not unlock primary seed, quitting")
	}

	// get secondary balance
	secFunds, err := xlm.GetNativeBalance(user.SecondaryWallet.PublicKey)
	if err != nil {
		return errors.Wrap(err, "could not get xlm balance of secondary account")
	}

	secFundsWithMinbal := utils.FtoS(utils.StoF(secFunds) - 5)
	// send the tx over
	_, txhash, err := xlm.SendXLM(user.PublicKey, secFundsWithMinbal, secSeed, "fund transfer to secondary")
	if err != nil {
		return errors.Wrap(err, "error while transferring funds to secondary account, quitting")
	}

	log.Println("transfer sec-prim tx hash: ", txhash)
	return nil
}

// AddEmail stores the passed email as the user's email.
func (a *User) AddEmail(email string) error {
	// call this function when a user wants to get notifications. Ask on frontend whether
	// it wants to
	a.Email = email
	a.Notification = true
	err := a.Save()
	if err != nil {
		return errors.Wrap(err, "error while saving investor")
	}
	return a.Save()
}

func (a *User) SetBan(userIndex int) error {
	if !a.Inspector {
		return fmt.Errorf("user not authorized to ban a user")
	}

	if a.Index == userIndex {
		return fmt.Errorf("can't ban yourself, quitting!")
	}

	user, err := RetrieveUser(userIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't  find user to ban, quitting")
	}

	if user.Banned {
		return errors.Wrap(err, "user already banned, not setitng another ban")
	}

	user.Banned = true
	return user.Save()
}

func (a *User) GiveFeedback(userIndex int, feedback int) error {
	user, err := RetrieveUser(userIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve user from db while giving feedback")
	}

	if len(user.StarRating) == 0 {
		// no one has given t3his user a starr rating before, so create a new map
		user.StarRating = make(map[int]int)
	}

	user.StarRating[a.Index] = feedback
	log.Println("STARRATING: ", user.StarRating, user.Name)
	err = user.Save()
	if err != nil {
		return errors.Wrap(err, "couldn't save feedback provided on user")
	}

	if len(a.GivenStarRating) == 0 {
		// no one has given t3his user a starr rating before, so create a new map
		a.GivenStarRating = make(map[int]int)
	}

	a.GivenStarRating[user.Index] = feedback
	return a.Save()
}
