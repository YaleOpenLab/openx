package aes

// the aes package implements AES-256 GCM encryption and decrpytion functions
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"

	utils "github.com/YaleOpenLab/openx/utils"
)

// Encrypt encrypts a given data stream with a given passphrase
func Encrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(utils.SHA3hash(passphrase)[96:128]) // last 32 characters in hash
	block, _ := aes.NewCipher(key)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return data, errors.Wrap(err, "Error while opening new GCM block")
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return data, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt decrypts a given data stream with a given passphrase
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	if len(data) == 0 || len(passphrase) == 0 {
		return data, errors.New("Length of data is zero, can't decrpyt!")
	}
	tempParam := utils.SHA3hash(passphrase)
	key := []byte(tempParam[96:128]) // last 32 characters in hash
	block, err := aes.NewCipher(key)
	if err != nil {
		return data, errors.Wrap(err, "Error while initializing new cipher")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return data, errors.Wrap(err, "failed to initialize new gcm block")
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return plaintext, errors.Wrap(err, "Error while opening gcm mode")
	}
	return plaintext, nil
}

// EncryptFile encrypts a given file with the given passphrase
func EncryptFile(filename string, data []byte, passphrase string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "Error while creating file")
	}
	defer f.Close()
	data, err = Encrypt(data, passphrase)
	if err != nil {
		return errors.Wrap(err, "Error while encrypting file")
	}
	f.Write(data)
	return nil
}

// DecryptFile encrypts a given file with the given passphrase
func DecryptFile(filename string, passphrase string) ([]byte, error) {
	data, _ := ioutil.ReadFile(filename)
	return Decrypt(data, passphrase)
}
