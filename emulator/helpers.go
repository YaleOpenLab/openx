package main

import (
	"fmt"
	"log"

	utils "github.com/Varunram/essentials/utils"
	"github.com/chzyer/readline"
)

func autoComplete() readline.AutoCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("help",
			readline.PcItem("ping"),
			readline.PcItem("exchange"),
			readline.PcItem("ipfs"),
			readline.PcItem("send"),
			readline.PcItem("receive"),
			readline.PcItem("create"),
			readline.PcItem("kyc"),
			readline.PcItem("increasetrust"),
			readline.PcItem("vote"),
			readline.PcItem("invest"),
			readline.PcItem("unlock"),
			readline.PcItem("payback"),
			readline.PcItem("finalize"),
			readline.PcItem("originate"),
			readline.PcItem("calculate"),
			readline.PcItem("propose"),
			readline.PcItem("myproposed"),
			readline.PcItem("addcollateral"),
			readline.PcItem("mystage0"),
			readline.PcItem("mystage1"),
			readline.PcItem("preoriginate"),
		),
		readline.PcItem("display",
			readline.PcItem("balance",
				readline.PcItem("xlm"),
				readline.PcItem("asset"),
			),
			readline.PcItem("info"),
		),
		readline.PcItem("exchange",
			readline.PcItem("amount"),
		),
		readline.PcItem("ipfs",
			readline.PcItem("string"),
		),
		readline.PcItem("send",
			readline.PcItem("amount"),
		),
		readline.PcItem("receive",
			readline.PcItem("xlm"),
			readline.PcItem("asset"),
		),
		readline.PcItem("create",
			readline.PcItem("asset"),
		),
		readline.PcItem("kyc",
			readline.PcItem("user"),
		),
		readline.PcItem("increasetrust",
			readline.PcItem("trustlimit"),
		),
		readline.PcItem("vote",
			readline.PcItem("project"),
		),
		readline.PcItem("invest",
			readline.PcItem("projIndex"),
			readline.PcItem("amount"),
			readline.PcItem("platform"),
		),
		readline.PcItem("unlock",
			readline.PcItem("projIndex"),
		),
		readline.PcItem("payback",
			readline.PcItem("projIndex"),
			readline.PcItem("amount"),
		),
		readline.PcItem("finalize",
			readline.PcItem("projIndex"),
		),
		readline.PcItem("originate",
			readline.PcItem("projIndex"),
		),
		readline.PcItem("calculate",
			readline.PcItem("projIndex"),
		),
		readline.PcItem("addcollateral",
			readline.PcItem("data"),
		),
	)
}

func displayHelper(input []string, username string, pwhash string, role string) {
	// display is a  broad command and needs to have a subcommand
	if len(input) == 1 {
		// only display was given, so display help command
		log.Println("<display><balance, profile, projects>")
		return
	}
	subcommand := input[1]
	switch subcommand {
	case "balance":
		if len(input) == 2 {
			log.Println("Calling balances API")
			balances, err := GetBalances(username, pwhash)
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
			balance, err := GetXLMBalance(username, pwhash)
			if err != nil {
				log.Println(err)
				break
			}
			ColorOutput("BALANCE: "+balance, MagentaColor)
		case "all":
			balances, err := GetBalances(username, pwhash)
			if err != nil {
				log.Println(err)
				break
			}
			PrintBalances(balances)
		default:
			balance, err := GetAssetBalance(username, pwhash, subcommand)
			if err != nil {
				log.Println(err)
				break
			}
			ColorOutput("BALANCE: "+balance, MagentaColor)
		}
	case "profile":
		log.Println("Displaying Profile")
		switch role {
		case "investor":
			PrintInvestor(LocalInvestor)
		case "recipient":
			PrintRecipient(LocalRecipient)
		case "contractor":
			PrintEntity(LocalContractor)
		case "originator":
			PrintEntity(LocalOriginator)
		}
	case "projects":
		if len(input) != 4 {
			// only display was given, so display help command
			log.Println("display projects <platform> <stageNumber>")
			break
		}
		platform := input[2]
		switch platform {
		case "opzones":
			log.Println("OPZONES PLATFORM")
			subsubcommand := input[3]
			switch subsubcommand {
			case "cbonds":
				log.Println("PRINTGING ALL OPEN Construction Bonds")
			case "lucoops":
				log.Println("PRINTGING ALL OPEN Living unit coops")
			}
		case "opensolar":
			subsubcommand := input[3]
			index, err := utils.StoICheck(subsubcommand)
			if err != nil {
				log.Println("Input not int, not retrieving!")
				return
			}
			arr, err := RetrieveProject(index)
			if err != nil {
				log.Println(err)
				break
			}
			PrintProjects(arr)
		}
	} // end of display
}

func exchangeHelper(input []string, username string, pwhash string, seed string) {
	if len(input) != 2 {
		// only display was given, so display help command
		log.Println("<exchange> amount")
		return
	}
	amount, err := utils.StoICheck(input[1])
	if err != nil {
		log.Println(err)
		return
	}
	// convert this to int and check if int
	fmt.Println("Exchanging", amount, "XLM for STABLEUSD")
	response, err := GetStableCoin(username, pwhash, input[1])
	if err != nil {
		log.Println(err)
		return
	}
	if response.Code == 200 {
		ColorOutput("SUCCESSFUL, CHECK BALANCES", GreenColor)
	} else {
		ColorOutput("RESPONSE STATUS: "+utils.ItoS(response.Code), GreenColor)
	}
}

func ipfsHelper(input []string, username string, pwhash string) {
	if len(input) != 2 {
		log.Println("<ipfs> string")
		return
	}
	inputString := input[1]
	hashString, err := GetIpfsHash(username, pwhash, inputString)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("IPFS HASH", hashString)
	// end of ipfs
}

func pingHelper() {
	err := PingRpc()
	if err != nil {
		log.Println(err)
		return
	}
}

func sendHelper(input []string, username string, pwhash string) {
	var err error
	if len(input) == 1 {
		log.Println("send <asset>")
		return
	}
	subcommand := input[1]
	switch subcommand {
	case "asset":
		if len(input) != 5 {
			log.Println("send asset <assetName> <destination> <amount>")
			return
		}

		assetName := input[2]
		destination := input[3]
		amount := input[4]

		txhash, err := SendLocalAsset(username, pwhash,
			LocalSeedPwd, assetName, destination, amount)
		if err != nil {
			log.Println(err)
		}
		ColorOutput("TX HASH: "+txhash, MagentaColor)
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
		txhash, err := SendXLM(username, pwhash, LocalSeedPwd, destination, amount, memo)
		if err != nil {
			log.Println(err)
		}
		ColorOutput("TX HASH: "+txhash, MagentaColor)
	}
}

func receiveHelper(input []string, username string, pwhash string) {
	// we can either receive from the faucet or trust issuers to receive assets
	var err error
	if len(input) == 1 {
		log.Println("receive <xlm, asset>")
		return
	}
	subcommand := input[1]
	switch subcommand {
	case "xlm":
		status, err := AskXLM(username, pwhash)
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

		status, err := TrustAsset(username, pwhash, assetName, issuerPubkey, limit, LocalSeedPwd)
		if err != nil {
			log.Println(err)
			break
		}
		if status.Code == 200 {
			ColorOutput("COIN REQUEST SUCCESSFUL, CHECK EMAIL", GreenColor)
		} else {
			ColorOutput("COIN REQUEST NOT SUCCESSFUL", RedColor)
		}
	} // end of receive
}

func createHelper(input []string, username string, pwhash string, pubkey string) {
	// create enables you to create tokens on stellar that you can excahnge with third parties.
	if len(input) == 1 {
		log.Println("create <asset>")
		return
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
		status, err := CreateAssetInv(username, pwhash, assetName, pubkey)
		if err != nil {
			log.Println(err)
			return
		}
		if status.Code == 200 {
			ColorOutput("INVESTMENT SUCCESSFUL, CHECK EMAIL", GreenColor)
		} else {
			ColorOutput("INVESTMENT NOT SUCCESSFUL", RedColor)
		}
	} // end of create
}

func kycHelper(input []string, username string, pwhash string, inspector bool) {
	var err error
	if !inspector {
		ColorOutput("YOU ARE NOT A KYC INSPECTOR", RedColor)
		return
	}
	if len(input) == 1 {
		log.Println("kyc <auth, view>")
		return
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
		status, err := AuthKyc(input[1], username, pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		if status.Code == 200 {
			ColorOutput("USER KYC'D!", GreenColor)
		} else {
			ColorOutput("USER NOT KYC'D", RedColor)
		}
	case "notdone":
		users, err := NotKycView(username, pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		PrintUsers(users)
		// print all the users who have kyc'd
	case "done":
		users, err := KycView(username, pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		PrintUsers(users)
		// print all the users who have kyc'd
	}
	// end of kyc
}

func increaseTrustHelper(input []string, username string, pwhash string) {
	if len(input) == 1 {
		log.Println("<increasetrust> trustlimit")
		return
	}
	trustLimit, err := utils.StoFWithCheck(input[1])
	if err != nil {
		log.Println(err, "input is not a string!")
	}
	if trustLimit == 0 {
		log.Println("Can't increaset trustlimti by zero, quitting!")
		return
	}
	response, err := IncreaseTrustLimit(username, pwhash, LocalSeedPwd, utils.FtoS(trustLimit))
	if err != nil {
		log.Println(err)
	}
	if response.Code == 200 {
		ColorOutput("SUCCESSFULLY INCREASED STABELCOIN TRUST LIMIT", GreenColor)
	} else {
		ColorOutput("COULD NOT INCREASE STABELCOIN TRUST LIMIT", RedColor)
	}
}

func sendSharesEmailHelper(input []string, username string, pwhash string) {
	if len(input) != 4 {
		log.Println("<sendshares> email1 email2 email3")
		return
	}
	email1 := input[1]
	email2 := input[2]
	email3 := input[3]

	response, err := SendSharesEmail(username, pwhash, email1, email2, email3)
	if err != nil {
		log.Println(err)
	}
	if response.Code == 200 {
		ColorOutput("SUCCESSFULLY SENT SHARES", GreenColor)
	} else {
		ColorOutput("COULD NOT SEND SHARES OUT TO PARTIES", RedColor)
	}
}

func genNewSharesHelper(input []string, username string, pwhash string, seedpwd string) {
	if len(input) != 4 {
		log.Println("<newshares> email1 email2 email3")
		return
	}

	email1 := input[1]
	email2 := input[2]
	email3 := input[3]

	response, err := SendNewSharesEmail(username, pwhash, seedpwd, email1, email2, email3)
	if err != nil {
		log.Println(err)
	}

	if response.Code == 200 {
		ColorOutput("SUCCESSFULLY SENT NEW SHARES", GreenColor)
	} else {
		ColorOutput("COULD NOT SEND NEW SHARES OUT TO PARTIES", RedColor)
	}
}
