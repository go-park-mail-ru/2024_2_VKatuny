package usecase_test

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/usecase"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSearchVacancies(t *testing.T) {
	t.Parallel()
	type repo struct {
		vacancies *mock.MockIVacanciesRepository
	}
	tests := []struct {
		name                                          string
		offsetStr, numStr, searchStr, group, searchBy string
		vacancy                                       []*dto.JSONVacancy
		prepare                                       func(
			repo *repo,
			offsetStr, numStr, searchStr, group, searchBy string,
			vacancy []*dto.JSONVacancy,
		) (string, string, string, string, string, []*dto.JSONVacancy)
	}{
		{
			name:      "Create: ok",
			vacancy:   make([]*dto.JSONVacancy, 0),
			offsetStr: "0",
			numStr:    "1",
			searchStr: "Художник",
			group:     "Художник",
			searchBy:  "position",
			prepare: func(
				repo *repo,
				offsetStr, numStr, searchStr, group, searchBy string,
				vacancy []*dto.JSONVacancy) (string, string, string, string, string, []*dto.JSONVacancy) {
				model := []*dto.JSONVacancy{
					{
						ID:                   0,
						EmployerID:           0,
						Salary:               1000,
						Position:             "Position",
						Location:             "Location",
						Description:          "Description",
						Avatar:               "Avatar",
						WorkType:             "WorkType",
						CompressedAvatar:     "CompressedAvatar",
						PositionCategoryName: "PositionCategoryName",
						CreatedAt:            "CreatedAt",
						UpdatedAt:            "UpdatedAt",
					},
				}

				repo.vacancies.
					EXPECT().
					SearchAll(gomock.Any(), uint64(0), uint64(1), searchStr, group, searchBy).
					Return(model, nil)
				rvacancy := []*dto.JSONVacancy{
					{
						ID:                   model[0].ID,
						EmployerID:           model[0].EmployerID,
						Salary:               model[0].Salary,
						Position:             model[0].Position,
						Location:             model[0].Location,
						WorkType:             model[0].WorkType,
						Description:          model[0].Description,
						Avatar:               model[0].Avatar,
						CompressedAvatar:     model[0].CompressedAvatar,
						PositionCategoryName: model[0].PositionCategoryName,
						CreatedAt:            model[0].CreatedAt,
						UpdatedAt:            model[0].UpdatedAt,
					},
				}
				return offsetStr, numStr, searchStr, group, searchBy, rvacancy
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			vacancies: mock.NewMockIVacanciesRepository(ctrl),
		}
		var offsetStr, numStr, searchStr, group, searchBy string
		offsetStr, numStr, searchStr, group, searchBy, tt.vacancy = tt.prepare(repo, tt.offsetStr, tt.numStr, tt.searchStr, tt.group, tt.searchBy, tt.vacancy)

		repositories := &internal.Repositories{
			VacanciesRepository: repo.vacancies,
		}
		uc := usecase.NewVacanciesUsecase(logrus.New(), repositories)
		vacancy, _ := uc.SearchVacancies(context.Background(), offsetStr, numStr, searchStr, group, searchBy)

		require.Equal(t, tt.vacancy, vacancy)
	}
}

func TestGetEmployerVacancies(t *testing.T) {
	t.Parallel()
	type repo struct {
		vacancies *mock.MockIVacanciesRepository
	}
	tests := []struct {
		name    string
		ID      uint64
		vacancy []*dto.JSONGetEmployerVacancy
		prepare func(
			repo *repo,
			ID uint64,
			vacancy []*dto.JSONGetEmployerVacancy,
		) (uint64, []*dto.JSONGetEmployerVacancy)
	}{
		{
			name:    "Create: ok",
			vacancy: make([]*dto.JSONGetEmployerVacancy, 0),
			ID:      1,
			prepare: func(
				repo *repo,
				ID uint64,
				vacancy []*dto.JSONGetEmployerVacancy) (uint64, []*dto.JSONGetEmployerVacancy) {
				model := []*dto.JSONVacancy{
					{
						ID:                   0,
						EmployerID:           0,
						Salary:               1000,
						Position:             "Position",
						Location:             "Location",
						Description:          "Description",
						Avatar:               "Avatar",
						WorkType:             "WorkType",
						CompressedAvatar:     "CompressedAvatar",
						PositionCategoryName: "PositionCategoryName",
						CreatedAt:            "CreatedAt",
						UpdatedAt:            "UpdatedAt",
					},
				}

				repo.vacancies.
					EXPECT().
					GetVacanciesByEmployerID(gomock.Any(), uint64(1)).
					Return(model, nil)
				rvacancy := []*dto.JSONGetEmployerVacancy{
					{
						ID:                   model[0].ID,
						EmployerID:           model[0].EmployerID,
						Salary:               model[0].Salary,
						Position:             model[0].Position,
						Location:             model[0].Location,
						WorkType:             model[0].WorkType,
						Description:          model[0].Description,
						Avatar:               model[0].Avatar,
						CompressedAvatar:     model[0].CompressedAvatar,
						PositionCategoryName: model[0].PositionCategoryName,
						CreatedAt:            model[0].CreatedAt,
						UpdatedAt:            model[0].UpdatedAt,
					},
				}
				return 1, rvacancy
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			vacancies: mock.NewMockIVacanciesRepository(ctrl),
		}
		var ID uint64
		ID, tt.vacancy = tt.prepare(repo, tt.ID, tt.vacancy)

		repositories := &internal.Repositories{
			VacanciesRepository: repo.vacancies,
		}
		uc := usecase.NewVacanciesUsecase(logrus.New(), repositories)
		vacancy, _ := uc.GetVacanciesByEmployerID(context.Background(), ID)

		require.Equal(t, tt.vacancy, vacancy)
	}
}

func TestGetApplicantFavoriteVacancies(t *testing.T) {
	t.Parallel()
	type repo struct {
		vacancies *mock.MockIVacanciesRepository
	}
	tests := []struct {
		name    string
		ID      uint64
		vacancy []*dto.JSONGetEmployerVacancy
		prepare func(
			repo *repo,
			ID uint64,
			vacancy []*dto.JSONGetEmployerVacancy,
		) (uint64, []*dto.JSONGetEmployerVacancy)
	}{
		{
			name:    "Create: ok",
			vacancy: make([]*dto.JSONGetEmployerVacancy, 0),
			ID:      1,
			prepare: func(
				repo *repo,
				ID uint64,
				vacancy []*dto.JSONGetEmployerVacancy) (uint64, []*dto.JSONGetEmployerVacancy) {
				model := []*dto.JSONVacancy{
					{
						ID:                   0,
						EmployerID:           0,
						Salary:               1000,
						Position:             "Position",
						Description:          "Description",
						Avatar:               "Avatar",
						WorkType:             "WorkType",
						CompressedAvatar:     "CompressedAvatar",
						PositionCategoryName: "PositionCategoryName",
						CreatedAt:            "CreatedAt",
						UpdatedAt:            "UpdatedAt",
					},
				}

				repo.vacancies.
					EXPECT().
					GetApplicantFavoriteVacancies(gomock.Any(), uint64(1)).
					Return(model, nil)
				rvacancy := []*dto.JSONGetEmployerVacancy{
					{
						ID:                   model[0].ID,
						EmployerID:           model[0].EmployerID,
						Salary:               model[0].Salary,
						Position:             model[0].Position,
						WorkType:             model[0].WorkType,
						Description:          model[0].Description,
						Avatar:               model[0].Avatar,
						CompressedAvatar:     model[0].CompressedAvatar,
						PositionCategoryName: model[0].PositionCategoryName,
						CreatedAt:            model[0].CreatedAt,
						UpdatedAt:            model[0].UpdatedAt,
					},
				}
				return 1, rvacancy
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			vacancies: mock.NewMockIVacanciesRepository(ctrl),
		}
		var ID uint64
		ID, tt.vacancy = tt.prepare(repo, tt.ID, tt.vacancy)

		repositories := &internal.Repositories{
			VacanciesRepository: repo.vacancies,
		}
		uc := usecase.NewVacanciesUsecase(logrus.New(), repositories)
		vacancy, _ := uc.GetApplicantFavoriteVacancies(context.Background(), ID)
		require.Equal(t, tt.vacancy, vacancy)
	}
}

func TestGetVacancySubscribers(t *testing.T) {
	t.Parallel()
	type repo struct {
		vacancies *mock.MockIVacanciesRepository
	}
	tests := []struct {
		name    string
		ID      uint64
		profile *dto.JSONVacancySubscribers
		user    *dto.UserFromSession
		prepare func(
			repo *repo,
			ID uint64, user *dto.UserFromSession,
			profile *dto.JSONVacancySubscribers,
		) (uint64, *dto.UserFromSession, *dto.JSONVacancySubscribers)
	}{
		{
			name:    "Create: ok",
			profile: &dto.JSONVacancySubscribers{},
			ID:      1,
			prepare: func(
				repo *repo,
				ID uint64, user *dto.UserFromSession,
				profile *dto.JSONVacancySubscribers,
			) (uint64, *dto.UserFromSession, *dto.JSONVacancySubscribers) {
				ruser := &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}
				rprofile := []*models.Applicant{
					{
						ID:                  0,
						FirstName:           "FirstName",
						LastName:            "LastName",
						CityName:            "City",
						BirthDate:           "BirthDate",
						PathToProfileAvatar: "Avatar",
						CompressedAvatar:    "CompressedAvatar",
						Contacts:            "Contacts",
						Education:           "Education",
					},
				}
				rprofile1 := []*dto.JSONGetApplicantProfile{
					{
						ID:               0,
						FirstName:        "FirstName",
						LastName:         "LastName",
						City:             "City",
						BirthDate:        "BirthDate",
						Avatar:           "Avatar",
						CompressedAvatar: "CompressedAvatar",
						Contacts:         "Contacts",
						Education:        "Education",
					},
				}
				vacacy := &dto.JSONVacancy{
					ID:         0,
					EmployerID: 1,
					Salary:     1000,
				}
				repo.vacancies.
					EXPECT().
					GetByID(gomock.Any(), uint64(1)).
					Return(vacacy, nil)
				repo.vacancies.
					EXPECT().
					GetSubscribersList(gomock.Any(), uint64(1)).
					Return(rprofile, nil)

				model := &dto.JSONVacancySubscribers{
						ID:          1,
						Subscribers: rprofile1,
				}
				return 1, ruser, model
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			vacancies: mock.NewMockIVacanciesRepository(ctrl),
		}
		var ID uint64
		var user *dto.UserFromSession
		ID, user, tt.profile = tt.prepare(repo, tt.ID, tt.user, tt.profile)

		repositories := &internal.Repositories{
			VacanciesRepository: repo.vacancies,
		}
		uc := usecase.NewVacanciesUsecase(logrus.New(), repositories)
		profile, _ := uc.GetVacancySubscribers(context.Background(), ID, user)
		require.Equal(t, tt.profile, profile)
	}
}
