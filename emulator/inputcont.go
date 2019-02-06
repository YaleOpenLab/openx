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
				balances, err := GetBalances(LocalContractor.U.Username, LocalContractor.U.Pwhash)
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
				balance, err := GetXLMBalance(LocalContractor.U.Username, LocalContractor.U.Pwhash)
				if err != nil {
					log.Println(err)
					break
				}
				ColorOutput("BALANCE: "+balance, MagentaColor)
				break
			case "all":
				balances, err := GetBalances(LocalContractor.U.Username, LocalContractor.U.Pwhash)
				if err != nil {
					log.Println(err)
					break
				}
				PrintBalances(balances)
				break
			default:
				balance, err := GetAssetBalance(LocalContractor.U.Username, LocalContractor.U.Pwhash, subcommand)
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
			PrintEntity(LocalContractor)
			break
		case "projects":
			if len(input) != 3 {
				// only display was given, so display help command
				log.Println("<display projects> <preorigin, origin, seed, proposed, open, funded, installed, power, fin>")
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
			log.Println("<exchange> amount")
			break
		}
		_, err = utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		// convert this to int and check if int
		amount := input[1]

		fmt.Println("Exchanging", amount, "XLM for STABLEUSD")
		response, err := GetStableCoin(LocalContractor.U.Username, LocalContractor.U.Pwhash, LocalSeed, amount)
		if err != nil {
			log.Println(err)
			break
		}
		if response.Code == 200 {
			ColorOutput("SUCCESSFUL, CHECK BALANCES", GreenColor)
		} else {
			ColorOutput("RESPONSE STATUS: "+utils.ItoS(response.Code), GreenColor)
		}
		// end of exchange
	case "ipfs":
		if len(input) != 2 {
			log.Println("<ipfs> string")
			break
		}
		inputString := input[1]
		hashString, err := GetIpfsHash(LocalContractor.U.Username, LocalContractor.U.Pwhash, inputString)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println("IPFS HASH", hashString)
		// end of ipfs
		// start of conttractor only functions
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
			status, err := CreateAssetInv(LocalContractor.U.Username, LocalContractor.U.Pwhash,
				assetName, LocalContractor.U.PublicKey)
			if err != nil {
				log.Println(err)
				break
			}
			if status.Code == 200 {
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

			txhash, err := SendLocalAsset(LocalContractor.U.Username, LocalContractor.U.Pwhash,
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
			txhash, err := SendXLM(LocalContractor.U.Username, LocalContractor.U.Pwhash,
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
			status, err := AskXLM(LocalContractor.U.Username, LocalContractor.U.Pwhash)
			if err != nil {
				log.Println(err)
				break
			}
			if status.Code == 200 {
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

			status, err := TrustAsset(LocalContractor.U.Username, LocalContractor.U.Pwhash, assetName, issuerPubkey, limit, LocalSeedPwd)
			if err != nil {
				log.Println(err)
				break
			}
			if status.Code == 200 {
				ColorOutput("COIN REQUEST SUCCESSFUL, CHECK EMAIL", GreenColor)
			} else {
				ColorOutput("COIN REQUEST NOT SUCCESSFUL", RedColor)
			}
			break
		} // end of receive
	}
	return nil
}
