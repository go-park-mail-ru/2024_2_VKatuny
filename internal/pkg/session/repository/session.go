package repository

import (
	"fmt"
)

// TODO: use session/session.go instead of this
type SessionRepository interface {
	Create(uint64, string) error
	GetUserIdBySession(string) (uint64, error)
	Delete(string) error
}

var (
	ErrSessionAlreadyExists = fmt.Errorf("session already exists")
	ErrSessionNotFound      = fmt.Errorf("session not found")
)
