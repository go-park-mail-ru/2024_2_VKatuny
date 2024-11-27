package models

type Applicant struct {
	ID                  uint64
	FirstName           string
	LastName            string
	CityName            string
	BirthDate           string
	PathToProfileAvatar string
	CompressedAvatar    string
	Contacts            string
	Education           string
	Email               string
	PasswordHash        string
	CreatedAt           string
	UpdatedAt           string
}
