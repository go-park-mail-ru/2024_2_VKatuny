package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/inmemorydb"
)

// SetSecurityAndOptionsHeaders Accepts funcion next and sets up CORS and content-type headers
// Returns wrapped function next with headers
func SetSecurityAndOptionsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set up CORS headers
		w.Header().Set("Access-Control-Allow-Origin", inmemorydb.FRONTENDIP)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

		if r.Method == http.MethodOptions {
			// should i return http.StatusOK?
			return
		}
		// Set up content-type header
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
