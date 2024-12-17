package usecase_test

import (
	"os"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading/mock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading/usecase"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFindCompressedFile(t *testing.T) {
	t.Parallel()
	type repo struct {
		fileLoadingRepo *mock.MockIFileLoadingRepository
	}
	tests := []struct {
		name    string
		path    string
		prepare func(
			repo *repo,
			path string,
		) string
	}{
		{
			name: "Create: ok",
			path: "favicon.ico",
			prepare: func(
				repo *repo,
				path string) string {
				return "/media/Compressed/favicon.ico"
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := &repo{
			fileLoadingRepo: mock.NewMockIFileLoadingRepository(ctrl),
		}
		//var userID uint64
		tt.path = tt.prepare(repo, tt.path)

		repositories := &internal.Repositories{
			FileLoadingRepository: repo.fileLoadingRepo,
		}
		microservices := &internal.Microservices{}
		pwd, _ := os.Getwd()
		newPwd := ""
		for _, i := range strings.Split(pwd, "/") {
			newPwd += i + "/"
			if i == "2024_2_VKatuny" {
				break
			}
		}
		newPwd += "configs/conf.yml"
		conf := configs.ReadConfig(newPwd)
		uc := usecase.NewFileLoadingUsecase(logrus.New(), repositories, microservices, conf)

		st1 := uc.FindCompressedFile("favicon.ico")

		require.Equal(t, tt.path, st1)
	}
}
