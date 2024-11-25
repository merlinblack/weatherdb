// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package weather

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Measurement struct {
	ID          int64
	RecordedAt  time.Time
	Temperature float64
	Humidity    float64
	Pressure    float64
	Location    string
	Hour        pgtype.Int8
}

type Trend struct {
	Temperature string
	Humidity    string
	Pressure    string
}
