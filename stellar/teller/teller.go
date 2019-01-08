package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"log"
	"strings"

	consts "github.com/YaleOpenLab/smartPropertyMVP/stellar/consts"
)

// package teller contains the remote client code that would be run on the client's
// side and communicate information with us and with atonomi and other partners.
// for now, we need a client that can start up and generate a pk/seed pair. This would
// be stored in the project struct and if anyone wants to see the status of this
// node, they can check the blockchain for the node's updates (node should update
// the blockchain at frequient intervals with power generation data in the memo
// field of the tx.
// Polling interval would be 5 minutes, 1440/5 = 288 updates a day

var PublicKey string
var Seed string

func main() {
	CreateHomeDir()
	CreateFile()
	promptColor := color.New(color.FgHiYellow).SprintFunc()
	whiteColor := color.New(color.FgHiWhite).SprintFunc()
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      promptColor("teller") + whiteColor("# "),
		HistoryFile: consts.TellerHomeDir + "/history.txt",
		// AutoComplete: lc.NewAutoCompleter(),
	})

	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()
	// main shell loop
	for {
		// setup reader with max 4K input chars
		msg, err := rl.Readline()
		if err != nil {
			break
		}
		msg = strings.TrimSpace(msg)
		if len(msg) == 0 {
			continue
		}
		rl.SaveHistory(msg)

		cmdslice := strings.Fields(msg)
		ColorOutput("entered command: "+msg, YellowColor)

		err = ParseInput(cmdslice)
		if err != nil {
			fmt.Println(err)
		}
	}
}
