package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/carmo-evan/strtotime"
	"github.com/merlinblack/weatherdb/internal/repository/weather"
)

func Trends(w http.ResponseWriter, r *http.Request, weatherdb *weather.Queries) {

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

	rows := make(map[string]any)

	for index, period := range periods {
		row := make(map[string]string)

		row[`temperature`] = trends[index].Temperature
		row[`humidity`] = trends[index].Humidity
		row[`pressure`] = trends[index].Pressure

		rows[period] = row
	}

	jsonResponse(w, http.StatusOK, rows)

}
