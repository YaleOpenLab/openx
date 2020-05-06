// +build all anchor

package rpc

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"testing"

	xlm "github.com/Varunram/essentials/xlm"
)

func TestAnchorEndpoints(t *testing.T) {
	_, pubkey, err := xlm.GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	email := "random@test.com"
	body := "https://sandbox-api.anchorusd.com/transfer/deposit?account=" + pubkey +
		"&asset_code=USD&email_address=" + email

	var x AnchorIntentResponse
	dataJson, err := GetRequest(body)
	if err != nil {
		t.Fatal(err)
	}
	// now data is in byte, we need the other structure now
	err = json.Unmarshal(dataJson, &x)
	if err != nil {
		t.Fatal(err)
	}
	depositIdentifier := x.Identifier

	body = "https://sandbox-api.anchorusd.com/api/register"
	data := url.Values{}
	data.Set("identifier", depositIdentifier)
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

	dataPost, err := PostRequest(body, payload)
	if err != nil {
		t.Fatal(err)
	}
	var xy kycDepositResponse
	err = json.Unmarshal(dataPost, &xy)
	if err != nil {
		t.Fatal(err)
		return
	}

	log.Println(xy)
	// if we fall down till here, deposit functions should be fine

	body = "https://sandbox-api.anchorusd.com/transfer/withdraw?type=bank_account&asset_code=USD&account=" + pubkey +
		"&email_address=" + email

	dataJson, err = GetRequest(body)
	if err != nil {
		t.Fatal(err)
	}
	// now data is in byte, we need the other structure now
	err = json.Unmarshal(dataJson, &x)
	if err != nil {
		t.Fatal(err)
	}

	withdrawIdentifier := x.Identifier

	body = "https://sandbox-api.anchorusd.com/api/register"
	data = url.Values{}
	data.Set("identifier", withdrawIdentifier)
	payload = strings.NewReader(data.Encode())

	dataPost, err = PostRequest(body, payload)
	if err != nil {
		t.Fatal(err)
	}
	var xz kycDepositResponse
	err = json.Unmarshal(dataPost, &xz)
	if err != nil {
		t.Fatal(err)
		return
	}
	log.Println(xz)
}
