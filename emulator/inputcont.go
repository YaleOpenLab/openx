package main

import (
	"fmt"
	"log"
	"strings"

	consts "github.com/YaleOpenLab/openx/consts"
	utils "github.com/YaleOpenLab/openx/utils"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

// inputcont.go contains all the relevant emulator commands for the contractor
func LoopCont() error {
	// This loop is exclusive to a contractor
	promptColor := color.New(color.FgHiYellow).SprintFunc()
	whiteColor := color.New(color.FgHiWhite).SprintFunc()
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      promptColor("emulator") + whiteColor("# "),
		HistoryFile: consts.TellerHomeDir + "/history.txt",
		// AutoComplete: lc.NewAutoCompleter(),
	})

	ColorOutput("YOUR SEED IS: "+LocalSeed, RedColor)

	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	for {
		// setup reader with max 4K input chars
		msg, err := rl.Readline()
		if err != nil {
			log.Println(err)
			return err
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
			log.Println(err)
			return err
		}
	}
	return nil
}

func ParseInputCont(input []string) error {
	var err error
	// Various command supported for the recipient
	if len(input) == 0 {
		// shouldn't happen, still
		return fmt.Errorf("Length of input array is zero, quitting!")
	}
	// input is greater than length 1 which means we can parse according to the command given
	command := input[0]
	switch command {
	case "help":
		fmt.Println("LIST OF SUPPORTED COMMANDS: ")
		fmt.Println("ping, display, exchange, ipfs, create, send, receive, originate, " +
			"propose, myproposed, addcollateral, myoriginated, mypreoriginated")
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
		createHelper(input, LocalContractor.U.Username, LocalContractor.U.Pwhash, LocalContractor.U.PublicKey)
	case "kyc":
		kycHelper(input, LocalContractor.U.Username, LocalContractor.U.Pwhash, LocalContractor.U.Inspector)
	// Contractor only functions
	case "propose":
		fmt.Println("Proposing a contract can be done only through the opensolar webui" +
			"since that involves document verification")
		break
		// end of propose
	case "myproposed":
		x, err := GetProposedContracts(LocalContractor.U.Username, LocalContractor.U.Pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
		break
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
		break
	case "mypreoriginated":
		x, err := GetPreOriginatedContracts(LocalContractor.U.Username, LocalContractor.U.Pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
		break
		// end of myoriginated
	case "myoriginated": // if the contractor acts as an originator sometime. Bool setting would be weird,
		// but I guess there's nothing that prevents a contractor from acting as an originator, so we allow this.
		x, err := GetOriginatedContracts(LocalContractor.U.Username, LocalContractor.U.Pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		PrintProjects(x)
		break
		// end of myoriginated
		// end of originate
	}
	return nil
}
