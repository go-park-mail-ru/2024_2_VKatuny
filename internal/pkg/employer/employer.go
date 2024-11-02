// Package employer is a core element of project
package employer

import (
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// Repository is an interface for Employer.
// Now implemented as a in memory db.
// Implementation locates in ./repository
type Repository interface {
	//rename to Add
	// probably shouldn't commit model to Create method
	Create(*dto.JSONEmployerRegistrationForm) (uint64, error)
	GetByID(id uint64) (*models.Employer, error)
	GetByEmail(email string) (*models.Employer, error)
}

var (
	// ErrNoUserExist error means that user doesn't exist
	ErrNoUserExist = fmt.Errorf("user doesnt't exist")
)
