package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
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

func (s *PostgreSQLVacanciesStorage) GetVacanciesByEmployerID(employerID uint64) ([]*dto.JSONVacancy, error) {

	Vacancies := make([]*dto.JSONVacancy, 0)

	rows, err := s.db.Query(`select vacancy.id, city.city_name, position, vacancy_description, salary, employer_id, work_type.work_type_name,
		path_to_company_avatar, vacancy.created_at, vacancy.updated_at from vacancy 
		left join work_type on vacancy.work_type_id=work_type.id left join city on vacancy.city_id=city.id where vacancy.employer_id = $1`, employerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var Vacancy dto.JSONVacancy
		if err := rows.Scan(&Vacancy.ID, &Vacancy.Location, &Vacancy.Position, &Vacancy.Description, &Vacancy.Salary,
			&Vacancy.EmployerID, &Vacancy.WorkType, &Vacancy.Avatar, &Vacancy.CreatedAt, &Vacancy.UpdatedAt); err != nil {
			return nil, err
		}
		Vacancies = append(Vacancies, &Vacancy)
		fmt.Println(Vacancy)
	}
	return Vacancies, nil
}

func (s *PostgreSQLVacanciesStorage) GetWithOffset(offset uint64, num uint64) ([]*dto.JSONVacancy, error) {

	Vacancies := make([]*dto.JSONVacancy, 0)

	rows, err := s.db.Query(`select vacancy.id, city.city_name, position, vacancy_description, salary, employer_id, work_type.work_type_name,
		path_to_company_avatar, vacancy.created_at, vacancy.updated_at from vacancy left join work_type on vacancy.work_type_id=work_type.id
		left join city on vacancy.city_id=city.id ORDER BY created_at desc limit $1 offset $2`, num, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var Vacancy dto.JSONVacancy
		if err := rows.Scan(&Vacancy.ID, &Vacancy.Location, &Vacancy.Position, &Vacancy.Description, &Vacancy.Salary, &Vacancy.EmployerID,
			&Vacancy.WorkType, &Vacancy.Avatar, &Vacancy.CreatedAt, &Vacancy.UpdatedAt); err != nil {
			return nil, err
		}
		Vacancies = append(Vacancies, &Vacancy)
		fmt.Println(Vacancy)
	}

	return Vacancies, nil
}

func (s *PostgreSQLVacanciesStorage) Create(vacancy *dto.JSONVacancy) (uint64, error) {
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
	var CityID int
	row = s.db.QueryRow(`select id from city where city_name=$1`, vacancy.Location)
	if err := row.Scan(&CityID); err != nil {
		switch err {
		case sql.ErrNoRows:
			row = s.db.QueryRow(`insert into city (city_name) VALUES ($1) returning id`, vacancy.Location)
			err = row.Scan(&CityID)
			if err != nil {
				return 0, err
			}
		default:
			return 0, err
		}
	}
	var VacancyId uint64
	row = s.db.QueryRow(`insert into vacancy (position, vacancy_description, salary, employer_id, work_type_id,
		path_to_company_avatar, city_id) VALUES ($1, $2, $3, $4, $5, $6, &7) returning id`, vacancy.Position, vacancy.Description, vacancy.Salary, vacancy.EmployerID, WorkTypeID, vacancy.Avatar, CityID)
	err := row.Scan(&VacancyId)
	if err != nil {
		return 0, err
	}
	return VacancyId, err
}

func (s *PostgreSQLVacanciesStorage) GetByID(ID uint64) (*dto.JSONVacancy, error) {

	row := s.db.QueryRow(`select vacancy.id, city.city_name, position, vacancy_description, salary, employer_id, work_type.work_type_name,
		path_to_company_avatar, vacancy.created_at, vacancy.updated_at from vacancy 
		left join work_type on vacancy.work_type_id=work_type.id left join city on vacancy.city_id=city.id where vacancy.id = $1`, ID)
	var oneVacancy dto.JSONVacancy

	err := row.Scan(
		&oneVacancy.ID,
		&oneVacancy.Position,
		&oneVacancy.Description,
		&oneVacancy.Salary,
		&oneVacancy.EmployerID,
		&oneVacancy.CreatedAt,
		&oneVacancy.Avatar,
		&oneVacancy.CreatedAt,
		&oneVacancy.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &oneVacancy, err
}

func (s *PostgreSQLVacanciesStorage) Update(ID uint64, updatedVacancy *dto.JSONVacancy) (*dto.JSONVacancy, error) {

	var CityId int
	row := s.db.QueryRow(`select id from city where city_name=$1`, updatedVacancy.Location)
	if err := row.Scan(&CityId); err != nil {
		switch err {
		case sql.ErrNoRows:
			row = s.db.QueryRow(`insert into city (city_name) VALUES ($1) returning id`, updatedVacancy.Location)
			err = row.Scan(&CityId)
			if err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}
	var WorkTypeID int
	row = s.db.QueryRow(`select id from work_type where work_type_name=$1`, updatedVacancy.WorkType)
	if err := row.Scan(&WorkTypeID); err != nil {
		switch err {
		case sql.ErrNoRows:
			row = s.db.QueryRow(`insert into work_type (work_type_name) VALUES ($1) returning id`, updatedVacancy.WorkType)
			err = row.Scan(&WorkTypeID)
			if err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}
	row = s.db.QueryRow(`update applicant
		set employer_id = $1, salary = $2, position = $3, city_id = $4, vacancy_description = $4,
		work_type_id = $5, path_to_company_avatar = $5 where id=$7 returning id, position, vacancy_description, 
		salary, employer_id, created_at, path_to_company_avatar, created_at, updated_at`,
		updatedVacancy.ID, updatedVacancy.EmployerID, updatedVacancy.Salary, updatedVacancy.Position,
		CityId, updatedVacancy.Description, WorkTypeID, updatedVacancy.Avatar, ID)

	var oneVacancy dto.JSONVacancy

	err := row.Scan(
		&oneVacancy.ID,
		&oneVacancy.Position,
		&oneVacancy.Description,
		&oneVacancy.Salary,
		&oneVacancy.EmployerID,
		&oneVacancy.CreatedAt,
		&oneVacancy.Avatar,
		&oneVacancy.CreatedAt,
		&oneVacancy.UpdatedAt,
	)
	oneVacancy.WorkType = updatedVacancy.WorkType
	oneVacancy.Location = updatedVacancy.Location
	if err != nil {
		return nil, err
	}
	return &oneVacancy, err
}

func (s *PostgreSQLVacanciesStorage) Delete(ID uint64) error {
	_, err := s.db.Exec(`delete from vacancy where id = $1`, ID)
	return err
}

func (s *PostgreSQLVacanciesStorage) Subscribe(ID uint64, applicantID uint64) error {
	_, err := s.db.Exec(`insert into vacancy_subscriber (vacancy_id, applicant_id) VALUES ($1, $2)`, ID, applicantID)
	return err
}

func (s *PostgreSQLVacanciesStorage) GetSubscriptionStatus(ID uint64, applicantID uint64) (bool, error) {
	var rowID uint64
	row := s.db.QueryRow(`select id from vacancy_subscriber where applicant_id=$1 and vacancy_id=$2`, ID, applicantID)
	if err := row.Scan(&rowID); err != nil {
		return false, err
	}
	return true, nil
}

func (s *PostgreSQLVacanciesStorage) GetSubscribersCount(ID uint64) (uint64, error) {
	var rowCount uint64
	row := s.db.QueryRow(`select count(id) from vacancy_subscriber where vacancy_id=$1`, ID)
	if err := row.Scan(&rowCount); err != nil {
		return rowCount, err
	}
	return rowCount, nil
}

func (s *PostgreSQLVacanciesStorage) GetSubscribersList(ID uint64) ([]*models.Applicant, error) {
	Applicants := make([]*models.Applicant, 0)

	rows, err := s.db.Query(`select applicant_id, first_name, last_name, city.city_name, birth_date, path_to_profile_avatar,
		contacts, education, email, password_hash , applicant.created_at, applicant.updated_at
		from vacancy_subscriber	left join applicant on applicant.id = applicant_id
		left join city on city.id =applicant.id where vacancy_subscriber.vacancy_id = $1`, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var oneApplicant models.Applicant
		if err := rows.Scan(&oneApplicant.ID, &oneApplicant.FirstName, &oneApplicant.LastName, &oneApplicant.CityName,
			&oneApplicant.BirthDate, &oneApplicant.PathToProfileAvatar, &oneApplicant.Contacts, &oneApplicant.Education,
			&oneApplicant.Email, &oneApplicant.PasswordHash, &oneApplicant.CreatedAt, &oneApplicant.UpdatedAt); err != nil {
			return nil, err
		}
		Applicants = append(Applicants, &oneApplicant)
		fmt.Println(oneApplicant)
	}
	return Applicants, nil
}

func (s *PostgreSQLVacanciesStorage) Unsubscribe(ID uint64, applicantID uint64) error {
	_, err := s.db.Exec(`delete from vacancy_subscriber where applicant_id=$1 and vacancy_id=$2`, applicantID, ID)
	return err
}
