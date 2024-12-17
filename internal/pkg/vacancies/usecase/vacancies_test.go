package usecase_test

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
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
					&dto.JSONVacancy{
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
					&dto.JSONVacancy{
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
					&dto.JSONVacancy{
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
					&dto.JSONGetEmployerVacancy{
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
