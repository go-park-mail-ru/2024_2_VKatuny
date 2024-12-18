package delivery_test

import (
	//"bytes"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"

	//"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	applicant_mock "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	file_loading_mock "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/delivery"
	vacancies_mock "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/mock"
	notifications_grpc "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/generated"
	grpc_mock "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/mock"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	//"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	//"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateVacancyHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}
	type dependencies struct {
		vacanciesUsecase   *vacancies_mock.MockIVacanciesUsecase
		fileLoadingUsecase *file_loading_mock.MockIFileLoadingUsecase
		logger             *logrus.Logger

		vacancy     interface{}
		currentUser dto.UserFromSession

		args args
	}
	tests := []struct {
		name               string
		prepare            func(f *dependencies)
		wantErr            bool
		expectedStatusCode int
	}{
		{
			name: "VacanciesHandler.CreateVacancyHandler successful generation of vacancy",
			prepare: func(f *dependencies) {
				f.currentUser = dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}
				f.vacancy = &dto.JSONVacancy{
					EmployerID:           0,
					Salary:               1111,
					Position:             "Mock Position",
					Location:             "Mock Location",
					Description:          "Mock Description",
					WorkType:             "Mock Work Type",
					PositionCategoryName: "Mock Position Category",
					//Avatar:               "Mock Avatar",
				}

				// disable logging
				f.logger = logrus.New()
				f.logger.Out = io.Discard

				f.vacanciesUsecase.
					EXPECT().
					CreateVacancy(gomock.Any(), f.vacancy, &f.currentUser).
					Return(f.vacancy, nil)

				//body, _ := json.Marshal(f.vacancy)
				var buf bytes.Buffer
				multipartWriter := multipart.NewWriter(&buf)
				defer multipartWriter.Close()
				// add form field
				CreateOneMultipart(multipartWriter, "salary", "1111")
				CreateOneMultipart(multipartWriter, "position", "Mock Position")
				CreateOneMultipart(multipartWriter, "location", "Mock Location")
				CreateOneMultipart(multipartWriter, "description", "Mock Description")
				CreateOneMultipart(multipartWriter, "workType", "Mock Work Type")
				CreateOneMultipart(multipartWriter, "group", "Mock Position Category")
				//filePart, _ := multipartWriter.CreateFormFile("company_avatar", "11.jpg")
				//file, _ := os.Open("11.jpg")
				//defer file.Close()
				//data := make([]byte, 1)
				//file.Read(data)
				//fmt.Println("!!!!!", data)
				//filePart.Write(data)
				f.args.r = httptest.NewRequest(
					http.MethodPost,
					"/api/v1/vacancy",
					&buf,
				).WithContext(
					context.WithValue(
						context.Background(),
						dto.UserContextKey,
						&f.currentUser,
					),
				)
				f.args.r.Header.Set("Content-Type", "multipart/form-data"+"; boundary="+multipartWriter.Boundary())
				f.args.w = httptest.NewRecorder()
			},
			wantErr:            false,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "VacanciesHandler.CreateVacancyHandler successful generation of vacancy",
			prepare: func(f *dependencies) {
				f.vacancy = &dto.JSONVacancy{
					EmployerID:           0,
					Salary:               1111,
					Position:             "Mock Position",
					Location:             "Mock Location",
					Description:          "Mock Description",
					WorkType:             "Mock Work Type",
					PositionCategoryName: "Mock Position Category",
					//Avatar:               "Mock Avatar",
				}

				// disable logging
				f.logger = logrus.New()
				f.logger.Out = io.Discard

				//body, _ := json.Marshal(f.vacancy)
				var buf bytes.Buffer
				multipartWriter := multipart.NewWriter(&buf)
				defer multipartWriter.Close()
				// add form field
				CreateOneMultipart(multipartWriter, "salary", "1111")
				CreateOneMultipart(multipartWriter, "position", "Mock Position")
				CreateOneMultipart(multipartWriter, "location", "Mock Location")
				CreateOneMultipart(multipartWriter, "description", "Mock Description")
				CreateOneMultipart(multipartWriter, "workType", "Mock Work Type")
				CreateOneMultipart(multipartWriter, "group", "Mock Position Category")
				//filePart, _ := multipartWriter.CreateFormFile("company_avatar", "file.jpg")
				//filePart.Write([]byte("0x5fd3bc"))
				f.args.r = httptest.NewRequest(
					http.MethodPost,
					"/api/v1/vacancy",
					&buf,
				)
				f.args.r.Header.Set("Content-Type", "multipart/form-data"+"; boundary="+multipartWriter.Boundary())
				f.args.w = httptest.NewRecorder()
			},
			wantErr:            true,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "VacanciesHandler.CreateVacancyHandler got invalid json",
			prepare: func(f *dependencies) {
				f.vacancy = struct {
					ID string `json:"id"`
				}{ID: "1"}

				f.logger = logrus.New()
				f.logger.Out = io.Discard

				body, _ := json.Marshal(f.vacancy)

				f.args.r = httptest.NewRequest(
					http.MethodPost,
					"/api/v1/vacancy",
					bytes.NewReader(body),
				)
				f.args.w = httptest.NewRecorder()
			},
			wantErr:            true,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "VacanciesHandler.CreateVacancyHandler successful generation of vacancy",
			prepare: func(f *dependencies) {
				f.currentUser = dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}
				f.vacancy = &dto.JSONVacancy{
					EmployerID:           0,
					Salary:               1111,
					Position:             "Mock Position",
					Location:             "Mock Location",
					Description:          "Mock Description",
					WorkType:             "Mock Work Type",
					PositionCategoryName: "Mock Position Category",
					//Avatar:               "Mock Avatar",
				}

				// disable logging
				f.logger = logrus.New()
				f.logger.Out = io.Discard

				f.vacanciesUsecase.
					EXPECT().
					CreateVacancy(gomock.Any(), f.vacancy, &f.currentUser).
					Return(nil, fmt.Errorf(dto.MsgDataBaseError))

				//body, _ := json.Marshal(f.vacancy)
				var buf bytes.Buffer
				multipartWriter := multipart.NewWriter(&buf)
				defer multipartWriter.Close()
				// add form field
				CreateOneMultipart(multipartWriter, "salary", "1111")
				CreateOneMultipart(multipartWriter, "position", "Mock Position")
				CreateOneMultipart(multipartWriter, "location", "Mock Location")
				CreateOneMultipart(multipartWriter, "description", "Mock Description")
				CreateOneMultipart(multipartWriter, "workType", "Mock Work Type")
				CreateOneMultipart(multipartWriter, "group", "Mock Position Category")
				f.args.r = httptest.NewRequest(
					http.MethodPost,
					"/api/v1/vacancy",
					&buf,
				).WithContext(
					context.WithValue(
						context.Background(),
						dto.UserContextKey,
						&f.currentUser,
					),
				)
				f.args.r.Header.Set("Content-Type", "multipart/form-data"+"; boundary="+multipartWriter.Boundary())
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
				vacanciesUsecase: vacancies_mock.NewMockIVacanciesUsecase(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&d)
			}

			app := &internal.App{
				Usecases: &internal.Usecases{
					VacanciesUsecase: d.vacanciesUsecase,
				},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
				Logger: d.logger,
			}

			h := delivery.NewVacanciesHandlers(app)

			r := mux.NewRouter()
			r.HandleFunc("/api/v1/vacancy", h.CreateVacancy).Methods(http.MethodPost)

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
					err = easyjson.Unmarshal(jsonData, gotJson)
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
					err = easyjson.Unmarshal(jsonData, gotJson)
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
					err = easyjson.Unmarshal(jsonData, gotJson)
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
func CreateOneMultipart(multipartWriter *multipart.Writer, name, content string) {
	position, _ := multipartWriter.CreateFormField(name)
	position.Write([]byte(content))
}
func TestGetVacancyHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}
	type dependencies struct {
		vacanciesUsecase   *vacancies_mock.MockIVacanciesUsecase
		fileLoadingUsecase *file_loading_mock.MockIFileLoadingUsecase

		vacancy dto.JSONVacancy

		args args
	}
	tests := []struct {
		name               string
		prepare            func(f *dependencies)
		expectedStatusCode int
		expectedBody       *dto.JSONResponse
	}{
		{
			name: "VacanciesHandler.GetVacancyHandler successfully got the vacancy",
			prepare: func(f *dependencies) {
				IDslug := uint64(1)

				f.vacancy = dto.JSONVacancy{
					EmployerID:           0,
					Salary:               1111,
					Position:             "Mock Position",
					Location:             "Mock Location",
					Description:          "Mock Description",
					WorkType:             "Mock Work Type",
					PositionCategoryName: "Mock Position Category",
					//Avatar:               "Mock Avatar",
				}

				f.vacanciesUsecase.
					EXPECT().
					GetVacancy(gomock.Any(), IDslug).
					Return(&f.vacancy, nil)
				f.fileLoadingUsecase.
					EXPECT().
					FindCompressedFile(f.vacancy.Avatar).
					Return("")

				f.args.r = httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/vacancy/%d", IDslug),
					nil,
				)
				f.args.w = httptest.NewRecorder()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "VacancyHandler.GetVacancyHandler db err",
			prepare: func(f *dependencies) {
				IDslug := uint64(1)

				f.vacanciesUsecase.
					EXPECT().
					GetVacancy(gomock.Any(), IDslug).
					Return(nil, fmt.Errorf(dto.MsgDataBaseError))

				f.args.r = httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/vacancy/%d", IDslug),
					nil,
				)
				f.args.w = httptest.NewRecorder()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "VacanciesHandler.GetVacancyHandler successfully got the vacancy",
			prepare: func(f *dependencies) {
				//IDslug := uint64(1)

				f.vacancy = dto.JSONVacancy{
					EmployerID:           0,
					Salary:               1111,
					Position:             "Mock Position",
					Location:             "Mock Location",
					Description:          "Mock Description",
					WorkType:             "Mock Work Type",
					PositionCategoryName: "Mock Position Category",
					//Avatar:               "Mock Avatar",
				}

				f.args.r = httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/vacancy/%s", "111111111111111111111111111111111111111"),
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
				vacanciesUsecase:   vacancies_mock.NewMockIVacanciesUsecase(ctrl),
				fileLoadingUsecase: file_loading_mock.NewMockIFileLoadingUsecase(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(&d)
			}

			logger := logger.NewLogrusLogger()
			logger.Out = io.Discard

			app := &internal.App{
				Logger: logger,
				Usecases: &internal.Usecases{
					VacanciesUsecase:   d.vacanciesUsecase,
					FileLoadingUsecase: d.fileLoadingUsecase,
				},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
			}

			h := delivery.NewVacanciesHandlers(app)
			r := mux.NewRouter()
			r.HandleFunc("/api/v1/vacancy/{id:[0-9]+}", h.GetVacancy).Methods(http.MethodGet)
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

func TestUpdateVacancyHandler(t *testing.T) {
	t.Parallel()
	type in struct {
		updatedVacancy *dto.JSONVacancy
		currentUser    *dto.UserFromSession
		slug           string
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
		vacanciesUsecase *vacancies_mock.MockIVacanciesUsecase
	}
	tests := []struct {
		name    string
		prepare func(in *in, out *outExpected, usecase *usecaseMock, args *args)
	}{
		{
			name: "VacancysHandler.UpdateVacancyHandler success",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.updatedVacancy = &dto.JSONVacancy{
					EmployerID:           0,
					Salary:               1111,
					Position:             "Mock Position",
					Location:             "Mock Location",
					Description:          "Mock Description",
					WorkType:             "Mock Work Type",
					PositionCategoryName: "Mock Position Category",
					//Avatar:               "Mock Avatar",
				}
				in.currentUser = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}
				in.slug = "1"

				multipartForm, contentType := createMultipartFormJSONVacancy(in.updatedVacancy, "1111")

				slugInt, _ := strconv.Atoi(in.slug)

				expectedVacancy := map[string]interface{}{
					"id":               float64(0),
					"employer":         float64(0),
					"salary":           float64(1111),
					"position":         "Mock Position",
					"location":         "Mock Location",
					"description":      "Mock Description",
					"workType":         "Mock Work Type",
					"positionGroup":    "Mock Position Category",
					"compressedAvatar": "",
					"updatedAt":        "",
					"createdAt":        "",
					"avatar":           "",
					"companyName":      "",
				}
				out.response = &dto.JSONResponse{
					HTTPStatus: http.StatusOK,
					Body:       expectedVacancy,
				}
				out.status = http.StatusOK

				usecase.vacanciesUsecase.
					EXPECT().
					UpdateVacancy(gomock.Any(), uint64(slugInt), in.updatedVacancy, in.currentUser).
					Return(in.updatedVacancy, nil)

				args.r = httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/vacancy/%s", in.slug),
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
			name: "VacancysHandler.UpdateVacancyHandler success",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.updatedVacancy = &dto.JSONVacancy{
					EmployerID:           0,
					Salary:               1111,
					Position:             "Mock Position",
					Location:             "Mock Location",
					Description:          "Mock Description",
					WorkType:             "Mock Work Type",
					PositionCategoryName: "Mock Position Category",
					//Avatar:               "Mock Avatar",
				}
				in.currentUser = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}
				in.slug = "1"

				multipartForm, contentType := createMultipartFormJSONVacancy(in.updatedVacancy, "1111.111")

				//slugInt, _ := strconv.Atoi(in.slug)

				out.response = &dto.JSONResponse{
					HTTPStatus: http.StatusBadRequest,
					Error:      dto.MsgInvalidJSON,
				}
				out.status = http.StatusBadRequest

				args.r = httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/vacancy/%s", in.slug),
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
			name: "VacancysHandler.UpdateVacancyHandler can't get user from context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"

				out.status = http.StatusUnauthorized
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgUnableToGetUserFromContext,
				}

				args.r = httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/vacancy/%s", in.slug),
					nil,
				)
				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "VacancysHandler.UpdateVacancyHandler usecase returned internal error",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.updatedVacancy = &dto.JSONVacancy{
					EmployerID:           0,
					Salary:               1111,
					Position:             "Mock Position",
					Location:             "Mock Location",
					Description:          "Mock Description",
					WorkType:             "Mock Work Type",
					PositionCategoryName: "Mock Position Category",
					//Avatar:               "Mock Avatar",
				}
				in.currentUser = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}
				in.slug = "1451566565115655515615645645656156156156156156156156156"

				multipartForm, contentType := createMultipartFormJSONVacancy(in.updatedVacancy, "1111")

				// slugInt, _ := strconv.Atoi(in.slug)

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
				}

				args.r = httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/vacancy/%s", in.slug),
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
			name: "VacancysHandler.UpdateVacancyHandler success",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.updatedVacancy = &dto.JSONVacancy{
					EmployerID:           0,
					Salary:               1111,
					Position:             "Mock Position",
					Location:             "Mock Location",
					Description:          "Mock Description",
					WorkType:             "Mock Work Type",
					PositionCategoryName: "Mock Position Category",
					//Avatar:               "Mock Avatar",
				}
				in.currentUser = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}
				in.slug = "1"

				multipartForm, contentType := createMultipartFormJSONVacancy(in.updatedVacancy, "1111")

				slugInt, _ := strconv.Atoi(in.slug)

				out.response = &dto.JSONResponse{
					HTTPStatus: http.StatusInternalServerError,
					Error:      dto.MsgDataBaseError,
				}
				out.status = http.StatusInternalServerError

				usecase.vacanciesUsecase.
					EXPECT().
					UpdateVacancy(gomock.Any(), uint64(slugInt), in.updatedVacancy, in.currentUser).
					Return(nil, fmt.Errorf(dto.MsgDataBaseError))

				args.r = httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/vacancy/%s", in.slug),
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

			usecase.vacanciesUsecase = vacancies_mock.NewMockIVacanciesUsecase(ctrl)

			tt.prepare(in, out, usecase, args)

			logger := logrus.New()
			logger.Out = io.Discard

			app := &internal.App{
				Logger: logger,
				Usecases: &internal.Usecases{
					VacanciesUsecase: usecase.vacanciesUsecase,
				},
				Repositories: &internal.Repositories{},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
			}

			h := delivery.NewVacanciesHandlers(app)
			r := mux.NewRouter()
			r.HandleFunc("/api/v1/vacancy/{id:[0-9]+}", h.UpdateVacancy)

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
func createMultipartFormJSONVacancy(jsonForm *dto.JSONVacancy, salary string) (*bytes.Buffer, string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	defer writer.Close()
	_ = writer.WriteField("salary", salary)
	_ = writer.WriteField("position", jsonForm.Position)
	_ = writer.WriteField("location", jsonForm.Location)
	_ = writer.WriteField("description", jsonForm.Description)
	_ = writer.WriteField("workType", jsonForm.WorkType)
	_ = writer.WriteField("group", jsonForm.PositionCategoryName)
	return &buf, writer.FormDataContentType()
}

func TestDeleteVacancyHandler(t *testing.T) {
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
		vacanciesUsecase *vacancies_mock.MockIVacanciesUsecase
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
			name: "VacancyHandler.DeleteVacancyHandler success",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"
				in.user = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}

				out.status = http.StatusOK
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
				}

				slugInt, _ := strconv.Atoi(in.slug)

				usecase.vacanciesUsecase.
					EXPECT().
					DeleteVacancy(gomock.Any(), uint64(slugInt), in.user).
					Return(nil)

				args.r = httptest.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/api/v1/vacancy/%s", in.slug),
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
			name: "VacancyHandler.DeleteVacancyHandler bad slug",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "11111111111111111111111111111"

				out.status = http.StatusNotFound
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
				}

				args.r = httptest.NewRequest(
					"",
					fmt.Sprintf("/api/v1/vacancy/11111111111111111111111111111"),
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
			name: "VacancyHandler.DeleteVacancyHandler no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"

				out.status = http.StatusUnauthorized
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgUnableToGetUserFromContext,
				}

				args.r = httptest.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/api/v1/vacancy/%s", in.slug),
					nil,
				)
				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "VacancyHandler.DeleteVacancyHandler no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"
				in.user = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgDataBaseError,
				}

				slugInt, _ := strconv.Atoi(in.slug)

				usecase.vacanciesUsecase.
					EXPECT().
					DeleteVacancy(gomock.Any(), uint64(slugInt), in.user).
					Return(fmt.Errorf(dto.MsgDataBaseError))

				args.r = httptest.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/api/v1/vacancy/%s", in.slug),
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

			usecase.vacanciesUsecase = vacancies_mock.NewMockIVacanciesUsecase(ctrl)

			tt.prepare(in, out, usecase, args)

			logger := logrus.New()
			logger.Out = io.Discard

			app := &internal.App{
				Logger: logger,
				Usecases: &internal.Usecases{
					VacanciesUsecase: usecase.vacanciesUsecase,
				},
				Repositories: &internal.Repositories{},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
			}

			testMux := mux.NewRouter()
			h := delivery.NewVacanciesHandlers(app)
			testMux.HandleFunc("/api/v1/vacancy/{id:[0-9]+}", h.DeleteVacancy)

			require.NotNil(t, args.r, "request is nil")
			require.NotNil(t, args.w, "response is nil")

			testMux.ServeHTTP(args.w, args.r)

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

func TestGetVacancySubscription(t *testing.T) {
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
		vacanciesUsecase *vacancies_mock.MockIVacanciesUsecase
	}
	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}
	tests := []struct {
		name    string
		prepare func(in *in, out *outExpected, usecase *usecaseMock, args *args)
	}{
		// {
		// 	name: "VacancyHandler.GetVacancySubscription success",
		// 	prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
		// 		in.slug = "1"
		// 		in.user = &dto.UserFromSession{
		// 			ID:       1,
		// 			UserType: dto.UserTypeApplicant,
		// 		}

		// 		out.status = http.StatusOK
		// 		out.response = &dto.JSONResponse{
		// 			HTTPStatus: out.status,
		// 		}

		// 		slugInt, _ := strconv.Atoi(in.slug)

		// 		usecase.vacanciesUsecase.
		// 			EXPECT().
		// 			GetSubscriptionInfo(gomock.Any(), uint64(1), uint64(slugInt)).
		// 			Return(nil, nil)

		// 		args.r = httptest.NewRequest(
		// 			http.MethodGet,
		// 			fmt.Sprintf("/api/v1/vacancy/%s/subscription", in.slug),
		// 			nil,
		// 		).WithContext(
		// 			context.WithValue(
		// 				context.Background(),
		// 				dto.UserContextKey,
		// 				in.user,
		// 			),
		// 		)
		// 		args.w = httptest.NewRecorder()
		// 	},
		// },
		{
			name: "VacancyHandler.GetVacancySubscription bad slug",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "142526575673463457814521467851672457"

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
				}

				args.r = httptest.NewRequest(
					"",
					fmt.Sprintf("/api/v1/vacancy/142526575673463457814521467851672457/subscription"),
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
			name: "VacancyHandler.GetVacancySubscription no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"

				out.status = http.StatusUnauthorized
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgUnableToGetUserFromContext,
				}

				args.r = httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/vacancy/%s/subscription", in.slug),
					nil,
				)
				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "VacancyHandler.DeleteVacancyHandler no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"
				in.user = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgDataBaseError,
				}

				slugInt, _ := strconv.Atoi(in.slug)

				usecase.vacanciesUsecase.
					EXPECT().
					GetSubscriptionInfo(gomock.Any(), uint64(slugInt), uint64(slugInt)).
					Return(nil, fmt.Errorf(dto.MsgDataBaseError))

				args.r = httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/vacancy/%s/subscription", in.slug),
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

			usecase.vacanciesUsecase = vacancies_mock.NewMockIVacanciesUsecase(ctrl)

			tt.prepare(in, out, usecase, args)

			logger := logrus.New()
			logger.Out = io.Discard

			app := &internal.App{
				Logger: logger,
				Usecases: &internal.Usecases{
					VacanciesUsecase: usecase.vacanciesUsecase,
				},
				Repositories: &internal.Repositories{},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
			}

			testMux := mux.NewRouter()
			h := delivery.NewVacanciesHandlers(app)
			testMux.HandleFunc("/api/v1/vacancy/{id:[0-9]+}/subscription", h.GetVacancySubscription)

			require.NotNil(t, args.r, "request is nil")
			require.NotNil(t, args.w, "response is nil")

			testMux.ServeHTTP(args.w, args.r)

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

func TestSubscribeVacancy(t *testing.T) {
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
		vacanciesUsecase  *vacancies_mock.MockIVacanciesUsecase
		applicantUsecase  *applicant_mock.MockIApplicantUsecase
		notificationsGRPC *grpc_mock.MockNotificationsServiceClient
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
			name: "VacancyHandler.GetVacancySubscription success",
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

				usecase.vacanciesUsecase.
					EXPECT().
					SubscribeOnVacancy(gomock.Any(), uint64(slugInt), in.user).
					Return(nil)
					// go func() {
				vacancyJSON := &dto.JSONVacancy{
					ID:         1,
					EmployerID: 1,
					Position:   "Position",
				}
				usecase.vacanciesUsecase.
					EXPECT().
					GetVacancy(gomock.Any(), uint64(slugInt)).
					Return(vacancyJSON, nil)
				applicant := &dto.JSONGetApplicantProfile{
					ID:        in.user.ID,
					FirstName: "John",
					LastName:  "Doe",
				}
				usecase.applicantUsecase.
					EXPECT().
					GetApplicantProfile(context.Background(), in.user.ID).
					Return(applicant, nil)
				input := &notifications_grpc.CreateEmployerNotificationInput{
					ApplicantID:   in.user.ID,
					VacancyID:     uint64(slugInt),
					EmployerID:    vacancyJSON.EmployerID,
					ApplicantInfo: applicant.FirstName + " " + applicant.LastName,
					VacancyInfo:   vacancyJSON.Position,
				}
				usecase.notificationsGRPC.
					EXPECT().
					CreateEmployerNotification(gomock.Any(), input).
					Return(&notifications_grpc.Nothing{}, nil)
				// }()

				args.r = httptest.NewRequest(
					http.MethodPost,
					fmt.Sprintf("/api/v1/vacancy/%s/subscription", in.slug),
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
			name: "VacancyHandler.GetVacancySubscription bad slug",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "142526575673463457814521467851672457"

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
				}

				args.r = httptest.NewRequest(
					"",
					fmt.Sprintf("/api/v1/vacancy/142526575673463457814521467851672457/subscription"),
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
			name: "VacancyHandler.GetVacancySubscription no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"

				out.status = http.StatusUnauthorized
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgUnableToGetUserFromContext,
				}

				args.r = httptest.NewRequest(
					http.MethodPost,
					fmt.Sprintf("/api/v1/vacancy/%s/subscription", in.slug),
					nil,
				)
				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "VacancyHandler.DeleteVacancyHandler no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"
				in.user = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgDataBaseError,
				}

				slugInt, _ := strconv.Atoi(in.slug)

				usecase.vacanciesUsecase.
					EXPECT().
					SubscribeOnVacancy(gomock.Any(), uint64(slugInt), in.user).
					Return(fmt.Errorf(dto.MsgDataBaseError))

				args.r = httptest.NewRequest(
					http.MethodPost,
					fmt.Sprintf("/api/v1/vacancy/%s/subscription", in.slug),
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

			usecase.vacanciesUsecase = vacancies_mock.NewMockIVacanciesUsecase(ctrl)
			usecase.applicantUsecase = applicant_mock.NewMockIApplicantUsecase(ctrl)
			//usecase.notificationsGRPC = notifications_mock.NewMockINotificationsUsecase(ctrl)
			usecase.notificationsGRPC = grpc_mock.NewMockNotificationsServiceClient(ctrl)

			tt.prepare(in, out, usecase, args)

			logger := logrus.New()
			logger.Out = io.Discard

			app := &internal.App{
				Logger: logger,
				Usecases: &internal.Usecases{
					VacanciesUsecase: usecase.vacanciesUsecase,
					ApplicantUsecase: usecase.applicantUsecase,
				},
				Repositories: &internal.Repositories{},
				Microservices: &internal.Microservices{
					Compress:      nil,
					Notifications: usecase.notificationsGRPC,
				},
			}

			testMux := mux.NewRouter()
			h := delivery.NewVacanciesHandlers(app)
			testMux.HandleFunc("/api/v1/vacancy/{id:[0-9]+}/subscription", h.SubscribeVacancy)

			require.NotNil(t, args.r, "request is nil")
			require.NotNil(t, args.w, "response is nil")

			testMux.ServeHTTP(args.w, args.r)

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

func TestUnsubscribeVacancy(t *testing.T) {
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
		vacanciesUsecase *vacancies_mock.MockIVacanciesUsecase
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
			name: "VacancyHandler.GetVacancySubscription success",
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

				usecase.vacanciesUsecase.
					EXPECT().
					UnsubscribeFromVacancy(gomock.Any(), uint64(slugInt), in.user).
					Return(nil)

				args.r = httptest.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/api/v1/vacancy/%s/subscription", in.slug),
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
			name: "VacancyHandler.GetVacancySubscription bad slug",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "142526575673463457814521467851672457"

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
				}

				args.r = httptest.NewRequest(
					"",
					fmt.Sprintf("/api/v1/vacancy/142526575673463457814521467851672457/subscription"),
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
			name: "VacancyHandler.GetVacancySubscription no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"

				out.status = http.StatusUnauthorized
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgUnableToGetUserFromContext,
				}

				args.r = httptest.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/api/v1/vacancy/%s/subscription", in.slug),
					nil,
				)
				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "VacancyHandler.DeleteVacancyHandler no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"
				in.user = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgDataBaseError,
				}

				slugInt, _ := strconv.Atoi(in.slug)

				usecase.vacanciesUsecase.
					EXPECT().
					UnsubscribeFromVacancy(gomock.Any(), uint64(slugInt), in.user).
					Return(fmt.Errorf(dto.MsgDataBaseError))

				args.r = httptest.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/api/v1/vacancy/%s/subscription", in.slug),
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

			usecase.vacanciesUsecase = vacancies_mock.NewMockIVacanciesUsecase(ctrl)

			tt.prepare(in, out, usecase, args)

			logger := logrus.New()
			logger.Out = io.Discard

			app := &internal.App{
				Logger: logger,
				Usecases: &internal.Usecases{
					VacanciesUsecase: usecase.vacanciesUsecase,
				},
				Repositories: &internal.Repositories{},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
			}

			testMux := mux.NewRouter()
			h := delivery.NewVacanciesHandlers(app)
			testMux.HandleFunc("/api/v1/vacancy/{id:[0-9]+}/subscription", h.UnsubscribeVacancy)

			require.NotNil(t, args.r, "request is nil")
			require.NotNil(t, args.w, "response is nil")

			testMux.ServeHTTP(args.w, args.r)

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

func TestGetVacancySubscribers(t *testing.T) {
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
		vacanciesUsecase   *vacancies_mock.MockIVacanciesUsecase
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
		// {
		// 	name: "VacancyHandler.GetVacancySubscription success",
		// 	prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
		// 		in.slug = "1"
		// 		in.user = &dto.UserFromSession{
		// 			ID:       1,
		// 			UserType: dto.UserTypeApplicant,
		// 		}

		// 		out.status = http.StatusOK
		// 		out.response = &dto.JSONResponse{
		// 			HTTPStatus: out.status,
		// 			Body:       map[string]interface {}{"subscribers":interface {}(nil), "vacancyID":float64(0)},
		// 		}

		// 		slugInt, _ := strconv.Atoi(in.slug)

		// 		subs := &dto.JSONVacancySubscribers{

		// 		}
		// 		usecase.vacanciesUsecase.
		// 			EXPECT().
		// 			GetVacancySubscribers(gomock.Any(), uint64(slugInt), in.user).
		// 			Return(subs, nil)

		// 		args.r = httptest.NewRequest(
		// 			http.MethodGet,
		// 			fmt.Sprintf("/api/v1/vacancy/%s/subscribers", in.slug),
		// 			nil,
		// 		).WithContext(
		// 			context.WithValue(
		// 				context.Background(),
		// 				dto.UserContextKey,
		// 				in.user,
		// 			),
		// 		)
		// 		args.w = httptest.NewRecorder()
		// 	},
		// },
		{
			name: "VacancyHandler.GetVacancySubscription bad slug",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "142526575673463457814521467851672457"

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
				}

				args.r = httptest.NewRequest(
					"",
					fmt.Sprintf("/api/v1/vacancy/142526575673463457814521467851672457/subscribers"),
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
			name: "VacancyHandler.GetVacancySubscription no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"

				out.status = http.StatusUnauthorized
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgUnableToGetUserFromContext,
				}

				args.r = httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/vacancy/%s/subscribers", in.slug),
					nil,
				)
				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "VacancyHandler.DeleteVacancyHandler no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"
				in.user = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgDataBaseError,
				}

				slugInt, _ := strconv.Atoi(in.slug)

				usecase.vacanciesUsecase.
					EXPECT().
					GetVacancySubscribers(gomock.Any(), uint64(slugInt), in.user).
					Return(nil, fmt.Errorf(dto.MsgDataBaseError))

				args.r = httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/vacancy/%s/subscribers", in.slug),
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

			usecase.vacanciesUsecase = vacancies_mock.NewMockIVacanciesUsecase(ctrl)
			usecase.fileLoadingUsecase = file_loading_mock.NewMockIFileLoadingUsecase(ctrl)

			tt.prepare(in, out, usecase, args)

			logger := logrus.New()
			logger.Out = io.Discard

			app := &internal.App{
				Logger: logger,
				Usecases: &internal.Usecases{
					VacanciesUsecase:   usecase.vacanciesUsecase,
					FileLoadingUsecase: usecase.fileLoadingUsecase,
				},
				Repositories: &internal.Repositories{},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
			}

			testMux := mux.NewRouter()
			h := delivery.NewVacanciesHandlers(app)
			testMux.HandleFunc("/api/v1/vacancy/{id:[0-9]+}/subscribers", h.GetVacancySubscribers)

			require.NotNil(t, args.r, "request is nil")
			require.NotNil(t, args.w, "response is nil")

			testMux.ServeHTTP(args.w, args.r)

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

func TestAddVacancyIntoFavorite(t *testing.T) {
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
		vacanciesUsecase *vacancies_mock.MockIVacanciesUsecase
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
			name: "VacancyHandler.GetVacancyFavorite success",
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

				usecase.vacanciesUsecase.
					EXPECT().
					AddIntoFavorite(gomock.Any(), uint64(slugInt), in.user).
					Return(nil)

				args.r = httptest.NewRequest(
					http.MethodPost,
					fmt.Sprintf("/api/v1/applicant/%s/favorite-vacancy", in.slug),
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
			name: "VacancyHandler.GetVacancyFavorite bad slug",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "142526575673463457814521467851672457"

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
				}

				args.r = httptest.NewRequest(
					"",
					fmt.Sprintf("/api/v1/applicant/142526575673463457814521467851672457/favorite-vacancy"),
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
			name: "VacancyHandler.GetVacancyFavorite no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"

				out.status = http.StatusUnauthorized
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgUnableToGetUserFromContext,
				}

				args.r = httptest.NewRequest(
					http.MethodPost,
					fmt.Sprintf("/api/v1/applicant/%s/favorite-vacancy", in.slug),
					nil,
				)
				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "VacancyHandler.DeleteVacancyHandler no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"
				in.user = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgDataBaseError,
				}

				slugInt, _ := strconv.Atoi(in.slug)

				usecase.vacanciesUsecase.
					EXPECT().
					AddIntoFavorite(gomock.Any(), uint64(slugInt), in.user).
					Return(fmt.Errorf(dto.MsgDataBaseError))

				args.r = httptest.NewRequest(
					http.MethodPost,
					fmt.Sprintf("/api/v1/applicant/%s/favorite-vacancy", in.slug),
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

			usecase.vacanciesUsecase = vacancies_mock.NewMockIVacanciesUsecase(ctrl)

			tt.prepare(in, out, usecase, args)

			logger := logrus.New()
			logger.Out = io.Discard

			app := &internal.App{
				Logger: logger,
				Usecases: &internal.Usecases{
					VacanciesUsecase: usecase.vacanciesUsecase,
				},
				Repositories: &internal.Repositories{},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
			}

			testMux := mux.NewRouter()
			h := delivery.NewVacanciesHandlers(app)
			testMux.HandleFunc("/api/v1/applicant/{id:[0-9]+}/favorite-vacancy", h.AddVacancyIntoFavorite)

			require.NotNil(t, args.r, "request is nil")
			require.NotNil(t, args.w, "response is nil")

			testMux.ServeHTTP(args.w, args.r)

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
func TestDellVacancyFromFavorite(t *testing.T) {
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
		vacanciesUsecase *vacancies_mock.MockIVacanciesUsecase
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
			name: "VacancyHandler.GetVacancyFavorite success",
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

				usecase.vacanciesUsecase.
					EXPECT().
					Unfavorite(gomock.Any(), uint64(slugInt), in.user).
					Return(nil)

				args.r = httptest.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/api/v1/applicant/%s/favorite-vacancy", in.slug),
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
			name: "VacancyHandler.GetVacancyFavorite bad slug",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "142526575673463457814521467851672457"

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
				}

				args.r = httptest.NewRequest(
					"",
					fmt.Sprintf("/api/v1/applicant/142526575673463457814521467851672457/favorite-vacancy"),
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
			name: "VacancyHandler.GetVacancyFavorite no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"

				out.status = http.StatusUnauthorized
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgUnableToGetUserFromContext,
				}

				args.r = httptest.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/api/v1/applicant/%s/favorite-vacancy", in.slug),
					nil,
				)
				args.w = httptest.NewRecorder()
			},
		},
		{
			name: "VacancyHandler.DeleteVacancyHandler no user in context",
			prepare: func(in *in, out *outExpected, usecase *usecaseMock, args *args) {
				in.slug = "1"
				in.user = &dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeEmployer,
				}

				out.status = http.StatusInternalServerError
				out.response = &dto.JSONResponse{
					HTTPStatus: out.status,
					Error:      dto.MsgDataBaseError,
				}

				slugInt, _ := strconv.Atoi(in.slug)

				usecase.vacanciesUsecase.
					EXPECT().
					Unfavorite(gomock.Any(), uint64(slugInt), in.user).
					Return(fmt.Errorf(dto.MsgDataBaseError))

				args.r = httptest.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/api/v1/applicant/%s/favorite-vacancy", in.slug),
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

			usecase.vacanciesUsecase = vacancies_mock.NewMockIVacanciesUsecase(ctrl)

			tt.prepare(in, out, usecase, args)

			logger := logrus.New()
			logger.Out = io.Discard

			app := &internal.App{
				Logger: logger,
				Usecases: &internal.Usecases{
					VacanciesUsecase: usecase.vacanciesUsecase,
				},
				Repositories: &internal.Repositories{},
				Microservices: &internal.Microservices{
					Compress: nil,
				},
			}

			testMux := mux.NewRouter()
			h := delivery.NewVacanciesHandlers(app)
			testMux.HandleFunc("/api/v1/applicant/{id:[0-9]+}/favorite-vacancy", h.DellVacancyFromFavorite)

			require.NotNil(t, args.r, "request is nil")
			require.NotNil(t, args.w, "response is nil")

			testMux.ServeHTTP(args.w, args.r)

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
