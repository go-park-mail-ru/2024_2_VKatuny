package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/sirupsen/logrus"
)

// AllowMethods checks if request method is allowed.
// Accepts http.Handler and array of allowed methods.
func AllowMethods(next http.Handler, allowedMethods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn := "middleware.AllowMethods"

		logger := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if logger == nil {
			fmt.Printf("function %s: unable to get logger from context. Processing without it", fn)
		}

		for _, allowedMethod := range allowedMethods {
			// compares methods without case
			if strings.EqualFold(r.Method, allowedMethod) {
				logger.Debugf("function %s: method %s is allowed", fn, r.Method)
				next.ServeHTTP(w, r)
				return
			}
		}
		
		logger.Debugf("function %s: method %s is not allowed", fn, r.Method)
		UniversalMarshal(w, http.StatusMethodNotAllowed, dto.JSONResponse{
			HTTPStatus: http.StatusMethodNotAllowed,
			Error:      http.StatusText(http.StatusMethodNotAllowed),
		})
	})
}
