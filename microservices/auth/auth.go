package auth

import (
	"context"
	"fmt"

	gen "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
)

var sessionTTL = 86400 // seconds = 24 hours

type IAuthorizationDelivery interface {
	Auth(ctx context.Context, req *gen.AuthRequest) (*gen.AuthResponse, error)
	CheckAuth(ctx context.Context, req *gen.CheckAuthRequest) (*gen.CheckAuthResponse, error)
	Deauth(ctx context.Context, req *gen.DeauthRequest) (*gen.DeauthResponse, error)
}

type IAuthorizationRepository interface {
	GetUser(userType, email string) (*User, error)
	CreateSession(uint64, string) error
	GetUserIdBySession(string) (uint64, error)
	DeleteSession(string) error
}

type User struct {
	UserType     string
	ID           uint64
	Email        string
	PasswordHash string
}

var (
	ErrBadUserType = fmt.Errorf("bad user type")
	ErrNoUserExist = fmt.Errorf("there is no such user")
	ErrNoSessionExist = fmt.Errorf("there is no such session")
)
