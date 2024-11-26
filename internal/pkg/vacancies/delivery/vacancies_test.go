package delivery_test

import (
	//"bytes"

	"context"
	"encoding/json"
	"fmt"
	"strconv"

	//"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/delivery"
	vacancies_mock "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/mock"
	"github.com/gorilla/mux"

	//"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	//"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

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
					DeleteVacancy(uint64(slugInt), in.user).
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
					DeleteVacancy(uint64(slugInt), in.user).
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
				Repositories: &internal.Repositories{
					SessionApplicantRepository: nil,
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
			err := json.NewDecoder(args.w.Result().Body).Decode(jsonResonse)
			require.NoError(t, err)

			require.Equalf(t, out.response, jsonResonse,
				"got response %v, expected %v",
				jsonResonse,
				out.response,
			)
		})
	}
}
