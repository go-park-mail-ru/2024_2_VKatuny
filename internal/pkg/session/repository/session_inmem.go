// Package repository is a repository layer of session
package repository

import (
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/session"
)

// sessionApplicantRepo is a in-memory implementation for applicant session
type sessionApplicantRepo struct {
	lastID   uint64
	sessions map[string]*models.SessionApplicant
}

// sessionEmployerRepo is a in-memory implementation for employer session
type sessionEmployerRepo struct {
	lastID   uint64
	sessions map[string]*models.SessionEmployer
}

// NewRepo returns new in-memory repository for sessions
func NewRepo() (*sessionApplicantRepo, *sessionEmployerRepo) {
	return &sessionApplicantRepo{
			lastID:   1,
			sessions: make(map[string]*models.SessionApplicant),
		}, &sessionEmployerRepo{
			lastID:   1,
			sessions: make(map[string]*models.SessionEmployer),
		}
}

// Add adds new applicantsession.
// Accepts applicant's id and session's id.
// Returns error if session already exists or nil in case of success
func (sa *sessionApplicantRepo) Add(userId uint64, sessionID string) error {
	if _, ok := sa.sessions[sessionID]; ok {
		return session.ErrSessionAlreadyExists
	}
	sa.sessions[sessionID] = &models.SessionApplicant{
		ID:          sa.lastID,
		ApplicantID: userId,
		CookieToken: sessionID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	sa.lastID++
	return nil
}

func (sa *sessionApplicantRepo) GetUserIdBySession(sessionId string) (uint64, error) {
	if session, ok := sa.sessions[sessionId]; ok {
		return session.ApplicantID, nil
	}
	return 0, session.ErrSessionNotFound
}

func (sa *sessionApplicantRepo) Delete(sessionId string) error {
	delete(sa.sessions, sessionId)
	return nil
}

func (se *sessionEmployerRepo) Add(userId uint64, sessionID string) error {
	if _, ok := se.sessions[sessionID]; ok {
		return session.ErrSessionAlreadyExists
	}
	se.sessions[sessionID] = &models.SessionEmployer{
		ID:          se.lastID,
		EmployerID:  userId,
		CookieToken: sessionID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	se.lastID++
	return nil
}

func (se *sessionEmployerRepo) GetUserIdBySession(sessionId string) (uint64, error) {
	if session, ok := se.sessions[sessionId]; ok {
		return session.EmployerID, nil
	}
	return 0, session.ErrSessionNotFound
}

func (se *sessionEmployerRepo) Delete(sessionID string) error {
	delete(se.sessions, sessionID)
	return nil
}
