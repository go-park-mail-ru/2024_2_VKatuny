package delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/delivery"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/mock"
	file_loading_mock "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading/mock"
	vacancies_mock "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/mock"
	auth_grpc "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
	grpc_mock "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/mock"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createMultipartFormJSON(jsonForm *dto.JSONUpdateEmployerProfile) (*bytes.Buffer, string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	defer writer.Close()
	writer.WriteField("firstName", jsonForm.FirstName)
	writer.WriteField("lastName", jsonForm.LastName)
	writer.WriteField("city", jsonForm.City)
	writer.WriteField("contacts", jsonForm.Contacts)
	return &buf, writer.FormDataContentType()
}

func TestGetProfileHandler(t *testing.T) {
	t.Parallel()
	type usecase struct {
		profile            *mock.MockIEmployerUsecase
		fileLoadingUsecase *file_loading_mock.MockIFileLoadingUsecase
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
			name:         "TestGetProfile: bad slug",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusInternalServerError,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				slug := "123534657548574856785785346346367542151341354756869568"
				nr := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/employer/%s/profile", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "TestProfile: bad usecase",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusInternalServerError,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				slug := uint64(1)
				nr := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/employer/%d/profile", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				usecase.profile.
					EXPECT().
					GetEmployerProfile(gomock.Any(), slug).
					Return(nil, fmt.Errorf("error"))
				return nw, nr
			},
		},
		{
			name:         "TestProfile: bad usecase",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusOK,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				slug := uint64(1)
				nr := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/employer/%d/profile", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				employerProfile := &dto.JSONGetEmployerProfile{
					ID:        slug,
					FirstName: "John",
					LastName:  "Doe",
				}
				usecase.profile.
					EXPECT().
					GetEmployerProfile(gomock.Any(), slug).
					Return(employerProfile, nil)
				usecase.fileLoadingUsecase.
					EXPECT().
					FindCompressedFile(employerProfile.Avatar).
					Return("")
				return nw, nr
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := &usecase{
				profile: mock.NewMockIEmployerUsecase(ctrl),
				fileLoadingUsecase: file_loading_mock.NewMockIFileLoadingUsecase(ctrl),
			}
			tt.w, tt.r = tt.prepare(tt.r, tt.w, usecase)

			app := &internal.App{
				Logger:         logrus.New(),
				BackendAddress: "http://localhost:8080",
				Usecases: &internal.Usecases{
					EmployerUsecase:    usecase.profile,
					VacanciesUsecase:   nil,
					FileLoadingUsecase: usecase.fileLoadingUsecase,
				},
				Microservices: &internal.Microservices{
					Auth: nil,
				},
			}

			h := delivery.NewEmployerHandlers(app)
			require.NotNil(t, h)
			require.NotNil(t, tt.r)
			require.NotNil(t, tt.w)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/employer/{id:[0-9]+}/profile", h.GetProfile).Methods(http.MethodGet)

			r.ServeHTTP(tt.w, tt.r)

			require.Equal(t, tt.codeExpected, tt.w.Code)
		})

	}
}

func TestUpdateProfile(t *testing.T) {
	t.Parallel()
	type usecase struct {
		profile *mock.MockIEmployerUsecase
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
			name:         "UpdateProfile: bad slug",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusInternalServerError,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				slug := "123534657548574856785785346346367542151341354756869568"
				nr := httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/employer/%s/profile", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "GetProfile: bad usecase",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusInternalServerError,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				slug := uint64(1)
				newProfile := &dto.JSONUpdateEmployerProfile{
					FirstName: "John",
					LastName:  "Doe",
					City:      "Chicago",
					Contacts:  "+1234567890",
				}
				mulipartForm, contentType := createMultipartFormJSON(newProfile)

				nr := httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/employer/%d/profile", slug),
					mulipartForm,
				)
				nr.Header.Set("Content-Type", contentType)

				nw := httptest.NewRecorder()
				usecase.profile.
					EXPECT().
					UpdateEmployerProfile(gomock.Any(), slug, newProfile).
					Return(fmt.Errorf("error"))
				return nw, nr
			},
		},
		{
			name:         "GetProfile: ok",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusOK,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				slug := uint64(1)
				newProfile := &dto.JSONUpdateEmployerProfile{
					FirstName: "John",
					LastName:  "Doe",
					City:      "Chicago",
					Contacts:  "+1234567890",
				}
				mulipartForm, contentType := createMultipartFormJSON(newProfile)

				nr := httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/employer/%d/profile", slug),
					mulipartForm,
				)
				nr.Header.Set("Content-Type", contentType)

				nw := httptest.NewRecorder()
				usecase.profile.
					EXPECT().
					UpdateEmployerProfile(gomock.Any(), slug, newProfile).
					Return(nil)
				return nw, nr
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := &usecase{
				profile: mock.NewMockIEmployerUsecase(ctrl),
			}
			tt.w, tt.r = tt.prepare(tt.r, tt.w, usecase)

			app := &internal.App{
				Logger:         logrus.New(),
				BackendAddress: "http://localhost:8080",
				Usecases: &internal.Usecases{
					EmployerUsecase:    usecase.profile,
					VacanciesUsecase:   nil,
					FileLoadingUsecase: nil,
				},
				Microservices: &internal.Microservices{
					Auth: nil,
				},
			}

			h := delivery.NewEmployerHandlers(app)
			require.NotNil(t, h)
			require.NotNil(t, tt.r)
			require.NotNil(t, tt.w)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/employer/{id:[0-9]+}/profile", h.UpdateProfile).Methods(http.MethodPut)

			r.ServeHTTP(tt.w, tt.r)

			require.Equal(t, tt.codeExpected, tt.w.Code)
		})
	}
}

func TestGetVacancies(t *testing.T) {
	t.Parallel()
	type usecase struct {
		profile *vacancies_mock.MockIVacanciesUsecase
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
			name:         "GetVacancies: bad slug",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusInternalServerError,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				slug := "123534657548574856785785346346367542151341354756869568"
				nr := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/employer/%s/vacancies", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "GetProfile: bad usecase",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusInternalServerError,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				slug := uint64(1)

				nr := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/employer/%d/vacancies", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				usecase.profile.
					EXPECT().
					GetVacanciesByEmployerID(gomock.Any(), slug).
					Return(nil, fmt.Errorf("error"))
				return nw, nr
			},
		},
		{
			name:         "GetProfile: ok",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusOK,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				slug := uint64(1)

				var vacancies = []*dto.JSONGetEmployerVacancy{
					{
						ID:                   1,
						Position:             "chemist",
						Salary:               1000,
						Description:          "need a organic chemist",
						Location:             "Texas",
						PositionCategoryName: "chemist",
					},
				}
				nr := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/employer/%d/vacancies", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				usecase.profile.
					EXPECT().
					GetVacanciesByEmployerID(gomock.Any(), slug).
					Return(vacancies, nil)
				return nw, nr
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := &usecase{
				profile: vacancies_mock.NewMockIVacanciesUsecase(ctrl),
			}
			tt.w, tt.r = tt.prepare(tt.r, tt.w, usecase)

			app := &internal.App{
				Logger:         logrus.New(),
				BackendAddress: "http://localhost:8080",
				Usecases: &internal.Usecases{
					EmployerUsecase:    nil,
					VacanciesUsecase:   usecase.profile,
					FileLoadingUsecase: nil,
				},
				Microservices: &internal.Microservices{
					Auth: nil,
				},
			}

			h := delivery.NewEmployerHandlers(app)
			require.NotNil(t, h)
			require.NotNil(t, tt.r)
			require.NotNil(t, tt.w)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/employer/{id:[0-9]+}/vacancies", h.GetEmployerVacancies).Methods(http.MethodGet)

			r.ServeHTTP(tt.w, tt.r)

			require.Equal(t, tt.codeExpected, tt.w.Code)
		})
	}
}

func TestRegistration(t *testing.T) {
	t.Parallel()
	type usecase struct {
		registration *mock.MockIEmployerUsecase
		grpc         *grpc_mock.MockAuthorizationClient
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
			name:         "Registration: bad json",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusBadRequest,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				jsonForm, _ := json.Marshal(
					`{
						"unknown_field": "unknown_value",
						broken field 22r#25,
						[
					}`,
				)
				nr := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/employer/registration",
					bytes.NewReader(jsonForm),
				)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "Registration: bad create usecase",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusInternalServerError,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				form := &dto.JSONEmployerRegistrationForm{
					FirstName:          "Ivan",
					LastName:           "Ivanov",
					Position:           "HR",
					Company:            "Texas Instruments",
					CompanyDescription: "Microelectronics company",
					CompanyWebsite:     "texasinstruments.com",
					Email:              "iY5sG@example.com",
					Password:           "password",
				}

				jsonForm, _ := easyjson.Marshal(form)
				usecase.registration.
					EXPECT().
					Create(gomock.Any(), form).
					Return(nil, errors.New("bad create usecase"))
				nr := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/employer/registration",
					bytes.NewReader(jsonForm),
				)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "Registration: bad grpc",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusInternalServerError,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				form := &dto.JSONEmployerRegistrationForm{
					FirstName:          "Ivan",
					LastName:           "Ivanov",
					Position:           "HR",
					Company:            "Texas Instruments",
					CompanyDescription: "Microelectronics company",
					CompanyWebsite:     "texasinstruments.com",
					Email:              "iY5sG@example.com",
					Password:           "password",
				}

				requestID := "1234567890"
				jsonForm, _ := easyjson.Marshal(form)
				grpc_request := &auth_grpc.AuthRequest{
					RequestID: requestID,
					UserType:  dto.UserTypeEmployer,
					Email:     form.Email,
					Password:  form.Password,
				}
				user := &dto.JSONUser{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}
				usecase.registration.
					EXPECT().
					Create(gomock.Any(), form).
					Return(user, nil)
				usecase.grpc.
					EXPECT().
					AuthUser(gomock.Any(), grpc_request).
					Return(nil, errors.New("bad grpc"))
				nr := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/employer/registration",
					bytes.NewReader(jsonForm),
				).WithContext(
					context.WithValue(r.Context(), dto.RequestIDContextKey, requestID),
				)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "Registration: ok",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusOK,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				form := &dto.JSONEmployerRegistrationForm{
					FirstName:          "Ivan",
					LastName:           "Ivanov",
					Position:           "HR",
					Company:            "Texas Instruments",
					CompanyDescription: "Microelectronics company",
					CompanyWebsite:     "texasinstruments.com",
					Email:              "iY5sG@example.com",
					Password:           "password",
				}

				requestID := "1234567890"
				jsonForm, _ := easyjson.Marshal(form)
				grpc_request := &auth_grpc.AuthRequest{
					RequestID: requestID,
					UserType:  dto.UserTypeEmployer,
					Email:     form.Email,
					Password:  form.Password,
				}
				user := &dto.JSONUser{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}
				grpc_response := &auth_grpc.AuthResponse{
					UserData: &auth_grpc.User{
						UserType: dto.UserTypeEmployer,
						ID:       user.ID,
					},
					Session: &auth_grpc.SessionToken{
						ID: requestID,
					},
				}
				usecase.registration.
					EXPECT().
					Create(gomock.Any(), form).
					Return(user, nil)
				usecase.grpc.
					EXPECT().
					AuthUser(gomock.Any(), grpc_request).
					Return(grpc_response, nil)
				nr := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/employer/registration",
					bytes.NewReader(jsonForm),
				).WithContext(
					context.WithValue(r.Context(), dto.RequestIDContextKey, requestID),
				)
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
				registration: mock.NewMockIEmployerUsecase(ctrl),
				grpc:         grpc_mock.NewMockAuthorizationClient(ctrl),
			}
			tt.w, tt.r = tt.prepare(tt.r, tt.w, usecase)

			app := &internal.App{
				Logger:         logrus.New(),
				BackendAddress: "http://localhost:8080",
				Usecases: &internal.Usecases{
					EmployerUsecase:    usecase.registration,
					VacanciesUsecase:   nil,
					FileLoadingUsecase: nil,
				},
				Microservices: &internal.Microservices{
					Auth: usecase.grpc,
				},
			}

			h := delivery.NewEmployerHandlers(app)
			require.NotNil(t, h)
			require.NotNil(t, tt.r)
			require.NotNil(t, tt.w)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/employer/registration", h.Registration).Methods(http.MethodPost)

			r.ServeHTTP(tt.w, tt.r)

			require.Equal(t, tt.codeExpected, tt.w.Code)
		})
	}
}
