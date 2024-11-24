package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
)

// GetCVsHandler returns list of CVs
// GetCVs godoc
// @Summary     Gets list of CVs
// @Description Accepts offset and number of CVs with id >= offset. Returns CVs
// @Tags        CVs
// @Produce     json
// @Param       offset query int true "offset"
// @Param       num    query int true "num"
// @Success     200
// @Failure     400
// @Failure     405
// @Failure     500
// @Router      /CVs [get]
func (h *CVsHandler) SearchCVs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "CvsHandlers.GetCVs"
	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	queryParams := r.URL.Query()
	h.logger.Debugf("%s; Query params read: %v", fn, queryParams)

	offsetStr := queryParams.Get("offset")
	numStr := queryParams.Get("num")
	searchStr := queryParams.Get("searchQuery")
	group := queryParams.Get("group")
	searchBy := queryParams.Get("searchBy")
	CVs, err := h.cvsUsecase.SearchCVs(offsetStr, numStr, searchStr, group, searchBy)

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
