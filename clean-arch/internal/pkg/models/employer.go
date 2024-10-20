package models

type Employer struct {
	ID       uint32 `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"lastName"`
	Company  string `json:"company"`
	Position string `json:"position"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
