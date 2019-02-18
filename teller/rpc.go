package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	database "github.com/YaleOpenLab/openx/database"
	solar "github.com/YaleOpenLab/openx/platforms/opensolar"
	rpc "github.com/YaleOpenLab/openx/rpc"
	utils "github.com/YaleOpenLab/openx/utils"
	geo "github.com/martinlindhe/google-geolocate"
)

func GetLocation(mapskey string) string {
	// see https://developers.google.com/maps/documentation/geolocation/intro on how
	// to improve location accuracy
	client := geo.NewGoogleGeo(mapskey)
	res, _ := client.Geolocate()
	location := fmt.Sprintf("Lat%fLng%f", res.Lat, res.Lng) // some random format, can be improved upon if necessary
	DeviceLocation = location
	return location
}

func PingRpc() error {
	// make a curl request out to lcoalhost and get the ping response
	data, err := rpc.GetRequest(ApiUrl + "/ping")
	if err != nil {
		return err
	}
	var x rpc.StatusResponse
	// now data is in byte, we need the other strucutre now
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	// the result would be the status of the platform
	ColorOutput("PLATFORM STATUS: "+utils.ItoS(x.Code), GreenColor)
	return nil
}

func GetProjectIndex(assetName string) (int, error) {
	data, err := rpc.GetRequest(ApiUrl + "/project/funded")
	if err != nil {
		log.Println(err)
		return -1, err
	}
	var x []solar.Project
	err = json.Unmarshal(data, &x)
	if err != nil {
		return -1, err
	}
	for _, elem := range x {
		if elem.DebtAssetCode == assetName {
			return elem.Index, nil
		}
	}
	return -1, fmt.Errorf("Not found")
}

func LoginToPlatform(username string, pwhash string) error {
	data, err := rpc.GetRequest(ApiUrl + "/recipient/validate?" + "username=" + username + "&pwhash=" + pwhash)
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

func ProjectPayback(assetName string, amount string) error {
	// retrieve project index
	data, err := rpc.GetRequest(ApiUrl + "/recipient/payback?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&projIndex=" + LocalProjIndex + "&assetName=" + LocalProject.DebtAssetCode + "&seedpwd=" +
		LocalSeedPwd + "&amount=" + amount)
	if err != nil {
		return err
	}
	var x rpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("PAID!", GreenColor)
		return nil
	}
	return fmt.Errorf("Errored out")
}

func SetDeviceId(username string, pwhash string, deviceId string) error {
	data, err := rpc.GetRequest(ApiUrl + "/recipient/deviceId?" + "username=" + username +
		"&pwhash=" + pwhash + "&deviceid=" + deviceId)
	if err != nil {
		return err
	}
	var x rpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("PAID!", GreenColor)
		return nil
	}
	return fmt.Errorf("Errored out, didn't receive 200")
}

func StoreStartTime() error {
	data, err := rpc.GetRequest(ApiUrl + "/recipient/startdevice?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&start=" + utils.I64toS(utils.Unix()))
	if err != nil {
		return err
	}

	var x rpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("LOGGED START TIME SUCCESSFULLY!", GreenColor)
		return nil
	}
	return fmt.Errorf("Errored out, didn't receive 200")
}

func StoreLocation(mapskey string) error {
	location := GetLocation(mapskey)
	data, err := rpc.GetRequest(ApiUrl + "/recipient/storelocation?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&location=" + location)
	if err != nil {
		return err
	}

	var x rpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("LOGGED LOCATION SUCCESSFULLY!", GreenColor)
		return nil
	}
	return fmt.Errorf("Errored out, didn't receive 200")
}

func GetPlatformEmail() error {
	data, err := rpc.GetRequest(ApiUrl + "/platformemail?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash)
	if err != nil {
		log.Println(err)
		return err
	}

	var x rpc.PlatformEmailResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		return err
	}
	PlatformEmail = x.Email
	ColorOutput("PLATFORMEMAIL: "+PlatformEmail, GreenColor)
	return nil
}

func SendDeviceShutdownEmail(tx1 string, tx2 string) error {

	data, err := rpc.GetRequest(ApiUrl + "/tellershutdown?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&projIndex=" + LocalProjIndex + "&deviceId=" + DeviceId +
		"&tx1=" + tx1 + "&tx2=" + tx2)
	if err != nil {
		log.Println(err)
		return err
	}

	var x rpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("SENT STOP EMAIL SUCCESSFULLY", GreenColor)
		return nil
	}
	return fmt.Errorf("Errored out, didn't receive 200")
}

func GetIpfsHash(inputString string) (string, error) {
	body := ApiUrl + "/ipfs/hash?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&string=" + inputString

	body = strings.Replace(body, " ", "%20", -1)
	data, err := rpc.GetRequest(body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var x string
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return x, nil
}

func GetLocalProjectDetails(projIndex string) (solar.Project, error) {

	var x solar.Project
	body := ApiUrl + "/project/get?index=" + projIndex
	data, err := rpc.GetRequest(body)
	if err != nil {
		log.Println(err)
		return x, err
	}

	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		return x, err
	}

	return x, nil
}

func SendDevicePaybackFailedEmail() error {

	data, err := rpc.GetRequest(ApiUrl + "/tellerpayback?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&projIndex=" + LocalProjIndex + "&deviceId=" + DeviceId)
	if err != nil {
		log.Println(err)
		return err
	}

	var x rpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("SENT FAILED PAYBACK EMAIL", RedColor)
		return nil
	}
	return fmt.Errorf("Errored out, didn't receive 200")
}

func StoreStateHistory(hash string) error {
	data, err := rpc.GetRequest(ApiUrl + "/recipient/ssh?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&hash=" + hash)
	if err != nil {
		log.Println(err)
		return err
	}

	var x rpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("SENT FAILED PAYBACK EMAIL", RedColor)
		return nil
	}
	return fmt.Errorf("Errored out, didn't receive 200")
}
