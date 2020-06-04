package database

import (
	"encoding/base32"
	"log"
	"strings"

	"github.com/pkg/errors"

	aes "github.com/Varunram/essentials/aes"
	algorand "github.com/Varunram/essentials/algorand"
	googauth "github.com/Varunram/essentials/googauth"
	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/Varunram/essentials/xlm"
	assets "github.com/Varunram/essentials/xlm/assets"
	wallet "github.com/Varunram/essentials/xlm/wallet"
	consts "github.com/YaleOpenLab/openx/consts"
	recovery "github.com/bithyve/research/sss"
)

// User defines a base layer structure that can be used by entities on platforms built on openx
type User struct {
	// Index is an incremental index maintained to easily retrieve users
	Index int
	// Name is the real name of the user
	Name string
	// Description contains some other stuff that a user might wish to describe about themselves
	Description string
	// Image is an optional profile image that users might upload
	Image string
	// FirstSignedUp contains the date on which the users signed up on openx
	FirstSignedUp string
	// Address is the address of the user
	Address string
	// City denotes the city of residence of the user
	City string
	// ZipCode is the zipcode iof the user's residence
	ZipCode string
	// Country is the country of residence of the user. Some platforms may be restricted to a certain subset of countries
	Country string
	// RecoveryPhone used to send recovery codes or contact the user in the event of an emergency
	RecoveryPhone string
	// Email is used to send users notifications on their actions on openx based platforms
	Email string
	// Notification is a bool which denotes whether the user wants to receive notifications related to the openx platform
	Notification bool
	// StellarWallet contains a publickey and encrypted seed that can be used to interact with the Stellar blockchain
	StellarWallet StellWallet
	// AlgorandWallet contains a publickey and privatekey pair which can be used to interact with the Algorand blockchain
	AlgorandWallet algorand.Wallet
	// Username denoteds the username of the user to log on to openx
	Username string
	// Pwhash is the 512 byte SHA-3 hash of the user's password
	Pwhash string
	// Kyc denotes whether the user has passed KYC
	Kyc bool
	// Admin denotes whether the user has passed kyc or not
	Admin bool
	// Inspector denotes whether the user is a KYC inspector ie whether they're authorized to approve other's KYC requests
	Inspector bool
	// Banned is true if the user is banned on openx
	Banned bool
	// Reputation is a float which denotes the reputation of a user on the openx platform
	Reputation float64
	// LocalAssets is a list of P2P assets belonging to the user
	LocalAssets []string
	// RecoveryShares is a collection of shares that a user can distribute to aid recovery of their seed later on
	RecoveryShares []string
	// PwdResetCode is a code that's set when a user wants to reset their password
	PwdResetCode string
	// SecondaryWallet is a secondary wallet where people can store their funds in
	SecondaryWallet StellWallet
	// PendingDocuments is a list of pending documents which a user must upload before progressing on to the next stage
	PendingDocuments map[string]string
	// KYC contains KYC information required by ComplyAdvantage
	KYC KycStruct
	// StarRating is a star rating similar to popular platforms which users can use to rate each other
	StarRating map[int]int
	// GivenStarRating contains a list of users whom this user has rated
	GivenStarRating map[int]int
	// TwoFASecret is the secret associated with Google 2FA that users can enable while logging on openx
	TwoFASecret string
	// AnchorKYC contains KYC information required by AnchorUSD
	AnchorKYC AnchorKYCHelper
	// AccessToken is the access token that will be used for authenticating RPC requests made to the server
	AccessToken map[string]int64
	// Mailbox is a mailbox where admins can send you messages or updated on your invested / interested projects
	Mailbox []MailboxHelper
	// Legal is a bool which is set when the user accepts the terms and conditions
	Legal bool
	// ProfileProgress is a float which denotes user profile completeness on the frontend
	ProfileProgress float64
	// Verified marks a person as verified
	Verified bool
	// VerifyReq requests verification
	VerifyReq bool
	// VerifiedBy stores the index of the admin who verified the user
	VerifiedBy int
	// VerifiedTime stores when the user was verified
	VerifiedTime string
	// ConfToken is the confirmation token sent to users to confirm their registration on openx
	ConfToken string
	// Conf is a bool that is set to true when users confirm their tokens
	Conf bool
}

// MailboxHelper is a helper struct that can be used to send admin notifications to users
type MailboxHelper struct {
	Subject string // the subject can be used to send push notifications
	Message string // the message
}

// KycStruct contains the parameters required by ComplyAdvantage
type KycStruct struct {
	PassportPhoto  string
	IDCardPhoto    string
	DriversLicense string
	PersonalPhoto  string
}

// AnchorKYCHelper contains the KYC parameters required by Anchor
type AnchorKYCHelper struct {
	Name     string
	Birthday struct {
		Month string
		Day   string
		Year  string
	}
	Tax struct {
		Country string
		ID      string
	}
	Address struct {
		Street  string
		City    string
		Postal  string
		Region  string
		Country string
		Phone   string
	}
	PrimaryPhone       string
	Gender             string
	DepositIdentifier  string
	WithdrawIdentifier string
	URL                string
	AccountID          string
}

// StellWallet hold the Stellar Publickey and Encrypted Seed
type StellWallet struct {
	PublicKey     string
	EncryptedSeed []byte
	SeedPwhash    string
}

// NewUser creates a new user, stores it in the openx database and returns a user struct
func NewUser(uname string, pwhash string, seedpwd string, email string) (User, error) {
	var a User

	_, err := CheckUsernameCollision(uname)
	if err != nil {
		return a, errors.Wrap(err, "username collision: "+uname+", quitting")
	}

	lim, err := RetrieveAllUsersLim()
	if err != nil {
		return a, errors.Wrap(err, "Error while retrieving all users from database")
	}
	a.Index = lim + 1

	err = a.GenKeys(seedpwd)
	if err != nil {
		return a, errors.Wrap(err, "Error while generating public and private keys")
	}

	a.Email = email
	a.Username = uname
	a.Pwhash = pwhash
	a.FirstSignedUp = utils.Timestamp()
	a.Kyc = false
	a.Notification = false
	a.ConfToken = strings.ToUpper(utils.GetRandomString(8))
	log.Println("saving: ", uname, pwhash, seedpwd, email)
	err = a.Save()
	return a, err
}

// RetrieveAllUsersWithoutKyc retrieves all users without kyc
func RetrieveAllUsersWithoutKyc() ([]User, error) {
	var arr []User

	users, err := RetrieveAllUsers()
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		if !user.Kyc {
			arr = append(arr, user)
		}
	}

	return arr, nil
}

// RetrieveAllUsersWithKyc retrieves all users with kyc
func RetrieveAllUsersWithKyc() ([]User, error) {
	var arr []User

	users, err := RetrieveAllUsers()
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		if user.Kyc {
			arr = append(arr, user)
		}
	}

	return arr, nil
}

// ValidateSeedpwd validates a user and their seedpwd
func ValidateSeedpwd(name string, pwhash string, seedpwd string) (User, error) {
	user, err := ValidatePwhash(name, pwhash)
	if err != nil {
		return user, errors.Wrap(err, "could not validate user")
	}
	seed, err := wallet.DecryptSeed(user.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return user, errors.Wrap(err, "failed to decrypt user's seed")
	}
	pubkey, err := wallet.ReturnPubkey(seed)
	if err != nil {
		return user, errors.Wrap(err, "could not decrypt seed")
	}
	if pubkey != user.StellarWallet.PublicKey {
		return user, errors.New("pubkeys don't match, quitting")
	}
	return user, nil
}

// ValidateSeedpwdAuthToken validates a user and their seedpwd using their accesstoken
func ValidateSeedpwdAuthToken(name string, token string, seedpwd string) (User, error) {
	user, err := ValidateAccessToken(name, token)
	if err != nil {
		return user, errors.Wrap(err, "could not validate user")
	}
	seed, err := wallet.DecryptSeed(user.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return user, errors.Wrap(err, "failed to decrypt user's seed")
	}
	pubkey, err := wallet.ReturnPubkey(seed)
	if err != nil {
		return user, errors.Wrap(err, "could not decrypt seed")
	}
	if pubkey != user.StellarWallet.PublicKey {
		return user, errors.New("pubkeys don't match, quitting")
	}
	return user, nil
}

// GenKeys generates a keypair for the user and takes in options on which blockchain to generate keys for
func (a *User) GenKeys(seedpwd string, options ...string) error {
	if len(options) == 1 {
		if consts.Mainnet {
			return errors.New("only stellar supported in mainnet mode, quitting")
		}
		chain := options[0]
		switch chain {
		/*
			case "algorand":
				log.Println("Generating Algorand wallet")

				var err error
				password := seedpwd

				a.AlgorandWallet, err = algorand.GenNewWallet("algowl", password)
				if err != nil {
					return errors.Wrap(err, "couldn't create new wallet id, quitting")
				}

				err = a.Save()
				if err != nil {
					return err
				}

				backupPhrase, err := algorand.GenerateBackup(a.AlgorandWallet.WalletName, password)
				if err != nil {
					return err
				}

				tmp, err := recovery.Create(2, 3, backupPhrase)
				if err != nil {
					return errors.Wrap(err, "error while storing recovery shares")
				}

				a.RecoveryShares = append(a.RecoveryShares, tmp...)
		*/
		default:
			log.Println("Chain not supported, please feel free to add support in aanew Pull Request")
			return errors.New("chain not supported, returning")
		} // end of switch
	} else if len(options) == 0 {
		// if no option is provided, default to Stellar
		var err error
		var seed string
		seed, a.StellarWallet.PublicKey, err = xlm.GetKeyPair()
		if err != nil {
			return errors.Wrap(err, "error while generating public and private key pair")
		}

		a.StellarWallet.EncryptedSeed, err = aes.Encrypt([]byte(seed), seedpwd)
		if err != nil {
			return errors.Wrap(err, "error while encrypting seed")
		}

		tmp, err := recovery.Create(2, 3, seed)
		if err != nil {
			return errors.Wrap(err, "error while storing recovery shares")
		}

		a.RecoveryShares = append(a.RecoveryShares, tmp...)
	}

	secSeed, secPubkey, err := xlm.GetKeyPair()
	if err != nil {
		return errors.Wrap(err, "could not generate secondary keypair")
	}

	a.SecondaryWallet.PublicKey = secPubkey
	a.SecondaryWallet.EncryptedSeed, err = aes.Encrypt([]byte(secSeed), seedpwd)
	if err != nil {
		return errors.Wrap(err, "error while encrypting seed")
	}

	err = a.Save()
	return err
}

// CheckUsernameCollision checks if a passed username collides with someone who's already
// on the platform. If a collision does exist, return the existing user in the database
func CheckUsernameCollision(uname string) (User, error) {
	var dummy User
	users, err := RetrieveAllUsers()
	if err != nil {
		return dummy, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		if user.Username == uname {
			return user, errors.New("username collision observed, quitting")
		}
	}

	return dummy, nil
}

// Authorize sets the Kyc flag on a user. Can only be called by Inspectors
func (a *User) Authorize(userIndex int) error {
	// we don't really mind who this user is since all we need to verify is his identity
	if !a.Inspector && !a.Admin {
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

// AddInspector sets the Inspector flag on a user
func AddInspector(userIndex int) error {
	// this should only be called by the platform itself and not open to others
	user, err := RetrieveUser(userIndex)
	if err != nil {
		return errors.Wrap(err, "error while retrieving user from database")
	}
	user.Inspector = true
	return user.Save()
}

// ChangeReputation changes the reputation associated with a user
func (a *User) ChangeReputation(reputation float64) error {
	a.Reputation += reputation
	return a.Save()
}

// IncreaseTrustLimit increases the trust limit of a user towards the in house stablecoin
func (a *User) IncreaseTrustLimit(seedpwd string, trust float64) error {

	seed, err := wallet.DecryptSeed(a.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "couldn't decrypt seed, quitting!")
	}

	if !consts.Mainnet {
		_, err = assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, trust+consts.StablecoinTrustLimit, seed)
		if err != nil {
			return errors.Wrap(err, "couldn't trust asset, quitting!")
		}
	} else {
		_, err = assets.TrustAsset(consts.AnchorUSDCode, consts.AnchorUSDAddress, trust+consts.AnchorUSDTrustLimit, seed)
		if err != nil {
			return errors.Wrap(err, "couldn't trust asset, quitting!")
		}
	}

	return nil
}

// SearchWithEmailID searches for a user given their email id
func SearchWithEmailID(email string) (User, error) {
	var dummy User
	users, err := RetrieveAllUsers()
	if err != nil {
		return dummy, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	return dummy, errors.New("could not find user with requested email id, quitting")
}

// MoveFundsFromSecondaryWallet moves XLM from the secondary wallet to the primary wallet
func (a *User) MoveFundsFromSecondaryWallet(amount float64, seedpwd string) error {
	// unlock secondary account
	secSeed, err := wallet.DecryptSeed(a.SecondaryWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "could not unlock secondary seed, quitting")
	}

	secFunds := xlm.GetNativeBalance(a.SecondaryWallet.PublicKey)
	if amount > secFunds {
		return errors.New("amount to be transferred is greater than the funds available in the secondary account, quitting")
	}

	_, txhash, err := xlm.SendXLM(a.StellarWallet.PublicKey, amount, secSeed, "fund transfer to secondary")
	if err != nil {
		return errors.Wrap(err, "error while transferring funds to secondary account, quitting")
	}

	log.Println("transfer sec-prim tx hash: ", txhash)
	return nil
}

// SweepSecondaryWallet sweeps XLM from the secondary account to the primary account
func (a *User) SweepSecondaryWallet(seedpwd string) error {
	// unlock secondary account

	secSeed, err := wallet.DecryptSeed(a.SecondaryWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "could not unlock primary seed, quitting")
	}

	secFunds := xlm.GetNativeBalance(a.SecondaryWallet.PublicKey)
	_, txhash, err := xlm.SendXLM(a.StellarWallet.PublicKey, secFunds-5, secSeed, "fund transfer to secondary")
	if err != nil {
		return errors.Wrap(err, "error while transferring funds to secondary account, quitting")
	}

	log.Println("transfer sec-prim tx hash: ", txhash)
	return nil
}

// AddEmail adds the email field to a given user
func (a *User) AddEmail(email string) error {
	a.Email = email
	a.Notification = true
	return a.Save()
}

// SetBan sets the Banned flag on a particular user
func (a *User) SetBan(userIndex int) error {
	user, err := RetrieveUser(userIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't  find user to ban, quitting")
	}

	if !a.Admin {
		return errors.New("user not authorized to ban a user")
	}

	if a.Index == userIndex {
		return errors.New("can't ban yourself, quitting")
	}

	if user.Banned {
		return errors.Wrap(err, "user already banned, not setitng another ban")
	}

	user.Banned = true
	return user.Save()
}

// GiveFeedback is used to rate another user
func (a *User) GiveFeedback(userIndex int, feedback int) error {
	user, err := RetrieveUser(userIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve user from db while giving feedback")
	}

	if len(user.StarRating) == 0 {
		user.StarRating = make(map[int]int)
	}

	if feedback > 5 || feedback < 0 {
		log.Println("feedback greater than 5 or less than 0, quitting")
		return errors.New("feedback greater than 5, quitting")
	}

	user.StarRating[a.Index] = feedback
	log.Println("STARRATING: ", user.StarRating, user.Name)
	err = user.Save()
	if err != nil {
		return errors.Wrap(err, "couldn't save feedback provided on user")
	}

	if len(a.GivenStarRating) == 0 {
		a.GivenStarRating = make(map[int]int)
	}

	a.GivenStarRating[user.Index] = feedback
	return a.Save()
}

// Generate2FA generates a new 2FA secret for the given user
func (a *User) Generate2FA() (string, error) {
	secret := utils.GetRandomString(35)
	secretBase32 := base32.StdEncoding.EncodeToString([]byte(secret))
	otpc := &googauth.OTPConfig{
		Secret:     secretBase32,
		WindowSize: 1,
		UTC:        true,
	}
	otpString, err := otpc.GenerateURI(a.Name)
	if err != nil {
		return otpString, err
	}
	if err != nil {
		return otpString, err
	}
	a.TwoFASecret = secret
	err = a.Save()
	if err != nil {
		return otpString, err
	}
	return otpString, nil
}

// Authenticate2FA authenticates the given password against the user's stored 2fA secret
func (a *User) Authenticate2FA(password string) (bool, error) {
	secretBase32 := base32.StdEncoding.EncodeToString([]byte(a.TwoFASecret))
	otpc := &googauth.OTPConfig{
		Secret:     secretBase32,
		WindowSize: 1,
		UTC:        true,
	}

	return otpc.Authenticate(password)
}

// ImportSeed can be used to import an ecrypted seed
func (a *User) ImportSeed(encryptedSeed []byte, pubkey string, seedpwd string) error {
	seed, err := wallet.DecryptSeed(encryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "could not decrypt seed")
	}
	checkPubkey, err := wallet.ReturnPubkey(seed)
	if err != nil {
		return errors.Wrap(err, "could not get pubkey from encrypted seed")
	}
	if pubkey != checkPubkey {
		return errors.New("decrypted pubkey does not match with provided pubkey")
	}
	a.StellarWallet.EncryptedSeed = encryptedSeed
	a.StellarWallet.PublicKey = pubkey
	return a.Save()
}

// GenAccessToken generates a new access token for the user
func (a *User) GenAccessToken() (string, error) {
	timeNow := utils.Unix()
	if len(a.AccessToken) == 0 {
		a.AccessToken = make(map[string]int64)
	} else {
		// delete expired tokens
		for token, timeout := range a.AccessToken {
			if timeNow-timeout >= consts.AccessTokenLife {
				delete(a.AccessToken, token)
			}
		}

		if len(a.AccessToken) == 5 { // all 5 tokens are valid, delete oldest token
			min := int64(0)
			minToken := ""
			for token, timeout := range a.AccessToken {
				if timeout > min {
					min = timeout
					minToken = token
				}
			}
			delete(a.AccessToken, minToken) // delete the oldest token
		}
	}

	token := utils.GetRandomString(consts.AccessTokenLength)
	a.AccessToken[token] = timeNow

	err := a.Save()
	if err != nil {
		return "", errors.Wrap(err, "could not save user to database")
	}
	return token, nil
}

// AllLogout invalidates the user access token
func (a *User) AllLogout() error {
	for token := range a.AccessToken {
		delete(a.AccessToken, token)
	}
	return a.Save()
}

// AddtoMailbox adds a message to a user's mailbox
func (a *User) AddtoMailbox(subject string, message string) error {
	var x MailboxHelper
	x.Subject = subject
	x.Message = message
	a.Mailbox = append(a.Mailbox, x)
	return a.Save()
}

// VerReq requests account verification
func (a *User) VerReq() error {
	a.VerifyReq = true
	return a.Save()
}

// UnverReq un-requests account verification
func (a *User) UnverReq() error {
	a.VerifyReq = false
	return a.Save()
}
