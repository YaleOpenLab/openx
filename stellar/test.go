package main

import (
	"log"

	accounts "github.com/Varunram/smartPropertyMVP/stellar/accounts"
)

func main() {
	var issuer accounts.Account
	var recipient accounts.Account
	err := (&issuer).New()
	if err != nil {
		log.Fatal(err)
	}

	err = (&recipient).New()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("My stellar public key and private key are: ", issuer.PublicKey, " ", issuer.Seed)

	err = issuer.GetCoins() // get coins for issuer
	if err != nil {
		log.Fatal(err)
	}

	err = recipient.GetCoins() // get coins for recipient
	if err != nil {
		log.Fatal(err)
	}

	// err = issuer.Balance()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	confHeight, txHash, err := issuer.SendCoins(recipient.PublicKey, "3.33") // send some coins from the issuer to the recipient
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Confirmation height is: ", confHeight, " and txHash is: ", txHash)

	asset := issuer.CreateAsset("PetroDlr") // create the asset that we want

	trustLimit := "100" // trust only 100 barrels of oil from Petro
	err = recipient.TrustAsset(asset, trustLimit)
	if err != nil {
		log.Println("Trust limit is in the wrong format")
		log.Fatal(err)
	}

	err = issuer.SendAsset("PetroDlr", recipient.PublicKey, "3.35")
	if err != nil {
		log.Fatal(err)
	}
}
