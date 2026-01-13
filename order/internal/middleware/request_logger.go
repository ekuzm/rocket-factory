package middleware

import (
	"log"
	"net/http"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		url := r.URL
		method := r.Method

		log.Printf("Handler at %v URL with %v Method run", url, method)

		next.ServeHTTP(w, r)

		log.Printf("Handler at %v URL with %v Method has worked for %v", url, method, time.Since(startTime))
	})
}
