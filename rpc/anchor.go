package rpc

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func setupAnchorHandlers() {
	intentAnchor()
	kycAnchor()
}

type AnchorIntentResponse struct {
	Type       string `json:"type"`
	Url        string `json:"url"`
	Identifier string `json:"identifier"`
}

type KycAnchorResponse struct {
	Error  string
	Result string
	Url    string
}

// PostAndSend is a handler that POSTs data and returns the response
func PostAndSend(w http.ResponseWriter, r *http.Request, body string, payload io.Reader) {
	data, err := PostRequest(body, payload)
	if err != nil {
		log.Println("did not receive success response", err)
		responseHandler(w, StatusBadRequest)
		return
	}
	var x KycAnchorResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println("did not unmarshal json", err)
		responseHandler(w, StatusInternalServerError)
		return
	}
	MarshalSend(w, x)
}

func intentAnchor() {
	// curl 'https://sandbox-api.anchorusd.com/transfer/deposit?account=GBP3XOFYC6TWUIRZAB7MB6MTUZBCREAYB4E7XKE3OWDP75VU5JB74ZF6&asset_code=USD&email_address=j%40anchorusd.com
	http.HandleFunc("/user/anchorusd/intent", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}

		body := "https://sandbox-api.anchorusd.com/transfer/deposit?account=" + prepUser.PublicKey +
			"&asset_code=USD&email_address=" + prepUser.Email
		var x AnchorIntentResponse
		GetAndSendJson(w, r, body, x)
	})
}

func kycAnchor() {
	http.HandleFunc("/user/anchorusd/kyc", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}

		body := "https://sandbox-api.anchorusd.com/api/register"
		data := url.Values{}
		data.Set("identifier", "b5836f518685c372fe5a013475ccbf37")
		data.Set("name", "Test User")
		data.Set("birthday[month]", "6")
		data.Set("birthday[day]", "8")
		data.Set("birthday[year]", "1993")
		data.Set("tax-country", "US")
		data.Set("tax-id-number", "111111111")
		data.Set("address[street-1]", "123 4 Street")
		data.Set("address[city]", "Anytown")
		data.Set("address[postal-code]", "94107")
		data.Set("address[region]", "CA")
		data.Set("address[country]", "US")
		data.Set("primary-phone-number", "+14151111111")
		data.Set("gender", "male")

		payload := strings.NewReader(data.Encode())
		PostAndSend(w, r, body, payload)
	})
}
