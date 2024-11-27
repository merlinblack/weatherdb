package handlers

import (
	"context"
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
			log.Printf("[%s] [%s] Bad value for limit sent: %v", r.Method, r.URL.Path, limitParam)
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

		row[`recordedAt`] = formatTime(measurement.RecordedAt)
		row[`temperature`] = formatFloat(measurement.Temperature)
		row[`humidity`] = formatFloat(measurement.Humidity)
		row[`pressure`] = formatFloat(measurement.Pressure)

		rows = append(rows, row)
	}

	jsonResponse(w, http.StatusOK, rows)

}
