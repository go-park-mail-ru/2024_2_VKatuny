package repository

import "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"

type ICVsRepository interface {
	GetCVsByApplicantID(applicantID uint64) ([]*models.CV, error)
}
