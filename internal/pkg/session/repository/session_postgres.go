package repository

import (
	"database/sql"
)

// PostgreSQLBoardStorage
// Хранилище досок в PostgreSQL
type PostgreSQLEmployerSession struct {
	db *sql.DB
}

type PostgreSQLApplicantSession struct {
	db *sql.DB
}

func NewSessionStorage(db *sql.DB) (*PostgreSQLApplicantSession, *PostgreSQLEmployerSession) {
	return &PostgreSQLApplicantSession{
			db: db,
		}, &PostgreSQLEmployerSession{
			db: db,
		}
}
func (s *PostgreSQLEmployerSession) GetUserIdBySession(sessionId string) (uint64, error) {
	row := s.db.QueryRow(`select employer_id from employer_session  where session_token = $1`, sessionId)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}
func (s *PostgreSQLApplicantSession) GetUserIdBySession(sessionId string) (uint64, error) {
	row := s.db.QueryRow(`select applicant_id from applicant_session  where session_token = $1`, sessionId)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (s *PostgreSQLEmployerSession) Create(userId uint64, sessionID string) error {
	_, err := s.db.Exec(`insert into employer_session (employer_id, session_token) VALUES ($1, $2)`, userId, sessionID)
	return err
}
func (s *PostgreSQLApplicantSession) Create(userId uint64, sessionID string) error {
	_, err := s.db.Exec(`insert into applicant_session (applicant_id, session_token) VALUES ($1, $2)`, userId, sessionID)
	return err
}

func (s *PostgreSQLEmployerSession) Delete(sessionId string) error {
	_, err := s.db.Exec(`delete from employer_session where session_token = $1`, sessionId)
	return err
}
func (s *PostgreSQLApplicantSession) Delete(sessionId string) error {
	_, err := s.db.Exec(`delete from applicant_session where session_token = $1`, sessionId)
	return err
}
