package main

import (
	"fmt"

	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	"github.com/stellar/go/protocols/horizon"
)

// PrintOrder pretty prints orders
func PrintOrders(orders []database.Order) {
	for _, order := range orders {
		PrintOrder(order)
	}
}

// PrintOrder pretty prints orders
func PrintOrder(order database.Order) {
	fmt.Println("    ORDER NUMBER: ", order.Index)
	fmt.Println("          Panel Size: ", order.PanelSize)
	fmt.Println("          Total Value: ", order.TotalValue)
	fmt.Println("          Location: ", order.Location)
	fmt.Println("          Money Raised: ", order.MoneyRaised)
	fmt.Println("          Metadata: ", order.Metadata)
	fmt.Println("          Years: ", order.Years)
	// if order.Live {
	fmt.Println("          Investor Asset Code: ", order.INVAssetCode)
	fmt.Println("          Debt Asset Code: ", order.DEBAssetCode)
	fmt.Println("          Payback Asset Code: ", order.PBAssetCode)
	fmt.Println("          Balance Left: ", order.BalLeft)
	// }
	fmt.Println("          Date Initiated: ", order.DateInitiated)
	// if order.Live {
	fmt.Println("          Date Last Paid: ", order.DateLastPaid)
	// }
	fmt.Println("          Recipient: ", order.OrderRecipient)
	fmt.Println("          Investors: ", order.OrderInvestors)
	fmt.Println("          Votes: ", order.Votes)
	fmt.Println("          Order Stage: ", order.Stage)
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
	fmt.Println("          Your Received Assets are: ", recipient.ReceivedOrders)
	fmt.Println("          Your Username is: ", recipient.U.LoginUserName)
	fmt.Println("          Your Password hash is: ", recipient.U.LoginPassword)
}

// PrintOrder pretty prints orders
func PrintPBOrders(orders []database.Order) {
	for _, order := range orders {
		if !order.PaidOff {
			fmt.Println("    ORDER NUMBER: ", order.Index)
			fmt.Println("          Panel Size: ", order.PanelSize)
			fmt.Println("          Total Value: ", order.TotalValue)
			fmt.Println("          Location: ", order.Location)
			fmt.Println("          Money Raised: ", order.MoneyRaised)
			fmt.Println("          Metadata: ", order.Metadata)
			fmt.Println("          Years: ", order.Years)
			fmt.Println("          Investor Asset Code: ", order.INVAssetCode)
			fmt.Println("          Debt Asset Code: ", order.DEBAssetCode)
			fmt.Println("          Payback Asset Code: ", order.PBAssetCode)
			fmt.Println("          Balance Left: ", order.BalLeft)
			fmt.Println("          Date Initiated: ", order.DateInitiated)
			fmt.Println("          Date Last Paid: ", order.DateLastPaid)
			fmt.Println("          Investors: ", order.OrderInvestors)
		}
	}
}

// PrintOrder pretty prints orders
func PrintPBOrder(order database.Order) {
	fmt.Println("    ORDER NUMBER: ", order.Index)
	fmt.Println("          Panel Size: ", order.PanelSize)
	fmt.Println("          Total Value: ", order.TotalValue)
	fmt.Println("          Location: ", order.Location)
	fmt.Println("          Money Raised: ", order.MoneyRaised)
	fmt.Println("          Metadata: ", order.Metadata)
	fmt.Println("          Years: ", order.Years)
	fmt.Println("          Investor Asset Code: ", order.INVAssetCode)
	fmt.Println("          Debt Asset Code: ", order.DEBAssetCode)
	fmt.Println("          Payback Asset Code: ", order.PBAssetCode)
	fmt.Println("          Balance Left: ", order.BalLeft)
	fmt.Println("          Date Initiated: ", order.DateInitiated)
	fmt.Println("          Date Last Paid: ", order.DateLastPaid)
}

func PrintDEB(orders []database.Order) {
	for _, order := range orders {
		fmt.Println("          Debt Asset Code: ", order.DEBAssetCode)
	}
}

// PrintOrder pretty prints orders
func PrintProposedContract(order database.Order) {
	fmt.Println("          Proposed Contract: ")
	fmt.Println("          Panel Size: ", order.PanelSize)
	fmt.Println("          Total Value: ", order.TotalValue)
	fmt.Println("          Location: ", order.Location)
	fmt.Println("          Metadata: ", order.Metadata)
	fmt.Println("          Years: ", order.Years)
	fmt.Println("          Date Initiated: ", order.DateInitiated)
	fmt.Println("          Recipient: ", order.OrderRecipient)
	fmt.Println("          Investors: ", order.OrderInvestors)
}

func PrintUser(user database.User) {
	fmt.Println("    WELCOME BACK ", user.Name)
	fmt.Println("          Your Public Key is: ", user.PublicKey)
	fmt.Println("          Your Seed is: ", user.Seed)
	fmt.Println("          Your Username is: ", user.LoginUserName)
	fmt.Println("          Your Password hash is: ", user.LoginPassword)
}

func PrintContractEntity(a database.ContractEntity) {
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
