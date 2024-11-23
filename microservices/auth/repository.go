package auth

import (
	"database/sql"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

type AuthorizationRepository struct {
	db *sql.DB
}

func NewAuthorizationRepository(db *sql.DB) *AuthorizationRepository {
	return &AuthorizationRepository{
		db: db,
	}
}

// TODO: передавать request_id через контекст
func (r *AuthorizationRepository) GetUser(userType, email string) (*User, error) {
	var row *sql.Row
	switch userType {
	case dto.UserTypeApplicant:
		row = r.db.QueryRow(
			`SELECT id, password_hash
			 FROM employer
			 WHERE email = $1`, email)
	case dto.UserTypeEmployer:
		row = r.db.QueryRow(
			`SELECT id, password_hash
			FROM employer
			WHERE email = $1`, email)
	default:
		return nil, BadUserType
	}

	user := new(User)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, NoUserExist
	} else if err != nil {
		return nil, err
	}

	user.Email = email
	user.UserType = userType
	return user, nil
}

func (r *AuthorizationRepository) CreateSession(userID uint64, sessionToken string) error {
	// Temporary all sessions will write into applicant_session
	// Так как у нас тип пользователя кодируется 1 символом строки, то нет смысла использовать 2 таблицы сессий
	_, err := r.db.Exec(
		`INSERT INTO applicant_session (applicant_id, session_token)
		 VALUES ($1, $2)`, userID, sessionToken)
	return err
}

func (r *AuthorizationRepository) GetUserIdBySession(sessionToken string) (uint64, error) {
	row := r.db.QueryRow(
		`SELECT applicant_id
		 FROM applicant_session
		 WHERE session_token = $1`, sessionToken)
	var userID uint64
	err := row.Scan(&userID)
	if err == sql.ErrNoRows {
		return 0, NoUserExist
	} else if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *AuthorizationRepository) DeleteSession(sessionToken string) error {
	_, err := r.db.Exec(
		`DELETE FROM applicant_session
		 WHERE session_token = $1`, sessionToken)
	return err
}
