package employer

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"
)

// Interface for Employer.
// Now implemented as a in memory db.
// Implementation locates in ./repository
type Repository interface {
	//rename to Add
	// probably shouldn't commit model to Create method
	Create(*models.Employer) (uint64, error)
	GetByID(id uint64) (*models.Employer, error)
	GetByEmail(email string) (*models.Employer, error)
}

var (
	ErrNoUserExist = fmt.Errorf("user doesnt't exist")
)
