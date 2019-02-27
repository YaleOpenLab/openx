package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"

	database "github.com/YaleOpenLab/openx/database"
	solar "github.com/YaleOpenLab/openx/platforms/opensolar"
	rpc "github.com/YaleOpenLab/openx/rpc"
	utils "github.com/YaleOpenLab/openx/utils"
	geo "github.com/martinlindhe/google-geolocate"
)

// GetLocation gets the teller's location
func GetLocation(mapskey string) string {
	// see https://developers.google.com/maps/documentation/geolocation/intro on how
	// to improve location accuracy
	client := geo.NewGoogleGeo(mapskey)
	res, _ := client.Geolocate()
	location := fmt.Sprintf("Lat%fLng%f", res.Lat, res.Lng) // some random format, can be improved upon if necessary
	DeviceLocation = location
	return location
}

// PingRpc pings the platform to see if its up
func PingRpc() error {
	// make a curl request out to lcoalhost and get the ping response
	data, err := rpc.GetRequest(ApiUrl + "/ping")
	if err != nil {
		return err
	}
	var x rpc.StatusResponse
	// now data is in byte, we need the other structure now
	err = x.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	// the result would be the status of the platform
	ColorOutput("PLATFORM STATUS: "+utils.ItoS(x.Code), GreenColor)
	return nil
}

// GetProjectIndex gets a specific project's index
func GetProjectIndex(assetName string) (int, error) {
	data, err := rpc.GetRequest(ApiUrl + "/project/funded")
	if err != nil {
		log.Println(err)
		return -1, err
	}
	var x solar.SolarProjectArray
	err = x.UnmarshalJSON(data)
	if err != nil {
		return -1, err
	}
	for _, elem := range x {
		if elem.DebtAssetCode == assetName {
			return elem.Index, nil
		}
	}
	return -1, errors.New("Not found")
}

// LoginToPlatform logs on to the platform
func LoginToPlatform(username string, pwhash string) error {
	data, err := rpc.GetRequest(ApiUrl + "/recipient/validate?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return err
	}
	var x database.Recipient
	err = x.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	ColorOutput("AUTHENTICATED RECIPIENT", GreenColor)
	LocalRecipient = x
	return nil
}

// ProjectPayback pays back to the platform
func ProjectPayback(assetName string, amount string) error {
	// retrieve project index
	data, err := rpc.GetRequest(ApiUrl + "/recipient/payback?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&projIndex=" + LocalProjIndex + "&assetName=" + LocalProject.DebtAssetCode + "&seedpwd=" +
		LocalSeedPwd + "&amount=" + amount)
	if err != nil {
		return err
	}
	var x rpc.StatusResponse
	err = x.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("PAID!", GreenColor)
		return nil
	}
	return errors.New("Errored out")
}

// SetDeviceId sets the device id of the teller
func SetDeviceId(username string, pwhash string, deviceId string) error {
	data, err := rpc.GetRequest(ApiUrl + "/recipient/deviceId?" + "username=" + username +
		"&pwhash=" + pwhash + "&deviceid=" + deviceId)
	if err != nil {
		return err
	}
	var x rpc.StatusResponse
	err = x.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("PAID!", GreenColor)
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// StoreStartTime stores that start time of this particular instance
func StoreStartTime() error {
	data, err := rpc.GetRequest(ApiUrl + "/recipient/startdevice?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&start=" + utils.I64toS(utils.Unix()))
	if err != nil {
		return err
	}

	var x rpc.StatusResponse
	err = x.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("LOGGED START TIME SUCCESSFULLY!", GreenColor)
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// StoreLocation stores the location of the teller
func StoreLocation(mapskey string) error {
	location := GetLocation(mapskey)
	data, err := rpc.GetRequest(ApiUrl + "/recipient/storelocation?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&location=" + location)
	if err != nil {
		return err
	}

	var x rpc.StatusResponse
	err = x.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("LOGGED LOCATION SUCCESSFULLY!", GreenColor)
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// GetPlatformEmail gets the email of the platform
func GetPlatformEmail() error {
	data, err := rpc.GetRequest(ApiUrl + "/platformemail?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash)
	if err != nil {
		log.Println(err)
		return err
	}

	var x rpc.PlatformEmailResponse
	err = x.UnmarshalJSON(data)
	if err != nil {
		log.Println(err)
		return err
	}
	PlatformEmail = x.Email
	ColorOutput("PLATFORMEMAIL: "+PlatformEmail, GreenColor)
	return nil
}

// SendDeviceShutdownEmail sends a shutdown notice to the platform
func SendDeviceShutdownEmail(tx1 string, tx2 string) error {

	data, err := rpc.GetRequest(ApiUrl + "/tellershutdown?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&projIndex=" + LocalProjIndex + "&deviceId=" + DeviceId +
		"&tx1=" + tx1 + "&tx2=" + tx2)
	if err != nil {
		log.Println(err)
		return err
	}

	var x rpc.StatusResponse
	err = x.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("SENT STOP EMAIL SUCCESSFULLY", GreenColor)
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// GetLocalProjectDetails gets the details of the local project
func GetLocalProjectDetails(projIndex string) (solar.Project, error) {

	var x solar.Project
	body := ApiUrl + "/project/get?index=" + projIndex
	data, err := rpc.GetRequest(body)
	if err != nil {
		log.Println(err)
		return x, err
	}

	err = x.UnmarshalJSON(data)
	if err != nil {
		log.Println(err)
		return x, err
	}

	return x, nil
}

// SendDevicePaybackFailedEmail sends a notification if the payback routine breaks in its execution
func SendDevicePaybackFailedEmail() error {

	data, err := rpc.GetRequest(ApiUrl + "/tellerpayback?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&projIndex=" + LocalProjIndex + "&deviceId=" + DeviceId)
	if err != nil {
		log.Println(err)
		return err
	}

	var x rpc.StatusResponse
	err = x.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("SENT FAILED PAYBACK EMAIL", RedColor)
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// StoreStateHistory stores state history in the data file
func StoreStateHistory(hash string) error {
	data, err := rpc.GetRequest(ApiUrl + "/recipient/ssh?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&hash=" + hash)
	if err != nil {
		log.Println(err)
		return err
	}

	var x rpc.StatusResponse
	err = x.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("SENT FAILED PAYBACK EMAIL", RedColor)
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}
