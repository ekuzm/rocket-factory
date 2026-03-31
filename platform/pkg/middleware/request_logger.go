package middleware

import (
	"net/http"
	"time"

	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	"github.com/sirupsen/logrus"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.WithFields(logrus.Fields{
			"URL":         r.URL,
			"HTTP Method": r.Method,
		}).Debug("Starting HTTP handler...")

		start := time.Now()

		next.ServeHTTP(w, r)

		logger.WithFields(logrus.Fields{
			"URL":         r.URL,
			"HTTP Method": r.Method,
		}).Debug("Handler has worked for ", time.Since(start))
	})
}
