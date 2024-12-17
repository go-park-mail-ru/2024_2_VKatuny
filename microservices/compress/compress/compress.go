package compressmicroserviceinterface

import "fmt"

// Interface for Compress.
type ICompressRepository interface {
	//SaveFile(filename string, fileType string, file []byte) error
	DeleteFile(filename string) error
	ScanDir() error
}

type ICompressUsecase interface {
	//CompressAndSaveFile(filename string, fileType string, file []byte) error
	//DeleteFile(filename string) error
	ScanDir() error
}

var (
	AllowedTypes = []string{
		"image/jpeg",
		"image/png",
	}
)

var (
	NotAllowedType = fmt.Errorf("not allowed type")
	WrongData = fmt.Errorf("wrong data in input")
)
