// Package repository is a repository layer of session
package repository

import (
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// SessionApplicantRepo is a in-memory implementation for applicant session
type SessionApplicantRepo struct {
	lastID   uint64
	sessions map[string]*models.SessionApplicant
}

// SessionEmployerRepo is a in-memory implementation for employer session
type SessionEmployerRepo struct {
	lastID   uint64
	sessions map[string]*models.SessionEmployer
}

// NewRepo returns new in-memory repository for sessions
func NewRepo() (*SessionApplicantRepo, *SessionEmployerRepo) {
	return &SessionApplicantRepo{
			lastID:   1,
			sessions: make(map[string]*models.SessionApplicant),
		}, &SessionEmployerRepo{
			lastID:   1,
			sessions: make(map[string]*models.SessionEmployer),
		}
}

// Add adds new applicantsession.
// Accepts applicant's id and session's id.
// Returns error if session already exists or nil in case of success
func (sa *SessionApplicantRepo) Create(userId uint64, sessionID string) error {
	if _, ok := sa.sessions[sessionID]; ok {
		return ErrSessionAlreadyExists
	}
	sa.sessions[sessionID] = &models.SessionApplicant{
		ID:          sa.lastID,
		ApplicantID: userId,
		CookieToken: sessionID,
		CreatedAt:   time.Now().String(),
		UpdatedAt:   time.Now().String(),
	}
	sa.lastID++
	return nil
}

func (sa *SessionApplicantRepo) GetUserIdBySession(sessionId string) (uint64, error) {
	if session, ok := sa.sessions[sessionId]; ok {
		return session.ApplicantID, nil
	}
	return 0, ErrSessionNotFound
}

func (sa *SessionApplicantRepo) Delete(sessionId string) error {
	delete(sa.sessions, sessionId)
	return nil
}

func (se *SessionEmployerRepo) Create(userId uint64, sessionID string) error {
	if _, ok := se.sessions[sessionID]; ok {
		return ErrSessionAlreadyExists
	}
	se.sessions[sessionID] = &models.SessionEmployer{
		ID:          se.lastID,
		EmployerID:  userId,
		CookieToken: sessionID,
		CreatedAt:   time.Now().String(),
		UpdatedAt:   time.Now().String(),
	}
	se.lastID++
	return nil
}

func (se *SessionEmployerRepo) GetUserIdBySession(sessionId string) (uint64, error) {
	if session, ok := se.sessions[sessionId]; ok {
		return session.EmployerID, nil
	}
	return 0, ErrSessionNotFound
}

func (se *SessionEmployerRepo) Delete(sessionID string) error {
	delete(se.sessions, sessionID)
	return nil
}
