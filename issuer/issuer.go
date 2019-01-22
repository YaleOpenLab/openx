package issuer

import (
	"log"
	"os"

	consts "github.com/OpenFinancing/openfinancing/consts"
	utils "github.com/OpenFinancing/openfinancing/utils"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
)

// issuer defines the issuer of asset for a specific project. We should generate a
// new seed and public key pair for each project that is at stage 3, so this would be
// automated at that stage. Once an investor has finished investing in the project,
// we need to send the recipient DebtAssets and then set all weights to zero in order
// to lock the account and prevent any further transactions from being authorized.
// One can stil send fund to the frozen account but the account can not use them
// this serves our purpose since we only want receipt of debt assets and want to freeze
// issuance so that anybody who hacks us can not print more tokens

// CreatePath returns the path of a specific project
func CreatePath(projIndex int) string {
	return consts.OpenSolarIssuerDir + utils.ItoS(projIndex) + ".key"
}

// CreateFile creates a new empty keyfile
func CreateFile(projIndex int) string {
	path := CreatePath(projIndex)
	// we need to create this file
	os.Create(path)
	return path
}

// InitIssuer creates a new keypair and stores it in a file
func InitIssuer(projIndex int, seedpwd string) error {
	// init a new pk and seed pair
	seed, _, err := xlm.GetKeyPair()
	if err != nil {
		log.Println(err)
		return err
	}
	// store this seed in home/projects/projIndex.hex
	// we need a password for encrypting the seed
	path := CreateFile(projIndex)
	err = wallet.StoreSeed(seed, seedpwd, path)
	if err != nil {
		return err
	}
	return nil
}

// DeleteIssuer deletes the keyfile
// But this is not needed since once the account is frozen, an attacker who does
// have access to the seed can not aim to achieve anything since the account is locked
func DeleteIssuer(projIndex int) error {
	path := CreatePath(projIndex)
	return os.Remove(path)
}

// FundIssuer funds the issuer using the platform's seed. This is the only way we
// can know if a specific account is an issuer or not. Also helps accountability
// of how many issuers are in existence and how many projects have been invested in.
func FundIssuer(projIndex int, seedpwd string, platformSeed string) error {
	// need to read the seed from the file using the seedpwd
	path := CreatePath(projIndex)
	pubkey, seed, err := wallet.RetrieveSeed(path, seedpwd)
	if err != nil {
		return err
	}
	log.Printf("Project Index: %d, Seed: %s, Address: %s", projIndex, seed, pubkey)
	// func SendXLMCreateAccount(destination string, amount string, Seed string) (int32, string, error) {
	_, txhash, err := xlm.SendXLMCreateAccount(pubkey, "100", platformSeed)
	if err != nil {
		return err
	}
	log.Printf("Txhash for setting up Project Issuer for project %d is %s", projIndex, txhash)
	_, txhash, err = xlm.SetAuthImmutable(seed)
	if err != nil {
		return err
	}
	log.Printf("Txhash for setting Auth Immutable on project %d is %s", projIndex, txhash)
	return nil
}

// FreezeIssuer freezes the issuer account and makes it not possible for someone
// to sign new transactions or send funds from the given account
func FreezeIssuer(projIndex int, seedpwd string) (string, error) {
	path := CreatePath(projIndex)
	_, seed, err := wallet.RetrieveSeed(path, seedpwd)
	if err != nil {
		return "", nil
	}
	_, txhash, err := xlm.FreezeAccount(seed)
	if err != nil {
		return "", nil
	}
	log.Println("Tx hash for freezing account is: ", txhash)
	return txhash, nil
}
