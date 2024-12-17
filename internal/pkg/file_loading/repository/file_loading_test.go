package repository_test

import (
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading/repository"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestCVtoPDF(t *testing.T) {
	t.Parallel()
	type repo struct {
		fileLoadingRepo *mock.MockIFileLoadingRepository
	}
	tests := []struct {
		name      string
		cv        *dto.JSONCv
		applicant *dto.JSONGetApplicantProfile
		prepare   func(cv *dto.JSONCv, applicant *dto.JSONGetApplicantProfile) (*dto.JSONCv, *dto.JSONGetApplicantProfile)
	}{
		{
			name: "Create: ok",
			cv: &dto.JSONCv{
				ID: 1,
			},
			applicant: &dto.JSONGetApplicantProfile{
				ID: 1,
			},
			prepare: func(cv *dto.JSONCv, applicant *dto.JSONGetApplicantProfile) (*dto.JSONCv, *dto.JSONGetApplicantProfile) {

				return cv, applicant

			},
		},
	}

	for _, tt := range tests {
		//var userID uint64
		tt.cv, tt.applicant = tt.prepare(tt.cv, tt.applicant)
		uc := repository.NewFileLoadingStorage(logrus.New(), "media/Uncompressed/", "media/CVinPDF/", "templates/")

		st1, _ := uc.CVtoPDF(tt.cv, tt.applicant)

		require.Equal(t, "media/CVinPDF/1&&0.pdf", st1)
	}
}
