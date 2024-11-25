package fileloading

import (
	"mime/multipart"
)

type IFileLoadingRepository interface {
	WriteFileOnDisk(filename string, header *multipart.FileHeader, file multipart.File) error
}

type IFileLoadingUsecase interface {
	WriteImage(file multipart.File, header *multipart.FileHeader) (string, error)
}
