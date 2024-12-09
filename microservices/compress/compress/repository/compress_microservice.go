package repository

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"github.com/sirupsen/logrus"
)

type CompressRepository struct {
	compressedDir   string
	uncompressedDir string
	logger          *logrus.Entry
}

func NewCompressRepository(compressedDir, uncompressedDir string, logger *logrus.Logger) *CompressRepository {
	return &CompressRepository{
		compressedDir:   compressedDir,
		uncompressedDir: uncompressedDir,
		logger:          &logrus.Entry{Logger: logger},
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

func (cr *CompressRepository) DeleteFile(filePath string) error {
	funcName := "CompressRepository.DeleteFile"
	cr.logger.Debugf("%s: got request: %s", funcName, filePath)
	err := os.Remove(filePath)
	if err != nil {
		cr.logger.Errorf("%s: got err: %s", funcName, err.Error())
		return err
	}
	cr.logger.Errorf("%s: file %s deleted", funcName, filePath)
	return nil
}

func (cr *CompressRepository) ScanDir() error {
	funcName := "CompressRepository.ScanDir"
	cr.logger.Debugf("%s: working", funcName)
	compressed, err := os.ReadDir(cr.compressedDir)
	if err != nil {
		cr.logger.Errorf("%s: got err: %s", funcName, err.Error())
	}
	uncompressed, err := os.ReadDir(cr.uncompressedDir)
	if err != nil {
		cr.logger.Errorf("%s: got err: %s", funcName, err.Error())
	}
	compressedList := []string{}
	uncompressedList := []string{}
	for _, file := range compressed {
		if !file.IsDir() {
			compressedList = append(compressedList, file.Name())
		}
	}
	for _, file := range uncompressed {
		if !file.IsDir() {
			uncompressedList = append(uncompressedList, file.Name())
		}
	}
	slices.Sort(compressedList)
	slices.Sort(uncompressedList)
	if slices.Equal(compressedList, uncompressedList) {
		cr.logger.Debugf("%s: Nothing to compress or delete", funcName)
		return nil
	} else {
		newCompressedList := []string{}
		for _, file := range compressedList {
			if !slices.Contains(uncompressedList, file) {
				err := cr.DeleteFile(cr.compressedDir + file)
				if err != nil {
					cr.logger.Errorf("%s: got err: %s", funcName, err.Error())
				}
			} else {
				newCompressedList = append(newCompressedList, file)
			}
		}
		for _, file := range uncompressedList {
			if !slices.Contains(newCompressedList, file) {
				err := cr.CompressAndWriteFile(cr.uncompressedDir + file)
				if err != nil {
					cr.logger.Errorf("%s: got err: %s", funcName, err.Error())
				}
			}
		}
	}
	return nil
}

func (cr *CompressRepository) CompressAndWriteFile(filePath string) error {
	funcName := "CompressRepository.CompressAndWriteFile"
	cr.logger.Debugf("%s: got request: %s", funcName, filePath)

	cr.logger.Errorf("%s: file %s compressed", funcName, filePath)
	return nil
}
