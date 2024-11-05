package models

// Vacancy is a fields of vacancy DB
type Vacancy struct {
	ID          uint64 `json:"id"`
	Position    string `json:"position"`
	Description string `json:"description"`
	Salary      string `json:"salary"`
	EmployerID  uint64 `json:"employerID"` // need id
	WorkType    string `json:"workType"`
	CreatedAt   string `json:"createdAt"`
	Logo        string `json:"logo"`
}
