package rpc

import (
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
		checkOrigin(w, r)
		allBonds, err := opzones.RetrieveAllLivingUnitCoops()
		if err != nil {
			log.Println("did not retrieve all bonds", err)
			responseHandler(w, StatusInternalServerError)
			return
		}
		MarshalSend(w, allBonds)
	})
}

// getCoopDetails gets the details of a particular coop
func getCoopDetails() {
	http.HandleFunc("/coop/get", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		// get the details of a specific bond by key
		if r.URL.Query()["index"] == nil {
			responseHandler(w, StatusBadRequest)
			return
		}
		uKey := utils.StoI(r.URL.Query()["index"][0])
		bond, err := opzones.RetrieveLivingUnitCoop(uKey)
		if err != nil {
			log.Println("did not retrieve coop", err)
			responseHandler(w, StatusBadRequest)
			return
		}
		bondJson, err := bond.MarshalJSON()
		if err != nil {
			log.Println("did not marhsal json", err)
			responseHandler(w, StatusInternalServerError)
			return
		}
		WriteToHandler(w, bondJson)
	})
}

// getBondDetails gets the details of a particular bond
func getBondDetails() {
	http.HandleFunc("/bond/get", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		// get the details of a specific bond by key
		if r.URL.Query()["index"] == nil {
			responseHandler(w, StatusInternalServerError)
			return
		}
		uKey := utils.StoI(r.URL.Query()["index"][0])
		bond, err := opzones.RetrieveConstructionBond(uKey)
		if err != nil {
			log.Println("did not retrieve bond", err)
			responseHandler(w, StatusInternalServerError)
			return
		}
		MarshalSend(w, bond)
	})
}

// GetAllBonds gets the list of all bonds that are registered on the platfomr
func GetAllBonds() {
	http.HandleFunc("/bond/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		allBonds, err := opzones.RetrieveAllConstructionBonds()
		if err != nil {
			log.Println("did not retrieve all bonds", err)
			responseHandler(w, StatusInternalServerError)
			return
		}
		MarshalSend(w, allBonds)
	})
}

// Search can be used for searching bonds and coops to a limited capacity
func Search() {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		// search for coop / bond  and return accordingly
		if r.URL.Query()["q"] == nil {
			responseHandler(w, StatusBadRequest)
			return
		}
		searchString := r.URL.Query()["q"][0]
		if strings.Contains(searchString, "bond") {
			allBonds, err := opzones.RetrieveAllConstructionBonds()
			if err != nil {
				log.Println("did not retrieve all bonds", err)
				responseHandler(w, StatusInternalServerError)
				return
			}
			MarshalSend(w, allBonds)
			// do bond stuff
		} else if strings.Contains(searchString, "coop") {
			// do coop stuff
			allCoops, err := opzones.RetrieveAllLivingUnitCoops()
			if err != nil {
				log.Println("did not retrieve bond", err)
				responseHandler(w, StatusInternalServerError)
				return
			}
			MarshalSend(w, allCoops)
		} else {
			responseHandler(w, StatusInternalServerError)
			return
		}
	})
}
