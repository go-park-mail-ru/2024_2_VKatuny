package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

// PostgreSQLBoardStorage
// Хранилище досок в PostgreSQL
type PostgreSQLEmployerStorage struct {
	db *sql.DB
	logger *logrus.Entry
}

func NewEmployerStorage(db *sql.DB, logger *logrus.Logger) *PostgreSQLEmployerStorage {
	return &PostgreSQLEmployerStorage{
		db: db,
		logger: logrus.NewEntry(logger),
	}
}

// GetByID
// находит доску и связанные с ней списки и задания по id
// или возвращает ошибки ...
func (s *PostgreSQLEmployerStorage) GetByID(ctx context.Context, id uint64) (*models.Employer, error) {
	funcName := "PostgreSQLEmployerStorage.GetById"
	s.logger = utils.SetLoggerRequestID(ctx, s.logger)
	s.logger.Debugf("%s: entering", funcName)

	row := s.db.QueryRow(`select employer.id, first_name, last_name, city.city_name, position, company.company_name, company_description, company_website, path_to_profile_avatar, contacts, 
	email, password_hash, employer.created_at, employer.updated_at 
	from employer left join city on employer.city_id = city.id left join company on employer.company_name_id = company.id where employer.id = $1`, id)

	var employerWithNull dto.EmployerWithNull
	err := row.Scan(
		&employerWithNull.ID,
		&employerWithNull.FirstName,
		&employerWithNull.LastName,
		&employerWithNull.CityName,
		&employerWithNull.Position,
		&employerWithNull.CompanyName,
		&employerWithNull.CompanyDescription,
		&employerWithNull.CompanyWebsite,
		&employerWithNull.PathToProfileAvatar,
		&employerWithNull.Contacts,
		&employerWithNull.Email,
		&employerWithNull.PasswordHash,
		&employerWithNull.CreatedAt,
		&employerWithNull.UpdatedAt,
	)
	s.logger.Debugf("%s: got from sql db %v", funcName, employerWithNull)

	employer := models.Employer{
		ID:                  employerWithNull.ID,
		FirstName:           employerWithNull.FirstName,
		LastName:            employerWithNull.LastName,
		CityName:            employerWithNull.CityName.String,
		Position:            employerWithNull.Position,
		CompanyName:         employerWithNull.CompanyName,
		CompanyDescription:  employerWithNull.CompanyDescription,
		CompanyWebsite:      employerWithNull.CompanyWebsite,
		PathToProfileAvatar: employerWithNull.PathToProfileAvatar,
		Contacts:            employerWithNull.Contacts.String,
		Email:               employerWithNull.Email,
		PasswordHash:        employerWithNull.PasswordHash,
		CreatedAt:           employerWithNull.CreatedAt,
		UpdatedAt:           employerWithNull.UpdatedAt,
	}

	if err != nil {
		s.logger.Errorf("%s: got error %v", funcName, err)
		return nil, err
	}

	return &employer, nil
}

func (s *PostgreSQLEmployerStorage) GetByEmail(ctx context.Context, email string) (*models.Employer, error) {
	funcName := "PostgreSQLEmployerStorage.GetByEmail"
	s.logger = utils.SetLoggerRequestID(ctx, s.logger)
	s.logger.Debugf("%s: entering", funcName)

	row := s.db.QueryRow(`select employer.id, first_name, last_name, city.city_name, position, company.company_name, company_description, company_website, path_to_profile_avatar, contacts, 
	email, password_hash, employer.created_at, employer.updated_at 
	from employer left join city on employer.city_id = city.id left join company on employer.company_name_id = company.id where employer.email = $1`, email)

	var employerWithNull dto.EmployerWithNull
	err := row.Scan(
		&employerWithNull.ID,
		&employerWithNull.FirstName,
		&employerWithNull.LastName,
		&employerWithNull.CityName,
		&employerWithNull.Position,
		&employerWithNull.CompanyName,
		&employerWithNull.CompanyDescription,
		&employerWithNull.CompanyWebsite,
		&employerWithNull.PathToProfileAvatar,
		&employerWithNull.Contacts,
		&employerWithNull.Email,
		&employerWithNull.PasswordHash,
		&employerWithNull.CreatedAt,
		&employerWithNull.UpdatedAt,
	)
	s.logger.Debugf("%s: got from sql db %v", funcName, employerWithNull)

	employer := models.Employer{
		ID:                  employerWithNull.ID,
		FirstName:           employerWithNull.FirstName,
		LastName:            employerWithNull.LastName,
		CityName:            employerWithNull.CityName.String,
		Position:            employerWithNull.Position,
		CompanyName:         employerWithNull.CompanyName,
		CompanyDescription:  employerWithNull.CompanyDescription,
		CompanyWebsite:      employerWithNull.CompanyWebsite,
		PathToProfileAvatar: employerWithNull.PathToProfileAvatar,
		Contacts:            employerWithNull.Contacts.String,
		Email:               employerWithNull.Email,
		PasswordHash:        employerWithNull.PasswordHash,
		CreatedAt:           employerWithNull.CreatedAt,
		UpdatedAt:           employerWithNull.UpdatedAt,
	}

	if err != nil {
		s.logger.Errorf("%s: got error %v", funcName, err)
		return nil, err
	}
	return &employer, nil
}

func (s *PostgreSQLEmployerStorage) Create(ctx context.Context, employerInput *dto.EmployerInput) (*models.Employer, error) {
	funcName := "PostgreSQLEmployerStorage.Create"
	s.logger = utils.SetLoggerRequestID(ctx, s.logger)
	s.logger.Debugf("%s: entering", funcName)

	var CompanyNameId int
	row := s.db.QueryRow(`select id from company where company_name = $1`, employerInput.CompanyName)
	if err := row.Scan(&CompanyNameId); err != nil {
		switch err {
		case sql.ErrNoRows:
			s.logger.Debugf("%s: got empty result: %s", funcName, sql.ErrNoRows.Error())
			row = s.db.QueryRow(`insert into company (company_name) VALUES ($1) returning id`, employerInput.CompanyName)
			err = row.Scan(&CompanyNameId)
			if err != nil {
				s.logger.Errorf("%s: got error %v", funcName, err)
				return nil, err
			}
		default:
			s.logger.Errorf("%s: got error %v", funcName, err)
			return nil, err
		}
	}
	row = s.db.QueryRow(`insert into employer (first_name, last_name, position, company_name_id, company_description, company_website, email, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) returning id, first_name, last_name, position, company_description, 
		company_website, path_to_profile_avatar, contacts, email, password_hash, employer.created_at, employer.updated_at, compressed_image`,
		employerInput.FirstName, employerInput.LastName, employerInput.Position, CompanyNameId, employerInput.CompanyDescription,
		employerInput.CompanyWebsite, employerInput.Email, employerInput.Password)

	var employerWithNull dto.EmployerWithNull
	err := row.Scan(
		&employerWithNull.ID,
		&employerWithNull.FirstName,
		&employerWithNull.LastName,
		&employerWithNull.Position,
		&employerWithNull.CompanyDescription,
		&employerWithNull.CompanyWebsite,
		&employerWithNull.PathToProfileAvatar,
		&employerWithNull.Contacts,
		&employerWithNull.Email,
		&employerWithNull.PasswordHash,
		&employerWithNull.CreatedAt,
		&employerWithNull.UpdatedAt,
		&employerWithNull.CompressedAvatar,
	)
	s.logger.Debugf("%s: got from sql db %v", funcName, employerWithNull)
	employer := models.Employer{
		ID:                  employerWithNull.ID,
		FirstName:           employerWithNull.FirstName,
		LastName:            employerWithNull.LastName,
		Position:            employerWithNull.Position,
		CompanyDescription:  employerWithNull.CompanyDescription,
		CompanyWebsite:      employerWithNull.CompanyWebsite,
		PathToProfileAvatar: employerWithNull.PathToProfileAvatar,
		Contacts:            employerWithNull.Contacts.String,
		Email:               employerWithNull.Email,
		PasswordHash:        employerWithNull.PasswordHash,
		CreatedAt:           employerWithNull.CreatedAt,
		UpdatedAt:           employerWithNull.UpdatedAt,
		CompressedAvatar:    employerWithNull.CompressedAvatar.String,
	}
	employer.CityName = employerInput.CityName
	employer.CompanyName = employerInput.CompanyName

	if err != nil {
		s.logger.Errorf("%s: got error %v", funcName, err)
		return nil, err
	}

	return &employer, nil
}

func (s *PostgreSQLEmployerStorage) Update(ctx context.Context, ID uint64, newEmployerData *dto.JSONUpdateEmployerProfile) (*models.Employer, error) {
	funcName := "PostgreSQLEmployerStorage.Update"
	s.logger = utils.SetLoggerRequestID(ctx, s.logger)
	s.logger.Debugf("%s: entering", funcName)

	var CityId int
	row := s.db.QueryRow(`select id from city where city_name=$1`, newEmployerData.City)
	if err := row.Scan(&CityId); err != nil {
		switch err {
		case sql.ErrNoRows:
			s.logger.Debugf("%s: got empty result: %s", funcName, sql.ErrNoRows.Error())
			row = s.db.QueryRow(`insert into city (city_name) VALUES ($1) returning id`, newEmployerData.City)
			err = row.Scan(&CityId)
			if err != nil {
				s.logger.Errorf("%s: got error %v", funcName, err)
				return nil, err
			}
		default:
			s.logger.Errorf("%s: got error %v", funcName, err)
			return nil, err
		}
	}
	if newEmployerData.Avatar == "" {
		row = s.db.QueryRow(`update employer
		set first_name = $1, last_name = $2, city_id = $3,
		contacts = $4 where id=$5 returning id, first_name, last_name, position, company_description, 
		company_website, path_to_profile_avatar, contacts, email, password_hash, employer.created_at, employer.updated_at, employer.compressed_image`,
			newEmployerData.FirstName, newEmployerData.LastName, CityId, newEmployerData.Contacts, ID)
	} else {
		row = s.db.QueryRow(`update employer
		set first_name = $1, last_name = $2, city_id = $3,
		contacts = $4, path_to_profile_avatar=$5, compressed_image=$6 where id=$7 returning id, first_name, last_name, position, company_description, 
		company_website, path_to_profile_avatar, contacts, email, password_hash, employer.created_at, employer.updated_at, employer.compressed_image`,
			newEmployerData.FirstName, newEmployerData.LastName, CityId, newEmployerData.Contacts, newEmployerData.Avatar, newEmployerData.CompressedAvatar, ID)
	}
	var employerWithNull dto.EmployerWithNull
	err := row.Scan(
		&employerWithNull.ID,
		&employerWithNull.FirstName,
		&employerWithNull.LastName,
		&employerWithNull.Position,
		&employerWithNull.CompanyDescription,
		&employerWithNull.CompanyWebsite,
		&employerWithNull.PathToProfileAvatar,
		&employerWithNull.Contacts,
		&employerWithNull.Email,
		&employerWithNull.PasswordHash,
		&employerWithNull.CreatedAt,
		&employerWithNull.UpdatedAt,
		&employerWithNull.CompressedAvatar,
	)
	s.logger.Debugf("%s: got from sql db %v", funcName, employerWithNull)

	employer := models.Employer{
		ID:                  employerWithNull.ID,
		FirstName:           employerWithNull.FirstName,
		LastName:            employerWithNull.LastName,
		Position:            employerWithNull.Position,
		CompanyDescription:  employerWithNull.CompanyDescription,
		CompanyWebsite:      employerWithNull.CompanyWebsite,
		PathToProfileAvatar: employerWithNull.PathToProfileAvatar,
		Contacts:            employerWithNull.Contacts.String,
		Email:               employerWithNull.Email,
		PasswordHash:        employerWithNull.PasswordHash,
		CreatedAt:           employerWithNull.CreatedAt,
		UpdatedAt:           employerWithNull.UpdatedAt,
		CompressedAvatar:    employerWithNull.CompressedAvatar.String,
	}
	employer.CityName = newEmployerData.City
	
	if err != nil {
		s.logger.Errorf("%s: got error %v", funcName, err)
		return nil, err
	}

	return &employer, nil
}
