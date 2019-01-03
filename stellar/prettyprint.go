package main

import (
	"fmt"

	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	"github.com/stellar/go/protocols/horizon"
)

func PrintProjects(projects []database.Project) {
	for _, project := range projects {
		PrintProject(project)
	}
}

func PrintProject(project database.Project) {
	PrintParams(project.Params)
	fmt.Println(" PROJECT INDEX: ", project.Params.Index)
	fmt.Println(" PROJECT ORIGINATOR: ", project.Originator)
	fmt.Println(" PROJECT CONTRACTOR: ", project.Contractor)
	fmt.Println(" PROJECT STAGE: ", project.Stage)
}

// PrintParams pretty prints projects
func PrintParams(params database.DBParams) {
	fmt.Println("    ORDER NUMBER: ", params.Index)
	fmt.Println("          Panel Size: ", params.PanelSize)
	fmt.Println("          Total Value: ", params.TotalValue)
	fmt.Println("          Location: ", params.Location)
	fmt.Println("          Money Raised: ", params.MoneyRaised)
	fmt.Println("          Metadata: ", params.Metadata)
	fmt.Println("          Years: ", params.Years)
	// if project.Live {
	fmt.Println("          Investor Asset Code: ", params.INVAssetCode)
	fmt.Println("          Debt Asset Code: ", params.DEBAssetCode)
	fmt.Println("          Payback Asset Code: ", params.PBAssetCode)
	fmt.Println("          Balance Left: ", params.BalLeft)
	// }
	fmt.Println("          Date Initiated: ", params.DateInitiated)
	// if project.Live {
	fmt.Println("          Date Last Paid: ", params.DateLastPaid)
	// }
	fmt.Println("          Recipient: ", params.ProjectRecipient)
	fmt.Println("          Investors: ", params.ProjectInvestors)
	fmt.Println("          Votes: ", params.Votes)
}

// PrintInvestor pretty prints investors
func PrintInvestor(investor database.Investor) {
	fmt.Println("    WELCOME BACK ", investor.U.Name)
	fmt.Println("          Your Public Key is: ", investor.U.PublicKey)
	fmt.Println("          Your Seed is: ", investor.U.Seed)
	fmt.Println("          Your Voting Balance is: ", investor.VotingBalance)
	fmt.Println("          You have Invested: ", investor.AmountInvested)
	fmt.Println("          Your Invested Assets are: ", investor.InvestedAssets)
	fmt.Println("          Your Username is: ", investor.U.LoginUserName)
	fmt.Println("          Your Password hash is: ", investor.U.LoginPassword)
}

// PrintRecipient pretty prints recipients
func PrintRecipient(recipient database.Recipient) {
	fmt.Println("    WELCOME BACK ", recipient.U.Name)
	fmt.Println("          Your Public Key is: ", recipient.U.PublicKey)
	fmt.Println("          Your Seed is: ", recipient.U.Seed)
	fmt.Println("          Your Received Assets are: ", recipient.ReceivedProjects)
	fmt.Println("          Your Username is: ", recipient.U.LoginUserName)
	fmt.Println("          Your Password hash is: ", recipient.U.LoginPassword)
}

// PrintParams pretty prints projects
func PrintPBProjects(projects []database.DBParams) {
	for _, project := range projects {
		if !project.PaidOff {
			fmt.Println("    ORDER NUMBER: ", project.Index)
			fmt.Println("          Panel Size: ", project.PanelSize)
			fmt.Println("          Total Value: ", project.TotalValue)
			fmt.Println("          Location: ", project.Location)
			fmt.Println("          Money Raised: ", project.MoneyRaised)
			fmt.Println("          Metadata: ", project.Metadata)
			fmt.Println("          Years: ", project.Years)
			fmt.Println("          Investor Asset Code: ", project.INVAssetCode)
			fmt.Println("          Debt Asset Code: ", project.DEBAssetCode)
			fmt.Println("          Payback Asset Code: ", project.PBAssetCode)
			fmt.Println("          Balance Left: ", project.BalLeft)
			fmt.Println("          Date Initiated: ", project.DateInitiated)
			fmt.Println("          Date Last Paid: ", project.DateLastPaid)
			fmt.Println("          Investors: ", project.ProjectInvestors)
		}
	}
}

// PrintParams pretty prints projects
func PrintPBProject(project database.DBParams) {
	fmt.Println("    ORDER NUMBER: ", project.Index)
	fmt.Println("          Panel Size: ", project.PanelSize)
	fmt.Println("          Total Value: ", project.TotalValue)
	fmt.Println("          Location: ", project.Location)
	fmt.Println("          Money Raised: ", project.MoneyRaised)
	fmt.Println("          Metadata: ", project.Metadata)
	fmt.Println("          Years: ", project.Years)
	fmt.Println("          Investor Asset Code: ", project.INVAssetCode)
	fmt.Println("          Debt Asset Code: ", project.DEBAssetCode)
	fmt.Println("          Payback Asset Code: ", project.PBAssetCode)
	fmt.Println("          Balance Left: ", project.BalLeft)
	fmt.Println("          Date Initiated: ", project.DateInitiated)
	fmt.Println("          Date Last Paid: ", project.DateLastPaid)
}

func PrintDEB(projects []database.DBParams) {
	for _, project := range projects {
		fmt.Println("          Debt Asset Code: ", project.DEBAssetCode)
	}
}

func PrintUser(user database.User) {
	fmt.Println("    WELCOME BACK ", user.Name)
	fmt.Println("          Your Public Key is: ", user.PublicKey)
	fmt.Println("          Your Seed is: ", user.Seed)
	fmt.Println("          Your Username is: ", user.LoginUserName)
	fmt.Println("          Your Password hash is: ", user.LoginPassword)
}

func PrintEntity(a database.Entity) {
	fmt.Println("    WELCOME BACK ", a.U.Name)
	fmt.Println("    			 Your Index is ", a.U.Index)
	fmt.Println("          Your Public Key is: ", a.U.PublicKey)
	fmt.Println("          Your Seed is: ", a.U.Seed)
	fmt.Println("          Your Username is: ", a.U.LoginUserName)
	fmt.Println("          Your Password hash is: ", a.U.LoginPassword)
	fmt.Println("          Your Address is: ", a.U.Address)
	fmt.Println("          Your Description is: ", a.U.Description)
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
