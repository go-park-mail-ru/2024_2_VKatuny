package vacancies

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"
)

// Interface for Vacancies.
// Now implemented as a in memory db.
// Implementation locates in ./repository
type Repository interface {
	Add(vacancy *models.Vacancy) (uint64, error)
	GetWithOffset(offset uint64, num uint64) ([]*models.Vacancy, error)
}
