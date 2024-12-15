package usecase

import (
	"context"
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
	var buff []byte
	file.Read(buff)
	vu.logger.Debugf("Start compression")
	_, err = vu.CompressGRPC.CompressAndSaveFile(
		context.Background(),
		&compressmicroservice.CompressAndSaveFileInput{
			FileName: filename + header.Filename,
			FileType: header.Header["Content-Type"][0],
			File:     buff,
		},
	)
	if err != nil {
		return "", "", err
	}
	return dir + fileAddress, vu.conf.CompressMicroservice.CompressedMediaDir + fileAddress, nil
}

func (vu *FileLoadingUsecase) CVtoPDF(CV *dto.JSONCv, applicant *dto.JSONGetApplicantProfile) (*dto.CVPDFFile, error) {
	name, err := vu.FileLoadingRepository.CVtoPDF(CV, applicant)
	if err != nil {
		return nil, err
	}
	return &dto.CVPDFFile{FileName: name}, nil
}
