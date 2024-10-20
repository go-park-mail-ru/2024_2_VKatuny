package models

// Worker is columns of worker db
type Worker struct {
	ID        uint64 `json:"id"`
	Name      string `json:"workerFirstName"`
	LastName  string `json:"workerLastName"`
	BirthDate string `json:"workerBirthDate"`
	Email     string `json:"workerEmail"`
	Password  string `json:"-"`
}
