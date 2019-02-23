// +build all

package aes

import (
	"os"
	"testing"
)

// Start benchmarking stuff here so that we can measuer how fast our functions run
func BenchmarkEncrypt(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = Encrypt([]byte("Hello World"), "Cool")
	}
}

// need to encrypt and decrypt here since we don't have the encrypted sipher text at hand
// and if we have a statice one, it would not be measuring anything
func BenchmarkDecrypt(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		ciphertext, _ := Encrypt([]byte("Hello World"), "cool")
		_, _ = Decrypt(ciphertext, "cool")
	}
}

func BenchmarkEncryptFile(b *testing.B) {
	os.MkdirAll("test_files", os.ModePerm)
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		EncryptFile("test_files/text.txt", []byte("Hello World"), "cool")
	}
	b.StopTimer()
	os.RemoveAll("test_files")
}

func BenchmarkEncryptDecryptFile(b *testing.B) {
	os.MkdirAll("test_files", os.ModePerm)
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		EncryptFile("test_files/text.txt", []byte("Hello World"), "cool")
		DecryptFile("test_files/text.txt", "cool")
	}
	b.StopTimer()
	os.RemoveAll("test_files")
}
