package auth_test

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	gen "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

func TestNewAuthorization(t *testing.T) {
	type args struct {
		dbConn    *sql.DB
		redisConn redis.Conn
		logger    *logrus.Logger
	}
	tests := []struct {
		name string
		args args
		want *AuthorizationDelivery
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthorization(tt.args.dbConn, tt.args.redisConn, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthorization() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorizationDelivery_AuthUser(t *testing.T) {
	type fields struct {
		UnimplementedAuthorizationServer gen.UnimplementedAuthorizationServer
		authRepo                         IAuthorizationRepository
		logger                           *logrus.Entry
		tokenLength                      uint
		sessionTTL                       time.Duration
	}
	type args struct {
		ctx context.Context
		req *gen.AuthRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *gen.AuthResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationDelivery{
				UnimplementedAuthorizationServer: tt.fields.UnimplementedAuthorizationServer,
				authRepo:                         tt.fields.authRepo,
				logger:                           tt.fields.logger,
				tokenLength:                      tt.fields.tokenLength,
				sessionTTL:                       tt.fields.sessionTTL,
			}
			got, err := a.AuthUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthorizationDelivery.AuthUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthorizationDelivery.AuthUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorizationDelivery_CheckAuth(t *testing.T) {
	type fields struct {
		UnimplementedAuthorizationServer gen.UnimplementedAuthorizationServer
		authRepo                         IAuthorizationRepository
		logger                           *logrus.Entry
		tokenLength                      uint
		sessionTTL                       time.Duration
	}
	type args struct {
		ctx context.Context
		req *gen.CheckAuthRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *gen.CheckAuthResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationDelivery{
				UnimplementedAuthorizationServer: tt.fields.UnimplementedAuthorizationServer,
				authRepo:                         tt.fields.authRepo,
				logger:                           tt.fields.logger,
				tokenLength:                      tt.fields.tokenLength,
				sessionTTL:                       tt.fields.sessionTTL,
			}
			got, err := a.CheckAuth(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthorizationDelivery.CheckAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthorizationDelivery.CheckAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorizationDelivery_DeauthUser(t *testing.T) {
	type fields struct {
		UnimplementedAuthorizationServer gen.UnimplementedAuthorizationServer
		authRepo                         IAuthorizationRepository
		logger                           *logrus.Entry
		tokenLength                      uint
		sessionTTL                       time.Duration
	}
	type args struct {
		ctx context.Context
		req *gen.DeauthRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *gen.DeauthResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationDelivery{
				UnimplementedAuthorizationServer: tt.fields.UnimplementedAuthorizationServer,
				authRepo:                         tt.fields.authRepo,
				logger:                           tt.fields.logger,
				tokenLength:                      tt.fields.tokenLength,
				sessionTTL:                       tt.fields.sessionTTL,
			}
			got, err := a.DeauthUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthorizationDelivery.DeauthUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthorizationDelivery.DeauthUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateToken(t *testing.T) {
	type args struct {
		tokenLength uint
		userType    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateToken(tt.args.tokenLength, tt.args.userType)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
