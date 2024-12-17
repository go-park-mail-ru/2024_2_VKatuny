package applicant

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// Interface for Worker.
// Now implemented as a in memory db.
// Implementation locates in ./repository
type IApplicantRepository interface { // TODO: rename to IApplicantRepository
	// Can we send dto to Repository?
	Create(ctx context.Context, applicant *dto.ApplicantInput) (*models.Applicant, error)
	Update(ctx context.Context, ID uint64, newApplicantData *dto.JSONUpdateApplicantProfile) (*models.Applicant, error)
	GetByID(ctx context.Context, ID uint64) (*models.Applicant, error)
	GetByEmail(ctx context.Context, email string) (*models.Applicant, error)
	GetAllCities(ctx context.Context, namePat string) ([]string, error)
}

type IApplicantUsecase interface {
	Create(ctx context.Context, applicant *dto.JSONApplicantRegistrationForm) (*dto.JSONUser, error)
	GetByID(ctx context.Context, ID uint64) (*dto.JSONApplicantOutput, error)
	GetApplicantProfile(ctx context.Context, userID uint64) (*dto.JSONGetApplicantProfile, error)
	UpdateApplicantProfile(
		ctx context.Context,
		applicantID uint64,
		newProfileData *dto.JSONUpdateApplicantProfile,
	) error
	GetAllCities(ctx context.Context, namePart string) ([]string, error)
}

var (
	ErrNoUserExist = fmt.Errorf("such user doesn't exist")
)
