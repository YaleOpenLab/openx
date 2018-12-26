package utils

// utils contains utility functions that are needed commonly in packages
import (
	//"log"
	"bufio"
	"os"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"syscall"
	"time"

	clients "github.com/stellar/go/clients/horizon"
	"golang.org/x/crypto/sha3"
	"golang.org/x/crypto/ssh/terminal"
)

var DefaultTestNetClient = &clients.Client{
	URL:  "https://horizon-testnet.stellar.org",
	HTTP: http.DefaultClient,
}

func ScanForInt() (int, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return -1, fmt.Errorf("Couldn't read user input")
	}
	num := scanner.Text()
	numI, err := strconv.Atoi(num)
	if err != nil {
		return -1, fmt.Errorf("Input not a number")
	}
	return numI, nil
}

func ScanForString() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", fmt.Errorf("Couldn't read user input")
	}
	inputString := scanner.Text()
	return inputString, nil
}

func ScanForStringWithCheckI() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", fmt.Errorf("Couldn't read user input")
	}
	inputString := scanner.Text()
	_, err := strconv.Atoi(inputString) // check whether input string is a number (for payback)
	if err != nil {
		return "", err
	}
	return inputString, nil
}


func ScanForStringWithCheckF() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", fmt.Errorf("Couldn't read user input")
	}
	inputString := scanner.Text()
	if StoF(inputString) == 0 {
		fmt.Println("Amount entered is not a float, quitting")
		return "", fmt.Errorf("Amount entered is not a float, quitting")
	}
	return inputString, nil
}

func ScanForPassword() (string, error) {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	tempString := string(bytePassword)
	hashedPassword := SHA3hash(tempString)
	return hashedPassword, nil
}

func ScanRawPassword() (string, error) {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	password := string(bytePassword)
	return password, nil
}

func Timestamp() string {
	return time.Now().Format(time.RFC850)
}

func ItoB(a int) []byte {
	// need to convert int to a byte array for indexing
	string1 := strconv.Itoa(a)
	return []byte(string1)
}

func ItoS(a int) string {
	aStr := strconv.Itoa(a)
	return aStr
}

func BToI(a []byte) int {
	x, _ := strconv.Atoi(string(a))
	return x
}

func FtoS(a float64) string {
	return fmt.Sprintf("%f", a)
}

func StoF(a string) float64 {
	x, _ := strconv.ParseFloat(a, 32)
	// ignore this error since we hopefully call this in the right place
	return x
}

func StoI(a string) int {
	// convert string to int
	aInt, _ := strconv.Atoi(a)
	return aInt
}

func SHA3hash(inputString string) string {
	byteString := sha3.Sum512([]byte(inputString))
	return hex.EncodeToString(byteString[:])
	// so now we have a SHA3hash that we can use to assign unique ids to our assets
}
