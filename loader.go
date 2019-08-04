package main

import (
	// "github.com/pkg/errors"
	"log"

	edb "github.com/Varunram/essentials/database"
	// utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
)

// imagine the loader like in a retro game, loading mainnet

func InitMainnet() error {
	var err error
	// initialize mainnet config stuff and users
	consts.DbDir = consts.HomeDir + "/database/mainnet" // set mainnet db to open in spearate folder
	edb.CreateDirs(consts.DbDir)                        // creates the additional mainnet sub folder

	var inv database.Investor
	allUsers, err := database.RetrieveAllUsers()
	if err != nil {
		return err
	}
	if len(allUsers) == 0 {
		// Create an admin investor
		inv, err = database.NewInvestor("dci@mit.edu", "p", "x", "dci")
		if err != nil {
			return err
		}
		inv.U.Inspector = true
		inv.U.Kyc = true
		inv.U.Admin = true // no handlers for the admin bool, just set it wherever needed.
		inv.U.Reputation = 100000
		inv.U.Notification = true
		err = inv.U.Save()
		if err != nil {
			return err
		}
		err = inv.U.AddEmail("varunramganesh@gmail.com") // change this to something more official later
		if err != nil {
			return err
		}
		err = inv.Save()
		if err != nil {
			return err
		}
		log.Println("Please seed DCI pubkey: ", inv.U.StellarWallet.PublicKey, " with funds")

		// Create an admin recipient
		recp, err := database.NewRecipient("varunramganesh@gmail.com", "p", "x", "vg")
		if err != nil {
			return err
		}
		recp.U.Inspector = true
		recp.U.Kyc = true
		recp.U.Admin = true // no handlers for the admin bool, just set it wherever needed.
		recp.U.Reputation = 100000
		recp.U.Notification = true
		err = recp.U.Save()
		if err != nil {
			return err
		}
		err = recp.U.AddEmail("varunramganesh@gmail.com")
		if err != nil {
			return err
		}
		err = recp.Save()
		if err != nil {
			return err
		}
		log.Println("Please seed Varunram's pubkey: ", recp.U.StellarWallet.PublicKey, " with funds")

		orig, err := opensolar.NewOriginator("martin", "p", "x", "Martin Wainstein", "California", "Project Originator")
		if err != nil {
			return err
		}

		log.Println("Please seed Martin's pubkey: ", orig.U.StellarWallet.PublicKey, " with funds")

		contractor, err := opensolar.NewContractor("samuel", "p", "x", "Samuel Visscher", "Georgia", "Project Contractor")
		if err != nil {
			return err
		}

		log.Println("Please seed Samuel's pubkey: ", contractor.U.StellarWallet.PublicKey, " with funds")

		var project opensolar.Project
		project.Index = 1
		project.TotalValue = 8000
		project.Name = "SU Pasto School, Aibonito"
		project.Metadata = "MIT/Yale Pilot 2"
		project.OriginatorIndex = orig.U.Index
		project.ContractorIndex = contractor.U.Index
		project.EstimatedAcquisition = 5
		project.Stage = 4
		project.MoneyRaised = 0
		// add stuff in here as necessary
		err = project.Save()
		if err != nil {
			return err
		}
	}
	return nil
}
