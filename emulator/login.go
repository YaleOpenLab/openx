package main

import (
	"encoding/json"
	"fmt"
	"log"

	database "github.com/OpenFinancing/openfinancing/database"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
	rpc "github.com/OpenFinancing/openfinancing/rpc"
	scan "github.com/OpenFinancing/openfinancing/scan"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
)

func Login(username string, pwhash string) (string, error) {
	var wString string
	data, err := rpc.GetRequest(ApiUrl + "/user/validate?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return wString, err
	}
	var x rpc.ValidateParams
	err = json.Unmarshal(data, &x)
	if err != nil {
		return wString, err
	}
	switch x.Role {
	case "Investor":
		wString = "Investor"
		data, err = rpc.GetRequest(ApiUrl + "/investor/validate?" + "username=" + username + "&pwhash=" + pwhash)
		if err != nil {
			return wString, err
		}
		var inv database.Investor
		err = json.Unmarshal(data, &inv)
		if err != nil {
			return wString, err
		}
		LocalInvestor = inv
		ColorOutput("ENTER YOUR SEEDPWD: ", CyanColor)
		LocalSeedPwd, err = scan.ScanRawPassword()
		if err != nil {
			log.Println(err)
			return wString, err
		}
		LocalSeed, err = wallet.DecryptSeed(LocalInvestor.U.EncryptedSeed, LocalSeedPwd)
		if err != nil {
			log.Println(err)
			return wString, err
		}
	case "Recipient":
		wString = "Recipient"
		data, err = rpc.GetRequest(ApiUrl + "/recipient/validate?" + "username=" + username + "&pwhash=" + pwhash)
		if err != nil {
			return wString, err
		}
		var recp database.Recipient
		err = json.Unmarshal(data, &recp)
		if err != nil {
			return wString, err
		}
		LocalRecipient = recp
		ColorOutput("ENTER YOUR SEEDPWD: ", CyanColor)
		LocalSeedPwd, err = scan.ScanRawPassword()
		if err != nil {
			log.Println(err)
			return wString, err
		}
		LocalSeed, err = wallet.DecryptSeed(LocalRecipient.U.EncryptedSeed, LocalSeedPwd)
		if err != nil {
			log.Println(err)
			return wString, err
		}
	case "Entity":
		data, err = rpc.GetRequest(ApiUrl + "/entity/validate?" + "username=" + username + "&pwhash=" + pwhash)
		if err != nil {
			return wString, err
		}
		var entity solar.Entity
		err = json.Unmarshal(data, &entity)
		if err != nil {
			return wString, err
		}
		if entity.Contractor {
			LocalContractor = entity
			wString = "Contractor"
		} else if entity.Originator {
			LocalOriginator = entity
			wString = "Originator"
		} else {
			return wString, fmt.Errorf("Not a contractor")
		}
		ColorOutput("ENTER YOUR SEEDPWD: ", CyanColor)
		LocalSeedPwd, err = scan.ScanRawPassword()
		if err != nil {
			log.Println(err)
			return wString, err
		}
		if entity.Contractor {
			LocalSeed, err = wallet.DecryptSeed(LocalContractor.U.EncryptedSeed, LocalSeedPwd)
			if err != nil {
				log.Println("error while decrpyting seed, quitting!", err)
				return wString, err
			}
		} else if entity.Originator {
			LocalSeed, err = wallet.DecryptSeed(LocalOriginator.U.EncryptedSeed, LocalSeedPwd)
			if err != nil {
				log.Println("error while decrpyting seed, quitting!", err)
				return wString, err
			}
		}
	}

	ColorOutput("AUTHENTICATED USER, YOUR ROLE IS: "+wString, GreenColor)
	return wString, nil
}
