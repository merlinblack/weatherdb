package middleware

import (
	"log"
	"net/http"
	"time"
)

type logging struct {
	handler http.Handler
}

func (m logging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	m.handler.ServeHTTP(w, r)

	elapsed := time.Since(start)

	log.Printf("[%s] [%s] [%s]\n", r.Method, r.URL.Path, elapsed)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return logging{next}
}