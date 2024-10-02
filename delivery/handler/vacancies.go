package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
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
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"status": 405, "error": "http request method isn't a GET"}`))
			return
		}

		queryParams := r.URL.Query()

		offsetStr := queryParams.Get("offset")
		if offsetStr == "" {
			log.Println("status 400 offset is empty")
			w.WriteHeader(http.StatusBadRequest) // HTTP 400
			w.Write([]byte(`{"status": 400, "error": "offset is empty"}`))
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			log.Println("status 400 offset isn't number")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"status": 400, "error": "offset isn't number"}`))
			return
		}

		numStr := queryParams.Get("num")
		if numStr == "" {
			log.Println("status 400 num is empty")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"status": 400, "error": "num is empty"}`))
			return
		}

		num, err := strconv.Atoi(numStr)
		if err != nil {
			log.Println("status 400 num isn't number")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"status": 400, "error": "num isn't number"}`))
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

		vacanciesTable.Mutex.Lock()
		vacancies := vacanciesTable.Vacancy[leftBound:rightBound]
		vacanciesTable.Mutex.Unlock()

		response := map[string]interface{}{
			"status":    200,
			"vacancies": vacancies,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("status 500")
			w.Write([]byte(`{"status": 500, "error": "encoding error"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
	return HttpHeadersWrapper(http.HandlerFunc(fn))
}
