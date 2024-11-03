package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

// PostgreSQLBoardStorage
// Хранилище досок в PostgreSQL
type PostgreSQLEmployerStorage struct {
	db *sql.DB
}

func NewEmployerStorage(db *sql.DB) *PostgreSQLEmployerStorage {
	return &PostgreSQLEmployerStorage{
		db: db,
	}
}

// GetByID
// находит доску и связанные с ней списки и задания по id
// или возвращает ошибки ...
func (s *PostgreSQLEmployerStorage) GetByID(id uint64) (*models.Employer, error) {
	// funcName := "PostgreSQLemployer.GetById"
	// logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
	// if !ok {
	// 	fmt.Printf("function %s: can't get logger from context\n", funcName)
	// }

	//requestID := ctx.Value(dto.RequestIDKey).(uuid.UUID)

	//logger.DebugFmt("Built query\n\t"+boardSql+"\nwith args\n\t"+fmt.Sprintf("%+v", args), requestID.String(), funcName, nodeName)
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
	employer := models.Employer{
		ID:                  employerWithNull.ID,
		FirstName:           employerWithNull.FirstName,
		LastName:            employerWithNull.LastName,
		CityName:            employerWithNull.CityName,
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
		return nil, err
	}
	//logger.DebugFmt(fmt.Sprintf("%+v", board), requestID.String(), funcName, nodeName)

	return &employer, nil
}

func (s *PostgreSQLEmployerStorage) GetByEmail(email string) (*models.Employer, error) {
	//log.Println("Looking for user with login", login.Value)

	//log.Println("Built query:", sql, "\nwith args:", args)

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
	employer := models.Employer{
		ID:                  employerWithNull.ID,
		FirstName:           employerWithNull.FirstName,
		LastName:            employerWithNull.LastName,
		CityName:            employerWithNull.CityName,
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
	//log.Println(user)
	//log.Println(err)
	if err != nil {
		fmt.Println("Error", err.Error())
		return nil, err
	}
	return &employer, nil
}

func (s *PostgreSQLEmployerStorage) Create(employerInput *dto.EmployerInput) (*models.Employer, error) {

	employerCityId := 1
	var CompanyNameId int
	row := s.db.QueryRow(`select id from company where company_name = $1`, employerInput.CompanyName)
	err := row.Scan(CompanyNameId)

	fmt.Println("!", CompanyNameId, "!", employerInput.CompanyName)
	CompanyNameId = 1
	if employerInput.PathToProfileAvatar == "" {
		employerInput.PathToProfileAvatar = "static/default_profile.png"
	}
	_, err = s.db.Exec(`insert into employer (first_name, last_name, city_id, position, company_name_id, company_description, company_website, path_to_profile_avatar, contacts, email, password_hash)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		employerInput.FirstName, employerInput.LastName, employerCityId, employerInput.Position, CompanyNameId, employerInput.CompanyDescription,
		employerInput.CompanyWebsite, employerInput.PathToProfileAvatar, employerInput.Contacts, employerInput.Email, employerInput.Password)

	if err != nil {
		return nil, err
	}

	employer, err := s.GetByEmail(employerInput.Email)
	if err != nil {
		return nil, err
	}

	return employer, nil
}
