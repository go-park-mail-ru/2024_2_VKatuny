package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

// UniversalMarshal marshal any struct to json (use it for any answer from handlers).
// Writes http status to http.ResponseWriter
func UniversalMarshal(w http.ResponseWriter, status int, body interface{}) error {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return fmt.Errorf(dto.MsgUnableToMarshalJSON)
	}
	return nil
}
