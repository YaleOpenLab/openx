package ipfs

// this package contains the ipfs interacting parts
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

	consts "github.com/OpenFinancing/openfinancing/consts"
	utils "github.com/OpenFinancing/openfinancing/utils"
	shell "github.com/ipfs/go-ipfs-api"
)

// RetrieveShell retrieves the ipfs shell for use by other functions
// the path must be set to the rpc port used by the local / remote host
func RetrieveShell() *shell.Shell {
	// this is the api endpoint of the ipfs daemon
	return shell.NewShell("localhost:5001")
}

// AddStringToIpfs stores the given s tring in ipfs and returns
// the hash of the string
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

// GetFileFromIpfs gets back the contents of an ipfs hash and stores them
// in the required extension format. This has to match with the extension
// format that the original file had or else one would not be able to view
// the file
func GetFileFromIpfs(hash string, extension string) error {
	// extension can be pdf, txt, ppt and others
	sh := RetrieveShell()
	// generate a random fileName and then return the file to the user
	fileName := utils.GetRandomString(consts.IpfsFileLength) + "." + extension
	log.Println("DECRYPTED FILE NAME is: ", fileName)
	return sh.Get(hash, fileName)
}

// GetStringFromIpfs gets back the contents of an ipfs hash as a string
func GetStringFromIpfs(hash string) (string, error) {
	sh := RetrieveShell()
	// since ipfs doesn't provide a method to read the string directly, we create a
	// random fiel at tmp/, decrypt contents to that fiel and then read the file
	// contents from there
	tmpFileDir := "/tmp/" + utils.GetRandomString(consts.IpfsFileLength) // using the same length here for consistency
	sh.Get(hash, tmpFileDir)
	data, err := ioutil.ReadFile(tmpFileDir)
	if err != nil {
		return "", err
	}
	os.Remove(tmpFileDir)
	return string(data), nil
}

// ReadfromPdf reads a pdf and returns the datastream
func ReadfromPdf(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

// IpfsHashPdf returns the ipfs hash of a pdf file
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
