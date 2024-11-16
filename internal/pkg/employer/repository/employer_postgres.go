package repository

import (
	"database/sql"

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
	//log.Println(user)
	//log.Println(err)
	return &employer, err
}

func (s *PostgreSQLEmployerStorage) Create(employerInput *dto.EmployerInput) (*models.Employer, error) {

	var CompanyNameId int
	row := s.db.QueryRow(`select id from company where company_name = $1`, employerInput.CompanyName)
	if err := row.Scan(&CompanyNameId); err != nil {
		switch err {
		case sql.ErrNoRows:
			row = s.db.QueryRow(`insert into company (company_name) VALUES ($1) returning id`, employerInput.CompanyName)
			err = row.Scan(&CompanyNameId)
			if err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}
	row = s.db.QueryRow(`insert into employer (first_name, last_name, position, company_name_id, company_description, company_website, email, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) returning id, first_name, last_name, position, company_description, 
		company_website, path_to_profile_avatar, contacts, email, password_hash, employer.created_at, employer.updated_at`,
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
	)
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
	}
	employer.CityName = employerInput.CityName
	employer.CompanyName = employerInput.CompanyName
	//log.Println(user)
	//log.Println(err)
	return &employer, err
}

func (s *PostgreSQLEmployerStorage) Update(ID uint64, newEmployerData *dto.JSONUpdateEmployerProfile) (*models.Employer, error) {

	var CityId int
	row := s.db.QueryRow(`select id from city where city_name=$1`, newEmployerData.City)
	if err := row.Scan(&CityId); err != nil {
		switch err {
		case sql.ErrNoRows:
			row = s.db.QueryRow(`insert into city (city_name) VALUES ($1) returning id`, newEmployerData.City)
			err = row.Scan(&CityId)
			if err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}
	row = s.db.QueryRow(`update employer
		set first_name = $1, last_name = $2, city_id = $3,
		contacts = $4 where id=$5 returning id, first_name, last_name, position, company_description, 
		company_website, path_to_profile_avatar, contacts, email, password_hash, employer.created_at, employer.updated_at`,
		newEmployerData.FirstName, newEmployerData.LastName, CityId, newEmployerData.Contacts, ID)
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
	)
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
	}
	employer.CityName = newEmployerData.City
	//log.Println(user)
	//log.Println(err)
	return &employer, err
}
