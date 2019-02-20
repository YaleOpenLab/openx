// +build all travis

package utils

import (
	"log"
	"os/user"
	"testing"
	"time"
)

func TestUtils(t *testing.T) {
	// test out stuff here
	testString := "10"
	testStringFloat := "10.000000"
	testInt := 10
	testByte := []byte("10")
	testFloat := 10.0
	var err error

	testSlice := ItoB(10)
	log.Println("hi", testString, testInt)
	// slcies can be compared only with nil
	for i, char := range testSlice {
		if char != testByte[i] {
			t.Fatalf("ItoB deosn't work as expected, quitting!")
		}
	}

	testSlice = []byte("01")
	check := false
	for i, char := range testSlice {
		if char == testByte[i] {
			check = true
		}
	}

	if check {
		t.Fatalf("Failed to catch error while comparing two different slcies")
	}

	if ItoS(testInt) != testString {
		t.Fatalf("ItoS deosn't work as expected, quitting!")
	}

	if ItoS(testInt) == "" {
		t.Fatalf("ItoS deosn't work as expected, quitting!")
	}

	if BToI(testByte) != testInt {
		t.Fatalf("BToI deosn't work as expected, quitting!")
	}

	if BToI(testByte) == 9 {
		t.Fatalf("BToI deosn't work as expected, quitting!")
	}

	if FtoS(testFloat) != testStringFloat {
		log.Println(FtoS(testFloat))
		t.Fatalf("FtoS deosn't work as expected, quitting!")
	}

	if FtoS(testFloat) == testString {
		log.Println(FtoS(testFloat))
		t.Fatalf("FtoS deosn't work as expected, quitting!")
	}

	if StoF(testStringFloat) != testFloat {
		log.Println(StoF(testStringFloat))
		t.Fatalf("StoF deosn't work as expected, quitting!")
	}

	if StoF(testStringFloat) == 9.0 {
		log.Println(StoF(testStringFloat))
		t.Fatalf("StoF deosn't work as expected, quitting!")
	}

	if StoI(testString) != testInt {
		t.Fatalf("StoI deosn't work as expected, quitting!")
	}

	if StoI(testString) == 9 {
		t.Fatalf("StoI deosn't work as expected, quitting!")
	}

	if SHA3hash("password") != "e9a75486736a550af4fea861e2378305c4a555a05094dee1dca2f68afea49cc3a50e8de6ea131ea521311f4d6fb054a146e8282f8e35ff2e6368c1a62e909716" {
		t.Fatalf("SHA3 doesn't work as expected, quitting!")
	}

	if SHA3hash("blah") == "e9a75486736a550af4fea861e2378305c4a555a05094dee1dca2f68afea49cc3a50e8de6ea131ea521311f4d6fb054a146e8282f8e35ff2e6368c1a62e909716" {
		t.Fatalf("SHA3 doesn't work as expected, quitting!")
	}

	hd, err := GetHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	usr, err := user.Current()
	if usr.HomeDir != hd || err != nil {
		t.Fatalf("Home directories don't match, quitting!")
	}

	rs := GetRandomString(10)
	if len(rs) != 10 {
		t.Fatalf("Random string length not equal to what is expected")
	}

	if time.Now().Format(time.RFC850) != Timestamp() {
		t.Fatalf("Timestamps don't match, quitting!")
	}

	if time.Now().Unix() != Unix() {
		t.Fatalf("Timestamps don't match, quitting!")
	}

	if I64toS(123412341234) != "123412341234" {
		t.Fatalf("I64 to string doesn't work, quitting!")
	}

	_, err = StoFWithCheck("blah")
	if err == nil {
		t.Fatalf("Not able to catch invalid string to float error!")
	}

	_, err = StoICheck("blah")
	if err == nil {
		t.Fatalf("Not able to catch invalid string to float error!")
	}
}
