package delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	applicant_mock "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/delivery"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	file_loading_mock "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading/mock"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createMultipartFormJSONCV(jsonForm *dto.JSONCv) (*bytes.Buffer, string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	defer writer.Close()
	writer.WriteField("positionRu", jsonForm.PositionRu)
	writer.WriteField("positionEn", jsonForm.PositionEn)
	writer.WriteField("description", jsonForm.Description)
	writer.WriteField("jobSearchStatus", jsonForm.JobSearchStatusName)
	writer.WriteField("workingExperience", jsonForm.WorkingExperience)
	writer.WriteField("group", jsonForm.PositionCategoryName)
	return &buf, writer.FormDataContentType()
}

func TestCreateCVHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}
	type dependencies struct {
		cvsUsecase *mock.MockICVsUsecase
		logger     *logrus.Logger

		cv          *dto.JSONCv
		currentUser *dto.UserFromSession

		args args
	}
	tests := []struct {
		name               string
		prepare            func(f *dependencies)
		wantErr            bool
		expectedStatusCode int
	}{
		{
			name: "CVsHandler.CreateCVHandler successful generation of cv",
			prepare: func(f *dependencies) {
				f.currentUser = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeApplicant,
				}
				f.cv = &dto.JSONCv{
					ID:                  1,
					ApplicantID:         1,
					PositionRu:          "Мок Должность",
					PositionEn:          "Mock Position",
					Description:         "Mock Description",
					JobSearchStatusName: "Ищу",
					WorkingExperience:   "1 год",
					Avatar:              "Mock Avatar",
				}

				multipartForm, contentType := createMultipartFormJSONCV(f.cv)

				// disable logging
				f.logger = logrus.New()
				f.logger.Out = io.Discard

				f.cvsUsecase.
					EXPECT().
					CreateCV(gomock.Any(), f.currentUser).
					Return(f.cv, nil)

				f.args.r = httptest.NewRequest(
					http.MethodPost,
					"/api/v1/cv",
					multipartForm,
				).WithContext(
					context.WithValue(
						context.Background(),
						dto.UserContextKey,
						f.currentUser,
					),
				)
				f.args.r.Header.Set("Content-Type", contentType)

				f.args.w = httptest.NewRecorder()
			},
			wantErr:            false,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "CVsHandler.CreateCVHandler got invalid multipart",
			prepare: func(f *dependencies) {
				f.cv = &dto.JSONCv{}
				f.currentUser = &dto.UserFromSession{}

				f.logger = logrus.New()
				f.logger.Out = io.Discard

				multipartForm, contentType := createMultipartFormJSONCV(f.cv)
				f.args.r = httptest.NewRequest(
					http.MethodPost,
					"/api/v1/cv",
					multipartForm,
				).WithContext(
					context.WithValue(
						context.Background(),
						dto.UserContextKey,
						f.currentUser,
					),
				)

				f.cvsUsecase.
					EXPECT().
					CreateCV(gomock.Any(), f.currentUser).
					Return(f.cv, fmt.Errorf(dto.MsgDataBaseError))

				f.args.r.Header.Set("Content-Type", contentType)
				f.args.w = httptest.NewRecorder()
			},
			wantErr:            true,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "CVsHandler.CreateCVHandler no user provided",
			prepare: func(f *dependencies) {
				f.cv = &dto.JSONCv{
					ID:                  1,
					ApplicantID:         1,
					PositionRu:          "Мок Должность",
					PositionEn:          "Mock Position",
					Description:         "Mock Description",
					JobSearchStatusName: "Ищу",
					WorkingExperience:   "1 год",
					Avatar:              "Mock Avatar",
				}

				// disable logging
				f.logger = logrus.New()
				f.logger.Out = io.Discard

				multipartForm, contentType := createMultipartFormJSONCV(f.cv)

				f.args.r = httptest.NewRequest(
					http.MethodPost,
					"/api/v1/cv",
					multipartForm,
				)
				f.args.r.Header.Set("Content-Type", contentType)

				f.args.w = httptest.NewRecorder()
			},
			wantErr:            true,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "CVsHandler.CreateCVHandler usecase returns error",
			prepare: func(f *dependencies) {
				f.currentUser = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeApplicant,
				}
				f.cv = &dto.JSONCv{
					ID:                  1,
					ApplicantID:         1,
					PositionRu:          "Мок Должность",
					PositionEn:          "Mock Position",
					Description:         "Mock Description",
					JobSearchStatusName: "Ищу",
					WorkingExperience:   "1 год",
					Avatar:              "Mock Avatar",
				}

				// disable logging
				f.logger = logrus.New()
				f.logger.Out = io.Discard

				f.cvsUsecase.
					EXPECT().
					CreateCV(gomock.Any(), f.currentUser).
					Return(nil, fmt.Errorf(dto.MsgDataBaseError))

				multipartForm, contentType := createMultipartFormJSONCV(f.cv)

				f.args.r = httptest.NewRequest(
					http.MethodPost,
					"/api/v1/cv",
					multipartForm,
				).WithContext(
					context.WithValue(
						context.Background(),
						dto.UserContextKey,
						f.currentUser,
					),
				)
				f.args.r.Header.Set("Content-Type", contentType)

				f.args.w = httptest.NewRecorder()
			},
			wantErr:            true,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d := dependencies{
				cvsUsecase: mock.NewMockICVsUsecase(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&d)
			}

			app := &internal.App{
				Usecases: &internal.Usecases{
					CVUsecase: d.cvsUsecase,
				},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
				Logger: d.logger,
			}

			h := delivery.NewCVsHandler(app)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/cv", h.CreateCV).Methods(http.MethodPost)

			r.ServeHTTP(d.args.w, d.args.r)

			status := d.args.w.Result().StatusCode
			require.EqualValuesf(t, tt.expectedStatusCode, status,
				"got status %d, expected %d. Response body: %s",
				status,
				tt.expectedStatusCode,
				d.args.w.Body.String(),
			)
			if tt.wantErr {
				switch tt.expectedStatusCode {
				case http.StatusBadRequest:
					body := d.args.w.Result().Body
					jsonData, err := io.ReadAll(body)
					require.NoError(t, err)

					gotJson := new(dto.JSONResponse)
					err = json.Unmarshal(jsonData, gotJson)
					require.NoError(t, err)

					require.EqualExportedValues(t, &dto.JSONResponse{
						HTTPStatus: http.StatusBadRequest,
						Error:      dto.MsgInvalidJSON,
					}, gotJson, "got unexpected json response body")
				case http.StatusUnauthorized:
					body := d.args.w.Result().Body
					jsonData, err := io.ReadAll(body)
					require.NoError(t, err)

					gotJson := new(dto.JSONResponse)
					err = json.Unmarshal(jsonData, gotJson)
					require.NoError(t, err)

					require.EqualExportedValues(t, &dto.JSONResponse{
						HTTPStatus: http.StatusUnauthorized,
						Error:      dto.MsgUnableToGetUserFromContext,
					}, gotJson, "expected response that there is no user from session")
				case http.StatusInternalServerError:
					body := d.args.w.Result().Body
					jsonData, err := io.ReadAll(body)
					require.NoError(t, err)

					gotJson := new(dto.JSONResponse)
					err = json.Unmarshal(jsonData, gotJson)
					require.NoError(t, err)

					require.EqualExportedValues(t, &dto.JSONResponse{
						HTTPStatus: http.StatusInternalServerError,
						Error:      dto.MsgDataBaseError,
					}, gotJson, "expected usecase error")
				}
			}
		})
	}
}

func TestGetCVHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}
	type dependencies struct {
		cvsUsecase *mock.MockICVsUsecase

		cv dto.JSONCv

		args args
	}
	tests := []struct {
		name               string
		prepare            func(f *dependencies)
		expectedStatusCode int
		expectedBody       *dto.JSONResponse
	}{
		{
			name: "CVsHandler.GetCVHandler successfully got the cv",
			prepare: func(f *dependencies) {
				IDslug := uint64(1)

				f.cv = dto.JSONCv{
					ID:                  1,
					ApplicantID:         1,
					PositionRu:          "Мок Должность",
					PositionEn:          "Mock Position",
					Description:         "Mock Description",
					JobSearchStatusName: "Ищу",
					WorkingExperience:   "1 год",
					Avatar:              "Mock Avatar",
					CreatedAt:           "2022-02-02",
					UpdatedAt:           "2022-02-05",
				}

				f.cvsUsecase.
					EXPECT().
					GetCV(IDslug).
					Return(&f.cv, nil)

				f.args.r = httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/cv/%d", IDslug),
					nil,
				)
				f.args.w = httptest.NewRecorder()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "CVsHandler.GetCVHandler db err",
			prepare: func(f *dependencies) {
				IDslug := uint64(1)

				f.cvsUsecase.
					EXPECT().
					GetCV(IDslug).
					Return(nil, fmt.Errorf(dto.MsgDataBaseError))

				f.args.r = httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/cv/%d", IDslug),
					nil,
				)
				f.args.w = httptest.NewRecorder()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d := dependencies{
				cvsUsecase: mock.NewMockICVsUsecase(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(&d)
			}

			logger := logger.NewLogrusLogger()
			logger.Out = io.Discard

			app := &internal.App{
				Logger: logger,
				Usecases: &internal.Usecases{
					CVUsecase: d.cvsUsecase,
				},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
			}

			h := delivery.NewCVsHandler(app)
			r := mux.NewRouter()
			r.HandleFunc("/api/v1/cv/{id:[0-9]+}", h.GetCV).Methods(http.MethodGet)
			r.ServeHTTP(d.args.w, d.args.r)

			status := d.args.w.Result().StatusCode

			require.EqualValuesf(t, tt.expectedStatusCode, status,
				"got status %d, expected %d. Response body: %s",
				status,
				tt.expectedStatusCode,
				d.args.w.Body.String(),
			)
		})
	}
}

func TestUpdateCVHandler(t *testing.T) {
	t.Parallel()
	type in struct {
		updatedCV   *dto.JSONCv
		currentUser *dto.UserFromSession
		slug        string
	}
	type outExpected struct {
		response *dto.JSONResponse
		status   int
	}
	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}
	type usecaseMock struct {
		cvsUsecase *mock.MockICVsUsecase
	}
	tests := []struct {
		name    string
		prepare func(in *in, out *outExpected, usecase *usecaseMock, args *args)
	}{
		{
			name: "CVsHandler.UpdateCVHandler success",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.updatedCV = &dto.JSONCv{
					ID:                   1,
					ApplicantID:          1,
					PositionRu:           "Мок Должность",
					PositionEn:           "Mock Position",
					Description:          "Mock Description",
					JobSearchStatusName:  "Больше не ищу",
					PositionCategoryName: "group",
					WorkingExperience:    "1 год",
				}
				in.currentUser = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeApplicant,
				}
				in.slug = "1"

				multipartForm, contentType := createMultipartFormJSONCV(in.updatedCV)

				slugInt, _ := strconv.Atoi(in.slug)

				expectedCV := map[string]interface{}{
					"id":                float64(1),
					"applicant":         float64(1),
					"positionRu":        "Мок Должность",
					"positionEn":        "Mock Position",
					"description":       "Mock Description",
					"jobSearchStatus":   "Больше не ищу",
					"positionGroup":     "group",
					"workingExperience": "1 год",
					"avatar":            "",
					"compressedAvatar":  "",
					"createdAt":         "",
					"updatedAt":         "",
				}
				out.response = &dto.JSONResponse{
					HTTPStatus: http.StatusOK,
					Body:       expectedCV,
				}
				out.status = http.StatusOK

				usecase.cvsUsecase.
					EXPECT().
					UpdateCV(uint64(slugInt), in.currentUser, gomock.Any()).
					Return(in.updatedCV, nil)

				args.r = httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/cv/%s", in.slug),
					multipartForm,
				).WithContext(
					context.WithValue(
						context.Background(),
						dto.UserContextKey,
						in.currentUser,
					),
				)
				args.r.Header.Set("Content-Type", contentType)

				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "CVsHandler.UpdateCVHandler can't get user from context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"

				out.status = http.StatusUnauthorized
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgUnauthorized,
				}

				args.r = httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/cv/%s", in.slug),
					nil,
				)
				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "CVsHandler.UpdateCVHandler usecase returned internal error",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.updatedCV = &dto.JSONCv{
					ID:                  1,
					PositionRu:          "Мок Должность",
					PositionEn:          "Mock Position",
					Description:         "Mock Description",
					JobSearchStatusName: "Больше не ищу",
					WorkingExperience:   "1 год",
					Avatar:              "Mock Avatar",
					CreatedAt:           "2022-02-02",
					UpdatedAt:           "2022-02-05",
				}
				in.currentUser = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeApplicant,
				}
				in.slug = "1451566565115655515615645645656156156156156156156156156"

				multipartForm, contentType := createMultipartFormJSONCV(in.updatedCV)

				// slugInt, _ := strconv.Atoi(in.slug)

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
				}

				// updatedCVJSON, _ := json.Marshal(in.updatedCV)

				// usecase.cvsUsecase.
				// 	EXPECT().
				// 	UpdateCV(uint64(slugInt), in.currentUser, in.updatedCV).
				// 	Return(nil, fmt.Errorf(dto.MsgDataBaseError))

				args.r = httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/cv/%s", in.slug),
					multipartForm,
				).WithContext(
					context.WithValue(
						context.Background(),
						dto.UserContextKey,
						in.currentUser,
					),
				)
				args.r.Header.Set("Content-Type", contentType)

				args.w = httptest.NewRecorder()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			in, out, usecase, args := new(in), new(outExpected), new(usecaseMock), new(args)

			usecase.cvsUsecase = mock.NewMockICVsUsecase(ctrl)

			tt.prepare(in, out, usecase, args)

			logger := logrus.New()
			logger.Out = io.Discard

			app := &internal.App{
				Logger: logger,
				Usecases: &internal.Usecases{
					CVUsecase: usecase.cvsUsecase,
				},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
			}

			h := delivery.NewCVsHandler(app)
			r := mux.NewRouter()
			r.HandleFunc("/api/v1/cv/{id:[0-9]+}", h.UpdateCV)

			require.NotNil(t, args.r, "request is nil")
			require.NotNil(t, args.w, "response is nil")

			r.ServeHTTP(args.w, args.r)

			require.EqualValuesf(t, out.status, args.w.Result().StatusCode,
				"got status %d, expected %d",
				args.w.Result().StatusCode,
				out.status,
			)

			jsonResonse := new(dto.JSONResponse)
			err := easyjson.UnmarshalFromReader(args.w.Result().Body, jsonResonse)
			require.NoError(t, err)

			require.Equalf(t, out.response, jsonResonse,
				"got response %v, expected %v",
				jsonResonse,
				out.response,
			)
		})
	}
}

func TestDeleteCVHandler(t *testing.T) {
	t.Parallel()

	type in struct {
		slug string
		user *dto.UserFromSession
	}
	type outExpected struct {
		response *dto.JSONResponse
		status   int
	}
	type usecaseMock struct {
		cvsUsecase *mock.MockICVsUsecase
	}
	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}
	tests := []struct {
		name    string
		prepare func(in *in, out *outExpected, usecase *usecaseMock, args *args)
	}{
		{
			name: "CVHandler.DeleteCVHandler success",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"
				in.user = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeApplicant,
				}

				out.status = http.StatusOK
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
				}

				slugInt, _ := strconv.Atoi(in.slug)

				usecase.cvsUsecase.
					EXPECT().
					DeleteCV(uint64(slugInt), in.user).
					Return(nil)

				args.r = httptest.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/api/v1/cv/%s", in.slug),
					nil,
				).WithContext(
					context.WithValue(
						context.Background(),
						dto.UserContextKey,
						in.user,
					),
				)
				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "CVHandler.DeleteCVHandler no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"

				out.status = http.StatusUnauthorized
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgUnableToGetUserFromContext,
				}

				args.r = httptest.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/api/v1/cv/%s", in.slug),
					nil,
				)
				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "CVHandler.DeleteCVHandler no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"
				in.user = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeApplicant,
				}

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgDataBaseError,
				}

				slugInt, _ := strconv.Atoi(in.slug)

				usecase.cvsUsecase.
					EXPECT().
					DeleteCV(uint64(slugInt), in.user).
					Return(fmt.Errorf(dto.MsgDataBaseError))

				args.r = httptest.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/api/v1/cv/%s", in.slug),
					nil,
				).WithContext(
					context.WithValue(
						context.Background(),
						dto.UserContextKey,
						in.user,
					),
				)
				args.w = httptest.NewRecorder()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			in, out, usecase, args := new(in), new(outExpected), new(usecaseMock), new(args)

			usecase.cvsUsecase = mock.NewMockICVsUsecase(ctrl)

			tt.prepare(in, out, usecase, args)

			logger := logrus.New()
			logger.Out = io.Discard

			app := &internal.App{
				Logger: logger,
				Usecases: &internal.Usecases{
					CVUsecase: usecase.cvsUsecase,
				},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
			}

			h := delivery.NewCVsHandler(app)
			r := mux.NewRouter()
			r.HandleFunc("/api/v1/cv/{id:[0-9]+}", h.DeleteCV)

			require.NotNil(t, args.r, "request is nil")
			require.NotNil(t, args.w, "response is nil")

			r.ServeHTTP(args.w, args.r)

			require.EqualValuesf(t, out.status, args.w.Result().StatusCode,
				"got status %d, expected %d",
				args.w.Result().StatusCode,
				out.status,
			)

			jsonResonse := new(dto.JSONResponse)
			err := easyjson.UnmarshalFromReader(args.w.Result().Body, jsonResonse)
			require.NoError(t, err)

			require.Equalf(t, out.response, jsonResonse,
				"got response %v, expected %v",
				jsonResonse,
				out.response,
			)
		})
	}
}

func TestCVtoPDF(t *testing.T) {
	t.Parallel()

	type in struct {
		slug string
		user *dto.UserFromSession
	}
	type outExpected struct {
		response *dto.JSONResponse
		status   int
	}
	type usecaseMock struct {
		cvsUsecase         *mock.MockICVsUsecase
		applicantUsecase   *applicant_mock.MockIApplicantUsecase
		fileLoadingUsecase *file_loading_mock.MockIFileLoadingUsecase
	}
	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}
	tests := []struct {
		name    string
		prepare func(in *in, out *outExpected, usecase *usecaseMock, args *args)
	}{
		{
			name: "CVHandler.DeleteCVHandler success",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"
				applicant := &dto.JSONGetApplicantProfile{
					ID:        0,
					FirstName: "Mock Applicant",
				}
				cv := &dto.JSONCv{
					ID:          0,
					ApplicantID: 1,
					PositionRu:  "Мок Должность",
					PositionEn:  "Mock Position",
					Description: "Мок Описание",
				}
				res := &dto.CVPDFFile{
					FileName: "Mock CV.pdf",
				}
				in.user = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeApplicant,
				}

				out.status = http.StatusOK
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Body: map[string]interface{}{
						"FileName": res.FileName,
					},
				}
				slugInt, _ := strconv.Atoi(in.slug)

				usecase.cvsUsecase.
					EXPECT().
					GetCV(uint64(slugInt)).
					Return(cv, nil)
				usecase.applicantUsecase.
					EXPECT().
					GetApplicantProfile(context.Background(), uint64(cv.ApplicantID)).
					Return(applicant, nil)
				usecase.fileLoadingUsecase.
					EXPECT().
					CVtoPDF(cv, applicant).
					Return(res, nil)

				args.r = httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/cv-to-pdf/%s", in.slug),
					nil,
				)
				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "CVHandler.DeleteCVHandler success",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.user = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeApplicant,
				}

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
				}

				args.r = httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/cv-to-pdf/%s", "111111111111111111111111111111111111111"),
					nil,
				)
				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "CVHandler.DeleteCVHandler success",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"
				in.user = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeApplicant,
				}

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      "not found",
				}
				slugInt, _ := strconv.Atoi(in.slug)

				usecase.cvsUsecase.
					EXPECT().
					GetCV(uint64(slugInt)).
					Return(nil, errors.New("not found"))

				args.r = httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/cv-to-pdf/%s", in.slug),
					nil,
				)
				args.w = httptest.NewRecorder()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			in, out, usecase, args := new(in), new(outExpected), new(usecaseMock), new(args)

			usecase.cvsUsecase = mock.NewMockICVsUsecase(ctrl)
			usecase.applicantUsecase = applicant_mock.NewMockIApplicantUsecase(ctrl)
			usecase.fileLoadingUsecase = file_loading_mock.NewMockIFileLoadingUsecase(ctrl)

			tt.prepare(in, out, usecase, args)

			logger := logrus.New()
			logger.Out = io.Discard

			app := &internal.App{
				Logger: logger,
				Usecases: &internal.Usecases{
					CVUsecase:          usecase.cvsUsecase,
					ApplicantUsecase:   usecase.applicantUsecase,
					FileLoadingUsecase: usecase.fileLoadingUsecase,
				},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
			}

			h := delivery.NewCVsHandler(app)
			r := mux.NewRouter()
			r.HandleFunc("/api/v1/cv-to-pdf/{id:[0-9]+}", h.CVtoPDF)

			require.NotNil(t, args.r, "request is nil")
			require.NotNil(t, args.w, "response is nil")

			r.ServeHTTP(args.w, args.r)

			require.EqualValuesf(t, out.status, args.w.Result().StatusCode,
				"got status %d, expected %d",
				args.w.Result().StatusCode,
				out.status,
			)

			jsonResonse := new(dto.JSONResponse)
			err := easyjson.UnmarshalFromReader(args.w.Result().Body, jsonResonse)
			require.NoError(t, err)
			require.Equalf(t, out.response, jsonResonse,
				"got response %v, expected %v",
				jsonResonse,
				out.response,
			)
		})
	}
}

// TODO: implement tests for SearchCVHandler
