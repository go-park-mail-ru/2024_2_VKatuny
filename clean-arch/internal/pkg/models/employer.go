package models

type Employer struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	LastName    string `json:"lastName"`
	Position    string `json:"position"`
	CompanyName string `json:"companyName"`
	Description string `json:"companyDescription"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Password    string `json:"-"`
}
