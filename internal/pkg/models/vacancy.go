package models

// Vacancy is a fields of vacancy DB
type Vacancy struct {
	ID          uint64
	Position    string
	Description string
	Salary      string
	EmployerID  uint64
	LocationID  uint64
	WorkType    string
	Logo        string
	CreatedAt   string
	UpdatedAt   string
}
