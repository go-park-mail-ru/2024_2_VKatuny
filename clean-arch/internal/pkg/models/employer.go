package models

type Employer struct {
	ID          uint64 `json:"id"`
	Name        string `json:"employerName"`
	LastName    string `json:"employerLastName"`
	Position    string `json:"employerPosition"`
	CompanyName string `json:"companyName"`
	Description string `json:"companyDescription"`
	Website     string `json:"website"`
	Email       string `json:"employerEmail"`
	Password    string `json:"-"`
}
