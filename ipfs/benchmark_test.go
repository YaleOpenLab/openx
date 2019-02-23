// +build all

package ipfs

import (
	"testing"
)

func BenchmarkAddString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AddStringToIpfs("Hello, this is a test from ipfs to see if it works")
	}
}

func BenchmarkAddAndRetrieveString(b *testing.B) {
	hash, err := AddStringToIpfs("Hello, this is a test from ipfs to see if it works")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetStringFromIpfs(hash)
	}
}

func BenchmarkReadFromFile(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ReadfromFile("files/test.pdf") // get the data from the pdf as a datastream
	}
}

func BenchmarkHashFile(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = IpfsHashFile("files/test.pdf")
	}
}

func BenchmarkHashBytestring(b *testing.B) {
	dummy := []byte("Hello World")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = IpfsHashData(dummy)
	}
}
