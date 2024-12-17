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
	jsonBytes, err := easyjson.Marshal(body)
	if err != nil {
		return err
	}

	_, err = w.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil

	// w.WriteHeader(status)
	// var err error
	// if body.Body == nil {
	// 	body.Body = dto.JSONLogoutForm{}
	// }
	// _, _, err = easyjson.MarshalToHTTPResponseWriter(body, w)

	// return err
}
