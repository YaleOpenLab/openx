package main

import (
	"fmt"
	"log"
	"strings"

	scan "github.com/Varunram/essentials/scan"
	xlm "github.com/Varunram/essentials/xlm"
	assets "github.com/Varunram/essentials/xlm/assets"
	multisig "github.com/Varunram/essentials/xlm/multisig"
	consts "github.com/YaleOpenLab/openx/consts"
)

// rescue mode contains a list of handlers that can be used when we need to login as any account and perform emergency resuce fns
// call this after starting the platform so we don't have to do the boring stuff again if we're the platform
// this will no have any sort of input parsing since only admins will be using it in the event of an emergency

// RescueMode starts openx in rescue mode as the platform in order to be able to salvage funds
// in case of a hack / disruption.
func RescueMode() {
	seed := consts.PlatformSeed
	pubkey := consts.PlatformPublicKey
	fmt.Println(`
		██████╗ ███████╗███████╗ ██████╗██╗   ██╗███████╗    ███╗   ███╗ ██████╗ ██████╗ ███████╗
		██╔══██╗██╔════╝██╔════╝██╔════╝██║   ██║██╔════╝    ████╗ ████║██╔═══██╗██╔══██╗██╔════╝
		██████╔╝█████╗  ███████╗██║     ██║   ██║█████╗       ██╔████╔██║██║   ██║██║  ██║█████╗
		██╔══██╗██╔══╝  ╚════██║██║     ██║   ██║██╔══╝       ██║╚██╔╝██║██║   ██║██║  ██║██╔══╝
		██║  ██║███████╗███████║╚██████╗╚██████╔╝███████╗    ██║ ╚═╝ ██║╚██████╔╝██████╔╝███████╗
		╚═╝  ╚═╝╚══════╝╚══════╝ ╚═════╝ ╚═════╝ ╚══════╝    ╚═╝     ╚═╝ ╚═════╝ ╚═════╝ ╚══════╝
		`)
	for {
		log.Println("WELCOME TO RESCUE MODE")
		log.Println("LIST OF FUNCTIONS AVAILABLE")
		log.Println("1. Send funds to another address")
		log.Println("2. Sweep funds from platform")
		log.Println("3. View platform balances")
		log.Println("4. Sweep XLM from project escrow")
		log.Println("5. Send USD to another address")
		choice, err := scan.Int()
		if err != nil {
			log.Fatal(err)
		}
		switch choice {
		case 1:
			log.Println("Enter address")
			address, err := scan.String()
			if err != nil {
				log.Println("!!!" + strings.ToUpper(err.Error()) + "!!!")
				break
			}
			log.Println("Enter amount")
			amount, err := scan.Float()
			if err != nil {
				log.Println("!!!" + strings.ToUpper(err.Error()) + "!!!")
				break
			}
			log.Println("address: ", address, "amount: ", amount)
			if !xlm.AccountExists(address) {
				_, _, err = xlm.SendXLMCreateAccount(address, amount, seed)
				if err != nil {
					log.Println("!!!" + strings.ToUpper(err.Error()) + "!!!")
					break
				}
			} else {
				_, _, err = xlm.SendXLM(address, amount, seed, "rescue mode")
				if err != nil {
					log.Println("!!!" + strings.ToUpper(err.Error()) + "!!!")
					break
				}
			}

		case 2:
			log.Println("Enter sweep address")
			address, err := scan.String()
			if err != nil {
				log.Println("!!!" + strings.ToUpper(err.Error()) + "!!!")
				break
			}

			amount := xlm.GetNativeBalance(pubkey)
			log.Println("SWEEP AMOUNT IS: ", amount)
			// send the tx over
			_, _, err = xlm.SendXLM(address, amount-5, seed, "rescue mode sweep")
			if err != nil {
				log.Println("error while transferring funds to secondary account, quitting")
				break
			}

		case 3:
			balances, err := xlm.GetAllBalances(pubkey)
			if err != nil {
				log.Println("!!!" + strings.ToUpper(err.Error()) + "!!!")
				continue
			}
			log.Println(balances)
		case 4:
			log.Println("Enter escrow address")
			source, err := scan.String()
			if err != nil {
				log.Println("!!!" + strings.ToUpper(err.Error()) + "!!!")
				break
			}
			amount := xlm.GetNativeBalance(source)
			log.Println("Enter sweep address")
			destination, err := scan.String()
			if err != nil {
				log.Println("!!!" + strings.ToUpper(err.Error()) + "!!!")
				break
			}
			log.Println("Enter other signer's seed")
			otherSeed, err := scan.String()
			if err != nil {
				log.Println("!!!" + strings.ToUpper(err.Error()) + "!!!")
				break
			}
			err = multisig.Tx2of2(source, destination, otherSeed, seed, amount, "escrow sweep")
			if err != nil {
				log.Println("!!!" + strings.ToUpper(err.Error()) + "!!!")
				break
			}
		case 5:
			log.Println("Trustline must already exiust for specified seed")
			log.Println("Enter address")
			address, err := scan.String()
			if err != nil {
				log.Println("!!!" + strings.ToUpper(err.Error()) + "!!!")
				break
			}

			log.Println("Enter amount")
			amount, err := scan.Float()
			if err != nil {
				log.Println("!!!" + strings.ToUpper(err.Error()) + "!!!")
				break
			}
			log.Println("AMOUNT IS: ", amount)

			balance := xlm.GetAssetBalance(pubkey, consts.AnchorUSDCode)
			log.Println("Available balance: ", balance)
			if balance < amount {
				log.Println("insufficient balance")
				break
			}

			_, _, err = assets.SendAsset(consts.AnchorUSDCode, consts.AnchorUSDAddress, address, amount,
				seed, "tf")
			if err != nil {
				log.Println(err)
				break
			}
		}
	}
}
