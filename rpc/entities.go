package rpc

import (
	"log"
	"net/http"

	database "github.com/OpenFinancing/openfinancing/database"
	solar "github.com/OpenFinancing/openfinancing/platforms/solar"
)

func setupEntityRPCs() {
	validateEntity()
}

func validateEntity() {
	http.HandleFunc("/entity/validate", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		if r.URL.Query() == nil || r.URL.Query()["username"] == nil || r.URL.Query()["pwhash"] == nil ||
			len(r.URL.Query()["pwhash"][0]) != 128 { // sha 512 length
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		prepUser, err := database.ValidateUser(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		// we now have the user, retreive the entity
		prepEntity, err := solar.RetrieveEntity(prepUser.Index)
		if err != nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		log.Println("Prepared Entity:", prepEntity)
		MarshalSend(w, r, prepEntity)
	})
}
