package middleware

import (
	"log"
	"net/http"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		log.Printf("Handler at %v URL with %v Method run", r.URL, r.Method)

		next.ServeHTTP(w, r)

		log.Printf("Handler has worked: %v", time.Since(startTime))
	})
}
