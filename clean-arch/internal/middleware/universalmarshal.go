package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// UniversalMarshal marshal any struct to json (use it for any answer from handlers).
// Writes http status to http.ResposeWriter
func UniversalMarshal(w http.ResponseWriter, status int, body interface{}) error {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return fmt.Errorf("Err while marshal")
	}
	return nil
}
