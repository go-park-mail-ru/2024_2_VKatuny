package repository

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/sirupsen/logrus"
)

type FileLoadingStorage struct {
	logger      *logrus.Logger
	mediaDir    string
	cvinPDFdir  string
	templateDir string
}

func NewFileLoadingStorage(logger *logrus.Logger, mediaDir, CVinPDFDir, templateDir string) *FileLoadingStorage {
	return &FileLoadingStorage{
		logger:      logger,
		mediaDir:    mediaDir,
		cvinPDFdir:  CVinPDFDir,
		templateDir: templateDir,
	}
}

func (s *FileLoadingStorage) WriteFileOnDisk(filename string, header *multipart.FileHeader, file multipart.File) (string, string, error) {
	fn := "FileLoadingStorage.WriteFileOnDisk"
	s.logger.Debugf("%s: entering", fn)
	fmt.Println(s.mediaDir + filename + header.Filename)
	dst, err := os.Create(s.mediaDir + filename + header.Filename)
	if err != nil {
		s.logger.Errorf("%s: got err %s", fn, err)
		return "", "", fmt.Errorf("error creating file")
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		s.logger.Errorf("%s: got error copying file", fn)
		return "", "", fmt.Errorf("error copying file")
	}
	s.logger.Debugf("%s: done with name: %s and %s", fn, s.mediaDir, filename+header.Filename)
	return s.mediaDir, filename + header.Filename, nil
}

func (s *FileLoadingStorage) CVtoPDF(CV *dto.JSONCv, applicant *dto.JSONGetApplicantProfile) (string, error) {
	fn := "FileLoadingStorage.CVtoPDF"
	s.logger.Debugf("%s: entering", fn)
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		s.logger.Errorf("%s: got err %s", fn, err)
		return "", err
	}
	pwd, err := os.Getwd()
	pwd += "/"
	if err != nil {
		s.logger.Errorf("%s: got err %s", fn, err)
		return "", err
	}
	html, err := os.ReadFile("templates/template.html")
	if err != nil {
		s.logger.Errorf("%s: got err %s", fn, err)
		return "", err
	}
	htmlText := string(html)
	fmt.Println(pwd + "/" + s.mediaDir + CV.Avatar)

	htmlText = strings.Replace(htmlText, "template.css", pwd+s.templateDir+"template.css", 1)
	htmlText = strings.Replace(htmlText, "profile-avatar.jpg", pwd+s.mediaDir+CV.Avatar, 1)
	htmlText = strings.Replace(htmlText, "FirstName", applicant.FirstName, 1)
	htmlText = strings.Replace(htmlText, "LastName", applicant.LastName, 1)
	htmlText = strings.Replace(htmlText, "Contacts", applicant.Contacts, 1)
	htmlText = strings.Replace(htmlText, "PositionRu", CV.PositionRu, 1)
	htmlText = strings.Replace(htmlText, "PositionEn", CV.PositionEn, 1)
	htmlText = strings.Replace(htmlText, "PositionCategoryName", CV.PositionCategoryName, 1)
	htmlText = strings.Replace(htmlText, "WorkingExperience", CV.WorkingExperience, 1)
	htmlText = strings.Replace(htmlText, "Education", applicant.Education, 1)
	htmlText = strings.Replace(htmlText, "BirthDate", applicant.BirthDate[:9], 1)
	htmlText = strings.Replace(htmlText, "City", applicant.City, 1)
	htmlText = strings.Replace(htmlText, "CreatedAt", CV.CreatedAt[:9], 1)
	

	page := wkhtml.NewPageReader(strings.NewReader(htmlText))
	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)
	err = pdfg.Create()
	if err != nil {
		s.logger.Errorf("%s: got err %s", fn, err)
		return "", err
	}
	name := s.cvinPDFdir + strconv.Itoa(int(CV.ID)) + "&&" + strconv.Itoa(int(CV.ApplicantID)) + ".pdf"
	err = pdfg.WriteFile(name)
	if err != nil {
		s.logger.Errorf("%s: got err %s", fn, err)
		return "", err
	}
	s.logger.Debugf("%s: done with name: %s", fn, name)
	return name, nil
}
