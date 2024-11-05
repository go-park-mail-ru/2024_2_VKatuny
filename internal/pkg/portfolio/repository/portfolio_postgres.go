package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

type PostgreSQLPortfolioStorage struct {
	db *sql.DB
}

func NewPortfolioStorage(db *sql.DB) *PostgreSQLPortfolioStorage {
	return &PostgreSQLPortfolioStorage{
		db: db,
	}
}

func (s *PostgreSQLPortfolioStorage) GetPortfoliosByApplicantID(applicantID uint64) ([]*models.Portfolio, error) {
	//row := s.db.QueryRow(`select id, applicant_id, portfolio_name, created_at, updated_at  from portfolio where portfolio.applicant_id = $1`, applicantID)

	portfolios := make([]*models.Portfolio, 0)

	rows, err := s.db.Query(`select id, applicant_id, portfolio_name, created_at, updated_at  from portfolio where portfolio.applicant_id = $1`, applicantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var portfolio models.Portfolio
		if err := rows.Scan(&portfolio.ID, &portfolio.ApplicantID, &portfolio.Name, &portfolio.CreatedAt, &portfolio.UpdatedAt); err != nil {
			return nil, err
		}
		portfolios = append(portfolios, &portfolio)
		fmt.Println(portfolio)
	}

	return portfolios, nil
}
