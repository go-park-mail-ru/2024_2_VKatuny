package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Panic recovers http.handler's panic.
// Accepts covering http.Handler.
// Returns http.Handler.
func Panic(next http.Handler, logger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcName := "midleware.Panic"
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
