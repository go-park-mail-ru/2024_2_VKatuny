package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func AccessLogger(next http.Handler, logger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			// "status": w.Status(), how to get response status?
			"execution": time.Since(start),
		}).Info("Request")
	})
}
