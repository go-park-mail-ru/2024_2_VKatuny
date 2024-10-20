package worker

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"
)

// Interface for Worker.
// Now implemented as a in memory db.
// Implementation locates in ./repository
type Repository interface {
	Add(worker *models.Worker) (uint64, error)
	GetByID(ID uint64) (*models.Worker, error)
}
