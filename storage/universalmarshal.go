package storage

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func UniversalMarshal(w http.ResponseWriter, status int, body interface{}) error {

	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Err while marshal")
	} else {
		w.WriteHeader(http.StatusOK)
		return nil
	}

}
