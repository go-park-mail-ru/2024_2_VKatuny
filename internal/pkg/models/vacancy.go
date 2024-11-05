package models

// Vacancy is a fields of vacancy DB
type Vacancy struct {
	ID          uint64 `json:"id"`
	Position    string `json:"position"`
	Description string `json:"description"`
	Salary      string `json:"salary"`
	Employer    string `json:"employer"` // need id
	Location    string `json:"location"`
	CreatedAt   string `json:"createdAt"`
	Logo        string `json:"logo"`
}
