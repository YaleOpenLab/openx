package main

import (
	//"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	//	"io/ioutil"
	"log"
	"math/big"
	// "net/http"
	// "strings"
)

func main() {

	//connect to rinkeby through infura
	ec, err := ethclient.Dial("https://ropsten.infura.io/")
	if err != nil {
		log.Fatal(err)
	}

	chainID := big.NewInt(3) //Ropsten

	//private key of sender
	privateKey, err := crypto.HexToECDSA("0000000000000000000000000000000000000000000000000")
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
	fmt.Printf("nonce: %i\n", nonce)
	fmt.Printf("amount: %i\n", amount)
	fmt.Printf("gasLimit: %s\n", gasLimit)
	fmt.Printf("gasPrice: %s\n", gasPrice)
	fmt.Printf("data: %s\n", data)

	//create raw transaction
	transaction := types.NewTransaction(nonce, recipient, amount, gasLimit, gasPrice, data)

	//sign transaction for rinkeby network
	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(transaction, signer, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(signedTx)
	//broadcast transaction
	err = ec.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
	//jsonData := fmt.Sprintf(` {"jsonrpc":"2.0", "method":"eth_sendRawTransaction", "params": ["0x%x"], "id":4}`, buff.Bytes())
	//params := buff.String()
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
