package models

import (
	"fmt"
	"log"
	"time"

	assets "github.com/YaleOpenLab/openx/assets"
	consts "github.com/YaleOpenLab/openx/consts"
	utils "github.com/YaleOpenLab/openx/utils"
	wallet "github.com/YaleOpenLab/openx/wallet"
	xlm "github.com/YaleOpenLab/openx/xlm"
)

// the models package won't be imported directly in any place but would be imported
// by all the investment models that exist
func SendUSDToPlatform(invSeed string, invAmount string, memo string) (string, error) {
	// send stableusd to the platform (not the issuer) since the issuer will be locked
	// and we can't use the funds. We also need ot be able to redeem the stablecoin for fiat
	// so we can't burn them

	invPubkey, err := wallet.ReturnPubkey(invSeed)
	if err != nil {
		return "", err
	}

	var oldPlatformBalance string
	oldPlatformBalance, err = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.Code)
	if err != nil {
		// platform does not have stablecoin, shouldn't arrive here ideally
		oldPlatformBalance = "0"
	}

	_, txhash, err := assets.SendAsset(consts.Code, consts.StablecoinPublicKey, consts.PlatformPublicKey, invAmount, invSeed, invPubkey, memo)
	if err != nil {
		log.Println("Sending stableusd to platform failed", consts.PlatformPublicKey, invAmount, invSeed, invPubkey)
		return txhash, err
	}

	log.Println("Sent STABLEUSD to platform, confirmation: ", txhash)
	time.Sleep(5 * time.Second) // wait for a block

	newPlatformBalance, err := xlm.GetAssetBalance(consts.PlatformPublicKey, consts.Code)
	if err != nil {
		return txhash, err
	}

	if utils.StoF(newPlatformBalance)-utils.StoF(oldPlatformBalance) < utils.StoF(invAmount)-1 {
		return txhash, fmt.Errorf("Sent amount doesn't match with investment amount")
	}
	return txhash, nil
}
