package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/repository"
	"github.com/sirupsen/logrus"
)

type IPortfolioUsecase interface {
	GetApplicantPortfolios(applicantID uint64) ([]*dto.JSONGetApplicantPortfolio, error)
}

type PortfolioUsecase struct {
	logger        *logrus.Logger
	portfolioRepo repository.IPortfolioRepository
}

func NewPortfolioUsecase(logger *logrus.Logger, repositories *internal.Repositories) *PortfolioUsecase {
	PortfolioRepository, ok := repositories.PortfolioRepository.(repository.IPortfolioRepository)
	if !ok {
		return nil
	}
	return &PortfolioUsecase{
		logger:        logger,
		portfolioRepo: PortfolioRepository,
	}
}

func (pu *PortfolioUsecase) GetApplicantPortfolios(applicantID uint64) ([]*dto.JSONGetApplicantPortfolio, error) {
	fn := "PortfolioUsecase.GetApplicantPortfolio"

	pu.logger.Debugf("function: %s; applicant id: %d. Trying to get applicant portfolio", fn, applicantID)
	portfoliosModel, err := pu.portfolioRepo.GetPortfoliosByApplicantID(applicantID)
	if err != nil {
		pu.logger.Errorf("function: %s; got err: %s", fn, err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}

	pu.logger.Debugf("function: %s; successfully got applicant portfolios: %d", fn, len(portfoliosModel))
	portfolio := make([]*dto.JSONGetApplicantPortfolio, 0, len(portfoliosModel))
	for _, portfolioModel := range portfoliosModel {
		portfolio = append(portfolio, &dto.JSONGetApplicantPortfolio{
			ID:          portfolioModel.ID,
			ApplicantID: portfolioModel.ApplicantID,
			Name:        portfolioModel.Name,
			CreatedAt:   portfolioModel.CreatedAt,
		})
	}

	return portfolio, nil
}
