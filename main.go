package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
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

type Measurement struct {
	Id          int
	RecordedAt  time.Time
	Temperature float64
	Humidity    float64
	Pressure    float64
	Location    string
}

type MeasurementJSON struct {
	RecordedAt  Time    `json:"recorded_at"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
}

func (m *Measurement) String() string {
	return m.asJSON()
}

func (m *Measurement) asJSON() string {
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

func MeasurementFromJSON(data string) Measurement {
	jm := MeasurementJSON{}

	err := json.Unmarshal([]byte(data), &jm)
	if err != nil {
		quitOnError("Problem unmarshalling JSON", err)
	}

	m := Measurement{}
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
	user := "nigel"
	password := os.Getenv("WEATHERDB_PASSWORD")
	host := `octavo.local`
	port := 5432
	database := "weather_test"

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, database)
}

func getConnection(dsn string) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		quitOnError("Unable to create connection pool", err)
	}

	return db
}

func getRecentMeasurements(db *pgxpool.Pool, limit int) []*Measurement {
	var measurements []*Measurement

	err := pgxscan.Select(context.Background(), db, &measurements, `select * from measurements order by recorded_at desc limit $1`, limit)
	if err != nil {
		quitOnError("pgxscan failed", err)
	}

	return measurements
}

func insertMeasurement(db *pgxpool.Pool, measurement *Measurement) error {

	measurement.RecordedAt = measurement.RecordedAt.Round(time.Second)

	_, err := db.Exec(context.Background(),
		`insert into measurements (recorded_at, temperature, humidity, pressure, location) values ($1, $2, $3, $4, $5)`,
		measurement.RecordedAt,
		measurement.Temperature,
		measurement.Humidity,
		measurement.Pressure,
		measurement.Location)

	return err
}

func main() {
	dsn := getDsn()
	db := getConnection(dsn)
	defer db.Close()

	data := `{"recorded_at":"2024-04-02 23:24", "temperature": 24.3, "humidity": 67.32, "pressure": 1019.2}`

	new := MeasurementFromJSON(data)
	new.Location = `testing`

	if err := insertMeasurement(db, &new); err != nil {
		fmt.Fprintf(os.Stderr, "Problem inserting row: %v\n", err)
	}

	measurements := getRecentMeasurements(db, 10)

	first := true
	fmt.Print("[\n")
	for _, measurement := range measurements {
		if !first {
			fmt.Print(",\n")
		} else {
			first = false
		}
		fmt.Printf("  %v", measurement)
	}
	fmt.Print("\n]\n")

	db.Exec(context.Background(), "delete from measurements where location = 'testing'")
}
