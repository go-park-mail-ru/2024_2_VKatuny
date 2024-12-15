package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/usecase"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type repo struct {
		portfolio *mock.MockIPortfolioRepository
	}
	tests := []struct {
		name       string
		portfolios []*dto.JSONGetApplicantPortfolio
		prepare    func(
			repo *repo,
			portfolios []*dto.JSONGetApplicantPortfolio,
		) (uint64, []*dto.JSONGetApplicantPortfolio)
	}{
		{
			name:       "Create: bad repository",
			portfolios: nil,
			prepare: func(
				repo *repo,
				portfolios []*dto.JSONGetApplicantPortfolio) (uint64, []*dto.JSONGetApplicantPortfolio) {
				userID := uint64(1)
				repo.portfolio.
					EXPECT().
					GetPortfoliosByApplicantID(userID).
					Return(nil, fmt.Errorf(dto.MsgDataBaseError))
				return userID, portfolios
			},
		},
		{
			name:       "Create: ok",
			portfolios: make([]*dto.JSONGetApplicantPortfolio, 0),
			prepare: func(
				repo *repo,
				portfolios []*dto.JSONGetApplicantPortfolio) (uint64, []*dto.JSONGetApplicantPortfolio) {
				userID := uint64(1)
				model := []*models.Portfolio{&models.Portfolio{
					ID:   userID,
					Name: "Name",
				},
				}

				repo.portfolio.
					EXPECT().
					GetPortfoliosByApplicantID(userID).
					Return(model, nil)
				rportfolios := []*dto.JSONGetApplicantPortfolio{

					&dto.JSONGetApplicantPortfolio{
						ID:   model[0].ID,
						Name: model[0].Name,
					},
				}
				return userID, rportfolios
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			portfolio: mock.NewMockIPortfolioRepository(ctrl),
		}
		var userID uint64
		userID, tt.portfolios = tt.prepare(repo, tt.portfolios)

		repositories := &internal.Repositories{
			PortfolioRepository: repo.portfolio,
		}
		uc := usecase.NewPortfolioUsecase(logrus.New(), repositories)
		portfolios, _ := uc.GetApplicantPortfolios(context.Background(), userID)

		require.Equal(t, tt.portfolios, portfolios)
	}
}
