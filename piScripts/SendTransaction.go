package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
	"os"
	"net/http"
	"strings"
	"io/ioutil"
)

type Message struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

// var (
// 	pin = rpio.Pin(15)
// )

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	var text string

	// //let user enter timestamp
	for text != "q" {
		fmt.Print("Hit Enter to Request (q to quit) ")
		scanner.Scan()
		text := scanner.Text()
		//fmt.Println(text)
		if text == "q" {
			return
		}
		//getProposedDeploymentDetails()
		//sendRawTransaction("89118DDA2B6C0F382D35905A766B0F46D8841FF6F0B7FAEA67545138C1E07940", "0x8abaf071687cbbd1b3dfbd6aa6c572a41e36d7ce", "makePayment(int256)", 100, text)
		sendEnergyData("B5858A3A04FEAA0D3C9953EADCB4D458D68ED85B7E5BC698F7208C0930D398D3", "0xa9f9fa6c881ef37865a85063ffa03cf4de992b16", "energyConsumed(uint256,uint256)", 10, "1", "10")
	}
}


func getProposedDeploymentDetails() string {
	transferFnSignature := []byte("getSolarSystemDetails(uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	argumentAmount := new(big.Int)
	argumentAmount.SetString("1", 10) //
	paddedAmount := common.LeftPadBytes(argumentAmount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAmount...)

	//fmt.Printf("data: %x\n", data)

	//user jsonrpc to do ethcall
	jsonData := fmt.Sprintf(` {"jsonrpc":"2.0", "method":"eth_call", "params": [{"from": "0x717d97A81e9aFF8748B23859eca81A4fE26d8165", "to": "0xa9f9fa6c881ef37865a85063ffa03cf4de992b16","gas": "0x7530", "data": "0x%x"}, "latest"], "id":3}`, data)
	//params := buff.String()
	//fmt.Printf("%s\n", jsonData)
	response, err := http.Post("https://ropsten.infura.io/gnNuNKvHFmjf9xkJ0StE", "application/json", strings.NewReader(jsonData))
	if err != nil {
		log.Fatal("Infura request failed", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		fmt.Println("INFURA response:")
		fmt.Println(string(data))
		return string(data)
		//c <- string(data)
	}
	return ""

}



//Sign and send a raw transaction using the private key of an account
//The function retrieves the necessary
//Args (all strings): private key, recipient address, method name, argument amount
//ex call: sendRawTransaction("f6c649c0e891b19df822730a0d773a7a54cc4e5dcaebe1a8543591f211e05cb5", "0x86a64d840ab2665c137335af9c354f3d57c189d9", "setPower(uint8)", "2")
//Does not return any data

//value: eth that is sent as payment
//argAmount: number that is passed to be used in contract handling (this function handles one argument)
func sendRawTransaction(_privateKey string, recipientAddress string, methodName string, value int64, argAmount string) {
	//connect to ropsten through infura
	ec, err := ethclient.Dial("https://ropsten.infura.io/")
	if err != nil {
		log.Fatal(err)
	}

	chainID := big.NewInt(3) //Ropsten

	//private key of sender
	//TODO: hide key when actual system is implemented
	privateKey, err := crypto.HexToECDSA(_privateKey)
	if err != nil {
		log.Fatal(err)
	}

	//get Public Key of sender
	publicKey := privateKey.Public()
	publicKey_ECDSA, valid := publicKey.(*ecdsa.PublicKey)
	if !valid {
		log.Fatal("error casting public key to ECDSA")
	}

	//get address of sender
	fromAddress := crypto.PubkeyToAddress(*publicKey_ECDSA)

	//get nonce of address
	nonce, err := ec.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	//get recipient address
	recipient := common.HexToAddress(recipientAddress)

	amount := big.NewInt(value) // 0 ether
	gasLimit := uint64(2000000)
	gasPrice, err := ec.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	transferFnSignature := []byte(methodName)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	//fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	argumentAmount := new(big.Int)
	argumentAmount.SetString(argAmount, 10) //
	paddedAmount := common.LeftPadBytes(argumentAmount.Bytes(), 32)
	//fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAmount...)
	//data := []byte("0x5c22b6b60000000000000000000000000000000000000000000000000000000000000007")
	// fmt.Printf("nonce: %i\n", nonce)
	// fmt.Printf("amount: %i\n", amount)
	// fmt.Printf("gasLimit: %s\n", gasLimit)
	// fmt.Printf("gasPrice: %s\n", gasPrice)
	fmt.Printf("data: %x\n", data)

	//create raw transaction
	transaction := types.NewTransaction(nonce, recipient, amount, gasLimit, gasPrice, data)

	//sign transaction for ropsten network
	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(transaction, signer, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// var buff bytes.Buffer
	// signedTx.EncodeRLP(&buff)
	// fmt.Printf("0x%x\n", buff.Bytes())

	//fmt.Println(signedTx)
	//broadcast transaction
	err = ec.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())

	// jsonData := fmt.Sprintf(` {"jsonrpc":"2.0", "method":"eth_sendRawTransaction", "params": ["0x%x"], "id":4}`, buff.Bytes())
	// //params := buff.String()
	// fmt.Printf("%s\n", jsonData)
	// response, err := http.Post("https://rinkeby.infura.io/gnNuNKvHFmjf9xkJ0StE", "application/json", strings.NewReader(jsonData))
	// if err != nil {

	// 	fmt.Printf("Request to INFURA failed with an error: %s\n", err)
	// 	fmt.Println()

	// } else {
	// 	data, _ := ioutil.ReadAll(response.Body)

	// 	fmt.Println("INFURA response:")
	// 	fmt.Println(string(data))
	// }
}


//value: eth that is sent as payment
//arg1: _ssAddress
//arg2: _energyConsumed
func sendEnergyData(_privateKey string, recipientAddress string, methodName string, value int64, arg1 string, arg2 string) {
	//connect to ropsten through infura
	ec, err := ethclient.Dial("https://ropsten.infura.io/")
	if err != nil {
		log.Fatal(err)
	}

	chainID := big.NewInt(3) //Ropsten

	//private key of sender
	//TODO: hide key when actual system is implemented
	privateKey, err := crypto.HexToECDSA(_privateKey)
	if err != nil {
		log.Fatal(err)
	}

	//get Public Key of sender
	publicKey := privateKey.Public()
	publicKey_ECDSA, valid := publicKey.(*ecdsa.PublicKey)
	if !valid {
		log.Fatal("error casting public key to ECDSA")
	}

	//get address of sender
	fromAddress := crypto.PubkeyToAddress(*publicKey_ECDSA)

	//get nonce of address
	nonce, err := ec.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	//get recipient address
	recipient := common.HexToAddress(recipientAddress)

	amount := big.NewInt(value) // 0 ether
	gasLimit := uint64(2000000)
	gasPrice, err := ec.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	transferFnSignature := []byte(methodName)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	//fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	argumentAmount1 := new(big.Int)
	argumentAmount1.SetString(arg1, 10) //
	paddedAmount1 := common.LeftPadBytes(argumentAmount1.Bytes(), 32)

	argumentAmount2 := new(big.Int)
	argumentAmount2.SetString(arg2, 10) //
	paddedAmount2 := common.LeftPadBytes(argumentAmount2.Bytes(), 32)
	//fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAmount1...)
	data = append(data, paddedAmount2...)
	//data := []byte("0x5c22b6b60000000000000000000000000000000000000000000000000000000000000007")
	// fmt.Printf("nonce: %i\n", nonce)
	// fmt.Printf("amount: %i\n", amount)
	// fmt.Printf("gasLimit: %s\n", gasLimit)
	// fmt.Printf("gasPrice: %s\n", gasPrice)
	fmt.Printf("data: %x\n", data)

	//create raw transaction
	transaction := types.NewTransaction(nonce, recipient, amount, gasLimit, gasPrice, data)

	//sign transaction for ropsten network
	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(transaction, signer, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// var buff bytes.Buffer
	// signedTx.EncodeRLP(&buff)
	// fmt.Printf("0x%x\n", buff.Bytes())

	//fmt.Println(signedTx)
	//broadcast transaction
	err = ec.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())

	// jsonData := fmt.Sprintf(` {"jsonrpc":"2.0", "method":"eth_sendRawTransaction", "params": ["0x%x"], "id":4}`, buff.Bytes())
	// //params := buff.String()
	// fmt.Printf("%s\n", jsonData)
	// response, err := http.Post("https://rinkeby.infura.io/gnNuNKvHFmjf9xkJ0StE", "application/json", strings.NewReader(jsonData))
	// if err != nil {

	// 	fmt.Printf("Request to INFURA failed with an error: %s\n", err)
	// 	fmt.Println()

	// } else {
	// 	data, _ := ioutil.ReadAll(response.Body)

	// 	fmt.Println("INFURA response:")
	// 	fmt.Println(string(data))
	// }
}
