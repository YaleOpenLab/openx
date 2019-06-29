package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"

	xlm "github.com/Varunram/essentials/crypto/xlm"
	scan "github.com/Varunram/essentials/scan"
	database "github.com/YaleOpenLab/openx/database"
	platform "github.com/YaleOpenLab/openx/platforms"
	solar "github.com/YaleOpenLab/openx/platforms/opensolar"
)

// StartPlatform starts the platform
func StartPlatform() (string, string, error) {
	var publicKey string
	var seed string
	database.CreateHomeDir()
	allContracts, err := solar.RetrieveAllProjects()
	if err != nil {
		log.Println("Error retrieving all projects from the database")
		return publicKey, seed, err
	}

	if len(allContracts) == 0 {
		log.Println("Populating database with test values")
		err = InsertDummyData(false)
		if err != nil {
			return publicKey, seed, err
		}
	}
	err = platform.InitializePlatform()
	return publicKey, seed, err
}

// NewUserPrompt is an ofcli helper function
func NewUserPrompt() (string, string, string, string, error) {
	realName, err := scan.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return "", "", "", "", err
	}
	fmt.Printf("%s: ", "ENTER YOUR USERNAME")
	loginUserName, err := scan.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return "", "", "", "", err
	}

	_, err = database.CheckUsernameCollision(loginUserName)
	if err != nil {
		fmt.Printf("%s", "username already taken, please choose a different one")
		return "", "", "", "", errors.New("username already taken, please choose a different one")
	}
	fmt.Printf("%s: ", "ENTER DESIRED PASSWORD, YOU WILL NOT BE ASKED TO CONFIRM THIS")
	loginPassword, err := scan.ScanForPassword()
	if err != nil {
		fmt.Println("Couldn't read password")
		return "", "", "", "", err
	}
	fmt.Printf("%s: ", "ENTER SEED PASSWORD, YOU WILL NOT BE ASKED TO CONFIRM THIS")
	seedPassword, err := scan.ScanForPassword()
	return realName, loginUserName, loginPassword, seedPassword, err
}

// NewInvestorPrompt is an ofcli helper function
func NewInvestorPrompt() error {
	log.Println("You have chosen to create a new investor account, welcome")
	loginUserName, loginPassword, realName, seedpwd, err := NewUserPrompt()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = database.NewInvestor(loginUserName, loginPassword, seedpwd, realName)
	if err != nil {
		log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
		return err
	}
	return err
}

// NewRecipientPrompt is an ofcli helper function
func NewRecipientPrompt() error {
	log.Println("You have chosen to create a new recipient account, welcome")
	loginUserName, loginPassword, realName, seedpwd, err := NewUserPrompt()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = database.NewRecipient(loginUserName, loginPassword, seedpwd, realName)
	if err != nil {
		log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
		return err
	}
	return err
}

// LoginPrompt is an ofcli helper function
func LoginPrompt() (database.Investor, database.Recipient, solar.Entity, bool, bool, error) {
	rbool := false
	cbool := false
	var investor database.Investor
	var recipient database.Recipient
	var contractor solar.Entity
	fmt.Println("---------SELECT YOUR ROLE---------")
	fmt.Println(" i. INVESTOR")
	fmt.Println(" r. RECIPIENT")
	fmt.Println(" c. CONTRACTOR")
	optS, err := scan.ScanForString()
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
		return investor, recipient, contractor, rbool, cbool, errors.New("INVALID INPUT, EXITING!")
	}
	// ask for username and password combo here
	fmt.Printf("%s: ", "ENTER YOUR USERNAME")
	loginUserName, err := scan.ScanForString()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return investor, recipient, contractor, rbool, cbool, err
	}

	fmt.Printf("%s: ", "ENTER YOUR PASSWORD: ")
	loginPassword, err := scan.ScanForPassword()
	if err != nil {
		fmt.Println("Couldn't read password")
		return investor, recipient, contractor, rbool, cbool, err
	}
	user, err := database.ValidateUser(loginUserName, loginPassword)
	if err != nil {
		fmt.Println("Couldn't read password")
		return investor, recipient, contractor, rbool, cbool, err
	}
	if rbool {
		recipient, err = database.RetrieveRecipient(user.Index)
		if err != nil {
			return investor, recipient, contractor, rbool, cbool, err
		}
	} else if cbool {
		contractor, err = solar.RetrieveEntity(user.Index)
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

// OriginContractPrompt is an ofcli helper function
func OriginContractPrompt(contractor *solar.Entity) error {
	fmt.Println("YOU HAVE DECIDED TO PROPOSE A NEW CONTRACT")
	fmt.Println("ENTER THE PANEL SIZE")
	panelSize, err := scan.ScanForString()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE COST OF PROJECT")
	totalValue, err := scan.ScanForFloat()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE LOCATION OF PROJECT")
	location, err := scan.ScanForString()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE ESTIMATED NUMBER OF YEARS FOR COMPLETION")
	years, err := scan.ScanForInt()
	if err != nil {
		return err
	}
	fmt.Println("ENTER METADATA REGARDING THE PROJECT")
	metadata, err := scan.ScanForString()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE RECIPIENT'S USER ID")
	recIndex, err := scan.ScanForInt()
	if err != nil {
		return err
	}
	originContract, err := contractor.Originate(panelSize, totalValue, location, years, metadata, recIndex, "blind")
	if err != nil {
		return err
	}
	// project insertion is done by the  above function, so we needn't call the database to do it again for us
	PrintProject(originContract)
	return nil
}

// ProposeContractPrompt is an ofcli helper function
func ProposeContractPrompt(contractor *solar.Entity) error {
	fmt.Println("YOU HAVE DECIDED TO PROPOSE A NEW CONTRACT")
	fmt.Println("ENTER THE PROJECT INDEX")
	contractIndex, err := scan.ScanForInt()
	if err != nil {
		return err
	}
	// we need to check if this contract index exists and retrieve
	rContract, err := solar.RetrieveProject(contractIndex)
	if err != nil {
		return err
	}
	log.Println("YOUR CONTRACT IS: ")
	PrintProject(rContract)
	if rContract.Index == 0 || rContract.Stage != 1 {
		// prevent people form porposing contracts for non originated contracts
		return errors.New("Invalid contract index")
	}
	panelSize := rContract.PanelSize
	location := rContract.State
	fmt.Println("ENTER THE COST OF PROJECT")
	totalValue, err := scan.ScanForFloat()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE ESTIMATED NUMBER OF YEARS FOR COMPLETION")
	years, err := scan.ScanForInt()
	if err != nil {
		return err
	}
	fmt.Println("ENTER METADATA REGARDING THE PROJECT")
	metadata, err := scan.ScanForString()
	if err != nil {
		return err
	}
	fmt.Println("ENTER THE RECIPIENT'S USER ID")
	recIndex, err := scan.ScanForInt()
	if err != nil {
		return err
	}
	originContract, err := contractor.Propose(panelSize, totalValue, location, years, metadata, recIndex, contractIndex, "blind")
	if err != nil {
		return err
	}
	// project insertion is done by the  above function, so we needn't call the database to do it again for us
	PrintProject(originContract)
	return nil
}

// Stage3ProjectsDisplayPrompt is an ofcli helper function
func Stage3ProjectsDisplayPrompt() {
	fmt.Println("------------LIST OF ALL AVAILABLE PROJECTS------------")
	allProjects, err := solar.RetrieveProjectsAtStage(solar.Stage3.Number)
	if err != nil {
		log.Println("Error retrieving all projects from the database")
	} else {
		PrintProjects(allProjects)
	}
}

// DisplayOriginProjects is an ofcli helper function
func DisplayOriginProjects() {
	fmt.Println("PRINTING ALL ORIGINATED PROJECTS: ")
	x, err := solar.RetrieveProjectsAtStage(solar.Stage0.Number)
	if err != nil {
		log.Println(err)
	} else {
		PrintProjects(x)
	}
}

// ExitPrompt is an ofcli helper function
func ExitPrompt() {
	// check whether he wants to go back to the display all screen again
	fmt.Println("DO YOU REALLY WANT TO EXIT? (PRESS Y TO CONFIRM)")
	exitOpt, err := scan.ScanForString()
	if err != nil {
		log.Println(err)
	}
	if exitOpt == "Y" || exitOpt == "y" {
		fmt.Println("YOU HAVE DECIDED TO EXIT")
		log.Fatal("")
	}
}

// BalanceDisplayPrompt is an ofcli helper function
func BalanceDisplayPrompt(publicKey string) {
	balances, err := xlm.GetAllBalances(publicKey)
	if err != nil {
		log.Println(err)
	} else {
		PrintBalances(balances)
	}
}
