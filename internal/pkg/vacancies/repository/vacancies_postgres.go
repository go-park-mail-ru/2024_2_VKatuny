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

	rows, err := s.db.Query(`select vacancy.id, position, vacancy_description, salary, employer_id, work_type_name,
		path_to_company_avatar, vacancy.created_at from vacancy left join work_type on vacancy.work_type_id=work_type.id where vacancy.employer_id = $1`, employerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var Vacancy models.Vacancy
		if err := rows.Scan(&Vacancy.ID, &Vacancy.Position, &Vacancy.Description, &Vacancy.Salary, &Vacancy.EmployerID, &Vacancy.WorkType, &Vacancy.Logo, &Vacancy.CreatedAt); err != nil {
			return nil, err
		}
		Vacancies = append(Vacancies, &Vacancy)
		fmt.Println(Vacancy)
	}

	return Vacancies, nil
}

func (s *PostgreSQLVacanciesStorage) GetWithOffset(offset uint64, num uint64) ([]*models.Vacancy, error) {

	Vacancies := make([]*models.Vacancy, 0)

	rows, err := s.db.Query(`select id, position, vacancy_description, salary, employer_id, work_type_name,
		path_to_company_avatar, created_at from vacancy ORDER BY created_at desc limit $1 offset $2`, num, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var Vacancy models.Vacancy
		if err := rows.Scan(&Vacancy.ID, &Vacancy.Position, &Vacancy.Description, &Vacancy.Salary, &Vacancy.EmployerID, &Vacancy.WorkType, &Vacancy.Logo, &Vacancy.CreatedAt); err != nil {
			return nil, err
		}
		Vacancies = append(Vacancies, &Vacancy)
		fmt.Println(Vacancy)
	}

	return Vacancies, nil
}

func (s *PostgreSQLVacanciesStorage) Create(vacancy *models.Vacancy) (uint64, error) {
	var WorkTypeID int
	row := s.db.QueryRow(`select id from work_type where work_type_name=$1`, vacancy.WorkType)
	if err := row.Scan(&WorkTypeID); err != nil {
		switch err {
		case sql.ErrNoRows:
			row = s.db.QueryRow(`insert into work_type (work_type_name) VALUES ($1) returning id`, vacancy.WorkType)
			err = row.Scan(&WorkTypeID)
			if err != nil {
				return 0, err
			}
		default:
			return 0, err
		}
	}
	var VacancyId uint64
	row = s.db.QueryRow(`insert into vacancy (position, vacancy_description, salary, employer_id, work_type_id,
		path_to_company_avatar) VALUES ($1, $2, $3, $4, $5) returning id`, vacancy.Position, vacancy.Description, vacancy.Salary, vacancy.EmployerID, WorkTypeID, vacancy.Logo)
	err := row.Scan(&VacancyId)
	return VacancyId, err
}
