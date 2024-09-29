package BD

import (
	"sync"
)

const FRONTAPI = "127.0.0.1"

type WorkerHandlers struct {
	Sessions map[string]uint64
	Users    map[string]Worker
	Mu       *sync.RWMutex
	Amount   uint64
}

type EmployerHandlers struct {
	Sessions map[string]uint64
	Users    map[string]Employer
	Mu       *sync.RWMutex
	Amount   uint64
}

var HandlersWorker = WorkerHandlers{
	Sessions: make(map[string]uint64, 0),
	Users:    make(map[string]Worker, 0),
	Mu:       &sync.RWMutex{},
	Amount:   0,
}

var HandlersEmployer = EmployerHandlers{
	Sessions: make(map[string]uint64, 0),
	Users:    make(map[string]Employer, 0),
	Mu:       &sync.RWMutex{},
	Amount:   0,
}

type WorkerInput struct {
	WorkerName      string // `json:"workerName"`
	WorkerLastName  string // `json:"workerLastName"`
	WorkerBirthDate string // `json:"workerBirthDate"`
	WorkerEmail     string // `json:"workerEmail"`
	WorkerPassword  string // `json:"workerPassword"`
}

type Worker struct {
	ID              uint64 //`json:"id"`
	WorkerName      string // `json:"workerFirstName"`
	WorkerLastName  string // `json:"workerLastName"`
	WorkerBirthDate string // `json:"workerBirthDate"`
	WorkerEmail     string // `json:"workerEmail"`
	WorkerPassword  string `json:"-"`
}

type EmployerInput struct {
	EmployerName       string `json:"employerName"`
	EmployerLastName   string `json:"employerLastName"`
	EmployerPosition   string `json:"employerPosition"`
	CompanyName        string `json:"companyName"`
	CompanyDescription string `json:"companyDescription"`
	Website            string `json:"website"`
	EmployerEmail      string `json:"employerEmail"`
	EmployerPassword   string `json:"employerPassword"`
}

type Employer struct {
	ID                 uint64 `json:"employerName"`
	EmployerName       string `json:"employerLastName"`
	EmployerLastName   string `json:"employerPosition"`
	EmployerPosition   string `json:"companyName"`
	CompanyName        string `json:"companyDescription"`
	CompanyDescription string `json:"website"`
	Website            string `json:"employerEmail"`
	EmployerEmail      string `json:"employerPassword"`
	EmployerPassword   string `json:"-"`
}

type UserInput struct {
	TypeUser string `json:"userType"`
	Email    string `json:"login"`
	Password string `json:"password"`
}
