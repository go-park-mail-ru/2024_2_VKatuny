package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

const (
	defaultVacanciesOffset = 0
	defaultVacanciesNum    = 10
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

	queryParams := r.URL.Query()
	h.logger.Debugf("%s; Query params read: %v", fn, queryParams)

	offsetStr := queryParams.Get("offset")
	numStr := queryParams.Get("num")
	offset, num, err := h.vacanciesUsecase.ValidateQueryParameters(offsetStr, numStr)

	// at lest one param is incorrect or not presented
	// using default values
	if err != nil {
		h.logger.Debugf("%s; got err: %s - processing with default values", fn, err)
		offset = defaultVacanciesOffset
		num = defaultVacanciesNum
	}
	h.logger.Debugf("%s; got offset %d and num %d", fn, offset, num)

	vacancies, err := h.vacanciesUsecase.GetVacanciesWithOffset(offset, num)
	if err != nil {
		h.logger.Errorf("%s; got err while reading vacancies from db %s", fn, err)
		middleware.UniversalMarshal(
			w,
			http.StatusInternalServerError,
			dto.JSONResponse{
				HTTPStatus: http.StatusInternalServerError,
				Error:      err.Error(),
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
