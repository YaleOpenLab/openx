package main

import (
	"fmt"
	"log"

	utils "github.com/OpenFinancing/openfinancing/utils"
)

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
		if len(input) != 2 {
			log.Println("kyc <userIndex>")
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
		break
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
			log.Println("<ipfs> string")
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
		/*
			err := solar.UnlockProject(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword, 1, LocalSeedPwd)
			log.Fatal(err)
		*/
		// PAYBACK LOOP DOES NOT RUN
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
		assetName := LocalRecipient.ReceivedSolarProjects[0] // hardcode for now, TODO: change this
		status, err := Payback(input[1], LocalSeedPwd, LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword, assetName, input[2])
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
		status, err := FinalizeProject(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword, input[1])
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
		status, err := OriginateProject(LocalRecipient.U.LoginUserName, LocalRecipient.U.LoginPassword, input[1])
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
		amount, err := utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		// convert this to int and check if int
		fmt.Println("Exchanging", amount, "XLM for STABLEUSD")
		response, err := GetStableCoin(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword, LocalSeed, input[1])
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

		_ ,err = utils.StoFWithCheck(input[2])
		if err != nil {
			log.Println(err)
			break
		}

		response, err := AddCollateral(LocalContractor.U.LoginUserName, LocalContractor.U.LoginPassword, input[1], input[2])
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
		amount, err := utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		// convert this to int and check if int
		fmt.Println("Exchanging", amount, "XLM for STABLEUSD")
		response, err := GetStableCoin(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword, LocalSeed, input[1])
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

		_ ,err = utils.StoFWithCheck(input[2])
		if err != nil {
			log.Println(err)
			break
		}

		response, err := AddCollateral(LocalOriginator.U.LoginUserName, LocalOriginator.U.LoginPassword, input[1], input[2])
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
	}
	return nil
}
