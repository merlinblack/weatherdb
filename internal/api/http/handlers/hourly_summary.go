package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/merlinblack/weatherdb/internal/repository/weather"
)

func HourlySummary(w http.ResponseWriter, r *http.Request, weather *weather.Queries) {
	w.Header().Set(`Content-Type`, `application/json; charset=utf=8`)

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

	log.Printf("[%s] [%s] Using hours = %v\n", r.Method, r.URL, hours)

	measurements, err := weather.GetHourlySummary(context.Background(), int32(hours))
	if err != nil {
		internal500(w, `Could not get hourly summary records`, err)
		return
	}

	first := true
	fmt.Fprintf(w, "[\n")
	for _, measure := range measurements {
		if !first {
			fmt.Fprintf(w, ",\n")
		} else {
			first = false
		}
		fmt.Fprintf(w, `  {"time":"%s", "temperature":"%s", "humidity":"%s", "pressure":"%s" }`,
			measure.Hour,
			measure.Temperature,
			measure.Humidity,
			measure.Pressure,
		)
	}
	fmt.Fprintf(w, "\n]\n")
}
