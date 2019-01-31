package rpc

import (
	"fmt"
	"log"
	"net/http"

	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
)

func setupEntityRPCs() {
	validateEntity()
	getPreOriginatedContracts()
	getOriginatedContracts()
}

func EntityValidateHelper(w http.ResponseWriter, r *http.Request) (solar.Entity, error) {
	// first validate the investor or anyone would be able to set device ids
	checkGet(w, r)
	var prepInvestor solar.Entity
	// need to pass the pwhash param here
	if r.URL.Query() == nil || r.URL.Query()["username"] == nil ||
		len(r.URL.Query()["pwhash"][0]) != 128 {
		return prepInvestor, fmt.Errorf("Invalid params passed")
	}

	prepEntity, err := solar.ValidateEntity(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
	if err != nil {
		return prepEntity, err
	}

	return prepEntity, nil
}

func validateEntity() {
	http.HandleFunc("/entity/validate", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		log.Println("Prepared Entity:", prepEntity)
		MarshalSend(w, r, prepEntity)
	})
}

func getPreOriginatedContracts() {
	http.HandleFunc("/entity/getpreorigin", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		x, err := solar.RetrieveOriginatorProjects(solar.PreOriginProject, prepEntity.U.Index)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		MarshalSend(w, r, x)
	})
}

func getOriginatedContracts() {
	http.HandleFunc("/entity/getorigin", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		x, err := solar.RetrieveOriginatorProjects(solar.OriginProject, prepEntity.U.Index)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		MarshalSend(w, r, x)
	})
}
