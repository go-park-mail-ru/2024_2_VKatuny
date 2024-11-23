package compressmicroservice

import (
	"fmt"

	"context"

	compressinterfaces "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress"
	compress "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"
)

type CompressManager struct {
	//compressedDir string
	compress.UnsafeCompressServiceServer
	compressUsecase compressinterfaces.ICompressUsecase
}

func NewCompressManager(compressUsecase compressinterfaces.ICompressUsecase) *CompressManager {
	return &CompressManager{
		//compressedDir: "media/compressed/",
	}
}

func (cm *CompressManager) CompressAndSaveFile(ctx context.Context, in *compress.CompressAndSaveFileInput) (*compress.Nothing, error) {
	funcName := "CompressService.CompressAndSaveFile"
	fmt.Println(funcName)
	err := cm.compressUsecase.CompressAndSaveFile(in.FileName)
	return &compress.Nothing{}, err
}

func (cm *CompressManager) DeleteFile(ctx context.Context, in *compress.DeleteFileInput) (*compress.Nothing, error) {
	funcName := "CompressService.DeleteFile"
	fmt.Println(funcName)
	err := cm.compressUsecase.DeleteFile(in.FileName)
	return &compress.Nothing{}, err
}
