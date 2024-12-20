package employer

import (
	"context"
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
	Create(ctx context.Context, employer *dto.EmployerInput) (*models.Employer, error)
	Update(ctx context.Context, ID uint64, newEmployerData *dto.JSONUpdateEmployerProfile) (*models.Employer, error)
	GetByID(ctx context.Context, id uint64) (*models.Employer, error)
	GetByEmail(ctx context.Context, email string) (*models.Employer, error)
}

type IEmployerUsecase interface {
	Create(ctx context.Context, form *dto.JSONEmployerRegistrationForm) (*dto.JSONUser, error)
	GetByID(ctx context.Context, id uint64) (*dto.JSONEmployer, error) 
	GetEmployerProfile(ctx context.Context, employerID uint64) (*dto.JSONGetEmployerProfile, error)
	UpdateEmployerProfile(ctx context.Context, employerID uint64, employerProfile *dto.JSONUpdateEmployerProfile) error
}

var (
	// ErrNoUserExist error means that user doesn't exist
	ErrNoUserExist = fmt.Errorf("user doesn't exist")
)
