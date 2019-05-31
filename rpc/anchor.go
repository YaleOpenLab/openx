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
	intentDeposit()
	kycDeposit()
	intentWithdraw()
	kycWithdraw()
}

type AnchorIntentResponse struct {
	Type       string `json:"type"`
	Url        string `json:"url"`
	Identifier string `json:"identifier"`
}

type kycDepositResponse struct {
	Error  string
	Result string
	Url    string
}

// GetAndReturnIdentifier is a handler that makes a get request and returns json data
func GetAndReturnIdentifier(w http.ResponseWriter, r *http.Request, body string) (AnchorIntentResponse, error) {
	var x AnchorIntentResponse
	data, err := GetRequest(body)
	if err != nil {
		log.Println("did not get response", err)
		return x, err
	}
	// now data is in byte, we need the other structure now
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println("did not unmarshal json", err)
		return x, err
	}
	return x, nil
}

// PostAndSend is a handler that POSTs data and returns the response
func PostAndSend(w http.ResponseWriter, r *http.Request, body string, payload io.Reader) {
	data, err := PostRequest(body, payload)
	if err != nil {
		log.Println("did not receive success response", err)
		responseHandler(w, StatusBadRequest)
		return
	}
	var x kycDepositResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println("did not unmarshal json", err)
		responseHandler(w, StatusInternalServerError)
		return
	}
	MarshalSend(w, x)
}

func intentDeposit() {
	// curl 'https://sandbox-api.anchorusd.com/transfer/deposit?account=GBP3XOFYC6TWUIRZAB7MB6MTUZBCREAYB4E7XKE3OWDP75VU5JB74ZF6&asset_code=USD&email_address=j%40anchorusd.com
	http.HandleFunc("/user/anchorusd/deposit/intent", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}

		body := "https://sandbox-api.anchorusd.com/transfer/deposit?account=" + prepUser.PublicKey +
			"&asset_code=USD&email_address=" + prepUser.Email
		x, err := GetAndReturnIdentifier(w, r, body) // we could return the identifier and save it if we have to. But the user has to click through anyawy and we could call the other endpoint from the frontend, so would need to discuss before we do that here
		if err != nil {
			responseHandler(w, StatusInternalServerError)
			return
		}

		prepUser.AnchorKYC.DepositIdentifier = x.Identifier
		err = prepUser.Save()
		if err != nil {
			responseHandler(w, StatusInternalServerError)
			return
		}

		MarshalSend(w, x)
	})
}

func kycDeposit() {
	http.HandleFunc("/user/anchorusd/deposit/kyc", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}

		body := "https://sandbox-api.anchorusd.com/api/register"
		data := url.Values{}
		data.Set("identifier", prepUser.AnchorKYC.DepositIdentifier)
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

func intentWithdraw() {
	// curl 'https://sandbox-api.anchorusd.com/transfer/withdraw?type=bank_account&asset_code=USD&email_address=j%40anchorusd.com
	http.HandleFunc("/user/anchorusd/withdraw/intent", func(w http.ResponseWriter, r *http.Request) {
		// the withdraw endpoint doesn't return an identifier and we'd have to parse some stuff ourselves. Ugly hack and we shouldn't really ahve to do this, should be fixed by Anchor
		checkGet(w, r)
		checkOrigin(w, r)

		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}

		// amount can be chosen by the user in the flow on anchor, so no need to handle that here
		body := "https://sandbox-api.anchorusd.com/transfer/withdraw?type=bank_account&asset_code=USD&account=" + prepUser.PublicKey +
			"&email_address=" + prepUser.Email

		x, err := GetAndReturnIdentifier(w, r, body) // we could return the identifier and save it if we have to. But the user has to click through anyawy and we could call the other endpoint from the frontend, so would need to discuss before we do that here
		if err != nil {
			responseHandler(w, StatusInternalServerError)
			return
		}

		prepUser.AnchorKYC.WithdrawIdentifier = x.Identifier
		err = prepUser.Save()
		if err != nil {
			responseHandler(w, StatusInternalServerError)
			return
		}

		MarshalSend(w, x)
	})
}

func kycWithdraw() {
	http.HandleFunc("/user/anchorusd/withdraw/kyc", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		prepUser, err := UserValidateHelper(w, r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}

		body := "https://sandbox-api.anchorusd.com/api/register"
		data := url.Values{}
		data.Set("identifier", prepUser.AnchorKYC.WithdrawIdentifier) // TODO: the deposit API doesn't parse identifiers. Should be fixed on AnchorUSD's end.
		payload := strings.NewReader(data.Encode())
		PostAndSend(w, r, body, payload)
	})
}
