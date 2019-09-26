package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	flags "github.com/jessevdk/go-flags"
	btcutils "github.com/bithyve/research/utils"
)

var (
	LastBuilt string
	GithubSecret string
	Sha Shastruct
	OpenxHashes []string
	OpensolarHashes []string
	TellerHashes []string
)

func openx() {
	http.HandleFunc("/openx-darwinamd64", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "openx-darwinamd64.gz")
	})
	http.HandleFunc("/openx-linuxamd64", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "openx-linuxamd64.gz")
	})
	http.HandleFunc("/openx-linux386", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "openx-linux386.gz")
	})
	http.HandleFunc("/openx-arm64", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "openx-arm64.gz")
	})
	http.HandleFunc("/openx-arm", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "openx-arm.gz")
	})
}

func opensolar() {
	http.HandleFunc("/opensolar-darwinamd64", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "opensolar-darwinamd64.gz")
	})
	http.HandleFunc("/opensolar-linuxamd64", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "opensolar-linuxamd64.gz")
	})
	http.HandleFunc("/opensolar-linux386", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "opensolar-linux386.gz")
	})
	http.HandleFunc("/opensolar-arm64", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "opensolar-arm64.gz")
	})
	http.HandleFunc("/opensolar-arm", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "opensolar-arm.gz")
	})
}

func teller() {
	http.HandleFunc("/teller-darwinamd64", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "teller-darwinamd64.gz")
	})
	http.HandleFunc("/teller-linuxamd64", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "teller-linuxamd64.gz")
	})
	http.HandleFunc("/teller-linux386", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "teller-linux386.gz")
	})
	http.HandleFunc("/teller-arm64", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "teller-arm64.gz")
	})
	http.HandleFunc("/teller-arm", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		http.ServeFile(w, r, "teller-arm.gz")
	})
}


func lastbuilt() {
	http.HandleFunc("/lastbuilt", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		erpc.MarshalSend(w, LastBuilt)
	})
}

func shaEndpoint() {
	http.HandleFunc("/sha", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		erpc.MarshalSend(w, Sha)
	})
}

type HashesResponse struct {
	Openx []string
	Opensolar []string
	Teller []string
}

func hashesEndpoint() {
	http.HandleFunc("/hashes", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
		}

		var x HashesResponse
		x.Openx = OpenxHashes
		x.Opensolar = OpensolarHashes
		x.Teller = TellerHashes

		erpc.MarshalSend(w, x)
	})
}


// FileSystem custom file system handler
type FileSystem struct {
	fs http.FileSystem
}

// Open opens file
func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func frontend() {
	fileServer := http.FileServer(FileSystem{http.Dir("static")})
	http.Handle("/fe", http.StripPrefix(strings.TrimRight("/fe/", "/"), fileServer))
}

func StartServer(portx int, insecure bool) {
	erpc.SetupBasicHandlers()
	openx()
	opensolar()
	teller()
	frontend()
	lastbuilt()
	shaEndpoint()

	port, err := utils.ToString(portx)
	if err != nil {
		log.Fatal("Port not string")
	}

	log.Println("Starting RPC Server on Port: ", port)
	if insecure {
		log.Println("starting server in insecure mode")
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(":"+port, "certs/server.crt", "certs/server.key", nil))
	}
}

var opts struct {
	Port     int  `short:"p" description:"The port on which the server runs on" default:"8081"`
	Insecure bool `short:"i" description:"Start the API using http. Not recommended"`
}

func writeLastBuilt() {
	data := time.Now().String()
	err := ioutil.WriteFile("lastbuilt.txt", []byte(data), 0644)
	if err != nil {
		// don't return error
		log.Println(err)
	}
}

func readLastBuilt() {
	data, err := ioutil.ReadFile("lastbuilt.txt")
	if err != nil {
		// don't return error
		log.Println(err)
	}
	LastBuilt = string(data)
}

func readGhSecret() {
	data, err := ioutil.ReadFile("secret.txt")
	if err != nil {
		// don't return error
		log.Println(err)
	}
	GithubSecret = string(data)[0:40] // splice off the \n at the end
}

// GetRequest is a handler that makes it easy to send out GET requests
func GetRequest(url string) ([]byte, error) {
	var dummy []byte
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dummy, errors.Wrap(err, "did not create new GET request")
	}
	req.Header.Set("Origin", "localhost")
	req.SetBasicAuth("Varunram", GithubSecret)

	res, err := client.Do(req)
	if err != nil {
		return dummy, errors.Wrap(err, "did not make request")
	}

	defer func() {
		if ferr := res.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	return ioutil.ReadAll(res.Body)
}

type GithubSha struct {
	Sha string `json:"sha"`
}

type Shastruct struct {
	OpenxSha string
	OpensolarSha string
}

func updateShastruct() {
	var gh GithubSha

	data, err := GetRequest("https://api.github.com/repos/YaleOpenLab/openx/commits/master")
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(data, &gh)
	if err != nil {
		log.Println(err)
		return
	}

	Sha.OpenxSha = string(gh.Sha)

	data, err = GetRequest("https://api.github.com/repos/YaleOpenLab/opensolar/commits/master")
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(data, &gh)
	if err != nil {
		log.Println(err)
		return
	}
	Sha.OpensolarSha = string(gh.Sha)
	log.Println(Sha)
}

func updateShaHashes() {
	var openxFileNames = []string{"openx-darwinamd64.gz", "openx-linuxamd64.gz", "openx-linux386.gz", "openx-arm64.gz", "openx-arm.gz"}
	var opensolarFileNames = []string{"opensolar-darwinamd64.gz", "opensolar-linuxamd64.gz", "opensolar-linux386.gz", "opensolar-arm64.gz", "opensolar-arm.gz"}
	var tellerFileNames = []string{"teller-darwinamd64.gz", "teller-linuxamd64.gz", "teller-linux386.gz", "teller-arm64.gz", "teller-arm.gz"}

	for _, file := range openxFileNames {
		sha2Bytes, err := btcutils.Sha256File(file)
		if err != nil {
			log.Println(err)
		}
		OpenxHashes = append(OpenxHashes, hex.EncodeToString(sha2Bytes))
	}

	for _, file := range opensolarFileNames {
		sha2Bytes, err := btcutils.Sha256File(file)
		if err != nil {
			log.Println(err)
		}
		OpensolarHashes = append(OpensolarHashes, hex.EncodeToString(sha2Bytes))
	}

	for _, file := range tellerFileNames {
		sha2Bytes, err := btcutils.Sha256File(file)
		if err != nil {
			log.Println(err)
		}
		TellerHashes = append(TellerHashes, hex.EncodeToString(sha2Bytes))
	}
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	updateShaHashes()
	readGhSecret()
	readLastBuilt()
	hashesEndpoint()

	go func() {
		for {
			time.Sleep(24 * time.Hour)
			log.Println("triggering build script")
			_, err := exec.Command("./build.sh").Output()
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("build built succesfully")
			writeLastBuilt()
			updateShastruct()
		}
	}()

	StartServer(opts.Port, opts.Insecure)
}
