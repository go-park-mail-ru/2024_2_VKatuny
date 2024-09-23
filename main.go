package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type WorkerHandlers struct {
	users []Worker
	mu    *sync.Mutex
}

type WorkerInput struct {
	WorkerName      string // `json:"WorkerName"`
	WorkerLastName  string // `json:"WorkerLastName"`
	WorkerBirthDate string // `json:"WorkerBirthDate"`
	WorkerEmail     string // `json:"WorkerEmail"`
	WorkerPassword  string // `json:"WorkerPassword"`
}

type Worker struct {
	ID              uint64 //`json:"id"`
	WorkerName      string // `json:"WorkerName"`
	WorkerLastName  string // `json:"WorkerLastName"`
	WorkerBirthDate string // `json:"WorkerBirthDate"`
	WorkerEmail     string // `json:"WorkerEmail"`
	WorkerPassword  string // `json:"WorkerPassword"`
}

type EmployerHandlers struct {
	users []Employer
	mu    *sync.Mutex
}

type EmployerInput struct {
	EmployerName       string
	EmployerLastName   string
	EmployerPosition   string
	CompanyName        string
	CompanyDescription string
	Website            string
	EmployerEmail      string
	EmployerPassword   string
}

type Employer struct {
	ID                 uint64
	EmployerName       string
	EmployerLastName   string
	EmployerPosition   string
	CompanyName        string
	CompanyDescription string
	Website            string
	EmployerEmail      string
	EmployerPassword   string
}

func (h *WorkerHandlers) HandleCreateWorker(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	newUserInput := new(WorkerInput)
	err := decoder.Decode(newUserInput)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}

	fmt.Println(newUserInput)
	h.mu.Lock()

	var id uint64 = 0
	if len(h.users) > 0 {
		id = h.users[len(h.users)-1].ID + 1
	}

	h.users = append(h.users, Worker{
		ID:              id,
		WorkerName:      newUserInput.WorkerName,
		WorkerLastName:  newUserInput.WorkerLastName,
		WorkerBirthDate: newUserInput.WorkerBirthDate,
		WorkerEmail:     newUserInput.WorkerEmail,
		WorkerPassword:  newUserInput.WorkerPassword,
	})
	h.mu.Unlock()
}

func (h *EmployerHandlers) HandleCreateEmployer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	newUserInput := new(EmployerInput)
	err := decoder.Decode(newUserInput)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}

	fmt.Println(newUserInput)
	h.mu.Lock()

	var id uint64 = 0
	if len(h.users) > 0 {
		id = h.users[len(h.users)-1].ID + 1
	}

	h.users = append(h.users, Employer{
		ID:                 id,
		EmployerName:       newUserInput.EmployerName,
		EmployerLastName:   newUserInput.EmployerLastName,
		EmployerPosition:   newUserInput.EmployerPosition,
		CompanyName:        newUserInput.CompanyName,
		CompanyDescription: newUserInput.CompanyDescription,
		Website:            newUserInput.Website,
		EmployerEmail:      newUserInput.EmployerEmail,
		EmployerPassword:   newUserInput.EmployerPassword,
	})
	h.mu.Unlock()
}

func main() {

	handlersWorker := WorkerHandlers{
		users: make([]Worker, 0),
		mu:    &sync.Mutex{},
	}

	handlersEmployer := EmployerHandlers{
		users: make([]Employer, 0),
		mu:    &sync.Mutex{},
	}
	//fmt.Println("1")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	})

	http.HandleFunc("/registration/worker", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Println(r.URL.Path)
		//fmt.Println("11", r.Method)
		if r.Method == http.MethodPost {
			fmt.Println("2")
			handlersWorker.HandleCreateWorker(w, r)
			return
		}

	})

	http.HandleFunc("/registration/employer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Println(r.URL.Path)
		fmt.Println("11", r.Method)
		if r.Method == http.MethodPost {
			fmt.Println("3")
			handlersEmployer.HandleCreateEmployer(w, r)
			return
		}

	})

	http.ListenAndServe("0.0.0.0:8080", nil)
}
