package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/merlinblack/weatherdb/internal/api/http/routing"
	"github.com/merlinblack/weatherdb/internal/weather_repository"
)

func getDsn() string {
	user := `nigel`
	password := os.Getenv(`WEATHERDB_PASSWORD`)
	host := `octavo.local`
	port := 5432
	database := `weather_test`

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, database)
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

	conn := getConnection(getDsn())
	defer conn.Close()

	weather := weather_repository.New(conn)

	chain := routing.GetRouteChain(weather)

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
