package opensolar

import (
	"log"
	"os"

	consts "github.com/YaleOpenLab/openx/consts"
	assets "github.com/YaleOpenLab/openx/assets"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
	xlm "github.com/YaleOpenLab/openx/xlm"
	"github.com/pkg/errors"
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
func CreatePath(path string, projIndex int) string {
	return path + utils.ItoS(projIndex) + ".key"
}

// CreateFile creates a new empty keyfile
func CreateFile(issuerPath string, projIndex int) string {
	path := CreatePath(issuerPath, projIndex)
	// we need to create this file
	os.Create(path)
	return path
}

// InitEscrow creates a new keypair and stores it in a file
func InitEscrow(issuerPath string, projIndex int, seedpwd string) error {
	// init a new pk and seed pair
	seed, pubkey, err := xlm.GetKeyPair()
	if err != nil {
		return errors.Wrap(err, "Error while generating keypair")
	}
	// store this seed in home/projects/projIndex.hex
	// we need a password for encrypting the seed
	path := CreateFile(issuerPath, projIndex)
	err = wallet.StoreSeed(seed, seedpwd, path)
	if err != nil {
		return errors.Wrap(err, "Error while storing seed")
	}

	_, txhash, err := xlm.SendXLMCreateAccount(pubkey, "100", consts.PlatformSeed) // pass the platform seed to be the account that seeds the escrow
	if err != nil {
		return errors.Wrap(err, "Error while sending xlm to create account")
	}
	log.Printf("Txhash for setting up Project Issuer for project %d is %s", projIndex, txhash)
	_, txhash, err = xlm.SetAuthImmutable(seed)
	if err != nil {
		return errors.Wrap(err, "Error while setting auth immutable on account")
	}
	log.Printf("Txhash for setting Auth Immutable on project %d is %s", projIndex, txhash)

	// create a trustline with the stablecoin
	txhash, err = assets.TrustAsset(consts.Code, consts.StablecoinPublicKey, "10000000000", pubkey, seed)
	if err != nil {
		return err
	}

	log.Println("TRUST HASH FOR ESCROW TRUSTING STABLECOIN: ", txhash)
	return nil
}

// DeleteIssuer deletes the keyfile
// But this is not needed since once the account is frozen, an attacker who does
// have access to the seed can not aim to achieve anything since the account is locked
func DeleteEscrow(issuerPath string, projIndex int) error {
	path := CreatePath(issuerPath, projIndex)
	return os.Remove(path)
}
