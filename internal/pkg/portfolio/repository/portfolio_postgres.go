package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

type PostgreSQLPortfolioStorage struct {
	db *sql.DB
	logger *logrus.Entry
}

func NewPortfolioStorage(db *sql.DB, logger *logrus.Logger) *PostgreSQLPortfolioStorage {
	return &PostgreSQLPortfolioStorage{
		db: db,
		logger: logrus.NewEntry(logger),
	}
}

func (s *PostgreSQLPortfolioStorage) GetPortfoliosByApplicantID(ctx context.Context, applicantID uint64) ([]*models.Portfolio, error) {
	funcName := "PostgreSQLPortfolioStorage.GetPortfoliosByApplicantID"
	s.logger = utils.SetLoggerRequestID(ctx, s.logger)
	s.logger.Debugf("%s: entering", funcName)

	portfolios := make([]*models.Portfolio, 0)

	rows, err := s.db.Query(`select id, applicant_id, portfolio_name, created_at, updated_at  from portfolio where portfolio.applicant_id = $1`, applicantID)
	if err != nil {
		s.logger.Errorf("%s: got error %v", funcName, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var portfolio models.Portfolio
		if err := rows.Scan(&portfolio.ID, &portfolio.ApplicantID, &portfolio.Name, &portfolio.CreatedAt, &portfolio.UpdatedAt); err != nil {
			s.logger.Errorf("%s: got error %v", funcName, err)
			return nil, err
		}
		portfolios = append(portfolios, &portfolio)
		s.logger.Debugf("%s: got portfolio from db %v", funcName, portfolio)
	}
	s.logger.Debugf("%s: got portfolios, count = %d", funcName, len(portfolios))

	return portfolios, nil
}
