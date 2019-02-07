package rpc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	//utils "github.com/OpenFinancing/openfinancing/utils"
)

// we need to call the endpoitns and display the stuff returned from that endpoint.
// TODO: in the future, we must define returnde user structs so that we can parse them a bit easier

// function to setup all the particle related endpoints
func setupParticleHandlers() {
	listAllDevices()
	listProductInfo()
	getDeviceInfo()
	pingDevice()
	signalDevice()
	serialNumberInfo()
	getDiagnosticsLast()
	getAllDiagnostics()
	getParticleUserInfo()
	getAllSims()
}

// handler to make easy get requests to a remote host
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

func PutRequest(url string, body string) ([]byte, error) {

	// the body must be the param that you usually pass to curl's -d option
	respbody := strings.NewReader(body)
	var dummy []byte
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, respbody)
	if err != nil {
		return dummy, err
	}
	res, err := client.Do(req)
	if err != nil {
		return dummy, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

type ParticleDevice struct {
	Id                      string `json:"id"`
	Name                    string `json:"name"`
	Last_app                string `json:"last_app"`
	Last_ip_address         string `json:"last_ip_address"`
	Product_id              int    `json:"product_id"`
	Connected               bool   `json:"connected"`
	Platform_id             int    `json:"platform_id"`
	Cellular                bool   `json:"cellular"`
	Notes                   string `json:"notes"`
	Status                  string `json:"status"`
	Serial_number           string `json:"serial_number"`
	Current_build_target    string `json:"current_build_target"`
	System_firmware_version string `json:"system_firmware_version"`
	Default_build_target    string `json:"default_build_target"`
}

type ParticleProductDevice struct {
	Id                                string   `json:"id"`
	Product_id                        int      `json:"product_id"`
	Last_ip_address                   string   `json:"last_ip_address"`
	Last_handshake_at                 string   `json:"last_handshake_at"`
	User_id                           string   `json:"user_id"`
	Online                            bool     `json:"online"`
	Name                              string   `json:"name"`
	Platform_id                       int      `json:"platform_id"`
	Firmware_product_id               int      `json:"firmware_product_id"`
	Quarantined                       bool     `json:"quarantined"`
	Denied                            bool     `json:"denied"`
	Development                       bool     `json:"development"`
	Groups                            []string `json:"groups"`
	Targeted_firmware_release_version string   `json:"targeted_firmware_release_version"`
	System_firmware_version           string   `json:"system_firmware_version"`
	Serial_number                     string   `json:"serial_number"`
	Owner                             string   `json:"owner"`
}

type ParticleProductInfo struct {
	Devices []ParticleProductDevice
}

type ParticlePingResponse struct {
	Online bool `json:"online"`
	Ok     bool `json:"ok"`
}

type SignalResponse struct {
	Id        string `json:"id"`
	Connected bool   `json:"connected"`
	Signaling bool   `json:"signaling"`
}

type SerialNumberResponse struct {
	Ok          bool   `json:"ok"`
	Device_id   string `json:"device_id"`
	Platform_id int    `json:"platform_id"`
}

type ParticleUser struct {
	Username         string   `json:"username"`
	Subscription_ids []string `json:"subscription_ids"`
	AccountInfo      struct {
		First_name       string `json:"first_name"`
		Last_name        string `json:"last_name"`
		Company_name     string `json:"company_name"`
		Business_account bool   `json:"business_account"`
	} `json:"account_info"`
	Team_invites          []string `json:"team_invites"`
	Wifi_device_count     int      `json:"wifi_device_count"`
	Cellular_device_count int      `json:"cellular_device_count"`
}

func listAllDevices() {
	// make a curl request out to lcoalhost and get the ping response
	http.HandleFunc("/particle/devices", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		_, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["accessToken"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		data, err := GetRequest("https://api.particle.io/v1/devices?access_token=" + accessToken)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		var x []ParticleDevice
		// now data is in byte, we need the other strucutre now
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println(err)
			responseHandler(w, r, StatusBadRequest)
			return
		}
		MarshalSend(w, r, x)
	})
}

func listProductInfo() {
	http.HandleFunc("/particle/productinfo", func(w http.ResponseWriter, r *http.Request) {

		_, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["accessToken"] == nil || r.URL.Query()["productInfo"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		productInfo := r.URL.Query()["productInfo"][0]

		data, err := GetRequest("https://api.particle.io/v1/products/" + productInfo + "/devices?access_token=" + accessToken)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		var x ParticleProductInfo
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println(err)
			responseHandler(w, r, StatusBadRequest)
			return
		}
		MarshalSend(w, r, x)

	})
}

func getDeviceInfo() {
	http.HandleFunc("/particle/deviceinfo", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		_, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["accessToken"] == nil || r.URL.Query()["deviceId"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		deviceId := r.URL.Query()["deviceId"][0]

		data, err := GetRequest("https://api.particle.io/v1/devices/" + deviceId + "?access_token=" + accessToken)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		var x ParticleDevice
		// now data is in byte, we need the other strucutre now
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println(err)
			responseHandler(w, r, StatusBadRequest)
			return
		}
		MarshalSend(w, r, x)
	})
}

func pingDevice() {
	http.HandleFunc("/particle/deviceping", func(w http.ResponseWriter, r *http.Request) {

		_, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["accessToken"] == nil || r.URL.Query()["deviceId"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		deviceId := r.URL.Query()["deviceId"][0]
		body := "access_token=" + accessToken // the part which would be passed as -d in curl

		data, err := PutRequest("https://api.particle.io/v1/devices/"+deviceId, body)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		var x ParticlePingResponse
		// now data is in byte, we need the other strucutre now
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println(err)
			responseHandler(w, r, StatusBadRequest)
			return
		}
		MarshalSend(w, r, x)
	})
}

func signalDevice() {
	http.HandleFunc("/particle/devicesignal", func(w http.ResponseWriter, r *http.Request) {

		_, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["signal"] == nil || r.URL.Query()["accessToken"] == nil {
			log.Println("1")
			responseHandler(w, r, StatusBadRequest)
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		deviceId := r.URL.Query()["deviceId"][0]
		signal := r.URL.Query()["signal"][0]
		if signal != "on" && signal != "off" {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		var body string
		if signal == "ok" {
			body = "signal=" + "1" + "accessToken=" + accessToken
		} else {
			body = "signal=" + "0" + "accessToken=" + accessToken
		}

		data, err := PutRequest("https://api.particle.io/v1/devices/"+deviceId, body)
		if err != nil {
			log.Println("2")
			responseHandler(w, r, StatusBadRequest)
			return
		}
		var x SignalResponse
		// now data is in byte, we need the other strucutre now
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println("3")
			responseHandler(w, r, StatusBadRequest)
			return
		}
		MarshalSend(w, r, x)
	})
}

func serialNumberInfo() {
	http.HandleFunc("/particle/getdeviceid", func(w http.ResponseWriter, r *http.Request) {

		_, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["serialNumber"] == nil || r.URL.Query()["accessToken"] == nil {
			log.Println("1")
			responseHandler(w, r, StatusBadRequest)
			return
		}

		serialNumber := r.URL.Query()["serialNumber"][0]
		accessToken := r.URL.Query()["accessToken"][0]

		body := "https://api.particle.io/v1/serial_numbers/" + serialNumber + "?access_token=" + accessToken
		data, err := GetRequest(body)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		var x SerialNumberResponse
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println("3")
			responseHandler(w, r, StatusBadRequest)
			return
		}
		MarshalSend(w, r, x)
	})
}

func getDiagnosticsLast() {
	http.HandleFunc("/particle/diag/last", func(w http.ResponseWriter, r *http.Request) {

		_, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["accessToken"] == nil || r.URL.Query()["deviceId"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		deviceId := r.URL.Query()["deviceId"][0]

		body := "https://api.particle.io/v1/diagnostics/" + deviceId + "/last?access_token=" + accessToken
		log.Println(body)

		data, err := GetRequest(body)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		w.Write(data) // don't know the struct of this since we don't ahve anythign to test against
	})
}

func getAllDiagnostics() {
	http.HandleFunc("/particle/diag/all", func(w http.ResponseWriter, r *http.Request) {

		_, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["accessToken"] == nil || r.URL.Query()["deviceId"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		deviceId := r.URL.Query()["deviceId"][0]

		body := "https://api.particle.io/v1/diagnostics/" + deviceId + "?access_token=" + accessToken

		data, err := GetRequest(body)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		w.Write(data) // don't know the struct of this since we don't ahve anythign to test against
	})
}

func getParticleUserInfo() {
	http.HandleFunc("/particle/user/info", func(w http.ResponseWriter, r *http.Request) {

		_, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["accessToken"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]
		body := "https://api.particle.io/v1/user?access_token=" + accessToken

		data, err := GetRequest(body)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		var x ParticleUser
		err = json.Unmarshal(data, &x)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		MarshalSend(w, r, x)
	})
}

func getAllSims() {
	http.HandleFunc("/particle/sims", func(w http.ResponseWriter, r *http.Request) {

		_, err := UserValidateHelper(w, r)
		if err != nil || r.URL.Query()["accessToken"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		accessToken := r.URL.Query()["accessToken"][0]

		body := "https://api.particle.io/v1/sims?access_token=" + accessToken
		data, err := GetRequest(body)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		w.Write(data)
	})
}
