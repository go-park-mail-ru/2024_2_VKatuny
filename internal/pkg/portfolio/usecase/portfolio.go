package usecase

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/repository"
	"github.com/sirupsen/logrus"
)

type IPortfolioUsecase interface {
	GetApplicantPortfolios(applicantID uint64) ([]*dto.JSONGetApplicantPortfolio, error)
}

type PortfolioUsecase struct {
	logger *logrus.Logger
	portfolioRepository repository.IPortfolioRepository
}

func (pu *PortfolioUsecase) GetApplicantPortfolios(applicantID uint64) ([]*dto.JSONGetApplicantPortfolio, error) {
	fn := "PortfolioUsecase.GetApplicantPortfolio"

	pu.logger.Debugf("function: %s; applicant id: %d. Trying to get applicant portfolio", fn, applicantID)
	portfoliosModel, err := pu.portfolioRepository.GetPortfoliosByApplicantID(applicantID)
	if err != nil {
		pu.logger.Errorf("function: %s; got err: %s", fn, err)
		return nil, err
	}

	pu.logger.Debugf("function: %s; successfully got applicant portfolios: %d", fn, len(portfoliosModel))
	portfolio := make([]*dto.JSONGetApplicantPortfolio,  0,len(portfoliosModel))
	for _, portfolioModel := range portfoliosModel {
		portfolio = append(portfolio, &dto.JSONGetApplicantPortfolio{
			ID: portfolioModel.ID,
			ApplicantID: portfolioModel.ApplicantID,
			Title: portfolioModel.Name,
			CreatedAt: portfolioModel.CreatedAt.Format("2006.01.02 15:04:05"),
		})
	}

	return portfolio, nil
}
