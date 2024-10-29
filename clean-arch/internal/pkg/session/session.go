package session

import (
	"fmt"
)

type Repository interface {
	Add(uint64, string) error
	GetUserIdBySession(string) (uint64, error)
	Delete(string) error
}

var (
	ErrSessionAlreadyExists = fmt.Errorf("session already exists")
	ErrSessionNotFound      = fmt.Errorf("session not found")
)
