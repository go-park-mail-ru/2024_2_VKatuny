package models

import "time"

type CV struct {
	ID                  uint64
	ApplicantID         uint64
	PositionRus         string
	PositionEng         string
	JobSearchStatusID   int
	WorkingExperience   string  
	PathToProfileAvatar string // Для Олега на русском. оно нам надо? Разве там не аватарка пользователя?
	CreatedAt           time.Time  // И еще нужно описание вакансии
	UpdatedAt           time.Time
}
