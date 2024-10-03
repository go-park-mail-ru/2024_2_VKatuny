package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
	// "github.com/go-park-mail-ru/2024_2_VKatuny/storage"
)

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
func VacanciesHandler(vacanciesTable *BD.VacanciesHandler) http.Handler {
	return HttpHeadersWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != http.MethodGet {
			storage.UniversalMarshal(w, http.StatusMethodNotAllowed, BD.ErrorMessages{405, "http request method isn't a GET"})
			return
		}

		queryParams := r.URL.Query()

		offsetStr := queryParams.Get("offset")
		if offsetStr == "" {
			log.Println("status 400 offset is empty")
			storage.UniversalMarshal(w, http.StatusBadRequest, BD.ErrorMessages{405, "offset is empty"})
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			log.Println("status 400 offset isn't number")
			storage.UniversalMarshal(w, http.StatusBadRequest, BD.ErrorMessages{405, "offset isn't number"})
			return
		}

		numStr := queryParams.Get("num")
		if numStr == "" {
			log.Println("status 400 num is empty")
			storage.UniversalMarshal(w, http.StatusBadRequest, BD.ErrorMessages{405, "num is empty"})
			return
		}

		num, err := strconv.Atoi(numStr)
		if err != nil {
			log.Println("status 400 num isn't number")
			storage.UniversalMarshal(w, http.StatusBadRequest, BD.ErrorMessages{405, "num isn't number"})
			return
		}

		leftBound := offset
		rightBound := offset + num
		// covering cases when offset is out of slice bounds
		if leftBound > int(vacanciesTable.Count) {
			rightBound = leftBound
		} else if rightBound > int(vacanciesTable.Count) {
			rightBound = int(vacanciesTable.Count)
		}

		vacanciesTable.Mutex.RLock()
		vacancies := vacanciesTable.Vacancy[leftBound:rightBound]
		vacanciesTable.Mutex.RUnlock()

		response := map[string]interface{}{
			"status":    200,
			"vacancies": vacancies,
		}

		storage.UniversalMarshal(w, http.StatusOK, response)
	}))
}
