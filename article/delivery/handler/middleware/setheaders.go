package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/inmemorydb"
)

// SetSecureHeaders sets secure headers
// Deprecated. Do not use!
func SetSecureHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", inmemorydb.FRONTENDIP)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
}
