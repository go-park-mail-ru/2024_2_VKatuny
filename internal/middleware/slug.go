package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

var (
	ErrorSlug = fmt.Errorf("something bad with slug")
)

func GetIDSlugAtEnd(w http.ResponseWriter, r *http.Request, prefix string) (int, error) {
	url := r.URL.Path[len(prefix):]
	slug := strings.Split(url, "/")
	ID, err := strconv.Atoi(slug[0])
	if err != nil || len(slug) > 1 {
		// h.logger.Errorf("function %s: got err %s", fn, ErrorSlug)
		UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      ErrorSlug.Error(),
		})
		return 0, ErrorSlug
	}
	return ID, nil

}
