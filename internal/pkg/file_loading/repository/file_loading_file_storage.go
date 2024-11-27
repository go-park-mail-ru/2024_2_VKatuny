package repository

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
)

type FileLoadingStorage struct {
	dir string
}

func NewFileLoadingStorage(dir string) *FileLoadingStorage {
	return &FileLoadingStorage{
		dir: dir,
	}
}

func (s *FileLoadingStorage) WriteFileOnDisk(filename string, header *multipart.FileHeader, file multipart.File) (string, string, error) {
	fmt.Println(s.dir + filename + header.Filename)
	dst, err := os.Create(s.dir + filename + header.Filename)
	if err != nil {
		log.Println("error creating file", err)
		return "", "", fmt.Errorf("error creating file")
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return "", "", fmt.Errorf("error copying file")
	}
	return s.dir, filename + header.Filename, nil
}
