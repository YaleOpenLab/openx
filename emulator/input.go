package main

import (
	"fmt"
	"log"

	utils "github.com/OpenFinancing/openfinancing/utils"
)

func ParseInputInv(input []string) error {
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
	// have a list of all the supported commands here and then check for length, etc inside
	case "display":
		// display is a  broad command and needs to have a subcommand
		if len(input) == 1 {
			// only display was given, so display help command
			log.Println("HELP COMMANDS")
			break
		}
		subcommand := input[1]
		switch subcommand {
		case "balances":
			log.Println("Calling balances API")
			break
		case "profile":
			log.Println("Displaying Profile")
			break
		case "projects":
			if len(input) == 2 {
				// only display was given, so display help command
				log.Println("PROJECTS HELP COMMANDS")
				break
			}
			subsubcommand := input[2]
			switch subsubcommand {
			case "originated":
				log.Println("Displaying all originated (stage 1) projects")
				break
			case "seed":
				log.Println("Displaying all seed (stage 1.5) projects")
				break
			case "proposed":
				log.Println("Displaying all proposed (stage 2) projects")
				break
			case "open":
				log.Println("Displaying open (stage 3) projects")
				break
			case "funded":
				log.Println("Displaying funded (stage 4) projects")
				break
			case "installed":
				log.Println("Displaying installed (stage 5) projects")
				break
			case "power":
				log.Println("Displaying funded (stage 6) projects")
				break
			case "fin":
				log.Println("Displaying funded (stage 7) projects")
				break
			}
			break
		} // end of display
	case "exchange":
		if len(input) == 1 {
			// only display was given, so display help command
			log.Println("EXCHANGE HELP COMMANDS")
			break
		}
		amount, err := utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		// convert this to int and check if int
		fmt.Println("Exchanging", amount, "XLM for STABLEUSD")
		// end of exchange
	case "vote":
		if len(input) == 1 {
			log.Println("VOTE HELP COMMANDS")
			break
		}
		projIndex, err := utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println("Voting t owards proposed contract:", projIndex)
		// end of vote
	case "ipfs":
		if len(input) == 1 {
			log.Println("IPFS HELP COMMANDS")
			break
		}
		inputString := input[1]
		fmt.Println("HASHING", inputString, "using IPFS")
		// end of ipfs
	case "kyc":
		if LocalInvestor.U.Inspector {
			fmt.Println("WELCOME TO THE KYC MENU")
		}
		// end of kyc
	}
	return nil
}

func ParseInputRecp(input []string) error {
	// Various command supported for the recipient
	if len(input) == 0 {
		// shouldn't happen, still
		return fmt.Errorf("Length of input array is zero, quitting!")
	}
	// input is greater than length 1 which means we can parse according to the command given
	command := input[0]
	switch command {
	case "display":
		// display is a  broad command and needs to have a subcommand
		if len(input) == 1 {
			// only display was given, so display help command
			log.Println("HELP COMMANDS")
			break
		}
		subcommand := input[1]
		switch subcommand {
		case "balances":
			log.Println("Calling balances API")
			break
		case "profile":
			log.Println("Displaying Profile")
			break
		case "projects":
			if len(input) == 2 {
				// only display was given, so display help command
				log.Println("PROJECTS HELP COMMANDS")
				break
			}
			subsubcommand := input[2]
			switch subsubcommand {
			case "originated":
				log.Println("Displaying all originated (stage 1) projects")
				break
			case "seed":
				log.Println("Displaying all seed (stage 1.5) projects")
				break
			case "proposed":
				log.Println("Displaying all proposed (stage 2) projects")
				break
			case "open":
				log.Println("Displaying open (stage 3) projects")
				break
			case "locked":
				log.Println("Displaying all locked (stage 3) projects")
				break
			case "funded":
				log.Println("Displaying funded (stage 4) projects")
				break
			case "installed":
				log.Println("Displaying installed (stage 5) projects")
				break
			case "power":
				log.Println("Displaying funded (stage 6) projects")
				break
			case "fin":
				log.Println("Displaying funded (stage 7) projects")
				break
			}
			break
		} // end of display
	case "exchange":
		if len(input) == 1 {
			// only display was given, so display help command
			log.Println("EXCHANGE HELP COMMANDS")
			break
		}
		amount, err := utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		// convert this to int and check if int
		fmt.Println("Exchanging", amount, "XLM for STABLEUSD")
		// end of exchange
	case "ipfs":
		if len(input) == 1 {
			log.Println("IPFS HELP COMMANDS")
			break
		}
		inputString := input[1]
		fmt.Println("HASHING", inputString, "using IPFS")
		// end of ipfs
	case "finalize":
		if len(input) == 1 {
			log.Println("FINALIZE HELP COMMANDS")
			break
		}
		projIndex, err := utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println("FINALIZING PROJECT WITH INDEX:", projIndex)
		// end of payback
	case "payback":
		if len(input) == 1 {
			log.Println("PAYBACK HELP COMMANDS")
			break
		}
		projIndex, err := utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println("PAYING TOWARDS PROJECT WITH INDEX:", projIndex)
		// end of payback
	case "originate":
		fmt.Println("ORIGINATING PROJECT!")
		// end of payback
	}
	return nil
}

func ParseInputCont(input []string) error {
	// Various command supported for the recipient
	if len(input) == 0 {
		// shouldn't happen, still
		return fmt.Errorf("Length of input array is zero, quitting!")
	}
	// input is greater than length 1 which means we can parse according to the command given
	command := input[0]
	switch command {
	case "display":
		// display is a  broad command and needs to have a subcommand
		if len(input) == 1 {
			// only display was given, so display help command
			log.Println("HELP COMMANDS")
			break
		}
		subcommand := input[1]
		switch subcommand {
		case "balances":
			log.Println("Calling balances API")
			break
		case "profile":
			log.Println("Displaying Profile")
			break
		case "projects":
			if len(input) == 2 {
				// only display was given, so display help command
				log.Println("PROJECTS HELP COMMANDS")
				break
			}
			subsubcommand := input[2]
			switch subsubcommand {
			case "originated":
				log.Println("Displaying all originated (stage 1) projects")
				break
			case "seed":
				log.Println("Displaying all seed (stage 1.5) projects")
				break
			case "proposed":
				log.Println("Displaying all proposed (stage 2) projects")
				break
			case "open":
				log.Println("Displaying open (stage 3) projects")
				break
			case "locked":
				log.Println("Displaying all locked (stage 3) projects")
				break
			case "funded":
				log.Println("Displaying funded (stage 4) projects")
				break
			case "installed":
				log.Println("Displaying installed (stage 5) projects")
				break
			case "power":
				log.Println("Displaying funded (stage 6) projects")
				break
			case "fin":
				log.Println("Displaying funded (stage 7) projects")
				break
			}
			break
		} // end of display
	case "exchange":
		if len(input) == 1 {
			// only display was given, so display help command
			log.Println("EXCHANGE HELP COMMANDS")
			break
		}
		amount, err := utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		// convert this to int and check if int
		fmt.Println("Exchanging", amount, "XLM for STABLEUSD")
		// end of exchange
	case "ipfs":
		if len(input) == 1 {
			log.Println("IPFS HELP COMMANDS")
			break
		}
		inputString := input[1]
		fmt.Println("HASHING", inputString, "using IPFS")
		// end of ipfs
	case "propose":
		fmt.Println("PROPOSED PROJECT INTERFACE")
		break
		// end of propose
	case "myproposed":
		fmt.Println("VIEWING ALL MY PROPOSED PROJECTS")
		break
		// end of myproposed
	case "collateral":
		fmt.Println("CREATING A NEW TYPE OF COLLATERAL")
		break
		// end of collateral
	}
	return nil
}

func ParseInputCont(input []string) error {
	// Various command supported for the recipient
	if len(input) == 0 {
		// shouldn't happen, still
		return fmt.Errorf("Length of input array is zero, quitting!")
	}
	// input is greater than length 1 which means we can parse according to the command given
	command := input[0]
	switch command {
	case "display":
		// display is a  broad command and needs to have a subcommand
		if len(input) == 1 {
			// only display was given, so display help command
			log.Println("HELP COMMANDS")
			break
		}
		subcommand := input[1]
		switch subcommand {
		case "balances":
			log.Println("Calling balances API")
			break
		case "profile":
			log.Println("Displaying Profile")
			break
		case "projects":
			if len(input) == 2 {
				// only display was given, so display help command
				log.Println("PROJECTS HELP COMMANDS")
				break
			}
			subsubcommand := input[2]
			switch subsubcommand {
			case "originated":
				log.Println("Displaying all originated (stage 1) projects")
				break
			case "seed":
				log.Println("Displaying all seed (stage 1.5) projects")
				break
			case "proposed":
				log.Println("Displaying all proposed (stage 2) projects")
				break
			case "open":
				log.Println("Displaying open (stage 3) projects")
				break
			case "locked":
				log.Println("Displaying all locked (stage 3) projects")
				break
			case "funded":
				log.Println("Displaying funded (stage 4) projects")
				break
			case "installed":
				log.Println("Displaying installed (stage 5) projects")
				break
			case "power":
				log.Println("Displaying funded (stage 6) projects")
				break
			case "fin":
				log.Println("Displaying funded (stage 7) projects")
				break
			}
			break
		} // end of display
	case "exchange":
		if len(input) == 1 {
			// only display was given, so display help command
			log.Println("EXCHANGE HELP COMMANDS")
			break
		}
		amount, err := utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		// convert this to int and check if int
		fmt.Println("Exchanging", amount, "XLM for STABLEUSD")
		// end of exchange
	case "ipfs":
		if len(input) == 1 {
			log.Println("IPFS HELP COMMANDS")
			break
		}
		inputString := input[1]
		fmt.Println("HASHING", inputString, "using IPFS")
		// end of ipfs
	case "propose":
		fmt.Println("PROPOSING PRE-ORIGIN PROJECT INTERFACE")
		break
		// end of propose
	case "myproposed":
		fmt.Println("VIEWING ALL MY PRE-ORIGIN PROJECTS")
		break
		// end of myproposed
	case "myoriginated":
		fmt.Println("VIEWING ALL MY PRE-ORIGINATED PROJECTS")
		break
		// end of myoriginated
	}
	return nil
}
