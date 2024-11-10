// Package middleware ontains middleware of project
package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// AccessLogger logs incoming requests
func AccessLogger(next http.Handler, logger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("Request received")
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			// "status": w.Status(), how to get response status?
			"elapsed": time.Since(start),
		}).Info("Response")
	})
}
