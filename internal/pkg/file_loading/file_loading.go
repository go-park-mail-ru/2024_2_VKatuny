package fileloading

import (
	"mime/multipart"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

type IFileLoadingRepository interface {
	WriteFileOnDisk(filename string, header *multipart.FileHeader, file []byte) (string, string, error)
	CVtoPDF(CV *dto.JSONCv, applicant *dto.JSONGetApplicantProfile) (string, error)
}

type IFileLoadingUsecase interface {
	WriteImage(file []byte, header *multipart.FileHeader) (string, string, error)
	FindCompressedFile(filename string) string
	CVtoPDF(CV *dto.JSONCv, applicant *dto.JSONGetApplicantProfile) (*dto.CVPDFFile, error)
}
