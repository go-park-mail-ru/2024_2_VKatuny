package delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/delivery"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestCreateCVHandler(t *testing.T) {
	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}
	type dependencies struct {
		cvsUsecase *mock.MockICVsUsecase
		logger     *logrus.Logger

		cv          dto.JSONCv
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
			name: "CVsHandler.CreateCVHandler successful generation of cv",
			prepare: func(f *dependencies) {
				f.currentUser = dto.UserFromSession{
					ID:       1,
					UserType: dto.UserTypeApplicant,
				}
				f.cv = dto.JSONCv{
					ID:                  1,
					ApplicantID:         1,
					PositionRu:          "Мок Должность",
					PositionEn:          "Mock Position",
					Description:         "Mock Description",
					JobSearchStatusName: "Ищу",
					WorkingExperience:   "1 год",
					Avatar:              "Mock Avatar",
				}

				f.logger = logrus.New()
				f.logger.Out = io.Discard

				f.cvsUsecase.
					EXPECT().
					CreateCV(&f.cv, &f.currentUser).
					Return(&f.cv, nil)

				body, _ := json.Marshal(f.cv)

				f.args.r = httptest.NewRequest(
					http.MethodPost,
					"/api/v1/cv/",
					bytes.NewReader(body),
				).WithContext(
					context.WithValue(
						context.Background(),
						dto.UserContextKey,
						&f.currentUser,
					),
				)

				f.args.w = httptest.NewRecorder()
			},
			expectedStatusCode: http.StatusOK,
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
				Repositories: &internal.Repositories{
					SessionApplicantRepository: nil,
				},
				Logger: d.logger,
			}

			testMux := http.NewServeMux()
			h := delivery.NewCVsHandler(app)

			testMux.HandleFunc("/api/v1/cv/", h.CreateCVHandler)

			testMux.ServeHTTP(d.args.w, d.args.r)

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

func TestGetCVHandler(t *testing.T) {
	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}
	type dependencies struct {
		cvsUsecase *mock.MockICVsUsecase
		logger     *logrus.Logger

		cv          dto.JSONCv

		args args
	}
	tests := []struct {
		name               string
		prepare            func(f *dependencies)
		wantErr            bool
		expectedStatusCode int
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

				// disable logger
				f.logger = logrus.New()
				f.logger.Out = io.Discard

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
				Logger: d.logger,
				Usecases: &internal.Usecases{
					CVUsecase: d.cvsUsecase,
				},
				Repositories: &internal.Repositories{
					SessionApplicantRepository: nil,
				},
			}

			testMux := http.NewServeMux()
			h := delivery.NewCVsHandler(app)

			testMux.HandleFunc("/api/v1/cv/", h.GetCVsHandler)

			testMux.ServeHTTP(d.args.w, d.args.r)
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
