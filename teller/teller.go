package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"log"
	"os"
	"os/signal"
	"strings"

	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
)

// package teller contains the remote client code that would be run on the client's
// side and communicate information with us and with atonomi and other partners.
// that it belongs, the contract, recipient, and eg. person who installed it.
// Consider doing this with IoT partners, eg. Atonomi.

// Teller authenticates with the platform using a remote API and then retrieves
// credentials once authenticated. Both the teller and the project recipient on the
// platform are the same entity, just that the teller is associated with the hw device.
// hw device needs an id and stuff, hopefully Atonomi can give us that.
// Teller tracks whenever the device starts and goes off, so we know when exactly the device was
// switched off. This is enough as proof that the device was running in between. This also
// avoids needing to poll the blockchain often and saves on the (minimal, still) tx fee.

// Since we can't compile this directly on the raspberry pi, we need to cross compile he
// go executable and transfer it over to the raspberry pi
// the following should do the trick for us
// env GOOS=linux GOARCH=arm GOARM=5 go build
var (
	LocalRecipient    database.Recipient
	RecpSeed          string
	RecpPublicKey     string
	PlatformPublicKey string
	ApiUrl            string
	DeviceId          string
	DeviceLocation    string
	DeviceInfo        string
	StartHash         string
	NowHash           string
	LocalSeedPwd      string
	PlatformEmail     string
)

var cleanupDone chan struct{}

func main() {
	// Authenticate with the platform
	err := SetupConfig()
	if err != nil {
		log.Fatal(err)
	}
	ColorOutput("TELLER PUBKEY: "+RecpPublicKey, GreenColor)
	ColorOutput("DEVICE ID: "+DeviceId, GreenColor)
	// log.Fatal("Checks done") // REMOVE THIS BEFORE COMMIT
	go CheckPayback()
	StartHash, err = BlockStamp()
	if err != nil {
		log.Fatal(err)
	}
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

	DeviceInfo = "Raspberry Pi3 Model B+"
	go StartServer()

	// channels for preventing immediate sigint
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		fmt.Println("\nSigint received in quit function. not quitting!")
		close(cleanupDone)
	}()
	for {
		// setup reader with max 4K input chars
		msg, err := rl.Readline()
		if err != nil {
			var err error
			err = EndHandler() // error, user wants to quit
			for err != nil {
				log.Println(err)
				err = EndHandler()
				<-cleanupDone // to prevent user from quitting when sigint arrives
			}
			break
		}
		msg = strings.TrimSpace(msg)
		if len(msg) == 0 {
			continue
		}
		rl.SaveHistory(msg)

		cmdslice := strings.Fields(msg)
		ColorOutput("entered command: "+msg, YellowColor)

		ParseInput(cmdslice)
	}
}
