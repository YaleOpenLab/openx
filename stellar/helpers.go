package main

import (
	"fmt"
	"log"

	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
)

func ValidateInputs() {
	if (opts.RecYears != 0) && !(opts.RecYears == 3 || opts.RecYears == 5 || opts.RecYears == 7) {
		// right now payoff periods are limited, I guess they don't need to be,
		// but in this case just are. Call this fucntion later when orders are being
		// created. Maybe don't need to restrict this at all?
		log.Fatal(fmt.Errorf("Number of years not supported"))
	}
}

func StartPlatform() error {
	ValidateInputs()
	allOrders, err := database.RetrieveAllOrders()
	if err != nil {
		log.Println("Error retrieving all orders from the database")
		return err
	}

	if len(allOrders) == 0 {
		log.Println("Populating database with test values")
		err = database.InsertDummyData()
		if err != nil {
			return err
		}
	}
	// setup issuer account if the platform doesn't  already exist
	// check whether the platform exists
	tmpPlatform, err := database.RetrievePlatform()
	if err != nil {
		return err
	}
	if len(tmpPlatform.PublicKey) == 0 {
		// this is the first time, so we initialize a platform
		log.Println("Initializing a new platform")
		_, err := database.InitializePlatform()
		if err != nil {
			return err
		}
	} else {
		log.Println("Platform already exists, using existing one")
	}
	return nil
}

func DecryptSeed() (string, error) {
	fmt.Printf("%s: ", "ENTER PASSWORD TO UNLOCK THE PLATFORM")
	rawPassword, err := utils.ScanRawPassword()
	if err != nil {
		log.Println(err)
		return "", err
	}
	platformSeed := database.GetSeedFromEncryptedSeed("seed.hex", rawPassword)
	return platformSeed, nil
}

func NewInvestorPrompt() error {
	fmt.Printf("%s: ", "ENTER YOUR REAL NAME")
	realName, err := utils.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return err
	}
	fmt.Printf("%s: ", "ENTER YOUR USERNAME")
	loginUserName, err := utils.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return err
	}

	fmt.Printf("%s: ", "ENTER DESIRED PASSWORD, YOU WILL NOT BE ASKED TO CONFIRM THIS")
	loginPassword, err := utils.ScanForPassword()
	if err != nil {
		fmt.Println("Couldn't read password")
		return err
	}

	investor, err := database.NewInvestor(loginUserName, loginPassword, realName)
	if err != nil {
		log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
		return err
	}
	err = database.InsertInvestor(investor)
	if err != nil {
		log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
		return err
	}
	return nil
}

func NewRecipientPrompt() error {
	fmt.Printf("%s: ", "ENTER YOUR REAL NAME")
	realName, err := utils.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return err
	}
	fmt.Printf("%s: ", "ENTER YOUR USERNAME")
	loginUserName, err := utils.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return err
	}

	fmt.Printf("%s: ", "ENTER DESIRED PASSWORD, YOU WILL NOT BE ASKED TO CONFIRM THIS")
	loginPassword, err := utils.ScanForPassword()
	if err != nil {
		fmt.Println("Couldn't read password")
		return err
	}

	recipient, err := database.NewRecipient(loginUserName, loginPassword, realName)
	if err != nil {
		log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
		return err
	}
	err = database.InsertRecipient(recipient)
	if err != nil {
		log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
		return err
	}
	return nil
}

func LoginPrompt() (database.Investor, database.Recipient, bool, error) {
	rbool := false
	var investor database.Investor
	var recipient database.Recipient
	fmt.Printf("%s: ", "---ARE YOU AN INVESTOR (I) OR RECIPIENT (R)? ---")
	optS, err := utils.ScanForString()
	if err != nil {
		log.Println("Failed to read user input")
		return investor, recipient, rbool, err
	}
	if optS == "I" || optS == "i" {
		fmt.Println("WELCOME BACK INVESTOR")
	} else if optS == "R" || optS == "r" {
		fmt.Println("WELCOME BACK RECIPIENT")
		rbool = true
	} else {
		log.Println("INVALID INPUT, EXITING!")
		return investor, recipient, rbool, fmt.Errorf("INVALID INPUT, EXITING!")
	}
	// ask for username and password combo here
	fmt.Printf("%s: ", "ENTER YOUR USERNAME")
	loginUserName, err := utils.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return investor, recipient, rbool, err
	}

	fmt.Printf("%s: ", "ENTER YOUR PASSWORD: ")
	loginPassword, err := utils.ScanForPassword()
	if err != nil {
		fmt.Println("Couldn't read password")
		return investor, recipient, rbool, err
	}
	user, err := database.ValidateUser(loginUserName, loginPassword)
	if err != nil {
		fmt.Println("Couldn't read password")
		return investor, recipient, rbool, err
	}
	recipient, _ = database.RetrieveRecipient(user.Index)
	investor, err = database.RetrieveInvestor(user.Index)
	if err != nil {
		// there is no investor, means user is a recipient
		rbool = true
	}
	return investor, recipient, rbool, nil
}
