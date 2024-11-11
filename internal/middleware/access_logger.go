// Package middleware ontains middleware of project
package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

// AccessLogger logs incoming requests
func AccessLogger(next http.Handler, logger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		requestID, err := utils.GenerateRequestID()
		ctx := r.Context()
		if err != nil {
			logger.Errorf("can't generate request id with error: %s, ignoring it...", err)
			requestID = ""
		} else {
			ctx = context.WithValue(ctx, dto.RequestIDContextKey, requestID)
		}
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"request_id": requestID,
		}).Info("Request received")
		start := time.Now()
		next.ServeHTTP(w, r.WithContext(ctx))
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"request_id": requestID,
			"elapsed": time.Since(start),
		}).Info("Response")
	})
}
