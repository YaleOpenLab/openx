package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"strings"

	utils "github.com/YaleOpenLab/openx/utils"
	"github.com/chzyer/readline"
)

// inputrecp.go contains all the relevant emulator commands for the recipient
func LoopRecp(rl *readline.Instance) error {
	// This loop is exclusive to a recipient
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

		err = ParseInputRecp(cmdslice)
		if err != nil {
			return errors.Wrap(err, "could not parse user input")
		}
	}
	return nil
}

func ParseInputRecp(input []string) error {
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
		fmt.Println("ping, display, exchange, ipfs, create, send, receive, unlock, payback, finalize, originate")
	case "ping":
		pingHelper()
	case "display":
		displayHelper(input, LocalRecipient.U.Username, LocalRecipient.U.Pwhash, "recipient")
	case "exchange":
		exchangeHelper(input, LocalRecipient.U.Username, LocalRecipient.U.Pwhash, LocalSeed)
	case "ipfs":
		ipfsHelper(input, LocalRecipient.U.Username, LocalRecipient.U.Pwhash)
	case "send":
		sendHelper(input, LocalRecipient.U.Username, LocalRecipient.U.Pwhash)
	case "receive":
		receiveHelper(input, LocalRecipient.U.Username, LocalRecipient.U.Pwhash)
	case "create":
		createHelper(input, LocalRecipient.U.Username, LocalRecipient.U.Pwhash, LocalRecipient.U.PublicKey)
	case "kyc":
		kycHelper(input, LocalRecipient.U.Username, LocalRecipient.U.Pwhash, LocalRecipient.U.Inspector)
	// Recipient ONly functions
	case "unlock":
		if len(input) < 3 {
			log.Println("unlock <projIndex> <platform>")
			break
		}
		_, err = utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		platform := input[2]
		switch platform {
		case "opensolar":
			status, err := UnlockOpenSolar(LocalRecipient.U.Username, LocalRecipient.U.Pwhash, LocalSeedPwd, input[1])
			if err != nil {
				log.Println(err)
				break
			}
			if status.Code == 200 {
				ColorOutput("PAYBACK SUCCESSFUL, CHECK EMAIL", GreenColor)
			} else {
				ColorOutput("PAYBACK NOT SUCCESSFUL", RedColor)
			}
			break
		case "opzones":
			if len(input) < 4 {
				log.Println("unlock <projIndex> opzones <cbond, lucoop>")
				break
			}
			model := input[3]
			switch model {
			case "cbond":
				log.Println("CGOND OKS")
				status, err := UnlockCBond(LocalRecipient.U.Username, LocalRecipient.U.Pwhash, LocalSeedPwd, input[1])
				if err != nil {
					log.Println(err)
					break
				}
				if status.Code == 200 {
					ColorOutput("PAYBACK SUCCESSFUL, CHECK EMAIL", GreenColor)
				} else {
					ColorOutput("PAYBACK NOT SUCCESSFUL", RedColor)
				}
				break
			} // end of model switch
		}
		break
	case "payback":
		if len(input) != 3 {
			log.Println("payback <projIndex> <amount>")
			break
		}
		_, err = utils.StoICheck(input[1]) // projectIndex
		if err != nil {
			log.Println(err)
			break
		}
		_, err = utils.StoICheck(input[2]) // amount
		if err != nil {
			log.Println(err)
			break
		}

		projIndex := input[1]
		amount := input[2]

		assetName := LocalRecipient.ReceivedSolarProjects[0] // hardcode for now, TODO: change this
		status, err := Payback(projIndex, LocalSeedPwd, LocalRecipient.U.Username, LocalRecipient.U.Pwhash, assetName, amount)
		if err != nil {
			log.Println(err)
			break
		}
		if status.Code == 200 {
			ColorOutput("PAYBACK SUCCESSFUL, CHECK EMAIL", GreenColor)
		} else {
			ColorOutput("PAYBACK NOT SUCCESSFUL", RedColor)
		}
		break
		// end of payback
	case "finalize":
		if len(input) != 2 {
			log.Println("finalize <projIndex>")
			break
		}

		_, err = utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}

		projIndex := input[1]

		status, err := FinalizeProject(LocalRecipient.U.Username, LocalRecipient.U.Pwhash, projIndex)
		if err != nil {
			log.Println(err)
			break
		}
		if status.Code == 200 {
			ColorOutput("PAYBACK SUCCESSFUL, CHECK EMAIL", GreenColor)
		} else {
			ColorOutput("PAYBACK NOT SUCCESSFUL", RedColor)
		}
		break
	case "originate":
		if len(input) != 2 {
			log.Println("originate <projIndex>")
			break
		}

		_, err = utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}

		projIndex := input[1]

		status, err := OriginateProject(LocalRecipient.U.Username, LocalRecipient.U.Pwhash, projIndex)
		if err != nil {
			log.Println(err)
			break
		}
		if status.Code == 200 {
			ColorOutput("PAYBACK SUCCESSFUL, CHECK EMAIL", GreenColor)
		} else {
			ColorOutput("PAYBACK NOT SUCCESSFUL", RedColor)
		}
		break
		// end of originate
	case "calculate":
		if len(input) == 1 {
			log.Println("calculate <payback>")
			break
		}
		subcommand := input[1]
		switch subcommand {
		case "ownership":
			// calculate the balance of the debt asset here
			if len(input) == 1 {
				fmt.Println("payback <assetName>")
				break
			}

			assetName := LocalRecipient.ReceivedSolarProjects[0]

			limit, err := GetTrustLimit(LocalRecipient.U.Username, LocalRecipient.U.Pwhash, assetName)
			if err != nil {
				log.Println(err)
				break
			}

			limitF := utils.StoF(limit)
			// get balance of debt asset here
			debtBalance, err := GetAssetBalance(LocalRecipient.U.Username, LocalRecipient.U.Pwhash, assetName)
			if err != nil {
				log.Println(err)
				break
			}

			debtF := utils.StoF(debtBalance)
			ownership := (1 - debtF/(limitF/2)) * 100
			ColorOutput("YOUR PERCENTAGE OWNERSHIP OF THE ASSET: "+utils.FtoS(ownership), MagentaColor)
			break
			// end of payback
		}
	}
	return nil
}
