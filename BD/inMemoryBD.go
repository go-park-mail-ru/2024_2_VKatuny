package BD

import (
	"sync"
)

const FRONTAPI = "http://127.0.0.1:8000"

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


type VacanciesHandler struct {
	Vacancy []Vacancy
	Count   uint64
	Mutex   *sync.RWMutex 
}

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

var Vacancies = VacanciesHandler{
	Vacancy: make([]Vacancy, 0),
	Count: 0,
	Mutex: &sync.RWMutex{},
}

func MakeVacancies() {
	Vacancies.Count = 25
	for i := uint64(0); i < 25; i += 5 {
		Vacancies.Vacancy = append(Vacancies.Vacancy, Vacancy{
			ID: i,
			Position: "Продавец консультант",
			Description: `Ищем продавца на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на тяжелую работу. Своевременную оплату гарантируем.`,
			Salary: "Не указана",
			Employer: "X-Retail Group",
			Location: "Moscow",
			CreatedAt: "2024.09.29 16:55:00",  // YYYY.MM.DD HH:MM:SS
			Logo: "/img/picture_name.png",
		})
		Vacancies.Vacancy = append(Vacancies.Vacancy, Vacancy{
			ID: i + 1,
			Position: "Продавец",
			Description: `Ищем продавца на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на тяжелую работу. Своевременную оплату гарантируем.`,
			Salary: "80 000",
			Employer: "X-Retail Group",
			Location: "Moscow",
			CreatedAt: "2024.09.29 17:55:00",  // YYYY.MM.DD HH:MM:SS
			Logo: "/img/picture_name.png",
		})
		Vacancies.Vacancy = append(Vacancies.Vacancy, Vacancy{
			ID: i + 2,
			Position: "Администратор",
			Description: `Ищем продавца на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на тяжелую работу. Своевременную оплату гарантируем.`,
			Salary: "100 500",
			Employer: "X-Retail Group",
			Location: "Moscow",
			CreatedAt: "2024.09.29 18:55:00",  // YYYY.MM.DD HH:MM:SS
			Logo: "/img/picture_name.png",
		})
		Vacancies.Vacancy = append(Vacancies.Vacancy, Vacancy{
			ID: i + 3,
			Position: "Охранник",
			Description: `Ищем продавца на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на тяжелую работу. Своевременную оплату гарантируем.`,
			Salary: "Не указана",
			Employer: "X-Retail Group",
			Location: "Moscow",
			CreatedAt: "2024.09.29 19:55:00",  // YYYY.MM.DD HH:MM:SS
			Logo: "/img/picture_name.png",
		})
		Vacancies.Vacancy = append(Vacancies.Vacancy, Vacancy{
			ID: i + 4,
			Position: "Уборщик помещений",
			Description: `Ищем продавца на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на тяжелую работу. Своевременную оплату гарантируем.`,
			Salary: "50 000",
			Employer: "X-Retail Group",
			Location: "Moscow",
			CreatedAt: "2024.09.29 20:55:00",  // YYYY.MM.DD HH:MM:SS
			Logo: "/img/picture_name.png",
		})
	}
	
}


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
