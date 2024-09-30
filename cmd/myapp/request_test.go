package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/delivery/handler"
)

type TestCase1 struct {
	ID         string
	Response   string
	StatusCode int
}

// func GetUser(w http.ResponseWriter, r *http.Request) {
// 	key := r.FormValue("id")
// 	if key == "42" {
// 		w.WriteHeader(http.StatusOK)
// 		io.WriteString(w, `{"status": 200, "resp": {"user": 42}}`)
// 	} else {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		io.WriteString(w, `{"status": 500, "err": "db_error"}`)
// 	}
// }

func TestGetUser(t *testing.T) {
	cases := []TestCase1{
		TestCase1{
			ID:         "42",
			Response:   `{"status": 200, "resp": {"user": 42}}`,
			StatusCode: http.StatusOK,
		},
		TestCase1{
			ID:         "500",
			Response:   `{"status": 500, "err": "db_error"}`,
			StatusCode: http.StatusInternalServerError,
		},
	}
	for caseNum, item := range cases {
		url := "http://localhost:8080/api/v1/authorized" // + item.ID
		req := httptest.NewRequest("POST", url, nil)
		w := httptest.NewRecorder()

		//f := handler.AuthorizedHandler()
		handler.Fn(w, req)
		fmt.Println(w.Body)
		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)

		bodyStr := string(body)
		if bodyStr != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				caseNum, bodyStr, item.Response)
		}
	}
}
