package opensolar

import (
	"log"
	"os"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	multisig "github.com/YaleOpenLab/openx/multisig"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
	xlm "github.com/YaleOpenLab/openx/xlm"
	"github.com/pkg/errors"
)

// escrow defines the escrow of asset for a specific project. We should generate a
// new seed and public key pair for each project that is at stage 3, so this would be
// automated at that stage. Once an investor has finished investing in the project,
// we need to send the recipient DebtAssets and then set all weights to zero in order
// to lock the account and prevent any further transactions from being authorized.
// One can stil send fund to the frozen account but the account can not use them
// this serves our purpose since we only want receipt of debt assets and want to freeze
// issuance so that anybody who hacks us can not print more tokens.

// In financial terms, an escrow is a special purpose vehicle (kind of cool that we have SPV in finance)
// CreatePath returns the path of a specific project
func CreatePath(path string, projIndex int) string {
	return path + utils.ItoS(projIndex) + ".key"
}

// CreateFile creates a new empty keyfile
func CreateFile(escrowPath string, projIndex int) string {
	path := CreatePath(escrowPath, projIndex)
	// we need to create this file
	os.Create(path)
	return path
}

// InitEscrow creates a new keypair and stores it in a file
func InitEscrow(escrowPath string, projIndex int, seedpwd string, recpPubkey string, mySeed string) error {
	// init a new pk and seed pair
	// TODO: replace with the escrow here
	pubkey, err := initMultisigEscrow(recpPubkey)
	if err != nil {
		return errors.Wrap(err, "error while initalizing multisig escrow, quitting!")
	}

	// define two seeds that are needed for signing transactions from the escrow
	seed1 := consts.PlatformSeed
	seed2 := mySeed

	// store this seed in home/projects/projIndex.hex
	// we need a password for encrypting the seed
	path := CreateFile(escrowPath, projIndex)
	err = wallet.StoreSeed(recpPubkey, seedpwd, path) // store the recipient's pubkey encrypted with a blob. TODO: we're doing this for privacy reasons (a person accessing the server can not see whose escrow this is)
	if err != nil {
		return errors.Wrap(err, "Error while storing seed")
	}

	_, txhash, err := xlm.SendXLMCreateAccount(pubkey, "100", consts.PlatformSeed) // pass the platform seed to be the account that seeds the escrow
	if err != nil {
		return errors.Wrap(err, "Error while sending xlm to create account")
	}
	log.Printf("Txhash for setting up Project escrow for project %d is %s", projIndex, txhash)

	err = multisig.AuthImmutable2of2(pubkey, seed1, seed2)
	if err != nil {
		return errors.Wrap(err, "could not set auth immutable on account, quitting!")
	}

	multisig.TrustAssetTx(consts.Code, consts.StablecoinPublicKey, "10000000000", pubkey, seed1, seed2)
	if err != nil {
		return errors.Wrap(err, "could not trust stablecoin, quitting!")
	}

	log.Println("TRUST HASH FOR ESCROW TRUSTING STABLECOIN: ", txhash)
	return nil
}

// Deleteescrow deletes the keyfile
// But this is not needed since once the account is frozen, an attacker who does
// have access to the seed can not aim to achieve anything since the account is locked
func DeleteEscrow(escrowPath string, projIndex int) error {
	path := CreatePath(escrowPath, projIndex)
	return os.Remove(path)
}

func TransferFundsToEscrow(amount float64, projIndex int) error {
	// we need to transfer funds that hte investors invested in the platform to the specific escrow
	escrowPath := CreatePath(consts.EscrowDir, projIndex)
	escrowPubkey, _, err := wallet.RetrieveSeed(escrowPath, consts.EscrowPwd)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve escrow seed")
	}

	// we have the wallet pubkey, transfer funds to the escrow now
	_, txhash, err := assets.SendAsset(consts.Code, consts.StablecoinPublicKey, escrowPubkey,
		utils.FtoS(amount), consts.PlatformSeed, consts.PlatformPublicKey, "escrow init")
	if err != nil {
		return errors.Wrap(err, "could not fund escrow, quitting!")
	}

	log.Println("tx hash for funding project escrow is: ", txhash)

	return nil
}

// InitMultisigEscrow initializes a multisig escrow with one signer as the recipient and the other as the platform
func initMultisigEscrow(pubkey1 string) (string, error) {
	// recpPubkey is the public key of the recipient
	// the seed of the escrow is needed to init the first tx that will change options
	pubkey2 := consts.PlatformPublicKey
	// we now have the two public keys that are needed to authorize this transaction. Construct a 2of2 multisig
	return multisig.New2of2(pubkey1, pubkey2)
}
