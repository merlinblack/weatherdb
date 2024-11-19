package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/carmo-evan/strtotime"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/merlinblack/weatherdb/internal/api/http/middleware"
	"github.com/merlinblack/weatherdb/internal/weather_repository"
)

func quitOnError(message string, err error) {
	log.Fatalf("%s: %v\n", message, err)
}

func getDsn() string {
	user := `nigel`
	password := os.Getenv(`WEATHERDB_PASSWORD`)
	host := `octavo.local`
	port := 5432
	database := `weather_test`

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, database)
}

func getConnection(dsn string) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		quitOnError(`Unable to create connection pool`, err)
	}

	return db
}

func test() {
	dsn := getDsn()
	conn := getConnection(dsn)
	defer conn.Close()

	weather := weather_repository.New(conn)

	data := `{"recorded_at":"2024-04-02 23:24", "temperature": 24.3, "humidity": 67.32, "pressure": 1019.2}`

	new := MeasurementFromJSON(data)
	new.Location = `testing`

	_, err := weather.InsertMeasurement(context.Background(),
		weather_repository.InsertMeasurementParams{
			RecordedAt:  new.RecordedAt,
			Temperature: new.Temperature,
			Humidity:    new.Humidity,
			Pressure:    new.Pressure,
			Location:    new.Location,
		})
	if err != nil {
		log.Printf("Problem inserting row: %v\n", err)
	}

	conn.Exec(context.Background(), `delete from measurements where location = 'testing'`)
}

func recentMeasurements(w http.ResponseWriter, r *http.Request, weather *weather_repository.Queries) {
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
	for _, measurement := range measurements {
		if !first {
			fmt.Fprintf(w, ",\n")
		} else {
			first = false
		}
		fmt.Fprintf(w, "  %v", MeasurementToJSON(&measurement))
	}
	fmt.Fprintf(w, "\n]\n")
}

func trends(w http.ResponseWriter, r *http.Request, weather *weather_repository.Queries) {
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

func makeHandlerWithRepo(repo *weather_repository.Queries, fn func(w http.ResponseWriter, r *http.Request, repo *weather_repository.Queries)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, repo)
	})
}

func main() {

	dsn := getDsn()
	conn := getConnection(dsn)
	defer conn.Close()

	weather := weather_repository.New(conn)

	mux := http.NewServeMux()

	mux.Handle(`GET /weather`, makeHandlerWithRepo(weather, recentMeasurements))
	mux.Handle(`GET /trends`, makeHandlerWithRepo(weather, trends))
	mux.HandleFunc(`GET /ping`, func(w http.ResponseWriter, _ *http.Request) { fmt.Fprintln(w, `pong`) })

	chain := middleware.ChainFinal(mux)
	chain.Use(middleware.LoggingMiddleware)

	server := &http.Server{
		Addr:    `:3000`,
		Handler: chain,
	}

	go func() {
		log.Printf("Listening on %v\n", server.Addr)
		err := server.ListenAndServe()

		if err == http.ErrServerClosed {
			log.Println(`Server closed, no longer accepting connections`)
		} else {
			quitOnError(`Problem starting http server`, err)
		}
	}()

	// Wait for OS interrupt (pkill -2 weatherdb)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	// put a newline after the possible '^C' that is now displayed if the user pressed ^C
	fmt.Println(``)

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	log.Println(`Graceful shutdown, current requests have 30 seconds to finish`)
	if err := server.Shutdown(ctx); err != nil {
		quitOnError("Problem shutting down", err)
	}

	log.Println(`Bye!`)
}
