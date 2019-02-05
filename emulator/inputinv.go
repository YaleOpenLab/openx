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

// inputinv.go contains all the relevant emulator commands for the investor
// we have one parse input function for each entity on the platform
// ie investor, recipient, contractor, originator and potentially more depednign upon usage
// the input array contains the commands that we want to parse.
// first check the length of the input array and then define accordingly

func LoopInv() error {
	// this loop is for an investor
	// we have authenticated the user and stored the details in an appropriate structure
	// need to repeat this struct everywhere because having separate functions and importing
	// it doesn't seem to work
	// the problem with having a conditional statement inside the loop is that it checks
	// role each time and that's not nice performance wise
	// TOOD: look at alternatives if possible
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

		err = ParseInputInv(cmdslice)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func ParseInputInv(input []string) error {
	var err error
	if len(input) == 0 {
		// shouldn't happen, still
		return fmt.Errorf("Length of input array is zero, quitting!")
	}
	// input is greater than length 1 which means we can parse according to the command given
	command := input[0]
	switch command {
	case "help":
		fmt.Println("LIST OF SUPPORTED COMMANDS: ")
		fmt.Println("ping, display, exchange, ipfs, vote, kyc, invest, create, send, receive")
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
				balances, err := GetBalances(LocalInvestor.U.Username, LocalInvestor.U.Pwhash)
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
				balance, err := GetXLMBalance(LocalInvestor.U.Username, LocalInvestor.U.Pwhash)
				if err != nil {
					log.Println(err)
					break
				}
				ColorOutput("BALANCE: "+balance, MagentaColor)
				break
			case "all":
				balances, err := GetBalances(LocalInvestor.U.Username, LocalInvestor.U.Pwhash)
				if err != nil {
					log.Println(err)
					break
				}
				PrintBalances(balances)
				break
			default:
				balance, err := GetAssetBalance(LocalInvestor.U.Username, LocalInvestor.U.Pwhash, subcommand)
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
			PrintInvestor(LocalInvestor)
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
		amount, err := utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		// convert this to int and check if int
		fmt.Println("Exchanging", amount, "XLM for STABLEUSD")
		response, err := GetStableCoin(LocalInvestor.U.Username, LocalInvestor.U.Pwhash, LocalSeed, input[1])
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
			log.Println("<ipfs> string")
			break
		}
		inputString := input[1]
		hashString, err := GetIpfsHash(LocalInvestor.U.Username, LocalInvestor.U.Pwhash, inputString)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println("IPFS HASH", hashString)
		// end of ipfs
		// start cases which are unique to investor
	case "vote":
		if len(input) != 3 {
			log.Println("vote <projIndex> <amount>")
			break
		}
		_, err = utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		_, err = utils.StoICheck(input[2])
		if err != nil {
			log.Println(err)
			break
		}
		status, err := VoteTowardsProject(input[1], input[2], LocalInvestor.U.Username, LocalInvestor.U.Pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		if status.Status == 200 {
			ColorOutput("VOTE CAST!", GreenColor)
		} else {
			ColorOutput("VOTE NOT CAST", RedColor)
		}
		break
		// end of vote
	case "kyc":
		if !LocalInvestor.U.Inspector {
			ColorOutput("YOU ARE NOT A KYC INSPECTOR", RedColor)
			break
		}
		if len(input) == 1 {
			log.Println("kyc <auth, view>")
		}
		subcommand := input[1]
		switch subcommand {
		case "auth":
			if len(input) != 3 {
				log.Println("kyc auth <userIndex>")
				break
			}
			_, err = utils.StoICheck(input[1])
			if err != nil {
				log.Println(err)
				break
			}
			status, err := AuthKyc(input[1], LocalInvestor.U.Username, LocalInvestor.U.Pwhash)
			if err != nil {
				log.Println(err)
				break
			}
			if status.Status == 200 {
				ColorOutput("USER KYC'D!", GreenColor)
			} else {
				ColorOutput("USER NOT KYC'D", RedColor)
			}
			break
		case "notdone":
			users, err := NotKycView(LocalInvestor.U.Username, LocalInvestor.U.Pwhash)
			if err != nil {
				log.Println(err)
				break
			}
			PrintUsers(users)
			// print all the users who have kyc'd
			break
		case "done":
			users, err := KycView(LocalInvestor.U.Username, LocalInvestor.U.Pwhash)
			if err != nil {
				log.Println(err)
				break
			}
			PrintUsers(users)
			// print all the users who have kyc'd
			break
		}
		// end of kyc
	case "invest":
		log.Println("Invest Params: invest <proj_number> <amount>")
		if len(input) != 3 {
			log.Println("Extra / less params passed, please check!")
			break
		}
		_, err = utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		_, err = utils.StoICheck(input[2])
		if err != nil {
			log.Println(err)
			break
		}
		// now we need to invest in this project, call RPC
		status, err := InvestInProject(input[1], input[2], LocalInvestor.U.Username, LocalInvestor.U.Pwhash, LocalSeedPwd)
		if err != nil {
			log.Println(err)
			break
		}
		if status.Status == 200 {
			ColorOutput("INVESTMENT SUCCESSFUL, CHECK EMAIL", GreenColor)
		} else {
			ColorOutput("INVESTMENT NOT SUCCESSFUL", RedColor)
		}
		break // end of invest
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
			status, err := CreateAssetInv(LocalInvestor.U.Username, LocalInvestor.U.Pwhash,
				assetName, LocalInvestor.U.PublicKey)
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

			txhash, err := SendLocalAsset(LocalInvestor.U.Username, LocalInvestor.U.Pwhash,
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
			txhash, err := SendXLM(LocalInvestor.U.Username, LocalInvestor.U.Pwhash,
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
			status, err := AskXLM(LocalInvestor.U.Username, LocalInvestor.U.Pwhash)
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

			status, err := TrustAsset(LocalInvestor.U.Username, LocalInvestor.U.Pwhash, assetName, issuerPubkey, limit, LocalSeedPwd)
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
	} // end of investor
	return nil
}
