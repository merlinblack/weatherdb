package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/merlinblack/weatherdb/internal/measurement"
	"github.com/merlinblack/weatherdb/internal/weather_repository"
)

func RecentMeasurements(w http.ResponseWriter, r *http.Request, weather *weather_repository.Queries) {
	w.Header().Set(`Content-Type`, `application/json; charset=utf=8`)

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

	log.Printf("[%s] [%s] Using limit = %v\n", r.Method, r.URL, limit)

	measurements, err := weather.GetRecentWeather(context.Background(), int32(limit))
	if err != nil {
		quitOnError(`Could not get recent weather records`, err)
	}

	first := true
	fmt.Fprintf(w, "[\n")
	for _, measure := range measurements {
		if !first {
			fmt.Fprintf(w, ",\n")
		} else {
			first = false
		}
		fmt.Fprintf(w, "  %v", measurement.ToJSON(&measure))
	}
	fmt.Fprintf(w, "\n]\n")
}
