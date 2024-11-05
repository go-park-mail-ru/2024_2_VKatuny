package cvs

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// Здесь все интерфейсы взаимодействия с CVs

type ICVsRepository interface {
	GetCVsByApplicantID(applicantID uint64) ([]*models.CV, error)
	Create(cv *dto.JSONCv) (*models.CV, error)
	GetByID(ID uint64) (*models.CV, error)
	Update(ID uint64, updatedVacancy *dto.JSONCv) (*models.CV, error)
	Delete(ID uint64) error
}

type ICVsUsecase interface {
	GetApplicantCVs(applicantID uint64) ([]*dto.JSONGetApplicantCV, error)
	CreateCV(cv *dto.JSONCv) (*dto.JSONCv, error)
	GetCV(ID uint64) (*dto.JSONCv, error)
	UpdateCV(ID uint64, sessionID string, cv *dto.JSONCv) (*dto.JSONCv, error)
	DeleteCV(ID uint64, sessionID string) error
}
