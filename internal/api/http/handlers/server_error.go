package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func internal500(w http.ResponseWriter, errorMessage string, err error) {

	log.Printf("%s: %v\n", errorMessage, err)

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]string)

	resp["message"] = errorMessage
	resp["status"] = "500"

	jsonResp, err := json.Marshal(resp)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}
