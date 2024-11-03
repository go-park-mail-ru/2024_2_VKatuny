package repository

import (
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// Interface for Worker.
// Now implemented as a in memory db.
// Implementation locates in ./repository
type ApplicantRepository interface {
	Create(worker *dto.ApplicantInput) (*models.Applicant, error)
	GetByID(ID uint64) (*models.Applicant, error)
	GetByEmail(email string) (*models.Applicant, error)
}

var (
	ErrNoUserExist = fmt.Errorf("such user doesn't exist")
)
