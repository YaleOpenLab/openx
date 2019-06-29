package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"strings"

	utils "github.com/Varunram/essentials/utils"
	"github.com/chzyer/readline"
)

// inputcont.go contains all the relevant emulator commands for the contractor

// LoopCont defines a loop for the contractor
func LoopCont(rl *readline.Instance) error {
	// This loop is exclusive to a contractor
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

		err = ParseInputCont(cmdslice)
		if err != nil {
			return errors.Wrap(err, "could not parse input")
		}
	}
}

// ParseInputCont parses input for the contractor
func ParseInputCont(input []string) error {
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
		fmt.Println("ping, display, exchange, ipfs, create, send, receive, originate, " +
			"propose, myproposed, addcollateral, mystage1, mystage0")
	case "ping":
		pingHelper()
	case "display":
		displayHelper(input, LocalContractor.U.Username, LocalContractor.U.Pwhash, "contractor")
	case "exchange":
		exchangeHelper(input, LocalContractor.U.Username, LocalContractor.U.Pwhash, LocalSeed)
	case "ipfs":
		ipfsHelper(input, LocalContractor.U.Username, LocalContractor.U.Pwhash)
	case "send":
		sendHelper(input, LocalContractor.U.Username, LocalContractor.U.Pwhash)
	case "receive":
		receiveHelper(input, LocalContractor.U.Username, LocalContractor.U.Pwhash)
	case "create":
		createHelper(input, LocalContractor.U.Username, LocalContractor.U.Pwhash, LocalContractor.U.StellarWallet.PublicKey)
	case "kyc":
		kycHelper(input, LocalContractor.U.Username, LocalContractor.U.Pwhash, LocalContractor.U.Inspector)
	case "increasetrust":
		increaseTrustHelper(input, LocalContractor.U.Username, LocalContractor.U.Pwhash)
	// Contractor only functions
	case "propose":
		fmt.Println("Proposing a contract can be done only through the opensolar webui" +
			"since that involves document verification")
		// end of propose
	case "myproposed":
		x, err := GetStage2Contracts(LocalContractor.U.Username, LocalContractor.U.Pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
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

		response, err := AddCollateral(LocalContractor.U.Username, LocalContractor.U.Pwhash, collateral, amount)
		if err != nil {
			log.Println(err)
			break
		}

		if response.Code == 200 {
			ColorOutput("SUCCESSFULLY ADDED COLLATERAL", GreenColor)
		} else {
			ColorOutput("RESPONSE STATUS: "+utils.ItoS(response.Code), GreenColor)
		}
	case "mystage0":
		x, err := GetStage0Contracts(LocalContractor.U.Username, LocalContractor.U.Pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
	case "mystage1": // if the contractor acts as an originator sometime. Bool setting would be weird,
		// but I guess there's nothing that prevents a contractor from acting as an originator, so we allow this.
		x, err := GetStage1Contracts(LocalContractor.U.Username, LocalContractor.U.Pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		PrintProjects(x)
	}
	return nil
}
