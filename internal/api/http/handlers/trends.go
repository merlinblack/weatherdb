package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/carmo-evan/strtotime"
	"github.com/merlinblack/weatherdb/internal/weather_repository"
)

func Trends(w http.ResponseWriter, r *http.Request, weather *weather_repository.Queries) {
	w.Header().Set(`Content-Type`, `application/json; charset=utf=8`)

	periods := []string{`15 minutes`, `1 hour`, `12 hours`, `1 week`, `1 month`}
	trends := make([]weather_repository.Trend, 0, len(periods))

	for _, period := range periods {
		seconds, err := strtotime.Parse(period, 0)
		if err != nil {
			quitOnError(`Problem parsing duration`, err)
		} else {
			interval := time.Duration(seconds * int64(time.Second))
			trend, err := weather.GetWeatherTrend(context.Background(), interval)
			if err != nil {
				quitOnError(`Problem retrieving weather trends`, err)
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
