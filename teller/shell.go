package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	utils "github.com/YaleOpenLab/openx/utils"
	xlm "github.com/YaleOpenLab/openx/xlm"
)

func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Header().Add("Access-Control-Allow-Origin", "localhost")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonString)
}

func ParseInput(input []string) error {
	if len(input) == 0 {
		fmt.Println("List of commands: test, getcoins, getinv, getrec, update, gbh, bal, xlm")
		return fmt.Errorf("No command given")
	}
	if len(input) == 1 {
		log.Println("Command entered: ", input[0])
		command := input[0]
		switch command {
		case "help":
			// print out all the commands, would be nice if there was some automated
			// way to get all the commands
			fmt.Println("List of commands: test, getcoins, getinv, getrec, update, gbh, bal, xlm")
			break
		case "test":
			log.Println("COOL STUFF!")
			err := PingRpc()
			if err != nil {
				log.Println(err)
			}
		case "getcoins":
			err := xlm.GetXLM(RecpPublicKey)
			if err != nil {
				log.Println(err)
			}
		case "getinv":
			err := GetInvestors()
			if err != nil {
				log.Println(err)
			}
		case "getrec":
			err := GetRecipients()
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
			if subcommand == "help" {
				fmt.Println("USAGE: update <state>")
				fmt.Println("THis hashes the state ")
				break
			}
			// the second part should be the state
			// send _timestamp_ stroops to ourselves, we just pay the network fee of 100 stroops
			// this gives us 10**5 updates per xlm, which is pretty nice, considering that we
			// do about 288 updates a day, this amounts to 347 days' worth updates with
			// 1 XLM
			// memo field restricted to 28 bytes - AAAAAAAAAAAAAAAAAAAAAAAAAAAA
			// we can take the first 28 bytes of the hash and then feed it.
			// we only have 56 bit security, but we e cost of breaking 56 bit security is higher
			// than the payment required per month, so we should be good
			shaHash := strings.ToUpper(utils.SHA3hash(subcommand))
			// we could ideally send the smallest amount of 1 stroop but stellar allows you to
			// send yourself as much money as you want, so we can have any number here
			// we could also time this amount to be the state update number itself
			_, hash, err := xlm.SendXLM(RecpPublicKey, utils.I64toS(utils.Unix()), RecpSeed, shaHash[:28])
			if err != nil {
				log.Println(err)
			}
			ColorOutput("Updated State: "+hash, MagentaColor)
		case "gbh":
			if subcommand == "help" {
				fmt.Println("USAGE: gbh <block number>")
				break
			}
			hash, err := xlm.GetBlockHash(subcommand)
			if err != nil {
				log.Println(err)
			}
			ColorOutput("THE BLOCKHASH OF GIVEN BLOCK IS: "+hash, MagentaColor)
		case "bal":
			var balance string
			var err error
			if subcommand == "help" {
				fmt.Println("USAGE: bal <asset>, xlm for native balance")
				break
			}
			ColorOutput("Displaying balance in "+subcommand+" for user: ", WhiteColor)
			switch subcommand {
			case "xlm":
				balance, err = xlm.GetNativeBalance(RecpPublicKey)
			default:
				balance, err = xlm.GetAssetBalance(RecpPublicKey, subcommand)
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
