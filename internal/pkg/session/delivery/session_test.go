// Package delivery is a handlers layer of session
package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	auth_grpc "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
	grpc_mock "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/mock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestIsAuthorized(t *testing.T) {
	t.Parallel()
	type usecase struct {
		auth_grpc *grpc_mock.MockAuthorizationClient
	}
	tests := []struct {
		name         string
		r            *http.Request
		w            *httptest.ResponseRecorder
		usecase      *usecase
		codeExpected int

		prepare func(
			r *http.Request,
			w *httptest.ResponseRecorder,
			usecase *usecase,
		) (*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name:         "IsAuthorized: no cookie",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusUnauthorized,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				nr := httptest.NewRequest(
					http.MethodGet,
					"/api/v1/authorized",
					nil,
				)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "IsAuthorized: bad user type",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusUnauthorized,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				cookie := &http.Cookie{
					Name:  dto.SessionIDName,
					Value: "0123456789",
				}
				nr := httptest.NewRequest(
					http.MethodGet,
					"/api/v1/authorized",
					nil,
				)
				nr.AddCookie(cookie)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "IsAuthorized: no request id",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusUnauthorized,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				cookie := &http.Cookie{
					Name:  dto.SessionIDName,
					Value: "1234567890",
				}
				nr := httptest.NewRequest(
					http.MethodGet,
					"/api/v1/authorized",
					nil,
				)
				nr.AddCookie(cookie)
				request := &auth_grpc.CheckAuthRequest{
					Session: &auth_grpc.SessionToken{
						ID: cookie.Value,
					},
				}
				usecase.auth_grpc.
					EXPECT().
					CheckAuth(gomock.Any(), request).
					Return(nil, errors.New("bad grpc response"))
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "IsAuthorized: bad grpc",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusUnauthorized,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				cookie := &http.Cookie{
					Name:  dto.SessionIDName,
					Value: "1234567890",
				}
				nr := httptest.NewRequest(
					http.MethodGet,
					"/api/v1/authorized",
					nil,
				)
				nr.AddCookie(cookie)
				request := &auth_grpc.CheckAuthRequest{
					Session: &auth_grpc.SessionToken{
						ID: cookie.Value,
					},
				}
				usecase.auth_grpc.
					EXPECT().
					CheckAuth(gomock.Any(), request).
					Return(nil, errors.New("bad grpc response"))
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "IsAuthorized: ok",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusOK,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				cookie := &http.Cookie{
					Name:  dto.SessionIDName,
					Value: "1234567890",
				}
				nr := httptest.NewRequest(
					http.MethodGet,
					"/api/v1/authorized",
					nil,
				)
				nr.AddCookie(cookie)
				request := &auth_grpc.CheckAuthRequest{
					Session: &auth_grpc.SessionToken{
						ID: cookie.Value,
					},
				}
				response := &auth_grpc.CheckAuthResponse{
					UserData: &auth_grpc.User{
						UserType: dto.UserTypeApplicant,
						ID:       uint64(1),
					},
				}
				usecase.auth_grpc.
					EXPECT().
					CheckAuth(gomock.Any(), request).
					Return(response, nil)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := &usecase{
				auth_grpc: grpc_mock.NewMockAuthorizationClient(ctrl),
			}
			tt.w, tt.r = tt.prepare(tt.r, tt.w, usecase)

			app := &internal.App{
				Logger:         logrus.New(),
				BackendAddress: "http://localhost:8080",
				Microservices: &internal.Microservices{
					Auth: usecase.auth_grpc,
				},
			}

			h := NewSessionHandlers(app)
			require.NotNil(t, h)
			require.NotNil(t, tt.r)
			require.NotNil(t, tt.w)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/authorized", h.IsAuthorized).Methods(http.MethodGet)

			r.ServeHTTP(tt.w, tt.r)

			require.Equal(t, tt.codeExpected, tt.w.Code)
		})
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()
	type usecase struct {
		auth_grpc *grpc_mock.MockAuthorizationClient
	}
	tests := []struct {
		name         string
		r            *http.Request
		w            *httptest.ResponseRecorder
		usecase      *usecase
		codeExpected int

		prepare func(
			r *http.Request,
			w *httptest.ResponseRecorder,
			usecase *usecase,
		) (*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name:         "Login: no json",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusBadRequest,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				nr := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/login",
					nil,
				)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "Login: bad grpc",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusUnauthorized,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				form := &dto.JSONLoginForm{
					UserType: dto.UserTypeApplicant,
					Email:    "testes@test.ru",
					Password: "pass1234",
				}
				JSONform, _ := json.Marshal(form)
				nr := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/login",
					bytes.NewReader(JSONform),
				)
				request := &auth_grpc.AuthRequest{
					UserType: form.UserType,
					Email:    form.Email,
					Password: form.Password,
				}
				usecase.auth_grpc.
					EXPECT().
					AuthUser(gomock.Any(), request).
					Return(nil, errors.New("bad grpc"))
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "Login: ok",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusOK,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				form := &dto.JSONLoginForm{
					UserType: dto.UserTypeApplicant,
					Email:    "testes@test.ru",
					Password: "pass1234",
				}
				JSONform, _ := json.Marshal(form)
				nr := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/login",
					bytes.NewReader(JSONform),
				)
				request := &auth_grpc.AuthRequest{
					UserType: form.UserType,
					Email:    form.Email,
					Password: form.Password,
				}
				response := &auth_grpc.AuthResponse{
					Session: &auth_grpc.SessionToken{
						ID: "1234567890",
					},
					UserData: &auth_grpc.User{
						UserType: form.UserType,
						ID:       uint64(1),
					},
				}
				usecase.auth_grpc.
					EXPECT().
					AuthUser(gomock.Any(), request).
					Return(response, nil)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := &usecase{
				auth_grpc: grpc_mock.NewMockAuthorizationClient(ctrl),
			}
			tt.w, tt.r = tt.prepare(tt.r, tt.w, usecase)

			app := &internal.App{
				Logger:         logrus.New(),
				BackendAddress: "http://localhost:8080",
				Microservices: &internal.Microservices{
					Auth: usecase.auth_grpc,
				},
			}

			h := NewSessionHandlers(app)
			require.NotNil(t, h)
			require.NotNil(t, tt.r)
			require.NotNil(t, tt.w)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/login", h.Login).Methods(http.MethodPost)

			r.ServeHTTP(tt.w, tt.r)

			require.Equal(t, tt.codeExpected, tt.w.Code)
		})
	}
}

func TestLogout(t *testing.T) {
	t.Parallel()
	type usecase struct {
		auth_grpc *grpc_mock.MockAuthorizationClient
	}
	tests := []struct {
		name         string
		r            *http.Request
		w            *httptest.ResponseRecorder
		usecase      *usecase
		codeExpected int

		prepare func(
			r *http.Request,
			w *httptest.ResponseRecorder,
			usecase *usecase,
		) (*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name:         "Logout: no cookie",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusOK,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				nr := httptest.NewRequest(
					http.MethodGet,
					"/api/v1/logout",
					nil,
				)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "Logout: bad user type",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusOK,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				cookie := &http.Cookie{
					Name:  dto.SessionIDName,
					Value: "0123456789",
				}
				nr := httptest.NewRequest(
					http.MethodGet,
					"/api/v1/logout",
					nil,
				)
				nr.AddCookie(cookie)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "Logout: no request id",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusUnauthorized,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				cookie := &http.Cookie{
					Name:  dto.SessionIDName,
					Value: "1234567890",
				}
				nr := httptest.NewRequest(
					http.MethodGet,
					"/api/v1/logout",
					nil,
				)
				nr.AddCookie(cookie)
				request := &auth_grpc.DeauthRequest{
					Session: &auth_grpc.SessionToken{
						ID: "1234567890",
					},
				}
				usecase.auth_grpc.
					EXPECT().
					DeauthUser(gomock.Any(), request).
					Return(nil, errors.New("bad grpc response"))
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "Logout: ok",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusOK,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				cookie := &http.Cookie{
					Name:  dto.SessionIDName,
					Value: "1234567890",
				}
				nr := httptest.NewRequest(
					http.MethodGet,
					"/api/v1/logout",
					nil,
				)
				nr.AddCookie(cookie)
				request := &auth_grpc.DeauthRequest{
					Session: &auth_grpc.SessionToken{
						ID: "1234567890",
					},
				}
				response := new(auth_grpc.DeauthResponse)
				usecase.auth_grpc.
					EXPECT().
					DeauthUser(gomock.Any(), request).
					Return(response, nil)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := &usecase{
				auth_grpc: grpc_mock.NewMockAuthorizationClient(ctrl),
			}
			tt.w, tt.r = tt.prepare(tt.r, tt.w, usecase)

			app := &internal.App{
				Logger:         logrus.New(),
				BackendAddress: "http://localhost:8080",
				Microservices: &internal.Microservices{
					Auth: usecase.auth_grpc,
				},
			}

			h := NewSessionHandlers(app)
			require.NotNil(t, h)
			require.NotNil(t, tt.r)
			require.NotNil(t, tt.w)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/logout", h.Logout).Methods(http.MethodGet)

			r.ServeHTTP(tt.w, tt.r)

			require.Equal(t, tt.codeExpected, tt.w.Code)
		})
	}
}
