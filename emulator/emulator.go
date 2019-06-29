package main

import (
	"github.com/pkg/errors"
	"log"

	scan "github.com/Varunram/essentials/scan"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	solar "github.com/YaleOpenLab/openx/platforms/opensolar"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// package emulator is used to emulate the environment of the platform and make changes
// as one would expect t o in the frontend. This is not meaent to be run anywhere and
// should be used only for testing.

// have different entities that will be used across the files here
// emulator is intended to be a model for a frontend platform that would later be developed
// using the same backend that we have right now
var (
	// have a global variable for each entity
	LocalRecipient    database.Recipient
	LocalInvestor     database.Investor
	LocalContractor   solar.Entity
	LocalOriginator   solar.Entity
	LocalSeed         string
	LocalSeedPwd      string
	PlatformPublicKey string
)

// ApiUrl points to the platform instance's public endpoint
var ApiUrl = "http://localhost:8080"

// SetupConfig reads from the teller's config file and authenticates with the platform
func SetupConfig() (string, error) {
	var err error
	viper.SetConfigType("yaml")
	viper.SetConfigName("emulator")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		return "", errors.Wrap(err, "error while reading email values from config file")
	}

	PlatformPublicKey = viper.Get("PlatformPublicKey").(string)

	log.Println("WELCOME TO THE SMARTSOLAR EMULATOR")

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
	role, err := Login(username, pwhash)
	if err != nil {
		return "", errors.Wrap(err, "could not login to the platform")
	}
	return role, nil
}

func main() {

	role, err := SetupConfig()
	if err != nil {
		log.Fatal(err)
	}

	promptColor := color.New(color.FgHiYellow).SprintFunc()
	whiteColor := color.New(color.FgHiWhite).SprintFunc()
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       promptColor("emulator") + whiteColor("# "),
		HistoryFile:  consts.TellerHomeDir + "/history_emulator.txt",
		AutoComplete: autoComplete(),
	})
	ColorOutput("YOUR SEED IS: "+LocalSeed, RedColor)
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	switch role {
	// start loops for each role, would be nice if we could come up with an alternative to
	// duplication here
	case "Investor":
		log.Fatal(LoopInv(rl))
	case "Recipient":
		log.Fatal(LoopRecp(rl))
	case "Originator":
		log.Fatal(LoopOrig(rl))
	case "Contractor":
		log.Fatal(LoopCont(rl))
	default:
		log.Println("It should never come here")
	}
}
