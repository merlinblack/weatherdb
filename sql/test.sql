
CREATE TABLE if not exists public.trend (
    temperature text not null,
    humidity text not null,
    pressure text not null
);


CREATE or replace FUNCTION public.weather_trend(since interval) RETURNS SETOF public.trend
    LANGUAGE plpgsql
    AS $$
begin
  return query (
    select 
      case when temperature > 0 then 'increasing' when temperature < 0 then 'decreasing' else 'stable' end temperature,
      case when humidity > 0 then 'increasing' when humidity < 0 then 'decreasing' else 'stable' end humidity,
      case when pressure > 0 then 'increasing' when pressure < 0 then 'decreasing' else 'stable' end pressure
    from (
      select
        regr_slope(temperature, extract(epoch from recorded_at)) as temperature,
        regr_slope(humidity, extract(epoch from recorded_at)) as humidity,
        regr_slope(pressure, extract(epoch from recorded_at)) as pressure
      from public.measurements m 
      where recorded_at > now() - since
    )
  );
end; $$;

