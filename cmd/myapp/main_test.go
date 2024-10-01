package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	// "reflect"
	"sync"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/delivery/handler"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
	"github.com/stretchr/testify/require"
)

// Accepts jsonable string or struct that can be converted into the JSON format
// Returns readeble JSON with indents and error
// func readableJSON(jsonable interface{}) (string, error) {
// 	// 4 spaces for tabulation
// 	json, err := json.MarshalIndent(jsonable, "", "    ")
//     if err != nil {
//         return "", fmt.Errorf("during marshalling got: %v", err)
//     }
// 	return string(json), nil
// }


func TestCreateWorker(t *testing.T) {
	// t.Parallel()

	// disable logs
	log.SetOutput(io.Discard)

	h := &BD.WorkerHandlers{
		Users: map[string]BD.Worker{},
		Mu:    &sync.RWMutex{},
	}

	tests := map[string]struct {
		requestBody string
		expectedStatus int
		expectedUser   BD.Worker
	}{
		// rename test to more representative one
		"<test_name>": {
			requestBody: `{"WorkerName": "Vasia", "WorkerLastName": "Vasin",
	            "WorkerBirthDate": "12-12-2012", "WorkerEmail": "a@mail.ru", "WorkerPassword": "pass"}`,
			expectedStatus: http.StatusOK,
			expectedUser: BD.Worker{
				ID:              1,
				WorkerName:      "Vasia",
				WorkerLastName:  "Vasin",
				WorkerBirthDate: "12-12-2012",
				WorkerEmail:     "a@mail.ru",
				WorkerPassword:  "pass",
			},
		},
	}

	for testName, test := range tests {
		t.Run(fmt.Sprintf("testing %v", testName), func(t *testing.T) {
			// t.Parallel()

			body := bytes.NewReader([]byte(test.requestBody))

			request := httptest.NewRequest("POST", "/api/v1/registration/worker", body)
			response := httptest.NewRecorder()
			
			handler := handler.CreateWorkerHandler(h)
			handler.ServeHTTP(response, request)

			require.Equal(t, test.expectedStatus, response.Code,
				"handler returned response status %v, expected %v", response.Code, http.StatusOK)

			jsonExpected, err := json.Marshal(test.expectedUser)

			require.NoError(t, err, "while marshalling got: %v", err)
			
			require.Equal(t, string(jsonExpected), response.Body.String(),
				"handler returned response body %v, expected %v", response.Body.String(), string(jsonExpected))
		})
	}
}

func TestCreateEmployer(t *testing.T) {
	// t.Parallel()

	// disable logs
	log.SetOutput(io.Discard)

	h := &BD.EmployerHandlers{
		Users: map[string]BD.Employer{},
		Mu:    &sync.RWMutex{},
	}

	tests := map[string]struct {
		requestBody string
		expectedStatus int
		expectedUser   BD.Employer
	}{
		// rename test to more representative one
		"<test_name>": {
			requestBody: `{"EmployerName":"Vasily","EmployerLastName":"Vasin",
                "EmployerPosition":"CEO","CompanyName":"Vasia Entertainment","CompanyDescription":"Vasia best company",
                "Website":"vasia.com","EmployerEmail":"vasia@gmail.com","EmployerPassword":"pass"}`,
			expectedStatus: http.StatusOK,
			expectedUser: BD.Employer{
				ID:                 1,
				EmployerName:       "Vasily",
				EmployerLastName:   "Vasin",
				EmployerPosition:   "CEO",
				CompanyName:        "Vasia Entertainment",
				CompanyDescription: "Vasia best company",
				Website:            "vasia.com",
				EmployerEmail:      "vasia@gmail.com",
				EmployerPassword:   "pass",
			},
		},
	}

	for testName, test := range tests {
		t.Run(fmt.Sprintf("testing %v", testName), func(t *testing.T) {
			body := bytes.NewReader([]byte(test.requestBody))

			request := httptest.NewRequest("POST", "/api/registration/employer", body)
			response := httptest.NewRecorder()
			
			handler := handler.CreateEmployerHandler(h)
			handler.ServeHTTP(response, request)


			require.Equal(t, test.expectedStatus, response.Code,
				"handler returned response status %v, expected %v", response.Code, http.StatusOK)

			jsonExpected, err := json.Marshal(test.expectedUser)

			require.NoError(t, err, "while marshalling got: %v", err)
			
			require.Equal(t, string(jsonExpected), response.Body.String(),
				"handler returned response body %v, expected %v", response.Body.String(), string(jsonExpected))
		})
	}
}

// func TestLogin(t *testing.T) {
// 	// t.Parallel()

// 	// disable log
// 	log.SetOutput(io.Discard)

// 	h := &BD.EmployerHandlers{
// 		Sessions: map[string]uint64{},
// 		Mu:    &sync.RWMutex{},
// 	}

// 	// tests := map[string]struct{
// 	// 	table interface{}		
// 	// }{

// 	// }

// 	body := bytes.NewReader([]byte(`{"userType":"employer", "login": "vasya@gmail.com", "password": "pwd"}`))
// 	request := httptest.NewRequest("POST", "/api/v1/login", )

// }

func TestAuthorization(t *testing.T) {
	// t.Parallel()

	// diasble log
	log.SetOutput(io.Discard)

	t.Run("test authorized", func(t *testing.T) {
		table := &BD.HandlersWorker
		table.Sessions = make(map[string]uint64, 0)

		// saving selected session SID into the table
		// and sending it as a cookie
		SID := storage.RandStringRunes(32)
		table.Sessions[SID] = 1

		request := httptest.NewRequest("POST", "/api/v1/vacancies", nil)
		response := httptest.NewRecorder()

		cookie := &http.Cookie{
			Name: "session_id1",
			Value: SID,
			Expires: time.Now().Add(10 * time.Hour),
			HttpOnly: true,	
		}

		request.AddCookie(cookie)
		
		handler := handler.AuthorizedHandler()
		handler.ServeHTTP(response, request)

		// StatusOK means that user is authorized (sent cookie and saved cookie are equals)
		// in other case 401 means that user is not authorized 
		require.Equal(t, http.StatusOK, response.Code,
			"handler returned response status %v, expected %v", response.Code, http.StatusOK)
	})

	t.Run("test unauthorized", func(t *testing.T) {
		table := &BD.HandlersWorker
		table.Sessions = make(map[string]uint64, 0)

		// saving selected session SID into the table
		// and sending it as a cookie
		SID := storage.RandStringRunes(32)
		table.Sessions[SID] = 2

		request := httptest.NewRequest("POST", "/api/v1/vacancies", nil)
		response := httptest.NewRecorder()

		cookie := &http.Cookie{
			Name: "session_id1",
			Value: "dgahsgjgfgsjdhjfgjmghjmfghmnfhgn",  // random cookie
			Expires: time.Now().Add(10 * time.Hour),
			HttpOnly: true,	
		}

		request.AddCookie(cookie)
		
		handler := handler.AuthorizedHandler()
		handler.ServeHTTP(response, request)

		// expecting 401. user is unauthorized
		require.Equal(t, http.StatusUnauthorized, response.Code,
			"handler returned response status %v, expected %v", response.Code, http.StatusOK)
	})
}
