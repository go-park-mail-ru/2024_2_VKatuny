package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

func TestCreateWorker(t *testing.T) {
	t.Parallel()

	h := WorkerHandlers{
		users: map[string]Worker{},
		mu:    &sync.RWMutex{},
	}

	body := bytes.NewReader([]byte(`{"WorkerName":"Vasia","WorkerLastName":"Vasin",
	"WorkerBirthDate":"12-12-2012","WorkerEmail":"a@mail.ru","WorkerPassword":"pass"}`))

	expectedUsers := map[string]Worker{
		"a@mail.ru": {
			ID:              1,
			WorkerName:      "Vasia",
			WorkerLastName:  "Vasin",
			WorkerBirthDate: "12-12-2012",
			WorkerEmail:     "a@mail.ru",
			WorkerPassword:  "pass",
		},
	}

	r := httptest.NewRequest("POST", "/registration/worker", body)
	w := httptest.NewRecorder()

	h.HandleCreateWorker(w, r)
	//fmt.Println(h.users, expectedUsers)
	//fmt.Println(h.users[0].ID, expectedUsers[0].ID)
	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	if !reflect.DeepEqual(h.users, expectedUsers) {
		t.Errorf("got %q, want %q", h.users, expectedUsers)
	}
}

func TestCreateEmployer(t *testing.T) {
	t.Parallel()

	h := EmployerHandlers{
		users: map[string]Employer{},
		mu:    &sync.RWMutex{},
	}

	body := bytes.NewReader([]byte(`{"EmployerName":"Vasily","EmployerLastName":"Vasin",
	"EmployerPosition":"CEO","CompanyName":"Vasia Entertainment","CompanyDescription":"Vasia best company",
	"Website":"vasia.com","EmployerEmail":"vasia@gmail.com","EmployerPassword":"pass"}`))

	expectedUsers := map[string]Employer{
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

	r := httptest.NewRequest("POST", "/registration/employer", body)
	w := httptest.NewRecorder()

	h.HandleCreateEmployer(w, r)

	//fmt.Println(h.users, expectedUsers)
	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	if !reflect.DeepEqual(h.users, expectedUsers) {
		t.Errorf("got %q, want %q", h.users, expectedUsers)
	}
}

//var expectedJSON = `[{"id":1,"name":"Afanasiy"},{"id":2,"name":"Ka"}]`

// func TestGetUsers(t *testing.T) {

// 	h := Handlers{
// 		users: []User{
// 			{
// 				ID:       1,
// 				Name:     "Afanasiy",
// 				Password: "1234",
// 			},
// 			{
// 				ID:       2,
// 				Name:     "Ka",
// 				Password: "jdjfaljhfljehfs;l3345354",
// 			},
// 		},
// 		mu: &sync.Mutex{},
// 	}

// 	t.Parallel()

// 	r := httptest.NewRequest("GET", "/users/", nil)
// 	w := httptest.NewRecorder()

// 	h.HandleListUsers(w, r)

// 	if w.Code != http.StatusOK {
// 		t.Error("status is not ok")
// 	}

// 	bytes, _ := ioutil.ReadAll(w.Body)
// 	if strings.Trim(string(bytes), "\n") != expectedJSON {
// 		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
// 	}
// }
