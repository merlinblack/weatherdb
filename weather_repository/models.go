// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package weather_repository

import (
	"time"
)

type Measurement struct {
	ID          int64
	RecordedAt  time.Time
	Temperature float64
	Humidity    float64
	Pressure    float64
	Location    string
}

type Trend struct {
	Temperature string
	Humidity    string
	Pressure    string
}
