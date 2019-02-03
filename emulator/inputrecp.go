package main

import (
	"fmt"
	"log"
	"strings"

	consts "github.com/OpenFinancing/openfinancing/consts"
	utils "github.com/OpenFinancing/openfinancing/utils"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

// inputrecp.go contains all the relevant emulator commands for the recipient
func LoopRecp() error {
	// This loop is exclusive to a recipient
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

		err = ParseInputRecp(cmdslice)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func ParseInputRecp(input []string) error {
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
		fmt.Println("ping, display, exchange, ipfs, create, send, receive, unlock, payback, finalize, originate")
		break
	case "ping":
		err = PingRpc()
		if err != nil {
			log.Println(err)
			return err
		}
	// have a list of all the supported commands here and then check for length, etc inside
	case "display":
		// display is a  broad command and needs to have a subcommand
		if len(input) == 1 {
			// only display was given, so display help command
			log.Println("<display><balance, profile, projects>")
			break
		}
		subcommand := input[1]
		switch subcommand {
		case "balance":
			if len(input) == 2 {
				log.Println("Calling balances API")
				balances, err := GetBalances(LocalRecipient.U.Username, LocalRecipient.U.Pwhash)
				if err != nil {
					log.Println(err)
					break
				}
				PrintBalances(balances)
				break
			}
			subcommand := input[2]
			switch subcommand {
			case "xlm":
				// print xlm balance
				balance, err := GetXLMBalance(LocalRecipient.U.Username, LocalRecipient.U.Pwhash)
				if err != nil {
					log.Println(err)
					break
				}
				ColorOutput("BALANCE: "+balance, MagentaColor)
				break
			case "all":
				balances, err := GetBalances(LocalRecipient.U.Username, LocalRecipient.U.Pwhash)
				if err != nil {
					log.Println(err)
					break
				}
				PrintBalances(balances)
				break
			default:
				balance, err := GetAssetBalance(LocalRecipient.U.Username, LocalRecipient.U.Pwhash, subcommand)
				if err != nil {
					log.Println(err)
					break
				}
				ColorOutput("BALANCE: "+balance, MagentaColor)
				break
				// print asset balance
			}
		case "profile":
			log.Println("Displaying Profile")
			PrintRecipient(LocalRecipient)
			break
		case "projects":
			if len(input) != 3 {
				// only display was given, so display help command
				log.Println("display projects <preorigin, origin, seed, proposed, open, funded, installed, power, fin>")
				break
			}
			subsubcommand := input[2]
			var stage float64
			switch subsubcommand {
			case "preorigin":
				log.Println("Displaying all pre-originated (stage 0) projects")
				stage = 0
				break
			case "origin":
				log.Println("Displaying all originated (stage 1) projects")
				stage = 1
				break
			case "seed":
				log.Println("Displaying all seed (stage 1.5) projects")
				stage = 1.5
				break
			case "proposed":
				log.Println("Displaying all proposed (stage 2) projects")
				stage = 2
				break
			case "open":
				log.Println("Displaying open (stage 3) projects")
				stage = 3
				break
			case "funded":
				log.Println("Displaying funded (stage 4) projects")
				stage = 4
				break
			case "installed":
				log.Println("Displaying installed (stage 5) projects")
				stage = 5
				break
			case "power":
				log.Println("Displaying funded (stage 6) projects")
				stage = 6
				break
			case "fin":
				log.Println("Displaying funded (stage 7) projects")
				stage = 7
				break
			}
			arr, err := RetrieveProject(stage)
			if err != nil {
				log.Println(err)
				break
			}
			PrintProjects(arr)
			break
		} // end of display
	case "exchange":
		if len(input) != 2 {
			// only display was given, so display help command
			log.Println("exchange <amount>")
			break
		}
		amount, err := utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		// convert this to int and check if int
		fmt.Println("Exchanging", amount, "XLM for STABLEUSD")
		response, err := GetStableCoin(LocalRecipient.U.Username, LocalRecipient.U.Pwhash, LocalSeed, input[1])
		if err != nil {
			log.Println(err)
			break
		}
		if response.Status == 200 {
			ColorOutput("SUCCESSFUL, CHECK BALANCES", GreenColor)
		} else {
			ColorOutput("RESPONSE STATUS: "+utils.ItoS(response.Status), GreenColor)
		}
		// end of exchange
	case "ipfs":
		if len(input) != 2 {
			log.Println("ipfs <string>")
			break
		}
		inputString := input[1]
		hashString, err := GetIpfsHash(LocalRecipient.U.Username, LocalRecipient.U.Pwhash, inputString)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println("IPFS HASH", hashString)
		// end of ipfs
		// start of recipient only functions
	case "unlock":
		if len(input) != 2 {
			log.Println("unlock <projIndex>")
			break
		}
		_, err = utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		status, err := UnlockProject(LocalRecipient.U.Username, LocalRecipient.U.Pwhash, LocalSeedPwd, input[1])
		if err != nil {
			log.Println(err)
			break
		}
		if status.Status == 200 {
			ColorOutput("PAYBACK SUCCESSFUL, CHECK EMAIL", GreenColor)
		} else {
			ColorOutput("PAYBACK NOT SUCCESSFUL", RedColor)
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
		if status.Status == 200 {
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
		if status.Status == 200 {
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
		if status.Status == 200 {
			ColorOutput("PAYBACK SUCCESSFUL, CHECK EMAIL", GreenColor)
		} else {
			ColorOutput("PAYBACK NOT SUCCESSFUL", RedColor)
		}
		break
		// end of originate
	case "create":
		// create enables you to create tokens on stellar that you can excahnge with third parties.
		if len(input) == 1 {
			log.Println("create <asset>")
			break
		}
		subcommand := input[1]
		switch subcommand {
		case "asset":
			// create a new asset
			if len(input) != 3 {
				log.Println("create asset <name>")
				break
			}

			assetName := input[2]

			status, err := CreateAssetInv(LocalRecipient.U.Username, LocalRecipient.U.Pwhash,
				assetName, LocalRecipient.U.PublicKey)
			if err != nil {
				log.Println(err)
				break
			}
			if status.Status == 200 {
				ColorOutput("INVESTMENT SUCCESSFUL, CHECK EMAIL", GreenColor)
			} else {
				ColorOutput("INVESTMENT NOT SUCCESSFUL", RedColor)
			}
		} // end of create
	case "send":
		if len(input) == 1 {
			log.Println("send <asset>")
			break
		}
		subcommand := input[1]
		switch subcommand {
		case "asset":
			if len(input) != 5 {
				log.Println("send asset <assetName> <destination> <amount>")
				break
			}

			assetName := input[2]
			destination := input[3]
			amount := input[4]

			txhash, err := SendLocalAsset(LocalRecipient.U.Username, LocalRecipient.U.Pwhash,
				LocalSeedPwd, assetName, destination, amount)
			if err != nil {
				log.Println(err)
			}
			ColorOutput("TX HASH: "+txhash, MagentaColor)
			break
			// end of asset
		case "xlm":
			if len(input) < 4 {
				log.Println("send xlm <destination> <amount> <<memo>>")
				break
			}
			destination := input[2]
			_, err = utils.StoFWithCheck(input[3])
			if err != nil {
				log.Println(err)
				break
			}
			// send xlm overs
			amount := input[3]
			var memo string
			if len(input) > 4 {
				memo = input[4]
			}
			txhash, err := SendXLM(LocalRecipient.U.Username, LocalRecipient.U.Pwhash,
				LocalSeedPwd, destination, amount, memo)
			if err != nil {
				log.Println(err)
			}
			ColorOutput("TX HASH: "+txhash, MagentaColor)
		}
		// end of send
	case "receive":
		// we can either receive from the faucet or trust issuers to receive assets
		if len(input) == 1 {
			log.Println("receive <xlm, asset>")
			break
		}
		subcommand := input[1]
		switch subcommand {
		case "xlm":
			status, err := AskXLM(LocalRecipient.U.Username, LocalRecipient.U.Pwhash)
			if err != nil {
				log.Println(err)
				break
			}
			if status.Status == 200 {
				ColorOutput("COIN REQUEST SUCCESSFUL, CHECK EMAIL", GreenColor)
			} else {
				ColorOutput("COIN REQUEST NOT SUCCESSFUL", RedColor)
			}
			// ask for coins from the faucet
		case "asset":
			if len(input) != 5 {
				log.Println("receive asset <assetName> <issuerPubkey> <limit>")
				break
			}

			assetName := input[2]
			issuerPubkey := input[3]
			_, err = utils.StoFWithCheck(input[4])
			if err != nil {
				log.Println(err)
				break
			}

			limit := input[4]

			status, err := TrustAsset(LocalRecipient.U.Username, LocalRecipient.U.Pwhash, assetName, issuerPubkey, limit, LocalSeedPwd)
			if err != nil {
				log.Println(err)
				break
			}
			if status.Status == 200 {
				ColorOutput("COIN REQUEST SUCCESSFUL, CHECK EMAIL", GreenColor)
			} else {
				ColorOutput("COIN REQUEST NOT SUCCESSFUL", RedColor)
			}
			break
		} // end of receive
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
