// Package inmemorydb is for inmemory db and structures
package inmemorydb

import (
	"sync"
)

// IP of our server
const IP = "89.208.199.175"

// FRONTENDIP is our fronted server ip, later it was const IP = "127.0.0.1"
const FRONTENDIP = "http://" + IP

// BACKENDIP is our backend server ip
const BACKENDIP = IP + ":8080"

// WorkerHandlers struct for worker table
type WorkerHandlers struct {
	// key   -    (Cookie.Value)
	// value - ID (Worker.ID, Employer.ID)
	Sessions map[string]uint64
	Users    map[string]Worker
	Mu       *sync.RWMutex
	Amount   uint64
}

// EmployerHandlers struct for employer table
type EmployerHandlers struct {
	Sessions map[string]uint64
	Users    map[string]Employer
	Mu       *sync.RWMutex
	Amount   uint64
}

// HandlersWorker simulates worker db
var HandlersWorker = WorkerHandlers{
	Sessions: make(map[string]uint64, 0),
	Users:    make(map[string]Worker, 0),
	Mu:       &sync.RWMutex{},
	Amount:   0,
}

// HandlersEmployer simulates employer db
var HandlersEmployer = EmployerHandlers{
	Sessions: make(map[string]uint64, 0),
	Users:    make(map[string]Employer, 0),
	Mu:       &sync.RWMutex{},
	Amount:   0,
}

// WorkerInput data getting from registration worker
type WorkerInput struct {
	WorkerName      string `json:"workerName"`
	WorkerLastName  string `json:"workerLastName"`
	WorkerBirthDate string `json:"workerBirthDate"`
	WorkerEmail     string `json:"workerEmail"`
	WorkerPassword  string `json:"workerPassword"`
}

// Worker is columns of worker db
type Worker struct {
	ID              uint64 `json:"id"`
	WorkerName      string `json:"workerFirstName"`
	WorkerLastName  string `json:"workerLastName"`
	WorkerBirthDate string `json:"workerBirthDate"`
	WorkerEmail     string `json:"workerEmail"`
	WorkerPassword  string `json:"-"`
}

// EmployerInput data getting from registration employer
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

// Employer is columns of employer db
type Employer struct {
	ID                 uint64 `json:"id"`
	EmployerName       string `json:"employerName"`
	EmployerLastName   string `json:"employerLastName"`
	EmployerPosition   string `json:"employerPosition"`
	CompanyName        string `json:"companyName"`
	CompanyDescription string `json:"companyDescription"`
	Website            string `json:"website"`
	EmployerEmail      string `json:"employerEmail"`
	EmployerPassword   string `json:"-"`
}

// UserInput data getting from login
type UserInput struct {
	TypeUser string `json:"userType"`
	Email    string `json:"login"`
	Password string `json:"password"`
}

// VacanciesHandler struct for vacancy table
type VacanciesHandler struct {
	Vacancy []Vacancy
	Count   uint64
	Mutex   *sync.RWMutex
}

// Vacancy is columns of vacancy db
type Vacancy struct {
	ID          uint64 `json:"id"`
	Position    string `json:"position"`
	Description string `json:"description"`
	Salary      string `json:"salary"`
	Employer    string `json:"employer"`
	Location    string `json:"location"`
	CreatedAt   string `json:"createdAt"`
	Logo        string `json:"logo"`
}

// Vacancies simulate vacancies db
var Vacancies = VacanciesHandler{
	Vacancy: make([]Vacancy, 0),
	Count:   0,
	Mutex:   &sync.RWMutex{},
}

// type userType string

const (
	// WORKER global name
	WORKER = "applicant" //userType("worker")
	// EMPLOYER global name
	EMPLOYER = "employer" //userType("employer")
)

// UserAlreadyExist err struct for existing user
type UserAlreadyExist struct {
	UserAlreadyExist bool `json:"userAlreadyExist"`
}

// ErrorMessages err struct
type ErrorMessages struct {
	Status  int    `json:"status"`
	ErrText string `json:"errText"`
}

// AuthorizedUserFields struct for response of handler after login
type AuthorizedUserFields struct {
	ID         uint64 `json:"id"`
	TypeOfUser string //userType `json:"typeOfUser"`
}

// ReturnUserFields struct for response of handler after registration
type ReturnUserFields struct {
	StatusCode int                  `json:"statusCode"`
	User       AuthorizedUserFields `json:"user"`
}

// MakeVacancies fill db with test vacancies
func MakeVacancies() {
	Vacancies.Count = 25
	for i := uint64(0); i < 25; i += 5 {
		Vacancies.Vacancy = append(Vacancies.Vacancy, Vacancy{
			ID:       i,
			Position: "Продавец консультант",
			Description: `Ищем продавца на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на тяжелую работу. Своевременную оплату гарантируем.`,
			Salary:    "Не указана",
			Employer:  "X-Retail Group",
			Location:  "Moscow",
			CreatedAt: "2024.09.29 16:55:00", // YYYY.MM.DD HH:MM:SS
			Logo:      "img/picture_name1.png",
		})
		Vacancies.Vacancy = append(Vacancies.Vacancy, Vacancy{
			ID:       i + 1,
			Position: "Продавец",
			Description: `Ищем продавца на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на тяжелую работу. Своевременную оплату гарантируем.`,
			Salary:    "80 000",
			Employer:  "X-Retail Group",
			Location:  "Moscow",
			CreatedAt: "2024.09.29 17:55:00", // YYYY.MM.DD HH:MM:SS
			Logo:      "img/picture_name2.png",
		})
		Vacancies.Vacancy = append(Vacancies.Vacancy, Vacancy{
			ID:       i + 2,
			Position: "Администратор",
			Description: `Ищем администратора на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на продуктивную работу с людьми. Своевременную оплату гарантируем.`,
			Salary:    "100 500",
			Employer:  "X-Retail Group",
			Location:  "Moscow",
			CreatedAt: "2024.09.29 18:55:00", // YYYY.MM.DD HH:MM:SS
			Logo:      "img/picture_name3.png",
		})
		Vacancies.Vacancy = append(Vacancies.Vacancy, Vacancy{
			ID:       i + 3,
			Position: "Охранник",
			Description: `Ищем охранника на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на тяжелую посменную работу. Своевременную оплату гарантируем.`,
			Salary:    "Не указана",
			Employer:  "X-Retail Group",
			Location:  "Moscow",
			CreatedAt: "2024.09.29 19:55:00", // YYYY.MM.DD HH:MM:SS
			Logo:      "img/picture_name4.png",
		})
		Vacancies.Vacancy = append(Vacancies.Vacancy, Vacancy{
			ID:       i + 4,
			Position: "Уборщик помещений",
			Description: `Ищем уборщика на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на тяжелую кропотливую работу. Своевременную оплату гарантируем.`,
			Salary:    "50 000",
			Employer:  "X-Retail Group",
			Location:  "Moscow",
			CreatedAt: "2024.09.29 20:55:00", // YYYY.MM.DD HH:MM:SS
			Logo:      "img/picture_name5.png",
		})
	}

}

// MakeUsers fill db with test users
func MakeUsers() {
	HandlersWorker.Users["a@mail.ru"] = Worker{
		ID:              1,
		WorkerName:      "Vasia",
		WorkerLastName:  "Vasion",
		WorkerBirthDate: "12-12-2012",
		WorkerEmail:     "a@mail.ru",
		WorkerPassword:  "$2a$10$nOPg8rvfOOSNrYv.zfzz7eKVDJvcHhXXGciR/SuHTekTOYTjZr4oa", // pass1234
	}
	HandlersEmployer.Users["b@mail.ru"] = Employer{
		ID:                 1,
		EmployerName:       "Ilia",
		EmployerLastName:   "Ilin",
		EmployerPosition:   "CEO",
		CompanyName:        "Ilia Ilin Enterprices",
		CompanyDescription: "Ilia Ilin best company",
		Website:            "Ilin.com",
		EmployerEmail:      "b@mail.ru",
		EmployerPassword:   "$2a$10$aw9A84PCPKXvvUMD4eQtyulfXNnlhN3.Wts7PF9xiuJWJd0bV3o9i", // pass4321
	}

}
