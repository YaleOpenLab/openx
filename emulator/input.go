package main

import (
	"fmt"
	"log"

	utils "github.com/OpenFinancing/openfinancing/utils"
)

// TODO: reduce code duplication here, we can't create an interface and handle stuff
// since interfaces can't be indexed.
func ParseInputInv(input []string) error {
	var err error
	// we need to have one parse input function for each entity on the platform
	// ie investor, recipient, contractor, originator, etc
	// the input array contains the commands that we want to parse.
	// first check the length of the input array and then define accordingly
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
				balances, err := GetBalances(LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword)
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
				balance, err := GetXLMBalance(LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword)
				if err != nil {
					log.Println(err)
					break
				}
				ColorOutput("BALANCE: "+balance, MagentaColor)
				break
			case "all":
				balances, err := GetBalances(LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword)
				if err != nil {
					log.Println(err)
					break
				}
				PrintBalances(balances)
				break
			default:
				balance, err := GetAssetBalance(LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword, subcommand)
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
		response, err := GetStableCoin(LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword, LocalSeed, input[1])
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
		hashString, err := GetIpfsHash(LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword, inputString)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println("IPFS HASH", hashString)
		// end of ipfs
		// start cases which are unique to investor
	case "vote":
		// TODO
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
		status, err := VoteTowardsProject(input[1], input[2], LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword)
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
		// TODO
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
			status, err := AuthKyc(input[1], LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword)
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
			users, err := NotKycView(LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword)
			if err != nil {
				log.Println(err)
				break
			}
			PrintUsers(users)
			// print all the users who have kyc'd
			break
		case "done":
			users, err := KycView(LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword)
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
		status, err := InvestInProject(input[1], input[2], LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword, LocalSeedPwd)
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
			status, err := CreateAssetInv(LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword,
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

			txhash, err := SendLocalAsset(LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword,
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
			txhash, err := SendXLM(LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword,
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
			status, err := AskXLM(LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword)
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

			status, err := TrustAsset(LocalInvestor.U.LoginUserName, LocalInvestor.U.LoginPassword, assetName, issuerPubkey, limit, LocalSeedPwd)
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
				balances, err := GetBalances(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword)
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
				balance, err := GetXLMBalance(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword)
				if err != nil {
					log.Println(err)
					break
				}
				ColorOutput("BALANCE: "+balance, MagentaColor)
				break
			case "all":
				balances, err := GetBalances(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword)
				if err != nil {
					log.Println(err)
					break
				}
				PrintBalances(balances)
				break
			default:
				balance, err := GetAssetBalance(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword, subcommand)
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
		response, err := GetStableCoin(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword, LocalSeed, input[1])
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
		hashString, err := GetIpfsHash(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword, inputString)
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
		status, err := UnlockProject(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword, LocalSeedPwd, input[1])
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
		status, err := Payback(projIndex, LocalSeedPwd, LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword, assetName, amount)
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

		status, err := FinalizeProject(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword, projIndex)
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

		status, err := OriginateProject(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword, projIndex)
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

			status, err := CreateAssetInv(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword,
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

			txhash, err := SendLocalAsset(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword,
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
			txhash, err := SendXLM(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword,
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
			status, err := AskXLM(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword)
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

			status, err := TrustAsset(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword, assetName, issuerPubkey, limit, LocalSeedPwd)
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
				balances, err := GetBalances(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword)
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
				balance, err := GetXLMBalance(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword)
				if err != nil {
					log.Println(err)
					break
				}
				ColorOutput("BALANCE: "+balance, MagentaColor)
				break
			case "all":
				balances, err := GetBalances(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword)
				if err != nil {
					log.Println(err)
					break
				}
				PrintBalances(balances)
				break
			default:
				balance, err := GetAssetBalance(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword, subcommand)
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
		response, err := GetStableCoin(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword, LocalSeed, amount)
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
		hashString, err := GetIpfsHash(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword, inputString)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println("IPFS HASH", hashString)
		// end of ipfs
		// start of conttractor only functions
	case "propose": // TODO: maybe arrive at a way to do this sometime later?
		fmt.Println("Proposing a contract can be done only through the opensolar webui" +
			"since that involves document verification")
		break
		// end of propose
	case "myproposed":
		x, err := GetProposedContracts(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword)
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

		response, err := AddCollateral(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword, collateral, amount)
		if err != nil {
			log.Println(err)
			break
		}

		if response.Status == 200 {
			ColorOutput("SUCCESSFULLY ADDED COLLATERAL", GreenColor)
		} else {
			ColorOutput("RESPONSE STATUS: "+utils.ItoS(response.Status), GreenColor)
		}
		break
	case "mypreoriginated":
		x, err := GetPreOriginatedContracts(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
		break
		// end of myoriginated
	case "myoriginated": // if the contractor acts as an originator sometime. Bool setting would be weird,
		// but I guess there's nothing that prevents a contractor from acting as an originator, so we allow this.
		x, err := GetOriginatedContracts(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword)
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
			status, err := CreateAssetInv(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword,
				assetName, LocalContractor.U.PublicKey)
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

			txhash, err := SendLocalAsset(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword,
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
			txhash, err := SendXLM(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword,
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
			status, err := AskXLM(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword)
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

			status, err := TrustAsset(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword, assetName, issuerPubkey, limit, LocalSeedPwd)
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
	}
	return nil
}

func ParseInputOrig(input []string) error {
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
		fmt.Println("ping, display, exchange, ipfs, create, send, receive, propose, " +
			"preoriginate, myproposed, addcollateral, myoriginated, mypreoriginated")
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
				balances, err := GetBalances(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword)
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
				balance, err := GetXLMBalance(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword)
				if err != nil {
					log.Println(err)
					break
				}
				ColorOutput("BALANCE: "+balance, MagentaColor)
				break
			case "all":
				balances, err := GetBalances(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword)
				if err != nil {
					log.Println(err)
					break
				}
				PrintBalances(balances)
				break
			default:
				balance, err := GetAssetBalance(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword, subcommand)
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
			PrintEntity(LocalOriginator)
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

		response, err := GetStableCoin(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword, LocalSeed, amount)
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
		hashString, err := GetIpfsHash(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword, inputString)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println("IPFS HASH", hashString)
		// end of ipfs
	case "propose":
		fmt.Println("Proposing a contract can be done only through the opensolar webui" +
			"since that involves document verification")
		break
	case "preoriginate":
		fmt.Println("Pre originating a contract can be done only through the opensolar webui" +
			"since that involves document verification")
		break
		// end of propose
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

		response, err := AddCollateral(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword, collateral, amount)
		if err != nil {
			log.Println(err)
			break
		}

		if response.Status == 200 {
			ColorOutput("SUCCESSFULLY ADDED COLLATERAL", GreenColor)
		} else {
			ColorOutput("RESPONSE STATUS: "+utils.ItoS(response.Status), GreenColor)
		}
		break
		// end of addcollateral
	case "myproposed":
		x, err := GetProposedContracts(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
		break
		// end of myproposed
	case "mypreoriginated":
		x, err := GetPreOriginatedContracts(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
		break
		// end of myoriginated
	case "myoriginated":
		x, err := GetOriginatedContracts(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
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
			status, err := CreateAssetInv(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword,
				assetName, LocalOriginator.U.PublicKey)
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

			txhash, err := SendLocalAsset(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword,
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
			txhash, err := SendXLM(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword,
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
			status, err := AskXLM(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword)
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

			status, err := TrustAsset(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword, assetName, issuerPubkey, limit, LocalSeedPwd)
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
	}
	return nil
}
