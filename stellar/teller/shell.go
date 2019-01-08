package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
)

func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Header().Add("Access-Control-Allow-Origin", "localhost")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonString)
}

func ParseInput(input []string) error {
	if len(input) == 0 {
		return fmt.Errorf("No command given")
	}
	if len(input) == 1 {
		log.Println("Command entered: ", input[0])
		command := input[0]
		switch command {
		case "start":
			// idelly this should be a singel command that we can run in order to start all our
			// systems and report back to the platform that we have.
			log.Println("Checking if hardware is installed")
			log.Println("Connecting to said hardware")
		case "test":
			log.Println("COOL STUFF!")
			PingRpc()
		case "getcoins":
			err := xlm.GetXLM(PublicKey)
			if err != nil {
				log.Println(err)
			}
		case "getinv":
			err := PingInvestors()
			if err != nil {
				log.Println(err)
			}
		case "getrec":
			err := PingRecipients()
			if err != nil {
				log.Println(err)
			}
		default:
			// handle defaults here
			log.Println("Invalid command or need more parameters")
		}
	}
	if len(input) == 2 {
		command := input[0]
		subcommand := input[1]
		switch command {
		case "update":
			// the second part should be the state
			// send _timestamp_ stroops to ourselves, we just pay the network fee of 100 stroops
			// this gives us 10**5 updates per xlm, which is pretty nice, considering that we
			// do about 288 updates a day, this amounts to 347 days' worth updates with
			// 1 XLM
			// memo field restricted to 28 bytes - AAAAAAAAAAAAAAAAAAAAAAAAAAAA
			// we can take the first 28 bytes of the hash and then feed it.
			// TODO: we only have 56 bit security, but we the cost of breaking 56 bit security is higher
			// than payment per month, so we should be good.
			x := utils.Timestamp()
			shaHash := strings.ToUpper(utils.SHA3hash(x))
			// we could ideally send the smallest amoutn of 1 stroop but stellar allows you to
			// send yourself as much money as you want, so we can have any number here
			// we could also time this amount to be the state update number itself
			_, hash, err := xlm.SendXLM(PublicKey, utils.I64toS(utils.Unix()), Seed, shaHash[:28])
			if err != nil {
				log.Println(err)
			}
			log.Println("Updated state: ", hash)
		case "gbh":
			hash, err := xlm.GetBlockHash(subcommand)
			if err != nil {
				log.Println(err)
			}
			ColorOutput("THE BLOCKHASH OF GIVEN BLOCK IS: "+hash, MagentaColor)
		case "bal":
			var balance string
			var err error
			ColorOutput("Displaying balance in "+subcommand+" for user: ", WhiteColor)
			switch subcommand {
			case "xlm":
				balance, err = xlm.GetNativeBalance(PublicKey)
			default:
				balance, err = xlm.GetAssetBalance(PublicKey, subcommand)
			}
			if err != nil {
				fmt.Println(err)
				break
			}
			ColorOutput(balance, MagentaColor)
		}
	}
	return nil
}
