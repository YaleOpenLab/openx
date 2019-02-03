package scan

// package scan is not in utils since we can't test the below functions (which require
// user interaction) whereas functions in utils are essential
// and need to be tested in order for stuff to run properly

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"syscall"

	utils "github.com/OpenFinancing/openfinancing/utils"
	"golang.org/x/crypto/ssh/terminal"
)

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

func ScanForFloat() (float64, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return -1, fmt.Errorf("Couldn't read user input")
	}
	num := scanner.Text()
	x, err := strconv.ParseFloat(num, 32)
	// ignore this error since we hopefully call this in the right place
	return x, err
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
	if utils.StoF(inputString) == 0 {
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
	hashedPassword := utils.SHA3hash(tempString)
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