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
	"github.com/stellar/go/protocols/horizon"
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
	var x rpc.StatusResponse
	// now data is in byte, we need the other strucutre now
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	// the result would be the status of the platform
	ColorOutput("PLATFORM STATUS: "+utils.ItoS(x.Status), GreenColor)
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

func ProjectPayback(recpIndex string, assetName string,
	recipientSeed string, amount string) error {
	// retrieve project index
	projIndexI, err := GetProjectIndex(assetName)
	if err != nil {
		return fmt.Errorf("Couldn't pay")
	}
	projIndex := utils.ItoS(projIndexI)
	PlatformPublicKey := "GDULAIM6N6SIW7MWS3NDJPY3UIFOHSM4766WQ6O6EKFDBC7PF53VKYLY" // this will be public, so hardcode
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

func RetrieveProject(stage float64) ([]solar.Project, error) {
	// retrieve project at a particular stage
	var x []solar.Project
	switch stage {
	case 0:
		data, err := GetRequest(ApiUrl + "/project/preorigin")
		if err != nil {
			return x, err
		}
		var x []solar.Project
		err = json.Unmarshal(data, &x)
		if err != nil {
			return x, err
		}
		return x, nil
	case 1:
		data, err := GetRequest(ApiUrl + "/project/origin")
		if err != nil {
			return x, err
		}
		var x []solar.Project
		err = json.Unmarshal(data, &x)
		if err != nil {
			return x, err
		}
		return x, nil
	case 2:
		data, err := GetRequest(ApiUrl + "/project/proposed")
		if err != nil {
			return x, err
		}
		var x []solar.Project
		err = json.Unmarshal(data, &x)
		if err != nil {
			return x, err
		}
		return x, nil
	case 3:
		data, err := GetRequest(ApiUrl + "/project/final")
		if err != nil {
			return x, err
		}
		var x []solar.Project
		err = json.Unmarshal(data, &x)
		if err != nil {
			return x, err
		}
		return x, nil
	case 4:
		data, err := GetRequest(ApiUrl + "/project/funded")
		if err != nil {
			return x, err
		}
		var x []solar.Project
		err = json.Unmarshal(data, &x)
		if err != nil {
			return x, err
		}
		return x, nil
	}
	return x, nil
}

func GetBalances(username string, pwhash string) ([]horizon.Balance, error) {
	// get the balance from the balances API
	var x []horizon.Balance
	data, err := GetRequest(ApiUrl + "/user/balances?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func GetXLMBalance(username string, pwhash string) (string, error) {
	// get the balance from the balances API
	var x string
	data, err := GetRequest(ApiUrl + "/user/balance/xlm?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func GetAssetBalance(username string, pwhash string, asset string) (string, error) {
	// get the balance from the balances API
	var x string
	data, err := GetRequest(ApiUrl + "/user/balance/asset?" + "username=" + username + "&pwhash=" + pwhash + "&asset=" + asset)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func GetStableCoin(username string, pwhash string, seed string, amount string) (rpc.StatusResponse, error) {
	var x rpc.StatusResponse
	data, err := GetRequest(ApiUrl + "/stablecoin/get?" + "seed=" + seed + "&amount=" +
		amount + "&username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func GetIpfsHash(username string, pwhash string, hashString string) (string, error) {
	var x string
	data, err := GetRequest(ApiUrl + "/ipfs/hash?" + "string=" + hashString +
		"&username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func InvestInProject(projIndex string, amount string, username string, pwhash string, seedpwd string) (rpc.StatusResponse, error) {
	var x rpc.StatusResponse
	data, err := GetRequest(ApiUrl + "/investor/invest?" + "username=" + username + "&pwhash=" + pwhash +
		"&seedpwd=" + seedpwd + "&projIndex=" + projIndex + "&amount=" + amount)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func VoteTowardsProject(projIndex string, amount string, username string, pwhash string) (rpc.StatusResponse, error) {
	var x rpc.StatusResponse
	data, err := GetRequest(ApiUrl + "/investor/vote?" + "username=" + username + "&pwhash=" + pwhash +
		"&projIndex=" + projIndex + "&votes=" + amount)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func AuthKyc(userIndex string, username string, pwhash string) (rpc.StatusResponse, error) {
	var x rpc.StatusResponse
	data, err := GetRequest(ApiUrl + "/user/kyc?" + "username=" + username + "&pwhash=" + pwhash +
		"&userIndex=" + userIndex)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func Payback(projIndex string, seedpwd string, username string, pwhash string, assetName string,
	amount string) (rpc.StatusResponse, error) {
	var x rpc.StatusResponse
	data, err := GetRequest(ApiUrl + "/recipient/payback?" + "username=" + username + "&pwhash=" + pwhash +
		"&projIndex=" + projIndex + "&seedpwd=" + seedpwd + "&amount=" + amount + "&assetName=" + assetName +
		"&platformPublicKey=" + PlatformPublicKey)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func UnlockProject(username string, pwhash string, seedpwd string, projIndex string) (rpc.StatusResponse, error) {
	var x rpc.StatusResponse
	data, err := GetRequest(ApiUrl + "/recipient/unlock?" + "username=" + username + "&pwhash=" + pwhash +
		"&projIndex=" + projIndex + "&seedpwd=" + seedpwd)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func FinalizeProject(username string, pwhash string, projIndex string) (rpc.StatusResponse, error) {
	var x rpc.StatusResponse
	data, err := GetRequest(ApiUrl + "/recipient/finalize?" + "username=" + username + "&pwhash=" + pwhash +
		"&projIndex=" + projIndex)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func OriginateProject(username string, pwhash string, projIndex string) (rpc.StatusResponse, error) {
	var x rpc.StatusResponse
	data, err := GetRequest(ApiUrl + "/recipient/originate?" + "username=" + username + "&pwhash=" + pwhash +
		"&projIndex=" + projIndex)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func GetOriginatedContracts(username string, pwhash string) ([]solar.Project, error) {
	var x []solar.Project
	data, err := GetRequest(ApiUrl + "/entity/getorigin?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func GetPreOriginatedContracts(username string, pwhash string) ([]solar.Project, error) {
	var x []solar.Project
	data, err := GetRequest(ApiUrl + "/entity/getpreorigin?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func GetProposedContracts(username string, pwhash string) ([]solar.Project, error) {
	var x []solar.Project
	data, err := GetRequest(ApiUrl + "/entity/getproposed?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func AddCollateral(username string, pwhash string, collateral string, amount string) (rpc.StatusResponse, error) {
	var x rpc.StatusResponse
	data, err := GetRequest(ApiUrl + "/entity/addcollateral?" + "username=" + username + "&pwhash=" + pwhash +
		"&collateral=" + collateral + "&amount=" + amount)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func CreateAssetInv(username string, pwhash string, assetName string, pubkey string) (rpc.StatusResponse, error) {
	var x rpc.StatusResponse
	data, err := GetRequest(ApiUrl + "/investor/localasset?" + "username=" + username + "&pwhash=" + pwhash +
		"&assetName=" + assetName)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func SendLocalAsset(username string, pwhash string, seedpwd string, assetName string,
	destination string, amount string) (string, error) {
	var x string

	data, err := GetRequest(ApiUrl + "/investor/sendlocalasset?" + "username=" + username + "&pwhash=" + pwhash +
		"&assetName=" + assetName + "&destination=" + destination + "&amount=" + amount + "&seedpwd=" + seedpwd)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func SendXLM(username string, pwhash string, seedpwd string, destination string,
	amount string, memo string) (string, error) {
	var x string
	data, err := GetRequest(ApiUrl + "/user/sendxlm?" + "username=" + username + "&pwhash=" + pwhash +
		"&destination=" + destination + "&amount=" + amount + "&seedpwd=" + seedpwd)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func NotKycView(username string, pwhash string) ([]database.User, error) {
	var x []database.User
	data, err := GetRequest(ApiUrl + "/user/notkycview?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func KycView(username string, pwhash string) ([]database.User, error) {
	var x []database.User
	data, err := GetRequest(ApiUrl + "/user/kycview?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func AskXLM(username string, pwhash string) (rpc.StatusResponse, error) {
	var x rpc.StatusResponse
	data, err := GetRequest(ApiUrl + "/user/askxlm?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func TrustAsset(username string, pwhash string, assetName string, issuerPubkey string,
	limit string, seedpwd string) (rpc.StatusResponse, error) {
	var x rpc.StatusResponse
	data, err := GetRequest(ApiUrl + "/user/trustasset?" + "username=" + username + "&pwhash=" + pwhash +
		"&assetCode=" + assetName + "&assetIssuer=" + issuerPubkey + "&limit=" + limit + "&seedpwd=" + seedpwd)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}
