package main

import (
	"log"
	"net/http"
	"os"

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

func StartServer(portx int) {
	openx()
	opensolar()
	teller()

	port, err := utils.ToString(portx)
	if err != nil {
		log.Fatal("Port not string")
	}

	log.Println("Starting RPC Server on Port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var opts struct {
	Port int `short:"p" description:"The port on which the server runs on" default:"8081"`
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}
	StartServer(opts.Port)
}
