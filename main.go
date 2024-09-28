package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var HandlersWorker = WorkerHandlers{
	users:  make(map[string]Worker, 0),
	mu:     &sync.RWMutex{},
	amount: 0,
}

var HandlersEmployer = EmployerHandlers{
	users:  make(map[string]Employer, 0),
	mu:     &sync.RWMutex{},
	amount: 0,
}

type WorkerHandlers struct {
	users  map[string]Worker
	mu     *sync.RWMutex
	amount uint64
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
	users  map[string]Employer
	mu     *sync.RWMutex
	amount uint64
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
	decErr := decoder.Decode(newUserInput)
	if decErr != nil {
		log.Printf("error while unmarshalling JSON: %s", decErr)
		w.Write([]byte("{}"))
		return
	}

	fmt.Println(newUserInput)

	_, rErr := GetWorkerByEmail(h, newUserInput.WorkerEmail)

	if rErr != nil {
		h.mu.Lock()
		var id uint64 = h.amount + 1
		h.users[newUserInput.WorkerEmail] = Worker{
			ID:              id,
			WorkerName:      newUserInput.WorkerName,
			WorkerLastName:  newUserInput.WorkerLastName,
			WorkerBirthDate: newUserInput.WorkerBirthDate,
			WorkerEmail:     newUserInput.WorkerEmail,
			WorkerPassword:  HashPassword(newUserInput.WorkerPassword),
		}
		h.mu.Unlock()
	} else {
		log.Printf("error user with this email already exists: %s", newUserInput.WorkerEmail)
		w.Write([]byte("{}"))
		return
	}
}

func GetWorkerByEmail(table *WorkerHandlers, email string) (Worker, error) {
	table.mu.RLock()
	user, err := table.users[email]
	table.mu.RUnlock()
	if err == true {
		return user, nil
	} else {
		return Worker{}, fmt.Errorf("No user with such email")
	}
}

func (h *EmployerHandlers) HandleCreateEmployer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	newUserInput := new(EmployerInput)
	err := decoder.Decode(newUserInput)
	if err != nil {
		w.WriteHeader(400)
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}

	fmt.Println(newUserInput)
	h.mu.RLock()
	_, alreadyExist := h.users[newUserInput.EmployerEmail]
	h.mu.RUnlock()
	if alreadyExist == false {
		h.mu.Lock()

		var id uint64 = h.amount + 1
		h.users[newUserInput.EmployerEmail] = Employer{
			ID:                 id,
			EmployerName:       newUserInput.EmployerName,
			EmployerLastName:   newUserInput.EmployerLastName,
			EmployerPosition:   newUserInput.EmployerPosition,
			CompanyName:        newUserInput.CompanyName,
			CompanyDescription: newUserInput.CompanyDescription,
			Website:            newUserInput.Website,
			EmployerEmail:      newUserInput.EmployerEmail,
			EmployerPassword:   newUserInput.EmployerPassword,
		}
		h.mu.Unlock()
	} else {
		log.Printf("error user with this email already exists: %s", newUserInput.EmployerEmail)
		w.Write([]byte("{}"))
		return
	}
}

func HashPassword(password string) string {
	bytePassword := []byte(password)
	cost := 10
	hashedPassword, _ := bcrypt.GenerateFromPassword(bytePassword, cost)
	fmt.Println(string(password[:]), string(hashedPassword[:]))
	return string(hashedPassword[:])
}

func EqualHashedPasswords(passwordBD string, passwordFront string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordBD), []byte(passwordFront))
	if err == nil {
		return true
	} else {
		return false
	}
}

func main() {
	fmt.Println(EqualHashedPasswords(HashPassword("pass"), "pass"))
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	})

	http.HandleFunc("/api/registration/worker", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Println(r.URL.Path)
		//fmt.Println("11", r.Method)
		if r.Method == http.MethodPost {
			fmt.Println("2")
			HandlersWorker.HandleCreateWorker(w, r)
			return
		}

	})

	http.HandleFunc("/api/registration/employer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Println(r.URL.Path)
		fmt.Println("11", r.Method)
		if r.Method == http.MethodPost {
			fmt.Println("3")
			HandlersEmployer.HandleCreateEmployer(w, r)
			return
		}

	})
	http.ListenAndServe("0.0.0.0:8080", nil)
}
