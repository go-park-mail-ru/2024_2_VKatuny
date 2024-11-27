package compressmicroservice

import (
	"context"

	compressinterfaces "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress"
	compress "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"
	"github.com/sirupsen/logrus"
)

type CompressManager struct {
	compress.UnsafeCompressServiceServer
	compressUsecase compressinterfaces.ICompressUsecase
	logger          *logrus.Entry
}

func NewCompressManager(compressUsecase compressinterfaces.ICompressUsecase, logger *logrus.Logger) *CompressManager {
	return &CompressManager{
		compressUsecase: compressUsecase,
		logger:          &logrus.Entry{Logger: logger},
	}
}

func (cm *CompressManager) CompressAndSaveFile(ctx context.Context, in *compress.CompressAndSaveFileInput) (*compress.Nothing, error) {
	funcName := "CompressDelivery.CompressAndSaveFile"
	cm.logger.Debugf("%s: got request: %s", funcName, in)
	if in == nil {
		return &compress.Nothing{}, compressinterfaces.WrongData
	}
	err := cm.compressUsecase.CompressAndSaveFile(in.FileName, in.FileType, in.File)
	return &compress.Nothing{}, err
}

func (cm *CompressManager) DeleteFile(ctx context.Context, in *compress.DeleteFileInput) (*compress.Nothing, error) {
	funcName := "CompressDelivery.DeleteFile"
	cm.logger.Debugf("%s: got request: %s", funcName, in)
	if in == nil {
		return &compress.Nothing{}, compressinterfaces.WrongData
	}
	err := cm.compressUsecase.DeleteFile(in.FileName)
	return &compress.Nothing{}, err
}
