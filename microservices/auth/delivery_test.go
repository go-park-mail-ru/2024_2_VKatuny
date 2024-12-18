package auth_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth"
	gen "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/mock"
	"github.com/rafaeljusto/redigomock/v3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthUser(t *testing.T) {
	type repo struct {
		auth *mock.MockIAuthorizationRepository
	}
	tests := []struct {
		name    string
		ctx     context.Context
		req     *gen.AuthRequest
		res     *gen.AuthResponse
		err     error
		prepare func(
			ctx context.Context,
			req *gen.AuthRequest,
			res *gen.AuthResponse,
			repo *repo,
		) (context.Context, *gen.AuthRequest, *gen.AuthResponse, *repo)
	}{
		{
			name: "nil request",
			ctx:  context.Background(),
			req:  nil,
			res: &gen.AuthResponse{
				Status: gen.StatusCode_EmptyRequest,
			},
			err: status.Error(codes.InvalidArgument, "request is nil"),
			prepare: func(
				ctx context.Context,
				req *gen.AuthRequest,
				res *gen.AuthResponse,
				repo *repo,
			) (context.Context, *gen.AuthRequest, *gen.AuthResponse, *repo) {
				return ctx, req, res, repo
			},
		},
		{
			name: "get user error",
			ctx:  context.Background(),
			req:  new(gen.AuthRequest),
			res: &gen.AuthResponse{
				Status: gen.StatusCode_UnableToGetUser,
			},
			err: status.Error(codes.Internal, "failed to get user"),
			prepare: func(
				ctx context.Context,
				req *gen.AuthRequest,
				res *gen.AuthResponse,
				repo *repo,
			) (context.Context, *gen.AuthRequest, *gen.AuthResponse, *repo) {
				req.RequestID = "123"
				req.UserType = "student"
				req.Email = "vasya@mail.ru"
				repo.auth.
					EXPECT().
					GetUser(req.UserType, req.Email).
					Return(new(auth.User), auth.ErrBadUserType)
				return ctx, req, res, repo
			},
		},
		{
			name: "userType mismatch",
			ctx:  context.Background(),
			req:  new(gen.AuthRequest),
			res: &gen.AuthResponse{
				Status: gen.StatusCode_InvalidCredentials,
			},
			err: status.Error(codes.InvalidArgument, "wrong user type"),
			prepare: func(
				ctx context.Context,
				req *gen.AuthRequest,
				res *gen.AuthResponse,
				repo *repo,
			) (context.Context, *gen.AuthRequest, *gen.AuthResponse, *repo) {
				req.RequestID = "123"
				req.UserType = "employer"
				req.Email = "vasya@mail.ru"
				user := &auth.User{
					UserType: "applicant",
				}
				repo.auth.
					EXPECT().
					GetUser(req.UserType, req.Email).
					Return(user, nil)
				return ctx, req, res, repo
			},
		},
		{
			name: "email mismatch",
			ctx:  context.Background(),
			req:  new(gen.AuthRequest),
			res: &gen.AuthResponse{
				Status: gen.StatusCode_InvalidCredentials,
			},
			err: status.Error(codes.InvalidArgument, "wrong login or password"),
			prepare: func(
				ctx context.Context,
				req *gen.AuthRequest,
				res *gen.AuthResponse,
				repo *repo,
			) (context.Context, *gen.AuthRequest, *gen.AuthResponse, *repo) {
				req.RequestID = "123"
				req.UserType = "employer"
				req.Email = "vasya@mail.ru"
				user := &auth.User{
					UserType: "employer",
					Email: "ne_vasya@mail.ru",
				}
				repo.auth.
					EXPECT().
					GetUser(req.UserType, req.Email).
					Return(user, nil)
				return ctx, req, res, repo
			},
		},
		{
			name: "password mismatch",
			ctx:  context.Background(),
			req:  new(gen.AuthRequest),
			res: &gen.AuthResponse{
				Status: gen.StatusCode_InvalidCredentials,
			},
			err: status.Error(codes.InvalidArgument, "wrong login or password"),
			prepare: func(
				ctx context.Context,
				req *gen.AuthRequest,
				res *gen.AuthResponse,
				repo *repo,
			) (context.Context, *gen.AuthRequest, *gen.AuthResponse, *repo) {
				req.RequestID = "123"
				req.UserType = "employer"
				req.Email = "vasya@mail.ru"
				req.Password = "pass1234"
				user := &auth.User{
					UserType: "employer",
					Email: "vasya@mail.ru",
					PasswordHash: "hash_pass1234",
				}
				repo.auth.
					EXPECT().
					GetUser(req.UserType, req.Email).
					Return(user, nil)
				return ctx, req, res, repo
			},
		},
		{
			name: "session creation failed",
			ctx:  context.Background(),
			req:  new(gen.AuthRequest),
			res: &gen.AuthResponse{
				Status: gen.StatusCode_UnableToCreateSession,
			},
			err: status.Error(codes.Internal, "failed to create session"),
			prepare: func(
				ctx context.Context,
				req *gen.AuthRequest,
				res *gen.AuthResponse,
				repo *repo,
			) (context.Context, *gen.AuthRequest, *gen.AuthResponse, *repo) {
				req.RequestID = "123"
				req.UserType = "employer"
				req.Email = "vasya@mail.ru"
				req.Password = "pass1234"
				user := &auth.User{
					UserType: "employer",
					Email: "vasya@mail.ru",
					PasswordHash: "$2a$10$9VjOcWJ7.SF6uLKiVurMZ.2R5pllY/ayp6jMKvgl6CPd6El5PMqqC",
				}
				repo.auth.
					EXPECT().
					GetUser(req.UserType, req.Email).
					Return(user, nil)
				repo.auth.
					EXPECT().
					CreateSession(user.ID, gomock.Any()).
					Return(fmt.Errorf("session creation failed"))
				return ctx, req, res, repo
			},
		},
		{
			name: "OK",
			ctx:  context.Background(),
			req:  new(gen.AuthRequest),
			res: &gen.AuthResponse{
				Status: gen.StatusCode_OK,
				UserData: &gen.User{
					UserType: "employer",
					ID:       uint64(1),
				},
			},
			err: nil,
			prepare: func(
				ctx context.Context,
				req *gen.AuthRequest,
				res *gen.AuthResponse,
				repo *repo,
			) (context.Context, *gen.AuthRequest, *gen.AuthResponse, *repo) {
				req.RequestID = "123"
				req.UserType = "employer"
				req.Email = "vasya@mail.ru"
				req.Password = "pass1234"
				user := &auth.User{
					UserType: "employer",
					ID:       uint64(1),
					Email: "vasya@mail.ru",
					PasswordHash: "$2a$10$9VjOcWJ7.SF6uLKiVurMZ.2R5pllY/ayp6jMKvgl6CPd6El5PMqqC",
				}
				repo.auth.
					EXPECT().
					GetUser(req.UserType, req.Email).
					Return(user, nil)
				repo.auth.
					EXPECT().
					CreateSession(user.ID, gomock.Any()).
					Return(nil)
				return ctx, req, res, repo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := &repo{
				auth: mock.NewMockIAuthorizationRepository(ctrl),
			}
			tt.ctx, tt.req, tt.res, repo = tt.prepare(tt.ctx, tt.req, tt.res, repo)

			d := auth.NewAuthorization(new(sql.DB), redigomock.NewConn(), new(logrus.Logger))
			d.AuthRepo = repo.auth

			response, err := d.AuthUser(tt.ctx, tt.req)
			require.Equal(t, tt.err, err)
			response.Session = nil
			require.Equal(t, tt.res, response)
		})
	}
}

func TestCheckAuth(t *testing.T) {
	type repo struct {
		auth *mock.MockIAuthorizationRepository
	}
	tests := []struct {
		name    string
		ctx     context.Context
		req     *gen.CheckAuthRequest
		res     *gen.CheckAuthResponse
		err     error
		prepare func(
			ctx context.Context,
			req *gen.CheckAuthRequest,
			res *gen.CheckAuthResponse,
			repo *repo,
		) (context.Context, *gen.CheckAuthRequest, *gen.CheckAuthResponse, *repo)
	}{
		{
			name: "nil request",
			ctx:  context.Background(),
			req:  nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "request is nil"),
			prepare: func(
				ctx context.Context,
				req *gen.CheckAuthRequest,
				res *gen.CheckAuthResponse,
				repo *repo,
			) (context.Context, *gen.CheckAuthRequest, *gen.CheckAuthResponse, *repo) {
				return ctx, req, res, repo
			},
		},
		{
			name: "can't get user id by session",
			ctx:  context.Background(),
			req:  new(gen.CheckAuthRequest),
			res: &gen.CheckAuthResponse{
				Status: gen.StatusCode_NoSessionExist,
			},
			err: status.Error(codes.Internal, "failed to get user"),
			prepare: func(
				ctx context.Context,
				req *gen.CheckAuthRequest,
				res *gen.CheckAuthResponse,
				repo *repo,
			) (context.Context, *gen.CheckAuthRequest, *gen.CheckAuthResponse, *repo) {
				req.RequestID = "123"
				req.Session = &gen.SessionToken{
					ID: "1token",
				}

				repo.auth.
					EXPECT().
					GetUserIdBySession(req.Session.ID).
					Return(uint64(0), fmt.Errorf("can't get user id by session"))
				return ctx, req, res, repo
			},
		},
		{
			name: "OK",
			ctx:  context.Background(),
			req:  new(gen.CheckAuthRequest),
			res: &gen.CheckAuthResponse{
				Status: gen.StatusCode_OK,
				UserData: &gen.User{
					UserType: "applicant",
					ID:       uint64(1),
				},
			},
			err: nil,
			prepare: func(
				ctx context.Context,
				req *gen.CheckAuthRequest,
				res *gen.CheckAuthResponse,
				repo *repo,
			) (context.Context, *gen.CheckAuthRequest, *gen.CheckAuthResponse, *repo) {
				req.RequestID = "123"
				req.Session = &gen.SessionToken{
					ID: "1token",
				}

				repo.auth.
					EXPECT().
					GetUserIdBySession(req.Session.ID).
					Return(uint64(1), nil)
				return ctx, req, res, repo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := &repo{
				auth: mock.NewMockIAuthorizationRepository(ctrl),
			}
			tt.ctx, tt.req, tt.res, repo = tt.prepare(tt.ctx, tt.req, tt.res, repo)

			d := auth.NewAuthorization(new(sql.DB), redigomock.NewConn(), new(logrus.Logger))
			d.AuthRepo = repo.auth

			response, err := d.CheckAuth(tt.ctx, tt.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.res, response)
		})
	}
}

func TestDeauthUser(t *testing.T) {
	type repo struct {
		auth *mock.MockIAuthorizationRepository
	}
	tests := []struct {
		name    string
		ctx     context.Context
		req     *gen.DeauthRequest
		res     *gen.DeauthResponse
		err     error
		prepare func(
			ctx context.Context,
			req *gen.DeauthRequest,
			res *gen.DeauthResponse,
			repo *repo,
		) (context.Context, *gen.DeauthRequest, *gen.DeauthResponse, *repo)
	}{
		{
			name: "nil request",
			ctx:  context.Background(),
			req:  nil,
			res: &gen.DeauthResponse{
				Status: gen.StatusCode_EmptyRequest,
			},
			err: status.Error(codes.InvalidArgument, "request is nil"),
			prepare: func(
				ctx context.Context,
				req *gen.DeauthRequest,
				res *gen.DeauthResponse,
				repo *repo,
			) (context.Context, *gen.DeauthRequest, *gen.DeauthResponse, *repo) {
				return ctx, req, res, repo
			},
		},
		{
			name: "nil session",
			ctx:  context.Background(),
			req:  &gen.DeauthRequest{
				Session: nil,
			},
			res: &gen.DeauthResponse{
				Status: gen.StatusCode_NoSessionExist,
			},
			err: status.Error(codes.InvalidArgument, "session is nil"),
			prepare: func(
				ctx context.Context,
				req *gen.DeauthRequest,
				res *gen.DeauthResponse,
				repo *repo,
			) (context.Context, *gen.DeauthRequest, *gen.DeauthResponse, *repo) {
				return ctx, req, res, repo
			},
		},
		{
			name: "failed to delete session",
			ctx:  context.Background(),
			req: new(gen.DeauthRequest),
			res: &gen.DeauthResponse{
				Status: gen.StatusCode_UnableToDeleteSession,
			},
			err: status.Error(codes.Internal, "failed to delete session"),
			prepare: func(
				ctx context.Context,
				req *gen.DeauthRequest,
				res *gen.DeauthResponse,
				repo *repo,
			) (context.Context, *gen.DeauthRequest, *gen.DeauthResponse, *repo) {
				req.Session = &gen.SessionToken{
					ID: "1token",
				}
				repo.auth.
					EXPECT().
					DeleteSession(req.Session.ID).
					Return(fmt.Errorf("failed to delete session"))
				return ctx, req, res, repo
			},
		},
		{
			name: "OK",
			ctx:  context.Background(),
			req: new(gen.DeauthRequest),
			res: &gen.DeauthResponse{
				Status: gen.StatusCode_OK,
			},
			err: nil,
			prepare: func(
				ctx context.Context,
				req *gen.DeauthRequest,
				res *gen.DeauthResponse,
				repo *repo,
			) (context.Context, *gen.DeauthRequest, *gen.DeauthResponse, *repo) {
				req.Session = &gen.SessionToken{
					ID: "1token",
				}
				repo.auth.
					EXPECT().
					DeleteSession(req.Session.ID).
					Return(nil)
				return ctx, req, res, repo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := &repo{
				auth: mock.NewMockIAuthorizationRepository(ctrl),
			}
			tt.ctx, tt.req, tt.res, repo = tt.prepare(tt.ctx, tt.req, tt.res, repo)

			d := auth.NewAuthorization(new(sql.DB), redigomock.NewConn(), new(logrus.Logger))
			d.AuthRepo = repo.auth

			response, err := d.DeauthUser(tt.ctx, tt.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.res, response)
		})
	}
}   

