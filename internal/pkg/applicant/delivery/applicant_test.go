package delivery_test

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
