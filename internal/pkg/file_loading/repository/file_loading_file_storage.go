package repository

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
)

type FileLoadingStorage struct {
	mediaDir string
	cvinPDFdir string
}

func NewFileLoadingStorage(mediaDir, CVinPDFDir string) *FileLoadingStorage {
	return &FileLoadingStorage{
		mediaDir: mediaDir,
		cvinPDFdir: CVinPDFDir,
	}
}

func (s *FileLoadingStorage) WriteFileOnDisk(filename string, header *multipart.FileHeader, file multipart.File) (string, string, error) {
	fmt.Println(s.mediaDir + filename + header.Filename)
	dst, err := os.Create(s.mediaDir + filename + header.Filename)
	if err != nil {
		log.Println("error creating file", err)
		return "", "", fmt.Errorf("error creating file")
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return "", "", fmt.Errorf("error copying file")
	}
	return s.mediaDir, filename + header.Filename, nil
}

func (s *FileLoadingStorage) CVtoPDF(cvID uint64) (string, error) {
	return s.cvinPDFdir, nil
}
