// Package middleware is for universal things for handlers
package middleware

import (
	"net/http"
)

// IsOption is for OPTIONS requests
// Deprecated. Do not use!
// OPTIOS metods handler sets headers
func IsOption(w http.ResponseWriter, r *http.Request) bool {

	if r.Method == http.MethodOptions {
		SetSecureHeaders(w)
		return true
	}
	return false
}
