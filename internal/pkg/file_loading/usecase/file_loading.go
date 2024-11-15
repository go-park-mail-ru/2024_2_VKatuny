package usecase

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
)

func WriteFile(staticDir string, file multipart.File, header *multipart.FileHeader) (string, error) {
	filename := utils.GenerateSessionToken(utils.TokenLength+10, dto.UserTypeApplicant)
	dst, err := os.Create(staticDir + filename + header.Filename)
	if err != nil {
		log.Println("error creating file", err) // just do it
		return "", fmt.Errorf("error creating file")
	}
	defer dst.Close()
	var p []byte
	file.Read(p)
	file.Seek(0, 0)
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("error copying file")
	}
	log.Println("successfully uploaded file", err)
	return filename, nil
}
