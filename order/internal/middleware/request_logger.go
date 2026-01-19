package middleware

import (
	"log"
	"net/http"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handler at %s with %s http method start", r.URL.Path, r.Method)

		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("Handler has worket for %s", time.Since(start))
	})
}
