package usecase

import (
	"slices"

	compressinterfaces "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress"
	"github.com/sirupsen/logrus"
)

type CompressUsecase struct {
	compressRepo compressinterfaces.ICompressRepository
	logger       *logrus.Entry
}

func NewCompressUsecase(compressRepo compressinterfaces.ICompressRepository, logger *logrus.Logger) *CompressUsecase {
	return &CompressUsecase{
		compressRepo: compressRepo,
		logger:       &logrus.Entry{Logger: logger},
	}
}

func (cu *CompressUsecase) CompressAndSaveFile(filename string, fileType string, file []byte) error {
	funcName := "CompressUsecase.CompressAndSaveFile"
	cu.logger.Debugf("%s: got request: %s %s %s", funcName, filename, fileType, file)
	if !slices.Contains(compressinterfaces.AllowedTypes, fileType) {
		return compressinterfaces.NotAllowedType
	}
	err := cu.compressRepo.SaveFile(filename, fileType, file)
	return err
}

func (cu *CompressUsecase) DeleteFile(filename string) error {
	funcName := "CompressUsecase.DeleteFile"
	cu.logger.Debugf("%s: got request: %s", funcName, filename)
	err := cu.compressRepo.DeleteFile(filename)
	return err
}
