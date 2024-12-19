package repository

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/template"

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
	s.logger.Debugf("%s: entering with name: %s", fn, s.mediaDir+filename+header.Filename)
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
	dirList := strings.Split(s.mediaDir, "/")
	dirList = dirList[slices.Index(dirList, "2024_2_VKatuny")+1:]
	dirCut := strings.Join(dirList, "/") + "/"
	return dirCut, filename + header.Filename, nil
}

func (s *FileLoadingStorage) CVtoPDF(CV *dto.JSONCv, applicant *dto.JSONGetApplicantProfile) (string, error) {
	fn := "FileLoadingStorage.CVtoPDF"
	s.logger.Debugf("%s: entering", fn)

	templateDir := addslesh(s.templateDir)
	mediaDir := addslesh(s.mediaDir)
	cvinPDFdir := addslesh(s.cvinPDFdir)
	tmpl := template.Must(template.ParseFiles(templateDir + "template.html"))
	type And struct {
		CV        dto.JSONCv
		Applicant dto.JSONGetApplicantProfile
		IsImg     int
		Template  string
	}
	megaStruct := And{CV: *CV, Applicant: *applicant}
	if len(megaStruct.Applicant.BirthDate) > 10 {
		megaStruct.Applicant.BirthDate = megaStruct.Applicant.BirthDate[:10]
	}
	if len(megaStruct.CV.CreatedAt) > 10 {
		megaStruct.CV.CreatedAt = megaStruct.CV.CreatedAt[:10]
	}
	s.logger.Debugf("avatar: %s", mediaDir+CV.Avatar)
	s.logger.Debugf("template: %s", templateDir+"template.css")
	megaStruct.Template = templateDir + "template.css"
	megaStruct.CV.Avatar = mediaDir + strings.Split(CV.Avatar, "/")[len(strings.Split(CV.Avatar, "/"))-1]
	if CV.Avatar != "" {
		megaStruct.IsImg = 1
	} else {
		megaStruct.IsImg = 0
	}
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, megaStruct)
	if err != nil {
		s.logger.Errorf("%s: got err %s", fn, err)
		return "", err
	}

	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		s.logger.Errorf("%s: got err %s", fn, err)
		return "", err
	}
	page := wkhtml.NewPageReader(strings.NewReader(buf.String()))
	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)
	err = pdfg.Create()
	if err != nil {
		s.logger.Errorf("%s: got err %s", fn, err)
		return "", err
	}
	name := cvinPDFdir + strconv.Itoa(int(CV.ID)) + "&&" + strconv.Itoa(int(CV.ApplicantID)) + ".pdf"
	err = pdfg.WriteFile(name)
	if err != nil {
		s.logger.Errorf("%s: got err %s", fn, err)
		return "", err
	}
	s.logger.Debugf("%s: done with name: %s", fn, name)
	return name, nil
}

func addslesh(str string) string {
	if str[len(str)-1] != '/' {
		str += "/"
	}
	return str
}
