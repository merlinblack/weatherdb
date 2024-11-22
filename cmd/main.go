package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/merlinblack/weatherdb/internal/api/http/routing"
	"github.com/merlinblack/weatherdb/internal/config"
	"github.com/merlinblack/weatherdb/internal/repository/weather"

	_ "embed"
)

//go:embed VERSION
var gitversion string

func getDsn(cfg *config.Config) string {
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

func main() {

	var versionFlag bool
	var configPath string

	flag.BoolVar(&versionFlag, `version`, false, `display version and exit.`)
	flag.BoolVar(&versionFlag, `v`, false, `display version and exit.`)
	flag.StringVar(&configPath, `config`, `config.json`, "`path` to a configuration file")
	flag.StringVar(&configPath, `c`, `config.json`, "`path` to a configuration file")
	flag.Parse()

	fmt.Printf("WeatherDB\nVersion: %s\n", gitversion)
	if versionFlag {
		// Exit now if only showing the version.
		return
	}

	cfg := &config.Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		log.Fatalf("There was a problem reading te configuration file %s: %v\n", configPath, err)
	}

	log.Printf("Using DB: %s on %s\n", cfg.Database.Name, cfg.Database.Host)

	conn := getConnection(getDsn(cfg))
	defer conn.Close()

	weatherdb := weather.New(conn)

	chain := routing.GetRouteChain(cfg, weatherdb)

	server := &http.Server{
		Addr:    cfg.ListenAddress,
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
	// Put a newline after the possible '^C' that is now displayed if the user pressed ^C
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
