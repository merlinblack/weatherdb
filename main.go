package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"net/http"
	"strconv"

	"github.com/carmo-evan/strtotime"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/merlinblack/weatherdb/weather_repository"
)

type Time struct {
	time.Time
}

const timeJSONLayout = `2006-01-02 15:04`

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(t.Format(timeJSONLayout))
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}
	t.Time, err = time.Parse(timeJSONLayout, s)
	return
}

type MeasurementJSON struct {
	RecordedAt  Time    `json:"recorded_at"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
}

func MeasurementToJSON(m *weather_repository.Measurement) string {
	jm := MeasurementJSON{}
	jm.RecordedAt.Time = m.RecordedAt
	jm.Temperature = m.Temperature
	jm.Humidity = m.Humidity
	jm.Pressure = m.Pressure

	jsonString, err := json.Marshal(jm)
	if err != nil {
		quitOnError("Problem marshalling json", err)
	}

	return string(jsonString)
}

func MeasurementFromJSON(data string) weather_repository.Measurement {
	jm := MeasurementJSON{}

	err := json.Unmarshal([]byte(data), &jm)
	if err != nil {
		quitOnError("Problem unmarshalling JSON", err)
	}

	m := weather_repository.Measurement{}
	m.RecordedAt = jm.RecordedAt.Time
	m.Temperature = jm.Temperature
	m.Humidity = jm.Humidity
	m.Pressure = jm.Pressure

	return m
}

func quitOnError(message string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Problem inserting row: %v\n", err)
	}

	measurements, err := weather.GetRecentWeather(context.Background(), 10)
	if err != nil {
		quitOnError(`Could not get recent weather records`, err)
	}

	first := true
	fmt.Print("[\n")
	for _, measurement := range measurements {
		if !first {
			fmt.Print(",\n")
		} else {
			first = false
		}
		fmt.Printf("  %v", MeasurementToJSON(&measurement))
	}
	fmt.Print("\n]\n")

	seconds, err := strtotime.Parse(`2 hour`, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem parsing duration: %v\n", err)
	} else {
		interval := time.Duration(seconds * int64(time.Second))
		trend, err := weather.GetWeatherTrend(context.Background(), interval)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem retrieving weather trends: %v\n", err)
		} else {
			fmt.Printf("Temperature: %s, Humidity: %s, Pressure: %s\n", trend.Temperature, trend.Humidity, trend.Pressure)
		}
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
			fmt.Fprintf(os.Stderr, "Bad value for limit sent: %v", limitParam)
		} else {
			limit = i
		}
	}

	fmt.Printf("Using limit = %v\n", limit )

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

	for _,period := range periods {
		seconds, err := strtotime.Parse(period, 0)
		if err != nil {
			quitOnError("Problem parsing duration: %v\n", err)
		} else {
			interval := time.Duration(seconds * int64(time.Second))
			trend, err := weather.GetWeatherTrend(context.Background(), interval)
			if err != nil {
				quitOnError("Problem retrieving weather trends: %v\n", err)
			}
			trends = append(trends, trend)
		}
	}

	first := true
	fmt.Fprintf(w, "{\n")
	for index,period := range periods {
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


func makeDBHandlerClosure( repo *weather_repository.Queries, fn func(w http.ResponseWriter, r *http.Request, repo *weather_repository.Queries)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, repo)
	}
}

func main() {

	dsn := getDsn()
	conn := getConnection(dsn)
	defer conn.Close()

	weather := weather_repository.New(conn)

	mux := http.NewServeMux()

	mux.HandleFunc(`GET /api/weather`, makeDBHandlerClosure(weather, recentMeasurements))
	mux.HandleFunc(`GET /api/trends`, makeDBHandlerClosure(weather, trends))

	fmt.Println(`Listening on localhost:3000`)

	http.ListenAndServe(`localhost:3000`, mux )
}
