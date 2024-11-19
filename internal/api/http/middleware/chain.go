package middleware

import "net/http"

type Chain struct {
	next http.Handler
}

func ChainFinal(handler http.Handler) *Chain {
	return &Chain{handler}
}

func (c *Chain) Use(handler func(http.Handler) http.Handler) {
	c.next = handler(c.next)
}

func (c *Chain) UseGroup(handlers ...func(http.Handler) http.Handler) {
	for _, handler := range handlers {
		c.Use(handler)
	}
}

func (c *Chain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.next.ServeHTTP(w, r)
}
