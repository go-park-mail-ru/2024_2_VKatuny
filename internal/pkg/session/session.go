package session

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

type ISessionRepository interface {
	Create(uint64, string) error
	GetUserIdBySession(string) (uint64, error)
	Delete(string) error
}

type ISessionUsecase interface {
	// GetUserTypeFromToken(sessionID string) (string, error)
	CheckAuthorization(ctx context.Context, userType string, sessionID string) (uint64, error)
	Login(ctx context.Context, form *dto.JSONLoginForm) (*dto.UserWithSession, error)
	Logout(ctx context.Context, userType string, sessionID string) (*dto.JSONUser, error)
}
