package rpc

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	erpc "github.com/Varunram/essentials/rpc"
	consts "github.com/YaleOpenLab/openx/consts"
)

// AnchorRPC is a collection of all Anchor RPC endpoints and their required params
var AnchorRPC = map[int][]string{
	1: {"/user/anchorusd/deposit/intent", "GET"},  // GET
	2: {"/user/anchorusd/deposit/kyc", "GET"},     // GET
	3: {"/user/anchorusd/withdraw/intent", "GET"}, // GET
	4: {"/user/anchorusd/withdraw/kyc", "GET"},    // GET
	5: {"/user/anchorusd/kycinfo", "GET"},         // GET
	6: {"/user/anchorusd/kyc/register", "POST"},   // POST
}

// When a user wants to procure or deal with AnchorUSD, there are a couple things that they need to do:
// 1. Deposit Funds and get AnchorUSD:
// 1a. Create a deposit intent - this wouuld be returned with a 403 since the user has not
// gone through Anchor's KYC process
// 1b. Verify KYC - need to pass relevant parameters to Anchor so they can verfiy from their
// end that the user who wants stablecoin is not sanctioned or something
// 1c. Deposit funds to the relevant bank address sent to the email address of the user.
// AnchorUSD verifies this and deposits AnchorUSD after 2/3 business days.
// 1d. After this, the user heads back to the platform and invests in the platform using the
// AnchorUSD associated with his account
// 2. Withdraw funds denominated in AnchorUSD
// 2a. Create Withdraw intent - returns the identifier / an error if the user hasn't gone through KYC
// 2b. Go through KYC / Withdraw funds - this again takes 2 or 3 business days for AnchorUSD
// to deposit or withdraw funds.
// there are a couple problems with automation in between since there's a delay of 2/3 days with
// each associated fiat operation.

// setupAnchorHandlers sets up all anchorUSD related endpoints
func setupAnchorHandlers() {
	intentDeposit()
	kycDeposit()
	intentWithdraw()
	kycWithdraw()
}

// AnchorIntentResponse defines the intent response struct for AnchorUSD
type AnchorIntentResponse struct {
	Type       string `json:"type"`
	URL        string `json:"url"`
	Identifier string `json:"identifier"`
}

// kycDepositResponse defines the kyc response st ruct for AnchorUSD
type kycDepositResponse struct {
	Error  string
	Result string
	URL    string
}

// GetAndReturnIdentifier is a handler that makes a get request and returns json data
func GetAndReturnIdentifier(w http.ResponseWriter, r *http.Request, body string) (AnchorIntentResponse, error) {
	var x AnchorIntentResponse
	data, err := erpc.GetRequest(body)
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
	data, err := erpc.PostRequest(body, payload)
	if erpc.Err(w, err, erpc.StatusBadRequest, "did not receive success response") {
		return
	}
	var x kycDepositResponse
	err = json.Unmarshal(data, &x)
	if erpc.Err(w, err, erpc.StatusInternalServerError, "did not unmarshal json") {
		return
	}
	erpc.MarshalSend(w, x)
}

// intentDeposit creates an intent to deposit funds in order to fetch AnchorUSD
func intentDeposit() {
	// curl 'https://sandbox-api.anchorusd.com/transfer/deposit?account=GBP3XOFYC6TWUIRZAB7MB6MTUZBCREAYB4E7XKE3OWDP75VU5JB74ZF6&asset_code=USD&email_address=j%40anchorusd.com
	http.HandleFunc(AnchorRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, AnchorRPC[1][2:], AnchorRPC[1][1])
		if err != nil {
			return
		}

		body := consts.AnchorAPI + "transfer/deposit?account=" + prepUser.StellarWallet.PublicKey +
			"&asset_code=USD&email_address=" + prepUser.Email
		x, err := GetAndReturnIdentifier(w, r, body) // we could return the identifier and save it if we have to. But the user has to click through anyawy and we could call the other endpoint from the frontend, so would need to discuss before we do that here
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		prepUser.AnchorKYC.DepositIdentifier = x.Identifier
		prepUser.AnchorKYC.URL = x.URL
		err = prepUser.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, x)
	})
}

// kycDeposit is the kyc workflow involved when a user wants to obtain AnchorUSD
func kycDeposit() {
	http.HandleFunc(AnchorRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, AnchorRPC[2][2:], AnchorRPC[2][1])
		if err != nil {
			return
		}

		body := consts.AnchorAPI + "api/register"
		data := url.Values{}
		data.Set("identifier", prepUser.AnchorKYC.DepositIdentifier)
		data.Set("name", prepUser.AnchorKYC.Name)
		data.Set("birthday[month]", prepUser.AnchorKYC.Birthday.Month)
		data.Set("birthday[day]", prepUser.AnchorKYC.Birthday.Day)
		data.Set("birthday[year]", prepUser.AnchorKYC.Birthday.Year)
		data.Set("tax-country", prepUser.AnchorKYC.Tax.Country)
		data.Set("tax-id-number", prepUser.AnchorKYC.Tax.ID)
		data.Set("address[street-1]", prepUser.AnchorKYC.Address.Street)
		data.Set("address[city]", prepUser.AnchorKYC.Address.City)
		data.Set("address[postal-code]", prepUser.AnchorKYC.Address.Postal)
		data.Set("address[region]", prepUser.AnchorKYC.Address.Region)
		data.Set("address[country]", prepUser.AnchorKYC.Address.Country)
		data.Set("primary-phone-number", prepUser.AnchorKYC.PrimaryPhone)
		data.Set("gender", prepUser.AnchorKYC.Gender)

		payload := strings.NewReader(data.Encode())
		PostAndSend(w, r, body, payload)
	})
}

// intentWithdraw creates an intent to withdraw funds from AnchorUSD
func intentWithdraw() {
	// curl 'https://sandbox-api.anchorusd.com/transfer/withdraw?type=bank_account&asset_code=USD&email_address=j%40anchorusd.com
	http.HandleFunc(AnchorRPC[3][0], func(w http.ResponseWriter, r *http.Request) {
		// the withdraw endpoint doesn't return an identifier and we'd have to parse some stuff ourselves. Ugly hack and we shouldn't really have to do this, should be fixed by Anchor
		prepUser, err := userValidateHelper(w, r, AnchorRPC[3][2:], AnchorRPC[3][1])
		if err != nil {
			return
		}

		// amount can be chosen by the user in the flow on anchor, so no need to handle that here
		body := consts.AnchorAPI + "transfer/withdraw?type=bank_account&asset_code=USD&account=" + prepUser.StellarWallet.PublicKey +
			"&email_address=" + prepUser.Email

		x, err := GetAndReturnIdentifier(w, r, body) // we could return the identifier and save it if we have to. But the user has to click through anyawy and we could call the other endpoint from the frontend, so would need to discuss before we do that here
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		prepUser.AnchorKYC.WithdrawIdentifier = x.Identifier
		prepUser.AnchorKYC.URL = x.URL
		err = prepUser.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, x)
	})
}

// kycDeposit is the kyc workflow involved when a user wants to withdraw fiat
func kycWithdraw() {
	http.HandleFunc(AnchorRPC[4][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, AnchorRPC[4][2:], AnchorRPC[4][1])
		if err != nil {
			return
		}

		body := consts.AnchorAPI + "api/register"
		data := url.Values{}
		if !prepUser.Kyc {
			// we need to call the register KYC first before heading directly to withdrawals
			data.Set("identifier", prepUser.AnchorKYC.DepositIdentifier)
			data.Set("name", prepUser.AnchorKYC.Name)
			data.Set("birthday[month]", prepUser.AnchorKYC.Birthday.Month)
			data.Set("birthday[day]", prepUser.AnchorKYC.Birthday.Day)
			data.Set("birthday[year]", prepUser.AnchorKYC.Birthday.Year)
			data.Set("tax-country", prepUser.AnchorKYC.Tax.Country)
			data.Set("tax-id-number", prepUser.AnchorKYC.Tax.ID)
			data.Set("address[street-1]", prepUser.AnchorKYC.Address.Street)
			data.Set("address[city]", prepUser.AnchorKYC.Address.City)
			data.Set("address[postal-code]", prepUser.AnchorKYC.Address.Postal)
			data.Set("address[region]", prepUser.AnchorKYC.Address.Region)
			data.Set("address[country]", prepUser.AnchorKYC.Address.Country)
			data.Set("primary-phone-number", prepUser.AnchorKYC.PrimaryPhone)
			data.Set("gender", prepUser.AnchorKYC.Gender)
		} else {
			// user has already done kyc while depositing or while signing up, no need for us to do anything here
			data.Set("identifier", prepUser.AnchorKYC.WithdrawIdentifier)
		}
		payload := strings.NewReader(data.Encode())
		PostAndSend(w, r, body, payload)
		// after this, the workflow will be handled by Anchor's KYC and withdrawal system
	})
}

type kycReturn struct {
	AccountID string `json:"account_id"`
	KycStatus string `json:"kyc_status`
}

// getKycStatus gets the status of a user's verification in AnchorUSD
func getKycStatus() {
	http.HandleFunc(AnchorRPC[5][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, AnchorRPC[5][2:], AnchorRPC[5][1])
		if err != nil {
			return
		}

		body := consts.AnchorAPI + "api/accounts/" + prepUser.AnchorKYC.AccountID + "/kyc"

		data, err := erpc.GetRequest(body)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		var ret kycReturn
		err = json.Unmarshal(data, &ret)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		switch ret.KycStatus {
		case "passed":
			erpc.ResponseHandler(w, erpc.StatusOK)
		default:
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
		}

	})
}

type kycR struct {
	URL       string `json:"url"`
	AccountID string `json:"account_id"`
}

// kycRegister is used to register for KYC on AnchorUSD's platform
func kycRegister() {
	http.HandleFunc(AnchorRPC[6][0], func(w http.ResponseWriter, r *http.Request) {
		prepUser, err := userValidateHelper(w, r, AnchorRPC[6][2:], AnchorRPC[6][1])
		if err != nil {
			return
		}

		body := consts.AnchorAPI + "api/register"
		data := url.Values{}
		data.Set("identifier", prepUser.AnchorKYC.DepositIdentifier)
		data.Set("name", prepUser.AnchorKYC.Name)
		data.Set("birthday[month]", prepUser.AnchorKYC.Birthday.Month)
		data.Set("birthday[day]", prepUser.AnchorKYC.Birthday.Day)
		data.Set("birthday[year]", prepUser.AnchorKYC.Birthday.Year)
		data.Set("tax-country", prepUser.AnchorKYC.Tax.Country)
		data.Set("tax-id-number", prepUser.AnchorKYC.Tax.ID)
		data.Set("address[street-1]", prepUser.AnchorKYC.Address.Street)
		data.Set("address[city]", prepUser.AnchorKYC.Address.City)
		data.Set("address[postal-code]", prepUser.AnchorKYC.Address.Postal)
		data.Set("address[region]", prepUser.AnchorKYC.Address.Region)
		data.Set("address[country]", prepUser.AnchorKYC.Address.Country)
		data.Set("primary-phone-number", prepUser.AnchorKYC.PrimaryPhone)
		data.Set("gender", prepUser.AnchorKYC.Gender)

		payload := strings.NewReader(data.Encode())
		retdata, err := erpc.PostRequest(body, payload)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		var ret kycR
		err = json.Unmarshal(retdata, &ret)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		prepUser.AnchorKYC.AccountID = ret.AccountID
		err = prepUser.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, ret)
	})
}
