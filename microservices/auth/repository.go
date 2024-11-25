package auth

import (
	"database/sql"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type AuthorizationRepository struct {
	db     *sql.DB
	redis  redis.Conn
	logger *logrus.Entry
}

func NewAuthorizationRepository(logger *logrus.Logger, db *sql.DB, redisConn redis.Conn) *AuthorizationRepository {
	return &AuthorizationRepository{
		logger: logrus.NewEntry(logger),
		db:     db,
		redis:  redisConn,
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
	fn := "AuthorizationRepository.CreateSession"
	_, err := r.redis.Do("SET", sessionToken, userID, "EX", sessionTTL)
	if err != nil {
		r.logger.Errorf("%s: redis got err: %s", fn, err)
		return err
	}
	r.logger.Debugf("%s: created session: %s for user: %d", fn, sessionToken, userID)
	return nil
}

func (r *AuthorizationRepository) GetUserIdBySession(sessionToken string) (uint64, error) {
	fn := "AuthorizationRepository.GetUserIdBySession"
	userID, err := redis.Uint64(r.redis.Do("GET", sessionToken))
	if err != nil {
		r.logger.Errorf("%s: redis got err: %s", fn, err)
		return 0, err
	}
	r.logger.Debugf("%s: got user %d from session:  %s", fn, userID, sessionToken)
	return userID, nil
}

func (r *AuthorizationRepository) DeleteSession(sessionToken string) error {
	fn := "AuthorizationRepository.DeleteSession"
	_, err := r.redis.Do("DEL", sessionToken)
	if err != nil {
		r.logger.Errorf("%s: redis got err: %s", fn, err)
		return err
	}
	r.logger.Debugf("%s: session deleted successfully: %s", fn, sessionToken)
	return nil
}
