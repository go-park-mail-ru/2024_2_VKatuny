// Package middleware ontains middleware of project
package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/crw"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/metrics"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

// AccessLogger logs incoming requests
func AccessLogger(next http.Handler, logger *logrus.Logger, metrics *metrics.Metrics) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		requestID, err := utils.GenerateRequestID()

		ctx := r.Context()
		if err != nil {
			logger.Errorf("can't generate request id with error: %s, ignoring it...", err)
		} 
		ctx = context.WithValue(ctx, dto.RequestIDContextKey, requestID)

		ws := crw.NewResponseWriterStatus(w)

		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"request_id": requestID,
		}).Info("Request received")

		start := time.Now()
		next.ServeHTTP(ws, r.WithContext(ctx))
		end := time.Since(start)

		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"request_id": requestID,
			"elapsed": end,
			"status": ws.Status(),
		}).Info("Response")
		
		statusString := strconv.Itoa(ws.Status())
		metrics.Hits.WithLabelValues(r.Method, r.URL.Path, statusString).Inc()
		metrics.Timings.WithLabelValues(r.Method, r.URL.Path).Observe(end.Seconds())
	})
}
