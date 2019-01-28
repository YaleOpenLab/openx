package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	database "github.com/OpenFinancing/openfinancing/database"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	rpc "github.com/OpenFinancing/openfinancing/rpc"
	utils "github.com/OpenFinancing/openfinancing/utils"
)

func GetRequest(url string) ([]byte, error) {
	// make a curl request out to lcoalhost and get the ping response
	var dummy []byte
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dummy, err
	}
	req.Header.Set("Origin", "localhost")
	res, err := client.Do(req)
	if err != nil {
		return dummy, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func PingRpc() error {
	// make a curl request out to lcoalhost and get the ping response
	data, err := GetRequest(ApiUrl + "/ping")
	if err != nil {
		return err
	}
	var x rpc.PingResponse
	// now data is in byte, we need the other strucutre now
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	// the result would be the status of the platform
	ColorOutput("PLATFORM STATUS: "+x.Status, GreenColor)
	return nil
}

func GetInvestors() error {
	data, err := GetRequest(ApiUrl + "/investor/all")
	if err != nil {
		return err
	}
	var x []database.Investor
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	// the result would be the status of the platform
	ColorOutput("REUQEST SUCCEEDED", GreenColor)
	log.Println(x)
	return nil
}

func GetRecipients() error {
	data, err := GetRequest(ApiUrl + "/recipient/all")
	if err != nil {
		return err
	}
	var x []database.Recipient
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	ColorOutput("REQUEST SUCCEEDED", GreenColor)
	log.Println(x)
	return nil
}

func GetProjectIndex(assetName string) (int, error) {
	data, err := GetRequest(ApiUrl + "/project/funded")
	if err != nil {
		return -1, err
	}
	var x []solar.Project
	err = json.Unmarshal(data, &x)
	if err != nil {
		return -1, err
	}
	for _, elem := range x {
		if elem.Params.DebtAssetCode == assetName {
			return elem.Params.Index, nil
		}
	}
	return -1, fmt.Errorf("Not found")
}

func LoginToPlatForm(username string, pwhash string) error {
	data, err := GetRequest(ApiUrl + "/recipient/validate?" + "LoginUserName=" + username + "&LoginPassword=" + pwhash)
	if err != nil {
		return err
	}
	var x database.Recipient
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	ColorOutput("AUTHENTICATED RECIPIENT", GreenColor)
	LocalRecipient = x
	return nil
}

func ProjectPayback(recpIndex string, assetName string,
	recipientSeed string, amount string) error {
	// retrieve project index
	projIndexI, err := GetProjectIndex(assetName)
	if err != nil {
		return fmt.Errorf("Couldn't pay")
	}
	projIndex := utils.ItoS(projIndexI)
	data, err := GetRequest(ApiUrl + "/recipient/payback?" + "recpIndex=" + recpIndex +
		"&projIndex=" + projIndex + "&assetName=" + assetName + "&recipientSeed=" +
		recipientSeed + "&amount=" + amount + "&platformPublicKey=" + PlatformPublicKey)
	if err != nil {
		return err
	}
	var x rpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Status == 200 {
		ColorOutput("PAID!", GreenColor)
		return nil
	}
	return fmt.Errorf("Errored out")
}

func SetDeviceId(username string, pwhash string, deviceId string) error {
	data, err := GetRequest(ApiUrl + "/recipient/deviceId?" + "LoginUserName=" + username +
		"&LoginPassword=" + pwhash + "&deviceid=" + deviceId)
	if err != nil {
		return err
	}
	var x rpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Status == 200 {
		ColorOutput("PAID!", GreenColor)
		return nil
	}
	return fmt.Errorf("Errored out, didn't receive 200")
}

func StoreStartTime() error {
	data, err := GetRequest(ApiUrl + "/recipient/startdevice?" + "LoginUserName=" + LocalRecipient.U.LoginUserName +
		"&LoginPassword=" + LocalRecipient.U.LoginPassword + "&start=" + utils.I64toS(utils.Unix()))
	if err != nil {
		return err
	}
	log.Println("DATA+", data)
	var x rpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Status == 200 {
		ColorOutput("LOGGED START TIME SUCCESSFULLY!", GreenColor)
		return nil
	}
	return fmt.Errorf("Errored out, didn't receive 200")
}
