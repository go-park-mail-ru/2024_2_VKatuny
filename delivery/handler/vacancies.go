package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
)

// HTTP GET. ?offset=10&num=5
func VacanciesHandler(vacanciesTable *BD.VacanciesHandler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		isoption := storage.Isoption(w, r)
		if isoption {
			return
		}
		storage.SetSecureHeaders(w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"status": 405, "error": "http request method isn't a GET"}`))
			return
		}

		queryParams := r.URL.Query()

		offsetStr := queryParams.Get("offset")
		if offsetStr == "" {
			w.WriteHeader(http.StatusBadRequest)  // HTTP 400
			w.Write([]byte(`{"status": 400, "error": "offset is empty"}`))
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) 
			w.Write([]byte(`{"status": 400, "error": "offset isn't number"}`))
			return
		}

		numStr := queryParams.Get("num")
		if numStr == "" {
			w.WriteHeader(http.StatusBadRequest) 
			w.Write([]byte(`{"status": 400, "error": "num is empty"}`))
			return
		}

		num, err := strconv.Atoi(numStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) 
			w.Write([]byte(`{"status": 400, "error": "count isn't number"}`))
			return
		}

		leftBound  := offset
		rightBound := offset + num
		// covering cases when offset is out of slice bounds
		if leftBound > int(vacanciesTable.Count) {
			rightBound = leftBound
		} else if rightBound > int(vacanciesTable.Count) { 
			rightBound = int(vacanciesTable.Count)
		}

		vacanciesTable.Mutex.Lock()
		vacancies := vacanciesTable.Vacancy[leftBound : rightBound]
		vacanciesTable.Mutex.Unlock()

		response := map[string]interface{}{
			"status": 200,
			"vacancies": vacancies,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"status": 500, "error": "encoding error"}`))
			return
		}
	
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}