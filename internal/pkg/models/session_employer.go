package models

import "time"

// model for employer_session table in DB
type SessionEmployer struct {
	ID          uint64
	EmployerID uint64
	CookieToken string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
