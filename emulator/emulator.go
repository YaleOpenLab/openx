package main

import (
	"fmt"
	"log"

	database "github.com/OpenFinancing/openfinancing/database"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	scan "github.com/OpenFinancing/openfinancing/scan"
	"github.com/spf13/viper"
)

// package emulator is used to emulate the environment of the platform and make changes
// as one would expect t o in the frontend. This is not meaent to be run anywhere and
// should be used only for testing.

// have different entities that will be used across the files here
// emulator is intended to be a model for a frontend platform that would later be developed
// using the same backend that we have right now
var (
	LocalRecipient    database.Recipient
	LocalUser         database.User
	LocalInvestor     database.Investor
	LocalContractor   solar.Entity
	LocalOriginator   solar.Entity
	LocalSeed         string
	LocalSeedPwd      string
	PlatformPublicKey string
)

var ApiUrl = "http://localhost:8080"

// this platform public key would be public, so makes sense that we hardcode it

func SetupConfig() (string, error) {
	var err error
	viper.SetConfigType("yaml")
	viper.SetConfigName("emulator")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading email values from config file")
		return "", err
	}

	PlatformPublicKey = viper.Get("PlatformPublicKey").(string)

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
	role, err := Login(username, pwhash)
	if err != nil {
		return "", err
	}
	return role, nil
}

func main() {

	role, err := SetupConfig()
	if err != nil {
		log.Fatal(err)
	}
	switch role {
	// start loops for each role, would be nice if we could come up with an alternative to
	// duplication here
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
