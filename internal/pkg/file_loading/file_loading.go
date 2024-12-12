package fileloading

import (
	"mime/multipart"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

type IFileLoadingRepository interface {
	WriteFileOnDisk(filename string, header *multipart.FileHeader, file multipart.File) (string, string, error)
	CVtoPDF(cvID uint64) (string, error)
}

type IFileLoadingUsecase interface {
	WriteImage(file multipart.File, header *multipart.FileHeader) (string, string, error)
	CVtoPDF(cvID uint64, currentUser *dto.UserFromSession) (*dto.CVPDFFile, error)
}
