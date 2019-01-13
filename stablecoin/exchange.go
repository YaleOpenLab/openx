package stablecoin

import (
	assets "github.com/YaleOpenLab/smartPropertyMVP/stellar/assets"
	consts "github.com/YaleOpenLab/smartPropertyMVP/stellar/consts"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
	"log"
)

func Exchange(recipientPK string, recipientSeed string, convAmount string) error {
	hash, err := assets.TrustAsset(StableUSD, consts.StablecoinTrustLimit, recipientPK, recipientSeed)
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
