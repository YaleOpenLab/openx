package main

import (
	"fmt"

	database "github.com/OpenFinancing/openfinancing/database"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	"github.com/stellar/go/protocols/horizon"
)

func PrintProjects(projects []solar.Project) {
	for _, project := range projects {
		PrintProject(project)
	}
}

// this function pretty prints out some stuff that we need in main.go
func PrintProject(project solar.Project) {
	fmt.Println("          PROJECT INDEX: ", project.Params.Index)
	fmt.Println("          Panel Size: ", project.Params.PanelSize)
	fmt.Println("          Total Value: ", project.Params.TotalValue)
	fmt.Println("          Location: ", project.Params.Location)
	fmt.Println("          Money Raised: ", project.Params.MoneyRaised)
	fmt.Println("          Metadata: ", project.Params.Metadata)
	fmt.Println("          Years: ", project.Params.Years)
	fmt.Println("          Auction Type: ", project.AuctionType)
	fmt.Println("          PROJECT ORIGINATOR: ")
	PrintEntity(project.Originator)
	fmt.Println("          PROJECT STAGE: ", project.Stage)
	fmt.Println("          RECIPIENT: ")
	PrintRecipient(project.ProjectRecipient)
	if project.Stage >= 2 {
		fmt.Println("          PROJECT CONTRACTOR: ")
		PrintEntity(project.Contractor)
		fmt.Println("          Votes: ", project.Params.Votes)
	}
	if project.Stage >= 3 {
		fmt.Println("          Investor Asset Code: ", project.Params.InvestorAssetCode)
		fmt.Println("          INVESTORS: ")
		for _, investor := range project.ProjectInvestors {
			PrintInvestor(investor)
		}
	}
	if project.Stage == 4 {
		fmt.Println("          Debt Asset Code: ", project.Params.DebtAssetCode)
		fmt.Println("          Payback Asset Code: ", project.Params.PaybackAssetCode)
		fmt.Println("          Balance Left: ", project.Params.BalLeft)
		fmt.Println("          Date Initiated: ", project.Params.DateInitiated)
		fmt.Println("          Date Last Paid: ", project.Params.DateLastPaid)
	}
}

// PrintInvestor pretty prints investors
func PrintInvestor(investor database.Investor) {
	fmt.Println("          Your Public Key is: ", investor.U.PublicKey)
	fmt.Println("          Your Encrypted Seed is: ", investor.U.EncryptedSeed)
	fmt.Println("          Your Voting Balance is: ", investor.VotingBalance)
	fmt.Println("          You have Invested: ", investor.AmountInvested)
	fmt.Println("          Your Invested Assets are: ", investor.InvestedSolarProjects)
	fmt.Println("          Your Username is: ", investor.U.LoginUserName)
	fmt.Println("          Your Password hash is: ", investor.U.LoginPassword)
	fmt.Println("          Your Inspector status is: ", investor.U.Inspector)
	if investor.U.Notification {
		fmt.Println("         Your Email id is: ", investor.U.Email)
	}
	fmt.Println("         Your Local Assets are: ", investor.U.LocalAssets)
}

func PrintUsers(users []database.User) {
	for _, elem := range users {
		PrintUser(elem)
	}
}

func PrintUser(user database.User) {
	fmt.Println("    WELCOME BACK ", user.Name)
	fmt.Println("          Your user index is: ", user.Index)
	fmt.Println("          Your Public Key is: ", user.PublicKey)
	fmt.Println("          Your Encrypted Seed is: ", user.EncryptedSeed)
	fmt.Println("          Your Username is: ", user.LoginUserName)
	fmt.Println("          Your Password hash is: ", user.LoginPassword)
	fmt.Println("          Your KYC status is: ", user.Kyc)
	if user.Notification {
		fmt.Println("         Your Email id is: ", user.Email)
	}
}

// PrintRecipient pretty prints recipients
func PrintRecipient(recipient database.Recipient) {
	fmt.Println("          Your Index is: ", recipient.U.Index)
	fmt.Println("          Your Public Key is: ", recipient.U.PublicKey)
	fmt.Println("          Your Encrypted Seed is: ", recipient.U.EncryptedSeed)
	fmt.Println("          Your Received Assets are: ", recipient.ReceivedSolarProjects)
	fmt.Println("          Your Username is: ", recipient.U.LoginUserName)
	fmt.Println("          Your Password hash is: ", recipient.U.LoginPassword)
	fmt.Println("          Your KYC status is: ", recipient.U.Kyc)
	if recipient.U.Notification {
		fmt.Println("         Your Email id is: ", recipient.U.Email)
	}
	fmt.Println("          Your Device ID is: ", recipient.DeviceId)
	fmt.Println("          Your Device Start Times are: ", recipient.DeviceStarts)
}

func PrintEntity(a solar.Entity) {
	fmt.Println("    WELCOME BACK ", a.U.Name)
	fmt.Println("    			 Your Index is ", a.U.Index)
	fmt.Println("          Your Public Key is: ", a.U.PublicKey)
	fmt.Println("          Your Encrypted Seed is: ", a.U.EncryptedSeed)
	fmt.Println("          Your Username is: ", a.U.LoginUserName)
	fmt.Println("          Your Password hash is: ", a.U.LoginPassword)
	fmt.Println("          Your Address is: ", a.U.Address)
	fmt.Println("          Your Description is: ", a.U.Description)
	fmt.Println("          Your Collateral Amount is: ", a.Collateral)
	fmt.Println("          Your Collateral Data is: ", a.CollateralData)
}

func PrintBalances(balances []horizon.Balance) {
	fmt.Println("   LIST OF ALL YOUR BALANCES: ")
	for _, balance := range balances {
		if balance.Asset.Code == "" {
			fmt.Printf("    ASSET CODE: XLM, ASSET BALANCE: %s\n", balance.Balance)
			continue
		}
		fmt.Printf("    ASSET CODE: %s, ASSET BALANCE: %s\n", balance.Asset.Code, balance.Balance)
	}
}
