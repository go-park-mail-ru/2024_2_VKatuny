package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type vacanciesListResponse struct {
	Id          uint64 `json:id"` 
	Position    string `json:"position"`
	Discription string `json:"description"`
	Employer    string `json:"employer"`
	Location    string `json:"location"`
	Salary      string `json:"salary"`
	CreatedAt   string `json:"createdAt"`
}

// HTTP GET. ?offset=10&count=5
func VacanciesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		queryParams := r.URL.Query()

		offsetStr := queryParams.Get("offset")
		if offsetStr == "" {
			w.WriteHeader(http.StatusBadRequest)  // HTTP 400
			w.Write([]byte(`{"statusCode": 400, "error": "offset is empty"}`))
			return
		}

		if offset, err := strconv.Atoi(offsetStr); err != nil {
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

		if count, err := strconv.Atoi(queryParams.Get("count")); err != nil {
			w.WriteHeader(http.StatusBadRequest) 
			w.Write([]byte(`{"statusCode": 400, "error": "offset isn't number"}`))
			return
		}

		
	})
}
