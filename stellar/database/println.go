package database

import (
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
)

// PrettyPrintOrder pretty prints orders
func PrettyPrintOrders(orders []Order) {
	for _, order := range orders {
		fmt.Println("    ORDER NUMBER: ", order.Index)
		fmt.Println("          Panel Size: ", order.PanelSize)
		fmt.Println("          Total Value: ", order.TotalValue)
		fmt.Println("          Location: ", order.Location)
		fmt.Println("          Money Raised: ", order.MoneyRaised)
		fmt.Println("          Metadata: ", order.Metadata)
		fmt.Println("          Years: ", order.Years)
		if order.Live {
			fmt.Println("          Investor Asset Code: ", order.INVAssetCode)
			fmt.Println("          Debt Asset Code: ", order.DEBAssetCode)
			fmt.Println("          Payback Asset Code: ", order.PBAssetCode)
			fmt.Println("          Balance Left: ", order.BalLeft)
		}
		fmt.Println("          Date Initiated: ", order.DateInitiated)
		if order.Live {
			fmt.Println("          Date Last Paid: ", order.DateLastPaid)
		}
	}
}

// PrettyPrintOrder pretty prints orders
func PrettyPrintOrder(order Order) {
	fmt.Println("    ORDER NUMBER: ", order.Index)
	fmt.Println("          Panel Size: ", order.PanelSize)
	fmt.Println("          Total Value: ", order.TotalValue)
	fmt.Println("          Location: ", order.Location)
	fmt.Println("          Money Raised: ", order.MoneyRaised)
	fmt.Println("          Metadata: ", order.Metadata)
	fmt.Println("          Years: ", order.Years)
	if order.Live {
		fmt.Println("          Investor Asset Code: ", order.INVAssetCode)
		fmt.Println("          Debt Asset Code: ", order.DEBAssetCode)
		fmt.Println("          Payback Asset Code: ", order.PBAssetCode)
		fmt.Println("          Balance Left: ", order.BalLeft)
	}
	fmt.Println("          Date Initiated: ", order.DateInitiated)
	if order.Live {
		fmt.Println("          Date Last Paid: ", order.DateLastPaid)
	}
}

func InsertDummyData() error {
	var err error
	// populate database with dumym data
	var order1 Order

	order1.Index = 1
	order1.PanelSize = "100 1000 sq.ft homes each with their own private spaces for luxury"
	order1.TotalValue = 14000
	order1.Location = "India Basin, San Francisco"
	order1.MoneyRaised = 0
	order1.Metadata = "India Basin is an upcoming creative project based in San Francisco that seeks to invite innovators from all around to participate"
	order1.Live = false
	order1.INVAssetCode = ""
	order1.DEBAssetCode = ""
	order1.PBAssetCode = ""
	order1.DateInitiated = utils.Timestamp()
	order1.Years = 3
	order1.RecipientName = "Martin" // this is not the username of the recipient
	err = InsertOrder(order1)
	if err != nil {
		return fmt.Errorf("Error inserting order into db")
	}

	order1.Index = 2
	order1.PanelSize = "180 1200 sq.ft homes in a high rise building 0.1mi from Kendall Square"
	order1.TotalValue = 30000
	order1.Location = "Kendall Square, Boston"
	order1.MoneyRaised = 0
	order1.Metadata = "Kendall Square is set in the heart of Cambridge and is a popular startup IT hub"
	order1.Live = false
	order1.INVAssetCode = ""
	order1.DEBAssetCode = ""
	order1.PBAssetCode = ""
	order1.DateInitiated = utils.Timestamp()
	order1.Years = 5
	order1.RecipientName = "Martin" // this is not the username of the recipient

	err = InsertOrder(order1)
	if err != nil {
		return fmt.Errorf("Error inserting order into db")
	}

	order1.Index = 3
	order1.PanelSize = "260 1500 sq.ft homes set in a medieval cathedral style construction"
	order1.TotalValue = 40000
	order1.Location = "Trafalgar Square, London"
	order1.MoneyRaised = 0
	order1.Metadata = "Trafalgar Square is set in the heart of London's financial district, with big banks all over"
	order1.Live = false
	order1.INVAssetCode = ""
	order1.DEBAssetCode = ""
	order1.PBAssetCode = ""
	order1.DateInitiated = utils.Timestamp()
	order1.Years = 7
	order1.RecipientName = "Martin" // this is not the username of the recipient

	err = InsertOrder(order1)
	if err != nil {
		return fmt.Errorf("Error inserting order into db")
	}

	var inv Investor
	allInvs, err := RetrieveAllInvestors()
	if err != nil {
		log.Fatal(err)
	}
	if len(allInvs) == 0 {
		inv.Index = 1
		inv.LoginUserName = "john"
		inv.LoginPassword = "e9a75486736a550af4fea861e2378305c4a555a05094dee1dca2f68afea49cc3a50e8de6ea131ea521311f4d6fb054a146e8282f8e35ff2e6368c1a62e909716"
		inv.Name = "John"
		inv.Seed, inv.PublicKey, err = xlm.GetKeyPair()
		if err != nil {
			log.Fatal(err)
		}
		err = InsertInvestor(inv)
		if err != nil {
			log.Fatal(err)
		}
	} else if len(allInvs) == 1 {
		// don't do anything
	}

	allRecs, err := RetrieveAllRecipients()
	if err != nil {
		log.Fatal(err)
	}
	if len(allRecs) == 0 {
		var rec Recipient
		rec.Index = 1
		rec.LoginUserName = "martin"
		rec.LoginPassword = "e9a75486736a550af4fea861e2378305c4a555a05094dee1dca2f68afea49cc3a50e8de6ea131ea521311f4d6fb054a146e8282f8e35ff2e6368c1a62e909716"
		rec.Name = "Martin"
		rec.Seed, rec.PublicKey, err = xlm.GetKeyPair()
		if err != nil {
			log.Fatal(err)
		}
		err = InsertRecipient(rec)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// PrettyPrintInvestor pretty prints investors
func PrettyPrintInvestor(investor Investor) {
	fmt.Println("    WELCOME BACK ", investor.Name)
	fmt.Println("          Your Public Key is: ", investor.PublicKey)
	fmt.Println("          Your Seed is: ", investor.Seed)
	fmt.Println("          You have Invested: ", investor.AmountInvested)
	fmt.Println("          Your Invested Assets are: ", investor.InvestedAssets)
	fmt.Println("          Your Username is: ", investor.LoginUserName)
	fmt.Println("          Your Password hash is: ", investor.LoginPassword)
}

// PrettyPrintRecipient pretty prints recipients
func PrettyPrintRecipient(recipient Recipient) {
	fmt.Println("    WELCOME BACK ", recipient.Name)
	fmt.Println("          Your Public Key is: ", recipient.PublicKey)
	fmt.Println("          Your Seed is: ", recipient.Seed)
	fmt.Println("          Your Received Assets are: ", recipient.ReceivedOrders)
	fmt.Println("          Your Username is: ", recipient.LoginUserName)
	fmt.Println("          Your Password hash is: ", recipient.LoginPassword)
}
