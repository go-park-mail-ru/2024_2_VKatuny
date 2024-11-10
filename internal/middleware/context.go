package middleware

import (
	"context"
	"net/http"
	"github.com/sirupsen/logrus"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

// SetContext adds logger to the context
// and put into the handler
func SetContext(next http.Handler, logger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Setting up context")
		ctx := r.Context()
		ctx = context.WithValue(ctx, dto.LoggerContextKey, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
