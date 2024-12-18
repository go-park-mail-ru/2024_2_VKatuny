package usecase_test

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/usecase"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSearchCVs(t *testing.T) {
	t.Parallel()
	type repo struct {
		cvs *mock.MockICVsRepository
	}
	tests := []struct {
		name                                          string
		offsetStr, numStr, searchStr, group, searchBy string
		cv                                            []*dto.JSONGetApplicantCV
		prepare                                       func(
			repo *repo,
			offsetStr, numStr, searchStr, group, searchBy string,
			cv []*dto.JSONGetApplicantCV,
		) (string, string, string, string, string, []*dto.JSONGetApplicantCV)
	}{
		{
			name:      "Create: ok",
			cv:        make([]*dto.JSONGetApplicantCV, 0),
			offsetStr: "0",
			numStr:    "1",
			searchStr: "Художник",
			group:     "Художник",
			searchBy:  "position",
			prepare: func(
				repo *repo,
				offsetStr, numStr, searchStr, group, searchBy string,
				cv []*dto.JSONGetApplicantCV) (string, string, string, string, string, []*dto.JSONGetApplicantCV) {
				model := []*dto.JSONCv{
					{
						ID:                   0,
						ApplicantID:          0,
						PositionRu:           "PositionRu",
						PositionEn:           "PositionEn",
						Description:          "Description",
						JobSearchStatusName:  "JobSearchStatusName",
						WorkingExperience:    "WorkingExperience",
						Avatar:               "Avatar",
						CompressedAvatar:     "CompressedAvatar",
						PositionCategoryName: "PositionCategoryName",
						CreatedAt:            "CreatedAt",
						UpdatedAt:            "UpdatedAt",
					},
				}

				repo.cvs.
					EXPECT().
					SearchAll(gomock.Any(), uint64(0), uint64(1), searchStr, group, searchBy).
					Return(model, nil)
				rcv := []*dto.JSONGetApplicantCV{
					{
						ID:                   model[0].ID,
						ApplicantID:          model[0].ApplicantID,
						PositionRu:           model[0].PositionRu,
						PositionEn:           model[0].PositionEn,
						Description:          model[0].Description,
						JobSearchStatus:      model[0].JobSearchStatusName,
						WorkingExperience:    model[0].WorkingExperience,
						Avatar:               model[0].Avatar,
						CompressedAvatar:     model[0].CompressedAvatar,
						PositionCategoryName: model[0].PositionCategoryName,
						CreatedAt:            model[0].CreatedAt,
						UpdatedAt:            model[0].UpdatedAt,
					},
				}
				return offsetStr, numStr, searchStr, group, searchBy, rcv
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			cvs: mock.NewMockICVsRepository(ctrl),
		}
		var offsetStr, numStr, searchStr, group, searchBy string
		offsetStr, numStr, searchStr, group, searchBy, tt.cv = tt.prepare(repo, tt.offsetStr, tt.numStr, tt.searchStr, tt.group, tt.searchBy, tt.cv)

		repositories := &internal.Repositories{
			CVRepository: repo.cvs,
		}
		uc := usecase.NewCVsUsecase(logrus.New(), repositories)
		cv, _ := uc.SearchCVs(context.Background(), offsetStr, numStr, searchStr, group, searchBy)

		require.Equal(t, tt.cv, cv)
	}
}

func TestGetApplicantCVs(t *testing.T) {
	t.Parallel()
	type repo struct {
		cvs *mock.MockICVsRepository
	}
	tests := []struct {
		name    string
		ID      uint64
		cv      []*dto.JSONGetApplicantCV
		prepare func(
			repo *repo,
			ID uint64,
			cv []*dto.JSONGetApplicantCV,
		) (uint64, []*dto.JSONGetApplicantCV)
	}{
		{
			name: "Create: ok",
			cv:   make([]*dto.JSONGetApplicantCV, 0),
			ID:   1,
			prepare: func(
				repo *repo,
				ID uint64,
				cv []*dto.JSONGetApplicantCV) (uint64, []*dto.JSONGetApplicantCV) {
				model := []*dto.JSONCv{
					{
						ID:                   0,
						ApplicantID:          0,
						PositionRu:           "PositionRu",
						PositionEn:           "PositionEn",
						Description:          "Description",
						JobSearchStatusName:  "JobSearchStatusName",
						WorkingExperience:    "WorkingExperience",
						Avatar:               "Avatar",
						CompressedAvatar:     "CompressedAvatar",
						PositionCategoryName: "PositionCategoryName",
						CreatedAt:            "CreatedAt",
						UpdatedAt:            "UpdatedAt",
					},
				}

				repo.cvs.
					EXPECT().
					GetCVsByApplicantID(gomock.Any(), uint64(1)).
					Return(model, nil)
				rcv := []*dto.JSONGetApplicantCV{
					{
						ID:                   model[0].ID,
						ApplicantID:          model[0].ApplicantID,
						PositionRu:           model[0].PositionRu,
						PositionEn:           model[0].PositionEn,
						Description:          model[0].Description,
						JobSearchStatus:      model[0].JobSearchStatusName,
						WorkingExperience:    model[0].WorkingExperience,
						Avatar:               model[0].Avatar,
						CompressedAvatar:     model[0].CompressedAvatar,
						PositionCategoryName: model[0].PositionCategoryName,
						CreatedAt:            model[0].CreatedAt,
						UpdatedAt:            model[0].UpdatedAt,
					},
				}
				return 1, rcv
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			cvs: mock.NewMockICVsRepository(ctrl),
		}
		var ID uint64
		ID, tt.cv = tt.prepare(repo, tt.ID, tt.cv)

		repositories := &internal.Repositories{
			CVRepository: repo.cvs,
		}
		uc := usecase.NewCVsUsecase(logrus.New(), repositories)
		cv, _ := uc.GetApplicantCVs(context.Background(),ID)

		require.Equal(t, tt.cv, cv)
	}
}
