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
    location,
    hour
)
values
(
    $1, $2, $3, $4, $5, extract('epoch' from date_trunc('hour', timezone('Australia/Sydney', $1::timestamp)))
)
RETURNING *;

-- name: GetHourlySummary :many
select * from (select
    to_timestamp(hour)::timestamp as hour,
    round(avg(temperature)::numeric,1)::text as temperature,
    round(avg(humidity)::numeric,1)::text as humidity,
    round(avg(pressure)::numeric,2)::text as pressure
from (select * from measurements order by recorded_at desc limit sqlc.arg(hours)::int * 60) measurements
group by hour
order by hour desc
limit sqlc.arg(hours)::int) hour_measurements
order by hour;