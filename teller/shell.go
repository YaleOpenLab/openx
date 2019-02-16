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

func ParseInput(input []string) {
	if len(input) == 0 {
		fmt.Println("List of commands: ping, receive, display, update")
		return
	}

	command := input[0]
	switch command {
	case "qq":
		// handler to quit and test the teller without hashing the state and committing two transactions
		// each time we start the teller
		log.Fatal("cool")
	case "help":
		fmt.Println("List of commands: ping, receive, display, update")
		break
	case "ping":
		err := PingRpc()
		if err != nil {
			log.Println(err)
		}
	case "receive":
		if len(input) != 2 {
			fmt.Println("USAGE: receive xlm")
			return
		}
		err := xlm.GetXLM(RecpPublicKey)
		if err != nil {
			log.Println(err)
		}
	case "display":
		if len(input) < 3 {
			fmt.Println("USAGE: display <balance>")
			return
		}
		subcommand := input[1]
		switch subcommand {
		case "balance":
			if len(input) < 4 {
				fmt.Println("USAGE: display balance <xlm, asset>")
				return
			}

			subsubcommand := input[2]
			var balance string
			var err error
			ColorOutput("Displaying balance in "+subsubcommand+" for user: ", WhiteColor)

			switch subsubcommand {
			case "xlm":
				balance, err = xlm.GetNativeBalance(RecpPublicKey)
			default:
				balance, err = xlm.GetAssetBalance(RecpPublicKey, subcommand)
			}

			if err != nil {
				log.Println(err)
				return
			}
			ColorOutput(balance, MagentaColor)
		default:
			// handle defaults here
			log.Println("Invalid command or need more parameters")
		} // end of display
	case "update":
		if len(input) != 2 {
			fmt.Println("USAGE: update <state>")
			return
		}
		subcommand := input[1]
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
	}
}
