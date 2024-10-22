package models

type Employer struct {
	ID          uint64 `json:"id"`
	Name        string `json:"Name"`
	LastName    string `json:"LastName"`
	Position    string `json:"Position"`
	CompanyName string `json:"companyName"`
	Description string `json:"companyDescription"`
	Website     string `json:"website"`
	Email       string `json:"Email"`
	Password    string `json:"-"`
}
