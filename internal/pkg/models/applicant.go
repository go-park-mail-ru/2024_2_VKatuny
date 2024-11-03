package models

// Worker is columns of worker db
// type Worker struct {
// 	ID        uint64 `json:"id"`
// 	Name      string `json:"firstName"`
// 	LastName  string `json:"lastName"`
// 	BirthDate string `json:"birthDate"`
// 	Email     string `json:"email"`
// 	Password  string `json:"password"`
// }

type Applicant struct {
	ID                  uint64 `json:"id"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	CityName            string `json:"cityName"`
	BirthDate           string `json:"birthDate"`
	PathToProfileAvatar string `json:"pathToProfileAvatar"`
	Contacts            string `json:"contacts"`
	Education           string `json:"education"`
	Email               string `json:"email"`
	PasswordHash        string `json:"passwordHash"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
}
