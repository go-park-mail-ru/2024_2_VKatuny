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
	filename := utils.GenerateSessionToken(utils.TokenLength+10, dto.UserTypeApplicant)
	dst, err := os.Create(staticDir + filename + header.Filename)
	if err != nil {
		log.Println("error creating file", err) // just do it
		return "", fmt.Errorf("error creating file")
	}
	defer dst.Close()
	// p := make([]byte, 10000)
	// a, b := file.Read(p)
	// fmt.Println("!!", a, b, string(p))

	// mediatype, params, err := mime.ParseMediaType(string(p))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("type:", mediatype)
	// fmt.Println("charset:", params["charset"])

	// file.Seek(0, 0)
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("error copying file")
	}
	log.Println("successfully uploaded file", err)
	return filename, nil
}
