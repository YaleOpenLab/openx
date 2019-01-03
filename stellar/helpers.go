package main

import (
	"fmt"
	"log"

	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	platform "github.com/YaleOpenLab/smartPropertyMVP/stellar/platform"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
)

func ValidateInputs() {
	if (opts.RecYears != 0) && !(opts.RecYears == 3 || opts.RecYears == 5 || opts.RecYears == 7) {
		// right now payoff periods are limited, I guess they don't need to be,
		// but in this case just are. Call this function later when projects are being
		// created. Maybe don't need to restrict this at all?
		log.Fatal(fmt.Errorf("Number of years not supported"))
	}
}

func StartPlatform() (string, string, error) {
	var publicKey string
	var seed string
	ValidateInputs()
	allContracts, err := database.RetrieveAllProjects()
	if err != nil {
		log.Println("Error retrieving all projects from the database")
		return publicKey, seed, err
	}

	if len(allContracts) == 0 {
		log.Println("Populating database with test values")
		err = InsertDummyData()
		if err != nil {
			return publicKey, seed, err
		}
	}
	publicKey, seed, err = platform.InitializePlatform()
	return publicKey, seed, err
}

func NewUserPrompt() (string, string, string, error) {
	realName, err := utils.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return "", "", "", err
	}
	fmt.Printf("%s: ", "ENTER YOUR USERNAME")
	loginUserName, err := utils.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return "", "", "", err
	}

	fmt.Printf("%s: ", "ENTER DESIRED PASSWORD, YOU WILL NOT BE ASKED TO CONFIRM THIS")
	loginPassword, err := utils.ScanForPassword()
	if err != nil {
		fmt.Println("Couldn't read password")
		return "", "", "", err
	}
	return realName, loginUserName, loginPassword, err
}

func NewInvestorPrompt() error {
	fmt.Printf("%s: ", "ENTER YOUR REAL NAME")

	loginUserName, loginPassword, realName, err := NewUserPrompt()
	if err != nil {
		log.Println(err)
		return err
	}

	investor, err := database.NewInvestor(loginUserName, loginPassword, realName)
	if err != nil {
		log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
		return err
	}
	err = investor.Save()
	if err != nil {
		log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
		return err
	}
	return nil
}

func NewRecipientPrompt() error {

	loginUserName, loginPassword, realName, err := NewUserPrompt()
	if err != nil {
		log.Println(err)
		return err
	}
	recipient, err := database.NewRecipient(loginUserName, loginPassword, realName)
	if err != nil {
		log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
		return err
	}
	err = recipient.Save()
	if err != nil {
		log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
		return err
	}
	return nil
}

func LoginPrompt() (database.Investor, database.Recipient, database.Entity, bool, bool, error) {
	rbool := false
	cbool := false
	var investor database.Investor
	var recipient database.Recipient
	var contractor database.Entity
	fmt.Println("---------SELECT YOUR ROLE---------")
	fmt.Println(" i. INVESTOR")
	fmt.Println(" r. RECIPIENT")
	fmt.Println(" c. CONTRACTOR")
	optS, err := utils.ScanForString()
	if err != nil {
		log.Println("Failed to read user input")
		return investor, recipient, contractor, rbool, cbool, err
	}
	if optS == "I" || optS == "i" {
		fmt.Println("WELCOME BACK INVESTOR")
	} else if optS == "R" || optS == "r" {
		fmt.Println("WELCOME BACK RECIPIENT")
		rbool = true
	} else if optS == "C" || optS == "c" {
		cbool = true
		fmt.Println("WELCOME BACK CONTRACTOR")
	} else {
		log.Println("INVALID INPUT, EXITING!")
		return investor, recipient, contractor, rbool, cbool, fmt.Errorf("INVALID INPUT, EXITING!")
	}
	// ask for username and password combo here
	fmt.Printf("%s: ", "ENTER YOUR USERNAME")
	loginUserName, err := utils.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return investor, recipient, contractor, rbool, cbool, err
	}

	fmt.Printf("%s: ", "ENTER YOUR PASSWORD: ")
	loginPassword, err := utils.ScanForPassword()
	if err != nil {
		fmt.Println("Couldn't read password")
		return investor, recipient, contractor, rbool, cbool, err
	}
	user, err := database.ValidateUser(loginUserName, loginPassword)
	if err != nil {
		fmt.Println("Couldn't read password")
		return investor, recipient, contractor, rbool, cbool, err
	}
	log.Println("WATCH USER IDNEX: ", user.Index)
	if rbool {
		recipient, err = database.RetrieveRecipient(user.Index)
		if err != nil {
			return investor, recipient, contractor, rbool, cbool, err
		}
	} else if cbool {
		contractor, err = database.RetrieveEntity(user.Index)
		if err != nil {
			return investor, recipient, contractor, rbool, cbool, err
		}
	} else {
		investor, err = database.RetrieveInvestor(user.Index)
		if err != nil {
			return investor, recipient, contractor, rbool, cbool, err
		}
	}
	return investor, recipient, contractor, rbool, cbool, nil
}

func OriginContractPrompt(contractor *database.Entity) error {
	fmt.Println("YOU HAVE DECIDED TO PROPOSE A NEW CONTRACT")
	fmt.Println("ENTER THE PANEL SIZE")
	panelSize, err := utils.ScanForString()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE COST OF PROJECT")
	totalValue, err := utils.ScanForInt()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE LOCATION OF PROJECT")
	location, err := utils.ScanForString()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE ESTIMATED NUMBER OF YEARS FOR COMPLETION")
	years, err := utils.ScanForInt()
	if err != nil {
		return err
	}
	fmt.Println("ENTER METADATA REGARDING THE PROJECT")
	metadata, err := utils.ScanForString()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE RECIPIENT'S USER ID")
	recIndex, err := utils.ScanForInt()
	if err != nil {
		return err
	}
	originContract, err := contractor.OriginContract(panelSize, totalValue, location, years, metadata, recIndex)
	if err != nil {
		return err
	}
	// project insertion is done by the  above function, so we needn't call the database to do it again for us
	PrintProject(originContract)
	return nil
}

func PrintAllProposedContracts(contractor *database.Entity) error {
	fmt.Println("LIST OF ALL PROPOSED CONTRACTS: ")
	// the database would be updated each time the user has an originated
	// contract, so we need to retrieve the contractor struct again
	contractorDup, err := database.RetrieveEntity(contractor.U.Index)
	if err != nil {
		return err
	}
	for _, elem := range contractor.ProposedContracts {
		PrintParams(elem.Params)
	}
	contractor = &contractorDup
	return nil
}

func ProposeContractPrompt(contractor *database.Entity) error {
	fmt.Println("YOU HAVE DECIDED TO PROPOSE A NEW CONTRACT")
	fmt.Println("ENTER THE PROJECT INDEX")
	contractIndex, err := utils.ScanForInt()
	if err != nil {
		return err
	}
	// we need to check if this contract index exists and retrieve
	rContract, err := database.RetrieveProject(contractIndex)
	if err != nil {
		return err
	}
	// TODO: PRint CONTRACTS HERE
	log.Println("YOUR CONTRACT IS: ")
	PrintProject(rContract)
	if rContract.Params.Index == 0 || rContract.Stage != 1 {
		// prevent people form porposing contracts for non originated contracts
		return fmt.Errorf("Invalid contract index")
	}
	panelSize := rContract.Params.PanelSize
	location := rContract.Params.Location
	fmt.Println("ENTER THE COST OF PROJECT")
	totalValue, err := utils.ScanForInt()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE ESTIMATED NUMBER OF YEARS FOR COMPLETION")
	years, err := utils.ScanForInt()
	if err != nil {
		return err
	}
	fmt.Println("ENTER METADATA REGARDING THE PROJECT")
	metadata, err := utils.ScanForString()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE RECIPIENT'S USER ID")
	recIndex, err := utils.ScanForInt()
	if err != nil {
		return err
	}
	originContract, err := contractor.ProposeContract(panelSize, totalValue, location, years, metadata, recIndex, contractIndex)
	if err != nil {
		return err
	}
	// project insertion is done by the  above function, so we needn't call the database to do it again for us
	PrintParams(originContract.Params)
	return nil
}
