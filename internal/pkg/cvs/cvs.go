package cvs

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// Здесь все интерфейсы взаимодействия с CVs

type ICVsRepository interface {
	GetCVsByApplicantID(applicantID uint64) ([]*models.CV, error)
	Create(*dto.JSONCv) (*models.CV, error)
	Get(ID uint64) (*models.CV, error)
	Update(ID uint64, updatedVacancy *dto.JSONCv) (*models.CV, error)
	Delete(ID uint64) error
}

type ICVsUsecase interface {
	GetApplicantCVs(applicantID uint64) ([]*dto.JSONGetApplicantCV, error)
}
