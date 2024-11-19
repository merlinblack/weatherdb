package middleware

import "net/http"

type chain struct {
	next http.Handler
}

func Chain(handler http.Handler) *chain {
	return &chain{handler}
}

func (c *chain) Use(handler func(http.Handler) http.Handler) {
	c.next = handler(c.next)
}

func (c *chain) UseGroup(handlers ...func(http.Handler) http.Handler) {
	for _, handler := range handlers {
		c.Use(handler)
	}
}

func (c *chain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.next.ServeHTTP(w, r)
}
