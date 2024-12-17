package usecase_test

import (
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/usecase"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAlEmployerNotifications(t *testing.T) {
	t.Parallel()
	type repo struct {
		notificationRepo *mock.MockINotificationsRepository
	}
	tests := []struct {
		name    string
		profile []*dto.EmployerNotification
		prepare func(
			repo *repo,
			profile []*dto.EmployerNotification,
		) (uint64, []*dto.EmployerNotification)
	}{
		{
			name:    "Create: bad repository",
			profile: nil,
			prepare: func(
				repo *repo,
				profile []*dto.EmployerNotification) (uint64, []*dto.EmployerNotification) {
				userID := uint64(1)
				repo.notificationRepo.
					EXPECT().
					GetAlEmployerNotifications(userID).
					Return(nil, errors.New("bad repository"))
				return userID, profile
			},
		},
		{
			name:    "Create: ok",
			profile: make([]*dto.EmployerNotification, 0),
			prepare: func(
				repo *repo,
				profile []*dto.EmployerNotification) (uint64, []*dto.EmployerNotification) {
				userID := uint64(1)
				model := []*dto.EmployerNotification{
					&dto.EmployerNotification{
						ID: 1,
					},
				}

				repo.notificationRepo.
					EXPECT().
					GetAlEmployerNotifications(userID).
					Return(model, nil)
				rprofile := []*dto.EmployerNotification{
					model[0],
				}
				return userID, rprofile
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			notificationRepo: mock.NewMockINotificationsRepository(ctrl),
		}
		var userID uint64
		userID, tt.profile = tt.prepare(repo, tt.profile)

		uc := usecase.NewNotificationsUsecase(repo.notificationRepo, logrus.New())
		profile, _ := uc.GetAlEmployerNotifications(userID)

		require.Equal(t, tt.profile, profile)
	}
}

func TestMakeEmployerNotificationRead(t *testing.T) {
	t.Parallel()
	type repo struct {
		notificationRepo *mock.MockINotificationsRepository
	}
	tests := []struct {
		name    string
		err     error
		prepare func(
			repo *repo,
			errs error,
		) (uint64, error)
	}{
		{
			name: "Create: ok",
			err:  nil,
			prepare: func(
				repo *repo,
				err error) (uint64, error) {
				userID := uint64(1)

				repo.notificationRepo.
					EXPECT().
					MakeEmployerNotificationRead(userID).
					Return(nil)
				return userID, err
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			notificationRepo: mock.NewMockINotificationsRepository(ctrl),
		}
		var userID uint64
		userID, _ = tt.prepare(repo, tt.err)

		uc := usecase.NewNotificationsUsecase(repo.notificationRepo, logrus.New())
		profile := uc.MakeEmployerNotificationRead(userID)
		require.Equal(t, tt.err, profile)
	}
}

func TestCreateEmployerNotification(t *testing.T) {
	t.Parallel()
	type repo struct {
		notificationRepo *mock.MockINotificationsRepository
	}
	tests := []struct {
		name          string
		applicantID   uint64
		employerID    uint64
		vacancyID     uint64
		applicantInfo string
		vacancyInfo   string
		err           error
		prepare       func(
			repo *repo,
			applicantID uint64,
			employerID uint64,
			vacancyID uint64,
			applicantInfo string,
			vacancyInfo string,
		) (uint64, uint64, uint64, string, string)
	}{
		{
			name: "Create: ok",
			err:  nil,
			prepare: func(
				repo *repo,
				applicantID uint64,
			employerID uint64,
			vacancyID uint64,
			applicantInfo string,
			vacancyInfo string) (uint64, uint64, uint64, string, string) {
				repo.notificationRepo.
					EXPECT().
					CreateEmployerNotification(applicantID, employerID, vacancyID, gomock.Any()).
					Return(nil)
				return applicantID, employerID, vacancyID, applicantInfo, vacancyInfo
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			notificationRepo: mock.NewMockINotificationsRepository(ctrl),
		}
		applicantID, employerID, vacancyID, applicantInfo, vacancyInfo := tt.prepare(repo, tt.applicantID, tt.employerID, tt.vacancyID, tt.applicantInfo, tt.vacancyInfo)

		uc := usecase.NewNotificationsUsecase(repo.notificationRepo, logrus.New())
		profile := uc.CreateEmployerNotification(applicantID, employerID, vacancyID, applicantInfo, vacancyInfo)
		require.Equal(t, tt.err, profile)
	}
}
