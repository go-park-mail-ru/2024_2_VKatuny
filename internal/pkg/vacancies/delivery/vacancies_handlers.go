package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
)

// GetVacancies godoc
// @Summary Get a list of vacancies
// @Description Retrieves a list of vacancies based on the provided query parameters
// @Tags Vacancy
// @Produce json
// @Param offset query string false "Offset for pagination"
// @Param num query string false "Number of vacancies to retrieve"
// @Param searchQuery query string false "Search query string"
// @Param group query string false "Group category for filtering"
// @Param searchBy query string false "Field to search by"
// @Success 200 {object} dto.JSONResponse{body=[]dto.JSONVacancy}
// @Failure 500 {object} dto.JSONResponse{error=string} "Internal Server Error"
// @Failure 405 {object} dto.JSONResponse{error=string} "Method Not Allowed"
// @Router /api/v1/vacancies [get]
func (h *VacanciesHandlers) GetVacancies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "VacanciesHandlers.GetVacancies"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s; entering", fn)

	queryParams := r.URL.Query()
	h.logger.Debugf("%s; Query params read: %v", fn, queryParams)

	offsetStr := queryParams.Get("offset")
	numStr := queryParams.Get("num")
	searchStr := queryParams.Get("searchQuery")
	group := queryParams.Get("group")
	searchBy := queryParams.Get("searchBy")
	vacancies, err := h.vacanciesUsecase.SearchVacancies(r.Context(),offsetStr, numStr, searchStr, group, searchBy)
	for _, vacancy := range vacancies {
		vacancy.CompressedAvatar  = h.fileLoadingUsecase.FindCompressedFile(vacancy.Avatar)
	}
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

