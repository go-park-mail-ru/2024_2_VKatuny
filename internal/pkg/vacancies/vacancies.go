package vacancies

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// Interface for Vacancies.
// Now implemented as a in-memory db.
// Implementation locates in ./repository
type IVacanciesRepository interface { // TODO: rename to IVacanciesRepository
	Create(vacancy *dto.JSONVacancy) (uint64, error) // TODO: should accept DTO not a model
	GetWithOffset(offset uint64, num uint64) ([]*models.Vacancy, error)
	GetVacanciesByEmployerID(employerID uint64) ([]*models.Vacancy, error)
	GetByID(ID uint64) (*models.Vacancy, error)
	Update(ID uint64, updatedVacancy *dto.JSONVacancy) (*models.Vacancy, error)
	Delete(ID uint64) error
	Subscribe(ID uint64, applicantID uint64) error
	GetSubscriptionStatus(ID uint64, applicantID uint64) (bool, error)
	GetScribersCount(ID uint64) (uint64, error)
	GetSubscribersList(ID uint64) ([]uint64, error)  // TODO: ask Ilya
	Unsubscribe(ID uint64, applicantID uint64) error
}

type IVacanciesUsecase interface {
	GetVacanciesByEmployerID(employerID uint64) ([]*dto.JSONGetEmployerVacancy, error)
}
