package main

import (
	"bufio"
	//"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	//"encoding/json"

	consts "github.com/YaleOpenLab/openx/consts"
	ipfs "github.com/YaleOpenLab/openx/ipfs"
	oracle "github.com/YaleOpenLab/openx/oracle"
	utils "github.com/YaleOpenLab/openx/utils"
	//	rpc "github.com/YaleOpenLab/openx/rpc"
	xlm "github.com/YaleOpenLab/openx/xlm"
)

// BlockStamp gets the latest block hash
func BlockStamp() (string, error) {
	hash, err := xlm.GetLatestBlockHash()
	return hash, err
}

// RefreshLogin runs once every 5 minutes in order to fetch the latest recipient details
// for eg, if the recipient loads his balance on the platform, we need it to be reflected on
// the teller
func RefreshLogin(username string, pwhash string) error {
	var err error
	for {
		err = LoginToPlatform(username, pwhash)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(consts.TellerPollInterval)
	}
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
	_, tx, err := xlm.SendXLM(RecpPublicKey, "1", RecpSeed, firstHalf)
	if err != nil {
		return err
	}
	_, tx2, err := xlm.SendXLM(RecpPublicKey, "1", RecpSeed, secondHalf)
	if err != nil {
		return err
	}
	log.Printf("tx hash: %s, tx2 hash: %s", tx, tx2)
	err = SendDeviceShutdownEmail()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// so the teller will be run on the hub and has some data that the platform might need
// The teller must serve some data to other entities as well. So we need to run a server for that
// and it must be over tls for preventing mitm attacks
func CheckPayback() {
	for {
		log.Println("PAYBACK TIME")
		assetName := LocalProject.DebtAssetCode
		amount := oracle.MonthlyBill() // TODO: consumption data must be accumulated from zigbee in the future

		err := ProjectPayback(assetName, amount)
		if err != nil {
			log.Println("Error while paying amount back", err)
			SendDevicePaybackFailedEmail()
		}
		time.Sleep(time.Duration(LocalProject.PaybackPeriod*consts.OneWeekInSecond) * time.Second)
	}
}

// UpdateState hashes the current state of the teller into ipfs and commits the ipfs hash
// to the blockchain
func UpdateState() {
	for {
		subcommand := "Energy production data for this cycle: " + "100" + "W"
		// no spaces since this won't allow us to send in a requerst which has strings in it
		// TODO: replace this with real data rather than fake data that we have here
		// use rest api for ipfs since this may be too heavy to load on a pi. If not, we can shift
		// this to the pi as well to achieve a s tate of good decentralization of information.
		ipfsHash, err := GetIpfsHash("Device ID: " + DeviceId + " UPDATESTATE" + subcommand)
		if err != nil {
			log.Println("Error while fetching ipfs hash", err)
			time.Sleep(consts.TellerPollInterval * time.Second)
		}

		ipfsHash = "STATUPD: " + ipfsHash
		// send _timestamp_ stroops to ourselves, we just pay the network fee of 100 stroops
		// this gives us 10**5 updates per xlm, which is pretty nice, considering that we
		// do about 288 updates a day, this amounts to 347 days' worth updates with 1 XLM
		// memo field restricted to 28 bytes - AAAAAAAAAAAAAAAAAAAAAAAAAAAA
		// we could ideally send the smallest amount of 1 stroop but stellar allows you to
		// send yourself as much money as you want, so we can have any number here
		// we could also time this amount to be the state update number itself.
		// TODO: is this an ideal solution?

		// don't use platform RPCs for interacting with the blockchain
		// But we do need to track this somehow, so maybe hash the device id and "STATUPS: "
		// so we can track if but others viewing the blockchain can't (since the deviceId is assumed
		// to be unique)
		_, hash1, err := xlm.SendXLM(RecpPublicKey, utils.I64toS(utils.Unix()), RecpSeed, ipfsHash[:28])
		if err != nil {
			log.Println(err)
		}

		_, hash2, err := xlm.SendXLM(RecpPublicKey, utils.I64toS(utils.Unix()), RecpSeed, ipfsHash[29:])
		if err != nil {
			log.Println(err)
		}

		// we updated state as hash1 and hash2
		// send email to the platform for this?  maybe overkill
		// TODO: Define structures on the backend that would keep track of this state change
		ColorOutput("Updated State: "+hash1+" "+hash2, MagentaColor)
		time.Sleep(consts.TellerPollInterval * time.Second)
	}
}

// TODO and MWTODO: think upon this problem and arrive at a solution. Might be useful to do
// after we figure out how to store the data that flows into the Rpi on ipfs.
// the problem here is to prove that energy was produced without revealing other details which
// may be confidential. eg. say we produced 30 units of energy and we have all the suporting evidence
// in the form of energy data on our end. now we might not want all of this data to be public but
// the energy production data is required to be public since one should verify that only a given
// amount of REC tokens were generated based on the energy production data.
// One proposed way wouldd be to store the data for periods in ips based on some standard format:
// ipfs(START, DATE: 20/2/19, TIME: <unixtimestamp>, END, DATE: 21/2/19,  TIME: <unixtimestamp> data)

// stream data from the pilot particle instance and write to a file
func storeDataLocal() {
	// store the data that we observe in real time to a file and commit it during the various commit intervals
	// curl "https://api.particle.io/v1/events/PRWhite?access_token=3f7d69aa99956fd77c5466f3f52eb6132f500210"
	// listen to thisparticle endpoint and see if we can stroe this data stream in a file
	path := consts.TellerHomeDir + "/data.txt"

	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: transport}

	body := "https://api.particle.io/v1/devices/events?access_token=3f7d69aa99956fd77c5466f3f52eb6132f500210"
	resp, err := client.Get(body)
	if err != nil {
		log.Println("error while reading from streaming endpoint: ", err)
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	x := make([]byte, 200)
	// open and write to file
	_, err = os.Create(path)
	if err != nil {
		log.Println("error while opening file", err)
		return
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println("error while opening file", err)
		return
	}
	for {
		_, err = reader.Read(x)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("streaming data from particle board: ", string(x))
		_, err = file.Write(x)
		if err != nil {
			log.Println("error while writing to file", err)
			continue
		}
		size, err := file.Stat()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("File size is: ", size.Size())
		if size.Size() == 100000000 { // 100 MB, can be even 1GB but we need to have this file in mem, so cap at 100M.
			// we need to store this in ipfs, delete this file and then commit the ipfs hash as the first line in a new file
			log.Println("SIZE limit reached, taking action")
		}
	}
	// here we want to delete this file every day so that we can save space on the pi
	// commit the ipfs header of the previous block as the first line in the new doc for continuity
	// sort of like a blockchain, so pretty cool that way
}

func commitDataToIpfs() error {
	// retrieve the data from local storage
	path := consts.TellerHomeDir + "/data.txt"
	//time := utils.Timestamp() // don't need to know the start data since we cna get that from the last period's commit
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	dataString := string(data) // have the data as a string

	hash, err := ipfs.AddStringToIpfs(dataString)
	if err != nil {
		return err
	}
	log.Println("IPFS HASH: ", hash)
	// now we need to store this ipfs hash in the list of state updates that this teller
	// has had
	err = StoreStateHistory(hash)
	if err != nil {
		return err
	}

	return nil
}
