package main

import (
	"fmt"
	"log"
	"os"

	consts "github.com/YaleOpenLab/smartPropertyMVP/stellar/consts"
	scan "github.com/YaleOpenLab/smartPropertyMVP/stellar/scan"
	wallet "github.com/YaleOpenLab/smartPropertyMVP/stellar/wallet"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
)

func CreateHomeDir() {
	if _, err := os.Stat(consts.TellerHomeDir); os.IsNotExist(err) {
		// directory does not exist, create one
		log.Println("Creating home directory for teller")
		os.MkdirAll(consts.TellerHomeDir, os.ModePerm)
	}
}

func CreateFile() {
	seedPath := consts.TellerHomeDir + "/seed.hex"
	if _, err := os.Stat(seedPath); os.IsNotExist(err) {
		// directory does not exist, create one
		ColorOutput("Creating home directory for teller", RedColor)
		os.MkdirAll(consts.TellerHomeDir, os.ModePerm)
		fmt.Printf("Generating keys in home directory, please enter a password: ")
		password, err := scan.ScanRawPassword()
		if err != nil {
			log.Fatal(err)
		}
		seed, address, err := xlm.GetKeyPair()
		if err != nil {
			fmt.Println(err)
		}
		log.Printf("SEED: %s\nADDRESS:%s", seed, address)
		err = wallet.StoreSeed(seed, password, seedPath)
		if err != nil {
			log.Fatal(err)
		}
		PublicKey = address
		Seed = seed
	} else {
		var err error
		ColorOutput("Please enter the password to decrypt your account: ", RedColor)
		password, err := scan.ScanRawPassword()
		if err != nil {
			log.Fatal(err)
		}
		PublicKey, Seed, err = wallet.RetrieveSeed(consts.TellerHomeDir+"/tellerseed.hex", password)
		if err != nil {
			log.Fatal(err)
		}
	}
}
