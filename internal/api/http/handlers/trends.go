package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/carmo-evan/strtotime"
	"github.com/merlinblack/weatherdb/internal/repository/weather"
)

func Trends(w http.ResponseWriter, r *http.Request, weatherdb *weather.Queries) {
	w.Header().Set(`Content-Type`, `application/json; charset=utf=8`)

	periods := []string{`15 minutes`, `1 hour`, `12 hours`, `1 week`, `1 month`}
	trends := make([]weather.Trend, 0, len(periods))

	for _, period := range periods {
		seconds, err := strtotime.Parse(period, 0)
		if err != nil {
			log.Fatalf("Problem parsing duration: %v\n", err)
		} else {
			interval := time.Duration(seconds * int64(time.Second))
			trend, err := weatherdb.GetTrends(context.Background(), interval)
			if err != nil {
				internal500(w, `Problem retrieving weather trends`, err)
				return
			}
			trends = append(trends, trend)
		}
	}

	first := true
	fmt.Fprintf(w, "{\n")
	for index, period := range periods {
		if !first {
			fmt.Fprintf(w, ",\n")
		} else {
			first = false
		}
		fmt.Fprintf(w, `"%v":{"temperature":"%v","humidity":"%v","pressure":"%v"}`,
			period,
			trends[index].Temperature,
			trends[index].Humidity,
			trends[index].Pressure,
		)
	}
	fmt.Fprintf(w, "\n}\n")
}
