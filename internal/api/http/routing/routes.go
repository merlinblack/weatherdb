package routing

import (
	"fmt"
	"net/http"

	"github.com/merlinblack/weatherdb/internal/api/http/handlers"
	"github.com/merlinblack/weatherdb/internal/api/http/middleware"
	"github.com/merlinblack/weatherdb/internal/config"
	"github.com/merlinblack/weatherdb/internal/repository/weather"
)

func makeHandlerWithRepo(repo *weather.Queries, fn func(w http.ResponseWriter, r *http.Request, repo *weather.Queries)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, repo)
	})
}

func makeHandlerWithRepoAndConfig(
	cfg *config.Config,
	repo *weather.Queries,
	fn func(w http.ResponseWriter, r *http.Request, cfg *config.Config, repo *weather.Queries)) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, cfg, repo)
	})
}

func GetRouteChain(cfg *config.Config, weather *weather.Queries) http.Handler {

	mux := http.NewServeMux()

	mux.Handle(`GET /weather`, makeHandlerWithRepo(weather, handlers.RecentMeasurements))
	mux.Handle(`POST /weather`, makeHandlerWithRepoAndConfig(cfg, weather, handlers.InsertMeasurement))
	mux.Handle(`GET /trends`, makeHandlerWithRepo(weather, handlers.Trends))
	mux.Handle(`GET /summary`, makeHandlerWithRepo(weather, handlers.HourlySummary))
	mux.HandleFunc(`GET /ping`, func(w http.ResponseWriter, _ *http.Request) { fmt.Fprintln(w, `pong`) })

	chain := middleware.Chain(mux)
	chain.Use(middleware.LoggingMiddleware)

	return chain
}
