package delivery

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/inmemorydb"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/vacancies"
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

		if r.Method != http.MethodGet {
			log.Println("Not GET method on /api/v1/vacancies StatusMethodNotAllowed")
			middleware.UniversalMarshal(w, http.StatusMethodNotAllowed, inmemorydb.ErrorMessages{http.StatusMethodNotAllowed, "http request method isn't a GET"})
			return
		}

		queryParams := r.URL.Query()

		offsetStr := queryParams.Get("offset")
		if offsetStr == "" {
			log.Println("status 400 offset is empty")
			middleware.UniversalMarshal(w, http.StatusBadRequest, inmemorydb.ErrorMessages{http.StatusBadRequest, "offset is empty"})
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			log.Println("status 400 offset isn't number")
			middleware.UniversalMarshal(w, http.StatusBadRequest, inmemorydb.ErrorMessages{http.StatusBadRequest, "offset isn't number"})
			return
		}

		numStr := queryParams.Get("num")
		if numStr == "" {
			log.Println("status 400 num is empty")
			middleware.UniversalMarshal(w, http.StatusBadRequest, inmemorydb.ErrorMessages{http.StatusBadRequest, "num is empty"})
			return
		}

		num, err := strconv.Atoi(numStr)
		if err != nil {
			log.Println("status 400 num isn't number")
			middleware.UniversalMarshal(w, http.StatusBadRequest, inmemorydb.ErrorMessages{http.StatusBadRequest, "num isn't number"})
			return
		}

		// May cause a bug when the offset and num is negative
		// Because of type casting
		vacancies, err := repo.GetWithOffset(uint64(offset), uint64(offset+num))

		if err != nil {
			log.Println("status 500 inmemorydb has problems")
			middleware.UniversalMarshal(w, http.StatusBadRequest, inmemorydb.ErrorMessages{http.StatusInternalServerError, "inmemorydb has problems"})
		}

		response := map[string]interface{}{
			"status":    http.StatusOK,
			"vacancies": vacancies,
		}

		middleware.UniversalMarshal(w, http.StatusOK, response)
	})
}