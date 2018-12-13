package rpc

// the rpc package contains functions related to the server which will be interacting
// with the frontend. Not expanding on this too much since this will be changing quite often
import (
	"fmt"
	"log"
	// "net"
	// "io/ioutil"
	"encoding/json"
	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	"net/http"
	// "net/rpc"
)

type Dummy struct {
	Username string
	Password string
}

type PingResponse struct {
	Status string
}

func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Header().Add("Access-Control-Allow-Origin", "localhost")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonString)
}

func checkOrigin(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Origin") != "localhost" {
		// allow only our frontend UI to connect to our RPC instance
		http.Error(w, "404 page not found", http.StatusNotFound)
		return fmt.Errorf("Cross domain request error")
	}
	return nil
}

func checkGet(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		// reject wrong method entries
		http.Error(w, "404 page not found", http.StatusNotFound)
		return fmt.Errorf("Invalid Method error")
	}
	return nil
}

func checkPost(w http.ResponseWriter, r *http.Request) error {
	log.Println("CHEcking POST")
	if r.Method != "POST" {
		// reject wrong method entries
		http.Error(w, "404 page not found", http.StatusNotFound)
		return fmt.Errorf("Invalid Method error")
	}
	return nil
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 Page Not Found")
	}
}

func setupDefaultHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// default to 404 for every application not running on localhost
		errorHandler(w, r, http.StatusNotFound)
		return
	})
}

func setupPingHandler() {

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		err := checkOrigin(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		err = checkGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		var pr PingResponse
		pr.Status = "Alive"
		prJson, err := json.Marshal(pr)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, prJson)
	})
}

func setupTestHandler() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		var dummy Dummy
		log.Println("remote peer calling /cool endpoint")
		dummy.Username = "dummy"
		dummy.Password = "123"
		dummyJson, err := json.Marshal(dummy)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, dummyJson)
	})
}

func getOpenOrders() {
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		// make a call to the db to get all orders
		// while making this call, the rpc should not be aware of the db we are using
		// and stuff. So we need to have another route that would open the existing
		// db, without asking for one
		allOrders, err := database.RetrieveAllWithoutDB()
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		ordersJson, err := json.Marshal(allOrders)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, ordersJson)
	})
}

func getOrder() {
	// we need to read passed the key from the URL that the user calls
	http.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkGet(w, r)
		URLPath := r.URL.Path
		// so we now have the URL path
		// slice "/order/" off from the URLPath
		keyS := URLPath[7:]
		// we now need to get the order corresponding to keyS
		// the rpc accepts the key as uint32 though, so string -> uint32

		uKey := utils.StoUint32(keyS)
		order, err := database.RetrieveOrderRPC(uKey)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		orderJson, err := json.Marshal(order)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		WriteToHandler(w, orderJson)
	})
}

func parseOrder(r *http.Request) (database.Order, error) {
	// we need to create an instance of the Order
	// and then map values if they do exist
	// note that we just prepare the order here and don't invest in it
	// for that, we need new a new investor struct and a recipient struct
	// Index         uint32
	// PanelSize     string  required
	// TotalValue    int     required
	// Location      string  required
	// MoneyRaised   int     default to 0
	// Metadata      string  optional
	// Live          bool    will be set to true if this is called
	// INVAssetCode  string  should be set by subsequent calls
	// DEBAssetCode  string  should be set by subsequent calls
	// PBAssetCode   string  should be set by subsequent calls
	// BalLeft       float64 should be equal to  totalValue sicne htis is a new order
	// DateInitiated string  auto
	// DateLastPaid  string  don't set

	var prepOrder database.Order
	err := r.ParseForm()
	if err != nil {
		return prepOrder, err
	}
	// if we're inserting this in, we need to get the next index number
	// so that we can set this without causing some weird bugs
	allOrders, err := database.RetrieveAllWithoutDB()
	if err != nil {
		return prepOrder, err
	}
	prepOrder.Index = uint32(len(allOrders) + 1)
	if r.FormValue("PanelSize") != "" {
		prepOrder.PanelSize = r.FormValue("PanelSize")
	} else {
		// if this is not defined, error out
		return prepOrder, err
	}

	if r.FormValue("TotalValue") != "" {
		// the totlaValue passed here is a string, we need to convert to an int
		totalValueS := utils.StoI(r.FormValue("TotalValue"))
		prepOrder.TotalValue = totalValueS
	} else {
		return prepOrder, err
	}

	if r.FormValue("Location") != "" {
		prepOrder.Location = r.FormValue("Location")
	} else {
		return prepOrder, err
	}

	prepOrder.MoneyRaised = 0

	if r.FormValue("Metadata") != "" {
		prepOrder.Metadata = r.FormValue("Metadata")
	}

	prepOrder.Live = true
	// set the codes later while setting up stuff, need rpc calls for those as well
	prepOrder.BalLeft = float64(0)
	prepOrder.DateInitiated = utils.Timestamp()
	return prepOrder, nil
}

func insertOrder() {
	// this should be a post method since you want to accetp an order and then insert
	// that into the database
	http.HandleFunc("/insert", func(w http.ResponseWriter, r *http.Request) {
		checkOrigin(w, r)
		checkPost(w, r)
		var prepOrder database.Order
		prepOrder, err := parseOrder(r)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		log.Println("Prepared Order:", prepOrder)
	})
}

func StartServer() {
	// this runs on the server side ie the server with the frontend.
	// having to define specific endpoints for this because this
	// is the system that would be used by the backend, so has to be secure.

	// the idea is that we have a unique handler to each of these routes, which will
	// then return the appropriate data to be used by the frotnend
	// we also have to process data from the pi itself, and that should have its own
	// functions somewhere else that can be accessed by the rpc.

	// setup a couple test handlers that we can remove later
	setupDefaultHandler()
	setupPingHandler()
	setupTestHandler()
	// lets setup an endpoint that would retrieve all orders
	getOpenOrders()
	// so the problem with this route is that we need to set this up for each
	// key that exists in our database, is hte re some other nice way to do that?
	// or handle this in a better way?
	getOrder()
	insertOrder()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
