package main

import (
	//"bytes"
	"bufio"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stianeikeland/go-rpio"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

var (
	pin = rpio.Pin(15)
)

func main() {

	//intiailize Pin settings
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()
	pin.Output()

	scanner := bufio.NewScanner(os.Stdin)
	var text string

	//let user enter timestamp
	for text != "q" {
		fmt.Print("Hit Enter to Request (q to quit) ")
		scanner.Scan()
		text := scanner.Text()
		res := checkIfLiquid(text)
		fmt.Println(res)

		//decode json
		m := Message{}
		err := json.Unmarshal([]byte(res), &m)

		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println(m.Result)
		//parse return data (currently pulling last 10 digits of hex string)
		last10 := m.Result[len(m.Result)-10:]
		hexString := new(big.Int)
		hexString.SetString(last10, 16)
		fmt.Println(hexString.Int64())
		if hexString.Int64() == 1 {
			fmt.Println("liquid")
			pin.Low()
		} else {
			fmt.Println("notliquid")
			pin.High()
		}
	}
}

func getBalance() {
	jsonData := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x3789d2e02925d9b91f1c00a5c8e5bac33e50c786", "latest"],"id":4}`)
	response, err := http.Post("https://ropsten.infura.io/gnNuNKvHFmjf9xkJ0StE", "application/json", strings.NewReader(jsonData))

	if err != nil {

		fmt.Printf("Request to INFURA failed with an error: %s\n", err)
		fmt.Println()

	} else {
		data, _ := ioutil.ReadAll(response.Body)

		fmt.Println("INFURA response:")
		fmt.Println(string(data))
	}
}

func checkIfLiquid(date string) string {
	transferFnSignature := []byte("stillLiquid(int256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	//fmt.Println(hexutil.Encode(methodID)) //

	argumentAmount := new(big.Int)
	argumentAmount.SetString(date, 10) //
	paddedAmount := common.LeftPadBytes(argumentAmount.Bytes(), 32)
	//fmt.Println(hexutil.Encode(paddedAmount))

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAmount...)
	//fmt.Printf("data: %x\n", data)

	jsonData := fmt.Sprintf(` {"jsonrpc":"2.0", "method":"eth_call", "params": [{"from": "0x6ca9a0f319ec632fc21d4a16998f750923a50b32", "to": "0x3789d2e02925d9b91f1c00a5c8e5bac33e50c786","gas": "0x7530", "data": "0x%x"}, "latest"], "id":3}`, data)
	//params := buff.String()
	//fmt.Printf("%s\n", jsonData)
	response, err := http.Post("https://ropsten.infura.io/gnNuNKvHFmjf9xkJ0StE", "application/json", strings.NewReader(jsonData))
	if err != nil {

		fmt.Printf("Request to INFURA failed with an error: %s\n", err)
		fmt.Println()

	} else {
		data, _ := ioutil.ReadAll(response.Body)

		fmt.Println("INFURA response:")
		//fmt.Println(string(data))
		return string(data)
	}
	return ""

}

func sendRawTransaction() {
	//connect to rinkeby through infura
	ec, err := ethclient.Dial("https://ropsten.infura.io/")
	if err != nil {
		log.Fatal(err)
	}

	chainID := big.NewInt(3) //Ropsten

	//private key of sender
	privateKey, err := crypto.HexToECDSA("f6c649c0e891b19df822730a0d773a7a54cc4e5dcaebe1a8543591f211e05cb5")
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
	recipient := common.HexToAddress("0x86a64d840ab2665c137335af9c354f3d57c189d9")

	amount := big.NewInt(0) // 0 ether
	gasLimit := uint64(2000000)
	gasPrice, err := ec.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	transferFnSignature := []byte("setPower(uint8)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	argumentAmount := new(big.Int)
	argumentAmount.SetString("2", 10) //
	paddedAmount := common.LeftPadBytes(argumentAmount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	//TODO: format data to accept inputs from various functions
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAmount...)
	//data := []byte("0x5c22b6b60000000000000000000000000000000000000000000000000000000000000007")
	// fmt.Printf("nonce: %i\n", nonce)
	// fmt.Printf("amount: %i\n", amount)
	// fmt.Printf("gasLimit: %s\n", gasLimit)
	// fmt.Printf("gasPrice: %s\n", gasPrice)
	// fmt.Printf("data: %s\n", data)

	//create raw transaction
	transaction := types.NewTransaction(nonce, recipient, amount, gasLimit, gasPrice, data)

	//sign transaction for rinkeby network
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
