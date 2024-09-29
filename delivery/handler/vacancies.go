package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
)

// HTTP GET. ?offset=10&count=5
func VacanciesHandler(vacanciesTable *BD.VacanciesHandler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		queryParams := r.URL.Query()

		offsetStr := queryParams.Get("offset")
		if offsetStr == "" {
			w.WriteHeader(http.StatusBadRequest)  // HTTP 400
			w.Write([]byte(`{"statusCode": 400, "error": "offset is empty"}`))
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) 
			w.Write([]byte(`{"statusCode": 400, "error": "offset isn't number"}`))
			return
		}

		countStr := queryParams.Get("count")
		if countStr == "" {
			w.WriteHeader(http.StatusBadRequest) 
			w.Write([]byte(`{"statusCode": 400, "error": "count is empty"}`))
			return
		}

		count, err := strconv.Atoi(queryParams.Get("count"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) 
			w.Write([]byte(`{"statusCode": 400, "error": "offset isn't number"}`))
			return
		}

		responseJson := "{"
		for _, vacancy := range vacanciesTable.Vacancy[offset : offset + count] {
			vacanciesTable.Mutex.Lock()

			res, _ := json.Marshal(vacancy)
			responseJson += (string(res) + ",")
			
			vacanciesTable.Mutex.Unlock()
		}
		responseJson += "}"
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseJson))
	}
	return http.HandlerFunc(fn)
}
