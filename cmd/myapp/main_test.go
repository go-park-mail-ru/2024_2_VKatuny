package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
)

func TestCreateWorker(t *testing.T) {
	t.Parallel()

	h := BD.WorkerHandlers{
		Users: map[string]BD.Worker{},
		Mu:    &sync.RWMutex{},
	}

	body := bytes.NewReader([]byte(`{"WorkerName":"Vasia","WorkerLastName":"Vasin",
	"WorkerBirthDate":"12-12-2012","WorkerEmail":"a@mail.ru","WorkerPassword":"pass"}`))

	expectedUsers := map[string]BD.Worker{
		"a@mail.ru": {
			ID:              1,
			WorkerName:      "Vasia",
			WorkerLastName:  "Vasin",
			WorkerBirthDate: "12-12-2012",
			WorkerEmail:     "a@mail.ru",
			WorkerPassword:  "pass",
		},
	}

	r := httptest.NewRequest("POST", "api/registration/worker", body)
	w := httptest.NewRecorder()
	h.HandleCreateWorker(w, r)
	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}
	if !reflect.DeepEqual(h.Users, expectedUsers) {
		t.Errorf("got %q, want %q", h.Users, expectedUsers)
	}
}

func TestCreateEmployer(t *testing.T) {
	t.Parallel()

	h := BD.EmployerHandlers{
		Users: map[string]BD.Employer{},
		Mu:    &sync.RWMutex{},
	}

	body := bytes.NewReader([]byte(`{"EmployerName":"Vasily","EmployerLastName":"Vasin",
	"EmployerPosition":"CEO","CompanyName":"Vasia Entertainment","CompanyDescription":"Vasia best company",
	"Website":"vasia.com","EmployerEmail":"vasia@gmail.com","EmployerPassword":"pass"}`))

	expectedUsers := map[string]BD.Employer{
		"vasia@gmail.com": {
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
	}

	r := httptest.NewRequest("POST", "api/registration/employer", body)
	w := httptest.NewRecorder()
	//workerHandler := handler.CreateWorkerHandler(&BD.HandlersWorker)

	h.HandleCreateEmployer(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	if !reflect.DeepEqual(h.Users, expectedUsers) {
		t.Errorf("got %q, want %q", h.Users, expectedUsers)
	}
}
