package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	consts "github.com/YaleOpenLab/openx/consts"
)

// deviceid sets the deviceid and stores it in a retrievable location

func GenerateRandomString(n int) (string, error) {
	// generate a crypto secure random string
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func GenerateDeviceID() (string, error) {
	// this function is supposed to set the device id of this particular instance.
	// This should run only once on startup and save the device ID to some local
	// file so that we can reference it later
	rs, err := GenerateRandomString(16)
	if err != nil {
		return "", err
	}
	upperCase := strings.ToUpper(rs)
	return upperCase, nil
}

func CreateHomeDir() {
	if _, err := os.Stat(consts.TellerHomeDir); os.IsNotExist(err) {
		// directory does not exist, create one
		log.Println("Creating home directory for teller")
		os.MkdirAll(consts.TellerHomeDir, os.ModePerm)
	}
}

func CheckDeviceID() error {
	// checks whether there is a device id set on this device beforehand
	if _, err := os.Stat(consts.TellerHomeDir); os.IsNotExist(err) {
		// directory does not exist, create a device id
		log.Println("Creating home directory for teller")
		os.MkdirAll(consts.TellerHomeDir, os.ModePerm)
		path := consts.TellerHomeDir + "/deviceid.hex"
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		deviceId, err := GenerateDeviceID()
		if err != nil {
			return err
		}
		log.Println("GENERATED UNIQUE DEVICE ID: ", deviceId)
		_, err = file.Write([]byte(deviceId))
		if err != nil {
			return err
		}
		file.Close()
		err = SetDeviceId(LocalRecipient.U.Username, LocalRecipient.U.Pwhash, deviceId)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDeviceID() (string, error) {
	path := consts.TellerHomeDir + "/deviceid.hex"
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	// read the hex string from the file
	data := make([]byte, 32)
	numInt, err := file.Read(data)
	if err != nil {
		return "", err
	}
	if numInt != 32 {
		log.Println("NUMINT: ", numInt)
		return "", fmt.Errorf("Length of strings doesn't match, quitting!")
	}
	return string(data), nil
}
