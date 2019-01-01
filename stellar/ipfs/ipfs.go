package ipfs

// this package contains the ipfs interacting parts
// to install, ipfs, download a release from https://github.com/ipfs/go-ipfs/releases
// and then run the install.sh script there.
// In case you face an issue with migration, you might need to run fs-repo-migrations
// from the ipfs github in order to migrate your existing ipfs files to a newer
// version. If you don't have anything valuable there, you can go ahead and delete
// the directory and then run ipfs init again. You need to keep your peer key
// stored in a safe place for future reference
// Then you need to start ipfs using ipfs daemon and then you can test if it worked
// by creating a test file test.txt and then doing ipfs add test.txt. The resultant
// hash van be decrypted using curl "http://127.0.0.1:8080/ipfs/hash" where 8080
// is the endpoint of the ipfs server or by doing cat /ipfs/hash directly

// when we are adding a file to ipfs, we either could use the javascript handler
// to call the ipfs api and then use the hash ourselves to decrypt it. Or we need to
// process a pdf file (ie build an xref table) and then convert that into an ipfs file
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	shell "github.com/ipfs/go-ipfs-api"
)

func RetrieveShell() *shell.Shell {
	// this is the api endpoint of the ipfs daemon
	return shell.NewShell("localhost:5001")
}

func AddStringToIpfs(a string) (string, error) {
	sh := RetrieveShell()
	hash, err := sh.Add(strings.NewReader(a)) // input must be an io.Reader
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Println("Added text: ", a, " to ipfs, hash: ", hash)
	return hash, nil
}

func GetFileFromIpfs(hash string, extension string) error {
	// extension can be pdf, txt, ppt and others
	sh := RetrieveShell()
	// generate a random fileName and then return the file to the user
	fileName := utils.GetRandomString(10) + "." + extension
	log.Println("DECRYPTED FILE NAME is: ", fileName)
	return sh.Get(hash, fileName)
}

func GetStringFromIpfs(hash string) (string, error) {
	sh := RetrieveShell()
	// since ipfs doesn't provide a method to read the string directly, we create a
	// random fiel at tmp/, decrypt contents to that fiel and then read the file
	// contents from there
	tmpFileDir := "/tmp/" + utils.GetRandomString(10)
	sh.Get(hash, tmpFileDir)
	data, err := ioutil.ReadFile(tmpFileDir)
	if err != nil {
		return "", err
	}
	os.Remove(tmpFileDir)
	return string(data), nil
}

// Read from pdf reads the pdf and returns the datastream
func ReadfromPdf(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

func IpfsHashPdf(filepath string) (string, error) {
	var dummy string
	dataStream, err := ReadfromPdf(filepath)
	if err != nil {
		return dummy, err
	}
	// need to get the ifps hash of this data stream and return hash
	reader := bytes.NewReader(dataStream)
	sh := RetrieveShell()
	hash, err := sh.Add(reader)
	if err != nil {
		return dummy, err
	}
	return hash, nil
}

/*
useful for writing tests in the future:
hash, err := ipfs.AddStringToIpfs("Hello, this is a test from ipfs to see if it works")
if err != nil {
	log.Fatal(err)
}
log.Println("HASH: ", hash)
string1, err := ipfs.GetStringFromIpfs(hash)
if err != nil {
	log.Fatal(err)
}
err = ipfs.GetFileFromIpfs("/ipfs/QmSjvpAbHtAkFNV7SRmsV5pRALwsr8waEoWL4NesCvUdpw", "pdf")
if err != nil {
	log.Fatal(err)
}
log.Println("DECRYPTED STRING IS: ", string1)
_, err = ipfs.ReadfromPdf("test.pdf") // get the data from the pdf as a datastream
if err != nil {
	log.Fatal(err)
}

hash, err := ipfs.IpfsHashPdf("test.pdf")
if err != nil {
	log.Fatal(err)
}
log.Println("HASH IS: ", hash)
err = ipfs.GetFileFromIpfs(hash, "pdf")
if err != nil {
	log.Fatal(err)
}
log.Fatal("TEST")
*/
