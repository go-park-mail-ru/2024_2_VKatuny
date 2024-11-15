package employer

import (
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// EmployerRepository is an interface for Employer.
// Now implemented as a in memory db.
// Implementation locates in ./repository
type IEmployerRepository interface {
	//rename to Add
	// probably shouldn't commit model to Create method
	Create(*dto.EmployerInput) (*models.Employer, error)
	Update(ID uint64, newEmployerData *dto.JSONUpdateEmployerProfile) error
	GetByID(id uint64) (*models.Employer, error)
	GetByEmail(email string) (*models.Employer, error)
}

type IEmployerUsecase interface {
	CreateEmployer(form *dto.JSONEmployerRegistrationForm) (*dto.JSONUser, error)
	GetByID(id uint64) (*models.Employer, error) 
	GetEmployerProfile(employerID uint64) (*dto.JSONGetEmployerProfile, error)
	UpdateEmployerProfile(employerID uint64, employerProfile *dto.JSONUpdateEmployerProfile) error
}

var (
	// ErrNoUserExist error means that user doesn't exist
	ErrNoUserExist = fmt.Errorf("user doesn't exist")
)
