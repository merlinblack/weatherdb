package measurement

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/merlinblack/weatherdb/internal/repository/weather"
)

type measurementTime struct {
	time.Time
}

const timeJSONLayout = `2006-01-02 15:04`

func (t measurementTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(t.Format(timeJSONLayout))
}

func (t *measurementTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}
	t.Time, err = time.Parse(timeJSONLayout, s)
	return
}

type measurementJSON struct {
	RecordedAt  measurementTime `json:"recorded_at"`
	Temperature float64         `json:"temperature"`
	Humidity    float64         `json:"humidity"`
	Pressure    float64         `json:"pressure"`
}

func ToJSON(m *weather.Measurement) string {
	jm := measurementJSON{}
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

func FromJSON(data string) weather.Measurement {
	jm := measurementJSON{}

	err := json.Unmarshal([]byte(data), &jm)
	if err != nil {
		log.Fatalf("Problem unmarshalling JSON: %v\n", err)
	}

	m := weather.Measurement{}
	m.RecordedAt = jm.RecordedAt.Time
	m.Temperature = jm.Temperature
	m.Humidity = jm.Humidity
	m.Pressure = jm.Pressure

	return m
}
