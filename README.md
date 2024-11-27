WeatherDB
---------

A Go replacement of my old php based back end for storing and retreiving weather data collected from a RaspberryPI.

- Uses standard lib http.
- pgx for connecting to Postgres
- sqlc for generating the 'repo' code, but only really because I wanted to play with that
- cleanenv for nice and easy config file reading

Started as a project to learn Go, it took awhile to gain momentum, but now is fully functional replacement for my old php code that runs an order of magnitude faster with less RAM usage.

Copying a single file to update is nice ;)
