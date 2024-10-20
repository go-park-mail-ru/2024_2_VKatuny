package employer

import "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"

// Interface for Employer.
// Now implemented as a in memory db.
// Implementation locates in ./repository
type Repository interface {
	Create(*models.Employer) (uint32, error)
	GetByID(id uint32) (*models.Employer, error)
}
