package models

// model for applicant_session table in DB
type SessionApplicant struct {
	ID          uint64
	ApplicantID uint64
	CookieToken string
	CreatedAt   string
	UpdatedAt   string
}
