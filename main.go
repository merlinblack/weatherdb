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

type MeasurementDB struct {
	Id          int `json:"-"`
	RecordedAt  time.Time
	Temperature float64
	Humidity    float64
	Pressure    float64
	Location    string `json:"-"`
}

func (m *MeasurementDB) String() string {
	return fmt.Sprintf("%d - %f, %f, %f %s %s", m.Id, m.Temperature, m.Humidity, m.Pressure, m.RecordedAt, m.Location)
}

type MeasurementJSON struct {
	MeasurementDB
	RecordedAt Time
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

func getRecentMeasurements(db *pgxpool.Pool) []*MeasurementDB {
	var measurements []*MeasurementDB

	err := pgxscan.Select(context.Background(), db, &measurements, `select * from measurements order by recorded_at desc limit 5`)
	if err != nil {
		quitOnError("pgxscan failed", err)
	}

	return measurements
}

func insertMeasurement(db *pgxpool.Pool, measurement *MeasurementDB) error {

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

	new := MeasurementJSON{}

	err := json.Unmarshal([]byte(data), &new)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem unmarshalling JSON: %v\n", err)
		os.Exit(1)
	}

	new.Location = `testing`

	newdb := MeasurementDB{}
	newdb.RecordedAt = new.RecordedAt.Time
	newdb.Temperature = new.Temperature
	newdb.Humidity = new.Humidity
	newdb.Pressure = new.Pressure
	newdb.Location = new.Location

	err = insertMeasurement(db, &newdb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem inserting row: %v\n", err)
	}

	measurements := getRecentMeasurements(db)

	for _, measurement := range measurements {
		fmt.Printf("%v\n", measurement)
	}

	for _, measurement := range measurements {
		jm := MeasurementJSON{}
		jm.Id = measurement.Id
		jm.RecordedAt.Time = measurement.RecordedAt
		jm.Temperature = measurement.Temperature
		jm.Humidity = measurement.Humidity
		jm.Pressure = measurement.Pressure
		jm.Location = measurement.Location

		jsonString, err := json.Marshal(jm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem marshalling json: %v\n", err)
		}

		fmt.Printf("%s\n", jsonString)
	}
}
