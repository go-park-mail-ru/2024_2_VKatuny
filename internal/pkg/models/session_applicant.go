package models

import "time"

// model for applicant_session table in DB
type SessionApplicant struct {
	ID          uint64
	ApplicantID uint64
	CookieToken string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
