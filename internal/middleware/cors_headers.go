package middleware

import (
	"fmt"
	"net/http"
)

// SetSecurityAndOptionsHeaders Accepts function next and sets up CORS and content-type headers
// Returns wrapped function next with headers
func SetSecurityAndOptionsHeaders(next http.Handler, frontURI string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(1)
		// Set up CORS headers
		w.Header().Set("Access-Control-Allow-Origin", frontURI)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, X-CSRF-Token")
		w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")

		if r.Method == http.MethodOptions {
			return
		}
		// Set up content-type header
		w.Header().Set("Content-Type", "application/json")
		
		next.ServeHTTP(w, r)
	})
}
