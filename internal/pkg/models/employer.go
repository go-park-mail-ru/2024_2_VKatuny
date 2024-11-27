package models

type Employer struct {
	ID                  uint64
	FirstName           string
	LastName            string
	CityName            string
	Position            string
	CompanyName         string
	CompanyDescription  string
	CompanyWebsite      string
	PathToProfileAvatar string
	CompressedAvatar    string
	Contacts            string
	Email               string
	PasswordHash        string
	CreatedAt           string
	UpdatedAt           string
}
