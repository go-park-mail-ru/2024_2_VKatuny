package delivery

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies"

	vacanciesUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/usecase"

	"github.com/sirupsen/logrus"
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
func GetVacanciesHandler(repo vacancies.IVacanciesRepository) http.Handler { //vacanciesTable *inmemorydb.VacanciesHandler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		funcName := "vacancies.GetVacanciesHandler"
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("%s: can't get logger from context\n", funcName)
			return
		}

		queryParams := r.URL.Query()
		logger.Debugf("function: %s; Query params read: %v", funcName, queryParams)

		offsetStr := queryParams.Get("offset")
		numStr := queryParams.Get("num")
		offset, num, err := vacanciesUsecase.ValidateRequestParams(offsetStr, numStr)

		// at lest one param is incorrect or not presented
		// using default values
		if err != nil {
			logger.Debugf("function: %s; got err: %s - processing with default values", funcName, err)
			offset = defaultVacanciesOffset
			num = defaultVacanciesNum
		}

		logger.Debugf("function: %s; going to db for vacancies", funcName)
		// TODO: use usecase
		vacancies, err := repo.GetWithOffset(offset, offset+num)

		if err != nil {
			logger.Errorf("function: %s; got err while reading vacancies from db %s", funcName, err)
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

		logger.Debugf("function: %s; got vacancies %+v", funcName, vacancies)
		middleware.UniversalMarshal(
			w,
			http.StatusOK,
			dto.JSONResponse{
				HTTPStatus: http.StatusOK,
				Body:       vacancies,
			},
		)
	})
}
