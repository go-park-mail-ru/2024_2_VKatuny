package portfolio

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

type IPortfolioRepository interface {
	// Add()
	// TODO: need right now
	GetPortfoliosByApplicantID(applicantID uint64) ([]*models.Portfolio, error)
}

type IPortfolioUsecase interface {
	GetApplicantPortfolios(applicantID uint64) ([]*dto.JSONGetApplicantPortfolio, error)
}