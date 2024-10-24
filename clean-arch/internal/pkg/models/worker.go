package models

// Worker is columns of worker db
type Worker struct {
	ID        uint64 `json:"id"`
	Name      string `json:"firstName"`
	LastName  string `json:"lastName"`
	BirthDate string `json:"birthDate"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
