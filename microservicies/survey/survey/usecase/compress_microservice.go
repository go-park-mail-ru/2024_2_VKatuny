package usecase

import (
	compressinterfaces "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress"
)

type CompressUsecase struct {
	compressRepo compressinterfaces.ICompressRepository
}

func NewCompressUsecase(compressRepo compressinterfaces.ICompressRepository) *CompressUsecase {
	return &CompressUsecase{
		compressRepo: compressRepo,
	}
}

func (cu *CompressUsecase) CompressAndSaveFile(filename string) error {
	err := cu.compressRepo.SaveFile(filename)
	return err
}

func (cu *CompressUsecase) DeleteFile(filename string) error {
	err := cu.compressRepo.DeleteFile(filename)
	return err
}
