package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/merlinblack/weatherdb/internal/api/http/routing"
	"github.com/merlinblack/weatherdb/internal/repository/weather"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"WEATHERDB_HOST" env-default:"localhost"`
	Port     string `json:"port" env:"WEATHERDB_PORT" env-default:"5432"`
	Username string `json:"username" env:"WEATHERDB_USERNAME" env-default:"weather"`
	Password string `json:"password" env:"WEATHERDB_PASSWORD" env-default:"weather"`
	Name     string `json:"name" env:"WEATHERDB_NAME" env-default:"weather"`
}

type APIConfig struct {
	WritePassword string `json:"password" env:"WEATHERDB_API_PASS" env-default:"weather"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
	API      APIConfig      `json:"API"`
}

func getDsn(cfg *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
}

func getConnection(dsn string) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	return db
}

/*
func test(weather *weather_repository.Queries) {

	data := `{"recorded_at":"2024-04-02 23:24", "temperature": 24.3, "humidity": 67.32, "pressure": 1019.2}`

	new := MeasurementFromJSON(data)
	new.Location = `testing`

	_, err := weather.InsertMeasurement(context.Background(),
		weather_repository.InsertMeasurementParams{
			RecordedAt:  new.RecordedAt,
			Temperature: new.Temperature,
			Humidity:    new.Humidity,
			Pressure:    new.Pressure,
			Location:    new.Location,
		})

	if err != nil {
		log.Printf("Problem inserting row: %v\n", err)
	}

	conn.Exec(context.Background(), `delete from measurements where location = 'testing'`)
}
*/

func main() {

	cfg := &Config{}

	err := cleanenv.ReadConfig(`config.json`, cfg)
	if err != nil {
		log.Printf("There was a problem reading te configuration file config.json: %v\n", err)
	}

	log.Printf("Configuration: %v\n", cfg)

	conn := getConnection(getDsn(cfg))
	defer conn.Close()

	weatherdb := weather.New(conn)

	chain := routing.GetRouteChain(weatherdb)

	server := &http.Server{
		Addr:    `:3000`,
		Handler: chain,
	}

	go func() {
		log.Printf("Listening on %v\n", server.Addr)
		err := server.ListenAndServe()

		if err == http.ErrServerClosed {
			log.Println(`Server closed, no longer accepting connections`)
		} else {
			log.Fatalf("Problem starting http server: %v\n", err)
		}
	}()

	// Wait for OS interrupt (pkill -2 weatherdb)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	// put a newline after the possible '^C' that is now displayed if the user pressed ^C
	fmt.Println(``)

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	log.Println(`Graceful shutdown, current requests have 30 seconds to finish`)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Problem shutting down: %v\n", err)
	}

	log.Println(`Bye!`)
}
