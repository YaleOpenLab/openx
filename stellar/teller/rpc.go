package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	rpc "github.com/YaleOpenLab/smartPropertyMVP/stellar/rpc"
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
	data, err := GetRequest("http://localhost:8080/ping")
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

func PingInvestors() error {
	data, err := GetRequest("http://localhost:8080/investor/all")
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

func PingRecipients() error {
	data, err := GetRequest("http://localhost:8080/recipient/all")
	if err != nil {
		return err
	}
	var x []database.Recipient
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	ColorOutput("REUQEST SUCCEEDED", GreenColor)
	log.Println(x)
	return nil
}
