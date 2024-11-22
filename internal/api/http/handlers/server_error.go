package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func internalError(w http.ResponseWriter, statusCode int, errorMessage string) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := make(map[string]string)

	resp["message"] = errorMessage
	resp["status"] = strconv.Itoa(statusCode)

	jsonResp, err := json.Marshal(resp)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

func internal500(w http.ResponseWriter, errorMessage string, err error) {

	// Only log err, client does not see it.

	log.Printf("%s: %v\n", errorMessage, err)
	internalError(w, http.StatusInternalServerError, errorMessage)
}

func internal400(w http.ResponseWriter, errorMessage string, err error) {

	// Send both message and err (if given) to client
	if err != nil {
		errorMessage = fmt.Sprintf("%s: %v", errorMessage, err)
	}
	log.Println(errorMessage)
	internalError(w, http.StatusBadRequest, errorMessage)
}
