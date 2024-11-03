package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

// PostgreSQLBoardStorage
// Хранилище досок в PostgreSQL
type PostgreSQLApplicantStorage struct {
	db *sql.DB
}

func NewApplicantStorage(db *sql.DB) *PostgreSQLApplicantStorage {
	return &PostgreSQLApplicantStorage{
		db: db,
	}
}

// GetByID
// находит доску и связанные с ней списки и задания по id
// или возвращает ошибки ...
func (s *PostgreSQLApplicantStorage) GetByID(id uint64) (*models.Applicant, error) {
	// funcName := "PostgreSQLapplicant.GetById"
	// logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
	// if !ok {
	// 	fmt.Printf("function %s: can't get logger from context\n", funcName)
	// }

	//requestID := ctx.Value(dto.RequestIDKey).(uuid.UUID)

	//logger.DebugFmt("Built query\n\t"+boardSql+"\nwith args\n\t"+fmt.Sprintf("%+v", args), requestID.String(), funcName, nodeName)
	row := s.db.QueryRow(`select applicant.id, first_name, last_name, city.city_name, birth_date, path_to_profile_avatar, contacts, 
	education, email, password_hash, applicant.created_at, applicant.updated_at 
	from applicant left join city on applicant.city_id = city.id where applicant.id = $1`, id)

	var applicantWithNull dto.ApplicantWithNull
	err := row.Scan(
		&applicantWithNull.ID,
		&applicantWithNull.FirstName,
		&applicantWithNull.LastName,
		&applicantWithNull.CityName,
		&applicantWithNull.BirthDate,
		&applicantWithNull.PathToProfileAvatar,
		&applicantWithNull.Contacts,
		&applicantWithNull.Education,
		&applicantWithNull.Email,
		&applicantWithNull.PasswordHash,
		&applicantWithNull.CreatedAt,
		&applicantWithNull.UpdatedAt,
	)
	applicant := models.Applicant{
		ID:                  applicantWithNull.ID,
		FirstName:           applicantWithNull.FirstName,
		LastName:            applicantWithNull.LastName,
		CityName:            applicantWithNull.CityName,
		BirthDate:           applicantWithNull.BirthDate,
		PathToProfileAvatar: applicantWithNull.PathToProfileAvatar,
		Contacts:            applicantWithNull.Contacts.String,
		Education:           applicantWithNull.Education.String,
		Email:               applicantWithNull.Email,
		PasswordHash:        applicantWithNull.PasswordHash,
		CreatedAt:           applicantWithNull.CreatedAt,
		UpdatedAt:           applicantWithNull.UpdatedAt,
	}

	if err != nil {
		return nil, err
	}
	//logger.DebugFmt(fmt.Sprintf("%+v", board), requestID.String(), funcName, nodeName)

	return &applicant, nil
}

func (s *PostgreSQLApplicantStorage) GetByEmail(email string) (*models.Applicant, error) {
	//log.Println("Looking for user with login", login.Value)

	//log.Println("Built query:", sql, "\nwith args:", args)

	row := s.db.QueryRow(`select applicant.id, first_name, last_name, city.city_name, birth_date, path_to_profile_avatar, contacts, 
	education, email, password_hash, applicant.created_at, applicant.updated_at 
	from applicant left join city on applicant.city_id = city.id where applicant.email=$1`, email)

	var applicantWithNull dto.ApplicantWithNull
	err := row.Scan(
		&applicantWithNull.ID,
		&applicantWithNull.FirstName,
		&applicantWithNull.LastName,
		&applicantWithNull.CityName,
		&applicantWithNull.BirthDate,
		&applicantWithNull.PathToProfileAvatar,
		&applicantWithNull.Contacts,
		&applicantWithNull.Education,
		&applicantWithNull.Email,
		&applicantWithNull.PasswordHash,
		&applicantWithNull.CreatedAt,
		&applicantWithNull.UpdatedAt,
	)
	applicant := models.Applicant{
		ID:                  applicantWithNull.ID,
		FirstName:           applicantWithNull.FirstName,
		LastName:            applicantWithNull.LastName,
		CityName:            applicantWithNull.CityName,
		BirthDate:           applicantWithNull.BirthDate,
		PathToProfileAvatar: applicantWithNull.PathToProfileAvatar,
		Contacts:            applicantWithNull.Contacts.String,
		Education:           applicantWithNull.Education.String,
		Email:               applicantWithNull.Email,
		PasswordHash:        applicantWithNull.PasswordHash,
		CreatedAt:           applicantWithNull.CreatedAt,
		UpdatedAt:           applicantWithNull.UpdatedAt,
	}
	//log.Println(user)
	//log.Println(err)
	if err != nil {
		fmt.Println("Error", err.Error())
		return nil, err
	}
	return &applicant, nil
}

func (s *PostgreSQLApplicantStorage) Create(applicantInput *dto.ApplicantInput) (*models.Applicant, error) {

	ApplicantCityId := 1 // do it and profileavatar
	if applicantInput.PathToProfileAvatar == "" {
		applicantInput.PathToProfileAvatar = "static/default_profile.png"
	}
	if applicantInput.Contacts == "" {
		applicantInput.Contacts = "no contacts yet"
	}
	_, err := s.db.Exec("insert into applicant (first_name, last_name, city_id, birth_date, path_to_profile_avatar, contacts, education, email, password_hash) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		applicantInput.FirstName, applicantInput.LastName, ApplicantCityId, applicantInput.BirthDate, applicantInput.PathToProfileAvatar, applicantInput.Contacts, applicantInput.Education, applicantInput.Email, applicantInput.Password)

	if err != nil {
		return nil, err
	}

	applicant, err := s.GetByEmail(applicantInput.Email)
	if err != nil {
		return nil, err
	}

	return applicant, nil
}
