package usecase

import (
	"fmt"
	"mime/multipart"
	"slices"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	fileloading "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	compressmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"
	"github.com/sirupsen/logrus"
)

type FileLoadingUsecase struct {
	logger                *logrus.Logger
	FileLoadingRepository fileloading.IFileLoadingRepository
	CompressGRPC          compressmicroservice.CompressServiceClient
	conf                  *configs.Config
}

func NewFileLoadingUsecase(logger *logrus.Logger, repositories *internal.Repositories, microservices *internal.Microservices, conf *configs.Config) *FileLoadingUsecase {
	return &FileLoadingUsecase{
		logger:                logger,
		FileLoadingRepository: repositories.FileLoadingRepository,
		CompressGRPC:          microservices.Compress,
		conf:                  conf,
	}
}

var allowedTypes = []string{"image/jpeg", "image/jpg", "image/svg", "image/svg+xml"}

func (vu *FileLoadingUsecase) WriteImage(file multipart.File, header *multipart.FileHeader) (string, string, error) {
	a := header.Header
	vu.logger.Debug(a["Content-Type"][0])
	for _, i := range a["Content-Type"] {
		if !slices.Contains(allowedTypes, i) {
			return "", "", fmt.Errorf(dto.MsgInvalidFile)
		}
	}
	filename := utils.GenerateSessionToken(utils.TokenLength+10, dto.UserTypeApplicant)
	dir, fileAddress, err := vu.FileLoadingRepository.WriteFileOnDisk(filename, header, file)
	if err != nil {
		return "", "", err
	}

	return dir + fileAddress, vu.conf.CompressMicroservice.CompressedMediaDir + fileAddress, nil
}
