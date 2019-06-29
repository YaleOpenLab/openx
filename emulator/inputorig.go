package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"strings"

	utils "github.com/Varunram/essentials/utils"
	"github.com/chzyer/readline"
)

// inputorig.go contains all the relevant emulator commands for the originator

// LoopOrig loops over originator input
func LoopOrig(rl *readline.Instance) error {
	// This loop is exclusive to an originator
	for {
		// setup reader with max 4K input chars
		msg, err := rl.Readline()
		if err != nil {
			return errors.Wrap(err, "could not read user input")
		}
		msg = strings.TrimSpace(msg)
		if len(msg) == 0 {
			continue
		}
		rl.SaveHistory(msg)

		cmdslice := strings.Fields(msg)
		ColorOutput("entered command: "+msg, YellowColor)

		err = ParseInputOrig(cmdslice)
		if err != nil {
			return errors.Wrap(err, "could not parse user input")
		}
	}
}

// ParseInputOrig parses orig input
func ParseInputOrig(input []string) error {
	var err error
	// Various command supported for the recipient
	if len(input) == 0 {
		// shouldn't happen, still
		return errors.New("Length of input array is zero, quitting!")
	}
	// input is greater than length 1 which means we can parse according to the command given
	command := input[0]
	switch command {
	case "help":
		fmt.Println("LIST OF SUPPORTED COMMANDS: ")
		fmt.Println("ping, display, exchange, ipfs, create, send, receive, propose, " +
			"newstage0, myproposed, addcollateral, mystage1, mystage0")
	case "ping":
		pingHelper()
	case "display":
		displayHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash, "originator")
	case "exchange":
		exchangeHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash, LocalSeed)
	case "ipfs":
		ipfsHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash)
	case "send":
		sendHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash)
	case "receive":
		receiveHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash)
	case "create":
		createHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash, LocalOriginator.U.StellarWallet.PublicKey)
	case "kyc":
		kycHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash, LocalOriginator.U.Inspector)
	case "increasetrust":
		increaseTrustHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash)
	// Originator only functions
	case "propose":
		fmt.Println("Proposing a contract can be done only through the opensolar webui" +
			"since that involves document verification")
	case "newstage0":
		fmt.Println("Pre-originating a contract can be done only through the opensolar webui" +
			"since that involves document verification")
	case "addcollateral":
		if len(input) != 3 {
			log.Println("<addcollateral> collateral amount")
			break
		}

		_, err = utils.StoFWithCheck(input[2])
		if err != nil {
			log.Println(err)
			break
		}

		collateral := input[1]
		amount := input[2]

		response, err := AddCollateral(LocalOriginator.U.Username, LocalOriginator.U.Pwhash, collateral, amount)
		if err != nil {
			log.Println(err)
			break
		}

		if response.Code == 200 {
			ColorOutput("SUCCESSFULLY ADDED COLLATERAL", GreenColor)
		} else {
			ColorOutput("RESPONSE STATUS: "+utils.ItoS(response.Code), GreenColor)
		}
	case "myproposed":
		x, err := GetStage2Contracts(LocalOriginator.U.Username, LocalOriginator.U.Pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
	case "mystage0":
		x, err := GetStage0Contracts(LocalOriginator.U.Username, LocalOriginator.U.Pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
	case "mystage1":
		x, err := GetStage1Contracts(LocalOriginator.U.Username, LocalOriginator.U.Pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
	}
	return nil
}
