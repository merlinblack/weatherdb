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

-- name: GetHourlySummary :many
select
    hour,
    round(avg(temperature)::numeric,1)::text as temperature,
    round(avg(humidity)::numeric,1)::text as humidity,
    round(avg(pressure)::numeric,2)::text as pressure
from (
 select date_trunc('hour',recorded_at)::text as hour, temperature, humidity, pressure
 from public.measurements
 order by hour desc
 limit sqlc.arg(hours)::int * 60
) data
group by hour
order by hour;