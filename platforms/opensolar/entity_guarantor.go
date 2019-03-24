package opensolar

import (
	"fmt"
	"log"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
	"github.com/pkg/errors"
)

// this should contain the future guarantor related functions once we define them concretely

// NewGuarantor returns a new guarantor
func NewGuarantor(uname string, pwd string, seedpwd string, Name string,
	Address string, Description string) (Entity, error) {
	return newEntity(uname, pwd, seedpwd, Name, Address, Description, "guarantor")
}

func (a *Entity) AddFirstLossGuarantee(seedpwd string, amount float64) error {
	a.FirstLossGuarantee = seedpwd
	a.FirstLossGuaranteeAmt = amount
	return a.Save()
}

// have functions in here for the guarantor to cover losses in the case that the recipient does not pay the investors
func CoverFirstLoss(projIndex int, entityIndex int, amount string) error {
	// cover first loss for the project specified
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "could not retrieve projects from database, quitting")
	}

	entity, err := RetrieveEntity(entityIndex)
	if err != nil {
		return errors.Wrap(err, "could not retrieve entity from database, quitting")
	}

	// we now have the entity and the project under question
	if project.Guarantor.U.Index != entity.U.Index {
		return fmt.Errorf("guarantor index does not match with entity's index in database")
	}

	if entity.FirstLossGuaranteeAmt < utils.StoF(amount) {
		log.Println("amount required greater than what guarantor agreed to provide, adjusting first loss to cover for what's available")
		amount = utils.FtoS(entity.FirstLossGuaranteeAmt)
	}
	// we now need to send funds from the gurantor's account to the escrow
	seed, err := wallet.DecryptSeed(entity.U.EncryptedSeed, entity.FirstLossGuarantee) //
	if err != nil {
		return errors.Wrap(err, "could not decrypt seed, quitting!")
	}

	// we now have the seed of the guarantor, shift the money to the escrow now
	escrowPath := CreatePath(consts.EscrowDir, projIndex)
	escrowPubkey, _, err := wallet.RetrieveSeed(escrowPath, consts.EscrowPwd)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve issuer seed")
	}

	// we have the escrow's pubkey, transfer funds to the escrow
	_, txhash, err := assets.SendAsset(consts.Code, consts.StableCoinAddress, escrowPubkey, amount, seed, entity.U.PublicKey, "first loss guarantee")
	if err != nil {
		return errors.Wrap(err, "could not transfer asset to escrow, quitting")
	}

	log.Println("txhash of guarantor kick in:", txhash)

	return nil
}
