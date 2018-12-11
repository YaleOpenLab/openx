package utils

import (
	//"log"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"

	"golang.org/x/crypto/sha3"
)

func Uint32toB(a uint32) []byte {
	// need to convert int to a byte array for indexing
	temp := make([]byte, 4)
	binary.LittleEndian.PutUint32(temp, a)
	return temp
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
