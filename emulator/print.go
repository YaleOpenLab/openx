package main

import (
	"fmt"

	database "github.com/YaleOpenLab/openx/database"
	solar "github.com/YaleOpenLab/openx/platforms/opensolar"
	"github.com/stellar/go/protocols/horizon"
)

// PrintProjects is a vanity prettyprint function
func PrintProjects(projects []solar.Project) {
	for _, project := range projects {
		PrintProject(project)
	}
}

// PrintProject pretty prints out some stuff that we need in main.go
func PrintProject(project solar.Project) {
	fmt.Println("          PROJECT INDEX: ", project.Index)
	fmt.Println("          Panel Size: ", project.PanelSize)
	fmt.Println("          Total Value: ", project.TotalValue)
	fmt.Println("          Location: ", project.State)
	fmt.Println("          Money Raised: ", project.MoneyRaised)
	fmt.Println("          Metadata: ", project.Metadata)
	fmt.Println("          ETA: ", project.EstimatedAcquisition)
	fmt.Println("          Auction Type: ", project.AuctionType)
	fmt.Println("          PROJECT ORIGINATOR:", project.OriginatorIndex)
	fmt.Println("          PROJECT STAGE: ", project.Stage)
	fmt.Println("          RECIPIENT: ", project.RecipientIndex)
	if project.Stage >= 2 {
		fmt.Println("          PROJECT CONTRACTOR: ", project.ContractorIndex)
		fmt.Println("          Votes: ", project.Votes)
	}
	if project.Stage >= 3 {
		fmt.Println("          Investor Asset Code: ", project.InvestorAssetCode)
		fmt.Println("          INVESTORS: ", project.InvestorIndices)
	}
	if project.Stage == 4 {
		fmt.Println("          Debt Asset Code: ", project.DebtAssetCode)
		fmt.Println("          Payback Asset Code: ", project.PaybackAssetCode)
		fmt.Println("          Balance Left: ", project.BalLeft)
		fmt.Println("          Date Initiated: ", project.DateInitiated)
		fmt.Println("          Date Last Paid: ", project.DateLastPaid)
	}
}

// PrintInvestor pretty prints investors
func PrintInvestor(investor database.Investor) {
	fmt.Println("          Your Public Key is: ", investor.U.StellarWallet.PublicKey)
	fmt.Println("          Your Encrypted Seed is: ", investor.U.StellarWallet.EncryptedSeed)
	fmt.Println("          Your Voting Balance is: ", investor.VotingBalance)
	fmt.Println("          You have Invested: ", investor.AmountInvested)
	fmt.Println("          Your Invested Assets are: ", investor.InvestedSolarProjects)
	fmt.Println("          Your Username is: ", investor.U.Username)
	fmt.Println("          Your Password hash is: ", investor.U.Pwhash)
	fmt.Println("          Your Inspector status is: ", investor.U.Inspector)
	if investor.U.Notification {
		fmt.Println("         Your Email id is: ", investor.U.Email)
	}
	fmt.Println("         Your Local Assets are: ", investor.U.LocalAssets)
	fmt.Println("          Your Recovery Shares are: ", investor.U.RecoveryShares)
	fmt.Println("          Your Secondary Account: ", investor.U.SecondaryWallet)
}

// PrintUsers is a vanity prettyprint function
func PrintUsers(users []database.User) {
	for _, elem := range users {
		PrintUser(elem)
	}
}

// PrintUser is a vanity prettyprint function
func PrintUser(user database.User) {
	fmt.Println("    WELCOME BACK ", user.Name)
	fmt.Println("          Your user index is: ", user.Index)
	fmt.Println("          Your Public Key is: ", user.StellarWallet.PublicKey)
	fmt.Println("          Your Encrypted Seed is: ", user.StellarWallet.EncryptedSeed)
	fmt.Println("          Your Username is: ", user.Username)
	fmt.Println("          Your Password hash is: ", user.Pwhash)
	fmt.Println("          Your KYC status is: ", user.Kyc)
	fmt.Println("          Your Recovery Shares are: ", user.RecoveryShares)
	if user.Notification {
		fmt.Println("         Your Email id is: ", user.Email)
	}
}

// PrintRecipient pretty prints recipients
func PrintRecipient(recipient database.Recipient) {
	fmt.Println("          Your Index is: ", recipient.U.Index)
	fmt.Println("          Your Public Key is: ", recipient.U.StellarWallet.PublicKey)
	fmt.Println("          Your Encrypted Seed is: ", recipient.U.StellarWallet.EncryptedSeed)
	fmt.Println("          Your Received Assets are: ", recipient.ReceivedSolarProjects)
	fmt.Println("          Your Username is: ", recipient.U.Username)
	fmt.Println("          Your Password hash is: ", recipient.U.Pwhash)
	fmt.Println("          Your KYC status is: ", recipient.U.Kyc)
	if recipient.U.Notification {
		fmt.Println("         Your Email id is: ", recipient.U.Email)
	}
	fmt.Println("          Your Device ID is: ", recipient.DeviceId)
	fmt.Println("          Your Device Start Times are: ", recipient.DeviceStarts)
	fmt.Println("          Your Device Location is: ", recipient.DeviceLocation)
	fmt.Println("          Your list of state hashes are: ", recipient.StateHashes)
	fmt.Println("          Your Recovery Shares are: ", recipient.U.RecoveryShares)
}

// PrintEntity is a vanity prettyprint function
func PrintEntity(a solar.Entity) {
	fmt.Println("    WELCOME BACK ", a.U.Name)
	fmt.Println("    			 Your Index is ", a.U.Index)
	fmt.Println("          Your Public Key is: ", a.U.StellarWallet.PublicKey)
	fmt.Println("          Your Encrypted Seed is: ", a.U.StellarWallet.EncryptedSeed)
	fmt.Println("          Your Username is: ", a.U.Username)
	fmt.Println("          Your Password hash is: ", a.U.Pwhash)
	fmt.Println("          Your Address is: ", a.U.Address)
	fmt.Println("          Your Description is: ", a.U.Description)
	fmt.Println("          Your Collateral Amount is: ", a.Collateral)
	fmt.Println("          Your Collateral Data is: ", a.CollateralData)
}

// PrintBalances is a vanity prettyprint function
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
