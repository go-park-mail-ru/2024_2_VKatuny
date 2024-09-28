package BD

import (
	"sync"
)

type WorkerHandlers struct {
	Users  map[string]Worker
	Mu     *sync.RWMutex
	Amount uint64
}

type EmployerHandlers struct {
	Users  map[string]Employer
	Mu     *sync.RWMutex
	Amount uint64
}

var HandlersWorker = WorkerHandlers{
	Users:  make(map[string]Worker, 0),
	Mu:     &sync.RWMutex{},
	Amount: 0,
}

var HandlersEmployer = EmployerHandlers{
	Users:  make(map[string]Employer, 0),
	Mu:     &sync.RWMutex{},
	Amount: 0,
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
