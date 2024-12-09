package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
)


// @Summary Search CVs
// @Description Get CVs with optional parameters
// @Tags CV
// @Accept json
// @Produce json
// @Param offset query integer false "Offset"
// @Param num query integer false "Number of objects to return"
// @Param searchQuery query string false "Search query"
// @Param group query string false "Group"
// @Param searchBy query string false "Search by"
// @Success 200 {object} dto.JSONResponse "CVs"
// @Failure 405 {object} dto.JSONResponse
// @Router /api/v1/cvs [get]
func (h *CVsHandler) SearchCVs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "CvsHandlers.GetCVs"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s; entering", fn)

	queryParams := r.URL.Query()
	h.logger.Debugf("%s; Query params read: %v", fn, queryParams)

	offsetStr := queryParams.Get("offset")
	numStr := queryParams.Get("num")
	searchStr := queryParams.Get("searchQuery")
	group := queryParams.Get("group")
	searchBy := queryParams.Get("searchBy")
	CVs, err := h.cvsUsecase.SearchCVs(offsetStr, numStr, searchStr, group, searchBy)
	for _, cv := range CVs {
		cv.CompressedAvatar  = h.fileLoadingUsecase.FindCompressedFile(cv.Avatar)
	}
	if err != nil {
		h.logger.Errorf("function: %s; got err while reading CVs from db %s", fn, err)
		middleware.UniversalMarshal(
			w,
			http.StatusInternalServerError,
			dto.JSONResponse{
				HTTPStatus: http.StatusInternalServerError,
				Error:      dto.MsgDataBaseError,
			},
		)
		return
	}

	h.logger.Debugf("%s; got CVs %v", fn, CVs)
	middleware.UniversalMarshal(
		w,
		http.StatusOK,
		dto.JSONResponse{
			HTTPStatus: http.StatusOK,
			Body:       CVs,
		},
	)
}
