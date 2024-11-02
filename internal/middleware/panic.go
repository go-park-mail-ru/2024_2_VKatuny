package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/sirupsen/logrus"
)

// Panic recovers http.handler's panic.
// Accepts covering http.Handler.
// Returns http.Handler.
func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcName := "midleware.Panic"
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("%s: can't get logger from context\n", funcName)
		}
		logger.Debugf("%s: entering", funcName)
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("%s: recovered panic: %v", funcName, err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
