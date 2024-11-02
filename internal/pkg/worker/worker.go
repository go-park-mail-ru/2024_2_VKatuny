package worker

import (
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// Interface for Worker.
// Now implemented as a in memory db.
// Implementation locates in ./repository
type Repository interface {
	Add(worker *models.Worker) (uint64, error)
	GetByID(ID uint64) (*models.Worker, error)
	GetByEmail(email string) (*models.Worker, error)
}

var (
	ErrNoUserExist = fmt.Errorf("such user doesn't exist")
)
