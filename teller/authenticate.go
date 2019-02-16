package main

import (
	"log"
	"time"

	consts "github.com/YaleOpenLab/openx/consts"
	ipfs "github.com/YaleOpenLab/openx/ipfs"
	oracle "github.com/YaleOpenLab/openx/oracle"
	xlm "github.com/YaleOpenLab/openx/xlm"
)

func BlockStamp() (string, error) {
	// get the latest  block here
	hash, err := xlm.GetLatestBlockHash()
	return hash, err
}

// EndHandler runs when the teller shuts down. Records the start time and location of the
// device in ipfs and commits it as two transactions to the blockchain
func EndHandler() error {
	log.Println("Gracefully shutting down, please do not press any button in the process")
	var err error
	NowHash, err = BlockStamp()
	if err != nil {
		return err
	}
	log.Printf("StartHash: %s, NowHash: %s", StartHash, NowHash)
	hashString := "Device Shutting down. Info: " + DeviceInfo + " Device Location: " + DeviceLocation + " Device Unique ID: " + DeviceId + " " + StartHash + NowHash
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
	err = SendDeviceShutdownEmail()
	if err != nil {
		log.Println(err)
	}
	return nil
}

// so the teller will be run on the hub and has some data that the platform might need
// the teller must serve some data to other entities as well. So we need a server for that
// and this must be over tls for preventing mitm attacks and a good tls certificate from an authorized
// provider
func CheckPayback() {
	for {
		log.Println("PAYBACK TIME")

		if len(LocalRecipient.ReceivedSolarProjects) == 0 {
			// if the recipient has no received solar projects, quit. This is done in order to test the teller
			return
		}
		// we only know the debt asset, so retrieve all projects and search for our debt asset
		assetName := LocalRecipient.ReceivedSolarProjects[0] // hardcode for now
		// also might not really be a problem since we assume one recipient per installed solar project
		amount := oracle.MonthlyBill() // TODO: consumption data must be accumulated from zigbee in the future
		err := ProjectPayback(assetName, amount)
		if err != nil {
			// payment failed for some reason, notify the developer of this platform
			// and the platform as well
			log.Println("Error while paying amount back", err)
		}
		time.Sleep(consts.PaybackInterval * time.Minute) // TODO" this is based on the agreed upon payback, change
	}
}
