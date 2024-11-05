package models

type CV struct {
	ID                  uint64
	ApplicantID         uint64
	PositionRus         string
	PositionEng         string
	JobSearchStatus     string
	Description         string
	WorkingExperience   string
	PathToProfileAvatar string
	CreatedAt           string
	UpdatedAt           string
}
