package database

import (
	"encoding/json"
	"fmt"
	"log"
	"syscall"

	aes "github.com/YaleOpenLab/smartPropertyMVP/stellar/aes"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
	"github.com/boltdb/bolt"
	"github.com/stellar/go/keypair"
	"golang.org/x/crypto/ssh/terminal"
)

// TODO: this strucutre assumes one platform for all assets, do we have a
// master st ruct which houses all these platforms for additional seed security?
// in this case, we could give the seed to a lawyer who can enforce the platform's
// behaviour in case of any dispute.
type Platform struct {
	// TODO: theoretically, we don't need this structure at all, since we can get the pubkey
	// anytime we want, so remove it
	Index uint32
	// ideally there should only be one platform
	PublicKey string
	// the publickey of the platform
	// the seed isn't stored in the database, so the only way
	// to access the seed would be through GetSeedFromEncryptedSeed
	DateInitiated string
	// date when the platform was created
	DateRestored string
	// date the platform was restored from its seed, useufl for auditing if a crash
	// did happen
	// We could have a multisig like scheme for hte platform between various
	// stakeholders to restore confidence that the platform is doing the right thing
	// as well, need to implement it the right way
}

func EncryptAndStoreSeed(seed string, password string) {
	// handler to store the seed over at seed.hex
	// person either needs to store this file and remember the password or has to
	// remember the seed in order to access the platform again
	aes.EncryptFile("seed.hex", []byte(seed), password)
	if seed != string(aes.DecryptFile("seed.hex", password)) {
		// somethign wrong wiht encryption, exit
		log.Fatal("Encrpytion and decryption seeds don't match, exiting!")
	}
	fmt.Println("Successfully encrypted your seed as seed.hex")
}

func NewPlatform() (Platform, error) {
	// this function is used to generate a new keypair which will be assigned to
	// the platform
	var nPlatform Platform
	var nPlatformSeed string // init eparately since we don't store this
	var err error
	nPlatform.Index = uint32(1) // only one platform, so this is fine
	nPlatformSeed, nPlatform.PublicKey, err = xlm.GetKeyPair()
	log.Printf("\nTHE PLATFORM SEED IS: %s\nAND YOUR PUBLIC KEY IS: %s\nKEEP IT SUPER SAFE OR YOU MIGHT NOT HAVE ACCESS TO THESE FUNDS AGAIN \n", nPlatformSeed, nPlatform.PublicKey)
	nPlatform.DateInitiated = utils.Timestamp()
	fmt.Println("Enter a password to encrypt your platform's master seed. Please store this in a very safe place. This prompt will not ask to confirm your password")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	fmt.Println() // needed to separate user input
	EncryptAndStoreSeed(nPlatformSeed, password) // store the seed in a secure location
	err = xlm.GetXLM(nPlatform.PublicKey)
	return nPlatform, err
}

func GetSeedFromEncryptedSeed(encrypted string, password string) (string) {
	// this function must be used for any handling within the code written here
	return string(aes.DecryptFile(encrypted, password))
}

func RestorePlatformFromSeed(seed string) (Platform, error) {
	// this function should be used when the platform admin remembers the seed but
	// does not possess the  encrypted file. The seed is what's needed to access
	// the account, so we don't restrict access
	var rPlatform Platform
	keyp, err := keypair.Parse(seed)
	if err != nil {
		return rPlatform, err
	} else {
		rPlatform.PublicKey = keyp.Address()
		rPlatform.DateRestored = utils.Timestamp()
	}
	fmt.Println("Enter a password to encrypt your platform's master seed. Please store this in a very safe place. This prompt will not ask to confirm your password")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	fmt.Println()
	EncryptAndStoreSeed(seed, password)
	log.Printf("restored platform public key %s from seed", rPlatform.PublicKey)
	return rPlatform, err
}

func RestorePlatformFromFile(path string, password string) (Platform, error){
	return RestorePlatformFromSeed(string(aes.DecryptFile(path, password)))
}

func InsertPlatform(a Platform) error {
	// inserts the PublicKey into the database to keep track of the PublicKey
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(PlatformBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(utils.Uint32toB(a.Index)), encoded)
	})
	return err
}

func RetrievePlatform() (Platform, error) {
	// retrieves the platforms (more like the publickey)
	var rPlatform Platform
	db, err := OpenDB()
	if err != nil {
		return rPlatform, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(PlatformBucket)
		i := uint32(1)
		x := b.Get(utils.Uint32toB(i))
		if x == nil {
			// this is where the key does not exist
			return nil
		}
		err := json.Unmarshal(x, &rPlatform)
		if err != nil {
			return nil
		}
		return nil
	})
	return rPlatform, err
}