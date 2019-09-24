package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	flags "github.com/jessevdk/go-flags"
)

func openx() {
	http.HandleFunc("/openx", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		http.ServeFile(w, r, "openx.gz")
	})
}

func opensolar() {
	http.HandleFunc("/opensolar", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		http.ServeFile(w, r, "opensolar.gz")
	})
}

func teller() {
	http.HandleFunc("/teller", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		http.ServeFile(w, r, "teller.gz")
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
	Port int `short:"p" description:"The port on which the server runs on" default:"8081"`
	Insecure  bool `short:"i" description:"Start the API using http. Not recommended"`
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}
	StartServer(opts.Port, opts.Insecure)
}
