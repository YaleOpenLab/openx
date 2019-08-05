package rpc

// anchor connects with AnchorUSD's endpoints and returns the relevant endpoints in order
// for us to parse correctly. Broadly when a user wants to procure or deal with AnchorUSD,
// there are a couple things that he needs to do:
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
// there are a couple problems with automation in betwee nsince there's a delay of 2/3 days with
// each associated fiat operation. Hopefully, we can solve this in some way or the other.
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
	if err != nil {
		log.Println("did not receive success response", err)
		erpc.ResponseHandler(w, erpc.StatusBadRequest)
		return
	}
	var x kycDepositResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println("did not unmarshal json", err)
		erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		return
	}
	erpc.MarshalSend(w, x)
}

func intentDeposit() {
	// curl 'https://sandbox-api.anchorusd.com/transfer/deposit?account=GBP3XOFYC6TWUIRZAB7MB6MTUZBCREAYB4E7XKE3OWDP75VU5JB74ZF6&asset_code=USD&email_address=j%40anchorusd.com
	http.HandleFunc("/user/anchorusd/deposit/intent", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepUser, err := CheckReqdParams(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		body := consts.AnchorAPI + "transfer/deposit?account=" + prepUser.StellarWallet.PublicKey +
			"&asset_code=USD&email_address=" + prepUser.Email
		x, err := GetAndReturnIdentifier(w, r, body) // we could return the identifier and save it if we have to. But the user has to click through anyawy and we could call the other endpoint from the frontend, so would need to discuss before we do that here
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		prepUser.AnchorKYC.DepositIdentifier = x.Identifier
		err = prepUser.Save()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

func kycDeposit() {
	http.HandleFunc("/user/anchorusd/deposit/kyc", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepUser, err := CheckReqdParams(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
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
		data.Set("tax-id-number", prepUser.AnchorKYC.Tax.Id)
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

func intentWithdraw() {
	// curl 'https://sandbox-api.anchorusd.com/transfer/withdraw?type=bank_account&asset_code=USD&email_address=j%40anchorusd.com
	http.HandleFunc("/user/anchorusd/withdraw/intent", func(w http.ResponseWriter, r *http.Request) {
		// the withdraw endpoint doesn't return an identifier and we'd have to parse some stuff ourselves. Ugly hack and we shouldn't really have to do this, should be fixed by Anchor
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepUser, err := CheckReqdParams(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		// amount can be chosen by the user in the flow on anchor, so no need to handle that here
		body := consts.AnchorAPI + "transfer/withdraw?type=bank_account&asset_code=USD&account=" + prepUser.StellarWallet.PublicKey +
			"&email_address=" + prepUser.Email

		x, err := GetAndReturnIdentifier(w, r, body) // we could return the identifier and save it if we have to. But the user has to click through anyawy and we could call the other endpoint from the frontend, so would need to discuss before we do that here
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		prepUser.AnchorKYC.WithdrawIdentifier = x.Identifier
		err = prepUser.Save()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

func kycWithdraw() {
	http.HandleFunc("/user/anchorusd/withdraw/kyc", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepUser, err := CheckReqdParams(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
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
			data.Set("tax-id-number", prepUser.AnchorKYC.Tax.Id)
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
		PostAndSend(w, r, body, payload) // send the payload and response, will be handled by Anchor's KYC and withdrawal system
	})
}
