package compressmicroserviceinterface

// Interface for Compress.
type ICompressRepository interface {
	SaveFile(filename string) error
	DeleteFile(filename string) error
}

type ICompressUsecase interface {
	CompressAndSaveFile(filename string) error
	DeleteFile(filename string) error
}
