-- name: GetRecentMeasurements :many
select * from measurements
order by recorded_at desc
limit $1;

-- name: GetTrends :one
select temperature, humidity, pressure 
from weather_trend(sqlc.arg(period));

-- name: InsertMeasurement :one
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
RETURNING *;
