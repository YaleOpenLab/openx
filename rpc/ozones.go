package rpc

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	opzones "github.com/YaleOpenLab/openx/platforms/ozones"
	utils "github.com/YaleOpenLab/openx/utils"
)

func setupCoopRPCs() {
	getCoopDetails()
	GetAllCoops()
}

func setupBondRPCs() {
	getBondDetails()
	Search()
	GetAllBonds()
}

// GetAllCoops gets a list of all the coops  that are registered on the platform
func GetAllCoops() {
	http.HandleFunc("/coop/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		allBonds, err := opzones.RetrieveAllLivingUnitCoops()
		if err != nil {
			log.Println("did not retrieve all bonds", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		MarshalSend(w, r, allBonds)
	})
}

// getCoopDetails gets the details of a particular coop
func getCoopDetails() {
	http.HandleFunc("/coop/get", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		// get the details of a specific bond by key
		if r.URL.Query()["index"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		uKey := utils.StoI(r.URL.Query()["index"][0])
		bond, err := opzones.RetrieveLivingUnitCoop(uKey)
		if err != nil {
			log.Println("did not retrieve coop", err)
			responseHandler(w, r, StatusBadRequest)
			return
		}
		bondJson, err := json.Marshal(bond)
		if err != nil {
			log.Println("did not marhsal json", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		WriteToHandler(w, bondJson)
	})
}

// getBondDetails gets the details of a particular bond
func getBondDetails() {
	http.HandleFunc("/bond/get", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		// get the details of a specific bond by key
		if r.URL.Query()["index"] == nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		uKey := utils.StoI(r.URL.Query()["index"][0])
		bond, err := opzones.RetrieveConstructionBond(uKey)
		if err != nil {
			log.Println("did not retrieve bond", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		MarshalSend(w, r, bond)
	})
}

// GetAllBonds gets the list of all bonds that are registered on the platfomr
func GetAllBonds() {
	http.HandleFunc("/bond/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		allBonds, err := opzones.RetrieveAllConstructionBonds()
		if err != nil {
			log.Println("did not retrieve all bonds", err)
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		MarshalSend(w, r, allBonds)
	})
}

// Search can be used for searching bonds and coops to a limited capacity
func Search() {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		// search for coop / bond  and return accordingly
		if r.URL.Query()["q"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		searchString := r.URL.Query()["q"][0]
		if strings.Contains(searchString, "bond") {
			allBonds, err := opzones.RetrieveAllConstructionBonds()
			if err != nil {
				log.Println("did not retrieve all bonds", err)
				responseHandler(w, r, StatusInternalServerError)
				return
			}
			MarshalSend(w, r, allBonds)
			// do bond stuff
		} else if strings.Contains(searchString, "coop") {
			// do coop stuff
			allCoops, err := opzones.RetrieveAllLivingUnitCoops()
			if err != nil {
				log.Println("did not retrieve bond", err)
				responseHandler(w, r, StatusInternalServerError)
				return
			}
			MarshalSend(w, r, allCoops)
		} else {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
	})
}
