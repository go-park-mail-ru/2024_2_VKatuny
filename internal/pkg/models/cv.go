package models

import "time"

type CV struct {
	ID                  uint64
	ApplicantID         uint64
	PositionRus         string
	PositionEng         string
	Description         string
	JobSearchStatus     string
	WorkingExperience   string
	PathToProfileAvatar string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
