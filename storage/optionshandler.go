package storage

import (
	"net/http"
)

// OPTIOS metods handler sets headers
func Isoption(w http.ResponseWriter, r *http.Request) bool {

	if r.Method == http.MethodOptions {
		SetSecureHeaders(w)
		return true
	}
	return false
}
