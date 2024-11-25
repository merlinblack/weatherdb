package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/merlinblack/weatherdb/internal/repository/weather"
)

func HourlySummary(w http.ResponseWriter, r *http.Request, weather *weather.Queries) {

	hours := 24
	hoursParam := r.URL.Query().Get(`hours`)

	if len(hoursParam) > 0 {
		i, err := strconv.Atoi(hoursParam)
		if err != nil {
			log.Printf("Bad value for hours sent: %v", hoursParam)
		} else {
			hours = i
		}
	}

	log.Printf("[%s] [%s] Using hours = %v\n", r.Method, r.URL.Path, hours)

	measurements, err := weather.GetHourlySummary(context.Background(), int32(hours))
	if err != nil {
		internal500(w, `Could not get hourly summary records`, err)
		return
	}

	var rows []map[string]any

	for _, measurement := range measurements {
		row := make(map[string]any)

		row[`time`] = measurement.Hour
		row[`temperature`] = measurement.Temperature
		row[`humidity`] = measurement.Humidity
		row[`pressure`] = measurement.Pressure

		rows = append(rows, row)
	}

	jsonResp, err := json.Marshal(rows)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Header().Set(`Content-Type`, `application/json; charset=utf=8`)
	w.Write(jsonResp)

}
