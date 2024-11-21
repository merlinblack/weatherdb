// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package weather

import (
	"context"
	"time"
)

const getHourlySummary = `-- name: GetHourlySummary :many
select
    hour,
    round(avg(temperature)::numeric,1)::text as temperature,
    round(avg(humidity)::numeric,1)::text as humidity,
    round(avg(pressure)::numeric,2)::text as pressure
from (
 select date_trunc('hour',recorded_at)::text as hour, temperature, humidity, pressure
 from public.measurements
 order by hour desc
 limit $1::int * 60
) data
group by hour
order by hour
`

type GetHourlySummaryRow struct {
	Hour        string
	Temperature string
	Humidity    string
	Pressure    string
}

func (q *Queries) GetHourlySummary(ctx context.Context, hours int32) ([]GetHourlySummaryRow, error) {
	rows, err := q.db.Query(ctx, getHourlySummary, hours)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetHourlySummaryRow
	for rows.Next() {
		var i GetHourlySummaryRow
		if err := rows.Scan(
			&i.Hour,
			&i.Temperature,
			&i.Humidity,
			&i.Pressure,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRecentMeasurements = `-- name: GetRecentMeasurements :many
select id, recorded_at, temperature, humidity, pressure, location from measurements
order by recorded_at desc
limit $1
`

func (q *Queries) GetRecentMeasurements(ctx context.Context, limit int32) ([]Measurement, error) {
	rows, err := q.db.Query(ctx, getRecentMeasurements, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Measurement
	for rows.Next() {
		var i Measurement
		if err := rows.Scan(
			&i.ID,
			&i.RecordedAt,
			&i.Temperature,
			&i.Humidity,
			&i.Pressure,
			&i.Location,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTrends = `-- name: GetTrends :one
select temperature, humidity, pressure 
from weather_trend($1)
`

func (q *Queries) GetTrends(ctx context.Context, period time.Duration) (Trend, error) {
	row := q.db.QueryRow(ctx, getTrends, period)
	var i Trend
	err := row.Scan(&i.Temperature, &i.Humidity, &i.Pressure)
	return i, err
}

const insertMeasurement = `-- name: InsertMeasurement :one
insert into measurements
(
    recorded_at, 
    temperature,
    humidity,
    pressure,
    location
)
values
(
    $1, $2, $3, $4, $5
)
RETURNING id, recorded_at, temperature, humidity, pressure, location
`

type InsertMeasurementParams struct {
	RecordedAt  time.Time
	Temperature float64
	Humidity    float64
	Pressure    float64
	Location    string
}

func (q *Queries) InsertMeasurement(ctx context.Context, arg InsertMeasurementParams) (Measurement, error) {
	row := q.db.QueryRow(ctx, insertMeasurement,
		arg.RecordedAt,
		arg.Temperature,
		arg.Humidity,
		arg.Pressure,
		arg.Location,
	)
	var i Measurement
	err := row.Scan(
		&i.ID,
		&i.RecordedAt,
		&i.Temperature,
		&i.Humidity,
		&i.Pressure,
		&i.Location,
	)
	return i, err
}
