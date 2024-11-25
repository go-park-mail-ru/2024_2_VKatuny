package auth

import (
	"database/sql"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/sirupsen/logrus"
)

type AuthorizationRepository struct {
	db     *sql.DB
	logger *logrus.Entry
}

func NewAuthorizationRepository(logger *logrus.Logger, db *sql.DB) *AuthorizationRepository {
	return &AuthorizationRepository{
		logger: logrus.NewEntry(logger),
		db:     db,
	}
}

// TODO: передавать request_id через контекст
func (r *AuthorizationRepository) GetUser(userType, email string) (*User, error) {
	fn := "AuthorizationRepository.GetUser"
	r.logger.Debugf("%s: trying to gey user with user type: %s, email: %s", fn, userType, email)
	var row *sql.Row
	switch userType {
	case dto.UserTypeApplicant:
		row = r.db.QueryRow(
			`SELECT id, password_hash
			 FROM applicant
			 WHERE email = $1`, email)
	case dto.UserTypeEmployer:
		row = r.db.QueryRow(
			`SELECT id, password_hash
			FROM employer
			WHERE email = $1`, email)
	default:
		return nil, ErrBadUserType
	}

	user := new(User)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, ErrNoUserExist
	} else if err != nil {
		return nil, err
	}

	user.Email = email
	user.UserType = userType
	r.logger.Debugf("%s: got user: %v", fn, user)
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
		return 0, ErrNoUserExist
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
