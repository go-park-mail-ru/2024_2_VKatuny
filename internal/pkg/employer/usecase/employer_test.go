package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/usecase"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"



)

func TestGet(t *testing.T) {
	t.Parallel()
	type repo struct {
		employer *mock.MockIEmployerRepository
	}
	tests := []struct {
		name    string
		profile *dto.JSONGetEmployerProfile
		prepare func(
			repo *repo,
			profile *dto.JSONGetEmployerProfile,
		) (uint64, *dto.JSONGetEmployerProfile)
	}{
		{
			name:    "Create: bad repository",
			profile: nil,
			prepare: func(
				repo *repo,
				profile *dto.JSONGetEmployerProfile) (uint64, *dto.JSONGetEmployerProfile) {
				userID := uint64(1)
				repo.employer.
					EXPECT().
					GetByID(userID).
					Return(nil, errors.New("bad repository"))
				return userID, profile
			},
		},
		{
			name:    "Create: ok",
			profile: new(dto.JSONGetEmployerProfile),
			prepare: func(
				repo *repo,
				profile *dto.JSONGetEmployerProfile) (uint64, *dto.JSONGetEmployerProfile) {
				userID := uint64(1)
				model := &models.Employer{
					ID:        userID,
					FirstName: "Ivan",
					LastName:  "Ivanov",
					CityName:  "Moscow",
				}

				repo.employer.
					EXPECT().
					GetByID(userID).
					Return(model, nil)
				rprofile := &dto.JSONGetEmployerProfile{
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
			employer: mock.NewMockIEmployerRepository(ctrl),
		}
		var userID uint64
		userID, tt.profile = tt.prepare(repo, tt.profile)

		repositories := &internal.Repositories{
			EmployerRepository: repo.employer,
		}
		uc := usecase.NewEmployerUsecase(logrus.New(), repositories)
		profile, _ := uc.GetEmployerProfile(context.Background(), userID)

		require.Equal(t, tt.profile, profile)
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	type repo struct {
		employer *mock.MockIEmployerRepository
	}
	tests := []struct {
		name    string
		profile *dto.JSONUpdateEmployerProfile
		expected error
		prepare func(
			repo *repo,
			profile *dto.JSONUpdateEmployerProfile,
		) (uint64, *dto.JSONUpdateEmployerProfile)
	}{
		{
			name:    "Create: bad repository",
			profile: nil,
			expected: errors.New("bad repository"),
			prepare: func(
				repo *repo,
				profile *dto.JSONUpdateEmployerProfile) (uint64, *dto.JSONUpdateEmployerProfile) {
				userID := uint64(1)
				repo.employer.
					EXPECT().
					Update(userID, gomock.Any()).
					Return(nil, errors.New("bad repository"))
				return userID, profile
			},
		},
		{
			name:    "Create: ok",
			profile: new(dto.JSONUpdateEmployerProfile),
			expected: nil,
			prepare: func(
				repo *repo,
				profile *dto.JSONUpdateEmployerProfile) (uint64, *dto.JSONUpdateEmployerProfile) {
				userID := uint64(1)
				model := &models.Employer{
					ID:        userID,
					FirstName: "Ivan",
					LastName:  "Ivanov",
					CityName:  "Moscow",
				}

				repo.employer.
					EXPECT().
					Update(userID, gomock.Any()).
					Return(nil, nil)
				rprofile := &dto.JSONUpdateEmployerProfile{
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
			employer: mock.NewMockIEmployerRepository(ctrl),
		}
		var userID uint64
		userID, tt.profile = tt.prepare(repo, tt.profile)

		repositories := &internal.Repositories{
			EmployerRepository: repo.employer,
		}
		uc := usecase.NewEmployerUsecase(logrus.New(), repositories)
		err := uc.UpdateEmployerProfile(context.Background(), userID, tt.profile)

		require.Equal(t, tt.expected ,err)
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()
	type repo struct {
		employer *mock.MockIEmployerRepository
	} 
	tests := []struct {
		name    string
		form *dto.JSONEmployerRegistrationForm
		user *dto.JSONUser
		prepare func(
			repo *repo,
			form *dto.JSONEmployerRegistrationForm,
			user *dto.JSONUser,
		) (*dto.JSONEmployerRegistrationForm, *dto.JSONUser)
	}{
		{
			name:    "Create: bad repository get",
			form: new(dto.JSONEmployerRegistrationForm),
			user: new(dto.JSONUser),
			prepare: func(
				repo *repo,
				form *dto.JSONEmployerRegistrationForm,
				user *dto.JSONUser,
			) (*dto.JSONEmployerRegistrationForm, *dto.JSONUser) {
				repo.employer.
					EXPECT().
					GetByEmail(form.Email).
					Return(nil, errors.New("sql: no rows in result sett"))
				return form, nil
			},
		},
		{
			name:    "Create: bad repository create",
			form: new(dto.JSONEmployerRegistrationForm),
			user: new(dto.JSONUser),
			prepare: func(
				repo *repo,
				form *dto.JSONEmployerRegistrationForm,
				user *dto.JSONUser,
			) (*dto.JSONEmployerRegistrationForm, *dto.JSONUser) {
				form = &dto.JSONEmployerRegistrationForm{
					FirstName: "Ivan",
					LastName:  "Ivanov",
					Email:     "ivanov@ya.ru",
					Password:  "123456",
				}
				employer := &models.Employer{
					ID:        1,
				}
				repo.employer.
					EXPECT().
					GetByEmail(form.Email).
					Return(nil, nil)
				repo.employer.
					EXPECT().
					Create(gomock.Any()).
					Return(employer, errors.New("bad repository"))
				return form, nil
			},
		},
		{
			name:    "Create: ok",
			form: new(dto.JSONEmployerRegistrationForm),
			user: new(dto.JSONUser),
			prepare: func(
				repo *repo,
				form *dto.JSONEmployerRegistrationForm,
				user *dto.JSONUser,
			) (*dto.JSONEmployerRegistrationForm, *dto.JSONUser) {
				form = &dto.JSONEmployerRegistrationForm{
					FirstName: "Ivan",
					LastName:  "Ivanov",
					Email:     "ivanov@ya.ru",
					Password:  "123456",
				}
				employer := &models.Employer{
					ID:        1,
				}
				repo.employer.
					EXPECT().
					GetByEmail(form.Email).
					Return(nil, nil)
				repo.employer.
					EXPECT().
					Create(gomock.Any()).
					Return(employer, nil)
				user = &dto.JSONUser{
					ID:        employer.ID,
					UserType:  dto.UserTypeEmployer,
				}
				return form, user
			},
		},
	}

	
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			employer: mock.NewMockIEmployerRepository(ctrl),
		}
		tt.form, tt.user =tt.prepare(repo, tt.form, tt.user)

		repositories := &internal.Repositories{
			EmployerRepository: repo.employer,
		}
		uc := usecase.NewEmployerUsecase(logrus.New(), repositories)
		user, _ := uc.Create(context.Background(), tt.form)

		require.Equal(t, tt.user, user)
	}
}

func TestGetByID(t *testing.T) {
	t.Parallel()
	type repo struct {
		employer *mock.MockIEmployerRepository
	}
	tests := []struct {
		name    string
		employer *dto.JSONEmployer
		prepare func(
			repo *repo,
			employer *dto.JSONEmployer,
		) (uint64, *dto.JSONEmployer)
	}{
		{
			name:    "Create: bad repository",
			employer: nil,
			prepare: func(
				repo *repo,
				employer *dto.JSONEmployer) (uint64, *dto.JSONEmployer) {
				userID := uint64(1)
				repo.employer.
					EXPECT().
					GetByID(userID).
					Return(nil, errors.New("bad repository"))
				return userID, employer
			},
		},
		{
			name:    "Create: bad repository",
			employer: new(dto.JSONEmployer),
			prepare: func(
				repo *repo,
				employer *dto.JSONEmployer) (uint64, *dto.JSONEmployer) {
				userID := uint64(1)
				EmployerModel := &models.Employer{
					ID:        1,
					FirstName: "Ivan",
					LastName: "Ivanov",
					CityName: "Moscow",
				}
				employer = &dto.JSONEmployer{
					UserType:  dto.UserTypeEmployer,
					ID:        EmployerModel.ID,
					FirstName: EmployerModel.FirstName,
					LastName:  EmployerModel.LastName,
					CityName:  EmployerModel.CityName,
				}
				repo.employer.
					EXPECT().
					GetByID(userID).
					Return(EmployerModel, nil)
				return userID, employer
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			employer: mock.NewMockIEmployerRepository(ctrl),
		}
		var userID uint64
		userID, tt.employer = tt.prepare(repo, tt.employer)

		repositories := &internal.Repositories{
			EmployerRepository: repo.employer,
		}
		uc := usecase.NewEmployerUsecase(logrus.New(), repositories)
		employer, _ := uc.GetByID(context.Background(), userID)

		require.Equal(t, tt.employer, employer)
	}
}
