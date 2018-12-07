package main

import (
	"log"

	accounts "github.com/Varunram/smartPropertyMVP/stellar/accounts"
)

func main() {
	var acc accounts.Account
	err := (&acc).New()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("My stellar public key and private key are: ", acc.PublicKey, " ", acc.Seed)

	err = acc.GetCoins()
	if err != nil {
		log.Fatal(err)
	}

	err = acc.Balance()
	if err != nil {
		log.Fatal(err)
	}

	destination := "GD4H5KBX6OL5VUBZDOC4DMCZIKGGUB4ZI3TUY4IPXJ4DLOD6BVT6GYSV" // a random address for now
	amount := "3.33"                                                          // weirdly enough, this is a string instead of uint or something

	confHeight, txHash, err := acc.SendCoins(destination, amount)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Confirmation height is: ", confHeight, " and txHash is: ", txHash)

	err = acc.CreateAsset("GD4H5KBX6OL5VUBZDOC4DMCZIKGGUB4ZI3TUY4IPXJ4DLOD6BVT6GYSV")
	if err != nil {
		log.Fatal(err)
	}
}
