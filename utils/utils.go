package utils

// utils contains utility functions that are needed commonly in packages
import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os/user"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)

func Timestamp() string {
	return time.Now().Format(time.RFC850)
}

func Unix() int64 {
	return time.Now().Unix()
}

func I64toS(a int64) string {
	return strconv.FormatInt(a, 10) // s == "97" (decimal)
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

func GetHomeDir() (string, error) {
	var homedir string
	usr, err := user.Current()
	if err != nil {
		return homedir, err
	}
	homedir = usr.HomeDir
	return homedir, nil
}

func GetRandomString(n int) string {
	// random string implementation courtesy: icza
	// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
