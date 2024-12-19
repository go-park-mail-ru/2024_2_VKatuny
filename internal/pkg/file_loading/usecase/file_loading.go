package usecase

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"strings"

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

var allowedTypes = []string{"image/jpeg", "image/svg+xml", "image/pjpeg", "image/webp"}

func (vu *FileLoadingUsecase) WriteImage(file []byte, header *multipart.FileHeader) (string, string, error) {
	ContentType := http.DetectContentType(file)
	vu.logger.Debug(ContentType)
	if !slices.Contains(allowedTypes, ContentType) || header.Size > 25<<21 || len(file) > 25<<21 {
		return "", "", fmt.Errorf(dto.MsgInvalidFile)
	}
	filename := utils.GenerateSessionToken(utils.TokenLength+10, dto.UserTypeApplicant)
	dir, fileName, err := vu.FileLoadingRepository.WriteFileOnDisk(filename, header, file)
	if err != nil {
		return "", "", err
	}
	return dir + fileName, vu.FindCompressedFile(fileName), nil
}

func (vu *FileLoadingUsecase) FindCompressedFile(filename string) string {
	filename = strings.Split(filename, "/")[len(strings.Split(filename, "/"))-1]
	vu.logger.Debugf("filename: %s", filename)
	dir := vu.conf.CompressMicroservice.CompressedMediaDir
	compressed, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	dirList := strings.Split(dir, "/")
	dirList = dirList[slices.Index(dirList, "2024_2_VKatuny")+1:]
	dirCut := strings.Join(dirList, "/") + "/"
	for _, file := range compressed {
		if file.Name()[:strings.Index(file.Name(), ".")] == filename[:strings.Index(filename, ".")] {
			return dirCut + file.Name()
		}
	}
	return ""
}

func (vu *FileLoadingUsecase) CVtoPDF(CV *dto.JSONCv, applicant *dto.JSONGetApplicantProfile) (*dto.CVPDFFile, error) {
	name, err := vu.FileLoadingRepository.CVtoPDF(CV, applicant)
	if err != nil {
		return nil, err
	}
	return &dto.CVPDFFile{FileName: name}, nil
}
