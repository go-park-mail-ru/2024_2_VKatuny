package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func UniversalMarshal(w http.ResponseWriter, status int, body interface{}) error {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return fmt.Errorf("Err while marshal")
	}
	return nil
}
