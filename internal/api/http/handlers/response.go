package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func formatTime(t time.Time) string {
	return t.Format(timeJSONLayout)
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 1, 64)
}

func jsonResponse(w http.ResponseWriter, statusCode int, data any) {

	jsonResp, err := json.Marshal(data)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Header().Set(`Content-Type`, `application/json; charset=utf-8`)
	w.WriteHeader(statusCode)
	w.Write(jsonResp)
}

func internalError(w http.ResponseWriter, statusCode int, errorMessage string) {

	resp := make(map[string]string)

	resp["message"] = errorMessage
	resp["status"] = strconv.Itoa(statusCode)

	jsonResponse(w, statusCode, resp)
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
