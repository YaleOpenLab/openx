package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"strings"

	utils "github.com/Varunram/essentials/utils"
	"github.com/chzyer/readline"
)

// inputinv.go contains all the relevant emulator commands for the investor
// we have one parse input function for each entity on the platform
// ie investor, recipient, contractor, originator and potentially more depednign upon usage
// the input array contains the commands that we want to parse.
// first check the length of the input array and then define accordingly

// LoopInv loops over investor input
func LoopInv(rl *readline.Instance) error {
	// this loop is for an investor
	// we have authenticated the user and stored the details in an appropriate structure
	// need to repeat this struct everywhere because having separate functions and importing
	// it doesn't seem to work
	// the problem with having a conditional statement inside the loop is that it checks
	// role each time and that's not nice performance wise
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

		err = ParseInputInv(cmdslice)
		if err != nil {
			return errors.Wrap(err, "could not parse user input")
		}
	}
}

// ParseInputInv parses investor input
func ParseInputInv(input []string) error {
	var err error
	if len(input) == 0 {
		// shouldn't happen, still
		return errors.New("Length of input array is zero, quitting!")
	}
	// input is greater than length 1 which means we can parse according to the command given
	command := input[0]
	switch command {
	case "help":
		fmt.Println("LIST OF SUPPORTED COMMANDS: ")
		fmt.Println("ping, display, exchange, ipfs, vote, kyc, invest, create, send, receive")
	case "ping":
		pingHelper()
	case "display":
		displayHelper(input, LocalInvestor.U.Username, LocalInvestor.U.Pwhash, "investor")
	case "exchange":
		exchangeHelper(input, LocalInvestor.U.Username, LocalInvestor.U.Pwhash, LocalSeed)
	case "ipfs":
		ipfsHelper(input, LocalInvestor.U.Username, LocalInvestor.U.Pwhash)
	case "send":
		sendHelper(input, LocalInvestor.U.Username, LocalInvestor.U.Pwhash)
	case "receive":
		receiveHelper(input, LocalInvestor.U.Username, LocalInvestor.U.Pwhash)
	case "create":
		createHelper(input, LocalInvestor.U.Username, LocalInvestor.U.Pwhash, LocalInvestor.U.StellarWallet.PublicKey)
	case "kyc":
		kycHelper(input, LocalInvestor.U.Username, LocalInvestor.U.Pwhash, LocalInvestor.U.Inspector)
	case "increasetrust":
		increaseTrustHelper(input, LocalInvestor.U.Username, LocalInvestor.U.Pwhash)
	case "sendshares":
		sendSharesEmailHelper(input, LocalInvestor.U.Username, LocalInvestor.U.Pwhash)
	case "newshares":
		genNewSharesHelper(input, LocalInvestor.U.Username, LocalInvestor.U.Pwhash, LocalSeedPwd)
	// Investor only functions
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
		if status.Code == 200 {
			ColorOutput("VOTE CAST!", GreenColor)
		} else {
			ColorOutput("VOTE NOT CAST", RedColor)
		}
	case "invest":
		if len(input) < 4 {
			log.Println("Invest Params: invest <proj_number> <amount> <platform>")
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
		platform := input[3]
		switch platform {
		case "opensolar":
			// now we need to invest in this project, call RPC
			status, err := InvestInProject(input[1], input[2], LocalInvestor.U.Username, LocalInvestor.U.Pwhash, LocalSeedPwd)
			if err != nil {
				log.Println(err)
				break
			}
			if status.Code == 200 {
				ColorOutput("INVESTMENT SUCCESSFUL, CHECK EMAIL", GreenColor)
			} else {
				ColorOutput("INVESTMENT NOT SUCCESSFUL", RedColor)
			}
		case "opzones":
			if len(input) < 5 {
				log.Println("Invest Params: invest <proj_number> <amount> opzones <cbond / lucoop>")
				break
			}
			investmentChoice := input[4]
			switch investmentChoice {
			case "cbond":
				// now we need to invest in this project, call RPC
				status, err := InvestInOpzoneCBond(input[1], input[2], LocalInvestor.U.Username, LocalInvestor.U.Pwhash, LocalSeedPwd)
				if err != nil {
					log.Println(err)
					break
				}
				if status.Code == 200 {
					ColorOutput("INVESTMENT SUCCESSFUL, CHECK EMAIL", GreenColor)
				} else {
					ColorOutput("INVESTMENT NOT SUCCESSFUL", RedColor)
				}
			case "lucoop":
				status, err := InvestInLivingUnitCoop(input[1], input[2], LocalInvestor.U.Username, LocalInvestor.U.Pwhash, LocalSeedPwd)
				if err != nil {
					log.Println(err)
					break
				}
				if status.Code == 200 {
					ColorOutput("INVESTMENT SUCCESSFUL, CHECK EMAIL", GreenColor)
				} else {
					ColorOutput("INVESTMENT NOT SUCCESSFUL", RedColor)
				}
			}
		}
	}
	return nil
}
