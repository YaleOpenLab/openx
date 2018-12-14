package utils

// utils contains utility functions that are needed commonly in packages
import (
	//"log"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)


func StoI(a string) int {
	temp, _ := strconv.Atoi(a)
	// ignore error sicne we assume that we'll call this in the right place
	return temp
}

func Timestamp() string {
	return time.Now().Format(time.RFC850)
}

func Uint32toB(a uint32) []byte {
	// need to convert int to a byte array for indexing
	temp := make([]byte, 4)
	binary.LittleEndian.PutUint32(temp, a)
	return temp
}

func Uint32toS(a uint32) string {
	// convert uint32 to int and then int to string
	aInt := int(a)
	aStr := strconv.Itoa(aInt)
	return aStr
}

func StoUint32(a string) uint32 {
	// convert string to int
	aInt, _ := strconv.Atoi(a)
	return uint32(aInt)
}

func BToUint32(a []byte) uint32 {
	return binary.LittleEndian.Uint32(a)
}

func FloatToString(a float64) string {
	return fmt.Sprintf("%f", a)
}

func StringToFloat(a string) float64 {
	x, _ := strconv.ParseFloat(a, 32)
	// ignore this error since we hopefully call this in the right place
	return x
}

func IntToString(a int) string {
	return strconv.Itoa(a)
}

func SHA3hash(inputString string) string {
	byteString := sha3.Sum512([]byte(inputString))
	return hex.EncodeToString(byteString[:])
	// so now we have a SHA3hash that we can use to assign unique ids to our assets
}
