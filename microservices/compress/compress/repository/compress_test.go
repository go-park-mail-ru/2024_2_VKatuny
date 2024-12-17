package repository_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress/repository"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/mock"
	"github.com/stretchr/testify/require"
)

func TestScanDir(t *testing.T) {
	t.Parallel()
	type repo struct {
		notificationRepo *mock.MockINotificationsRepository
	}
	tests := []struct {
		name    string
		prepare func() error
	}{
		{
			name: "Create: ok",
			prepare: func() error {
				// userID := uint64(1)
				// model := []*dto.EmployerNotification{
				// 	&dto.EmployerNotification{
				// 		ID: 1,
				// 	},
				// }

				// // repo.notificationRepo.
				// // 	EXPECT().
				// // 	GetAlEmployerNotifications(userID).
				// // 	Return(model, nil)
				// rprofile := []*dto.EmployerNotification{
				// 	model[0],
				// }
				return nil
			},
		},
	}

	for _, tt := range tests {
		//ctrl := gomock.NewController(t)
		//defer ctrl.Finish()

		//repo := &repo{
		//	notificationRepo: mock.NewMockINotificationsRepository(ctrl),
		//}
		logger := logger.NewLogrusLogger()
		var err error
		err = tt.prepare()
		pwd, err := os.Getwd()
		pwdList := strings.Split(pwd, "/")
		mainDir := ""
		for _, i := range pwdList {
			if i == "microservices" {
				break
			}
			mainDir += i + "/"
		}
		compressed := mainDir + "media/Compressed"
		uncompressed := mainDir + "media/Uncompressed"
		
		fmt.Println(compressed, uncompressed)
		fu := repository.NewCompressRepository(compressed, uncompressed, logger)
		errfu := fu.ScanDir()

		require.Equal(t, err, errfu)
	}
}
