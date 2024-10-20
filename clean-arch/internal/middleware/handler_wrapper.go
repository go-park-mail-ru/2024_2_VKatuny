package middleware

import (
	"net/http"
)

// HTTPHeadersWrapper Accepts wrappedHandlerFunc and sets up CORS and content-type headers
// Returns function with headers
func HTTPHeadersWrapper(wrappedHandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		wrappedHandlerFunc(w, r)
	}
}
