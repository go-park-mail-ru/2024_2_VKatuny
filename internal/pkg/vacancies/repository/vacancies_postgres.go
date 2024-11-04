package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

type PostgreSQLVacanciesStorage struct {
	db *sql.DB
}

func NewVacanciesStorage(db *sql.DB) *PostgreSQLVacanciesStorage {
	return &PostgreSQLVacanciesStorage{
		db: db,
	}
}

func (s *PostgreSQLVacanciesStorage) GetVacanciesByEmployerID(employerID uint64) ([]*models.Vacancy, error) {

	Vacancies := make([]*models.Vacancy, 0)

	rows, err := s.db.Query(`select id, position, description, salary, employer_id,
		path_to_company_avatar, created_at from vacancy where vacancy.applicant_id = $1`, employerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var Vacancy models.Vacancy
		if err := rows.Scan(&Vacancy.ID, &Vacancy.Position, &Vacancy.Description, &Vacancy.Salary, &Vacancy.Employer, &Vacancy.Logo, &Vacancy.CreatedAt); err != nil {
			return nil, err
		}
		fmt.Println(Vacancy)
	}
	if !rows.NextResultSet() {
		return nil, fmt.Errorf("err with rows count")
	}

	return Vacancies, nil
}

func (s *PostgreSQLVacanciesStorage) GetWithOffset(offset uint64, num uint64) ([]*models.Vacancy, error) {

	Vacancies := make([]*models.Vacancy, 0)

	rows, err := s.db.Query(`select id, position, description, salary, employer_id,
		path_to_company_avatar, created_at from vacancy ORDER BY created_at desc limit $1 offset $2`, num, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var Vacancy models.Vacancy
		if err := rows.Scan(&Vacancy.ID, &Vacancy.Position, &Vacancy.Description, &Vacancy.Salary, &Vacancy.Employer, &Vacancy.Logo, &Vacancy.CreatedAt); err != nil {
			return nil, err
		}
		fmt.Println(Vacancy)
	}
	if !rows.NextResultSet() {
		return nil, fmt.Errorf("err with rows count")
	}

	return Vacancies, nil
}

func (s *PostgreSQLVacanciesStorage) Add(vacancy *models.Vacancy) (uint64, error) {
	var VacancyId uint64
	row := s.db.QueryRow(`insert into vacancy (position, description, salary, employer_id,
		path_to_company_avatar) VALUES ($1, $2, $3, $4, $5) returning id`, vacancy.Position, vacancy.Description, vacancy.Salary, vacancy.Employer, vacancy.Logo)
	err := row.Scan(&VacancyId)
	return VacancyId, err
}
