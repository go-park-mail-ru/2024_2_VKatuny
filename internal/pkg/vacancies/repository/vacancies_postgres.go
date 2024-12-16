package repository

import (
	"database/sql"
	"fmt"
	"strconv"

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

	rows, err := s.db.Query(`select vacancy.id, city.city_name, vacancy.position, vacancy_description, salary, employer_id, work_type.work_type_name, path_to_company_avatar, vacancy.created_at, vacancy.updated_at, 
		company.company_name, position_category.category_name, vacancy.compressed_image from vacancy
		left join work_type on vacancy.work_type_id=work_type.id left join city on vacancy.city_id=city.id
		left join employer on vacancy.employer_id=employer.id left join company on employer.company_name_id=company.id
		left join position_category on vacancy.position_category_id = position_category.id
		where vacancy.employer_id = $1`, employerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var Vacancy dto.JSONVacancyWithNull
		if err := rows.Scan(
			&Vacancy.ID,
			&Vacancy.Location,
			&Vacancy.Position,
			&Vacancy.Description,
			&Vacancy.Salary,
			&Vacancy.EmployerID,
			&Vacancy.WorkType,
			&Vacancy.Avatar,
			&Vacancy.CreatedAt,
			&Vacancy.UpdatedAt,
			&Vacancy.CompanyName,
			&Vacancy.PositionCategoryName,
			&Vacancy.CompressedAvatar); err != nil {
			return nil, err
		}
		VacancyOk := dto.JSONVacancy{
			ID:                   Vacancy.ID,
			EmployerID:           Vacancy.EmployerID,
			Salary:               Vacancy.Salary,
			Position:             Vacancy.Position,
			Location:             Vacancy.Location,
			Description:          Vacancy.Description,
			WorkType:             Vacancy.WorkType,
			Avatar:               Vacancy.Avatar,
			CompanyName:          Vacancy.CompanyName,
			PositionCategoryName: Vacancy.PositionCategoryName.String,
			CompressedAvatar:     Vacancy.CompressedAvatar.String,
			CreatedAt:            Vacancy.CreatedAt,
			UpdatedAt:            Vacancy.UpdatedAt,
		}
		Vacancies = append(Vacancies, &VacancyOk)
		//TODO logger
		fmt.Println(VacancyOk)
	}
	return Vacancies, nil
}

func (s *PostgreSQLVacanciesStorage) SearchAll(offset uint64, num uint64, searchStr, group, searchBy string) ([]*dto.JSONVacancy, error) {
	Vacancies := make([]*dto.JSONVacancy, 0)
	iter := 1
	mainPart := `select vacancy.id, city.city_name, vacancy.position, vacancy_description, salary, employer_id, work_type.work_type_name, path_to_company_avatar, vacancy.created_at, vacancy.updated_at, 
		company.company_name, position_category.category_name, vacancy.compressed_image from vacancy
		left join work_type on vacancy.work_type_id=work_type.id left join city on vacancy.city_id=city.id
		left join employer on vacancy.employer_id=employer.id left join company on employer.company_name_id=company.id
		left join position_category on vacancy.position_category_id = position_category.id `
	categoryPart := ""
	if group != "" {
		categoryPart = "where position_category.category_name = $" + strconv.Itoa(iter)
		iter++
	}
	searchPart := ""
	if searchStr != "" {
		if iter != 1 {
			searchPart += " and "
		} else {
			searchPart += " where "
		}
		if searchBy == "company" {
			searchPart += "company.id in (select company.id from company where ts_rank_cd(company.fts, plainto_tsquery('russian', $" + strconv.Itoa(iter) + ")) <> 0)  order by ts_rank_cd(vacancy.fts, plainto_tsquery('russian', $" + strconv.Itoa(iter+1) + ")) desc "
		} else {
			weights := "'{0, 0, 1, 1}'"
			if searchBy == "position" {
				weights = "'{0, 0, 0, 1}'"
			} else if searchBy == "description" {
				weights = "'{0, 0, 1, 0}'"
			}
			searchPart += "ts_rank_cd(" + weights + ", vacancy.fts, plainto_tsquery('russian', $" + strconv.Itoa(iter) + ")) <> 0 order by ts_rank_cd(" + weights + ", vacancy.fts, plainto_tsquery('russian', $" + strconv.Itoa(iter+1) + ")) desc "
		}
		iter += 2
	} else {
		searchPart += " ORDER BY created_at desc "
	}

	lastPart := " limit $" + strconv.Itoa(iter) + " offset $" + strconv.Itoa(iter+1)
	iter += 2
	var rows *sql.Rows
	var err error
	if group != "" && searchStr != "" {
		rows, err = s.db.Query(mainPart+categoryPart+searchPart+lastPart, group, searchStr, searchStr, num, offset)
	} else if group == "" && searchStr != "" {
		rows, err = s.db.Query(mainPart+categoryPart+searchPart+lastPart, searchStr, searchStr, num, offset)
	} else if group == "" && searchStr == "" {
		rows, err = s.db.Query(mainPart+categoryPart+searchPart+lastPart, num, offset)
	} else if group != "" && searchStr == "" {
		rows, err = s.db.Query(mainPart+categoryPart+searchPart+lastPart, group, num, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var Vacancy dto.JSONVacancyWithNull
		if err := rows.Scan(&Vacancy.ID, &Vacancy.Location, &Vacancy.Position, &Vacancy.Description, &Vacancy.Salary, &Vacancy.EmployerID,
			&Vacancy.WorkType, &Vacancy.Avatar, &Vacancy.CreatedAt, &Vacancy.UpdatedAt, &Vacancy.CompanyName, &Vacancy.PositionCategoryName, &Vacancy.CompressedAvatar); err != nil {
			return nil, err
		}
		VacancyOk := dto.JSONVacancy{
			ID:                   Vacancy.ID,
			EmployerID:           Vacancy.EmployerID,
			Salary:               Vacancy.Salary,
			Position:             Vacancy.Position,
			Location:             Vacancy.Location,
			Description:          Vacancy.Description,
			WorkType:             Vacancy.WorkType,
			Avatar:               Vacancy.Avatar,
			CompanyName:          Vacancy.CompanyName,
			PositionCategoryName: Vacancy.PositionCategoryName.String,
			CreatedAt:            Vacancy.CreatedAt,
			UpdatedAt:            Vacancy.UpdatedAt,
			CompressedAvatar:     Vacancy.CompressedAvatar.String,
		}
		Vacancies = append(Vacancies, &VacancyOk)
		//TODO logger
		fmt.Println(VacancyOk)
	}
	//TODO logger
	fmt.Println(Vacancies)
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
	if vacancy.PositionCategoryName == "" {
		row = s.db.QueryRow(`insert into vacancy (position, vacancy_description, salary, employer_id, work_type_id,
		path_to_company_avatar, city_id, compressed_image) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) returning id`, vacancy.Position, vacancy.Description,
			vacancy.Salary, vacancy.EmployerID, WorkTypeID, vacancy.Avatar, CityID, vacancy.CompressedAvatar)
	} else {
		var PositionCategoryID int
		row = s.db.QueryRow(`select id from position_category where category_name=$1`, vacancy.PositionCategoryName)
		err := row.Scan(&PositionCategoryID)
		if err != nil {
			return 0, err
		}
		row = s.db.QueryRow(`insert into vacancy (position, vacancy_description, salary, employer_id, work_type_id,
		path_to_company_avatar, city_id, position_category_id, compressed_image) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`, vacancy.Position, vacancy.Description,
			vacancy.Salary, vacancy.EmployerID, WorkTypeID, vacancy.Avatar, CityID, PositionCategoryID, vacancy.CompressedAvatar)
	}
	err := row.Scan(&VacancyId)
	if err != nil {
		return 0, err
	}
	return VacancyId, nil
}

func (s *PostgreSQLVacanciesStorage) GetByID(ID uint64) (*dto.JSONVacancy, error) {

	row := s.db.QueryRow(`select vacancy.id, city.city_name, vacancy.position, vacancy_description, salary, employer_id, work_type.work_type_name, path_to_company_avatar, vacancy.created_at, vacancy.updated_at, 
		company.company_name, position_category.category_name, vacancy.compressed_image from vacancy
		left join work_type on vacancy.work_type_id=work_type.id left join city on vacancy.city_id=city.id
		left join employer on vacancy.employer_id=employer.id left join company on employer.company_name_id=company.id
		left join position_category on vacancy.position_category_id = position_category.id
		where vacancy.id = $1`, ID)
	var oneVacancy dto.JSONVacancyWithNull

	err := row.Scan(
		&oneVacancy.ID,
		&oneVacancy.Location,
		&oneVacancy.Position,
		&oneVacancy.Description,
		&oneVacancy.Salary,
		&oneVacancy.EmployerID,
		&oneVacancy.WorkType,
		&oneVacancy.Avatar,
		&oneVacancy.CreatedAt,
		&oneVacancy.UpdatedAt,
		&oneVacancy.CompanyName,
		&oneVacancy.PositionCategoryName,
		&oneVacancy.CompressedAvatar,
	)
	VacancyOk := dto.JSONVacancy{
		ID:                   oneVacancy.ID,
		EmployerID:           oneVacancy.EmployerID,
		Salary:               oneVacancy.Salary,
		Position:             oneVacancy.Position,
		Location:             oneVacancy.Location,
		Description:          oneVacancy.Description,
		WorkType:             oneVacancy.WorkType,
		Avatar:               oneVacancy.Avatar,
		CompanyName:          oneVacancy.CompanyName,
		PositionCategoryName: oneVacancy.PositionCategoryName.String,
		CreatedAt:            oneVacancy.CreatedAt,
		UpdatedAt:            oneVacancy.UpdatedAt,
		CompressedAvatar:     oneVacancy.CompressedAvatar.String,
	}
	if err != nil {
		return nil, err
	}
	return &VacancyOk, err
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
	var PositionCategoryID int
	if updatedVacancy.PositionCategoryName != "" {
		row = s.db.QueryRow(`select id from position_category where category_name=$1`, updatedVacancy.PositionCategoryName)
		err := row.Scan(&PositionCategoryID)
		if err != nil {
			return nil, err
		}
	}
	if updatedVacancy.Avatar != "" {
		if updatedVacancy.PositionCategoryName == "" {
			row = s.db.QueryRow(`update vacancy
				set employer_id = $1, salary = $2, position = $3, city_id = $4, vacancy_description = $5,
				work_type_id = $6, path_to_company_avatar = $7, compressed_image=$8 where id=$9 returning id, position, vacancy_description, 
				salary, employer_id, path_to_company_avatar, created_at, updated_at, compressed_image`,
				updatedVacancy.EmployerID, updatedVacancy.Salary, updatedVacancy.Position,
				CityId, updatedVacancy.Description, WorkTypeID, updatedVacancy.Avatar, updatedVacancy.CompressedAvatar, ID)
		} else {
			row = s.db.QueryRow(`update vacancy
				set employer_id = $1, salary = $2, position = $3, city_id = $4, vacancy_description = $5,
				work_type_id = $6, path_to_company_avatar = $7, position_category_id = $8, compressed_image=$9 where id=$10 returning id, position, vacancy_description, 
				salary, employer_id, path_to_company_avatar, created_at, updated_at, compressed_image`,
				updatedVacancy.EmployerID, updatedVacancy.Salary, updatedVacancy.Position,
				CityId, updatedVacancy.Description, WorkTypeID, updatedVacancy.Avatar, PositionCategoryID, updatedVacancy.CompressedAvatar, ID)
		}
	} else {
		if updatedVacancy.PositionCategoryName == "" {
			row = s.db.QueryRow(`update vacancy
				set employer_id = $1, salary = $2, position = $3, city_id = $4, vacancy_description = $5,
				work_type_id = $6 where id=$7 returning id, position, vacancy_description, 
				salary, employer_id, path_to_company_avatar, created_at, updated_at, compressed_image`,
				updatedVacancy.EmployerID, updatedVacancy.Salary, updatedVacancy.Position,
				CityId, updatedVacancy.Description, WorkTypeID, ID)
		} else {
			row = s.db.QueryRow(`update vacancy
				set employer_id = $1, salary = $2, position = $3, city_id = $4, vacancy_description = $5,
				work_type_id = $6, position_category_id = $7 where id=$8 returning id, position, vacancy_description, 
				salary, employer_id, path_to_company_avatar, created_at, updated_at, compressed_image`,
				updatedVacancy.EmployerID, updatedVacancy.Salary, updatedVacancy.Position,
				CityId, updatedVacancy.Description, WorkTypeID, PositionCategoryID, ID)
		}
	}

	var oneVacancy dto.JSONVacancyWithNull

	err := row.Scan(
		&oneVacancy.ID,
		&oneVacancy.Position,
		&oneVacancy.Description,
		&oneVacancy.Salary,
		&oneVacancy.EmployerID,
		&oneVacancy.Avatar,
		&oneVacancy.CreatedAt,
		&oneVacancy.UpdatedAt,
		&oneVacancy.CompressedAvatar,
	)
	VacancyOk := dto.JSONVacancy{
		ID:               oneVacancy.ID,
		EmployerID:       oneVacancy.EmployerID,
		Salary:           oneVacancy.Salary,
		Position:         oneVacancy.Position,
		Description:      oneVacancy.Description,
		Avatar:           oneVacancy.Avatar,
		CreatedAt:        oneVacancy.CreatedAt,
		UpdatedAt:        oneVacancy.UpdatedAt,
		CompressedAvatar: oneVacancy.CompressedAvatar.String,
	}
	if err != nil {
		return nil, err
	}
	row = s.db.QueryRow(`select company.company_name from employer left join company on company.id = employer.company_name_id where employer.id=$1`, oneVacancy.EmployerID)
	err = row.Scan(
		&VacancyOk.CompanyName)
	VacancyOk.WorkType = updatedVacancy.WorkType
	VacancyOk.Location = updatedVacancy.Location
	VacancyOk.PositionCategoryName = updatedVacancy.PositionCategoryName
	if err != nil {
		return nil, err
	}
	return &VacancyOk, err
}

func (s *PostgreSQLVacanciesStorage) Delete(ID uint64) error {
	_, err := s.db.Exec(`delete from vacancy where id = $1`, ID)
	return err
}

func (s *PostgreSQLVacanciesStorage) Subscribe(ID uint64, applicantID uint64) error {
	//TODO logger
	fmt.Println(ID, applicantID)
	_, err := s.db.Exec(`insert into vacancy_subscriber (vacancy_id, applicant_id) VALUES ($1, $2)`, ID, applicantID)
	return err
}

func (s *PostgreSQLVacanciesStorage) GetSubscriptionStatus(ID uint64, applicantID uint64) (bool, error) {
	var rowID uint64
	row := s.db.QueryRow(`select applicant_id from vacancy_subscriber where applicant_id=$1 and vacancy_id=$2`, applicantID, ID)
	if err := row.Scan(&rowID); err != nil {
		return false, nil
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
		contacts, education, email, password_hash , applicant.created_at, applicant.updated_at, applicant.compressed_image
		from vacancy_subscriber	left join applicant on applicant.id = applicant_id
		left join city on city.id = applicant.id where vacancy_subscriber.vacancy_id = $1`, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var applicantWithNull dto.ApplicantWithNull
		if err := rows.Scan(&applicantWithNull.ID, &applicantWithNull.FirstName, &applicantWithNull.LastName, &applicantWithNull.CityName,
			&applicantWithNull.BirthDate, &applicantWithNull.PathToProfileAvatar, &applicantWithNull.Contacts, &applicantWithNull.Education,
			&applicantWithNull.Email, &applicantWithNull.PasswordHash, &applicantWithNull.CreatedAt, &applicantWithNull.UpdatedAt, &applicantWithNull.CompressedAvatar); err != nil {
			return nil, err
		}
		oneApplicant := models.Applicant{
			ID:                  applicantWithNull.ID,
			FirstName:           applicantWithNull.FirstName,
			LastName:            applicantWithNull.LastName,
			CityName:            applicantWithNull.CityName.String,
			BirthDate:           applicantWithNull.BirthDate,
			PathToProfileAvatar: applicantWithNull.PathToProfileAvatar,
			Contacts:            applicantWithNull.Contacts.String,
			Education:           applicantWithNull.Education.String,
			Email:               applicantWithNull.Email,
			PasswordHash:        applicantWithNull.PasswordHash,
			CreatedAt:           applicantWithNull.CreatedAt,
			UpdatedAt:           applicantWithNull.UpdatedAt,
			CompressedAvatar:    applicantWithNull.CompressedAvatar.String,
		}

		Applicants = append(Applicants, &oneApplicant)
		//TODO logger
		fmt.Println(oneApplicant)
	}
	return Applicants, nil
}

func (s *PostgreSQLVacanciesStorage) Unsubscribe(ID uint64, applicantID uint64) error {
	_, err := s.db.Exec(`delete from vacancy_subscriber where applicant_id=$1 and vacancy_id=$2`, applicantID, ID)
	return err
}

func (s *PostgreSQLVacanciesStorage) GetApplicantFavoriteVacancies(applicantID uint64) ([]*dto.JSONVacancy, error) {

	Vacancies := make([]*dto.JSONVacancy, 0)

	rows, err := s.db.Query(`select vacancy.id, city.city_name, vacancy.position, vacancy_description, salary, employer_id, work_type.work_type_name, path_to_company_avatar, vacancy.created_at, vacancy.updated_at, 
		company.company_name, position_category.category_name, vacancy.compressed_image from vacancy
		left join work_type on vacancy.work_type_id=work_type.id left join city on vacancy.city_id=city.id
		left join employer on vacancy.employer_id=employer.id left join company on employer.company_name_id=company.id
		left join position_category on vacancy.position_category_id = position_category.id
		left join favorite_vacancy on favorite_vacancy.vacancy_id = vacancy.id
		left join applicant on applicant.id = favorite_vacancy.applicant_id 
		where applicant.id = $1`, applicantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var Vacancy dto.JSONVacancyWithNull
		if err := rows.Scan(
			&Vacancy.ID,
			&Vacancy.Location,
			&Vacancy.Position,
			&Vacancy.Description,
			&Vacancy.Salary,
			&Vacancy.EmployerID,
			&Vacancy.WorkType,
			&Vacancy.Avatar,
			&Vacancy.CreatedAt,
			&Vacancy.UpdatedAt,
			&Vacancy.CompanyName,
			&Vacancy.PositionCategoryName,
			&Vacancy.CompressedAvatar); err != nil {
			return nil, err
		}
		VacancyOk := dto.JSONVacancy{
			ID:                   Vacancy.ID,
			EmployerID:           Vacancy.EmployerID,
			Salary:               Vacancy.Salary,
			Position:             Vacancy.Position,
			Location:             Vacancy.Location,
			Description:          Vacancy.Description,
			WorkType:             Vacancy.WorkType,
			Avatar:               Vacancy.Avatar,
			CompanyName:          Vacancy.CompanyName,
			PositionCategoryName: Vacancy.PositionCategoryName.String,
			CompressedAvatar:     Vacancy.CompressedAvatar.String,
			CreatedAt:            Vacancy.CreatedAt,
			UpdatedAt:            Vacancy.UpdatedAt,
		}
		Vacancies = append(Vacancies, &VacancyOk)
		//TODO logger
		fmt.Println(VacancyOk)
	}
	return Vacancies, nil
}

func (s *PostgreSQLVacanciesStorage) MakeFavorite(ID uint64, applicantID uint64) error {
	//TODO logger
	fmt.Println(ID, applicantID)
	_, err := s.db.Exec(`insert into favorite_vacancy (applicant_id, vacancy_id) VALUES ($1, $2)`, applicantID, ID)
	return err
}

func (s *PostgreSQLVacanciesStorage) Unfavorite(ID uint64, applicantID uint64) error {
	//TODO logger
	fmt.Println(ID, applicantID)
	_, err := s.db.Exec(`delete from favorite_vacancy where applicant_id = $1 and vacancy_id = $2`, applicantID, ID)
	return err
}
