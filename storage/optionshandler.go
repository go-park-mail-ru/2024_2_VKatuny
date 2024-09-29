package storage

import (
	"net/http"
)

func Isoption(w http.ResponseWriter, r *http.Request) bool {

	if r.Method == http.MethodOptions {
		SetSecureHeaders(w)
		return true
	}
	return false
}
