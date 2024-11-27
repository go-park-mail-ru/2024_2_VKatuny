package compressmicroserviceinterface

import "fmt"

// Interface for Compress.
type ICompressRepository interface {
	SaveFile(filename string, fileType string, file []byte) error
	DeleteFile(filename string) error
}

type ICompressUsecase interface {
	CompressAndSaveFile(filename string, fileType string, file []byte) error
	DeleteFile(filename string) error
}

var (
	AllowedTypes = []string{
		"image/jpeg",
		"image/png",
	}
)

var (
	NotAllowedType = fmt.Errorf("not allowed type")
)
