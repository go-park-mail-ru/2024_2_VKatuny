package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
)

// GetVacanciesHandler returns list of vacancies
// GetVacancies godoc
// @Summary     Gets list of vacancies
// @Description Accepts offset and number of vacancies with id >= offset. Returns vacancies
// @Tags        Vacancies
// @Produce     json
// @Param       offset query int true "offset"
// @Param       num    query int true "num"
// @Success     200
// @Failure     400
// @Failure     405
// @Failure     500
// @Router      /vacancies [get]
func (h *VacanciesHandlers) GetVacancies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "VacanciesHandlers.GetVacancies"
	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	queryParams := r.URL.Query()
	h.logger.Debugf("%s; Query params read: %v", fn, queryParams)

	offsetStr := queryParams.Get("offset")
	numStr := queryParams.Get("num")
	searchStr := queryParams.Get("positionDescription")
	vacancies, err := h.vacanciesUsecase.SearchVacancies(offsetStr, numStr, searchStr)

	if err != nil {
		h.logger.Errorf("function: %s; got err while reading vacancies from db %s", fn, err)
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

	h.logger.Debugf("%s; got vacancies %v", fn, vacancies)
	middleware.UniversalMarshal(
		w,
		http.StatusOK,
		dto.JSONResponse{
			HTTPStatus: http.StatusOK,
			Body:       vacancies,
		},
	)
}
