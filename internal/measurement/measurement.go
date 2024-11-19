package measurement

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/merlinblack/weatherdb/internal/weather_repository"
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

func ToJSON(m *weather_repository.Measurement) string {
	jm := MeasurementJSON{}
	jm.RecordedAt.Time = m.RecordedAt
	jm.Temperature = m.Temperature
	jm.Humidity = m.Humidity
	jm.Pressure = m.Pressure

	jsonString, err := json.Marshal(jm)
	if err != nil {
		log.Fatalf("Problem marshalling json: %v\n", err)
	}

	return string(jsonString)
}

func FromJSON(data string) weather_repository.Measurement {
	jm := MeasurementJSON{}

	err := json.Unmarshal([]byte(data), &jm)
	if err != nil {
		log.Fatalf("Problem unmarshalling JSON: %v\n", err)
	}

	m := weather_repository.Measurement{}
	m.RecordedAt = jm.RecordedAt.Time
	m.Temperature = jm.Temperature
	m.Humidity = jm.Humidity
	m.Pressure = jm.Pressure

	return m
}
