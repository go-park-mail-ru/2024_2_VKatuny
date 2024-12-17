package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/mailru/easyjson"
)

// UniversalMarshal marshal any struct to json (use it for any answer from handlers).
// Writes http status to http.ResponseWriter
func UniversalMarshal(w http.ResponseWriter, status int, body dto.JSONResponse) error {
	w.WriteHeader(status)
	_, _, err := easyjson.MarshalToHTTPResponseWriter(body, w)
	return err
}
