package rpc

import (
	"fmt"
	//"log"
	"net/http"

	platform "github.com/YaleOpenLab/openx/platforms/opensolar"
	utils "github.com/YaleOpenLab/openx/utils"
)

func setupEntityRPCs() {
	validateEntity()
	getPreOriginatedContracts()
	getOriginatedContracts()
	getProposedContracts()
	addCollateral()
}

// EntityValidateHelper is a helper that helps validate an entity
func EntityValidateHelper(w http.ResponseWriter, r *http.Request) (platform.Entity, error) {
	// first validate the investor or anyone would be able to set device ids
	checkGet(w, r)
	var prepInvestor platform.Entity
	// need to pass the pwhash param here
	if r.URL.Query() == nil || r.URL.Query()["username"] == nil ||
		len(r.URL.Query()["pwhash"][0]) != 128 {
		return prepInvestor, fmt.Errorf("Invalid params passed")
	}

	prepEntity, err := platform.ValidateEntity(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
	if err != nil {
		return prepEntity, err
	}

	return prepEntity, nil
}

// validateEntity is an endpoint that vlaidates is a specific entity is registered on the platform
func validateEntity() {
	http.HandleFunc("/entity/validate", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}
		MarshalSend(w, r, prepEntity)
	})
}

// getPreOriginatedContracts gets a list of all the pre origianted contracts on the platform
func getPreOriginatedContracts() {
	http.HandleFunc("/entity/getpreorigin", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		x, err := platform.RetrieveOriginatorProjects(platform.PreOriginProject, prepEntity.U.Index)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		MarshalSend(w, r, x)
	})
}

// getOriginatedContracts gets a list of all the originated contracts on the platform
func getOriginatedContracts() {
	http.HandleFunc("/entity/getorigin", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		x, err := platform.RetrieveOriginatorProjects(platform.OriginProject, prepEntity.U.Index)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		MarshalSend(w, r, x)
	})
}

// getProposedContracts gets a list of all the proposed contracts on the platform
func getProposedContracts() {
	http.HandleFunc("/entity/getproposed", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		x, err := platform.RetrieveContractorProjects(platform.ProposedProject, prepEntity.U.Index)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}
		MarshalSend(w, r, x)
	})
}

// addCollateral is a route that can be used to add collateral to a specific contractor who wishes
// to propose a contract towards a specific originated project.
func addCollateral() {
	//func (contractor *Entity) AddCollateral(amount float64, data string) error {
	http.HandleFunc("/entity/addcollateral", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil || r.URL.Query()["amount"] == nil || r.URL.Query()["collateral"] == nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		collateralAmount, err := utils.StoFWithCheck(r.URL.Query()["amount"][0])
		if err != nil {
			responseHandler(w, r, StatusBadRequest)
			return
		}

		collateralData := r.URL.Query()["collateral"][0]
		err = prepEntity.AddCollateral(collateralAmount, collateralData)
		if err != nil {
			responseHandler(w, r, StatusInternalServerError)
			return
		}

		responseHandler(w, r, StatusOK)
	})
}
