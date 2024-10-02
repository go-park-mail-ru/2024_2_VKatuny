package handler

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
)

// Accepts wrappedHandlerFunc and sets up CORS and content-type headers
// Returns function with headers
func HttpHeadersWrapper(wrappedHandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set up CORS headers
		w.Header().Set("Access-Control-Allow-Origin", BD.FRONTENDIP)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		
		// Set up content-type header
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			// should i return http.StatusOK?
			return
		}

		wrappedHandlerFunc(w, r)
	}
}
