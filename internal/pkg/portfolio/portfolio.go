package portfolio

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

type IPortfolioRepository interface {
	// Add()
	// TODO: need right now
	GetPortfoliosByApplicantID(ctx context.Context, applicantID uint64) ([]*models.Portfolio, error)
}

type IPortfolioUsecase interface {
	GetApplicantPortfolios(ctx context.Context, applicantID uint64) ([]*dto.JSONGetApplicantPortfolio, error)
}
