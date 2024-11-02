package models

type Employer struct {
	ID                 uint64
	FirstName          string
	LastName           string
	Position           string
	Company            string
	CompanyDescription string
	CompanyWebsite     string
	Email              string
	Password           string
}
