package escrow

import (
	"log"

	stablecoin "github.com/Varunram/essentials/crypto/stablecoin"
	assets "github.com/Varunram/essentials/crypto/xlm/assets"
	multisig "github.com/Varunram/essentials/crypto/xlm/multisig"
	wallet "github.com/Varunram/essentials/crypto/xlm/wallet"
	"github.com/pkg/errors"
)

// escrow implements an escrow based off Stellar

// InitEscrow creates a new keypair and stores it in a file
func InitEscrow(projIndex int, seedpwd string, recpPubkey string, mySeed string, otherSeed string) (string, error) {
	otherPubkey, err := wallet.ReturnPubkey(otherSeed)
	if err != nil {
		return "", errors.Wrap(err, "could not get pubkey from seed")
	}

	pubkey, err := initMultisigEscrow(recpPubkey, otherPubkey)
	if err != nil {
		return pubkey, errors.Wrap(err, "error while initializing multisig escrow, quitting!")
	}

	log.Println("successfully initialized multisig escrow")
	// define two seeds that are needed for signing transactions from the escrow
	seed1 := otherSeed
	seed2 := mySeed

	log.Println("stored escrow pubkey successfully")
	err = multisig.AuthImmutable2of2(pubkey, seed1, seed2)
	if err != nil {
		return pubkey, errors.Wrap(err, "could not set auth immutable on account, quitting!")
	}

	log.Println("set auth immutable on account successfully")
	multisig.TrustAssetTx(stablecoin.StablecoinCode, stablecoin.StablecoinPublicKey, "10000000000", pubkey, seed1, seed2)
	if err != nil {
		return pubkey, errors.Wrap(err, "could not trust stablecoin, quitting!")
	}

	return pubkey, nil
}

// TransferFundsToEscrow transfers stablecoin to the escrow address from otherSeed
func TransferFundsToEscrow(amount float64, projIndex int, escrowPubkey string, otherSeed string) error {
	// we have the wallet pubkey, transfer funds to the escrow now
	_, txhash, err := assets.SendAsset(stablecoin.StablecoinCode, stablecoin.StablecoinPublicKey, escrowPubkey,
		amount, otherSeed, "escrow init")
	if err != nil {
		return errors.Wrap(err, "could not fund escrow, quitting!")
	}

	log.Println("tx hash for funding project escrow is: ", txhash)
	return nil
}

// InitMultisigEscrow initializes a multisig escrow
func initMultisigEscrow(pubkey1 string, pubkey2 string) (string, error) {
	return multisig.New2of2(pubkey1, pubkey2)
}

// SendFundsFromEscrow sends funds to a destination address from the project escrow
func SendFundsFromEscrow(escrowPubkey string, destination string, signer1 string, signer2 string, amount float64, memo string) error {
	return multisig.Tx2of2(escrowPubkey, destination, signer1, signer2, amount, memo)
}
