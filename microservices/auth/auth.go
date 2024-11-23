package auth

import (
	"context"
	"fmt"

	gen "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
)

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
	BadUserType = fmt.Errorf("Bad user type")
	NoUserExist = fmt.Errorf("There is no such user")
	NoSessionExist = fmt.Errorf("There is no such session")
)
