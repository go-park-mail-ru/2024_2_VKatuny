package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

func WriteFile(staticDir string, file multipart.File, header *multipart.FileHeader) (string, error) {
	a := header.Header
	fmt.Println(a["Content-Type"][0])
	for _, i := range a["Content-Type"] {
		switch i {
		case "image/jpeg":
		case "image/jpg":
		case "image/png":
		case "image/svg":
		case "image/svg+xml":
		default:
			return "", fmt.Errorf(dto.MsgInvalidFile)
		}
	}
	filename := GenerateSessionToken(TokenLength+10, dto.UserTypeApplicant)
	dst, err := os.Create(staticDir + filename + header.Filename)
	if err != nil {
		log.Println("error creating file", err) // just do it
		return "", fmt.Errorf("error creating file")
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("error copying file")
	}
	log.Println("successfully uploaded file", err)
	return filename, nil
}
