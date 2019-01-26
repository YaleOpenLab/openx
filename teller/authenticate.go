package main

import (
	"log"
	"time"

	consts "github.com/OpenFinancing/openfinancing/consts"
	ipfs "github.com/OpenFinancing/openfinancing/ipfs"
	oracle "github.com/OpenFinancing/openfinancing/oracle"
	utils "github.com/OpenFinancing/openfinancing/utils"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
)

type Client struct {
	Info     string
	Location string
	UniqueId string
}

// we should Authenticate whether the code we wanted the rpi to run is actually running
// in order to do this, we have a password + pk/sk based authentication system
var StartHash string
var NowHash string

// BlockStamp gets the latest blockhash from the API
func BlockStamp() (string, error) {
	// get the latest  block here
	hash, err := xlm.GetLatestBlockHash()
	return hash, err
}

// EndHandler runs when the teller shuts down
func EndHandler(t Client) error {
	log.Println("Gracefully shutting down, please do not press any buttons in the process")
	var err error
	NowHash, err = BlockStamp()
	if err != nil {
		return err
	}
	log.Printf("StartHash: %s, NowHash: %s", StartHash, NowHash)
	hashString := "Device Shutting down. Info: " + t.Info + " Device Location: " + t.Location + " Device Unique ID: " + t.UniqueId + " " + StartHash + NowHash
	// need to hash this with ipfs
	ipfsHash, err := ipfs.AddStringToIpfs(hashString)
	if err != nil {
		return err
	}
	log.Println("ipfs hash: ", ipfsHash)
	memoText := "IPFSHASH: " + ipfsHash
	// 10 + 46 (ipfs hash length) characters
	firstHalf := memoText[:28]
	secondHalf := memoText[28:]
	_, tx1, err := xlm.SendXLM(RecpPublicKey, "1", RecpSeed, firstHalf)
	if err != nil {
		return err
	}
	_, tx2, err := xlm.SendXLM(RecpPublicKey, "1", RecpSeed, secondHalf)
	if err != nil {
		return err
	}
	log.Printf("tx1 hash: %s, tx2 hash: %s", tx1, tx2)
	return nil
}

// so the teller will be run on the hub and has some data that the platform might need
// the teller must serve some data to other entities as well. So we need a server for that
// and this must be over tls for preventing mitm attacks and a good tls certificate from an authorized
// provider
func CheckPayback() {
	for {
		log.Println("PAYBACK TIME")

		recpIndex := utils.ItoS(LocalRecipient.U.Index)
		// we only know the debt asset, so retrieve all projects and search for our debt asset
		assetName := LocalRecipient.ReceivedSolarProjects[0] // hardcode for now
		// also might not really be a problem since we assume one recipient per installed solar project
		recipientSeed := RecpSeed
		amount := oracle.MonthlyBill() // TODO: this should be data accumulated from zigbee in the future
		// the platform RecpPublicKey will be static, so can be hardcoded
		// sleep for the interval we want to payback in
		err := ProjectPayback(recpIndex, assetName, recipientSeed, amount)
		if err != nil {
			// payment failed for some reason, notify the developer of this platform
			// and the platform as well
			log.Println(err)
			// <-cleanupDone // terminate and commit hash
		}
		time.Sleep(consts.PaybackInterval * time.Second)
	}
}
