package storage

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func UniversalMarshal(w http.ResponseWriter, status int, body interface{}) error {
	if body == nil {
		w.WriteHeader(status)
		return nil
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Err while marshal")
	} else {
		w.WriteHeader(status)
		return nil
	}

}
