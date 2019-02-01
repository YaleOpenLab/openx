package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"log"
	// "os"
	// "os/signal"
	"strings"

	consts "github.com/OpenFinancing/openfinancing/consts"
	database "github.com/OpenFinancing/openfinancing/database"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	scan "github.com/OpenFinancing/openfinancing/scan"
)

// package emulator is used to emulate the environment of the platform and make changes
// as one would expect t o in the frontend. This is not meaent to be run anywhere and
// should be used only for testing.

// have different entities that will be used across the files here
var (
	LocalRecipient  database.Recipient
	LocalUser       database.User
	LocalInvestor   database.Investor
	LocalContractor solar.Entity
	LocalOriginator solar.Entity
	LocalSeed       string
	LocalSeedPwd    string
)

var ApiUrl = "http://localhost:8080"
var PlatformPublicKey = "GDULAIM6N6SIW7MWS3NDJPY3UIFOHSM4766WQ6O6EKFDBC7PF53VKYLY"

func main() {

	fmt.Println("WELCOME TO THE SMARTSOLAR EMULATOR")
	ColorOutput("ENTER YOUR USERNAME: ", CyanColor)
	username, err := scan.ScanForString()
	if err != nil {
		log.Fatal(err)
	}
	ColorOutput("ENTER YOUR PASSWORD: ", CyanColor)
	pwhash, err := scan.ScanForPassword()
	if err != nil {
		log.Fatal(err)
	}
	// need to validate with the RPC here
	wString, err := Login(username, pwhash)
	if err != nil {
		log.Fatal(err)
	}

	switch wString {
	case "Investor":
		log.Fatal(LoopInv())
	case "Recipient":
		log.Fatal(LoopRecp())
	case "Originator":
		log.Fatal(LoopOrig())
	case "Contractor":
		log.Fatal(LoopCont())
	default:
		log.Println("It should never come here")
	}
}

func LoopInv() error {
	// this loop is for an investor
	// we have authenticated the user and stored the details in an appropriate structure
	// need to repeat this struct everywhere because having separate functions and importing
	// it doesn't seem to work
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

func LoopOrig() error {
	// This loop is exclusive to an originator
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

		err = ParseInputOrig(cmdslice)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

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
