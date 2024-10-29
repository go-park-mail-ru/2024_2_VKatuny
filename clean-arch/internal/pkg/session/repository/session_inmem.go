package repository

import (
	"time"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/session"
)

type sessionApplicantRepo struct {
	lastID   uint64
	sessions map[string]*models.SessionApplicant
}

type sessionEmployerRepo struct {
	lastID   uint64
	sessions map[string]*models.SessionEmployer
}

func NewRepo() (*sessionApplicantRepo, *sessionEmployerRepo) {
	return &sessionApplicantRepo{
		sessions: make(map[string]*models.SessionApplicant),
	}, &sessionEmployerRepo{
		sessions: make(map[string]*models.SessionEmployer),
	}
}

func (sa *sessionApplicantRepo) Add(userId uint64, sessionId string) error {
	if _, ok := sa.sessions[sessionId]; ok {
		return session.ErrSessionAlreadyExists
	}
	sa.sessions[sessionId] = &models.SessionApplicant{
		ID:          sa.lastID,
		ApplicantID: userId,
		CookieToken: sessionId,
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

func (se *sessionEmployerRepo) Add(userId uint64, sessionId string) error {
	if _, ok := se.sessions[sessionId]; ok {
		return session.ErrSessionAlreadyExists
	}
	se.sessions[sessionId] = &models.SessionEmployer{
		ID: se.lastID,
		EmployerID: userId,
		CookieToken: sessionId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

func (se *sessionEmployerRepo) Delete(sessionId string) error {
	delete(se.sessions, sessionId)
	return nil
}
