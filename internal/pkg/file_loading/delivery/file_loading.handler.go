package delivery

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	file_loading_usecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading/usecase"
)

var uploadFormTmpl = []byte(`
<html>
	<body>
	<form action="api/v1/upload" method="post" enctype="multipart/form-data">
		Image: <input type="file" name="my_file">
		<input type="submit" value="Upload">
	</form>
	</body>
</html>
`)

func CreateMainP() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(uploadFormTmpl)
	})
}

func CreateUploadHandler(staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(2.5 * (10 << 20)) // 25Mb
		file, header, err := r.FormFile("my_file")
		defer file.Close()
		defer r.MultipartForm.RemoveAll()
		if err != nil {
			log.Println("error retrieving file", err)
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      dto.MsgUnableToReadFile,
			})
			return
		}
		fileAddress, err := file_loading_usecase.WriteFile(staticDir, file, header)
		if err != nil {
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      dto.MsgUnableToUploadFile,
			})
		}
		middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
			HTTPStatus: http.StatusOK,
			Body:       fileAddress + header.Filename,
		})

	})
}
