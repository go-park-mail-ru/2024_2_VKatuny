package vacancies

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// Interface for Vacancies.
// Now implemented as a in-memory db.
// Implementation locates in ./repository
type IVacanciesRepository interface { // TODO: rename to IVacanciesRepository
	Create(vacancy *dto.JSONVacancy) (uint64, error) // TODO: should accept DTO not a model
	GetWithOffset(offset uint64, num uint64) ([]*dto.JSONVacancy, error)
	GetVacanciesByEmployerID(employerID uint64) ([]*dto.JSONVacancy, error)
	GetByID(ID uint64) (*dto.JSONVacancy, error)
	Update(ID uint64, updatedVacancy *dto.JSONVacancy) (*dto.JSONVacancy, error)
	Delete(ID uint64) error
	Subscribe(ID uint64, applicantID uint64) error
	GetSubscriptionStatus(ID uint64, applicantID uint64) (bool, error)
	GetSubscribersCount(ID uint64) (uint64, error)
	GetSubscribersList(ID uint64) ([]*models.Applicant, error)
	Unsubscribe(ID uint64, applicantID uint64) error
}

type IVacanciesUsecase interface {
	GetVacanciesByEmployerID(employerID uint64) ([]*dto.JSONGetEmployerVacancy, error)
	CreateVacancy(vacancy *dto.JSONVacancy, currentUser *dto.SessionUser) (*dto.JSONVacancy, error)
	GetVacancy(ID uint64) (*dto.JSONVacancy, error)
	UpdateVacancy(ID uint64, updatedVacancy *dto.JSONVacancy, currentUser *dto.SessionUser) (*dto.JSONVacancy, error)
	DeleteVacancy(ID uint64, currentUser *dto.SessionUser) error
	SubscribeOnVacancy(ID uint64, currentUser *dto.SessionUser) error
	UnsubscribeFromVacancy(ID uint64, currentUser *dto.SessionUser) error
	GetSubscriptionInfo(ID uint64, applicantID uint64) (*dto.JSONVacancySubscriptionStatus, error)
	GetVacancySubscribers(ID uint64, currentUser *dto.SessionUser) (*dto.JSONVacancySubscribers, error)
}
