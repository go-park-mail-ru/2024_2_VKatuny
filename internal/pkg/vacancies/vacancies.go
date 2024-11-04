package vacancies

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// Interface for Vacancies.
// Now implemented as a in-memory db.
// Implementation locates in ./repository
type Repository interface {   // TODO: rename to IVacanciesRepository
	Add(vacancy *models.Vacancy) (uint64, error)  // TODO: should accept DTO not a model
	GetWithOffset(offset uint64, num uint64) ([]*models.Vacancy, error)
	GetVacanciesByEmployerID(employerID uint64) ([]*models.Vacancy, error)
}
