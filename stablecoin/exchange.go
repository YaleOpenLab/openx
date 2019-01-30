package stablecoin

import (
	"fmt"
	"log"

	assets "github.com/OpenFinancing/openfinancing/assets"
	consts "github.com/OpenFinancing/openfinancing/consts"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
)

// TODO: in case the person does not have enough stablecoin on hand to invest but has xlm,
// we must offer to exchange their xlm for stablecoin and enable it by default.
func Exchange(recipientPK string, recipientSeed string, convAmount string) error {

	if !xlm.AccountExists(recipientPK) {
		return fmt.Errorf("Account does not exist, quitting!")
	}

	hash, err := assets.TrustAsset(Code, consts.StableCoinAddress, consts.StablecoinTrustLimit, recipientPK, recipientSeed)
	if err != nil {
		return err
	}
	log.Println("tx hash for trusting stableUSD: ", hash)
	// now send coins across and see if our tracker detects it
	_, hash, err = xlm.SendXLM(PublicKey, convAmount, recipientSeed, "Sending XLM to bootstrap")
	if err != nil {
		return err
	}
	log.Println("tx hash for sent xlm: ", hash, "pubkey: ", recipientPK)
	return nil
}
