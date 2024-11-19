package routing

import (
	"fmt"
	"net/http"

	"github.com/merlinblack/weatherdb/internal/api/http/handlers"
	"github.com/merlinblack/weatherdb/internal/api/http/middleware"
	"github.com/merlinblack/weatherdb/internal/weather_repository"
)

func makeHandlerWithRepo(repo *weather_repository.Queries, fn func(w http.ResponseWriter, r *http.Request, repo *weather_repository.Queries)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, repo)
	})
}

func GetRouteChain(weather *weather_repository.Queries) http.Handler {

	mux := http.NewServeMux()

	mux.Handle(`GET /weather`, makeHandlerWithRepo(weather, handlers.RecentMeasurements))
	mux.Handle(`GET /trends`, makeHandlerWithRepo(weather, handlers.Trends))
	mux.HandleFunc(`GET /ping`, func(w http.ResponseWriter, _ *http.Request) { fmt.Fprintln(w, `pong`) })

	chain := middleware.ChainFinal(mux)
	chain.Use(middleware.LoggingMiddleware)

	return chain
}
