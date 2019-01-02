package platform

import (
	"fmt"
	"log"
	"os"

	aes "github.com/YaleOpenLab/smartPropertyMVP/stellar/aes"
	consts "github.com/YaleOpenLab/smartPropertyMVP/stellar/consts"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
	"github.com/stellar/go/keypair"
)

// TODO: this structure assumes one platform for all assets, do we have a
// master st ruct which houses all these platforms for additional seed security?
// in this case, we could give the seed to a lawyer who can enforce the platform's
// behaviour in case of any dispute.
// TODO: should  this have its won database for security reasons?
type Platform struct {
	// TODO: theoretically, we don't need this structure at all, since we can get the pubkey
	// anytime we want, so remove it
	Index int
	// ideally there should only be one platform
	PublicKey string
	// the publickey of the platform
	// the seed isn't stored in the database, so the only way
	// to access the seed would be through GetSeedFromEncryptedSeed
}

// EncryptAndStoreSeed encrypts and stores the seed of the platform in a file
// named platformseed.hex at the root directory
func EncryptAndStoreSeed(seed string, password string) {
	// handler to store the seed over at platformseed.hex
	// person either needs to store this file and remember the password or has to
	// remember the seed in order to access the platform again
	aes.EncryptFile(consts.HomeDir+"/platformseed.hex", []byte(seed), password)
	if seed != string(aes.DecryptFile(consts.HomeDir+"/platformseed.hex", password)) {
		// something wrong with encryption, exit
		log.Fatal("Encryption and decryption seeds don't match, exiting!")
	}
	fmt.Println("Successfully encrypted your seed as platformseed.hex")
}

// NewPlatform creates a new platform and returns the platform struct
func NewPlatform() (string, string, error) {
	// this function is used to generate a new keypair which will be assigned to
	// the platform
	var seed string
	var publicKey string
	var err error

	seed, publicKey, err = xlm.GetKeyPair()
	if err != nil {
		return publicKey, seed, err
	}
	log.Printf("\nTHE PLATFORM SEED IS: %s\nAND YOUR PUBLIC KEY IS: %s\nKEEP IT SUPER SAFE OR YOU MIGHT NOT HAVE ACCESS TO THESE FUNDS AGAIN \n", seed, publicKey)
	fmt.Println("Enter a password to encrypt your platform's master seed. Please store this in a very safe place. This prompt will not ask to confirm your password")
	password, err := utils.ScanRawPassword()
	if err != nil {
		return publicKey, seed, err
	}
	EncryptAndStoreSeed(seed, password) // store the seed in a secure location
	err = xlm.GetXLM(publicKey)         // get funds for our platform
	return publicKey, seed, err
}

// GetSeedFromEncryptedSeed gets the unencrypted seed from the encrypted file
// stored on disk with the help of the password.
func GetSeedFromEncryptedSeed(encrypted string, password string) string {
	// this function must be used for any handling within the code written here
	return string(aes.DecryptFile(encrypted, password))
}

// GetPlatformPublicKeyAndStoreSeed restores the platform struct when passed the seed
func GetPlatformPublicKeyAndStoreSeed(seed string) (string, error) {
	// this function should be used when the platform admin remembers the seed but
	// does not possess the encrypted file. The seed is what's needed to access
	// the account, so we don't restrict access
	var publicKey string
	keyp, err := keypair.Parse(seed)
	if err != nil {
		return publicKey, err
	} else {
		publicKey = keyp.Address()
	}
	log.Println("ENTER A PASSWORD TO ENCRYPT YOUR SEED")
	password, err := utils.ScanRawPassword()
	if err != nil {
		return publicKey, err
	}
	EncryptAndStoreSeed(seed, password)
	return publicKey, nil
}

// RestorePlatformFromFile restores the platform struct directly from the file
func GetPlatformFromFile(path string, password string) (string, string, error) {
	var publicKey string
	var seed string
	seed = string(aes.DecryptFile(path, password))
	keyp, err := keypair.Parse(seed)
	if err != nil {
		return publicKey, seed, err
	} else {
		publicKey = keyp.Address()
	}
	return publicKey, seed, nil // this is only for the public key
}

// InitializePlatform returns the platform structure and the seed
func InitializePlatform() (string, string, error) {
	var publicKey string
	var seed string
	var err error

	if _, err := os.Stat(consts.HomeDir); os.IsNotExist(err) {
		// directory does not exist, create one
		log.Println(consts.HomeDir)
		log.Println("Creating home directory")
		os.MkdirAll(consts.HomeDir, os.ModePerm)
	}
	// now we can be sure we have the directory, check for seed
	if _, err := os.Stat(consts.HomeDir + "/platformseed.hex"); !os.IsNotExist(err) {
		// the seed exists
		fmt.Println("ENTER YOUR PASSWORD TO DECRYPT THE SEED FILE")
		password, err := utils.ScanRawPassword()
		if err != nil {
			log.Println(err)
			return publicKey, seed, err
		}
		publicKey, seed, err = GetPlatformFromFile(consts.HomeDir+"/platformseed.hex", password)
		return publicKey, seed, err
	}
	// platform doesn't exist or user doesn't have encrypted file. Ask
	fmt.Println("DO YOU HAVE YOUR RAW SEED? IF SO, ENTER SEED. ELSE ENTER N")
	seed, err = utils.ScanForString()
	if err != nil {
		log.Println(err)
		return publicKey, seed, err
	}
	if seed == "N" || seed == "n" {
		// the user doesn't have seed, so create a new platform
		publicKey, seed, err = NewPlatform()
		if err != nil {
			log.Println(err)
			return publicKey, seed, err
		}
	} else {
		// user has given us a seed, validate
		publicKey, err = GetPlatformPublicKeyAndStoreSeed(seed)
		if err != nil {
			log.Println(err)
			return publicKey, seed, err
		}
	}
	return publicKey, seed, nil
}
