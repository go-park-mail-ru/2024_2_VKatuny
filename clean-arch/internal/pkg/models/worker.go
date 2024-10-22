package models

// Worker is columns of worker db
type Worker struct {
	ID        uint64 `json:"id"`
	Name      string `json:"FirstName"`
	LastName  string `json:"LastName"`
	BirthDate string `json:"BirthDate"`
	Email     string `json:"Email"`
	Password  string `json:"-"`
}
