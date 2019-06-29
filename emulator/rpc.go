package main

import (
	"encoding/json"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	database "github.com/YaleOpenLab/openx/database"
	solar "github.com/YaleOpenLab/openx/platforms/opensolar"
	"github.com/stellar/go/protocols/horizon"
)

func PingRpc() error {
	// make a curl request out to lcoalhost and get the ping response
	data, err := erpc.GetRequest(ApiUrl + "/ping")
	if err != nil {
		return err
	}
	var x erpc.StatusResponse
	// now data is in byte, we need the other structure now
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	// the result would be the status of the platform
	ColorOutput("PLATFORM STATUS: "+utils.ItoS(x.Code), GreenColor)
	return nil
}

func RetrieveProject(index int) ([]solar.Project, error) {
	// retrieve project at a particular stage
	var x []solar.Project
	data, err := erpc.GetRequest(ApiUrl + "/projects?index=" + utils.ItoS(index))
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func GetBalances(username string, pwhash string) ([]horizon.Balance, error) {
	// get the balance from the balances API
	var x []horizon.Balance
	data, err := erpc.GetRequest(ApiUrl + "/user/balances?" + "username=" + username + "&pwhash=" + pwhash)
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
	data, err := erpc.GetRequest(ApiUrl + "/user/balance/xlm?" + "username=" + username + "&pwhash=" + pwhash)
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
	data, err := erpc.GetRequest(ApiUrl + "/user/balance/asset?" + "username=" + username + "&pwhash=" + pwhash + "&asset=" + asset)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func GetStableCoin(username string, pwhash string, amount string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	data, err := erpc.GetRequest(ApiUrl + "/stablecoin/get?" + "seedpwd=" + LocalSeedPwd + "&amount=" +
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
	data, err := erpc.GetRequest(ApiUrl + "/ipfs/hash?" + "string=" + hashString +
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

func InvestInProject(projIndex string, amount string, username string, pwhash string, seedpwd string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	data, err := erpc.GetRequest(ApiUrl + "/investor/invest?" + "username=" + username + "&pwhash=" + pwhash +
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

func VoteTowardsProject(projIndex string, amount string, username string, pwhash string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	data, err := erpc.GetRequest(ApiUrl + "/investor/vote?" + "username=" + username + "&pwhash=" + pwhash +
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

func AuthKyc(userIndex string, username string, pwhash string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	data, err := erpc.GetRequest(ApiUrl + "/user/kyc?" + "username=" + username + "&pwhash=" + pwhash +
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
	amount string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	data, err := erpc.GetRequest(ApiUrl + "/recipient/payback?" + "username=" + username + "&pwhash=" + pwhash +
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

func UnlockOpenSolar(username string, pwhash string, seedpwd string, projIndex string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	body := ApiUrl + "/recipient/unlock/opensolar?" + "username=" + username + "&pwhash=" + pwhash +
		"&projIndex=" + projIndex + "&seedpwd=" + seedpwd

	data, err := erpc.GetRequest(body)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func FinalizeProject(username string, pwhash string, projIndex string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	data, err := erpc.GetRequest(ApiUrl + "/recipient/finalize?" + "username=" + username + "&pwhash=" + pwhash +
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

func OriginateProject(username string, pwhash string, projIndex string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	data, err := erpc.GetRequest(ApiUrl + "/recipient/originate?" + "username=" + username + "&pwhash=" + pwhash +
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

func GetStage1Contracts(username string, pwhash string) ([]solar.Project, error) {
	var x []solar.Project
	data, err := erpc.GetRequest(ApiUrl + "/entity/stage1?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func GetStage0Contracts(username string, pwhash string) ([]solar.Project, error) {
	var x []solar.Project
	data, err := erpc.GetRequest(ApiUrl + "/entity/stage0?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func GetStage2Contracts(username string, pwhash string) ([]solar.Project, error) {
	var x []solar.Project
	data, err := erpc.GetRequest(ApiUrl + "/entity/stage2?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func AddCollateral(username string, pwhash string, collateral string, amount string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	data, err := erpc.GetRequest(ApiUrl + "/entity/addcollateral?" + "username=" + username + "&pwhash=" + pwhash +
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

func CreateAssetInv(username string, pwhash string, assetName string, pubkey string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	data, err := erpc.GetRequest(ApiUrl + "/investor/localasset?" + "username=" + username + "&pwhash=" + pwhash +
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

	data, err := erpc.GetRequest(ApiUrl + "/investor/sendlocalasset?" + "username=" + username + "&pwhash=" + pwhash +
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
	data, err := erpc.GetRequest(ApiUrl + "/user/sendxlm?" + "username=" + username + "&pwhash=" + pwhash +
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
	data, err := erpc.GetRequest(ApiUrl + "/user/notkycview?" + "username=" + username + "&pwhash=" + pwhash)
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
	data, err := erpc.GetRequest(ApiUrl + "/user/kycview?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func AskXLM(username string, pwhash string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	data, err := erpc.GetRequest(ApiUrl + "/user/askxlm?" + "username=" + username + "&pwhash=" + pwhash)
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
	limit string, seedpwd string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	data, err := erpc.GetRequest(ApiUrl + "/user/trustasset?" + "username=" + username + "&pwhash=" + pwhash +
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

func GetTrustLimit(username string, pwhash string, assetName string) (string, error) {
	var x string
	data, err := erpc.GetRequest(ApiUrl + "/recipient/trustlimit?" + "username=" + username + "&pwhash=" +
		pwhash + "&assetName=" + assetName)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func InvestInOpzoneCBond(projIndex string, amount string, username string, pwhash string, seedpwd string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	data, err := erpc.GetRequest(ApiUrl + "/constructionbond/invest?" + "username=" + username + "&pwhash=" + pwhash +
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

func InvestInLivingUnitCoop(projIndex string, amount string, username string, pwhash string, seedpwd string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	data, err := erpc.GetRequest(ApiUrl + "/livingunitcoop/invest?" + "username=" + username + "&pwhash=" + pwhash +
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

func UnlockCBond(username string, pwhash string, seedpwd string, projIndex string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	body := ApiUrl + "/recipient/unlock/opzones/cbond?" + "username=" + username + "&pwhash=" + pwhash +
		"&projIndex=" + projIndex + "&seedpwd=" + seedpwd

	data, err := erpc.GetRequest(body)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func IncreaseTrustLimit(username string, pwhash string, seedpwd string, trust string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	body := ApiUrl + "/user/increasetrustlimit?" + "username=" + username + "&pwhash=" + pwhash +
		"&seedpwd=" + seedpwd + "&trust=" + trust

	data, err := erpc.GetRequest(body)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func SendSharesEmail(username string, pwhash string, email1 string, email2 string, email3 string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	body := ApiUrl + "/user/sendrecovery?" + "username=" + username + "&pwhash=" + pwhash +
		"&email1=" + email1 + "&email2=" + email2 + "&email3=" + email3

	data, err := erpc.GetRequest(body)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}

func SendNewSharesEmail(username string, pwhash string, seedpwd string, email1 string, email2 string, email3 string) (erpc.StatusResponse, error) {
	var x erpc.StatusResponse
	body := ApiUrl + "/user/newsecrets?" + "username=" + username + "&pwhash=" + pwhash +
		"&seedpwd=" + seedpwd + "&email1=" + email1 + "&email2=" + email2 + "&email3=" + email3

	data, err := erpc.GetRequest(body)
	if err != nil {
		return x, err
	}
	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}
	return x, nil
}
