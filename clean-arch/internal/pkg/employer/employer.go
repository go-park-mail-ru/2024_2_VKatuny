package employer

import "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"

// Interface for Employer.
// Now implemented as a in memory db.
// Implementation locates in ./repository
type Repository interface {
	//rename to Add
	Create(*models.Employer) (uint64, error)
	GetByID(id uint64) (*models.Employer, error)
}
