package models

import "time"

type Portfolio struct {
	ID          uint64
	ApplicantID uint64
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
