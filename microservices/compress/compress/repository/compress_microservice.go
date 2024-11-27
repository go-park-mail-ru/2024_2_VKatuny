package repository

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

type CompressRepository struct {
	compressedDir string
	logger        *logrus.Entry
}

func NewCompressRepository(compressedDir string, logger *logrus.Logger) *CompressRepository {
	return &CompressRepository{
		compressedDir: compressedDir,
		logger:        &logrus.Entry{Logger: logger},
	}
}

func (cr *CompressRepository) SaveFile(filename string, fileType string, file []byte) error {
	funcName := "CompressRepository.SaveFile"
	cr.logger.Debugf("%s: got request: %s %s %s", funcName, filename, fileType, file)
	dst, err := os.Create(cr.compressedDir + filename)
	if err != nil {
		log.Println("error creating file", err)
		return fmt.Errorf("error creating file")
	}
	defer dst.Close()
	r := bytes.NewReader(file)
	if _, err := io.Copy(dst, r); err != nil {
		return fmt.Errorf("error copying file")
	}
	return nil
}

func (cr *CompressRepository) DeleteFile(filename string) error {
	funcName := "CompressRepository.DeleteFile"
	cr.logger.Debugf("%s: got request: %s", funcName, filename)
	return nil
}
