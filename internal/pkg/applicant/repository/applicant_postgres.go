package repository

import (
	"context"
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
	education, email, password_hash, applicant.created_at, applicant.updated_at, applicant.compressed_image
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
		&applicantWithNull.CompressedAvatar,
	)
	applicant := models.Applicant{
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
	education, email, password_hash, applicant.created_at, applicant.updated_at, applicant.compressed_image
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
		&applicantWithNull.CompressedAvatar,
	)
	applicant := models.Applicant{
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
	//log.Println(user)
	//log.Println(err)
	if err != nil {
		fmt.Println("Error", err.Error())
		return nil, err
	}
	return &applicant, nil
}

func (s *PostgreSQLApplicantStorage) Create(applicantInput *dto.ApplicantInput) (*models.Applicant, error) {
	row := s.db.QueryRow(`insert into applicant (first_name, last_name, birth_date, education, email, password_hash) VALUES ($1, $2, $3, $4, $5, $6)
	returning id, first_name, last_name, birth_date, path_to_profile_avatar, contacts, education, email, password_hash, created_at, updated_at, compressed_image`,
		applicantInput.FirstName, applicantInput.LastName, applicantInput.BirthDate, applicantInput.Education, applicantInput.Email, applicantInput.Password)

	var applicantWithNull dto.ApplicantWithNull
	err := row.Scan(
		&applicantWithNull.ID,
		&applicantWithNull.FirstName,
		&applicantWithNull.LastName,
		&applicantWithNull.BirthDate,
		&applicantWithNull.PathToProfileAvatar,
		&applicantWithNull.Contacts,
		&applicantWithNull.Education,
		&applicantWithNull.Email,
		&applicantWithNull.PasswordHash,
		&applicantWithNull.CreatedAt,
		&applicantWithNull.UpdatedAt,
		&applicantWithNull.CompressedAvatar,
	)
	applicant := models.Applicant{
		ID:                  applicantWithNull.ID,
		FirstName:           applicantWithNull.FirstName,
		LastName:            applicantWithNull.LastName,
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
	applicant.CityName = applicantInput.CityName

	return &applicant, err
}

func (s *PostgreSQLApplicantStorage) Update(ID uint64, newApplicantData *dto.JSONUpdateApplicantProfile) (*models.Applicant, error) {

	var CityId uint64
	row := s.db.QueryRow(`select id from city where city_name=$1`, newApplicantData.City)
	if err := row.Scan(&CityId); err != nil {
		switch err {
		case sql.ErrNoRows:
			row = s.db.QueryRow(`insert into city (city_name) VALUES ($1) returning id`, newApplicantData.City)
			err = row.Scan(&CityId)
			if err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}
	if newApplicantData.Avatar == "" {
		row = s.db.QueryRow(`update applicant
		set first_name = $1, last_name = $2, city_id = $3, birth_date=$4,
		contacts = $5, education = $6 where id=$7 returning id, first_name, last_name, birth_date, path_to_profile_avatar, contacts, education, email, password_hash, created_at, updated_at, compressed_image`,
			newApplicantData.FirstName, newApplicantData.LastName, CityId, newApplicantData.BirthDate, newApplicantData.Contacts, newApplicantData.Education, ID)
	} else {
		row = s.db.QueryRow(`update applicant
		set first_name = $1, last_name = $2, city_id = $3, birth_date=$4,
		contacts = $5, education = $6, path_to_profile_avatar=$7, compressed_image=$8 where id=$9 returning id, first_name, last_name, birth_date, path_to_profile_avatar, contacts, education, email, password_hash, created_at, updated_at, compressed_image`,
			newApplicantData.FirstName, newApplicantData.LastName, CityId, newApplicantData.BirthDate, newApplicantData.Contacts, newApplicantData.Education, newApplicantData.Avatar, newApplicantData.CompressedAvatar, ID)
	}
	var applicantWithNull dto.ApplicantWithNull
	err := row.Scan(
		&applicantWithNull.ID,
		&applicantWithNull.FirstName,
		&applicantWithNull.LastName,
		&applicantWithNull.BirthDate,
		&applicantWithNull.PathToProfileAvatar,
		&applicantWithNull.Contacts,
		&applicantWithNull.Education,
		&applicantWithNull.Email,
		&applicantWithNull.PasswordHash,
		&applicantWithNull.CreatedAt,
		&applicantWithNull.UpdatedAt,
		&applicantWithNull.CompressedAvatar,
	)
	applicant := models.Applicant{
		ID:                  applicantWithNull.ID,
		FirstName:           applicantWithNull.FirstName,
		LastName:            applicantWithNull.LastName,
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
	applicant.CityName = newApplicantData.City

	return &applicant, err
}

func (s *PostgreSQLApplicantStorage) GetAllCities(ctx context.Context, namePart string) ([]string, error) {
	Cities := make([]string, 0)
	rows, err := s.db.Query(`select city.city_name from city where city.city_name like $1`, "%"+namePart+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var oneCity string
		if err := rows.Scan(&oneCity); err != nil {
			return nil, err
		}
		Cities = append(Cities, oneCity)
		fmt.Println(oneCity)
	}

	return Cities, nil
}
