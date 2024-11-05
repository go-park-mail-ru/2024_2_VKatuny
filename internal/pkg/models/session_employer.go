package models

// model for employer_session table in DB
type SessionEmployer struct {
	ID          uint64
	EmployerID  uint64
	CookieToken string
	CreatedAt   string
	UpdatedAt   string
}
