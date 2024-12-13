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
	Create(applicant *dto.ApplicantInput) (*models.Applicant, error)
	Update(ID uint64, newApplicantData *dto.JSONUpdateApplicantProfile) (*models.Applicant, error)
	GetByID(ID uint64) (*models.Applicant, error)
	GetByEmail(email string) (*models.Applicant, error)
	GetAllCities(ctx context.Context) ([]*dto.City, error)
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
	GetAllCities(ctx context.Context) ([]*dto.City, error)
}

var (
	ErrNoUserExist = fmt.Errorf("such user doesn't exist")
)
