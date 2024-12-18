package vacancies

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// Interface for Vacancies.
// Now implemented as a in-memory db.
// Implementation locates in ./repository
type IVacanciesRepository interface {
	Create(ctx context.Context, vacancy *dto.JSONVacancy) (uint64, error)
	SearchAll(ctx context.Context, offset uint64, num uint64, searchStr, group, searchBy string) ([]*dto.JSONVacancy, error)
	GetVacanciesByEmployerID(ctx context.Context, employerID uint64) ([]*dto.JSONVacancy, error)
	GetByID(ctx context.Context, ID uint64) (*dto.JSONVacancy, error)
	Update(ctx context.Context, ID uint64, updatedVacancy *dto.JSONVacancy) (*dto.JSONVacancy, error)
	Delete(ctx context.Context, ID uint64) error
	Subscribe(ctx context.Context, ID uint64, applicantID uint64) error
	GetSubscriptionStatus(ctx context.Context, ID uint64, applicantID uint64) (bool, error)
	GetSubscribersCount(ctx context.Context, ID uint64) (uint64, error)
	GetSubscribersList(ctx context.Context, ID uint64) ([]*models.Applicant, error)
	Unsubscribe(ctx context.Context, ID uint64, applicantID uint64) error
	GetApplicantFavoriteVacancies(ctx context.Context, applicantID uint64) ([]*dto.JSONVacancy, error)
	MakeFavorite(ctx context.Context, ID uint64, applicantID uint64) error
	Unfavorite(ctx context.Context, ID uint64, applicantID uint64) error
}

type IVacanciesUsecase interface {
	GetVacanciesByEmployerID(ctx context.Context, employerID uint64) ([]*dto.JSONGetEmployerVacancy, error)
	CreateVacancy(ctx context.Context, vacancy *dto.JSONVacancy, currentUser *dto.UserFromSession) (*dto.JSONVacancy, error)
	GetVacancy(ctx context.Context, ID uint64) (*dto.JSONVacancy, error)
	ValidateQueryParameters(ctx context.Context, offset, num string) (uint64, uint64, error)
	UpdateVacancy(ctx context.Context, ID uint64, updatedVacancy *dto.JSONVacancy, currentUser *dto.UserFromSession) (*dto.JSONVacancy, error)
	DeleteVacancy(ctx context.Context, ID uint64, currentUser *dto.UserFromSession) error
	SubscribeOnVacancy(ctx context.Context, ID uint64, currentUser *dto.UserFromSession) error
	UnsubscribeFromVacancy(ctx context.Context, ID uint64, currentUser *dto.UserFromSession) error
	GetSubscriptionInfo(ctx context.Context, ID uint64, applicantID uint64) (*dto.JSONVacancySubscriptionStatus, error)
	SearchVacancies(ctx context.Context, offsetStr, numStr, searchStr, group, searchBy string) ([]*dto.JSONVacancy, error)
	GetVacancySubscribers(ctx context.Context, ID uint64, currentUser *dto.UserFromSession) (*dto.JSONVacancySubscribers, error)
	GetApplicantFavoriteVacancies(ctx context.Context, applicantID uint64) ([]*dto.JSONGetEmployerVacancy, error)
	AddIntoFavorite(ctx context.Context, ID uint64, currentUser *dto.UserFromSession) error
	Unfavorite(ctx context.Context, ID uint64, currentUser *dto.UserFromSession) error
}
