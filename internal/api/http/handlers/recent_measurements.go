package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/merlinblack/weatherdb/internal/repository/weather"
)

func RecentMeasurements(w http.ResponseWriter, r *http.Request, weather *weather.Queries) {

	limit := 10
	limitParam := r.URL.Query().Get(`limit`)

	if len(limitParam) > 0 {
		i, err := strconv.Atoi(limitParam)
		if err != nil {
			log.Printf("Bad value for limit sent: %v", limitParam)
		} else {
			limit = i
		}
	}

	log.Printf("[%s] [%s] Using limit = %v\n", r.Method, r.URL.Path, limit)

	measurements, err := weather.GetRecentMeasurements(context.Background(), int32(limit))
	if err != nil {
		internal500(w, `Could not get recent weather records`, err)
		return
	}

	var rows []map[string]any

	for _, measurement := range measurements {
		row := make(map[string]any)

		row[`recordedAt`] = measurement.RecordedAt.Format(timeJSONLayout)
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
