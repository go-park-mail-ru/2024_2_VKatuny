package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/usecase"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"



)

func TestGet(t *testing.T) {
	t.Parallel()
	type repo struct {
		applicant *mock.MockIApplicantRepository
	}
	tests := []struct {
		name    string
		profile *dto.JSONGetApplicantProfile
		prepare func(
			repo *repo,
			profile *dto.JSONGetApplicantProfile,
		) (uint64, *dto.JSONGetApplicantProfile)
	}{
		{
			name:    "Create: bad repository",
			profile: nil,
			prepare: func(
				repo *repo,
				profile *dto.JSONGetApplicantProfile) (uint64, *dto.JSONGetApplicantProfile) {
				userID := uint64(1)
				repo.applicant.
					EXPECT().
					GetByID(userID).
					Return(nil, errors.New("bad repository"))
				return userID, profile
			},
		},
		{
			name:    "Create: ok",
			profile: new(dto.JSONGetApplicantProfile),
			prepare: func(
				repo *repo,
				profile *dto.JSONGetApplicantProfile) (uint64, *dto.JSONGetApplicantProfile) {
				userID := uint64(1)
				model := &models.Applicant{
					ID:        userID,
					FirstName: "Ivan",
					LastName:  "Ivanov",
					CityName:  "Moscow",
				}

				repo.applicant.
					EXPECT().
					GetByID(userID).
					Return(model, nil)
				rprofile := &dto.JSONGetApplicantProfile{
					ID:        model.ID,
					FirstName: model.FirstName,
					LastName:  model.LastName,
					City:      model.CityName,
				}
				return userID, rprofile
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			applicant: mock.NewMockIApplicantRepository(ctrl),
		}
		var userID uint64
		userID, tt.profile = tt.prepare(repo, tt.profile)

		repositories := &internal.Repositories{
			ApplicantRepository: repo.applicant,
		}
		uc := usecase.NewApplicantUsecase(logrus.New(), repositories)
		profile, _ := uc.GetApplicantProfile(context.Background(), userID)

		require.Equal(t, tt.profile, profile)
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	type repo struct {
		applicant *mock.MockIApplicantRepository
	}
	tests := []struct {
		name    string
		profile *dto.JSONUpdateApplicantProfile
		expected error
		prepare func(
			repo *repo,
			profile *dto.JSONUpdateApplicantProfile,
		) (uint64, *dto.JSONUpdateApplicantProfile)
	}{
		{
			name:    "Create: bad repository",
			profile: nil,
			expected: errors.New("bad repository"),
			prepare: func(
				repo *repo,
				profile *dto.JSONUpdateApplicantProfile) (uint64, *dto.JSONUpdateApplicantProfile) {
				userID := uint64(1)
				repo.applicant.
					EXPECT().
					Update(userID, gomock.Any()).
					Return(nil, errors.New("bad repository"))
				return userID, profile
			},
		},
		{
			name:    "Create: ok",
			profile: new(dto.JSONUpdateApplicantProfile),
			expected: nil,
			prepare: func(
				repo *repo,
				profile *dto.JSONUpdateApplicantProfile) (uint64, *dto.JSONUpdateApplicantProfile) {
				userID := uint64(1)
				model := &models.Applicant{
					ID:        userID,
					FirstName: "Ivan",
					LastName:  "Ivanov",
					CityName:  "Moscow",
				}

				repo.applicant.
					EXPECT().
					Update(userID, gomock.Any()).
					Return(nil, nil)
				rprofile := &dto.JSONUpdateApplicantProfile{
					FirstName: model.FirstName,
					LastName:  model.LastName,
					City:      model.CityName,
				}
				return userID, rprofile
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			applicant: mock.NewMockIApplicantRepository(ctrl),
		}
		var userID uint64
		userID, tt.profile = tt.prepare(repo, tt.profile)

		repositories := &internal.Repositories{
			ApplicantRepository: repo.applicant,
		}
		uc := usecase.NewApplicantUsecase(logrus.New(), repositories)
		err := uc.UpdateApplicantProfile(context.Background(), userID, tt.profile)

		require.Equal(t, tt.expected ,err)
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()
	type repo struct {
		applicant *mock.MockIApplicantRepository
	}
	tests := []struct {
		name    string
		form *dto.JSONApplicantRegistrationForm
		user *dto.JSONUser
		prepare func(
			repo *repo,
			form *dto.JSONApplicantRegistrationForm,
			user *dto.JSONUser,
		) (*dto.JSONApplicantRegistrationForm, *dto.JSONUser)
	}{
		{
			name:    "Create: bad repository get",
			form: new(dto.JSONApplicantRegistrationForm),
			user: new(dto.JSONUser),
			prepare: func(
				repo *repo,
				form *dto.JSONApplicantRegistrationForm,
				user *dto.JSONUser,
			) (*dto.JSONApplicantRegistrationForm, *dto.JSONUser) {
				repo.applicant.
					EXPECT().
					GetByEmail(form.Email).
					Return(nil, errors.New("sql: no rows in result set"))
				return form, nil
			},
		},
		{
			name:    "Create: bad repository create",
			form: new(dto.JSONApplicantRegistrationForm),
			user: new(dto.JSONUser),
			prepare: func(
				repo *repo,
				form *dto.JSONApplicantRegistrationForm,
				user *dto.JSONUser,
			) (*dto.JSONApplicantRegistrationForm, *dto.JSONUser) {
				form = &dto.JSONApplicantRegistrationForm{
					FirstName: "Ivan",
					LastName:  "Ivanov",
					BirthDate: "2000-01-01",
					Email:     "ivanov@ya.ru",
					Password:  "123456",
				}
				applicant := &models.Applicant{
					ID:        1,
				}
				repo.applicant.
					EXPECT().
					GetByEmail(form.Email).
					Return(nil, nil)
				repo.applicant.
					EXPECT().
					Create(gomock.Any()).
					Return(applicant, errors.New("bad repository"))
				return form, nil
			},
		},
		{
			name:    "Create: ok",
			form: new(dto.JSONApplicantRegistrationForm),
			user: new(dto.JSONUser),
			prepare: func(
				repo *repo,
				form *dto.JSONApplicantRegistrationForm,
				user *dto.JSONUser,
			) (*dto.JSONApplicantRegistrationForm, *dto.JSONUser) {
				form = &dto.JSONApplicantRegistrationForm{
					FirstName: "Ivan",
					LastName:  "Ivanov",
					BirthDate: "2000-01-01",
					Email:     "ivanov@ya.ru",
					Password:  "123456",
				}
				applicant := &models.Applicant{
					ID:        1,
				}
				repo.applicant.
					EXPECT().
					GetByEmail(form.Email).
					Return(nil, nil)
				repo.applicant.
					EXPECT().
					Create(gomock.Any()).
					Return(applicant, nil)
				user = &dto.JSONUser{
					ID:        applicant.ID,
					UserType:  dto.UserTypeApplicant,
				}
				return form, user
			},
		},
	}

	
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			applicant: mock.NewMockIApplicantRepository(ctrl),
		}
		tt.form, tt.user =tt.prepare(repo, tt.form, tt.user)

		repositories := &internal.Repositories{
			ApplicantRepository: repo.applicant,
		}
		uc := usecase.NewApplicantUsecase(logrus.New(), repositories)
		user, _ := uc.Create(context.Background(), tt.form)

		require.Equal(t, tt.user, user)
	}
}

func TestGetByID(t *testing.T) {
	t.Parallel()
	type repo struct {
		applicant *mock.MockIApplicantRepository
	}
	tests := []struct {
		name    string
		applicant *dto.JSONApplicantOutput
		prepare func(
			repo *repo,
			applicant *dto.JSONApplicantOutput,
		) (uint64, *dto.JSONApplicantOutput)
	}{
		{
			name:    "Create: bad repository",
			applicant: nil,
			prepare: func(
				repo *repo,
				applicant *dto.JSONApplicantOutput) (uint64, *dto.JSONApplicantOutput) {
				userID := uint64(1)
				repo.applicant.
					EXPECT().
					GetByID(userID).
					Return(nil, errors.New("bad repository"))
				return userID, applicant
			},
		},
		{
			name:    "Create: bad repository",
			applicant: new(dto.JSONApplicantOutput),
			prepare: func(
				repo *repo,
				applicant *dto.JSONApplicantOutput) (uint64, *dto.JSONApplicantOutput) {
				userID := uint64(1)
				applicantModel := &models.Applicant{
					ID:        1,
					FirstName: "Ivan",
					LastName: "Ivanov",
					CityName: "Moscow",
				}
				applicant = &dto.JSONApplicantOutput{
					UserType:  dto.UserTypeApplicant,
					ID:        applicantModel.ID,
					FirstName: applicantModel.FirstName,
					LastName:  applicantModel.LastName,
					CityName:  applicantModel.CityName,
				}
				repo.applicant.
					EXPECT().
					GetByID(userID).
					Return(applicantModel, nil)
				return userID, applicant
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			applicant: mock.NewMockIApplicantRepository(ctrl),
		}
		var userID uint64
		userID, tt.applicant = tt.prepare(repo, tt.applicant)

		repositories := &internal.Repositories{
			ApplicantRepository: repo.applicant,
		}
		uc := usecase.NewApplicantUsecase(logrus.New(), repositories)
		applicant, _ := uc.GetByID(context.Background(), userID)

		require.Equal(t, tt.applicant, applicant)
	}
}

func TestGetAllCities(t *testing.T) {
	t.Parallel()
	type repo struct {
		applicant *mock.MockIApplicantRepository
	}
	tests := []struct {
		name    string
		profile []string
		prepare func(
			repo *repo, profile []string,
		) ([]string)
	}{
		{
			name:    "Create: ok",
			profile: make([]string, 0),
			prepare: func(
				repo *repo, profile []string) ([]string) {
				model := []string{
					"Moscow",
				}

				repo.applicant.
					EXPECT().
					GetAllCities(context.Background(), "Мос").
					Return(model, nil)
				rprofile := []string{
					"Moscow",
					
				}
				return rprofile
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			applicant: mock.NewMockIApplicantRepository(ctrl),
		}
		tt.profile = tt.prepare(repo, tt.profile)

		repositories := &internal.Repositories{
			ApplicantRepository: repo.applicant,
		}
		uc := usecase.NewApplicantUsecase(logrus.New(), repositories)
		profile, _ := uc.GetAllCities(context.Background(), "Мос")

		require.Equal(t, tt.profile, profile)
	}
}