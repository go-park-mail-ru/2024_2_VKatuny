package middleware

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/sirupsen/logrus"
)

// SetLogger adds logger to the context
// and put into the handler
func SetLogger(next http.Handler, logger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {	
		logger.Debug("Setting up logger")
		ctx := r.Context()
		ctx = context.WithValue(ctx, dto.LoggerContextKey, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
