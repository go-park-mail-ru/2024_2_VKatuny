package delivery_test

<<<<<<< HEAD
// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
// 	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/delivery"
// 	applicant_mock "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/mock"
// 	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
// 	session_mock "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/mock"
// 	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
// 	"github.com/golang/mock/gomock"
// 	"github.com/sirupsen/logrus"
// 	"github.com/stretchr/testify/require"
// )

// func TestApplicantRegistration(t *testing.T) {
// 	type in struct {
// 		registrationForm interface{} // dto.ApplicantRegistrationForm
// 		backendURI       string
// 	}
// 	type out struct {
// 		status    int
// 		response  *dto.JSONResponse
// 		sessionID string
// 	}
// 	type usecase struct {
// 		applicant *applicant_mock.MockIApplicantUsecase
// 		session   *session_mock.MockISessionUsecase
// 	}
// 	type args struct {
// 		r *http.Request
// 		w *httptest.ResponseRecorder
// 	}

// 	tests := []struct {
// 		name    string
// 		prepare func(in *in, out *out, usecase *usecase, args *args)
// 	}{
// 		{
// 			name: "ApplicantRegistration creating applicant success",
// 			prepare: func(in *in, out *out, usecase *usecase, args *args) {
// 				registrationForm := &dto.JSONApplicantRegistrationForm{
// 					FirstName: "John",
// 					LastName:  "Doe",
// 					BirthDate: "2000-01-01",
// 					Email:     "john.doe@bk.ru",
// 					Password:  "password",
// 				}
// 				in.registrationForm = registrationForm
// 				in.backendURI = "http://127.0.0.1"

// 				applicantID := uint64(1)
// 				user := &dto.JSONUser{
// 					ID:       applicantID,
// 					UserType: dto.UserTypeApplicant,
// 				}
// 				body := map[string]interface{}{
// 					"id":       float64(applicantID),
// 					"userType": "applicant",
// 				}

// 				out.status = http.StatusOK
// 				out.response = &dto.JSONResponse{
// 					HTTPStatus: http.StatusOK,
// 					Body:       body,
// 				}

// 				usecase.applicant.
// 					EXPECT().
// 					Create(in.registrationForm).
// 					Return(user, nil)

// 				loginForm := &dto.JSONLoginForm{
// 					UserType: dto.UserTypeApplicant,
// 					Email:    registrationForm.Email,
// 					Password: registrationForm.Password,
// 				}

// 				out.sessionID = utils.GenerateSessionToken(utils.TokenLength, dto.UserTypeApplicant)
// 				userWithSession := &dto.UserWithSession{
// 					ID:        applicantID,
// 					UserType:  dto.UserTypeApplicant,
// 					SessionID: out.sessionID,
// 				}

// 				usecase.session.
// 					EXPECT().
// 					Login(loginForm).
// 					Return(userWithSession, nil)

// 				requestBody, _ := json.Marshal(registrationForm)
// 				args.r = httptest.NewRequest(http.MethodPost,
// 					"/api/v1/registration/applicant",
// 					bytes.NewReader(requestBody),
// 				)
// 				args.w = httptest.NewRecorder()
// 			},
// 		},
// 		{
// 			name: "ApplicantRegistration invalid json",
// 			prepare: func(in *in, out *out, usecase *usecase, args *args) {
// 				out.status = http.StatusBadRequest
// 				out.response = &dto.JSONResponse{
// 					HTTPStatus: http.StatusBadRequest,
// 					Error:      dto.MsgInvalidJSON,
// 				}

// 				args.r = httptest.NewRequest(http.MethodPost,
// 					"/api/v1/registration/applicant",
// 					nil,
// 				)
// 				args.w = httptest.NewRecorder()
// 			},
// 		},
// 		{
// 			name: "ApplicantRegistration creating applicant failed",
// 			prepare: func(in *in, out *out, usecase *usecase, args *args) {
// 				registrationForm := &dto.JSONApplicantRegistrationForm{
// 					FirstName: "John",
// 					LastName:  "Doe",
// 					BirthDate: "2000-01-01",
// 					Email:     "john.doe@bk.ru",
// 					Password:  "password",
// 				}
// 				in.registrationForm = registrationForm

// 				out.status = http.StatusInternalServerError
// 				out.response = &dto.JSONResponse{
// 					HTTPStatus: http.StatusInternalServerError,
// 					Error:      dto.MsgUnableToCreateUser,
// 				}

// 				usecase.applicant.
// 					EXPECT().
// 					Create(in.registrationForm).
// 					Return(nil, fmt.Errorf(dto.MsgDataBaseError))

// 				requestBody, _ := json.Marshal(registrationForm)
// 				args.r = httptest.NewRequest(http.MethodPost,
// 					"/api/v1/registration/applicant",
// 					bytes.NewReader(requestBody),
// 				)
// 				args.w = httptest.NewRecorder()
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			logger := logrus.New()
// 			logger.Out = io.Discard

// 			in, out, usecase, args := new(in), new(out), new(usecase), new(args)
// 			usecase.applicant = applicant_mock.NewMockIApplicantUsecase(ctrl)
// 			usecase.session = session_mock.NewMockISessionUsecase(ctrl)

// 			tt.prepare(in, out, usecase, args)

// 			app := &internal.App{
// 				Logger:         logger,
// 				BackendAddress: in.backendURI,
// 				Usecases: &internal.Usecases{
// 					ApplicantUsecase: usecase.applicant,
// 					SessionUsecase:   usecase.session,
// 					PortfolioUsecase: nil,
// 					CVUsecase:        nil,
// 				},
// 			}

// 			mux := http.NewServeMux()
// 			handlers := delivery.NewApplicantProfileHandlers(app)
// 			mux.HandleFunc("/api/v1/registration/applicant", handlers.ApplicantRegistration)

// 			require.NotNil(t, args.r, "request can't be nil")
// 			require.NotNil(t, args.w, "response recorder can't be nil")
// 			mux.ServeHTTP(args.w, args.r)

// 			require.Equal(t, out.status, args.w.Code, "http status doesn't match")

// 			jsonResponse := new(dto.JSONResponse)
// 			err := json.NewDecoder(args.w.Result().Body).Decode(jsonResponse)
// 			require.NoError(t, err, "can't decode response body")

// 			require.Equal(t, out.response, jsonResponse, "responses doesn't match")

// 			if args.w.Code == http.StatusOK {
// 				cookies := args.w.Result().Cookies()
// 				var found bool
// 				var cookieValue string
// 				for _, cookie := range cookies {
// 					if cookie.Name == dto.SessionIDName {
// 						found = true
// 						cookieValue = cookie.Value
// 						break
// 					}
// 				}
// 				require.True(t, found, "cookie not found")
// 				require.Equal(t, out.sessionID, cookieValue, "cookies doesn't match")
// 			}
// 		})
// 	}
// }
=======
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
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/delivery"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/mock"
	cv_mock "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	portfolio_mock "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/mock"
	auth_grpc "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
	grpc_mock "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/mock"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createMultipartForm(jsonForm *dto.JSONUpdateApplicantProfile) (*bytes.Buffer, string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	defer writer.Close()
	writer.WriteField("firstName", jsonForm.FirstName)
	writer.WriteField("lastName", jsonForm.LastName)
	writer.WriteField("city", jsonForm.City)
	writer.WriteField("birthDate", jsonForm.BirthDate)
	writer.WriteField("education", jsonForm.Education)
	writer.WriteField("contacts", jsonForm.Contacts)
	return &buf, writer.FormDataContentType()
}

func TestGetProfileHandler(t *testing.T) {
	t.Parallel()
	type usecase struct {
		profile *mock.MockIApplicantUsecase
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
					fmt.Sprintf("/api/v1/applicant/%s/profile", slug),
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
					fmt.Sprintf("/api/v1/applicant/%d/profile", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				usecase.profile.
					EXPECT().
					GetApplicantProfile(gomock.Any(), slug).
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
					fmt.Sprintf("/api/v1/applicant/%d/profile", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				profile := &dto.JSONGetApplicantProfile{
					ID:        slug,
					FirstName: "John",
					LastName:  "Doe",
				}
				usecase.profile.
					EXPECT().
					GetApplicantProfile(gomock.Any(), slug).
					Return(profile, nil)
				return nw, nr
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := &usecase{
				profile: mock.NewMockIApplicantUsecase(ctrl),
			}
			tt.w, tt.r = tt.prepare(tt.r, tt.w, usecase)

			app := &internal.App{
				Logger:         logrus.New(),
				BackendAddress: "http://localhost:8080",
				Usecases: &internal.Usecases{
					ApplicantUsecase:   usecase.profile,
					CVUsecase:          nil,
					PortfolioUsecase:   nil,
					FileLoadingUsecase: nil,
				},
				Microservices: &internal.Microservices{
					Auth: nil,
				},
			}

			h := delivery.NewApplicantProfileHandlers(app)
			require.NotNil(t, h)
			require.NotNil(t, tt.r)
			require.NotNil(t, tt.w)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/applicant/{id:[0-9]+}/profile", h.GetProfile).Methods(http.MethodGet)
			r.ServeHTTP(tt.w, tt.r)

			require.Equal(t, tt.codeExpected, tt.w.Code)
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	t.Parallel()
	type usecase struct {
		profile *mock.MockIApplicantUsecase
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
					fmt.Sprintf("/api/v1/applicant/%s/profile", slug),
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
				newProfile := &dto.JSONUpdateApplicantProfile{
					FirstName: "John",
					LastName:  "Doe",
					City:      "Chicago",
					BirthDate: "2000-01-01",
					Contacts:  "+1234567890",
				}
				mulipartForm, contentType := createMultipartForm(newProfile)

				nr := httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/applicant/%d/profile", slug),
					mulipartForm,
				)
				nr.Header.Set("Content-Type", contentType)

				nw := httptest.NewRecorder()
				usecase.profile.
					EXPECT().
					UpdateApplicantProfile(gomock.Any(), slug, newProfile).
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
				newProfile := &dto.JSONUpdateApplicantProfile{
					FirstName: "John",
					LastName:  "Doe",
					City:      "Chicago",
					BirthDate: "2000-01-01",
					Contacts:  "+1234567890",
				}
				mulipartForm, contentType := createMultipartForm(newProfile)

				nr := httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/applicant/%d/profile", slug),
					mulipartForm,
				)
				nr.Header.Set("Content-Type", contentType)

				nw := httptest.NewRecorder()
				usecase.profile.
					EXPECT().
					UpdateApplicantProfile(gomock.Any(), slug, newProfile).
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
				profile: mock.NewMockIApplicantUsecase(ctrl),
			}
			tt.w, tt.r = tt.prepare(tt.r, tt.w, usecase)

			app := &internal.App{
				Logger:         logrus.New(),
				BackendAddress: "http://localhost:8080",
				Usecases: &internal.Usecases{
					ApplicantUsecase:   usecase.profile,
					CVUsecase:          nil,
					PortfolioUsecase:   nil,
					FileLoadingUsecase: nil,
				},
				Microservices: &internal.Microservices{
					Auth: nil,
				},
			}

			h := delivery.NewApplicantProfileHandlers(app)
			require.NotNil(t, h)
			require.NotNil(t, tt.r)
			require.NotNil(t, tt.w)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/applicant/{id:[0-9]+}/profile", h.UpdateProfile).Methods(http.MethodPut)

			r.ServeHTTP(tt.w, tt.r)

			require.Equal(t, tt.codeExpected, tt.w.Code)
		})
	}
}

func TestGetCVs(t *testing.T) {
	t.Parallel()
	type usecase struct {
		cv *cv_mock.MockICVsUsecase
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
					fmt.Sprintf("/api/v1/applicant/%s/cv", slug),
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
					fmt.Sprintf("/api/v1/applicant/%d/cv", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				usecase.cv.
					EXPECT().
					GetApplicantCVs(slug).
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

				var vacancies = []*dto.JSONGetApplicantCV{
					{
						ID:                   1,
						ApplicantID:          1,
						PositionRu:           "химик",
						PositionEn:           "chemist",
						Description:          "нужен химик",
						JobSearchStatus:      "searching",
						WorkingExperience:    "1 year",
						PositionCategoryName: "chemistry",
					},
				}
				nr := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/applicant/%d/cv", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				usecase.cv.
					EXPECT().
					GetApplicantCVs(slug).
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
				cv: cv_mock.NewMockICVsUsecase(ctrl),
			}
			tt.w, tt.r = tt.prepare(tt.r, tt.w, usecase)

			app := &internal.App{
				Logger:         logrus.New(),
				BackendAddress: "http://localhost:8080",
				Usecases: &internal.Usecases{
					ApplicantUsecase:   nil,
					CVUsecase:          usecase.cv,
					PortfolioUsecase:   nil,
					FileLoadingUsecase: nil,
				},
				Microservices: &internal.Microservices{
					Auth: nil,
				},
			}

			h := delivery.NewApplicantProfileHandlers(app)
			require.NotNil(t, h)
			require.NotNil(t, tt.r)
			require.NotNil(t, tt.w)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/applicant/{id:[0-9]+}/cv", h.GetCVs).Methods(http.MethodGet)

			r.ServeHTTP(tt.w, tt.r)

			require.Equal(t, tt.codeExpected, tt.w.Code)
		})
	}
}

func TestGetPortfolios(t *testing.T) {
	t.Parallel()
	type usecase struct {
		portfolio *portfolio_mock.MockIPortfolioUsecase
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
			name:         "GetPortfolio: bad slug",
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
					fmt.Sprintf("/api/v1/applicant/%s/portfolio", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				return nw, nr
			},
		},
		{
			name:         "GetPortfolio: bad usecase",
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
					fmt.Sprintf("/api/v1/applicant/%d/portfolio", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				usecase.portfolio.
					EXPECT().
					GetApplicantPortfolios(gomock.Any(), slug).
					Return(nil, fmt.Errorf("error"))
				return nw, nr
			},
		},
		{
			name:         "GetPortfolio: ok",
			r:            new(http.Request),
			w:            new(httptest.ResponseRecorder),
			codeExpected: http.StatusOK,
			prepare: func(
				r *http.Request,
				w *httptest.ResponseRecorder,
				usecase *usecase,
			) (*httptest.ResponseRecorder, *http.Request) {
				slug := uint64(1)

				var vacancies = []*dto.JSONGetApplicantPortfolio{
					{
						ID:          1,
						ApplicantID: 1,
						Name:        "bass guitar",
						Description: "my best guitar",
					},
				}
				nr := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/applicant/%d/portfolio", slug),
					nil,
				)
				nw := httptest.NewRecorder()
				usecase.portfolio.
					EXPECT().
					GetApplicantPortfolios(gomock.Any(), slug).
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
				portfolio: portfolio_mock.NewMockIPortfolioUsecase(ctrl),
			}
			tt.w, tt.r = tt.prepare(tt.r, tt.w, usecase)

			app := &internal.App{
				Logger:         logrus.New(),
				BackendAddress: "http://localhost:8080",
				Usecases: &internal.Usecases{
					ApplicantUsecase:   nil,
					CVUsecase:          nil,
					PortfolioUsecase:   usecase.portfolio,
					FileLoadingUsecase: nil,
				},
				Microservices: &internal.Microservices{
					Auth: nil,
				},
			}

			h := delivery.NewApplicantProfileHandlers(app)
			require.NotNil(t, h)
			require.NotNil(t, tt.r)
			require.NotNil(t, tt.w)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/applicant/{id:[0-9]+}/portfolio", h.GetPortfolios).Methods(http.MethodGet)

			r.ServeHTTP(tt.w, tt.r)

			require.Equal(t, tt.codeExpected, tt.w.Code)
		})
	}
}

func TestRegistration(t *testing.T) {
	t.Parallel()
	type usecase struct {
		registration *mock.MockIApplicantUsecase
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
					"/api/v1/applicant/registration",
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
				form := &dto.JSONApplicantRegistrationForm{
					FirstName: "Ivan",
					LastName:  "Ivanov",
					BirthDate: "2000-01-01",
					Email:     "iY5sG@example.com",
					Password:  "password",
				}

				jsonForm, _ := json.Marshal(form)
				usecase.registration.
					EXPECT().
					Create(gomock.Any(), form).
					Return(nil, errors.New("bad create usecase"))
				nr := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/applicant/registration",
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
				form := &dto.JSONApplicantRegistrationForm{
					FirstName: "Ivan",
					LastName:  "Ivanov",
					BirthDate: "2000-01-01",
					Email:     "iY5sG@example.com",
					Password:  "password",
				}

				requestID := "1234567890"
				jsonForm, _ := json.Marshal(form)
				grpc_request := &auth_grpc.AuthRequest{
					RequestID: requestID,
					UserType:  dto.UserTypeApplicant,
					Email:     form.Email,
					Password:  form.Password,
				}
				user := &dto.JSONUser{
					ID:       1,
					UserType: dto.UserTypeApplicant,
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
					"/api/v1/applicant/registration",
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
				form := &dto.JSONApplicantRegistrationForm{
					FirstName: "Ivan",
					LastName:  "Ivanov",
					BirthDate: "2000-01-01",
					Email:     "iY5sG@example.com",
					Password:  "password",
				}

				requestID := "1234567890"
				jsonForm, _ := json.Marshal(form)
				grpc_request := &auth_grpc.AuthRequest{
					RequestID: requestID,
					UserType:  dto.UserTypeApplicant,
					Email:     form.Email,
					Password:  form.Password,
				}
				user := &dto.JSONUser{
					ID:       1,
					UserType: dto.UserTypeApplicant,
				}
				grpc_response := &auth_grpc.AuthResponse{
					UserData: &auth_grpc.User{
						UserType: dto.UserTypeApplicant,
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
					"/api/v1/applicant/registration",
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
				registration: mock.NewMockIApplicantUsecase(ctrl),
				grpc:         grpc_mock.NewMockAuthorizationClient(ctrl),
			}
			tt.w, tt.r = tt.prepare(tt.r, tt.w, usecase)

			app := &internal.App{
				Logger:         logrus.New(),
				BackendAddress: "http://localhost:8080",
				Usecases: &internal.Usecases{
					ApplicantUsecase:   usecase.registration,
					CVUsecase:          nil,
					PortfolioUsecase:   nil,
					FileLoadingUsecase: nil,
				},
				Microservices: &internal.Microservices{
					Auth: usecase.grpc,
				},
			}

			h := delivery.NewApplicantProfileHandlers(app)
			require.NotNil(t, h)
			require.NotNil(t, tt.r)
			require.NotNil(t, tt.w)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/applicant/registration", h.ApplicantRegistration).Methods(http.MethodPost)

			r.ServeHTTP(tt.w, tt.r)

			require.Equal(t, tt.codeExpected, tt.w.Code)
		})
	}
}
>>>>>>> dev
