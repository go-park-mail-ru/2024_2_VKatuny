package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"math/big"
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	gen "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthorizationDelivery struct {
	gen.UnimplementedAuthorizationServer

	authRepo IAuthorizationRepository
	logger   *logrus.Entry

	tokenLength uint
	sessionTTL  time.Duration
}

func NewAuthorization(dbConn *sql.DB, redisConn redis.Conn, logger *logrus.Logger) *AuthorizationDelivery {
	fn := "NewAuthorizationDelivery"
	if dbConn == nil {
		logger.Fatal("db connection is nil")
	}
	if redisConn == nil {
		logger.Fatal("redis connection is nil")
	}
	logger.Infof("%s: initializing", fn)
	return &AuthorizationDelivery{
		authRepo:    NewAuthorizationRepository(logger, dbConn, redisConn),
		logger:      &logrus.Entry{Logger: logger},
		tokenLength: 32,
		sessionTTL:  24 * time.Hour,
	}
}

func (a *AuthorizationDelivery) AuthUser(ctx context.Context, req *gen.AuthRequest) (*gen.AuthResponse, error) {
	fn := "AuthDelivery.Auth"

	a.logger.Debugf("%s: got request: %s", fn, req)
	response := new(gen.AuthResponse)

	if req == nil {
		a.logger.Errorf("%s: request is nil", fn)
		response.Status = gen.StatusCode_EmptyRequest
		return response, status.Error(codes.InvalidArgument, "request is nil")
	}

	requestID := req.RequestID
	a.logger.Debugf("%s: got request with ID - %s", fn, requestID)

	token, err := generateToken(a.tokenLength, req.UserType)
	if err != nil {
		a.logger.Errorf("%s: failed to generate token: %s", fn, err)
		response.Status = gen.StatusCode_UnableToGenerateToken
		return response, status.Error(codes.Internal, "failed to generate token")
	}
	a.logger.Debugf("%s: generated token: %s", fn, token)

	userInfo, err := a.authRepo.GetUser(req.UserType, req.Email)
	if err != nil {
		a.logger.Errorf("%s: got err %s", fn, err)
		response.Status = gen.StatusCode_UnableToGetUser
		return response, status.Error(codes.Internal, "failed to get user")
	}

	if req.UserType != userInfo.UserType {
		a.logger.Errorf("%s: user type mismatch, expected %s, got %s", fn, userInfo.UserType, req.UserType)
		response.Status = gen.StatusCode_InvalidCredentials
		return response, status.Error(codes.InvalidArgument, "wrong user type")
	}
	if req.Email != userInfo.Email {
		a.logger.Errorf("%s: email mismatch, expected %s, got %s", fn, userInfo.Email, req.Email)
		response.Status = gen.StatusCode_InvalidCredentials
		return response, status.Error(codes.InvalidArgument, "wrong login or password")
	}

	if !utils.EqualHashedPasswords(userInfo.PasswordHash, req.Password) {
		a.logger.Errorf("%s: password comparison failed", fn)
		response.Status = gen.StatusCode_InvalidCredentials
		return response, status.Error(codes.InvalidArgument, "wrong login or password")
	}

	err = a.authRepo.CreateSession(userInfo.ID, token)
	if err != nil {
		a.logger.Errorf("%s: got err %s", fn, err)
		response.Status = gen.StatusCode_UnableToCreateSession
		return response, status.Error(codes.Internal, "failed to create session")
	}

	response.Status = gen.StatusCode_OK
	response.Session = &gen.SessionToken{
		ID:             token,
		ExpirationDate: timestamppb.New(time.Now().Add(a.sessionTTL)),
	}
	response.UserData = &gen.User{
		UserType: userInfo.UserType,
		ID:       userInfo.ID,
	}
	return response, nil
}

func (a *AuthorizationDelivery) CheckAuth(ctx context.Context, req *gen.CheckAuthRequest) (*gen.CheckAuthResponse, error) {
	fn := "AuthDelivery.CheckAuth"
	if req == nil {
		a.logger.Errorf("%s: request is nil", fn)
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	requestID := req.RequestID
	a.logger.Debugf("%s: got request with ID - %s", fn, requestID)

	response := new(gen.CheckAuthResponse)

	sessionID := req.Session.ID

	userID, err := a.authRepo.GetUserIdBySession(sessionID)
	if err != nil {
		a.logger.Errorf("%s: got err %s", fn, err)
		response.Status = gen.StatusCode_NoSessionExist
		return response, status.Error(codes.Internal, "failed to get user")
	}
	userType, _ := utils.CheckToken(sessionID)

	// TODO: DB should contain expiration date
	response.Status = gen.StatusCode_OK
	response.UserData = &gen.User{
		UserType: userType,
		ID:       userID,
	}
	return response, nil
}

func (a *AuthorizationDelivery) DeauthUser(ctx context.Context, req *gen.DeauthRequest) (*gen.DeauthResponse, error) {
	fn := "AuthDelivery.Deauth"

	response := new(gen.DeauthResponse)

	if req == nil {
		a.logger.Errorf("%s: request is nil", fn)
		response.Status = gen.StatusCode_EmptyRequest
		return response, status.Error(codes.InvalidArgument, "request is nil")
	}

	if req.Session == nil {
		a.logger.Errorf("%s: session is nil", fn)
		response.Status = gen.StatusCode_NoSessionExist
		return response, status.Error(codes.InvalidArgument, "session is nil")
	}

	err := a.authRepo.DeleteSession(req.Session.ID)
	if err != nil {
		a.logger.Errorf("%s: got err %s", fn, err)
		response.Status = gen.StatusCode_UnableToDeleteSession
		return response, status.Error(codes.Internal, "failed to delete session")
	}

	response.Status = gen.StatusCode_OK
	return response, nil
}

func generateToken(tokenLength uint, userType string) (string, error) {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	buf := make([]rune, tokenLength)
	for i := range buf {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			return "", err
		}
		buf[i] = letterRunes[j.Int64()]
	}
	if userType == dto.UserTypeApplicant {
		return "1" + string(buf), nil
	}
	return "2" + string(buf), nil
}
