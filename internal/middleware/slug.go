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
	slugID := strings.Split(url, "/")[0]
	ID, err := strconv.Atoi(slugID)
	if len(url) > 1 || err != nil {
		// h.logger.Errorf("function %s: got err %s", fn, ErrorSlug)
		UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      ErrorSlug.Error(),
		})
		return 0, ErrorSlug
	}
	return ID, nil

}
