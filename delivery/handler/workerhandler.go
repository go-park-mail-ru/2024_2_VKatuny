package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
	"github.com/go-park-mail-ru/2024_2_VKatuny/usecase/service"
)

// CreateWorker godoc
// @Summary     Creates a new user as a worker
// @Description -
// @Tags        Registration
// @Accept      json
// @Produce     json
// @Param       email    body string true "User's email"
// @Param       password body string true "User's password"
// @Success     200 {object} BD.UserInput
// @Failure     http.StatusBadRequest {object} nil
// @Router      /registration/worker/ [post]
func CreateWorkerHandler(h *BD.WorkerHandlers) http.Handler {
	return HttpHeadersWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)

		newUserInput := new(BD.WorkerInput)
		err := decoder.Decode(newUserInput)
		if err != nil {
			storage.UniversalMarshal(w, http.StatusBadRequest, nil)
			log.Printf("error while unmarshalling worker JSON: %s", err)
			return
		}
		if len(newUserInput.WorkerName) < 3 || len(newUserInput.WorkerLastName) < 3 ||
			strings.Index(newUserInput.WorkerEmail, "@") < 0 || len(newUserInput.WorkerPassword) < 4 {
			storage.UniversalMarshal(w, http.StatusBadRequest, nil)
			log.Printf("error while unmarshalling employer  JSON: %s", err)
			return
		}
		user, err := service.TryCreateWorker(h, newUserInput)
		if err == nil {
			storage.UniversalMarshal(w, http.StatusOK, user)
		} else {
			log.Println("!!!", err)
			storage.UniversalMarshal(w, http.StatusBadRequest, BD.UserAlreadyExist{true})
			log.Printf("error user with this email already exists: %s", newUserInput.WorkerEmail)
		}

	}))
}
