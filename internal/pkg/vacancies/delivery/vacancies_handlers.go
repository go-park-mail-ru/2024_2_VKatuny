package delivery

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/vacancies"

	"github.com/sirupsen/logrus"
)

// VacanciesHandler returns list of vacancies
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
func VacanciesHandler(repo vacancies.Repository) http.Handler { //vacanciesTable *inmemorydb.VacanciesHandler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		funcName := "VacanciesHandler"
		fmt.Println("s122134")
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		logger.Debugf("yo")
		if !ok {
			fmt.Printf("%s: can't get logger from context\n", funcName)
		}

		if r.Method != http.MethodGet {
			logger.Errorf("Got %s method; allowed %s", r.Method, http.MethodGet)
			middleware.UniversalMarshal(
				w,
				http.StatusMethodNotAllowed,
				dto.JsonResponse{
					HttpStatus: http.StatusMethodNotAllowed,
					Error:      "http request method isn't a GET",
				},
			)
			return
		}

		queryParams := r.URL.Query()
		logger.Debugf("function: %s; Query params read: %v", funcName, queryParams)

		offsetStr := queryParams.Get("offset")
		if offsetStr == "" {
			logger.Errorf("Offset is empty; %s", offsetStr)
			middleware.UniversalMarshal(
				w,
				http.StatusBadRequest,
				dto.JsonResponse{
					HttpStatus: http.StatusBadRequest,
					Error:      "query parameter offset is empty",
				},
			)
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			logger.Errorf("function: %s; got err: %s; offset = %s", funcName, err, offsetStr)
			middleware.UniversalMarshal(
				w,
				http.StatusBadRequest,
				dto.JsonResponse{
					HttpStatus: http.StatusBadRequest,
					Error:      "query parameter offset isn't a number",
				},
			)
			return
		}

		numStr := queryParams.Get("num")
		if numStr == "" {
			logger.Errorf("function: %s; num is empty; %s", funcName, numStr)
			middleware.UniversalMarshal(
				w,
				http.StatusBadRequest,
				dto.JsonResponse{
					HttpStatus: http.StatusBadRequest,
					Error:      "query parameter num is empty",
				},
			)
			return
		}

		num, err := strconv.Atoi(numStr)
		if err != nil {
			logger.Errorf("function: %s; got err: %s; num = %s", funcName, err, numStr)
			middleware.UniversalMarshal(
				w,
				http.StatusBadRequest,
				dto.JsonResponse{
					HttpStatus: http.StatusBadRequest,
					Error:      "query parameter num isn't a number",
				},
			)
			return
		}

		// May cause a bug when the offset and num is negative
		// Because of type casting
		logger.Debugf("function: %s; going to db for vacancies", funcName)
		vacancies, err := repo.GetWithOffset(uint64(offset), uint64(offset+num))

		if err != nil {
			logger.Errorf("function: %s; got err while reading vacancies from db %s", funcName, err)
			middleware.UniversalMarshal(
				w,
				http.StatusBadRequest,
				dto.JsonResponse{
					HttpStatus: http.StatusBadRequest,
					Error:      "can't get vacancies from db",
				},
			)
		}

		logger.Debugf("function: %s; got vacancies %+v", funcName, vacancies)
		middleware.UniversalMarshal(
			w,
			http.StatusOK,
			dto.JsonResponse{
				HttpStatus: http.StatusOK,
				Body:       vacancies,
			},
		)
	})
}
