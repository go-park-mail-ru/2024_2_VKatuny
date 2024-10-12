package storage

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
)

// Deprecated. Do not use!
func SetSecureHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", BD.FRONTENDIP)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
}
