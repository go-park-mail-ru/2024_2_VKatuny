package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

type PortfolioUsecase struct {
	logger        *logrus.Entry
	portfolioRepo portfolio.IPortfolioRepository
}

func NewPortfolioUsecase(logger *logrus.Logger, repositories *internal.Repositories) *PortfolioUsecase {
	return &PortfolioUsecase{
		logger:        logrus.NewEntry(logger),
		portfolioRepo: repositories.PortfolioRepository,
	}
}

func (pu *PortfolioUsecase) GetApplicantPortfolios(ctx context.Context, applicantID uint64) ([]*dto.JSONGetApplicantPortfolio, error) {
	fn := "PortfolioUsecase.GetApplicantPortfolio"
	pu.logger = utils.SetLoggerRequestID(ctx, pu.logger)


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
